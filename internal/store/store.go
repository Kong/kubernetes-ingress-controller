/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package store

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured/unstructuredscheme"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/tools/cache"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	ctrlutils "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

const (
	caCertKey = "konghq.com/ca-cert"
	// IngressClassKongController is the string used for the Controller field of a recognized IngressClass.
	IngressClassKongController = "ingress-controllers.konghq.com/kong"
)

// ErrNotFound error is returned when a lookup results in no resource.
// This type is meant to be used for error handling using `errors.As()`.
type ErrNotFound struct {
	Message string
}

func (e ErrNotFound) Error() string {
	if e.Message == "" {
		return "not found"
	}
	return e.Message
}

// Storer is the interface that wraps the required methods to gather information
// about ingresses, services, secrets and ingress annotations.
type Storer interface {
	GetSecret(namespace, name string) (*corev1.Secret, error)
	GetService(namespace, name string) (*corev1.Service, error)
	GetEndpointSlicesForService(namespace, name string) ([]*discoveryv1.EndpointSlice, error)
	GetKongIngress(namespace, name string) (*kongv1.KongIngress, error)
	GetKongPlugin(namespace, name string) (*kongv1.KongPlugin, error)
	GetKongClusterPlugin(name string) (*kongv1.KongClusterPlugin, error)
	GetKongConsumer(namespace, name string) (*kongv1.KongConsumer, error)
	GetKongConsumerGroup(namespace, name string) (*kongv1beta1.KongConsumerGroup, error)
	GetIngressClassName() string
	GetIngressClassV1(name string) (*netv1.IngressClass, error)
	GetIngressClassParametersV1Alpha1(ingressClass *netv1.IngressClass) (*kongv1alpha1.IngressClassParameters, error)
	GetGateway(namespace string, name string) (*gatewayv1.Gateway, error)

	ListIngressesV1() []*netv1.Ingress
	ListIngressClassesV1() []*netv1.IngressClass
	ListIngressClassParametersV1Alpha1() []*kongv1alpha1.IngressClassParameters
	ListHTTPRoutes() ([]*gatewayv1.HTTPRoute, error)
	ListUDPRoutes() ([]*gatewayv1alpha2.UDPRoute, error)
	ListTCPRoutes() ([]*gatewayv1alpha2.TCPRoute, error)
	ListTLSRoutes() ([]*gatewayv1alpha2.TLSRoute, error)
	ListGRPCRoutes() ([]*gatewayv1alpha2.GRPCRoute, error)
	ListReferenceGrants() ([]*gatewayv1beta1.ReferenceGrant, error)
	ListGateways() ([]*gatewayv1.Gateway, error)
	ListTCPIngresses() ([]*kongv1beta1.TCPIngress, error)
	ListUDPIngresses() ([]*kongv1beta1.UDPIngress, error)
	ListKnativeIngresses() ([]*knative.Ingress, error)
	ListGlobalKongPlugins() ([]*kongv1.KongPlugin, error)
	ListGlobalKongClusterPlugins() ([]*kongv1.KongClusterPlugin, error)
	ListKongPlugins() []*kongv1.KongPlugin
	ListKongClusterPlugins() []*kongv1.KongClusterPlugin
	ListKongConsumers() []*kongv1.KongConsumer
	ListKongConsumerGroups() []*kongv1beta1.KongConsumerGroup
	ListCACerts() ([]*corev1.Secret, error)
}

// Store implements Storer and can be used to list Ingress, Services
// and other resources from k8s APIserver. The backing stores should
// be synced and updated by the caller.
// It is ingressClass filter aware.
type Store struct {
	stores CacheStores

	ingressClass         string
	ingressClassMatching annotations.ClassMatching

	isValidIngressClass   func(objectMeta *metav1.ObjectMeta, annotation string, handling annotations.ClassMatching) bool
	isValidIngressV1Class func(ingress *netv1.Ingress, handling annotations.ClassMatching) bool

	logger logrus.FieldLogger
}

