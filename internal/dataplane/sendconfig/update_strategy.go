package sendconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/kong/deck/diff"
	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

// UpdateStrategy is the way we approach updating data-plane's configuration, depending on its type.
type UpdateStrategy interface {
	// Update applies targetConfig to the data-plane.
	Update(ctx context.Context, targetContent *file.Content) error

	// MetricsProtocol returns a string describing the update strategy type to be used in metrics.
	MetricsProtocol() string
}

func ProvideUpdateStrategy(
	client *adminapi.Client,
	config Config,
) UpdateStrategy {
	if !config.InMemory || client.IsKonnect() {
		return NewUpdateStrategyDBMode(
			client.AdminAPIClient(),
			dump.Config{
				SkipCACerts:         config.SkipCACertificates,
				SelectorTags:        config.FilterTags,
				KonnectRuntimeGroup: client.KonnectRuntimeGroup(),
			},
			config.Version,
			config.Concurrency,
		)
	}

	return NewUpdateStrategyInMemory(client.AdminAPIClient().Configs)
}

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
) *UpdateStrategyDBMode {
	return &UpdateStrategyDBMode{
		client:      client,
		dumpConfig:  dumpConfig,
		version:     version,
		concurrency: concurrency,
	}
}

func (s UpdateStrategyDBMode) Update(ctx context.Context, targetContent *file.Content) error {
	cs, err := currentState(ctx, s.client, s.dumpConfig)
	if err != nil {
		return fmt.Errorf("failed getting current state for %s: %w", s.client.BaseRootURL(), err)
	}

	ts, err := targetState(ctx, targetContent, cs, s.version, s.client, s.dumpConfig)
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

type UpdateStrategyInMemory struct {
	configService kong.AbstractConfigService
}

func NewUpdateStrategyInMemory(
	configService kong.AbstractConfigService,
) *UpdateStrategyInMemory {
	return &UpdateStrategyInMemory{
		configService: configService,
	}
}

func (s UpdateStrategyInMemory) Update(
	ctx context.Context,
	targetState *file.Content,
) error {
	// Kong will error out if this is set
	targetState.Info = nil

	// Kong errors out if `null`s are present in `config` of plugins
	deckgen.CleanUpNullsInPluginConfigs(targetState)

	config, err := json.Marshal(targetState)
	if err != nil {
		return fmt.Errorf("constructing kong configuration: %w", err)
	}

	if err = s.configService.ReloadDeclarativeRawConfig(ctx, bytes.NewReader(config), true); err != nil {
		return err
	}

	return nil
}

func (s UpdateStrategyInMemory) MetricsProtocol() string {
	return metrics.ProtocolDBLess
}
