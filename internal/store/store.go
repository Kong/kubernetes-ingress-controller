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
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
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
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	ctrlutils "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

const (
	caCertKey = "konghq.com/ca-cert"
	// IngressClassKongController is the string used for the Controller field of a recognized IngressClass
	IngressClassKongController = "ingress-controllers.konghq.com/kong"
)

// ErrNotFound error is returned when a lookup results in no resource.
// This type is meant to be used for error handling using `errors.As()`.
type ErrNotFound struct {
	message string
}

func (e ErrNotFound) Error() string {
	if e.message == "" {
		return "not found"
	}
	return e.message
}

// Storer is the interface that wraps the required methods to gather information
// about ingresses, services, secrets and ingress annotations.
type Storer interface {
	GetSecret(namespace, name string) (*corev1.Secret, error)
	GetService(namespace, name string) (*corev1.Service, error)
	GetEndpointsForService(namespace, name string) (*corev1.Endpoints, error)
	GetKongIngress(namespace, name string) (*kongv1.KongIngress, error)
	GetKongPlugin(namespace, name string) (*kongv1.KongPlugin, error)
	GetKongClusterPlugin(name string) (*kongv1.KongClusterPlugin, error)
	GetKongConsumer(namespace, name string) (*kongv1.KongConsumer, error)
	GetIngressClassV1(name string) (*networkingv1.IngressClass, error)

	ListIngressesV1beta1() []*networkingv1beta1.Ingress
	ListIngressesV1() []*networkingv1.Ingress
	ListIngressClassesV1() []*networkingv1.IngressClass
	ListHTTPRoutes() ([]*gatewayv1alpha2.HTTPRoute, error)
	ListUDPRoutes() ([]*gatewayv1alpha2.UDPRoute, error)
	ListTCPRoutes() ([]*gatewayv1alpha2.TCPRoute, error)
	ListTLSRoutes() ([]*gatewayv1alpha2.TLSRoute, error)
	ListReferencePolicies() ([]*gatewayv1alpha2.ReferencePolicy, error)
	ListTCPIngresses() ([]*kongv1beta1.TCPIngress, error)
	ListUDPIngresses() ([]*kongv1beta1.UDPIngress, error)
	ListKnativeIngresses() ([]*knative.Ingress, error)
	ListGlobalKongPlugins() ([]*kongv1.KongPlugin, error)
	ListGlobalKongClusterPlugins() ([]*kongv1.KongClusterPlugin, error)
	ListKongConsumers() []*kongv1.KongConsumer
	ListCACerts() ([]*corev1.Secret, error)
}

// Store implements Storer and can be used to list Ingress, Services
// and other resources from k8s APIserver. The backing stores should
// be synced and updated by the caller.
// It is ingressClass filter aware.
type Store struct {
	stores CacheStores

	ingressClass string

	ingressV1Beta1ClassMatching annotations.ClassMatching
	ingressV1ClassMatching      annotations.ClassMatching
	kongConsumerClassMatching   annotations.ClassMatching

	isValidIngressClass   func(objectMeta *metav1.ObjectMeta, annotation string, handling annotations.ClassMatching) bool
	isValidIngressV1Class func(ingress *networkingv1.Ingress, handling annotations.ClassMatching) bool

	logger logrus.FieldLogger
}

// CacheStores stores cache.Store for all Kinds of k8s objects that
// the Ingress Controller reads.
type CacheStores struct {
	// Core Kubernetes Stores
	IngressV1beta1 cache.Store
	IngressV1      cache.Store
	IngressClassV1 cache.Store
	Service        cache.Store
	Secret         cache.Store
	Endpoint       cache.Store

	// Gateway API Stores
	HTTPRoute       cache.Store
	UDPRoute        cache.Store
	TCPRoute        cache.Store
	TLSRoute        cache.Store
	ReferencePolicy cache.Store

	// Kong Stores
	Plugin        cache.Store
	ClusterPlugin cache.Store
	Consumer      cache.Store
	KongIngress   cache.Store
	TCPIngress    cache.Store
	UDPIngress    cache.Store

	// Knative Stores
	KnativeIngress cache.Store

	l *sync.RWMutex
}