var _ Storer = Store{}

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
	IngressClassParametersV1alpha1 cache.Store

	// Knative Stores
	KnativeIngress cache.Store

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
		IngressClassParametersV1alpha1: cache.NewStore(keyFunc),
		// Knative Stores
		KnativeIngress: cache.NewStore(keyFunc),

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
	case *gatewayv1.HTTPRoute:
		return c.HTTPRoute.Get(obj)
	case *gatewayv1alpha2.UDPRoute:
		return c.UDPRoute.Get(obj)
	case *gatewayv1alpha2.TCPRoute:
		return c.TCPRoute.Get(obj)
	case *gatewayv1alpha2.TLSRoute:
		return c.TLSRoute.Get(obj)
	case *gatewayv1alpha2.GRPCRoute:
		return c.GRPCRoute.Get(obj)
	case *gatewayv1beta1.ReferenceGrant:
		return c.ReferenceGrant.Get(obj)
	case *gatewayv1.Gateway:
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
	case *kongv1alpha1.IngressClassParameters:
		return c.IngressClassParametersV1alpha1.Get(obj)
	// ----------------------------------------------------------------------------
	// 3rd Party API Support
	// ----------------------------------------------------------------------------
	case *knative.Ingress:
		return c.KnativeIngress.Get(obj)
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
	case *gatewayv1.HTTPRoute:
		return c.HTTPRoute.Add(obj)
	case *gatewayv1alpha2.UDPRoute:
		return c.UDPRoute.Add(obj)
	case *gatewayv1alpha2.TCPRoute:
		return c.TCPRoute.Add(obj)
	case *gatewayv1alpha2.TLSRoute:
		return c.TLSRoute.Add(obj)
	case *gatewayv1alpha2.GRPCRoute:
		return c.GRPCRoute.Add(obj)
	case *gatewayv1beta1.ReferenceGrant:
		return c.ReferenceGrant.Add(obj)
	case *gatewayv1.Gateway:
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
	case *kongv1alpha1.IngressClassParameters:
		return c.IngressClassParametersV1alpha1.Add(obj)
	// ----------------------------------------------------------------------------
	// 3rd Party API Support
	// ----------------------------------------------------------------------------
	case *knative.Ingress:
		return c.KnativeIngress.Add(obj)
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
	case *gatewayv1.HTTPRoute:
		return c.HTTPRoute.Delete(obj)
	case *gatewayv1alpha2.UDPRoute:
		return c.UDPRoute.Delete(obj)
	case *gatewayv1alpha2.TCPRoute:
		return c.TCPRoute.Delete(obj)
	case *gatewayv1alpha2.TLSRoute:
		return c.TLSRoute.Delete(obj)
	case *gatewayv1alpha2.GRPCRoute:
		return c.GRPCRoute.Delete(obj)
	case *gatewayv1beta1.ReferenceGrant:
		return c.ReferenceGrant.Delete(obj)
	case *gatewayv1.Gateway:
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
	case *kongv1alpha1.IngressClassParameters:
		return c.IngressClassParametersV1alpha1.Delete(obj)
	// ----------------------------------------------------------------------------
	// 3rd Party API Support
	// ----------------------------------------------------------------------------
	case *knative.Ingress:
		return c.KnativeIngress.Delete(obj)
	default:
		return fmt.Errorf("cannot delete unsupported kind %q from the store", obj.GetObjectKind().GroupVersionKind())
	}
}

// New creates a new object store to be used in the ingress controller.
func New(cs CacheStores, ingressClass string, logger logrus.FieldLogger) Storer {
	return Store{
		stores:                cs,
		ingressClass:          ingressClass,
		ingressClassMatching:  annotations.ExactClassMatch,
		isValidIngressClass:   annotations.IngressClassValidatorFuncFromObjectMeta(ingressClass),
		isValidIngressV1Class: annotations.IngressClassValidatorFuncFromV1Ingress(ingressClass),
		logger:                logger,
	}
}

// GetSecret returns a Secret using the namespace and name as key.
func (s Store) GetSecret(namespace, name string) (*corev1.Secret, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	secret, exists, err := s.stores.Secret.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("Secret %v not found", key)}
	}
	return secret.(*corev1.Secret), nil
}

