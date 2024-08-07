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

// GatewayConfigApplyStatus stores the status of building Kong configuration and sending configuration to Kong gateways.
type GatewayConfigApplyStatus struct {
	// TranslationFailuresOccurred is true means Translation of some of Kubernetes objects failed.
	TranslationFailuresOccurred bool

	// Any error occurred when syncing with Gateways.
	ApplyConfigFailed bool
}

// KonnectConfigUploadStatus stores the status of uploading configuration to Konnect.
type KonnectConfigUploadStatus struct {
	Failed bool
}

// CalculateConfigStatus calculates a clients.ConfigStatus that sums up the configuration synchronisation result as
// a single enumerated value.
func CalculateConfigStatus(g GatewayConfigApplyStatus, k KonnectConfigUploadStatus) ConfigStatus {
	switch {
	case !g.ApplyConfigFailed && !g.TranslationFailuresOccurred && !k.Failed:
		return ConfigStatusOK
	case !g.ApplyConfigFailed && g.TranslationFailuresOccurred && !k.Failed:
		return ConfigStatusTranslationErrorHappened
	case g.ApplyConfigFailed && !k.Failed: // We don't care about translation failures if we can't apply to gateways.
		return ConfigStatusApplyFailed
	case !g.ApplyConfigFailed && !g.TranslationFailuresOccurred && k.Failed:
		return ConfigStatusOKKonnectApplyFailed
	case !g.ApplyConfigFailed && g.TranslationFailuresOccurred && k.Failed:
		return ConfigStatusTranslationErrorHappenedKonnectApplyFailed
	case g.ApplyConfigFailed && k.Failed: // We don't care about translation failures if we can't apply to gateways.
		return ConfigStatusApplyFailedKonnectApplyFailed
	}

	// Shouldn't happen.
	return ConfigStatusUnknown
}

type ConfigStatusNotifier interface {
	NotifyGatewayConfigStatus(context.Context, GatewayConfigApplyStatus)
	NotifyKonnectConfigStatus(context.Context, KonnectConfigUploadStatus)
}

type ConfigStatusSubscriber interface {
	SubscribeGatewayConfigStatus() chan GatewayConfigApplyStatus
	SubscribeKonnectConfigStatus() chan KonnectConfigUploadStatus
}

type NoOpConfigStatusNotifier struct{}

var _ ConfigStatusNotifier = NoOpConfigStatusNotifier{}

func (n NoOpConfigStatusNotifier) NotifyGatewayConfigStatus(_ context.Context, _ GatewayConfigApplyStatus) {
}

func (n NoOpConfigStatusNotifier) NotifyKonnectConfigStatus(_ context.Context, _ KonnectConfigUploadStatus) {
}

type ChannelConfigNotifier struct {
	gatewayStatusCh chan GatewayConfigApplyStatus
	konnectStatusCh chan KonnectConfigUploadStatus
	logger          logr.Logger
}

var _ ConfigStatusNotifier = &ChannelConfigNotifier{}

func NewChannelConfigNotifier(logger logr.Logger) *ChannelConfigNotifier {
	return &ChannelConfigNotifier{
		gatewayStatusCh: make(chan GatewayConfigApplyStatus),
		konnectStatusCh: make(chan KonnectConfigUploadStatus),
		logger:          logger,
	}
}

// NotifyGatewayConfigStatus notifies status of sending configuration to Kong gateway(s).
func (n *ChannelConfigNotifier) NotifyGatewayConfigStatus(ctx context.Context, status GatewayConfigApplyStatus) {
	const notifyTimeout = time.Second

	go func() {
		timeout := time.NewTimer(notifyTimeout)
		defer timeout.Stop()

		select {
		case n.gatewayStatusCh <- status:
		case <-ctx.Done():
			n.logger.Info("Context done, not notifying gateway config status", "status", status)
		case <-timeout.C:
			n.logger.Info("Timed out notifying gateway config status", "status", status)
		}
	}()
}

// NotifyKonnectConfigStatus notifies status of sending configuration to Konnect.
func (n *ChannelConfigNotifier) NotifyKonnectConfigStatus(ctx context.Context, status KonnectConfigUploadStatus) {
	const notifyTimeout = time.Second

	go func() {
		timeout := time.NewTimer(notifyTimeout)
		defer timeout.Stop()

		select {
		case n.konnectStatusCh <- status:
		case <-ctx.Done():
			n.logger.Info("Context done, not notifying Konnect config status", "status", status)
		case <-timeout.C:
			n.logger.Info("Timed out notifying Konnect config status", "status", status)
		}
	}()
}

func (n *ChannelConfigNotifier) SubscribeGatewayConfigStatus() chan GatewayConfigApplyStatus {
	// TODO: in case of multiple subscribers, we should use a fan-out pattern.
	return n.gatewayStatusCh
}

func (n *ChannelConfigNotifier) SubscribeKonnectConfigStatus() chan KonnectConfigUploadStatus {
	// TODO: in case of multiple subscribers, we should use a fan-out pattern.
	return n.konnectStatusCh
}