// NewCacheStores is a convenience function for CacheStores to initialize all attributes with new cache stores
func NewCacheStores() (c CacheStores) {
	c.ClusterPlugin = cache.NewStore(clusterResourceKeyFunc)
	c.Consumer = cache.NewStore(keyFunc)
	c.Endpoint = cache.NewStore(keyFunc)
	c.IngressV1 = cache.NewStore(keyFunc)
	c.IngressClassV1 = cache.NewStore(clusterResourceKeyFunc)
	c.IngressV1beta1 = cache.NewStore(keyFunc)
	c.HTTPRoute = cache.NewStore(keyFunc)
	c.UDPRoute = cache.NewStore(keyFunc)
	c.TCPRoute = cache.NewStore(keyFunc)
	c.TLSRoute = cache.NewStore(keyFunc)
	c.ReferencePolicy = cache.NewStore(keyFunc)
	c.KnativeIngress = cache.NewStore(keyFunc)
	c.Plugin = cache.NewStore(keyFunc)
	c.Secret = cache.NewStore(keyFunc)
	c.Service = cache.NewStore(keyFunc)
	c.TCPIngress = cache.NewStore(keyFunc)
	c.UDPIngress = cache.NewStore(keyFunc)
	c.KongIngress = cache.NewStore(keyFunc)
	c.l = &sync.RWMutex{}
	return
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
	case *extensions.Ingress:
		return c.IngressV1beta1.Get(obj)
	case *networkingv1beta1.Ingress:
		return c.IngressV1beta1.Get(obj)
	case *networkingv1.Ingress:
		return c.IngressV1.Get(obj)
	case *networkingv1.IngressClass:
		return c.IngressClassV1.Get(obj)
	case *corev1.Service:
		return c.Service.Get(obj)
	case *corev1.Secret:
		return c.Secret.Get(obj)
	case *corev1.Endpoints:
		return c.Endpoint.Get(obj)
	// ----------------------------------------------------------------------------
	// Kubernetes Gateway API Support
	// ----------------------------------------------------------------------------
	case *gatewayv1alpha2.HTTPRoute:
		return c.HTTPRoute.Get(obj)
	case *gatewayv1alpha2.UDPRoute:
		return c.UDPRoute.Get(obj)
	case *gatewayv1alpha2.TCPRoute:
		return c.TCPRoute.Get(obj)
	case *gatewayv1alpha2.TLSRoute:
		return c.TLSRoute.Get(obj)
	case *gatewayv1alpha2.ReferencePolicy:
		return c.ReferencePolicy.Get(obj)
	// ----------------------------------------------------------------------------
	// Kong API Support
	// ----------------------------------------------------------------------------
	case *kongv1.KongPlugin:
		return c.Plugin.Get(obj)
	case *kongv1.KongClusterPlugin:
		return c.ClusterPlugin.Get(obj)
	case *kongv1.KongConsumer:
		return c.Consumer.Get(obj)
	case *kongv1.KongIngress:
		return c.KongIngress.Get(obj)
	case *kongv1beta1.TCPIngress:
		return c.TCPIngress.Get(obj)
	case *kongv1beta1.UDPIngress:
		return c.UDPIngress.Get(obj)
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
	case *extensions.Ingress:
		return c.IngressV1beta1.Add(obj)
	case *networkingv1beta1.Ingress:
		return c.IngressV1beta1.Add(obj)
	case *networkingv1.Ingress:
		return c.IngressV1.Add(obj)
	case *networkingv1.IngressClass:
		return c.IngressClassV1.Add(obj)
	case *corev1.Service:
		return c.Service.Add(obj)
	case *corev1.Secret:
		return c.Secret.Add(obj)
	case *corev1.Endpoints:
		return c.Endpoint.Add(obj)
	// ----------------------------------------------------------------------------
	// Kubernetes Gateway API Support
	// ----------------------------------------------------------------------------
	case *gatewayv1alpha2.HTTPRoute:
		return c.HTTPRoute.Add(obj)
	case *gatewayv1alpha2.UDPRoute:
		return c.UDPRoute.Add(obj)
	case *gatewayv1alpha2.TCPRoute:
		return c.TCPRoute.Add(obj)
	case *gatewayv1alpha2.TLSRoute:
		return c.TLSRoute.Add(obj)
	case *gatewayv1alpha2.ReferencePolicy:
		return c.ReferencePolicy.Add(obj)
	// ----------------------------------------------------------------------------
	// Kong API Support
	// ----------------------------------------------------------------------------
	case *kongv1.KongPlugin:
		return c.Plugin.Add(obj)
	case *kongv1.KongClusterPlugin:
		return c.ClusterPlugin.Add(obj)
	case *kongv1.KongConsumer:
		return c.Consumer.Add(obj)
	case *kongv1.KongIngress:
		return c.KongIngress.Add(obj)
	case *kongv1beta1.TCPIngress:
		return c.TCPIngress.Add(obj)
	case *kongv1beta1.UDPIngress:
		return c.UDPIngress.Add(obj)
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
	case *extensions.Ingress:
		return c.IngressV1beta1.Delete(obj)
	case *networkingv1beta1.Ingress:
		return c.IngressV1beta1.Delete(obj)
	case *networkingv1.Ingress:
		return c.IngressV1.Delete(obj)
	case *networkingv1.IngressClass:
		return c.IngressClassV1.Delete(obj)
	case *corev1.Service:
		return c.Service.Delete(obj)
	case *corev1.Secret:
		return c.Secret.Delete(obj)
	case *corev1.Endpoints:
		return c.Endpoint.Delete(obj)
	// ----------------------------------------------------------------------------
	// Kubernetes Gateway API Support
	// ----------------------------------------------------------------------------
	case *gatewayv1alpha2.HTTPRoute:
		return c.HTTPRoute.Delete(obj)
	case *gatewayv1alpha2.UDPRoute:
		return c.UDPRoute.Delete(obj)
	case *gatewayv1alpha2.TCPRoute:
		return c.TCPRoute.Delete(obj)
	case *gatewayv1alpha2.TLSRoute:
		return c.TLSRoute.Delete(obj)
	case *gatewayv1alpha2.ReferencePolicy:
		return c.ReferencePolicy.Delete(obj)
	// ----------------------------------------------------------------------------
	// Kong API Support
	// ----------------------------------------------------------------------------
	case *kongv1.KongPlugin:
		return c.Plugin.Delete(obj)
	case *kongv1.KongClusterPlugin:
		return c.ClusterPlugin.Delete(obj)
	case *kongv1.KongConsumer:
		return c.Consumer.Delete(obj)
	case *kongv1.KongIngress:
		return c.KongIngress.Delete(obj)
	case *kongv1beta1.TCPIngress:
		return c.TCPIngress.Delete(obj)
	case *kongv1beta1.UDPIngress:
		return c.UDPIngress.Delete(obj)
	// ----------------------------------------------------------------------------
	// 3rd Party API Support
	// ----------------------------------------------------------------------------
	case *knative.Ingress:
		return c.KnativeIngress.Delete(obj)
	default:
		return fmt.Errorf("cannot delete unsupported kind %q from the store", obj.GetObjectKind().GroupVersionKind())
	}
}

