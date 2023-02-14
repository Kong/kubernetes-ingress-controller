package dataplane

// REVIEW: put the package here, or in internal/adminapi?

type ConfigStatus int

const (
	// ConfigStatusOK: no error happens in translation from
	ConfigStatusOK ConfigStatus = iota
	ConfigStatusTranslationErrorHappened
	ConfigStatusApplyFailed
)

type ConfigStatusNotifier interface {
	NotifyConfigStatus(ConfigStatus)
}

type NoOpConfigStatusNotifier struct {
}

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

func NewChannelConfigNotifier(ch chan ConfigStatus) ConfigStatusNotifier {
	return &ChannelConfigNotifier{
		ch: ch,
	}
}
