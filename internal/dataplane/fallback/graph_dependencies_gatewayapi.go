package fallback

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/6026
func resolveHTTPRouteDependencies(_ store.CacheStores, _ *gatewayapi.HTTPRoute) []client.Object {
	return nil
}

// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/6026
func resolveTLSRouteDependencies(_ store.CacheStores, _ *gatewayapi.TLSRoute) []client.Object {
	return nil
}

// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/6026
func resolveTCPRouteDependencies(_ store.CacheStores, _ *gatewayapi.TCPRoute) []client.Object {
	return nil
}

// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/6026
func resolveUDPRouteDependencies(_ store.CacheStores, _ *gatewayapi.UDPRoute) []client.Object {
	return nil
}

// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/6026
func resolveGRPCRouteDependencies(_ store.CacheStores, _ *gatewayapi.GRPCRoute) []client.Object {
	return nil
}
