package fallback_test

import (
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// mockGraphProvider is a mock implementation of the CacheGraphProvider interface.
type mockGraphProvider struct {
	graphToReturnOn   map[store.CacheStores]*fallback.ConfigGraph
	cacheToGraphCalls []store.CacheStores
}

// CacheToGraph returns the graph that was set on the mockGraphProvider. It also records the last store that was passed to it.
func (m *mockGraphProvider) CacheToGraph(c store.CacheStores) (*fallback.ConfigGraph, error) {
	m.cacheToGraphCalls = append(m.cacheToGraphCalls, c)
	if g, ok := m.graphToReturnOn[c]; ok {
		return g, nil
	}
	return nil, errors.New("unexpected call")
}

// ReturnGraphOn sets the graph that should be returned when CacheToGraph is called with the given cache stores.
func (m *mockGraphProvider) ReturnGraphOn(cache store.CacheStores, graph *fallback.ConfigGraph) {
	if m.graphToReturnOn == nil {
		m.graphToReturnOn = make(map[store.CacheStores]*fallback.ConfigGraph)
	}
	m.graphToReturnOn[cache] = graph
}

// CacheToGraphLastCalledWith returns the last cache stores that were passed to CacheToGraph.
func (m *mockGraphProvider) CacheToGraphLastCalledWith() store.CacheStores {
	if len(m.cacheToGraphCalls) == 0 {
		return store.CacheStores{}
	}
	return m.cacheToGraphCalls[len(m.cacheToGraphCalls)-1]
}

// CacheToGraphLastNCalledWith returns the last N cache stores that were passed to CacheToGraph.
func (m *mockGraphProvider) CacheToGraphLastNCalledWith(n int) []store.CacheStores {
	if maxLen := len(m.cacheToGraphCalls); n > maxLen {
		n = maxLen
	}
	return m.cacheToGraphCalls[len(m.cacheToGraphCalls)-n:]
}

func TestGenerator_GenerateExcludingBrokenObjects(t *testing.T) {
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

	graphProvider := &mockGraphProvider{}
	graphProvider.ReturnGraphOn(inputCacheStores, graph)
	g := fallback.NewGenerator(graphProvider, logr.Discard())

	t.Run("ingressClass is broken", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateExcludingBrokenObjects(inputCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(ingressClass)})
		require.NoError(t, err)
		require.Equal(t, inputCacheStores, graphProvider.CacheToGraphLastCalledWith(), "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, &inputCacheStores, &fallbackCache)
		require.Empty(t, fallbackCache.IngressClassV1.List(), "ingressClass should be excluded as it's broken")
		require.Empty(t, fallbackCache.Service.List(), "service should be excluded as it depends on ingressClass")
		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade should be excluded as it depends on service")
		require.ElementsMatch(t, fallbackCache.Plugin.List(), []any{plugin}, "plugin shouldn't be excluded as it doesn't depend on ingressClass")
	})

	t.Run("service is broken", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateExcludingBrokenObjects(inputCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(service)})
		require.NoError(t, err)
		require.Equal(t, inputCacheStores, graphProvider.CacheToGraphLastCalledWith(), "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, &inputCacheStores, &fallbackCache)
		require.Empty(t, fallbackCache.Service.List(), "service should be excluded as it's broken")
		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade should be excluded as it depends on service")
		require.ElementsMatch(t, fallbackCache.IngressClassV1.List(), []any{ingressClass}, "ingressClass shouldn't be excluded as it doesn't depend on service")
		require.ElementsMatch(t, fallbackCache.Plugin.List(), []any{plugin}, "plugin shouldn't be excluded as it doesn't depend on service")
	})

	t.Run("serviceFacade is broken", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateExcludingBrokenObjects(inputCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(serviceFacade)})
		require.NoError(t, err)
		require.Equal(t, inputCacheStores, graphProvider.CacheToGraphLastCalledWith(), "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, &inputCacheStores, &fallbackCache)
		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade should be excluded as it's broken")
		require.ElementsMatch(t, fallbackCache.IngressClassV1.List(), []any{ingressClass}, "ingressClass shouldn't be excluded as it doesn't depend on service")
		require.ElementsMatch(t, fallbackCache.Service.List(), []any{service}, "service shouldn't be excluded as it doesn't depend on serviceFacade")
		require.ElementsMatch(t, fallbackCache.Plugin.List(), []any{plugin}, "plugin shouldn't be excluded as it doesn't depend on serviceFacade")
	})

	t.Run("plugin is broken", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateExcludingBrokenObjects(inputCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(plugin)})
		require.NoError(t, err)
		require.Equal(t, inputCacheStores, graphProvider.CacheToGraphLastCalledWith(), "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, &inputCacheStores, &fallbackCache)
		require.Empty(t, fallbackCache.Plugin.List(), "plugin should be excluded as it's broken")
		require.ElementsMatch(t, fallbackCache.IngressClassV1.List(), []any{ingressClass}, "ingressClass shouldn't be excluded as it doesn't depend on plugin")
		require.ElementsMatch(t, fallbackCache.Service.List(), []any{service}, "service shouldn't be excluded as it doesn't depend on plugin")
		require.ElementsMatch(t, fallbackCache.KongServiceFacade.List(), []any{serviceFacade}, "serviceFacade shouldn't be excluded as it doesn't depend on plugin")
	})

	t.Run("multiple objects are broken", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateExcludingBrokenObjects(inputCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(ingressClass), fallback.GetObjectHash(service)})
		require.NoError(t, err)
		require.Equal(t, inputCacheStores, graphProvider.CacheToGraphLastCalledWith(), "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, &inputCacheStores, &fallbackCache)
		require.Empty(t, fallbackCache.IngressClassV1.List(), "ingressClass should be excluded as it's broken")
		require.Empty(t, fallbackCache.Service.List(), "service should be excluded as it's broken")
		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade should be excluded as it depends on service")
		require.ElementsMatch(t, fallbackCache.Plugin.List(), []any{plugin}, "plugin shouldn't be excluded as it doesn't depend on either ingressClass or service")
	})
}

