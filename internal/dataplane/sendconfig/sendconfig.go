package sendconfig

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
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
	version semver.Version,
	concurrency int,
	inMemory bool,
	reverseSync bool,
	skipCACertificates bool,
	targetContent *file.Content,
	selectorTags []string,
	oldSHA []byte,
	promMetrics *metrics.CtrlFuncMetrics,
) ([]byte, error, []failures.ResourceFailure) {
	newSHA, err := deckgen.GenerateSHA(targetContent)
	if err != nil {
		return oldSHA, err, []failures.ResourceFailure{}
	}

	// disable optimization if reverse sync is enabled
	if !reverseSync {
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
				return nil, err, []failures.ResourceFailure{}
			}

			if status.ConfigurationHash == initialHash {
				ready = false
			}

			if ready {
				log.Debug("no configuration change, skipping sync to kong")
				return oldSHA, nil, []failures.ResourceFailure{}
			}
			log.Debugf("starting to send configuration (hash: %s)", status.ConfigurationHash)
		}
	}

	var metricsProtocol string
	timeStart := time.Now()
	var errParseErr error
	var resourceErrors []ResourceError
	if inMemory {
		metricsProtocol = metrics.ProtocolDBLess
		err, resourceErrors, errParseErr = onUpdateInMemoryMode(ctx, log, targetContent, client.Client)
	} else {
		metricsProtocol = metrics.ProtocolDeck
		dumpConfig := dump.Config{SelectorTags: selectorTags, SkipCACerts: skipCACertificates}
		err = onUpdateDBMode(ctx, targetContent, client, dumpConfig, version, concurrency)
	}
	timeEnd := time.Now()

	if err != nil {
		// TODO TRM the collector model doesn't make much sense here since we generate all errors in one go and then toss
		// the instance--no immediate need to retain it, but you could. having it in parser is also a bit awkward, it
		// needs its own package. the translation name is no longer correct either
		failuresCollector, tfcErr := failures.NewResourceFailuresCollector(log)
		if errParseErr != nil {
			log.WithError(errParseErr).Error("could not parse error response from Kong")
		} else {
			if tfcErr != nil {
				log.WithError(errParseErr).Error("could not parse error response from Kong")
			}
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
					// TODO this object is incomplete and therefore breaks events a bit. they'll show up in the event
					// list with short fields populated, but won't appear in "describe resource" output. this requires
					// the UID in the reference, so we either need to get the object using the info given or store
					// the UID in tags.
					failuresCollector.PushResourceFailure(
						fmt.Sprintf("invalid %s: %s", field, problem),
						&obj)
					log.Info(fmt.Sprintf("adding failure for %s: %s = %s", ee.Name, field, problem)) // TODO remove
				}
			}
		}

		promMetrics.ConfigPushCount.With(prometheus.Labels{
			metrics.SuccessKey:       metrics.SuccessFalse,
			metrics.ProtocolKey:      metricsProtocol,
			metrics.FailureReasonKey: pushFailureReason(err),
		}).Inc()
		promMetrics.ConfigPushDuration.With(prometheus.Labels{
			metrics.SuccessKey:  metrics.SuccessFalse,
			metrics.ProtocolKey: metricsProtocol,
		}).Observe(float64(timeEnd.Sub(timeStart).Milliseconds()))
		return nil, err, failuresCollector.PopResourceFailures()
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
	return newSHA, nil, []failures.ResourceFailure{}
}

// -----------------------------------------------------------------------------
// Sendconfig - Private Functions
// -----------------------------------------------------------------------------

func onUpdateInMemoryMode(
	ctx context.Context,
	log logrus.FieldLogger,
	state *file.Content,
	client *kong.Client,
) (error, []ResourceError, error) {
	// Kong will error out if this is set
	state.Info = nil
	// Kong errors out if `null`s are present in `config` of plugins
	deckgen.CleanUpNullsInPluginConfigs(state)
	var parseErr error
	var resourceErrors []ResourceError
	var errBody []byte

	config, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("constructing kong configuration: %w", err), resourceErrors, parseErr
	}

	log.WithField("kong_url", client.BaseRootURL()).
		Debug("sending configuration to Kong Admin API")
	if errBody, err = client.ReloadDeclarativeRawConfig(ctx, bytes.NewReader(config), true); err != nil {
		resourceErrors, parseErr = parseFlatEntityErrors(errBody, log)
		return err, resourceErrors, parseErr
	}

	// TODO TRM ditto
	return nil, resourceErrors, parseErr
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
