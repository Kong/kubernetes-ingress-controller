package sendconfig

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/deck/diff"
	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/deck/state"
	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

const initialHash = "00000000000000000000000000000000"

// -----------------------------------------------------------------------------
// Sendconfig - Public Functions
// -----------------------------------------------------------------------------

// PerformUpdate writes `targetContent` to Kong Admin API specified by `kongConfig`.
func PerformUpdate(ctx context.Context,
	log logrus.FieldLogger,
	client *adminapi.Client,
	config Config,
	targetContent *file.Content,
	oldSHA []byte,
	promMetrics *metrics.CtrlFuncMetrics,
) ([]byte, error) {
	newSHA, err := deckgen.GenerateSHA(targetContent)
	if err != nil {
		return oldSHA, err
	}

	// disable optimization if reverse sync is enabled
	if !config.EnableReverseSync {
		// use the previous SHA to determine whether or not to perform an update
		if bytes.Equal(oldSHA, newSHA) {
			if !hasSHAUpdateAlreadyBeenReported(newSHA) {
				log.Debugf("sha %s has been reported", hex.EncodeToString(newSHA))
			}

			// we assume ready as not all Kong versions provide their configuration hash,
			// and their readiness state is always unknown
			ready := true

			status, err := client.Status(ctx)
			if err != nil {
				return nil, err
			}

			if status.ConfigurationHash == initialHash {
				ready = false
			}

			if ready {
				log.Debug("no configuration change, skipping sync to kong")
				return oldSHA, nil
			}
			log.Debugf("starting to send configuration (hash: %s)", status.ConfigurationHash)
		}
	}

	var metricsProtocol string
	timeStart := time.Now()
	if config.InMemory {
		metricsProtocol = metrics.ProtocolDBLess
		err = onUpdateInMemoryMode(ctx, log, targetContent, client.Client)
	} else {
		metricsProtocol = metrics.ProtocolDeck
		dumpConfig := dump.Config{SelectorTags: config.FilterTags, SkipCACerts: config.SkipCACertificates}
		err = onUpdateDBMode(ctx, targetContent, client, dumpConfig, config.Version, config.Concurrency)
	}
	timeEnd := time.Now()

	if err != nil {
		promMetrics.ConfigPushCount.With(prometheus.Labels{
			metrics.SuccessKey:       metrics.SuccessFalse,
			metrics.ProtocolKey:      metricsProtocol,
			metrics.FailureReasonKey: pushFailureReason(err),
		}).Inc()
		promMetrics.ConfigPushDuration.With(prometheus.Labels{
			metrics.SuccessKey:  metrics.SuccessFalse,
			metrics.ProtocolKey: metricsProtocol,
		}).Observe(float64(timeEnd.Sub(timeStart).Milliseconds()))
		return nil, err
	}

	promMetrics.ConfigPushCount.With(prometheus.Labels{
		metrics.SuccessKey:       metrics.SuccessTrue,
		metrics.ProtocolKey:      metricsProtocol,
		metrics.FailureReasonKey: "",
	}).Inc()
	promMetrics.ConfigPushDuration.With(prometheus.Labels{
		metrics.SuccessKey:  metrics.SuccessTrue,
		metrics.ProtocolKey: metricsProtocol,
	}).Observe(float64(timeEnd.Sub(timeStart).Milliseconds()))
	log.Info("successfully synced configuration to kong.")
	return newSHA, nil
}

// -----------------------------------------------------------------------------
// Sendconfig - Private Functions
// -----------------------------------------------------------------------------

func onUpdateInMemoryMode(
	ctx context.Context,
	log logrus.FieldLogger,
	state *file.Content,
	client *kong.Client,
) error {
	// Kong will error out if this is set
	state.Info = nil
	// Kong errors out if `null`s are present in `config` of plugins
	deckgen.CleanUpNullsInPluginConfigs(state)

	config, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("constructing kong configuration: %w", err)
	}

	log.WithField("kong_url", client.BaseRootURL()).
		Debug("sending configuration to Kong Admin API")
	if err = client.Configs.ReloadDeclarativeRawConfig(ctx, bytes.NewReader(config), true); err != nil {
		return err
	}

	return nil
}

func onUpdateDBMode(
	ctx context.Context,
	targetContent *file.Content,
	client *adminapi.Client,
	dumpConfig dump.Config,
	version semver.Version,
	concurrency int,
) error {
	cs, err := currentState(ctx, client.Client, dumpConfig)
	if err != nil {
		return fmt.Errorf("failed getting current state for %s: %w", client.BaseRootURL(), err)
	}

	ts, err := targetState(ctx, targetContent, cs, version, client.Client, dumpConfig)
	if err != nil {
		return deckConfigConflictError{err}
	}

	syncer, err := diff.NewSyncer(diff.SyncerOpts{
		CurrentState:    cs,
		TargetState:     ts,
		KongClient:      client.Client,
		SilenceWarnings: true,
	})
	if err != nil {
		return fmt.Errorf("creating a new syncer for %s: %w", client.BaseRootURL(), err)
	}

	_, errs := syncer.Solve(ctx, concurrency, false)
	if errs != nil {
		return deckutils.ErrArray{Errors: errs}
	}

	return nil
}

func currentState(ctx context.Context, kongClient *kong.Client, dumpConfig dump.Config) (*state.KongState, error) {
	rawState, err := dump.Get(ctx, kongClient, dumpConfig)
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
	kongClient *kong.Client,
	dumpConfig dump.Config,
) (*state.KongState, error) {
	rawState, err := file.Get(ctx, targetContent, file.RenderConfig{
		CurrentState: currentState,
		KongVersion:  version,
	}, dumpConfig, kongClient)
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