// New creates a new object store to be used in the ingress controller
func New(cs CacheStores, ingressClass string, processClasslessIngressV1Beta1 bool, processClasslessIngressV1 bool,
	processClasslessKongConsumer bool, logger logrus.FieldLogger) Storer {
	var ingressV1Beta1ClassMatching annotations.ClassMatching
	var ingressV1ClassMatching annotations.ClassMatching
	var kongConsumerClassMatching annotations.ClassMatching
	if processClasslessIngressV1Beta1 {
		ingressV1Beta1ClassMatching = annotations.ExactOrEmptyClassMatch
	} else {
		ingressV1Beta1ClassMatching = annotations.ExactClassMatch
	}
	if processClasslessIngressV1 {
		ingressV1ClassMatching = annotations.ExactOrEmptyClassMatch
	} else {
		ingressV1ClassMatching = annotations.ExactClassMatch
	}
	if processClasslessKongConsumer {
		kongConsumerClassMatching = annotations.ExactOrEmptyClassMatch
	} else {
		kongConsumerClassMatching = annotations.ExactClassMatch
	}
	return Store{
		stores:                      cs,
		ingressClass:                ingressClass,
		ingressV1Beta1ClassMatching: ingressV1Beta1ClassMatching,
		ingressV1ClassMatching:      ingressV1ClassMatching,
		kongConsumerClassMatching:   kongConsumerClassMatching,
		isValidIngressClass:         annotations.IngressClassValidatorFuncFromObjectMeta(ingressClass),
		isValidIngressV1Class:       annotations.IngressClassValidatorFuncFromV1Ingress(ingressClass),
		logger:                      logger,
	}
}

