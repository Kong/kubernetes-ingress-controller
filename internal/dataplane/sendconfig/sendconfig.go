package sendconfig

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/kong/deck/diff"
	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/deck/state"
	deckutils "github.com/kong/deck/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

const initialHash = "00000000000000000000000000000000"

// -----------------------------------------------------------------------------
// Sendconfig - Public Functions
// -----------------------------------------------------------------------------

// PerformUpdate writes `targetContent` and `customEntities` to Kong Admin API specified by `kongConfig`.
func PerformUpdate(ctx context.Context,
	log logrus.FieldLogger,
	kongConfig *Kong,
	inMemory bool,
	reverseSync bool,
	targetContent *file.Content,
	selectorTags []string,
	customEntities []byte,
	oldSHA []byte,
	promMetrics *metrics.CtrlFuncMetrics) ([]byte, error) {
	newSHA, err := deckgen.GenerateSHA(targetContent, customEntities)
	if err != nil {
		return oldSHA, err
	}
	// disable optimization if reverse sync is enabled
	if !reverseSync {
		// use the previous SHA to determine whether or not to perform an update
		if equalSHA(oldSHA, newSHA) {
			if !hasSHAUpdateAlreadyBeenReported(newSHA) {
				log.Debugf("sha %s has been reported", hex.EncodeToString(newSHA))
			}
			// we assume ready as not all Kong versions provide their configuration hash, and their readiness state
			// is always unknown
			ready := true
			status, err := kongConfig.Client.Status(ctx)
			if err != nil {
				log.WithError(err).Error("checking config status failed")
				log.Debug("configuration state unknown, skipping sync to kong")
				return oldSHA, nil
			}
			if status.ConfigurationHash == initialHash {
				ready = false
			}
			if ready {
				log.Debug("no configuration change, skipping sync to kong")
				return oldSHA, nil
			}
		}
	}

	var metricsProtocol string
	timeStart := time.Now()
	if inMemory {
		metricsProtocol = metrics.ProtocolDBLess
		err = onUpdateInMemoryMode(ctx, log, targetContent, customEntities, kongConfig)
	} else {
		metricsProtocol = metrics.ProtocolDeck
		err = onUpdateDBMode(ctx, targetContent, kongConfig, selectorTags)
	}
	timeEnd := time.Now()

	if err != nil {
		promMetrics.ConfigPushCount.With(prometheus.Labels{
			metrics.SuccessKey:  metrics.SuccessFalse,
			metrics.ProtocolKey: metricsProtocol,
		}).Inc()
		promMetrics.ConfigPushDuration.With(prometheus.Labels{
			metrics.SuccessKey:  metrics.SuccessFalse,
			metrics.ProtocolKey: metricsProtocol,
		}).Observe(float64(timeEnd.Sub(timeStart).Milliseconds()))
		return nil, err
	}

	promMetrics.ConfigPushCount.With(prometheus.Labels{
		metrics.SuccessKey:  metrics.SuccessTrue,
		metrics.ProtocolKey: metricsProtocol,
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

func renderConfigWithCustomEntities(log logrus.FieldLogger, state *file.Content,
	customEntitiesJSONBytes []byte) ([]byte, error) {

	var kongCoreConfig []byte
	var err error

	kongCoreConfig, err = json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("marshaling kong config into json: %w", err)
	}

	// fast path
	if len(customEntitiesJSONBytes) == 0 {
		return kongCoreConfig, nil
	}

	// slow path
	mergeMap := map[string]interface{}{}
	var result []byte
	var customEntities map[string]interface{}

	// unmarshal core config into the merge map
	err = json.Unmarshal(kongCoreConfig, &mergeMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling kong config into map[string]interface{}: %w", err)
	}

	// unmarshal custom entities config into the merge map
	err = json.Unmarshal(customEntitiesJSONBytes, &customEntities)
	if err != nil {
		// do not error out when custom entities are messed up
		log.WithError(err).Error("failed to unmarshal custom entities from secret data")
	} else {
		for k, v := range customEntities {
			if _, exists := mergeMap[k]; !exists {
				mergeMap[k] = v
			}
		}
	}

	// construct the final configuration
	result, err = json.Marshal(mergeMap)
	if err != nil {
		err = fmt.Errorf("marshaling final config into JSON: %w", err)
		return nil, err
	}

	return result, nil
}

func onUpdateInMemoryMode(ctx context.Context,
	log logrus.FieldLogger,
	state *file.Content,
	customEntities []byte,
	kongConfig *Kong,
) error {
	// Kong will error out if this is set
	state.Info = nil
	// Kong errors out if `null`s are present in `config` of plugins
	deckgen.CleanUpNullsInPluginConfigs(state)

	config, err := renderConfigWithCustomEntities(log, state, customEntities)
	if err != nil {
		return fmt.Errorf("constructing kong configuration: %w", err)
	}

	req, err := http.NewRequest("POST", kongConfig.URL+"/config",
		bytes.NewReader(config))
	if err != nil {
		return fmt.Errorf("creating new HTTP request for /config: %w", err)
	}
	req.Header.Add("content-type", "application/json")

	queryString := req.URL.Query()
	queryString.Add("check_hash", "1")

	req.URL.RawQuery = queryString.Encode()

	_, err = kongConfig.Client.Do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("posting new config to /config: %w", err)
	}

	return err
}

func onUpdateDBMode(ctx context.Context,
	targetContent *file.Content,
	kongConfig *Kong,
	selectorTags []string,
) error {
	dumpConfig := dump.Config{SelectorTags: selectorTags}
	// read the current state
	rawState, err := dump.Get(ctx, kongConfig.Client, dumpConfig)
	if err != nil {
		return fmt.Errorf("loading configuration from kong: %w", err)
	}
	currentState, err := state.Get(rawState)
	if err != nil {
		return err
	}

	// read the target state
	rawState, err = file.Get(ctx, targetContent, file.RenderConfig{
		CurrentState: currentState,
		KongVersion:  kongConfig.Version,
	}, dumpConfig, kongConfig.Client)
	if err != nil {
		return err
	}
	targetState, err := state.Get(rawState)
	if err != nil {
		return err
	}

	syncer, err := diff.NewSyncer(diff.SyncerOpts{
		CurrentState:    currentState,
		TargetState:     targetState,
		KongClient:      kongConfig.Client,
		SilenceWarnings: true,
	})
	if err != nil {
		return fmt.Errorf("creating a new syncer: %w", err)
	}
	_, errs := syncer.Solve(ctx, kongConfig.Concurrency, false)
	if errs != nil {
		return deckutils.ErrArray{Errors: errs}
	}
	return nil
}

func equalSHA(a, b []byte) bool {
	return reflect.DeepEqual(a, b)
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
//       but in the future we might configure rolling this into
//       some object/interface which has this functionality as an
//       inherent behavior.
func hasSHAUpdateAlreadyBeenReported(latestUpdateSHA []byte) bool {
	shaLock.Lock()
	defer shaLock.Unlock()
	if equalSHA(latestReportedSHA, latestUpdateSHA) {
		return true
	}
	latestReportedSHA = latestUpdateSHA
	return false
}