// GetService returns a Service using the namespace and name as key.
func (s Store) GetService(namespace, name string) (*corev1.Service, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	service, exists, err := s.stores.Service.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("Service %v not found", key)}
	}
	return service.(*corev1.Service), nil
}

// ListIngressesV1 returns the list of Ingresses in the Ingress v1 store.
func (s Store) ListIngressesV1() []*netv1.Ingress {
	// filter ingress rules
	var ingresses []*netv1.Ingress
	for _, item := range s.stores.IngressV1.List() {
		ing, ok := item.(*netv1.Ingress)
		if !ok {
			s.logger.Warnf("listIngressesV1: dropping object of unexpected type: %#v", item)
			continue
		}
		if ing.ObjectMeta.GetAnnotations()[annotations.IngressClassKey] != "" {
			if !s.isValidIngressClass(&ing.ObjectMeta, annotations.IngressClassKey, s.ingressClassMatching) {
				continue
			}
		} else if ing.Spec.IngressClassName != nil {
			if !s.isValidIngressV1Class(ing, s.ingressClassMatching) {
				continue
			}
		} else {
			class, err := s.GetIngressClassV1(s.ingressClass)
			if err != nil {
				s.logger.Debugf("IngressClass %s not found", s.ingressClass)
				continue
			}
			if !ctrlutils.IsDefaultIngressClass(class) {
				continue
			}
		}
		ingresses = append(ingresses, ing)
	}

	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name)) < 0
	})

	return ingresses
}

// ListIngressClassesV1 returns the list of Ingresses in the Ingress v1 store.
func (s Store) ListIngressClassesV1() []*netv1.IngressClass {
	// filter ingress rules
	var classes []*netv1.IngressClass
	for _, item := range s.stores.IngressClassV1.List() {
		class, ok := item.(*netv1.IngressClass)
		if !ok {
			s.logger.Warnf("listIngressClassesV1: dropping object of unexpected type: %#v", item)
			continue
		}
		if class.Spec.Controller != IngressClassKongController {
			continue
		}
		classes = append(classes, class)
	}

	sort.SliceStable(classes, func(i, j int) bool {
		return strings.Compare(classes[i].Name, classes[j].Name) < 0
	})

	return classes
}

// ListIngressClassParametersV1Alpha1 returns the list of IngressClassParameters in the Ingress v1alpha1 store.
func (s Store) ListIngressClassParametersV1Alpha1() []*kongv1alpha1.IngressClassParameters {
	var classParams []*kongv1alpha1.IngressClassParameters
	for _, item := range s.stores.IngressClassParametersV1alpha1.List() {
		classParam, ok := item.(*kongv1alpha1.IngressClassParameters)
		if !ok {
			s.logger.Warnf("listIngressClassParametersV1alpha1: dropping object of unexpected type: %#v", item)
			continue
		}
		classParams = append(classParams, classParam)
	}

	sort.SliceStable(classParams, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", classParams[i].Namespace, classParams[i].Name),
			fmt.Sprintf("%s/%s", classParams[j].Namespace, classParams[j].Name),
		) < 0
	})

	return classParams
}

// ListHTTPRoutes returns the list of HTTPRoutes in the HTTPRoute cache store.
func (s Store) ListHTTPRoutes() ([]*gatewayv1.HTTPRoute, error) {
	var httproutes []*gatewayv1.HTTPRoute
	if err := cache.ListAll(s.stores.HTTPRoute, labels.NewSelector(),
		func(ob interface{}) {
			httproute, ok := ob.(*gatewayv1.HTTPRoute)
			if ok {
				httproutes = append(httproutes, httproute)
			}
		},
	); err != nil {
		return nil, err
	}
	return httproutes, nil
}

// ListUDPRoutes returns the list of UDPRoutes in the UDPRoute cache store.
func (s Store) ListUDPRoutes() ([]*gatewayv1alpha2.UDPRoute, error) {
	var udproutes []*gatewayv1alpha2.UDPRoute
	if err := cache.ListAll(s.stores.UDPRoute, labels.NewSelector(),
		func(ob interface{}) {
			udproute, ok := ob.(*gatewayv1alpha2.UDPRoute)
			if ok {
				udproutes = append(udproutes, udproute)
			}
		},
	); err != nil {
		return nil, err
	}
	return udproutes, nil
}