func TestGenerator_GenerateBackfillingBrokenObjects(t *testing.T) {
	// We have to use real-world object types here as we're testing integration with store.CacheStores.
	ingressClass := testIngressClass(t, "ingressClass")
	service := testService(t, "service")
	serviceFacade := testKongServiceFacade(t, "serviceFacade")
	plugin := testKongPlugin(t, "kongPlugin")
	inputCacheStores := cacheStoresFromObjs(t, ingressClass, service, serviceFacade, plugin)

	// We'll annotate the last valid objects, so we can verify that they were recovered from the last valid cache snapshot.
	const lastValidAnnotationKey = "from-last-valid"
	annotatedLastValid := func(o client.Object) client.Object {
		o.SetAnnotations(map[string]string{lastValidAnnotationKey: "true"})
		return o
	}
	requireAnnotatedLastValid := func(t *testing.T, o client.Object) {
		require.Equal(t, "true", o.GetAnnotations()[lastValidAnnotationKey], "expected the object to be recovered from the last valid cache snapshot")
	}
	requireNotAnnotatedLastValid := func(t *testing.T, o client.Object) {
		require.NotContains(t, o.GetAnnotations(), lastValidAnnotationKey, "expected the object to not be recovered from the last valid cache snapshot")
	}

	lastValidIngressClass := annotatedLastValid(testIngressClass(t, "ingressClass"))
	lastValidService := annotatedLastValid(testService(t, "service"))
	lastValidPlugin := annotatedLastValid(testKongPlugin(t, "kongPlugin"))
	lastValidCacheStores := cacheStoresFromObjs(t, lastValidIngressClass, lastValidService, lastValidPlugin)

	// This graph doesn't reflect real dependencies between the objects - it's only used for testing purposes.
	// It will be injected into the Generator via the mockGraphProvider.
	// Dependency resolving between the objects is tested in TestResolveDependencies_* tests.
	//
	// Graph structure (edges define dependency -> dependant relationship):
	//   ┌────────────┐  ┌──────┐
	//   │ingressClass│  │plugin│
	//   └──────┬─────┘  └───┬──┘
	//          │            │
	//          ├────────────┘
	//          │
	//      ┌───▼───┐
	//      │service│
	//      └───┬───┘
	//          │
	//   ┌──────▼──────┐
	//   │serviceFacade│
	//   └─────────────┘
	currentGraph, err := NewGraphBuilder().
		WithVertices(ingressClass, service, serviceFacade, plugin).
		WithEdge(ingressClass, service).
		WithEdge(plugin, service).
		WithEdge(service, serviceFacade).
		Build()
	require.NoError(t, err)

	// Last valid graph differs from the input graph by lack of the serviceFacade.
	//   ┌────────────┐  ┌──────┐
	//   │ingressClass│  │plugin│
	//   └──────┬─────┘  └───┬──┘
	//          │            │
	//          ├────────────┘
	//          │
	//      ┌───▼───┐
	//      │service│
	//      └───────┘
	lastValidGraph, err := NewGraphBuilder().
		WithVertices(lastValidIngressClass, lastValidService, lastValidPlugin).
		WithEdge(lastValidIngressClass, lastValidService).
		WithEdge(lastValidPlugin, lastValidService).
		Build()
	require.NoError(t, err)

	graphProvider := &mockGraphProvider{}
	graphProvider.ReturnGraphOn(inputCacheStores, currentGraph)
	graphProvider.ReturnGraphOn(lastValidCacheStores, lastValidGraph)
	g := fallback.NewGenerator(graphProvider, logr.Discard())

	t.Run("ingressClass is broken", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateBackfillingBrokenObjects(inputCacheStores, &lastValidCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(ingressClass)})
		require.NoError(t, err)
		require.Equal(t, []store.CacheStores{inputCacheStores, lastValidCacheStores}, graphProvider.CacheToGraphLastNCalledWith(2),
			"expected the generator to call CacheToGraph with the input cache stores and the last valid cache stores")
		require.NotSame(t, &inputCacheStores, &fallbackCache)

		fallbackIngressClass, err := getFromStore[*netv1.IngressClass](fallbackCache.IngressClassV1, ingressClass)
		require.NoError(t, err)
		requireAnnotatedLastValid(t, fallbackIngressClass)

		fallbackService, err := getFromStore[*corev1.Service](fallbackCache.Service, service)
		require.NoError(t, err)
		requireAnnotatedLastValid(t, fallbackService)

		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade shouldn't be recovered as it wasn't in the last valid cache snapshot")

		fallbackPlugin, err := getFromStore[*kongv1.KongPlugin](fallbackCache.Plugin, plugin)
		require.NoError(t, err)
		requireNotAnnotatedLastValid(t, fallbackPlugin)
	})

	t.Run("service is broken", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateBackfillingBrokenObjects(inputCacheStores, &lastValidCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(service)})
		require.NoError(t, err)
		require.Equal(t, []store.CacheStores{inputCacheStores, lastValidCacheStores}, graphProvider.CacheToGraphLastNCalledWith(2),
			"expected the generator to call CacheToGraph with the input cache stores and fallback cache")
		require.NotSame(t, &inputCacheStores, &fallbackCache)

		fallbackService, err := getFromStore[*corev1.Service](fallbackCache.Service, service)
		require.NoError(t, err)
		requireAnnotatedLastValid(t, fallbackService)

		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade shouldn't be recovered as it wasn't in the last valid cache snapshot")

		fallbackPlugin, err := getFromStore[*kongv1.KongPlugin](fallbackCache.Plugin, plugin)
		require.NoError(t, err)
		requireNotAnnotatedLastValid(t, fallbackPlugin)

		fallbackIngressClass, err := getFromStore[*netv1.IngressClass](fallbackCache.IngressClassV1, ingressClass)
		require.NoError(t, err)
		requireNotAnnotatedLastValid(t, fallbackIngressClass)
	})

	t.Run("serviceFacade is broken", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateBackfillingBrokenObjects(inputCacheStores, &lastValidCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(serviceFacade)})
		require.NoError(t, err)
		require.Equal(t, []store.CacheStores{inputCacheStores, lastValidCacheStores}, graphProvider.CacheToGraphLastNCalledWith(2), "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, &inputCacheStores, &fallbackCache)

		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade should be excluded as it's broken and it's not present in the last valid cache snapshot")

		fallbackIngressClass, err := getFromStore[*netv1.IngressClass](fallbackCache.IngressClassV1, ingressClass)
		require.NoError(t, err)
		requireNotAnnotatedLastValid(t, fallbackIngressClass)

		fallbackService, err := getFromStore[*corev1.Service](fallbackCache.Service, service)
		require.NoError(t, err)
		requireNotAnnotatedLastValid(t, fallbackService)

		fallbackPlugin, err := getFromStore[*kongv1.KongPlugin](fallbackCache.Plugin, plugin)
		require.NoError(t, err)
		requireNotAnnotatedLastValid(t, fallbackPlugin)
	})

	t.Run("plugin is broken", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateBackfillingBrokenObjects(inputCacheStores, &lastValidCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(plugin)})
		require.NoError(t, err)
		require.Equal(t, []store.CacheStores{inputCacheStores, lastValidCacheStores}, graphProvider.CacheToGraphLastNCalledWith(2), "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, &inputCacheStores, &fallbackCache)

		fallbackPlugin, err := getFromStore[*kongv1.KongPlugin](fallbackCache.Plugin, plugin)
		require.NoError(t, err)
		requireAnnotatedLastValid(t, fallbackPlugin)

		fallbackIngressClass, err := getFromStore[*netv1.IngressClass](fallbackCache.IngressClassV1, ingressClass)
		require.NoError(t, err)
		requireNotAnnotatedLastValid(t, fallbackIngressClass)

		fallbackService, err := getFromStore[*corev1.Service](fallbackCache.Service, service)
		require.NoError(t, err, "service should be backfilled as it was present in the last valid state and is directly affected by broken plugin")
		requireAnnotatedLastValid(t, fallbackService)

		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade shouldn't be recovered as it wasn't in the last valid cache snapshot and it's indirectly affected by broken plugin")
	})

	t.Run("multiple objects are broken", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateBackfillingBrokenObjects(inputCacheStores, &lastValidCacheStores, []fallback.ObjectHash{fallback.GetObjectHash(ingressClass), fallback.GetObjectHash(service)})
		require.NoError(t, err)
		require.Equal(t, []store.CacheStores{inputCacheStores, lastValidCacheStores}, graphProvider.CacheToGraphLastNCalledWith(2), "expected the generator to call CacheToGraph with the input cache stores")
		require.NotSame(t, &inputCacheStores, &fallbackCache)

		fallbackIngressClass, err := getFromStore[*netv1.IngressClass](fallbackCache.IngressClassV1, ingressClass)
		require.NoError(t, err)
		requireAnnotatedLastValid(t, fallbackIngressClass)

		fallbackService, err := getFromStore[*corev1.Service](fallbackCache.Service, service)
		require.NoError(t, err)
		requireAnnotatedLastValid(t, fallbackService)

		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade shouldn't be recovered as it wasn't in the last valid cache snapshot")

		fallbackPlugin, err := getFromStore[*kongv1.KongPlugin](fallbackCache.Plugin, plugin)
		require.NoError(t, err)
		requireNotAnnotatedLastValid(t, fallbackPlugin)
	})

	t.Run("multiple objects are broken and last valid cache is nil", func(t *testing.T) {
		fallbackCache, _, err := g.GenerateBackfillingBrokenObjects(inputCacheStores, nil, []fallback.ObjectHash{fallback.GetObjectHash(ingressClass), fallback.GetObjectHash(service)})
		require.NoError(t, err)
		require.Equal(t, []store.CacheStores{inputCacheStores}, graphProvider.CacheToGraphLastNCalledWith(1), "expected the generator to call CacheToGraph with the input cache stores only")
		require.NotSame(t, &inputCacheStores, &fallbackCache)

		require.Empty(t, fallbackCache.IngressClassV1.List(), "ingressClass should be excluded as it's broken and there's no last valid cache snapshot")
		require.Empty(t, fallbackCache.Service.List(), "service should be excluded as it's broken and there's no last valid cache snapshot")
		require.Empty(t, fallbackCache.KongServiceFacade.List(), "serviceFacade should be excluded as it depends on service and there's no last valid cache snapshot")

		fallbackPlugin, err := getFromStore[*kongv1.KongPlugin](fallbackCache.Plugin, plugin)
		require.NoError(t, err)
		requireNotAnnotatedLastValid(t, fallbackPlugin)
	})
}

