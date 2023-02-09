package sendconfig

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/deck/diff"
	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/deck/state"
	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckerrors"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

// -----------------------------------------------------------------------------
// Sendconfig - Public Functions
// -----------------------------------------------------------------------------

// PerformUpdate writes `targetContent` to Kong Admin API specified by `kongConfig`.
func PerformUpdate(ctx context.Context,
	log logrus.FieldLogger,
	client *adminapi.Client,
	config Config,
	targetContent *file.Content,
	promMetrics *metrics.CtrlFuncMetrics,
) ([]byte, error) {
	oldSHA := client.LastConfigSHA()
	newSHA, err := deckgen.GenerateSHA(targetContent)
	if err != nil {
		return oldSHA, err
	}

	// disable optimization if reverse sync is enabled
	if !config.EnableReverseSync {
		configurationChanged, err := hasConfigurationChanged(ctx, oldSHA, newSHA, client, client.AdminAPIClient(), log)
		if err != nil {
			return nil, err
		}
		if !configurationChanged {
			log.Debug("no configuration change, skipping sync to Kong")
			return oldSHA, nil
		}
	}

	var metricsProtocol metrics.Protocol
	timeStart := time.Now()
	if config.InMemory {
		metricsProtocol = metrics.ProtocolDBLess
		err = onUpdateInMemoryMode(ctx, log, targetContent, client.AdminAPIClient())
	} else {
		metricsProtocol = metrics.ProtocolDeck
		dumpConfig := dump.Config{SelectorTags: config.FilterTags, SkipCACerts: config.SkipCACertificates}
		err = onUpdateDBMode(ctx, targetContent, client, dumpConfig, config)
	}
	duration := time.Since(timeStart)

	if err != nil {
		promMetrics.RecordPushFailure(metricsProtocol, duration, client.BaseRootURL(), err)
		return nil, err
	}

	promMetrics.RecordPushSuccess(metricsProtocol, duration, client.BaseRootURL())
	log.Info("successfully synced configuration to kong.")
	return newSHA, nil
}

// -----------------------------------------------------------------------------
// Sendconfig - Private Functions
// -----------------------------------------------------------------------------

type InMemoryClient interface {
	BaseRootURL() string
	ReloadDeclarativeRawConfig(ctx context.Context, config io.Reader, checkHash bool, flattenErrors bool) ([]byte, error)
}

func onUpdateInMemoryMode(
	ctx context.Context,
	log logrus.FieldLogger,
	state *file.Content,
	client InMemoryClient,
) error {
	// Kong will error out if this is set
	state.Info = nil
	// Kong errors out if `null`s are present in `config` of plugins
	deckgen.CleanUpNullsInPluginConfigs(state)

	config, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("constructing kong configuration: %w", err)
	}

	log.Debug("sending configuration to Kong Admin API")
	if _, err := client.ReloadDeclarativeRawConfig(ctx, bytes.NewReader(config), true, false); err != nil {
		return fmt.Errorf("failed sending declarative config to %s: %w", client.BaseRootURL(), err)
	}

	return nil
}

func onUpdateDBMode(
	ctx context.Context,
	targetContent *file.Content,
	client *adminapi.Client,
	dumpConfig dump.Config,
	config Config,
) error {
	cs, err := currentState(ctx, client, dumpConfig)
	if err != nil {
		return fmt.Errorf("failed getting current state for %s: %w", client.BaseRootURL(), err)
	}

	ts, err := targetState(ctx, targetContent, cs, config.Version, client, dumpConfig)
	if err != nil {
		return deckerrors.ConfigConflictError{Err: err}
	}

	syncer, err := diff.NewSyncer(diff.SyncerOpts{
		CurrentState:    cs,
		TargetState:     ts,
		KongClient:      client.AdminAPIClient(),
		SilenceWarnings: true,
	})
	if err != nil {
		return fmt.Errorf("creating a new syncer for %s: %w", client.BaseRootURL(), err)
	}

	_, errs := syncer.Solve(ctx, config.Concurrency, false)
	if errs != nil {
		return deckutils.ErrArray{Errors: errs}
	}

	return nil
}