// ListTCPRoutes returns the list of TCPRoutes in the TCPRoute cache store.
func (s Store) ListTCPRoutes() ([]*gatewayv1alpha2.TCPRoute, error) {
	var tcproutes []*gatewayv1alpha2.TCPRoute
	if err := cache.ListAll(s.stores.TCPRoute, labels.NewSelector(),
		func(ob interface{}) {
			tcproute, ok := ob.(*gatewayv1alpha2.TCPRoute)
			if ok {
				tcproutes = append(tcproutes, tcproute)
			}
		},
	); err != nil {
		return nil, err
	}
	return tcproutes, nil
}

// ListTLSRoutes returns the list of TLSRoutes in the TLSRoute cache store.
func (s Store) ListTLSRoutes() ([]*gatewayv1alpha2.TLSRoute, error) {
	var tlsroutes []*gatewayv1alpha2.TLSRoute
	if err := cache.ListAll(s.stores.TLSRoute, labels.NewSelector(),
		func(ob interface{}) {
			tlsroute, ok := ob.(*gatewayv1alpha2.TLSRoute)
			if ok {
				tlsroutes = append(tlsroutes, tlsroute)
			}
		},
	); err != nil {
		return nil, err
	}
	return tlsroutes, nil
}

// ListGRPCRoutes returns the list of GRPCRoutes in the GRPCRoute cache store.
func (s Store) ListGRPCRoutes() ([]*gatewayv1alpha2.GRPCRoute, error) {
	var grpcroutes []*gatewayv1alpha2.GRPCRoute
	if err := cache.ListAll(s.stores.GRPCRoute, labels.NewSelector(),
		func(ob interface{}) {
			tlsroute, ok := ob.(*gatewayv1alpha2.GRPCRoute)
			if ok {
				grpcroutes = append(grpcroutes, tlsroute)
			}
		},
	); err != nil {
		return nil, err
	}
	return grpcroutes, nil
}

// ListReferenceGrants returns the list of ReferenceGrants in the ReferenceGrant cache store.
func (s Store) ListReferenceGrants() ([]*gatewayv1beta1.ReferenceGrant, error) {
	var grants []*gatewayv1beta1.ReferenceGrant
	if err := cache.ListAll(s.stores.ReferenceGrant, labels.NewSelector(),
		func(ob interface{}) {
			grant, ok := ob.(*gatewayv1beta1.ReferenceGrant)
			if ok {
				grants = append(grants, grant)
			}
		},
	); err != nil {
		return nil, err
	}
	return grants, nil
}

// ListGateways returns the list of Gateways in the Gateway cache store.
func (s Store) ListGateways() ([]*gatewayv1.Gateway, error) {
	var gateways []*gatewayv1.Gateway
	if err := cache.ListAll(s.stores.Gateway, labels.NewSelector(),
		func(ob interface{}) {
			gw, ok := ob.(*gatewayv1.Gateway)
			if ok {
				gateways = append(gateways, gw)
			}
		},
	); err != nil {
		return nil, err
	}
	return gateways, nil
}

// ListTCPIngresses returns the list of TCP Ingresses from
// configuration.konghq.com group.
func (s Store) ListTCPIngresses() ([]*kongv1beta1.TCPIngress, error) {
	var ingresses []*kongv1beta1.TCPIngress
	err := cache.ListAll(s.stores.TCPIngress, labels.NewSelector(),
		func(ob interface{}) {
			ing, ok := ob.(*kongv1beta1.TCPIngress)
			if ok && s.isValidIngressClass(&ing.ObjectMeta, annotations.IngressClassKey, s.getIngressClassHandling()) {
				ingresses = append(ingresses, ing)
			}
		})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name)) < 0
	})
	return ingresses, nil
}

