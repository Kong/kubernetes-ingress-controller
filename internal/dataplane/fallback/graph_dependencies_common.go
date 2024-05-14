package fallback

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// resolveObjectDependenciesPlugin resolves KongPlugin and KongClusterPlugin dependencies for an arbitrary object
// that refers them in its annotations.
func resolveObjectDependenciesPlugin(cache store.CacheStores, obj client.Object) []client.Object {
	var dependencies []client.Object
	for _, pluginName := range annotations.ExtractKongPluginsFromAnnotations(obj.GetAnnotations()) {
		if plugin, exists, err := cache.Plugin.GetByKey(
			fmt.Sprintf("%s/%s", obj.GetNamespace(), pluginName),
		); err == nil && exists {
			dependencies = append(dependencies, plugin.(client.Object))
		}

		if plugin, exists, err := cache.ClusterPlugin.GetByKey(pluginName); err == nil && exists {
			dependencies = append(dependencies, plugin.(client.Object))
		}
	}
	return dependencies
}
