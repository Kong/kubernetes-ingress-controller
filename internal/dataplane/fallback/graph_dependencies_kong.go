package fallback

import (
	"fmt"
	"slices"

	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// resolveKongPluginDependencies resolves potential dependencies for a KongPlugin object:
// - Secret.
func resolveKongPluginDependencies(cache store.CacheStores, kongPlugin *kongv1.KongPlugin) []client.Object {
	var dependencies []client.Object
	if cf := kongPlugin.ConfigFrom; cf != nil {
		if s, ok := fetchSecret(
			cache,
			k8stypes.NamespacedName{
				Namespace: kongPlugin.Namespace,
				Name:      cf.SecretValue.Secret,
			},
		); ok {
			dependencies = append(dependencies, s)
		}
	}
	for _, cp := range kongPlugin.ConfigPatches {
		if s, ok := fetchSecret(
			cache,
			k8stypes.NamespacedName{
				Namespace: kongPlugin.Namespace,
				Name:      cp.ValueFrom.SecretValue.Secret,
			},
		); ok {
			dependencies = append(dependencies, s)
		}
	}
	return dependencies
}

// resolveKongClusterPluginDependencies resolves potential dependencies for a KongClusterPlugin object:
// - Secret.
func resolveKongClusterPluginDependencies(cache store.CacheStores, kongClusterPlugin *kongv1.KongClusterPlugin) []client.Object {
	var dependencies []client.Object
	if cf := kongClusterPlugin.ConfigFrom; cf != nil {
		if s, ok := fetchSecret(
			cache,
			k8stypes.NamespacedName{
				Namespace: cf.SecretValue.Namespace,
				Name:      cf.SecretValue.Secret,
			},
		); ok {
			dependencies = append(dependencies, s)
		}
	}
	for _, cp := range kongClusterPlugin.ConfigPatches {
		if s, ok := fetchSecret(
			cache,
			k8stypes.NamespacedName{
				Namespace: cp.ValueFrom.SecretValue.Namespace,
				Name:      cp.ValueFrom.SecretValue.Secret,
			},
		); ok {
			dependencies = append(dependencies, s)
		}
	}
	return dependencies
}

// resolveKongConsumerDependencies resolves potential dependencies for a KongConsumer object:
// - KongPlugin
// - KongClusterPlugin
// - Secret.
func resolveKongConsumerDependencies(cache store.CacheStores, kongConsumer *kongv1.KongConsumer) []client.Object {
	return slices.Concat(
		resolveObjectDependenciesPlugin(cache, kongConsumer),
		resolveKongConsumerSecretDependencies(cache, kongConsumer),
	)
}

// resolveKongConsumerSecretDependencies resolves Secret dependencies for a KongConsumer object.
func resolveKongConsumerSecretDependencies(cache store.CacheStores, kongConsumer *kongv1.KongConsumer) []client.Object {
	var dependencies []client.Object
	for _, credSecret := range kongConsumer.Credentials {
		secret, exists, err := cache.Secret.GetByKey(fmt.Sprintf("%s/%s", kongConsumer.Namespace, credSecret))
		if err == nil && exists {
			dependencies = append(dependencies, secret.(client.Object))
		}
	}
	return dependencies
}

// resolveKongConsumerGroupDependencies resolves potential dependencies for a KongConsumerGroup object:
// - KongPlugin
// - KongClusterPlugin.
func resolveKongConsumerGroupDependencies(cache store.CacheStores, kongConsumerGroup *kongv1beta1.KongConsumerGroup) []client.Object {
	return resolveObjectDependenciesPlugin(cache, kongConsumerGroup)
}

// resolveUDPIngressDependencies resolves potential dependencies for a UDPIngress object:
// - KongPlugin
// - KongClusterPlugin
// - Service.
func resolveUDPIngressDependencies(cache store.CacheStores, udpIngress *kongv1beta1.UDPIngress) []client.Object {
	dependencies := resolveObjectDependenciesPlugin(cache, udpIngress)
	for _, rule := range udpIngress.Spec.Rules {
		if service, exists, err := cache.Service.GetByKey(
			fmt.Sprintf("%s/%s", udpIngress.GetNamespace(), rule.Backend.ServiceName),
		); err == nil && exists {
			dependencies = append(dependencies, service.(client.Object))
		}
	}
	return dependencies
}

// resolveTCPIngressDependencies resolves potential dependencies for a TCPIngress object:
// - KongPlugin
// - KongClusterPlugin
// - Service.
func resolveTCPIngressDependencies(cache store.CacheStores, tcpIngress *kongv1beta1.TCPIngress) []client.Object {
	dependencies := resolveObjectDependenciesPlugin(cache, tcpIngress)
	for _, rule := range tcpIngress.Spec.Rules {
		if service, exists, err := cache.Service.GetByKey(
			fmt.Sprintf("%s/%s", tcpIngress.GetNamespace(), rule.Backend.ServiceName),
		); err == nil && exists {
			dependencies = append(dependencies, service.(client.Object))
		}
	}
	return dependencies
}

// resolveKongServiceFacadeDependencies resolves potential dependencies for a KongServiceFacade object:
// - KongPlugin
// - KongClusterPlugin
// - KongUpstreamPolicy.
func resolveKongServiceFacadeDependencies(cache store.CacheStores, kongServiceFacade *incubatorv1alpha1.KongServiceFacade) []client.Object {
	return resolveDependenciesForServiceLikeObj(cache, kongServiceFacade)
}

// resolveKongCustomEntityDependencies resolves potential dependencies for a KongCustomEntities object:
// - KongPlugin
// - KongClusterPlugin.
func resolveKongCustomEntityDependencies(cache store.CacheStores, obj *kongv1alpha1.KongCustomEntity) []client.Object {
	if obj.Spec.ParentRef == nil {
		return nil
	}

	parentRef := *obj.Spec.ParentRef
	groupMatches := parentRef.Group != nil && *parentRef.Group == kongv1.GroupVersion.Group

	if isKongPlugin := parentRef.Kind != nil && *parentRef.Kind == "KongPlugin" && groupMatches; isKongPlugin {
		// TODO: Cross-namespace references are not supported yet.
		if parentRef.Namespace != nil && *parentRef.Namespace != "" &&
			*parentRef.Namespace != obj.GetNamespace() {
			return nil
		}
		if plugin, exists, err := cache.Plugin.GetByKey(
			fmt.Sprintf("%s/%s", obj.GetNamespace(), parentRef.Name),
		); err == nil && exists {
			return []client.Object{plugin.(client.Object)}
		}
	}

	if isKongClusterPlugin := parentRef.Kind != nil && *parentRef.Kind == "KongClusterPlugin" && groupMatches; isKongClusterPlugin {
		if plugin, exists, err := cache.ClusterPlugin.GetByKey(parentRef.Name); err == nil && exists {
			return []client.Object{plugin.(client.Object)}
		}
	}

	return nil
}
