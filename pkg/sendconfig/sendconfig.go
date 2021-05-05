package sendconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/kong/deck/diff"
	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/deck/solver"
	"github.com/kong/deck/state"
	deckutils "github.com/kong/deck/utils"
	"github.com/kong/kubernetes-ingress-controller/pkg/deckgen"
	"github.com/sirupsen/logrus"
)

func equalSHA(a, b []byte) bool {
	return reflect.DeepEqual(a, b)
}

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
) ([]byte, error) {

	newSHA, err := deckgen.GenerateSHA(targetContent, customEntities)
	if err != nil {
		return oldSHA, err
	}
	// disable optimization if reverse sync is enabled
	if !reverseSync {
		// use the previous SHA to determine whether or not to perform an update
		if equalSHA(oldSHA, newSHA) {
			log.Info("no configuration change, skipping sync to kong")
			return oldSHA, nil
		}
	}

	if inMemory {
		err = onUpdateInMemoryMode(ctx, log, targetContent, customEntities, kongConfig)
	} else {
		err = onUpdateDBMode(targetContent, kongConfig, selectorTags)
	}
	if err != nil {
		return nil, err
	}
	log.Info("successfully synced configuration to kong")
	return newSHA, nil
}

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
		log.Errorf("failed to unmarshal custom entities from secret data: %v", err)
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

func onUpdateDBMode(
	targetContent *file.Content,
	kongConfig *Kong,
	selectorTags []string,
) error {
	// read the current state
	rawState, err := dump.Get(kongConfig.Client, dump.Config{
		SelectorTags: selectorTags,
	})
	if err != nil {
		return fmt.Errorf("loading configuration from kong: %w", err)
	}
	currentState, err := state.Get(rawState)
	if err != nil {
		return err
	}

	// read the target state
	rawState, err = file.Get(targetContent, file.RenderConfig{
		CurrentState: currentState,
		KongVersion:  kongConfig.Version,
	})
	if err != nil {
		return err
	}
	targetState, err := state.Get(rawState)
	if err != nil {
		return err
	}

	syncer, err := diff.NewSyncer(currentState, targetState)
	if err != nil {
		return fmt.Errorf("creating a new syncer: %w", err)
	}
	syncer.SilenceWarnings = true
	_, errs := solver.Solve(nil, syncer, kongConfig.Client, nil, kongConfig.Concurrency, false)
	if errs != nil {
		return deckutils.ErrArray{Errors: errs}
	}
	return nil
}
