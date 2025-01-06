package fallback

import (
	"fmt"
	"slices"

	"github.com/samber/lo"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// resolveHTTPRouteDependencies resolves potential dependencies for a given HTTPRoute object:
// - Service
// - KongPlugin
// - KongClusterPlugin.
func resolveHTTPRouteDependencies(cache store.CacheStores, route *gatewayapi.HTTPRoute) []client.Object {
	return slices.Concat(
		resolveGatewayAPIRouteDependenciesBackendRefs(cache, route, getHTTPRouteBackendRefs(route)),
		resolveObjectDependenciesPlugin(cache, route),
	)
}

// resolveTCPRouteDependencies resolves potential dependencies for a given TCPRoute object:
// - Service
// - KongPlugin
// - KongClusterPlugin.
func resolveTCPRouteDependencies(cache store.CacheStores, route *gatewayapi.TCPRoute) []client.Object {
	return slices.Concat(
		resolveGatewayAPIRouteDependenciesBackendRefs(cache, route, getTCPRouteBackendRefs(route)),
		resolveObjectDependenciesPlugin(cache, route),
	)
}

// resolveUDPRouteDependencies resolves potential dependencies for a given UDPRoute object:
// - Service
// - KongPlugin
// - KongClusterPlugin.
func resolveUDPRouteDependencies(cache store.CacheStores, route *gatewayapi.UDPRoute) []client.Object {
	return slices.Concat(
		resolveGatewayAPIRouteDependenciesBackendRefs(cache, route, getUDPRouteBackendRefs(route)),
		resolveObjectDependenciesPlugin(cache, route),
	)
}

// resolveTLSRouteDependencies resolves potential dependencies for a given TLSRoute object:
// - Service
// - KongPlugin
// - KongClusterPlugin.
func resolveTLSRouteDependencies(cache store.CacheStores, route *gatewayapi.TLSRoute) []client.Object {
	return slices.Concat(
		resolveGatewayAPIRouteDependenciesBackendRefs(cache, route, getTLSRouteBackendRefs(route)),
		resolveObjectDependenciesPlugin(cache, route),
	)
}

// resolveGRPCRouteDependencies resolves potential dependencies for a given GRPCRoute object:
// - Service
// - KongPlugin
// - KongClusterPlugin.
func resolveGRPCRouteDependencies(cache store.CacheStores, route *gatewayapi.GRPCRoute) []client.Object {
	return slices.Concat(
		resolveGatewayAPIRouteDependenciesBackendRefs(cache, route, getGRPCRouteBackendRefs(route)),
		resolveObjectDependenciesPlugin(cache, route),
	)
}

// gatewayAPIRoute is an interface that represents a GatewayAPI Route object.
type gatewayAPIRoute interface {
	client.Object
	*gatewayapi.HTTPRoute | *gatewayapi.TCPRoute | *gatewayapi.UDPRoute | *gatewayapi.TLSRoute | *gatewayapi.GRPCRoute
}

// resolveGatewayAPIRouteDependenciesBackendRefs resolves backend references for a given gatewayAPIRoute object.
func resolveGatewayAPIRouteDependenciesBackendRefs[T gatewayAPIRoute](cache store.CacheStores, route T, backendRefs []gatewayapi.BackendRef) []client.Object {
	var dependencies []client.Object
	for _, backendRef := range backendRefs {
		if !util.IsBackendRefGroupKindSupported(backendRef.Group, backendRef.Kind) {
			continue
		}
		ns := route.GetNamespace()
		if backendRef.Namespace != nil {
			ns = string(*backendRef.Namespace)
		}
		ingressClass, exists, err := cache.Service.GetByKey(fmt.Sprintf("%s/%s", ns, backendRef.Name))
		if err == nil && exists {
			dependencies = append(dependencies, ingressClass.(client.Object))
		}
	}
	return dependencies
}

func getHTTPRouteBackendRefs(route *gatewayapi.HTTPRoute) []gatewayapi.BackendRef {
	var backendRefs []gatewayapi.BackendRef
	for _, rule := range route.Spec.Rules {
		backendRefs = append(backendRefs, lo.Map(rule.BackendRefs, func(b gatewayapi.HTTPBackendRef, _ int) gatewayapi.BackendRef {
			return b.BackendRef
		})...)
	}
	return backendRefs
}

func getTCPRouteBackendRefs(route *gatewayapi.TCPRoute) []gatewayapi.BackendRef {
	var backendRefs []gatewayapi.BackendRef
	for _, rule := range route.Spec.Rules {
		backendRefs = append(backendRefs, rule.BackendRefs...)
	}
	return backendRefs
}

func getUDPRouteBackendRefs(route *gatewayapi.UDPRoute) []gatewayapi.BackendRef {
	var backendRefs []gatewayapi.BackendRef
	for _, rule := range route.Spec.Rules {
		backendRefs = append(backendRefs, rule.BackendRefs...)
	}
	return backendRefs
}

func getTLSRouteBackendRefs(route *gatewayapi.TLSRoute) []gatewayapi.BackendRef {
	var backendRefs []gatewayapi.BackendRef
	for _, rule := range route.Spec.Rules {
		backendRefs = append(backendRefs, rule.BackendRefs...)
	}
	return backendRefs
}

func getGRPCRouteBackendRefs(route *gatewayapi.GRPCRoute) []gatewayapi.BackendRef {
	var backendRefs []gatewayapi.BackendRef
	for _, rule := range route.Spec.Rules {
		backendRefs = append(backendRefs, lo.Map(rule.BackendRefs, func(b gatewayapi.GRPCBackendRef, _ int) gatewayapi.BackendRef {
			return b.BackendRef
		})...)
	}
	return backendRefs
}
