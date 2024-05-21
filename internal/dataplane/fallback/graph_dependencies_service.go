package fallback

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// resolveServiceDependencies resolves potential dependencies for a Service object:
// - KongPlugin
// - KongClusterPlugin.
func resolveServiceDependencies(cache store.CacheStores, service *corev1.Service) []client.Object {
	return resolveObjectDependenciesPlugin(cache, service)
}