func currentState(ctx context.Context, kongClient *adminapi.Client, dumpConfig dump.Config) (*state.KongState, error) {
	rawState, err := dump.Get(ctx, kongClient.AdminAPIClient(), dumpConfig)
	if err != nil {
		return nil, fmt.Errorf("loading configuration from kong: %w", err)
	}

	return state.Get(rawState)
}

func targetState(
	ctx context.Context,
	targetContent *file.Content,
	currentState *state.KongState,
	version semver.Version,
	kongClient *adminapi.Client,
	dumpConfig dump.Config,
) (*state.KongState, error) {
	rawState, err := file.Get(ctx, targetContent, file.RenderConfig{
		CurrentState: currentState,
		KongVersion:  version,
	}, dumpConfig, kongClient.AdminAPIClient())
	if err != nil {
		return nil, err
	}

	return state.Get(rawState)
}

var (
	latestReportedSHA []byte
	shaLock           sync.RWMutex
)

// hasSHAUpdateAlreadyBeenReported is a helper function to allow
// sendconfig internals to be aware of the last logged/reported
// update to the Kong Admin API. Given the most recent update SHA,
// it will return true/false whether or not that SHA has previously
// been reported (logged, e.t.c.) so that the caller can make
// decisions (such as staggering or stifling duplicate log lines).
//
// TODO: This is a bit of a hack for now to keep backwards compat,
// but in the future we might configure rolling this into
// some object/interface which has this functionality as an
// inherent behavior.
func hasSHAUpdateAlreadyBeenReported(latestUpdateSHA []byte) bool {
	shaLock.Lock()
	defer shaLock.Unlock()
	if bytes.Equal(latestReportedSHA, latestUpdateSHA) {
		return true
	}
	latestReportedSHA = latestUpdateSHA
	return false
}

type KonnectAwareClient interface {
	IsKonnect() bool
}

type StatusClient interface {
	Status(context.Context) (*kong.Status, error)
}

// hasConfigurationChanged verifies whether configuration has changed by comparing old and new config's SHAs.
// In case the SHAs are equal, it still can return true if a client is considered crashed based on its status.
func hasConfigurationChanged(
	ctx context.Context,
	oldSHA, newSHA []byte,
	client KonnectAwareClient,
	statusClient StatusClient,
	log logrus.FieldLogger,
) (bool, error) {
	if !bytes.Equal(oldSHA, newSHA) {
		return true, nil
	}
	if !hasSHAUpdateAlreadyBeenReported(newSHA) {
		log.Debugf("sha %s has been reported", hex.EncodeToString(newSHA))
	}
	// In case of Konnect, we skip further steps as it doesn't report its configuration hash.
	if client.IsKonnect() {
		return false, nil
	}

	hasNoConfiguration, err := kongHasNoConfiguration(ctx, statusClient, log)
	if err != nil {
		return false, fmt.Errorf("failed to verify kong readiness: %w", err)
	}
	// Kong instance has no configuration, we should push despite the oldSHA and newSHA being equal.
	if hasNoConfiguration {
		return true, nil
	}

	return false, nil
}

const wellKnownInitialHash = "00000000000000000000000000000000"

// kongHasNoConfiguration checks Kong's status endpoint and read its config hash.
// If the config hash reported by Kong is the known empty hash, it's considered crashed.
// This allows providing configuration to Kong instances that have unexpectedly crashed and
// lost their configuration.
func kongHasNoConfiguration(ctx context.Context, client StatusClient, log logrus.FieldLogger) (bool, error) {
	status, err := client.Status(ctx)
	if err != nil {
		return false, err
	}

	if hasNoConfig := status.ConfigurationHash == wellKnownInitialHash; hasNoConfig {
		log.Debugf("starting to send configuration (hash: %s)", status.ConfigurationHash)
		return true, nil
	}

	return false, nil
}
