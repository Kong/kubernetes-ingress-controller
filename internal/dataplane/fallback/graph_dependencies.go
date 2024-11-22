package fallback

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// ResolveDependencies resolves dependencies for a given object. Dependencies are all objects referenced by the
// given object. For example, an Ingress object might refer to an IngressClass, Services, Plugins, etc.
// Every supported object type should explicitly have a case in this function.
func ResolveDependencies(cache store.CacheStores, obj client.Object) ([]client.Object, error) {
	switch obj := obj.(type) {
	// Standard Kubernetes objects.
	case *corev1.Service:
		return resolveServiceDependencies(cache, obj), nil
	case *netv1.Ingress:
		return resolveIngressDependencies(cache, obj), nil
	// Gateway API objects.
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
	// Kong specific objects.
	case *kongv1.KongPlugin:
		return resolveKongPluginDependencies(cache, obj), nil
	case *kongv1.KongClusterPlugin:
		return resolveKongClusterPluginDependencies(cache, obj), nil
	case *kongv1.KongConsumer:
		return resolveKongConsumerDependencies(cache, obj), nil
	case *kongv1beta1.KongConsumerGroup:
		return resolveKongConsumerGroupDependencies(cache, obj), nil
	case *kongv1beta1.UDPIngress:
		return resolveUDPIngressDependencies(cache, obj), nil
	case *kongv1beta1.TCPIngress:
		return resolveTCPIngressDependencies(cache, obj), nil
	case *incubatorv1alpha1.KongServiceFacade:
		return resolveKongServiceFacadeDependencies(cache, obj), nil
	case *kongv1alpha1.KongCustomEntity:
		return resolveKongCustomEntityDependencies(cache, obj), nil
	// Object types that have no dependencies.
	case *netv1.IngressClass,
		*corev1.Secret,
		*discoveryv1.EndpointSlice,
		*gatewayapi.ReferenceGrant,
		*gatewayapi.Gateway,
		*gatewayapi.BackendTLSPolicy,
		*kongv1.KongIngress,
		*kongv1beta1.KongUpstreamPolicy,
		*kongv1alpha1.IngressClassParameters,
		*kongv1alpha1.KongVault:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported object type: %T", obj)
	}
}
