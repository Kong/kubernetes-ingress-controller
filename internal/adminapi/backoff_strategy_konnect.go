package adminapi

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
)

const (
	KonnectBackoffInitialInterval = time.Second * 3
	KonnectBackoffMaxInterval     = time.Minute * 15
	KonnectBackoffMultiplier      = 2
)

type Clock interface {
	Now() time.Time
}

// KonnectBackoffStrategy keeps track of Konnect config push backoffs.
//
// It takes into account:
// - a regular exponential backoff that is incremented on every Update failure,
// - a last failed configuration hash (where we skip Update until a config changes).
//
// It's important to note that KonnectBackoffStrategy can use the latter (config hash)
// because of the nature of the one-directional integration where KIC is the only
// component responsible for populating configuration of Konnect's Control Plane.
// In case that changes in the future (e.g. manual modifications to parts of the
// configuration are allowed on Konnect side for some reason), we might have to
// drop this part of the backoff strategy.
type KonnectBackoffStrategy struct {
	b                    *backoff.Backoff
	nextAttempt          time.Time
	clock                Clock
	lastFailedConfigHash []byte

	lock sync.RWMutex
}

func NewKonnectBackoffStrategy(clock Clock) *KonnectBackoffStrategy {
	exponentialBackoff := &backoff.Backoff{
		Min:    KonnectBackoffInitialInterval,
		Max:    KonnectBackoffMaxInterval,
		Factor: KonnectBackoffMultiplier,
	}
	exponentialBackoff.Reset()

	return &KonnectBackoffStrategy{
		b:     exponentialBackoff,
		clock: clock,
	}
}

func (s *KonnectBackoffStrategy) CanUpdate(configHash []byte) (bool, string) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// The exponential backoff duration is satisfied.
	// In case of the first attempt it will be satisfied as s.nextAttempt will be a zero value which is always in the past.
	timeLeft := s.nextAttempt.Sub(s.clock.Now())
	exponentialBackoffSatisfied := timeLeft <= 0

	// The configuration we're attempting to update is not the same faulty config we've already tried pushing.
	isTheSameFaultyConfig := s.lastFailedConfigHash != nil && bytes.Equal(s.lastFailedConfigHash, configHash)

	// In case both conditions are satisfied, we're good to make an attempt.
	if exponentialBackoffSatisfied && !isTheSameFaultyConfig {
		return true, ""
	}

	// Otherwise, we build a human-readable explanation of why the update cannot be performed at this point in time.
	return false, s.whyCannotUpdate(timeLeft, isTheSameFaultyConfig)
}

func (s *KonnectBackoffStrategy) RegisterUpdateFailure(err error, configHash []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()

	apiErrs := deckerrors.ExtractAPIErrors(err)
	tooManyRequestsErr, isTooManyRequests := lo.Find(apiErrs, func(err *kong.APIError) bool {
		return err.Code() == http.StatusTooManyRequests
	})
	if isTooManyRequests {
		s.handleTooManyRequests(tooManyRequestsErr)
		return
	}

	isClientError := lo.ContainsBy(apiErrs, func(err *kong.APIError) bool {
		return err.Code() >= 400 && err.Code() < 500
	})
	if isClientError {
		s.handleGenericClientError(configHash)
		return
	}

	// If it's neither of the specific cases above, we just increment the standard exponential backoff.
	s.incrementExponentialBackoff()
}

func (s *KonnectBackoffStrategy) RegisterUpdateSuccess() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.b.Reset()
	s.nextAttempt = time.Time{}
	s.lastFailedConfigHash = nil
}

func (s *KonnectBackoffStrategy) handleTooManyRequests(tooManyRequestsErr *kong.APIError) {
	if details, ok := tooManyRequestsErr.Details().(kong.ErrTooManyRequestsDetails); ok && details.RetryAfter != 0 {
		// In case we get 429 with details embedded, we just retry after the suggested Retry-After time.
		s.nextAttempt = s.clock.Now().Add(details.RetryAfter)
	} else {
		// In case the details for 429 are missing, we retry after the standard exponential backoff time.
		s.incrementExponentialBackoff()
	}

	// Despite whether we've got details or not, we prune the last failed config hash to not block update after the
	// period we set up above.
	s.lastFailedConfigHash = nil
}

func (s *KonnectBackoffStrategy) handleGenericClientError(configHash []byte) {
	// We increment the standard exponential backoff time and store the faulty config hash to prevent pushing it again.
	s.incrementExponentialBackoff()
	s.lastFailedConfigHash = configHash
}

func (s *KonnectBackoffStrategy) incrementExponentialBackoff() {
	// Backoff.Duration() call returns backoff time we need to wait until next attempt.
	// It also increments the internal attempts counter so the next time we call it, the
	// duration will be multiplied accordingly.
	timeLeft := s.b.Duration()

	// We're storing the exact point in time after which we'll be allowed to perform the next update attempt.
	s.nextAttempt = s.clock.Now().Add(timeLeft)
}

func (s *KonnectBackoffStrategy) whyCannotUpdate(
	timeLeft time.Duration,
	isTheSameFaultyConfig bool,
) string {
	var reasons []string

	if isTheSameFaultyConfig {
		reasons = append(reasons, fmt.Sprintf(
			"Config has to be changed: %q hash has already failed to be pushed with a client error",
			hex.EncodeToString(s.lastFailedConfigHash),
		))
	}

	if timeLeft > 0 {
		reasons = append(reasons, fmt.Sprintf("next attempt allowed in %s", timeLeft))
	}

	return strings.Join(reasons, ", ")
}
