package store

import (
	"fmt"
	"sync"

	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured/unstructuredscheme"
	"k8s.io/apimachinery/pkg/runtime"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/tools/cache"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

//go:generate go run ../../hack/generators/cache-stores-snapshot

// CacheStores stores cache.Store for all Kinds of k8s objects that
// the Ingress Controller reads.
type CacheStores struct {
	// Core Kubernetes Stores
	IngressV1      cache.Store
	IngressClassV1 cache.Store
	Service        cache.Store
	Secret         cache.Store
	EndpointSlice  cache.Store

	// Gateway API Stores
	HTTPRoute      cache.Store
	UDPRoute       cache.Store
	TCPRoute       cache.Store
	TLSRoute       cache.Store
	GRPCRoute      cache.Store
	ReferenceGrant cache.Store
	Gateway        cache.Store

	// Kong Stores
	Plugin                         cache.Store
	ClusterPlugin                  cache.Store
	Consumer                       cache.Store
	ConsumerGroup                  cache.Store
	KongIngress                    cache.Store
	TCPIngress                     cache.Store
	UDPIngress                     cache.Store
	KongUpstreamPolicy             cache.Store
	IngressClassParametersV1alpha1 cache.Store
	KongServiceFacade              cache.Store
	KongVault                      cache.Store

	l *sync.RWMutex
}

// NewCacheStores is a convenience function for CacheStores to initialize all attributes with new cache stores.
func NewCacheStores() CacheStores {
	return CacheStores{
		// Core Kubernetes Stores
		IngressV1:      cache.NewStore(keyFunc),
		IngressClassV1: cache.NewStore(clusterResourceKeyFunc),
		Service:        cache.NewStore(keyFunc),
		Secret:         cache.NewStore(keyFunc),
		EndpointSlice:  cache.NewStore(keyFunc),
		// Gateway API Stores
		HTTPRoute:      cache.NewStore(keyFunc),
		UDPRoute:       cache.NewStore(keyFunc),
		TCPRoute:       cache.NewStore(keyFunc),
		TLSRoute:       cache.NewStore(keyFunc),
		GRPCRoute:      cache.NewStore(keyFunc),
		ReferenceGrant: cache.NewStore(keyFunc),
		Gateway:        cache.NewStore(keyFunc),
		// Kong Stores
		Plugin:                         cache.NewStore(keyFunc),
		ClusterPlugin:                  cache.NewStore(clusterResourceKeyFunc),
		Consumer:                       cache.NewStore(keyFunc),
		ConsumerGroup:                  cache.NewStore(keyFunc),
		KongIngress:                    cache.NewStore(keyFunc),
		TCPIngress:                     cache.NewStore(keyFunc),
		UDPIngress:                     cache.NewStore(keyFunc),
		KongUpstreamPolicy:             cache.NewStore(keyFunc),
		IngressClassParametersV1alpha1: cache.NewStore(keyFunc),
		KongServiceFacade:              cache.NewStore(keyFunc),
		KongVault:                      cache.NewStore(clusterResourceKeyFunc),

		l: &sync.RWMutex{},
	}
}

// NewCacheStoresFromObjYAML provides a new CacheStores object given any number of byte arrays containing
// YAML Kubernetes objects. An error is returned if any provided YAML was not a valid Kubernetes object.
func NewCacheStoresFromObjYAML(objs ...[]byte) (c CacheStores, err error) {
	kobjs := make([]runtime.Object, 0, len(objs))
	sr := serializer.NewYAMLSerializer(
		yamlserializer.DefaultMetaFactory,
		unstructuredscheme.NewUnstructuredCreator(),
		unstructuredscheme.NewUnstructuredObjectTyper(),
	)
	for _, yaml := range objs {
		kobj, _, decodeErr := sr.Decode(yaml, nil, nil)
		if err = decodeErr; err != nil {
			return
		}
		kobjs = append(kobjs, kobj)
	}
	return NewCacheStoresFromObjs(kobjs...)
}

// NewCacheStoresFromObjs provides a new CacheStores object given any number of Kubernetes
// objects that should be pre-populated. This function will sort objects into the appropriate
// sub-storage (e.g. IngressV1, TCPIngress, e.t.c.) but will produce an error if any of the
// input objects are erroneous or otherwise unusable as Kubernetes objects.
func NewCacheStoresFromObjs(objs ...runtime.Object) (CacheStores, error) {
	c := NewCacheStores()
	for _, obj := range objs {
		typedObj, err := mkObjFromGVK(obj.GetObjectKind().GroupVersionKind())
		if err != nil {
			return c, err
		}

		if err := convUnstructuredObj(obj, typedObj); err != nil {
			return c, err
		}

		if err := c.Add(typedObj); err != nil {
			return c, err
		}
	}
	return c, nil
}

// Get checks whether or not there's already some version of the provided object present in the cache.
func (c CacheStores) Get(obj runtime.Object) (item interface{}, exists bool, err error) {
	c.l.RLock()
	defer c.l.RUnlock()

	switch obj := obj.(type) {
	// ----------------------------------------------------------------------------
	// Kubernetes Core API Support
	// ----------------------------------------------------------------------------
	case *netv1.Ingress:
		return c.IngressV1.Get(obj)
	case *netv1.IngressClass:
		return c.IngressClassV1.Get(obj)
	case *corev1.Service:
		return c.Service.Get(obj)
	case *corev1.Secret:
		return c.Secret.Get(obj)
	case *discoveryv1.EndpointSlice:
		return c.EndpointSlice.Get(obj)
	// ----------------------------------------------------------------------------
	// Kubernetes Gateway API Support
	// ----------------------------------------------------------------------------
	case *gatewayapi.HTTPRoute:
		return c.HTTPRoute.Get(obj)
	case *gatewayapi.UDPRoute:
		return c.UDPRoute.Get(obj)
	case *gatewayapi.TCPRoute:
		return c.TCPRoute.Get(obj)
	case *gatewayapi.TLSRoute:
		return c.TLSRoute.Get(obj)
	case *gatewayapi.GRPCRoute:
		return c.GRPCRoute.Get(obj)
	case *gatewayapi.ReferenceGrant:
		return c.ReferenceGrant.Get(obj)
	case *gatewayapi.Gateway:
		return c.Gateway.Get(obj)
	// ----------------------------------------------------------------------------
	// Kong API Support
	// ----------------------------------------------------------------------------
	case *kongv1.KongPlugin:
		return c.Plugin.Get(obj)
	case *kongv1.KongClusterPlugin:
		return c.ClusterPlugin.Get(obj)
	case *kongv1.KongConsumer:
		return c.Consumer.Get(obj)
	case *kongv1beta1.KongConsumerGroup:
		return c.ConsumerGroup.Get(obj)
	case *kongv1.KongIngress:
		return c.KongIngress.Get(obj)
	case *kongv1beta1.TCPIngress:
		return c.TCPIngress.Get(obj)
	case *kongv1beta1.UDPIngress:
		return c.UDPIngress.Get(obj)
	case *kongv1beta1.KongUpstreamPolicy:
		return c.KongUpstreamPolicy.Get(obj)
	case *kongv1alpha1.IngressClassParameters:
		return c.IngressClassParametersV1alpha1.Get(obj)
	case *incubatorv1alpha1.KongServiceFacade:
		return c.KongServiceFacade.Get(obj)
	case *kongv1alpha1.KongVault:
		return c.KongVault.Get(obj)
	}
	return nil, false, fmt.Errorf("%T is not a supported cache object type", obj)
}

// Add stores a provided runtime.Object into the CacheStore if it's of a supported type.
// The CacheStore must be initialized (see NewCacheStores()) or this will panic.
func (c CacheStores) Add(obj runtime.Object) error {
	c.l.Lock()
	defer c.l.Unlock()

	switch obj := obj.(type) {
	// ----------------------------------------------------------------------------
	// Kubernetes Core API Support
	// ----------------------------------------------------------------------------
	case *netv1.Ingress:
		return c.IngressV1.Add(obj)
	case *netv1.IngressClass:
		return c.IngressClassV1.Add(obj)
	case *corev1.Service:
		return c.Service.Add(obj)
	case *corev1.Secret:
		return c.Secret.Add(obj)
	case *discoveryv1.EndpointSlice:
		return c.EndpointSlice.Add(obj)
	// ----------------------------------------------------------------------------
	// Kubernetes Gateway API Support
	// ----------------------------------------------------------------------------
	case *gatewayapi.HTTPRoute:
		return c.HTTPRoute.Add(obj)
	case *gatewayapi.UDPRoute:
		return c.UDPRoute.Add(obj)
	case *gatewayapi.TCPRoute:
		return c.TCPRoute.Add(obj)
	case *gatewayapi.TLSRoute:
		return c.TLSRoute.Add(obj)
	case *gatewayapi.GRPCRoute:
		return c.GRPCRoute.Add(obj)
	case *gatewayapi.ReferenceGrant:
		return c.ReferenceGrant.Add(obj)
	case *gatewayapi.Gateway:
		return c.Gateway.Add(obj)
	// ----------------------------------------------------------------------------
	// Kong API Support
	// ----------------------------------------------------------------------------
	case *kongv1.KongPlugin:
		return c.Plugin.Add(obj)
	case *kongv1.KongClusterPlugin:
		return c.ClusterPlugin.Add(obj)
	case *kongv1.KongConsumer:
		return c.Consumer.Add(obj)
	case *kongv1beta1.KongConsumerGroup:
		return c.ConsumerGroup.Add(obj)
	case *kongv1.KongIngress:
		return c.KongIngress.Add(obj)
	case *kongv1beta1.TCPIngress:
		return c.TCPIngress.Add(obj)
	case *kongv1beta1.UDPIngress:
		return c.UDPIngress.Add(obj)
	case *kongv1beta1.KongUpstreamPolicy:
		return c.KongUpstreamPolicy.Add(obj)
	case *kongv1alpha1.IngressClassParameters:
		return c.IngressClassParametersV1alpha1.Add(obj)
	case *incubatorv1alpha1.KongServiceFacade:
		return c.KongServiceFacade.Add(obj)
	case *kongv1alpha1.KongVault:
		return c.KongVault.Add(obj)
	default:
		return fmt.Errorf("cannot add unsupported kind %q to the store", obj.GetObjectKind().GroupVersionKind())
	}
}

// Delete removes a provided runtime.Object from the CacheStore if it's of a supported type.
// The CacheStore must be initialized (see NewCacheStores()) or this will panic.
func (c CacheStores) Delete(obj runtime.Object) error {
	c.l.Lock()
	defer c.l.Unlock()

	switch obj := obj.(type) {
	// ----------------------------------------------------------------------------
	// Kubernetes Core API Support
	// ----------------------------------------------------------------------------
	case *netv1.Ingress:
		return c.IngressV1.Delete(obj)
	case *netv1.IngressClass:
		return c.IngressClassV1.Delete(obj)
	case *corev1.Service:
		return c.Service.Delete(obj)
	case *corev1.Secret:
		return c.Secret.Delete(obj)
	case *discoveryv1.EndpointSlice:
		return c.EndpointSlice.Delete(obj)
	// ----------------------------------------------------------------------------
	// Kubernetes Gateway API Support
	// ----------------------------------------------------------------------------
	case *gatewayapi.HTTPRoute:
		return c.HTTPRoute.Delete(obj)
	case *gatewayapi.UDPRoute:
		return c.UDPRoute.Delete(obj)
	case *gatewayapi.TCPRoute:
		return c.TCPRoute.Delete(obj)
	case *gatewayapi.TLSRoute:
		return c.TLSRoute.Delete(obj)
	case *gatewayapi.GRPCRoute:
		return c.GRPCRoute.Delete(obj)
	case *gatewayapi.ReferenceGrant:
		return c.ReferenceGrant.Delete(obj)
	case *gatewayapi.Gateway:
		return c.Gateway.Delete(obj)
	// ----------------------------------------------------------------------------
	// Kong API Support
	// ----------------------------------------------------------------------------
	case *kongv1.KongPlugin:
		return c.Plugin.Delete(obj)
	case *kongv1.KongClusterPlugin:
		return c.ClusterPlugin.Delete(obj)
	case *kongv1.KongConsumer:
		return c.Consumer.Delete(obj)
	case *kongv1beta1.KongConsumerGroup:
		return c.ConsumerGroup.Delete(obj)
	case *kongv1.KongIngress:
		return c.KongIngress.Delete(obj)
	case *kongv1beta1.TCPIngress:
		return c.TCPIngress.Delete(obj)
	case *kongv1beta1.UDPIngress:
		return c.UDPIngress.Delete(obj)
	case *kongv1beta1.KongUpstreamPolicy:
		return c.KongUpstreamPolicy.Delete(obj)
	case *kongv1alpha1.IngressClassParameters:
		return c.IngressClassParametersV1alpha1.Delete(obj)
	case *incubatorv1alpha1.KongServiceFacade:
		return c.KongServiceFacade.Delete(obj)
	case *kongv1alpha1.KongVault:
		return c.KongVault.Delete(obj)
	default:
		return fmt.Errorf("cannot delete unsupported kind %q from the store", obj.GetObjectKind().GroupVersionKind())
	}
}
