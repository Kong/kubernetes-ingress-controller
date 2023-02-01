package sendconfig

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/blang/semver/v4"
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
		configurationChanged, err := hasConfigurationChanged(ctx, oldSHA, newSHA, client, log)
		if err != nil {
			return nil, err
		}
		if !configurationChanged {
			log.Debug("no configuration change, skipping sync to kong")
			return oldSHA, nil
		}
	}

	updateStrategy := ProvideUpdateStrategy(
		client,
		config,
	)

	timeStart := time.Now()
	err = updateStrategy.Update(ctx, targetContent)
	timeEnd := time.Now()

	updateDuration := float64(timeEnd.Sub(timeStart).Milliseconds())
	metricsProtocol := updateStrategy.MetricsProtocol()

	if err != nil {
		promMetrics.ConfigPushCount.With(prometheus.Labels{
			metrics.SuccessKey:       metrics.SuccessFalse,
			metrics.ProtocolKey:      metricsProtocol,
			metrics.FailureReasonKey: pushFailureReason(err),
		}).Inc()
		promMetrics.ConfigPushDuration.With(prometheus.Labels{
			metrics.SuccessKey:  metrics.SuccessFalse,
			metrics.ProtocolKey: metricsProtocol,
		}).Observe(updateDuration)
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
	}).Observe(updateDuration)
	log.Info("successfully synced configuration to kong.")
	return newSHA, nil
}

// -----------------------------------------------------------------------------
// Sendconfig - Private Functions
// -----------------------------------------------------------------------------

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

// hasConfigurationChanged verifies whether configuration has changed by comparing old and new config's SHAs.
// In case the SHAs are equal, it still can return true if a client is considered crashed based on its status.
func hasConfigurationChanged(
	ctx context.Context,
	oldSHA, newSHA []byte,
	client *adminapi.Client,
	log logrus.FieldLogger,
) (bool, error) {
	if !bytes.Equal(oldSHA, newSHA) {
		return true, nil
	}
	if !hasSHAUpdateAlreadyBeenReported(newSHA) {
		log.Debugf("sha %s has been reported", hex.EncodeToString(newSHA))
	}
	if !client.IsKonnect() {
		crashed, err := hasCrashed(ctx, client, log)
		if err != nil {
			return false, fmt.Errorf("failed to verify kong readiness: %w", err)
		}
		// Kong instance has crashed, we should push config despite the oldSHA and newSHA being equal.
		if crashed {
			return true, nil
		}
	}

	return false, nil
}

// hasCrashed checks Kong's status endpoint and read its config hash.
// If the config hash reported by Kong is the known empty hash, it's considered crashed.
// This allows providing configuration to Kong instances that have unexpectedly crashed and
// lost their configuration.
func hasCrashed(ctx context.Context, client *adminapi.Client, log logrus.FieldLogger) (bool, error) {
	status, err := client.AdminAPIClient().Status(ctx)
	if err != nil {
		return false, err
	}

	const initialHash = "00000000000000000000000000000000"
	if crashed := status.ConfigurationHash == initialHash; crashed {
		log.Debugf("starting to send configuration (hash: %s)", status.ConfigurationHash)
		return true, nil
	}

	return false, nil
}
