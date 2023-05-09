package sendconfig

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

type ErrUpdateSkippedDueToBackoffStrategy struct {
	explanation string
}

func NewErrUpdateSkippedDueToBackoffStrategy(explanation string) ErrUpdateSkippedDueToBackoffStrategy {
	return ErrUpdateSkippedDueToBackoffStrategy{explanation: explanation}
}

func (e ErrUpdateSkippedDueToBackoffStrategy) Error() string {
	return fmt.Sprintf("update skipped due to a backoff strategy not being satisfied: %s", e.explanation)
}

func (e ErrUpdateSkippedDueToBackoffStrategy) Is(err error) bool {
	return errors.Is(err, ErrUpdateSkippedDueToBackoffStrategy{})
}

// UpdateBackoffStrategy keeps state of a backoff strategy.
type UpdateBackoffStrategy interface {
	// CanUpdate tells whether we're allowed to make an update attempt for a given config hash.
	// In case it returns false, the second return value is a human-readable explanation of why the update cannot
	// be updated at this point in time.
	CanUpdate([]byte) (bool, string)

	// RegisterUpdateSuccess resets the backoff strategy, effectively making it allow next update straight away.
	RegisterUpdateSuccess()

	// RegisterUpdateFailure registers an update failure along with its failure reason passed as a generic error, and
	// a config hash that we failed to push.
	RegisterUpdateFailure(failureReason error, configHash []byte)
}

// UpdateStrategyWithBackoff decorates any UpdateStrategy to respect a passed UpdateBackoffStrategy.
type UpdateStrategyWithBackoff struct {
	decorated       UpdateStrategy
	backoffStrategy UpdateBackoffStrategy
	log             logrus.FieldLogger
}

func NewUpdateStrategyWithBackoff(
	decorated UpdateStrategy,
	backoffStrategy UpdateBackoffStrategy,
	log logrus.FieldLogger,
) UpdateStrategyWithBackoff {
	return UpdateStrategyWithBackoff{
		decorated:       decorated,
		backoffStrategy: backoffStrategy,
		log:             log,
	}
}

// Update will ensure that the decorated UpdateStrategy.Update is called only when an underlying
// UpdateBackoffStrategy.CanUpdate is satisfied.
// In case it's not, it will return a predefined ErrUpdateSkippedDueToBackoffStrategy.
// In case it's, apart from calling UpdateStrategy.Update, it will also register a success or a failure of an update
// attempt so that the UpdateBackoffStrategy can keep track of it.
func (s UpdateStrategyWithBackoff) Update(ctx context.Context, targetContent ContentWithHash) (
	err error,
	resourceErrors []ResourceError,
	resourceErrorsParseErr error,
) {
	if canUpdate, whyNot := s.backoffStrategy.CanUpdate(targetContent.Hash); !canUpdate {
		s.log.Debug("Skipping update due to a backoff strategy")
		return NewErrUpdateSkippedDueToBackoffStrategy(whyNot), nil, nil
	}

	err, resourceErrors, resourceErrorsParseErr = s.decorated.Update(ctx, targetContent)
	if err != nil {
		s.log.WithError(err).Error("Update failed, registering it for backoff strategy")
		s.backoffStrategy.RegisterUpdateFailure(err, targetContent.Hash)
	} else {
		s.backoffStrategy.RegisterUpdateSuccess()
	}

	return err, resourceErrors, resourceErrorsParseErr
}

func (s UpdateStrategyWithBackoff) MetricsProtocol() metrics.Protocol {
	return s.decorated.MetricsProtocol()
}
