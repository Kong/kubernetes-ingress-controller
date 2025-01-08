package mocks

import (
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
)

// KonnectKongStateUpdater is a mock implementation of dataplane.KonnectKongStateUpdater.
type KonnectKongStateUpdater struct {
	calls []KonnectKongStateUpdaterCall
}

type KonnectKongStateUpdaterCall struct {
	KongState  *kongstate.KongState
	IsFallback bool
}

func (k *KonnectKongStateUpdater) UpdateKongState(kongState *kongstate.KongState, isFallback bool) {
	k.calls = append(k.calls, KonnectKongStateUpdaterCall{
		KongState:  kongState,
		IsFallback: isFallback,
	})
}

func (k *KonnectKongStateUpdater) Calls() []KonnectKongStateUpdaterCall {
	return k.calls
}
