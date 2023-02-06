package sendconfig

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/kong/deck/file"
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

	updateStrategy := ResolveUpdateStrategy(client, config)

	timeStart := time.Now()
	err = updateStrategy.Update(ctx, targetContent)
	timeEnd := time.Now()

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
