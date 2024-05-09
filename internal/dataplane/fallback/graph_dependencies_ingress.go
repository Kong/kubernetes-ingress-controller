package fallback

import (
	"fmt"
	"slices"

	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// resolveIngressDependencies resolves potential dependencies for an Ingress object:
// - IngressClass
// - Service
// - KongServiceFacade
// - KongUpstreamPolicy
// - KongPlugin
// - KongClusterPlugin.
func resolveIngressDependencies(cache store.CacheStores, ingress *netv1.Ingress) []client.Object {
	return slices.Concat(
		resolveIngressDependenciesIngressClass(cache, ingress),
		resolveIngressDependenciesService(cache, ingress),
		resolveIngressDependenciesKongUpstreamPolicy(cache, ingress),
		resolveIngressDependenciesPlugin(cache, ingress),
	)
}

// resolveIngressDependenciesIngressClass resolves IngressClass dependencies for an Ingress object.
func resolveIngressDependenciesIngressClass(cache store.CacheStores, ingress *netv1.Ingress) []client.Object {
	resolveIngressClassName := func() string {
		// IngressClass can be referenced by the Ingress object in two ways:
		// 1. Using the IngressClassName field.
		if ingress.Spec.IngressClassName != nil && *ingress.Spec.IngressClassName != "" {
			return *ingress.Spec.IngressClassName
		}
		// 2. Using the annotation in the Ingress object.
		if ingressClassName, ok := ingress.Annotations[annotations.IngressClassKey]; ok && ingressClassName != "" {
			return ingressClassName
		}
		return "" // No IngressClass reference was found.
	}

	var dependencies []client.Object
	if ingressClassName := resolveIngressClassName(); ingressClassName != "" {
		ingressClass, exists, err := cache.IngressClassV1.GetByKey(ingressClassName)
		if err == nil && exists {
			dependencies = append(dependencies, ingressClass.(client.Object))
		}
	}
	return dependencies
}

// resolveIngressDependenciesService resolves Service and KongServiceFacade dependencies for an Ingress object.
func resolveIngressDependenciesService(cache store.CacheStores, ingress *netv1.Ingress) []client.Object {
	var dependencies []client.Object
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}
		for _, path := range rule.HTTP.Paths {
			if path.Backend.Service != nil {
				service, exists, err := cache.Service.GetByKey(fmt.Sprintf("%s/%s", ingress.GetNamespace(), path.Backend.Service.Name))
				if err == nil && exists {
					dependencies = append(dependencies, service.(client.Object))
				}
			}

			if resource := path.Backend.Resource; resource != nil && subtranslator.IsKongServiceFacade(resource) {
				kongServiceFacade, exists, err := cache.KongServiceFacade.GetByKey(fmt.Sprintf("%s/%s", ingress.GetNamespace(), resource.Name))
				if err == nil && exists {
					dependencies = append(dependencies, kongServiceFacade.(client.Object))
				}
			}
		}
	}
	return dependencies
}

// resolveIngressDependenciesKongUpstreamPolicy resolves KongUpstreamPolicy dependencies for an Ingress object.
func resolveIngressDependenciesKongUpstreamPolicy(cache store.CacheStores, ingress *netv1.Ingress) []client.Object {
	upstreamPolicyName, ok := annotations.ExtractUpstreamPolicy(ingress.Annotations)
	if ok {
		upstreamPolicy, exists, err := cache.KongUpstreamPolicy.GetByKey(fmt.Sprintf("%s/%s", ingress.GetNamespace(), upstreamPolicyName))
		if err == nil && exists {
			return []client.Object{upstreamPolicy.(client.Object)}
		}
	}
	return nil
}

// resolveIngressDependenciesPlugin resolves KongPlugin and KongClusterPlugin dependencies for an Ingress object.
func resolveIngressDependenciesPlugin(cache store.CacheStores, ingress *netv1.Ingress) []client.Object {
	var dependencies []client.Object
	for _, pluginName := range annotations.ExtractKongPluginsFromAnnotations(ingress.Annotations) {
		if plugin, exists, err := cache.Plugin.GetByKey(
			fmt.Sprintf("%s/%s", ingress.GetNamespace(), pluginName),
		); err == nil && exists {
			dependencies = append(dependencies, plugin.(client.Object))
		}

		if plugin, exists, err := cache.ClusterPlugin.GetByKey(pluginName); err == nil && exists {
			dependencies = append(dependencies, plugin.(client.Object))
		}
	}
	return dependencies
}