// ListUDPIngresses returns the list of UDP Ingresses.
func (s Store) ListUDPIngresses() ([]*kongv1beta1.UDPIngress, error) {
	ingresses := []*kongv1beta1.UDPIngress{}
	if s.stores.UDPIngress == nil {
		// older versions of the KIC do not support UDPIngress so short circuit to maintain support with them
		return ingresses, nil
	}

	err := cache.ListAll(s.stores.UDPIngress, labels.NewSelector(),
		func(ob interface{}) {
			ing, ok := ob.(*kongv1beta1.UDPIngress)
			if ok && s.isValidIngressClass(&ing.ObjectMeta, annotations.IngressClassKey, s.getIngressClassHandling()) {
				ingresses = append(ingresses, ing)
			}
		})
	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name)) < 0
	})
	return ingresses, err
}

// ListKnativeIngresses returns the list of Knative Ingresses from
// ingresses.networking.internal.knative.dev group.
func (s Store) ListKnativeIngresses() ([]*knative.Ingress, error) {
	var ingresses []*knative.Ingress
	if s.stores.KnativeIngress == nil {
		return ingresses, nil
	}

	err := cache.ListAll(
		s.stores.KnativeIngress,
		labels.NewSelector(),
		func(ob interface{}) {
			ing, ok := ob.(*knative.Ingress)
			if ok {
				handlingClass := s.getIngressClassHandling()
				if s.isValidIngressClass(&ing.ObjectMeta, annotations.KnativeIngressClassKey, handlingClass) ||
					s.isValidIngressClass(&ing.ObjectMeta, annotations.KnativeIngressClassDeprecatedKey, handlingClass) {
					ingresses = append(ingresses, ing)
				}
			}
		})
	if err != nil {
		return nil, err
	}

	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name)) < 0
	})
	return ingresses, nil
}

// GetEndpointSlicesForService returns all EndpointSlices for service
// 'namespace/name' inside K8s.
func (s Store) GetEndpointSlicesForService(namespace, name string) ([]*discoveryv1.EndpointSlice, error) {
	// EndpointSlices are tied to a Service via a label.
	req, err := labels.NewRequirement(discoveryv1.LabelServiceName, selection.Equals, []string{name})
	if err != nil {
		return nil, err
	}
	var endpointSlices []*discoveryv1.EndpointSlice
	if err := cache.ListAll(
		s.stores.EndpointSlice, labels.NewSelector().Add(*req),
		func(obj interface{}) {
			// Ensure the EndpointSlice is for the Service from the requested namespace.
			if eps, ok := obj.(*discoveryv1.EndpointSlice); ok && eps.Namespace == namespace {
				endpointSlices = append(endpointSlices, eps)
			}
		},
	); err != nil {
		return nil, err
	}
	if len(endpointSlices) == 0 {
		return nil, ErrNotFound{fmt.Sprintf("EndpointSlices for Service %s/%s not found", namespace, name)}
	}
	return endpointSlices, nil
}

// GetKongPlugin returns the 'name' KongPlugin resource in namespace.
func (s Store) GetKongPlugin(namespace, name string) (*kongv1.KongPlugin, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.stores.Plugin.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("KongPlugin %v not found", key)}
	}
	return p.(*kongv1.KongPlugin), nil
}

// GetKongClusterPlugin returns the 'name' KongClusterPlugin resource.
func (s Store) GetKongClusterPlugin(name string) (*kongv1.KongClusterPlugin, error) {
	p, exists, err := s.stores.ClusterPlugin.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("KongClusterPlugin %v not found", name)}
	}
	return p.(*kongv1.KongClusterPlugin), nil
}

// GetKongIngress returns the 'name' KongIngress resource in namespace.
func (s Store) GetKongIngress(namespace, name string) (*kongv1.KongIngress, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.stores.KongIngress.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("KongIngress %v not found", name)}
	}
	return p.(*kongv1.KongIngress), nil
}

// GetKongConsumer returns the 'name' KongConsumer resource in namespace.
func (s Store) GetKongConsumer(namespace, name string) (*kongv1.KongConsumer, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.stores.Consumer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("KongConsumer %v not found", key)}
	}
	return p.(*kongv1.KongConsumer), nil
}

// GetKongConsumerGroup returns the 'name' KongConsumerGroup resource in namespace.
func (s Store) GetKongConsumerGroup(namespace, name string) (*kongv1beta1.KongConsumerGroup, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.stores.ConsumerGroup.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("KongConsumerGroup %v not found", key)}
	}
	return p.(*kongv1beta1.KongConsumerGroup), nil
}

