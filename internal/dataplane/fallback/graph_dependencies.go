package fallback

import (
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

// ResolveDependencies resolves dependencies for a given object. Dependencies are all objects referenced by the
// given object. For example, an Ingress object might refer to an IngressClass, Services, Plugins, etc.
func ResolveDependencies(cache store.CacheStores, obj client.Object) []client.Object {
	// TODO: Implement dependency resolution for all types below.
	// https://github.com/Kong/kubernetes-ingress-controller/issues/5929
	switch obj := obj.(type) {
	case *netv1.Ingress:
		return resolveIngressDependencies(cache, obj)
	case *gatewayapi.HTTPRoute:
	case *gatewayapi.TLSRoute:
	case *gatewayapi.TCPRoute:
	case *gatewayapi.UDPRoute:
	case *gatewayapi.GRPCRoute:
	case *kongv1.KongPlugin:
	case *kongv1.KongClusterPlugin:
	case *kongv1beta1.UDPIngress:
	case *kongv1beta1.TCPIngress:
	case *incubatorv1alpha1.KongServiceFacade:
	}

	// If there's no dependency resolution for the given object type, we assume there are no dependencies.
	return nil
}
