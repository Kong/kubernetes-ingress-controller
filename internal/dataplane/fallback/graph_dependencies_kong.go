package fallback

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/5929
func resolveKongPluginDependencies(_ store.CacheStores, _ *kongv1.KongPlugin) []client.Object {
	return nil
}

// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/5929
func resolveKongClusterPluginDependencies(_ store.CacheStores, _ *kongv1.KongClusterPlugin) []client.Object {
	return nil
}

// resolveKongConsumerDependencies resolves potential dependencies for a KongConsumer object:
// - KongPlugin
// - KongClusterPlugin.
func resolveKongConsumerDependencies(cache store.CacheStores, kongConsumer *kongv1.KongConsumer) []client.Object {
	return resolveObjectDependenciesPlugin(cache, kongConsumer)
}

// resolveKongConsumerGroupDependencies resolves potential dependencies for a KongConsumerGroup object:
// - KongPlugin
// - KongClusterPlugin.
func resolveKongConsumerGroupDependencies(cache store.CacheStores, kongConsumerGroup *kongv1beta1.KongConsumerGroup) []client.Object {
	return resolveObjectDependenciesPlugin(cache, kongConsumerGroup)
}

// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/5929
func resolveUDPIngressDependencies(_ store.CacheStores, _ *kongv1beta1.UDPIngress) []client.Object {
	return nil
}

// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/5929
func resolveTCPIngressDependencies(_ store.CacheStores, _ *kongv1beta1.TCPIngress) []client.Object {
	return nil
}

// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/5929
func resolveKongServiceFacadeDependencies(_ store.CacheStores, _ *incubatorv1alpha1.KongServiceFacade) []client.Object {
	return nil
}
