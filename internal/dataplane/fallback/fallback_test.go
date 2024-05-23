package fallback_test

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// mockGraphProvider is a mock implementation of the CacheGraphProvider interface.
type mockGraphProvider struct {
	graph               *fallback.ConfigGraph
	lastCalledWithStore store.CacheStores
}

// CacheToGraph returns the graph that was set on the mockGraphProvider. It also records the last store that was passed to it.
func (m *mockGraphProvider) CacheToGraph(s store.CacheStores) (*fallback.ConfigGraph, error) {
	m.lastCalledWithStore = s
	return m.graph, nil
}

func TestGenerator_GenerateExcludingAffected(t *testing.T) {
	// We have to use real-world object types here as we're testing integration with store.CacheStores.
	ingressClass := testIngressClass(t, "ingressClass")
	service := testService(t, "service")
	serviceFacade := testKongServiceFacade(t, "serviceFacade")
	plugin := testKongPlugin(t, "kongPlugin")
	inputCacheStores := cacheStoresFromObjs(t, ingressClass, service, serviceFacade, plugin)

	// This graph doesn't reflect real dependencies between the objects - it's only used for testing purposes.
	// It will be injected into the Generator via the mockGraphProvider.
	// Dependency resolving between the objects is tested in TestResolveDependencies_* tests.
	//
	// Graph structure (edges define dependency -> dependant relationship):
	//  ┌────────────┐  ┌──────┐
	//  │ingressClass│  │plugin│
	//  └──────┬─────┘  └──────┘
	//         │
	//     ┌───▼───┐
	//     │service│
	//     └───┬───┘
	//         │
	//  ┌──────▼──────┐
	//  │serviceFacade│
	//  └─────────────┘
	graph, err := NewGraphBuilder().
		WithVertices(ingressClass, service, serviceFacade, plugin).
		WithEdge(ingressClass, service).
		WithEdge(service, serviceFacade).
		Build()
	require.NoError(t, err)

	graphProvider := &mockGraphProvider{graph: graph}
	g := fallback.NewGenerator(graphProvider, logr.Discard())

	t.Run("ingressClass is broken", func(t *testing.T) {
		fallbackCache, err := g.GenerateExcludingAffected(inputCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(ingressClass)})
		require.NoError(t, err)
		require.Equal(t, inputCacheStores, graphProvider.lastCalledWithStore, "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, inputCacheStores, fallbackCache)
		require.Empty(t, fallbackCache.IngressClassV1.List(), "ingressClass should be excluded as it's broken")
		require.Empty(t, fallbackCache.Service.List(), "service should be excluded as it depends on ingressClass")
		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade should be excluded as it depends on service")
		require.ElementsMatch(t, fallbackCache.Plugin.List(), []any{plugin}, "plugin shouldn't be excluded as it doesn't depend on ingressClass")
	})

	t.Run("service is broken", func(t *testing.T) {
		fallbackCache, err := g.GenerateExcludingAffected(inputCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(service)})
		require.NoError(t, err)
		require.Equal(t, inputCacheStores, graphProvider.lastCalledWithStore, "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, inputCacheStores, fallbackCache)
		require.Empty(t, fallbackCache.Service.List(), "service should be excluded as it's broken")
		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade should be excluded as it depends on service")
		require.ElementsMatch(t, fallbackCache.IngressClassV1.List(), []any{ingressClass}, "ingressClass shouldn't be excluded as it doesn't depend on service")
		require.ElementsMatch(t, fallbackCache.Plugin.List(), []any{plugin}, "plugin shouldn't be excluded as it doesn't depend on service")
	})

	t.Run("serviceFacade is broken", func(t *testing.T) {
		fallbackCache, err := g.GenerateExcludingAffected(inputCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(serviceFacade)})
		require.NoError(t, err)
		require.Equal(t, inputCacheStores, graphProvider.lastCalledWithStore, "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, inputCacheStores, fallbackCache)
		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade should be excluded as it's broken")
		require.ElementsMatch(t, fallbackCache.IngressClassV1.List(), []any{ingressClass}, "ingressClass shouldn't be excluded as it doesn't depend on service")
		require.ElementsMatch(t, fallbackCache.Service.List(), []any{service}, "service shouldn't be excluded as it doesn't depend on serviceFacade")
		require.ElementsMatch(t, fallbackCache.Plugin.List(), []any{plugin}, "plugin shouldn't be excluded as it doesn't depend on serviceFacade")
	})

	t.Run("plugin is broken", func(t *testing.T) {
		fallbackCache, err := g.GenerateExcludingAffected(inputCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(plugin)})
		require.NoError(t, err)
		require.Equal(t, inputCacheStores, graphProvider.lastCalledWithStore, "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, inputCacheStores, fallbackCache)
		require.Empty(t, fallbackCache.Plugin.List(), "plugin should be excluded as it's broken")
		require.ElementsMatch(t, fallbackCache.IngressClassV1.List(), []any{ingressClass}, "ingressClass shouldn't be excluded as it doesn't depend on plugin")
		require.ElementsMatch(t, fallbackCache.Service.List(), []any{service}, "service shouldn't be excluded as it doesn't depend on plugin")
		require.ElementsMatch(t, fallbackCache.KongServiceFacade.List(), []any{serviceFacade}, "serviceFacade shouldn't be excluded as it doesn't depend on plugin")
	})

	t.Run("multiple objects are broken", func(t *testing.T) {
		fallbackCache, err := g.GenerateExcludingAffected(inputCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(ingressClass), fallback.GetObjectHash(service)})
		require.NoError(t, err)
		require.Equal(t, inputCacheStores, graphProvider.lastCalledWithStore, "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, inputCacheStores, fallbackCache)
		require.Empty(t, fallbackCache.IngressClassV1.List(), "ingressClass should be excluded as it's broken")
		require.Empty(t, fallbackCache.Service.List(), "service should be excluded as it's broken")
		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade should be excluded as it depends on service")
		require.ElementsMatch(t, fallbackCache.Plugin.List(), []any{plugin}, "plugin shouldn't be excluded as it doesn't depend on either ingressClass or service")
	})
}
