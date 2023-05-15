package adminapi

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckerrors"
)

const (
	KonnectBackoffInitialInterval = time.Second * 3
	KonnectBackoffMaxInterval     = time.Minute * 15
	KonnectBackoffMultiplier      = 2
)

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now() }

// KonnectBackoffStrategy keeps track of Konnect config push backoffs.
//
// It takes into account:
// - a regular exponential backoff that is incremented on every Update failure,
// - a last failed configuration hash (where we skip Update until a config changes).
//
// It's important to note that KonnectBackoffStrategy can use the latter (config hash)
// because of the nature of the one-directional integration where KIC is the only
// component responsible for populating configuration of Konnect's Runtime Group.
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

	if errs := deckerrors.ExtractAPIErrors(err); len(errs) > 0 {
		_, hasClientError := lo.Find(errs, func(item *kong.APIError) bool {
			return item.Code() >= 400 && item.Code() < 500
		})

		// We store the failed configuration hash only in case we receive a client error code [400, 500).
		// It's because we don't want to repeatedly try sending the config that we know is faulty on our side.
		// It only makes sense to retry when the config changes.
		if hasClientError {
			s.lastFailedConfigHash = configHash
		} else {
			s.lastFailedConfigHash = nil
		}
	}

	// Backoff.Duration() call returns backoff time we need to wait until next attempt.
	// It also increments the internal attempts counter so the next time we call it, the
	// duration will be multiplied accordingly.
	timeLeft := s.b.Duration()

	// We're storing the exact point in time after which we'll be allowed to perform the next update attempt.
	s.nextAttempt = s.clock.Now().Add(timeLeft)
}

func (s *KonnectBackoffStrategy) RegisterUpdateSuccess() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.b.Reset()
	s.nextAttempt = time.Time{}
	s.lastFailedConfigHash = nil
}

func (s *KonnectBackoffStrategy) whyCannotUpdate(
	timeLeft time.Duration,
	isTheSameFaultyConfig bool,
) string {
	var reasons []string

	if isTheSameFaultyConfig {
		reasons = append(reasons, fmt.Sprintf(
			"config has to be changed: %q hash has already failed to be pushed with a client error",
			string(s.lastFailedConfigHash),
		))
	}

	if timeLeft > 0 {
		reasons = append(reasons, fmt.Sprintf("next attempt allowed in %s", timeLeft))
	}

	return strings.Join(reasons, ", ")
}
