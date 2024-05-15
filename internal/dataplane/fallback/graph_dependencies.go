package fallback

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

// ResolveDependencies resolves dependencies for a given object. Dependencies are all objects referenced by the
// given object. For example, an Ingress object might refer to an IngressClass, Services, Plugins, etc.
// Every supported object type should explicitly have a case in this function.
func ResolveDependencies(cache store.CacheStores, obj client.Object) ([]client.Object, error) {
	switch obj := obj.(type) {
	case *netv1.Ingress:
		return resolveIngressDependencies(cache, obj), nil
	case *gatewayapi.HTTPRoute:
		return resolveHTTPRouteDependencies(cache, obj), nil
	case *gatewayapi.TLSRoute:
		return resolveTLSRouteDependencies(cache, obj), nil
	case *gatewayapi.TCPRoute:
		return resolveTCPRouteDependencies(cache, obj), nil
	case *gatewayapi.UDPRoute:
		return resolveUDPRouteDependencies(cache, obj), nil
	case *gatewayapi.GRPCRoute:
		return resolveGRPCRouteDependencies(cache, obj), nil
	case *kongv1.KongPlugin:
		return resolveKongPluginDependencies(cache, obj), nil
	case *kongv1.KongClusterPlugin:
		return resolveKongClusterPluginDependencies(cache, obj), nil
	case *kongv1beta1.UDPIngress:
		return resolveUDPIngressDependencies(cache, obj), nil
	case *kongv1beta1.TCPIngress:
		return resolveTCPIngressDependencies(cache, obj), nil
	case *incubatorv1alpha1.KongServiceFacade:
		return resolveKongServiceFacadeDependencies(cache, obj), nil
	case *netv1.IngressClass, // Object types that have no dependencies.
		*corev1.Service,
		*corev1.Secret,
		*discoveryv1.EndpointSlice,
		*gatewayapi.ReferenceGrant,
		*gatewayapi.Gateway,
		*kongv1.KongConsumer,
		*kongv1beta1.KongConsumerGroup,
		*kongv1.KongIngress,
		*kongv1beta1.KongUpstreamPolicy,
		*kongv1alpha1.IngressClassParameters,
		*kongv1alpha1.KongVault:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported object type: %T", obj)
	}
}