func (s Store) GetIngressClassName() string {
	return s.ingressClass
}

// GetIngressClassV1 returns the 'name' IngressClass resource.
func (s Store) GetIngressClassV1(name string) (*netv1.IngressClass, error) {
	p, exists, err := s.stores.IngressClassV1.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("IngressClass %v not found", name)}
	}
	return p.(*netv1.IngressClass), nil
}

// GetIngressClassParametersV1Alpha1 returns IngressClassParameters for provided
// IngressClass.
func (s Store) GetIngressClassParametersV1Alpha1(ingressClass *netv1.IngressClass) (*kongv1alpha1.IngressClassParameters, error) {
	if ingressClass == nil {
		return nil, fmt.Errorf("provided IngressClass is nil")
	}

	if ingressClass.Spec.Parameters == nil {
		return &kongv1alpha1.IngressClassParameters{}, nil
	}

	if ingressClass.Spec.Parameters.APIGroup == nil ||
		*ingressClass.Spec.Parameters.APIGroup != kongv1alpha1.GroupVersion.Group {
		return nil, fmt.Errorf(
			"IngressClass %s should reference parameters in apiGroup:%s",
			ingressClass.Name,
			kongv1alpha1.GroupVersion.Group,
		)
	}

	if ingressClass.Spec.Parameters.Kind != kongv1alpha1.IngressClassParametersKind {
		return nil, fmt.Errorf(
			"IngressClass %s should reference parameters with kind:%s",
			ingressClass.Name,
			kongv1alpha1.IngressClassParametersKind,
		)
	}

	if ingressClass.Spec.Parameters.Scope == nil || ingressClass.Spec.Parameters.Namespace == nil {
		return nil, fmt.Errorf("IngressClass %s should reference namespaced parameters", ingressClass.Name)
	}

	key := fmt.Sprintf("%v/%v", *ingressClass.Spec.Parameters.Namespace, ingressClass.Spec.Parameters.Name)
	params, exists, err := s.stores.IngressClassParametersV1alpha1.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("IngressClassParameters %v not found", ingressClass.Spec.Parameters.Name)}
	}
	return params.(*kongv1alpha1.IngressClassParameters), nil
}

// GetGateway returns gateway resource having specified namespace and name.
func (s Store) GetGateway(namespace string, name string) (*gatewayv1.Gateway, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	obj, exists, err := s.stores.Gateway.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("Gateway %v not found", name)}
	}
	return obj.(*gatewayv1.Gateway), nil
}

// ListKongConsumers returns all KongConsumers filtered by the ingress.class
// annotation.
func (s Store) ListKongConsumers() []*kongv1.KongConsumer {
	var consumers []*kongv1.KongConsumer
	for _, item := range s.stores.Consumer.List() {
		c, ok := item.(*kongv1.KongConsumer)
		if ok && s.isValidIngressClass(&c.ObjectMeta, annotations.IngressClassKey, s.getIngressClassHandling()) {
			consumers = append(consumers, c)
		}
	}

	return consumers
}

// ListKongConsumerGroups returns all KongConsumerGroups filtered by the ingress.class
// annotation.
func (s Store) ListKongConsumerGroups() []*kongv1beta1.KongConsumerGroup {
	var consumerGroups []*kongv1beta1.KongConsumerGroup
	for _, item := range s.stores.ConsumerGroup.List() {
		c, ok := item.(*kongv1beta1.KongConsumerGroup)
		if ok && s.isValidIngressClass(&c.ObjectMeta, annotations.IngressClassKey, s.getIngressClassHandling()) {
			consumerGroups = append(consumerGroups, c)
		}
	}

	return consumerGroups
}