func TestGenerator_ReturnsMetadata(t *testing.T) {
	ingressClass := testIngressClass(t, "ingressClass")
	service := testService(t, "service")
	inputCacheStores := cacheStoresFromObjs(t, ingressClass, service)

	lastValidCacheStores := cacheStoresFromObjs(t, ingressClass, service)

	configGraph, err := NewGraphBuilder().
		WithVertices(ingressClass, service).
		WithEdge(ingressClass, service).
		Build()
	require.NoError(t, err)

	graphProvider := &mockGraphProvider{}
	graphProvider.ReturnGraphOn(inputCacheStores, configGraph)
	graphProvider.ReturnGraphOn(lastValidCacheStores, configGraph)
	g := fallback.NewGenerator(graphProvider, logr.Discard())

	t.Run("on excluding", func(t *testing.T) {
		_, meta, err := g.GenerateExcludingBrokenObjects(inputCacheStores, []fallback.ObjectHash{
			fallback.GetObjectHash(ingressClass),
		})
		require.NoError(t, err)
		require.Len(t, meta.BrokenObjects, 1)
		require.Len(t, meta.ExcludedObjects, 2)
		require.Len(t, meta.BackfilledObjects, 0)
	})
	t.Run("on backfilling", func(t *testing.T) {
		_, meta, err := g.GenerateBackfillingBrokenObjects(inputCacheStores, &lastValidCacheStores, []fallback.ObjectHash{
			fallback.GetObjectHash(ingressClass),
			fallback.GetObjectHash(service),
		})
		require.NoError(t, err)
		require.Len(t, meta.BrokenObjects, 2)
		require.Len(t, meta.ExcludedObjects, 2)
		require.Len(t, meta.BackfilledObjects, 2)
	})
}
