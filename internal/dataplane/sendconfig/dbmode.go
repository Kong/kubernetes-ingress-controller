package sendconfig

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/blang/semver/v4"
	"github.com/kong/deck/diff"
	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/deck/state"
	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"

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

func (s UpdateStrategyDBMode) Update(ctx context.Context, targetContent *file.Content) error {
	cs, err := s.currentState(ctx)
	if err != nil {
		return fmt.Errorf("failed getting current state for %s: %w", s.client.BaseRootURL(), err)
	}

	ts, err := s.targetState(ctx, cs, targetContent)
	if err != nil {
		return deckConfigConflictError{err}
	}

	syncer, err := diff.NewSyncer(diff.SyncerOpts{
		CurrentState:    cs,
		TargetState:     ts,
		KongClient:      s.client,
		SilenceWarnings: true,
	})
	if err != nil {
		return fmt.Errorf("creating a new syncer for %s: %w", s.client.BaseRootURL(), err)
	}

	_, errs := syncer.Solve(ctx, s.concurrency, false)
	if errs != nil {
		return deckutils.ErrArray{Errors: errs}
	}

	return nil
}

func (s UpdateStrategyDBMode) MetricsProtocol() string {
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

// deckConfigConflictError is an error used to wrap deck config conflict errors returned from deck functions
// transforming KongRawState to KongState (e.g. state.Get, dump.Get).
type deckConfigConflictError struct {
	err error
}

func (e deckConfigConflictError) Error() string {
	return e.err.Error()
}

func (e deckConfigConflictError) Is(target error) bool {
	_, ok := target.(deckConfigConflictError)
	return ok
}

func (e deckConfigConflictError) Unwrap() error {
	return e.err
}

// pushFailureReason extracts config push failure reason from an error returned from onUpdateInMemoryMode or onUpdateDBMode.
func pushFailureReason(err error) string {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return metrics.FailureReasonNetwork
	}

	if isConflictErr(err) {
		return metrics.FailureReasonConflict
	}

	return metrics.FailureReasonOther
}

func isConflictErr(err error) bool {
	var apiErr *kong.APIError
	if errors.As(err, &apiErr) && apiErr.Code() == http.StatusConflict ||
		errors.Is(err, deckConfigConflictError{}) {
		return true
	}

	var deckErrArray deckutils.ErrArray
	if errors.As(err, &deckErrArray) {
		for _, err := range deckErrArray.Errors {
			if isConflictErr(err) {
				return true
			}
		}
	}

	return false
}
