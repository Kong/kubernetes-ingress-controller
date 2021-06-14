package sendconfig

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/pkg/deckgen"
	"github.com/kong/kubernetes-ingress-controller/pkg/parser"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
)

// -----------------------------------------------------------------------------
// Sendconfig - Workflow Functions
// -----------------------------------------------------------------------------

// UpdateKongAdminSimple is a helper function for the most common usage of PerformUpdate() with only minimal
// upfront configuration required. This function is specialized and highly opinionated.
//
// If you're implementation needs to expand on the configuration and usage of the following inner components:
//
//   - store.Storer
//   - kongstate.Kong
//   - deckgen.ToDeckContent()
//   - sendconfig.PerformUpdate()
//
// Or any other encapsulated components this function makes all of that opaque to the caller.
// Treat this function as a very specific "workflow" to update the Kong Admin API,
// and use it as a reference to implement the workflow you need.
func UpdateKongAdminSimple(ctx context.Context,
	lastConfigSHA []byte,
	cache *store.CacheStores,
	ingressClassName string,
	deprecatedLogger logrus.FieldLogger,
	kongConfig Kong,
	enableReverseSync bool,
) ([]byte, error) {
	fmt.Printf("\n#### UpdateKongAdminSimple 1111 ", cache.KnativeIngress.List()...)
	// build the kongstate object from the Kubernetes objects in the storer
	storer := store.New(*cache, ingressClassName, false, false, false, deprecatedLogger)
	kongstate, err := parser.Build(deprecatedLogger, storer)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n#### UpdateKongAdminSimple 2222 ")
	// generate the deck configuration to be applied to the admin API
	targetConfig := deckgen.ToDeckContent(ctx, deprecatedLogger, kongstate, kongConfig.PluginSchemaStore,
		kongConfig.FilterTags)

	fmt.Printf("\n#### UpdateKongAdminSimple 33333 ")
	// apply the configuration update in Kong
	timedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	fmt.Printf("\n#### UpdateKongAdminSimple 4444 \n")
	configSHA, err := PerformUpdate(timedCtx,
		deprecatedLogger, &kongConfig,
		kongConfig.InMemory, enableReverseSync,
		targetConfig, nil, nil, lastConfigSHA,
	)
	if err != nil {
		return nil, err
	}

	return configSHA, nil
}
