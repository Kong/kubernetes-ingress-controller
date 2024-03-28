package sendconfig

import (
	"context"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/kong/go-database-reconciler/pkg/diff"
	"github.com/kong/go-database-reconciler/pkg/dump"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-database-reconciler/pkg/state"
	deckutils "github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

// UpdateStrategyDBMode implements the UpdateStrategy interface. It updates Kong's data-plane
// configuration using decK's syncer.
type UpdateStrategyDBMode struct {
	client      *kong.Client
	dumpConfig  dump.Config
	version     semver.Version
	concurrency int
	isKonnect   bool
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

func NewUpdateStrategyDBModeKonnect(
	client *kong.Client,
	dumpConfig dump.Config,
	version semver.Version,
	concurrency int,
) UpdateStrategyDBMode {
	s := NewUpdateStrategyDBMode(client, dumpConfig, version, concurrency)
	s.isKonnect = true
	return s
}

func (s UpdateStrategyDBMode) Update(ctx context.Context, targetContent ContentWithHash) (
	err error,
	resourceErrors []ResourceError,
	// this is always empty for DB mode, as it does not have a single config error, and will instead return individual
	// errors as a deckutils.ErrArray. we need it for the interface signature, however.
	rawErrBody []byte,
	resourceErrorsParseErr error,
) {
	cs, err := s.currentState(ctx)
	if err != nil {
		return fmt.Errorf("failed getting current state for %s: %w", s.client.BaseRootURL(), err), nil, nil, nil
	}

	ts, err := s.targetState(ctx, cs, targetContent.Content)
	if err != nil {
		return deckerrors.ConfigConflictError{Err: err}, nil, nil, nil
	}

	syncer, err := diff.NewSyncer(diff.SyncerOpts{
		CurrentState:    cs,
		TargetState:     ts,
		KongClient:      s.client,
		SilenceWarnings: true,
		IsKonnect:       s.isKonnect,
	})
	if err != nil {
		return fmt.Errorf("creating a new syncer for %s: %w", s.client.BaseRootURL(), err), nil, nil, nil
	}

	_, errs, _ := syncer.Solve(ctx, s.concurrency, false, false)
	if errs != nil {
		return deckutils.ErrArray{Errors: errs}, nil, nil, nil
	}

	return nil, nil, nil, nil
}

func (s UpdateStrategyDBMode) MetricsProtocol() metrics.Protocol {
	return metrics.ProtocolDeck
}

func (s UpdateStrategyDBMode) Type() string {
	return "DBMode"
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
