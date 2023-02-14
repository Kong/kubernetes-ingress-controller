package dataplane

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
	NotifyConfigStatus(ConfigStatus)
}

type ConfigStatusSubscriber interface {
	SubscribeConfigStatus() chan ConfigStatus
}

type NoOpConfigStatusNotifier struct{}

var _ ConfigStatusNotifier = NoOpConfigStatusNotifier{}

func (n NoOpConfigStatusNotifier) NotifyConfigStatus(status ConfigStatus) {
}

type ChannelConfigNotifier struct {
	ch chan ConfigStatus
}

var _ ConfigStatusNotifier = &ChannelConfigNotifier{}

func (n *ChannelConfigNotifier) NotifyConfigStatus(status ConfigStatus) {
	n.ch <- status
}

func (n *ChannelConfigNotifier) SubscribeConfigStatus() chan ConfigStatus {
	return n.ch
}

func NewChannelConfigNotifier(ch chan ConfigStatus) *ChannelConfigNotifier {
	return &ChannelConfigNotifier{
		ch: ch,
	}
}
