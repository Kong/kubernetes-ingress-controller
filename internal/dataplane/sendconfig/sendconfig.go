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
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

// -----------------------------------------------------------------------------
// Sendconfig - Public Functions
// -----------------------------------------------------------------------------

// PerformUpdate writes `targetContent` to Kong Admin API specified by `kongConfig`.
func PerformUpdate(
	ctx context.Context,
	log logrus.FieldLogger,
	client *adminapi.Client,
	config Config,
	targetContent *file.Content,
	promMetrics *metrics.CtrlFuncMetrics,
) ([]byte, []failures.ResourceFailure, error) {
	oldSHA := client.LastConfigSHA()
	newSHA, err := deckgen.GenerateSHA(targetContent)
	if err != nil {
		return oldSHA, []failures.ResourceFailure{}, err
	}

	// disable optimization if reverse sync is enabled
	if !config.EnableReverseSync {
		configurationChanged, err := hasConfigurationChanged(ctx, oldSHA, newSHA, client, client.AdminAPIClient(), log)
		if err != nil {
			return nil, []failures.ResourceFailure{}, err
		}
		if !configurationChanged {
			log.Debug("no configuration change, skipping sync to Kong")
			return oldSHA, []failures.ResourceFailure{}, nil
		}
	}

	updateStrategy := ResolveUpdateStrategy(client, config, log)

	timeStart := time.Now()
	err, resourceErrors, resourceErrorsParseErr := updateStrategy.Update(ctx, targetContent)
	duration := time.Since(timeStart)

	metricsProtocol := updateStrategy.MetricsProtocol()
	if err != nil {
		resourceFailures := resourceErrorsToResourceFailures(resourceErrors, resourceErrorsParseErr, log)
		promMetrics.RecordPushFailure(metricsProtocol, duration, client.BaseRootURL(), err)
		return nil, resourceFailures, err
	}

	promMetrics.RecordPushSuccess(metricsProtocol, duration, client.BaseRootURL())
	log.Info("successfully synced configuration to kong")
	return newSHA, nil, nil
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
	// In case of Konnect, we skip further steps that are meant to detect Kong instances crash/reset
	// that are not relevant for Konnect.
	// We're sure that if oldSHA and newSHA are equal, we are safe to skip the update.
	if client.IsKonnect() {
		return false, nil
	}

	// Check if a Kong instance has no configuration yet (could mean it crashed, was rebooted, etc.).
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

// resourceErrorsToResourceFailures translates a slice of ResourceError to a slice of failures.ResourceFailure.
// In case of parseErr being not nil, it just returns a nil slice.
func resourceErrorsToResourceFailures(resourceErrors []ResourceError, parseErr error, log logrus.FieldLogger) []failures.ResourceFailure {
	if parseErr != nil {
		log.WithError(parseErr).Error("failed parsing resource errors")
		return nil
	}

	var out []failures.ResourceFailure
	for _, ee := range resourceErrors {
		obj := metav1.PartialObjectMetadata{
			TypeMeta: metav1.TypeMeta{
				Kind:       ee.Kind,
				APIVersion: ee.APIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ee.Namespace,
				Name:      ee.Name,
				UID:       types.UID(ee.UID),
			},
		}
		for field, problem := range ee.Problems {
			log.Debug(fmt.Sprintf("adding failure for %s: %s = %s", ee.Name, field, problem))
			resourceFailure, failureCreateErr := failures.NewResourceFailure(
				fmt.Sprintf("invalid %s: %s", field, problem),
				&obj,
			)
			if failureCreateErr != nil {
				log.WithError(failureCreateErr).Error("could create resource failure event")
			} else {
				out = append(out, resourceFailure)
			}
		}
	}

	return out
}