// GetSecret returns a Secret using the namespace and name as key
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

// GetService returns a Service using the namespace and name as key
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
func (s Store) ListIngressesV1() []*networkingv1.Ingress {
	// filter ingress rules
	var ingresses []*networkingv1.Ingress
	for _, item := range s.stores.IngressV1.List() {
		ing, ok := item.(*networkingv1.Ingress)
		if !ok {
			s.logger.Warnf("listIngressesV1: dropping object of unexpected type: %#v", item)
			continue
		}
		if ing.ObjectMeta.GetAnnotations()[annotations.IngressClassKey] != "" {
			if !s.isValidIngressClass(&ing.ObjectMeta, annotations.IngressClassKey, s.ingressV1ClassMatching) {
				continue
			}
		} else if ing.Spec.IngressClassName != nil {
			if !s.isValidIngressV1Class(ing, s.ingressV1ClassMatching) {
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
func (s Store) ListIngressClassesV1() []*networkingv1.IngressClass {
	// filter ingress rules
	var classes []*networkingv1.IngressClass
	for _, item := range s.stores.IngressClassV1.List() {
		class, ok := item.(*networkingv1.IngressClass)
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

// ListIngressesV1beta1 returns the list of Ingresses in the Ingress v1beta1 store.
func (s Store) ListIngressesV1beta1() []*networkingv1beta1.Ingress {
	// filter ingress rules
	var ingresses []*networkingv1beta1.Ingress
	for _, item := range s.stores.IngressV1beta1.List() {
		ing := s.networkingIngressV1Beta1(item)
		if !s.isValidIngressClass(&ing.ObjectMeta, annotations.IngressClassKey, s.ingressV1Beta1ClassMatching) {
			continue
		}
		ingresses = append(ingresses, ing)
	}
	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name)) < 0
	})
	return ingresses
}

// ListHTTPRoutes returns the list of HTTPRoutes in the HTTPRoute cache store.
func (s Store) ListHTTPRoutes() ([]*gatewayv1alpha2.HTTPRoute, error) {
	var httproutes []*gatewayv1alpha2.HTTPRoute
	if err := cache.ListAll(s.stores.HTTPRoute, labels.NewSelector(),
		func(ob interface{}) {
			httproute, ok := ob.(*gatewayv1alpha2.HTTPRoute)
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

// ListReferencePolicies returns the list of ReferencePolicies in the ReferencePolicy cache store.
func (s Store) ListReferencePolicies() ([]*gatewayv1alpha2.ReferencePolicy, error) {
	var policies []*gatewayv1alpha2.ReferencePolicy
	if err := cache.ListAll(s.stores.ReferencePolicy, labels.NewSelector(),
		func(ob interface{}) {
			policy, ok := ob.(*gatewayv1alpha2.ReferencePolicy)
			if ok {
				policies = append(policies, policy)
			}
		},
	); err != nil {
		return nil, err
	}
	return policies, nil
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

// ListUDPIngresses returns the list of UDP Ingresses
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

// GetEndpointsForService returns the internal endpoints for service
// 'namespace/name' inside k8s.
func (s Store) GetEndpointsForService(namespace, name string) (*corev1.Endpoints, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	eps, exists, err := s.stores.Endpoint.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("Endpoints for service %v not found", key)}
	}
	return eps.(*corev1.Endpoints), nil
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

// GetIngressClassV1 returns the 'name' IngressClass resource
func (s Store) GetIngressClassV1(name string) (*networkingv1.IngressClass, error) {
	p, exists, err := s.stores.IngressClassV1.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("IngressClass %v not found", name)}
	}
	return p.(*networkingv1.IngressClass), nil
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

// ListGlobalKongPlugins returns all KongPlugin resources
// filtered by the ingress.class annotation and with the
// label global:"true".
// Support for these global namespaced KongPlugins was removed in 0.10.0
// This function remains only to provide warnings to users with old configuration
func (s Store) ListGlobalKongPlugins() ([]*kongv1.KongPlugin, error) {

	var plugins []*kongv1.KongPlugin
	// var globalPlugins []*configurationv1.KongPlugin
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

func (s Store) networkingIngressV1Beta1(obj interface{}) *networkingv1beta1.Ingress {
	switch obj := obj.(type) {
	case *networkingv1beta1.Ingress:
		return obj

	case *extensions.Ingress:
		out, err := toNetworkingIngressV1Beta1(obj)
		if err != nil {
			s.logger.Errorf("cannot convert to networking v1beta1 Ingress: %v", err)
			return nil
		}
		return out

	default:
		s.logger.Errorf("cannot convert to networking v1beta1 Ingress: unsupported type: %v", reflect.TypeOf(obj))
		return nil
	}
}

// getIngressClassHandling returns annotations.ExactOrEmptyClassMatch if an IngressClass is the default class, or
// annotations.ExactClassMatch if the IngressClass is not default or does not exist
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

func toNetworkingIngressV1Beta1(obj *extensions.Ingress) (*networkingv1beta1.Ingress, error) {
	js, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize object of type %v: %w", reflect.TypeOf(obj), err)
	}
	var out networkingv1beta1.Ingress
	if err := json.Unmarshal(js, &out); err != nil {
		return nil, fmt.Errorf("failed to deserialize json: %w", err)
	}
	out.APIVersion = networkingv1beta1.SchemeGroupVersion.String()
	return &out, nil
}

// convUnstructuredObj is a convenience function to quickly convert any runtime.Object where the underlying type
// is an *unstructured.Unstructured (client-go's dynamic client type) and convert that object to a runtime.Object
// which is backed by the API type it represents. You can use the GVK of the runtime.Object to determine what type
// you want to convert to. This function is meant so that storer implementations can optionally work with YAML files
// for caller convenience when initializing new CacheStores objects.
//
// TODO: upon some searching I didn't find an analog to this over in client-go (https://github.com/kubernetes/client-go)
//       however I could have just missed it. We should switch if we find something better, OR we should contribute
//       this functionality upstream.
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
	case extensions.SchemeGroupVersion.WithKind("Ingress"):
		return &extensions.Ingress{}, nil
	case networkingv1.SchemeGroupVersion.WithKind("Ingress"):
		return &networkingv1.Ingress{}, nil
	case kongv1beta1.SchemeGroupVersion.WithKind("TCPIngress"):
		return &kongv1beta1.TCPIngress{}, nil
	case corev1.SchemeGroupVersion.WithKind("Service"):
		return &corev1.Service{}, nil
	case corev1.SchemeGroupVersion.WithKind("Secret"):
		return &corev1.Secret{}, nil
	case corev1.SchemeGroupVersion.WithKind("Endpoints"):
		return &corev1.Endpoints{}, nil
	// ----------------------------------------------------------------------------
	// Kubernetes Gateway APIs
	// ----------------------------------------------------------------------------
	case gatewayv1alpha2.SchemeGroupVersion.WithKind("HTTPRoutes"):
		return &gatewayv1alpha2.HTTPRoute{}, nil
	// ----------------------------------------------------------------------------
	// Kong APIs
	// ----------------------------------------------------------------------------
	case kongv1.SchemeGroupVersion.WithKind("KongIngress"):
		return &kongv1.KongIngress{}, nil
	case kongv1beta1.SchemeGroupVersion.WithKind("UDPIngress"):
		return &kongv1beta1.UDPIngress{}, nil
	case kongv1.SchemeGroupVersion.WithKind("KongPlugin"):
		return &kongv1.KongPlugin{}, nil
	case kongv1.SchemeGroupVersion.WithKind("KongClusterPlugin"):
		return &kongv1.KongClusterPlugin{}, nil
	case kongv1.SchemeGroupVersion.WithKind("KongConsumer"):
		return &kongv1.KongConsumer{}, nil
	// ----------------------------------------------------------------------------
	// Knative APIs
	// ----------------------------------------------------------------------------
	case knative.SchemeGroupVersion.WithKind("Ingress"):
		return &knative.Ingress{}, nil
	default:
		return nil, fmt.Errorf("%s is not a supported runtime.Object", gvk)
	}
}
