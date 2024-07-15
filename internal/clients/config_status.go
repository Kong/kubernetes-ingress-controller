package clients

import (
	"context"
	"time"

	"github.com/go-logr/logr"
)

// ConfigStatus is an enumerated type that represents the status of the configuration synchronisation.
// Look at CalculateConfigStatus for more details.
type ConfigStatus string

const (
	ConfigStatusOK                                         ConfigStatus = "OK"
	ConfigStatusTranslationErrorHappened                   ConfigStatus = "TranslationErrorHappened"
	ConfigStatusApplyFailed                                ConfigStatus = "ApplyFailed"
	ConfigStatusOKKonnectApplyFailed                       ConfigStatus = "OKKonnectApplyFailed"
	ConfigStatusTranslationErrorHappenedKonnectApplyFailed ConfigStatus = "TranslationErrorHappenedKonnectApplyFailed"
	ConfigStatusApplyFailedKonnectApplyFailed              ConfigStatus = "ApplyFailedKonnectApplyFailed"
	ConfigStatusUnknown                                    ConfigStatus = "Unknown"
)

// CalculateConfigStatusInput aggregates the input to CalculateConfigStatus.
type CalculateConfigStatusInput struct {
	// Any error occurred when syncing with Gateways.
	GatewaysFailed bool

	// Any error occurred when syncing with Konnect,
	KonnectFailed bool

	// Translation of some of Kubernetes objects failed.
	TranslationFailuresOccurred bool
}

func (i CalculateConfigStatusInput) SetKonnectFailed(failed bool) (CalculateConfigStatusInput, bool) {
	if i.KonnectFailed != failed {
		i.KonnectFailed = failed
		return i, true
	}
	return i, false
}

// CalculateConfigStatus calculates a clients.ConfigStatus that sums up the configuration synchronisation result as
// a single enumerated value.
func CalculateConfigStatus(i CalculateConfigStatusInput) ConfigStatus {
	switch {
	case !i.GatewaysFailed && !i.KonnectFailed && !i.TranslationFailuresOccurred:
		return ConfigStatusOK
	case !i.GatewaysFailed && !i.KonnectFailed && i.TranslationFailuresOccurred:
		return ConfigStatusTranslationErrorHappened
	case i.GatewaysFailed && !i.KonnectFailed: // We don't care about translation failures if we can't apply to gateways.
		return ConfigStatusApplyFailed
	case !i.GatewaysFailed && i.KonnectFailed && !i.TranslationFailuresOccurred:
		return ConfigStatusOKKonnectApplyFailed
	case !i.GatewaysFailed && i.KonnectFailed && i.TranslationFailuresOccurred:
		return ConfigStatusTranslationErrorHappenedKonnectApplyFailed
	case i.GatewaysFailed && i.KonnectFailed: // We don't care about translation failures if we can't apply to gateways.
		return ConfigStatusApplyFailedKonnectApplyFailed
	}

	// Shouldn't happen.
	return ConfigStatusUnknown
}

type ConfigStatusNotifier interface {
	NotifyConfigStatus(context.Context, ConfigStatus)
}

type ConfigStatusSubscriber interface {
	SubscribeConfigStatus() chan ConfigStatus
}

type NoOpConfigStatusNotifier struct{}

var _ ConfigStatusNotifier = NoOpConfigStatusNotifier{}

func (n NoOpConfigStatusNotifier) NotifyConfigStatus(_ context.Context, _ ConfigStatus) {
}

type ChannelConfigNotifier struct {
	ch     chan ConfigStatus
	logger logr.Logger
}

var _ ConfigStatusNotifier = &ChannelConfigNotifier{}

func NewChannelConfigNotifier(logger logr.Logger) *ChannelConfigNotifier {
	return &ChannelConfigNotifier{
		ch:     make(chan ConfigStatus),
		logger: logger,
	}
}

// NotifyConfigStatus sends the status in a separate goroutine. If the notification is not received in 1s, it's dropped.
func (n *ChannelConfigNotifier) NotifyConfigStatus(ctx context.Context, status ConfigStatus) {
	const notifyTimeout = time.Second

	go func() {
		timeout := time.NewTimer(notifyTimeout)
		defer timeout.Stop()

		select {
		case n.ch <- status:
		case <-ctx.Done():
			n.logger.Info("Context done, not notifying config status", "status", status)
		case <-timeout.C:
			n.logger.Info("Timed out notifying config status", "status", status)
		}
	}()
}

func (n *ChannelConfigNotifier) SubscribeConfigStatus() chan ConfigStatus {
	// TODO: in case of multiple subscribers, we should use a fan-out pattern.
	return n.ch
}
