package sendconfig

import (
	"context"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/kong/deck/diff"
	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/deck/state"
	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckerrors"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

// UpdateStrategyDBMode implements the UpdateStrategy interface. It updates Kong's data-plane
// configuration using decK's syncer.
type UpdateStrategyDBMode struct {
	client      *kong.Client
	dumpConfig  dump.Config
	version     semver.Version
	concurrency int
}

func NewUpdateStrategyDBMode(
	client *kong.Client,
	dumpConfig dump.Config,
	version semver.Version,
	concurrency int,
) UpdateStrategyDBMode {
	return UpdateStrategyDBMode{
		client:      client,
		dumpConfig:  dumpConfig,
		version:     version,
		concurrency: concurrency,
	}
}

func (s UpdateStrategyDBMode) Update(ctx context.Context, targetContent *file.Content) (
	err error,
	resourceErrors []ResourceError,
	resourceErrorsParseErr error,
) {
	cs, err := s.currentState(ctx)
	if err != nil {
		return fmt.Errorf("failed getting current state for %s: %w", s.client.BaseRootURL(), err), nil, nil
	}

	ts, err := s.targetState(ctx, cs, targetContent)
	if err != nil {
		return deckerrors.ConfigConflictError{Err: err}, nil, nil
	}

	syncer, err := diff.NewSyncer(diff.SyncerOpts{
		CurrentState:    cs,
		TargetState:     ts,
		KongClient:      s.client,
		SilenceWarnings: true,
	})
	if err != nil {
		return fmt.Errorf("creating a new syncer for %s: %w", s.client.BaseRootURL(), err), nil, nil
	}

	_, errs := syncer.Solve(ctx, s.concurrency, false)
	if errs != nil {
		return deckutils.ErrArray{Errors: errs}, nil, nil
	}

	return nil, nil, nil
}

func (s UpdateStrategyDBMode) MetricsProtocol() metrics.Protocol {
	return metrics.ProtocolDeck
}

func (s UpdateStrategyDBMode) currentState(ctx context.Context) (*state.KongState, error) {
	rawState, err := dump.Get(ctx, s.client, s.dumpConfig)
	if err != nil {
		return nil, fmt.Errorf("loading configuration from kong: %w", err)
	}

	return state.Get(rawState)
}

func (s UpdateStrategyDBMode) targetState(
	ctx context.Context,
	currentState *state.KongState,
	targetContent *file.Content,
) (*state.KongState, error) {
	rawState, err := file.Get(ctx, targetContent, file.RenderConfig{
		CurrentState: currentState,
		KongVersion:  s.version,
	}, s.dumpConfig, s.client)
	if err != nil {
		return nil, err
	}

	return state.Get(rawState)
}
