package clients

import (
	"context"
	"time"

	"github.com/go-logr/logr"
)

type ConfigStatus int

const (
	// ConfigStatusOK: no error happens in translation from k8s objects to kong configuration
	// and succeeded to apply kong configuration to kong gateway.
	ConfigStatusOK ConfigStatus = iota
	// ConfigStatusTranslationErrorHappened: error happened in translation of k8s objects
	// but succeeded to apply kong configuration for remaining objects.
	ConfigStatusTranslationErrorHappened
	// ConfigStatusApplyFailed: failed to apply kong configurations.
	ConfigStatusApplyFailed
)

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

func NewChannelConfigNotifier(logger logr.Logger) *ChannelConfigNotifier {
	return &ChannelConfigNotifier{
		ch:     make(chan ConfigStatus, 1),
		logger: logger,
	}
}
