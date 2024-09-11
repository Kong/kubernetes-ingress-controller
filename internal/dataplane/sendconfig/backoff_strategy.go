package sendconfig

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

type UpdateSkippedDueToBackoffStrategyError struct {
	explanation string
}

func NewUpdateSkippedDueToBackoffStrategyError(explanation string) UpdateSkippedDueToBackoffStrategyError {
	return UpdateSkippedDueToBackoffStrategyError{explanation: explanation}
}

func (e UpdateSkippedDueToBackoffStrategyError) Error() string {
	return fmt.Sprintf("update skipped due to a backoff strategy not being satisfied: %s", e.explanation)
}

// UpdateStrategyWithBackoff decorates any UpdateStrategy to respect a passed adminapi.UpdateBackoffStrategy.
type UpdateStrategyWithBackoff struct {
	decorated       UpdateStrategy
	backoffStrategy adminapi.UpdateBackoffStrategy
	logger          logr.Logger
}

func NewUpdateStrategyWithBackoff(
	decorated UpdateStrategy,
	backoffStrategy adminapi.UpdateBackoffStrategy,
	logger logr.Logger,
) UpdateStrategyWithBackoff {
	return UpdateStrategyWithBackoff{
		decorated:       decorated,
		backoffStrategy: backoffStrategy,
		logger:          logger,
	}
}

// Update will ensure that the decorated UpdateStrategy.Update is called only when an underlying
// UpdateBackoffStrategy.CanUpdate is satisfied.
// In case it's not, it will return a predefined ErrUpdateSkippedDueToBackoffStrategy.
// In case it is, apart from calling UpdateStrategy.Update, it will also register a success or a failure of an update
// attempt so that the UpdateBackoffStrategy can keep track of it.
func (s UpdateStrategyWithBackoff) Update(ctx context.Context, targetContent ContentWithHash) (err error) {
	if canUpdate, whyNot := s.backoffStrategy.CanUpdate(targetContent.Hash); !canUpdate {
		return NewUpdateSkippedDueToBackoffStrategyError(whyNot)
	}

	err = s.decorated.Update(ctx, targetContent)
	if err != nil {
		s.logger.V(logging.DebugLevel).Info("Update failed, registering it for backoff strategy", "reason", err.Error())
		s.backoffStrategy.RegisterUpdateFailure(err, targetContent.Hash)
		return err
	}

	s.backoffStrategy.RegisterUpdateSuccess()
	return nil
}

func (s UpdateStrategyWithBackoff) MetricsProtocol() metrics.Protocol {
	return s.decorated.MetricsProtocol()
}

func (s UpdateStrategyWithBackoff) Type() string {
	return fmt.Sprintf("WithBackoff(%s)", s.decorated.Type())
}
