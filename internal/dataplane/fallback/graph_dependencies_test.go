package fallback_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

// resolveDependenciesTestCase is a test case for the ResolveDependencies function.
type resolveDependenciesTestCase struct {
	name     string
	object   client.Object
	cache    store.CacheStores
	expected []client.Object
}

func runResolveDependenciesTest(t *testing.T, tc resolveDependenciesTestCase) {
	t.Run(tc.name, func(t *testing.T) {
		t.Run("when cache empty, does not panic and gives no dependencies", func(t *testing.T) {
			require.NotPanics(t, func() {
				dependencies := fallback.ResolveDependencies(store.NewCacheStores(), tc.object)
				require.Empty(t, dependencies, "expect no dependencies found in an empty cache")
			})
		})
		t.Run("when cache has objects, resolves dependencies as expected", func(t *testing.T) {
			dependencies := fallback.ResolveDependencies(tc.cache, tc.object)
			require.ElementsMatch(t, tc.expected, dependencies)
		})
	})
}

// cacheStoresFromObjs creates a CacheStores from the given objects.
// It assigns each object a type meta using the scheme.
func cacheStoresFromObjs(t *testing.T, objs ...runtime.Object) store.CacheStores {
	for i := range objs {
		obj := objs[i].(client.Object)
		obj = helpers.WithTypeMeta(t, obj)
		objs[i] = obj
	}
	s, err := store.NewCacheStoresFromObjs(objs...)
	require.NoError(t, err)
	return s
}