// ListGlobalKongPlugins returns all KongPlugin resources
// filtered by the ingress.class annotation and with the
// label global:"true".
// Support for these global namespaced KongPlugins was removed in 0.10.0
// This function remains only to provide warnings to users with old configuration.
func (s Store) ListGlobalKongPlugins() ([]*kongv1.KongPlugin, error) {
	var plugins []*kongv1.KongPlugin
	req, err := labels.NewRequirement("global", selection.Equals, []string{"true"})
	if err != nil {
		return nil, err
	}
	err = cache.ListAll(s.stores.Plugin,
		labels.NewSelector().Add(*req),
		func(ob interface{}) {
			p, ok := ob.(*kongv1.KongPlugin)
			if ok && s.isValidIngressClass(&p.ObjectMeta, annotations.IngressClassKey, s.getIngressClassHandling()) {
				plugins = append(plugins, p)
			}
		})
	if err != nil {
		return nil, err
	}
	return plugins, nil
}

// ListGlobalKongClusterPlugins returns all KongClusterPlugin resources
// filtered by the ingress.class annotation and with the
// label global:"true".
func (s Store) ListGlobalKongClusterPlugins() ([]*kongv1.KongClusterPlugin, error) {
	var plugins []*kongv1.KongClusterPlugin

	req, err := labels.NewRequirement("global", selection.Equals, []string{"true"})
	if err != nil {
		return nil, err
	}
	err = cache.ListAll(s.stores.ClusterPlugin,
		labels.NewSelector().Add(*req),
		func(ob interface{}) {
			p, ok := ob.(*kongv1.KongClusterPlugin)
			if ok && s.isValidIngressClass(&p.ObjectMeta, annotations.IngressClassKey, s.getIngressClassHandling()) {
				plugins = append(plugins, p)
			}
		})
	if err != nil {
		return nil, err
	}
	return plugins, nil
}

