package mocks

import (
	"context"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
)

type KonnectKongStateUpdater struct {
	calls []KonnectKongStateUpdaterCall
}

type KonnectKongStateUpdaterCall struct {
	KongState  *kongstate.KongState
	IsFallback bool
}

func (k *KonnectKongStateUpdater) UpdateKongState(_ context.Context, kongState *kongstate.KongState, isFallback bool) {
	k.calls = append(k.calls, KonnectKongStateUpdaterCall{
		KongState:  kongState,
		IsFallback: isFallback,
	})
}

func (k *KonnectKongStateUpdater) Calls() []KonnectKongStateUpdaterCall {
	return k.calls
}