// ListKongClusterPlugins lists all KongClusterPlugins that match expected ingress.class annotation.
func (s Store) ListKongClusterPlugins() []*kongv1.KongClusterPlugin {
	var plugins []*kongv1.KongClusterPlugin
	for _, item := range s.stores.ClusterPlugin.List() {
		p, ok := item.(*kongv1.KongClusterPlugin)
		if ok && s.isValidIngressClass(&p.ObjectMeta, annotations.IngressClassKey, s.getIngressClassHandling()) {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

// ListKongPlugins lists all KongPlugins.
func (s Store) ListKongPlugins() []*kongv1.KongPlugin {
	var plugins []*kongv1.KongPlugin
	for _, item := range s.stores.Plugin.List() {
		p, ok := item.(*kongv1.KongPlugin)
		if ok {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

// ListCACerts returns all Secrets containing the label
// "konghq.com/ca-cert"="true".
func (s Store) ListCACerts() ([]*corev1.Secret, error) {
	var secrets []*corev1.Secret
	req, err := labels.NewRequirement(caCertKey,
		selection.Equals, []string{"true"})
	if err != nil {
		return nil, err
	}
	err = cache.ListAll(s.stores.Secret,
		labels.NewSelector().Add(*req),
		func(ob interface{}) {
			p, ok := ob.(*corev1.Secret)
			if ok && s.isValidIngressClass(&p.ObjectMeta, annotations.IngressClassKey, s.getIngressClassHandling()) {
				secrets = append(secrets, p)
			}
		})
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

// getIngressClassHandling returns annotations.ExactOrEmptyClassMatch if an IngressClass is the default class, or
// annotations.ExactClassMatch if the IngressClass is not default or does not exist.
func (s Store) getIngressClassHandling() annotations.ClassMatching {
	class, err := s.GetIngressClassV1(s.ingressClass)
	if err != nil {
		s.logger.Debugf("IngressClass %s not found", s.ingressClass)
		return annotations.ExactClassMatch
	}
	if ctrlutils.IsDefaultIngressClass(class) {
		return annotations.ExactOrEmptyClassMatch
	}
	return annotations.ExactClassMatch
}

// convUnstructuredObj is a convenience function to quickly convert any runtime.Object where the underlying type
// is an *unstructured.Unstructured (client-go's dynamic client type) and convert that object to a runtime.Object
// which is backed by the API type it represents. You can use the GVK of the runtime.Object to determine what type
// you want to convert to. This function is meant so that storer implementations can optionally work with YAML files
// for caller convenience when initializing new CacheStores objects.
//
// TODO: upon some searching I didn't find an analog to this over in client-go (https://github.com/kubernetes/client-go)
// however I could have just missed it. We should switch if we find something better, OR we should contribute
// this functionality upstream.
func convUnstructuredObj(from, to runtime.Object) error {
	b, err := yaml.Marshal(from)
	if err != nil {
		return fmt.Errorf("failed to convert object %s to yaml: %w", from.GetObjectKind().GroupVersionKind(), err)
	}
	return yaml.Unmarshal(b, to)
}

// mkObjFromGVK is a factory function that returns a concrete implementation runtime.Object
// for the given GVK. Callers can then use `convert()` to convert an unstructured
// runtime.Object into a concrete one.
func mkObjFromGVK(gvk schema.GroupVersionKind) (runtime.Object, error) {
	switch gvk {
	// ----------------------------------------------------------------------------
	// Kubernetes Core APIs
	// ----------------------------------------------------------------------------
	case netv1.SchemeGroupVersion.WithKind("Ingress"):
		return &netv1.Ingress{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case corev1.SchemeGroupVersion.WithKind("Service"):
		return &corev1.Service{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case corev1.SchemeGroupVersion.WithKind("Secret"):
		return &corev1.Secret{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	// ----------------------------------------------------------------------------
	// Kubernetes Discovery APIs
	// ----------------------------------------------------------------------------
	case discoveryv1.SchemeGroupVersion.WithKind("EndpointSlice"):
		return &discoveryv1.EndpointSlice{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	// ----------------------------------------------------------------------------
	// Kubernetes Gateway APIs
	// ----------------------------------------------------------------------------
	case gatewayv1.SchemeGroupVersion.WithKind("HTTPRoute"):
		return &gatewayv1.HTTPRoute{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case gatewayv1alpha2.SchemeGroupVersion.WithKind("GRPCRoute"):
		return &gatewayv1alpha2.GRPCRoute{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case gatewayv1alpha2.SchemeGroupVersion.WithKind("TCPRoute"):
		return &gatewayv1alpha2.TCPRoute{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case gatewayv1alpha2.SchemeGroupVersion.WithKind("UDPRoute"):
		return &gatewayv1alpha2.UDPRoute{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case gatewayv1alpha2.SchemeGroupVersion.WithKind("TLSRoute"):
		return &gatewayv1alpha2.TLSRoute{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case gatewayv1.SchemeGroupVersion.WithKind("ReferenceGrant"):
		return &gatewayv1beta1.ReferenceGrant{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	// ----------------------------------------------------------------------------
	// Kong APIs
	// ----------------------------------------------------------------------------
	case kongv1.SchemeGroupVersion.WithKind("KongIngress"):
		return &kongv1.KongIngress{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case kongv1beta1.SchemeGroupVersion.WithKind("UDPIngress"):
		return &kongv1beta1.UDPIngress{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case kongv1beta1.SchemeGroupVersion.WithKind("TCPIngress"):
		return &kongv1beta1.TCPIngress{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case kongv1.SchemeGroupVersion.WithKind("KongPlugin"):
		return &kongv1.KongPlugin{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case kongv1.SchemeGroupVersion.WithKind("KongClusterPlugin"):
		return &kongv1.KongClusterPlugin{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case kongv1.SchemeGroupVersion.WithKind("KongConsumer"):
		return &kongv1.KongConsumer{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case kongv1beta1.SchemeGroupVersion.WithKind("KongConsumerGroup"):
		return &kongv1beta1.KongConsumerGroup{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	case kongv1alpha1.SchemeGroupVersion.WithKind("IngressClassParameters"):
		return &kongv1alpha1.IngressClassParameters{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	// ----------------------------------------------------------------------------
	// Knative APIs
	// ----------------------------------------------------------------------------
	case knative.SchemeGroupVersion.WithKind("Ingress"):
		return &knative.Ingress{
			TypeMeta: typeMetaFromGVK(gvk),
		}, nil
	default:
		return nil, fmt.Errorf("%s is not a supported runtime.Object", gvk)
	}
}

func typeMetaFromGVK(gvk schema.GroupVersionKind) metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
	}
}
