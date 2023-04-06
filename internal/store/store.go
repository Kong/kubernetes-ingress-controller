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
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured/unstructuredscheme"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/selection"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	GetSecret(ctx context.Context, namespace, name string) (*corev1.Secret, error)
	GetService(ctx context.Context, namespace, name string) (*corev1.Service, error)
	GetEndpointsForService(ctx context.Context, namespace, name string) (*corev1.Endpoints, error)
	GetKongIngress(ctx context.Context, namespace, name string) (*kongv1.KongIngress, error)
	GetKongPlugin(ctx context.Context, namespace, name string) (*kongv1.KongPlugin, error)
	GetKongClusterPlugin(ctx context.Context, name string) (*kongv1.KongClusterPlugin, error)
	GetKongConsumer(ctx context.Context, namespace, name string) (*kongv1.KongConsumer, error)
	GetIngressClassName() string
	GetIngressClassV1(ctx context.Context, name string) (*netv1.IngressClass, error)
	GetIngressClassParametersV1Alpha1(ctx context.Context, ingressClass *netv1.IngressClass) (*kongv1alpha1.IngressClassParameters, error)
	GetGateway(ctx context.Context, namespace string, name string) (*gatewayv1beta1.Gateway, error)

	ListIngressesV1beta1(ctx context.Context) []*netv1beta1.Ingress
	ListIngressesV1(ctx context.Context) []*netv1.Ingress
	ListIngressClassesV1(ctx context.Context) []*netv1.IngressClass
	ListIngressClassParametersV1Alpha1(ctx context.Context) []*kongv1alpha1.IngressClassParameters
	ListHTTPRoutes(ctx context.Context) ([]*gatewayv1beta1.HTTPRoute, error)
	ListUDPRoutes(ctx context.Context) ([]*gatewayv1alpha2.UDPRoute, error)
	ListTCPRoutes(ctx context.Context) ([]*gatewayv1alpha2.TCPRoute, error)
	ListTLSRoutes(ctx context.Context) ([]*gatewayv1alpha2.TLSRoute, error)
	ListGRPCRoutes(ctx context.Context) ([]*gatewayv1alpha2.GRPCRoute, error)
	ListReferenceGrants(ctx context.Context) ([]*gatewayv1beta1.ReferenceGrant, error)
	ListGateways(ctx context.Context) ([]*gatewayv1beta1.Gateway, error)
	ListTCPIngresses(ctx context.Context) ([]*kongv1beta1.TCPIngress, error)
	ListUDPIngresses(ctx context.Context) ([]*kongv1beta1.UDPIngress, error)
	ListKnativeIngresses(ctx context.Context) ([]*knative.Ingress, error)
	ListGlobalKongPlugins(ctx context.Context) ([]*kongv1.KongPlugin, error)
	ListGlobalKongClusterPlugins(ctx context.Context) ([]*kongv1.KongClusterPlugin, error)
	ListKongPlugins(ctx context.Context) []*kongv1.KongPlugin
	ListKongClusterPlugins(ctx context.Context) []*kongv1.KongClusterPlugin
	ListKongConsumers(ctx context.Context) []*kongv1.KongConsumer
	ListCACerts(ctx context.Context) ([]*corev1.Secret, error)
}

// Store implements Storer and can be used to list Ingress, Services
// and other resources from k8s APIserver. The backing stores should
// be synced and updated by the caller.
// It is ingressClass filter aware.
type Store struct {
	client client.Client

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
	l *sync.RWMutex

	client client.Client
}

// NewCacheStores is a convenience function for CacheStores to initialize all attributes with new cache stores.
func NewCacheStores(client client.Client) CacheStores {
	return CacheStores{
		client: client,

		// Core Kubernetes Stores
		// Gateway API Stores
		// Kong Stores
		// Knative Stores

		l: &sync.RWMutex{},
	}
}

// NewCacheStoresFromObjYAML provides a new CacheStores object given any number of byte arrays containing
// YAML Kubernetes objects. An error is returned if any provided YAML was not a valid Kubernetes object.
func NewCacheStoresFromObjYAML(client client.Client, objs ...[]byte) (c CacheStores, err error) {
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
	return NewCacheStoresFromObjs(client, kobjs...)
}

// NewCacheStoresFromObjs provides a new CacheStores object given any number of Kubernetes
// objects that should be pre-populated. This function will sort objects into the appropriate
// sub-storage (e.g. IngressV1, TCPIngress, e.t.c.) but will produce an error if any of the
// input objects are erroneous or otherwise unusable as Kubernetes objects.
func NewCacheStoresFromObjs(client client.Client, objs ...runtime.Object) (CacheStores, error) {
	c := NewCacheStores(client)
	for _, obj := range objs {
		typedObj, err := mkObjFromGVK(obj.GetObjectKind().GroupVersionKind())
		if err != nil {
			return c, err
		}

		if err := convUnstructuredObj(obj, typedObj); err != nil {
			return c, err
		}

		// TODO(pmalek)
		// if err := c.Add(typedObj); err != nil {
		// 	return c, err
		// }
	}
	return c, nil
}

// New creates a new object store to be used in the ingress controller.
func New(client client.Client, ingressClass string, logger logrus.FieldLogger) Storer {
	return Store{
		client:                client,
		ingressClass:          ingressClass,
		ingressClassMatching:  annotations.ExactClassMatch,
		isValidIngressClass:   annotations.IngressClassValidatorFuncFromObjectMeta(ingressClass),
		isValidIngressV1Class: annotations.IngressClassValidatorFuncFromV1Ingress(ingressClass),
		logger:                logger,
	}
}

// GetSecret returns a Secret using the namespace and name as key.
func (s Store) GetSecret(ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
	var secret corev1.Secret
	err := s.client.Get(ctx, key, &secret)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound{fmt.Sprintf("Secret %v not found", name)}
		}
		return nil, err
	}
	return &secret, nil
}

// GetService returns a Service using the namespace and name as key.
func (s Store) GetService(ctx context.Context, namespace, name string) (*corev1.Service, error) {
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
	var svc corev1.Service
	err := s.client.Get(ctx, key, &svc)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound{fmt.Sprintf("Service %v not found", name)}
		}
		return nil, err
	}
	return &svc, nil
}

// ListIngressesV1 returns the list of Ingresses in the Ingress v1 store.
func (s Store) ListIngressesV1(ctx context.Context) []*netv1.Ingress {
	var ingresses []*netv1.Ingress
	var ingressesV1 netv1.IngressList
	if err := s.client.List(ctx, &ingressesV1); err != nil {
		s.logger.Errorf("failed to list networking v1 Ingresses: %v", err)
	}

	for i := range ingressesV1.Items {
		ing := ingressesV1.Items[i]
		if ing.ObjectMeta.GetAnnotations()[annotations.IngressClassKey] != "" {
			if !s.isValidIngressClass(&ing.ObjectMeta, annotations.IngressClassKey, s.ingressClassMatching) {
				continue
			}
		} else if ing.Spec.IngressClassName != nil {
			if !s.isValidIngressV1Class(&ing, s.ingressClassMatching) {
				continue
			}
		} else {
			class, err := s.GetIngressClassV1(ctx, s.ingressClass)
			if err != nil {
				s.logger.Debugf("IngressClass %s not found", s.ingressClass)
				continue
			}
			if !ctrlutils.IsDefaultIngressClass(class) {
				continue
			}
		}
		ingresses = append(ingresses, &ing)
	}

	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name),
		) < 0
	})

	return ingresses
}

// ListIngressClassesV1 returns the list of IngressClasses.
func (s Store) ListIngressClassesV1(ctx context.Context) []*netv1.IngressClass {
	// filter ingress rules
	var classes []*netv1.IngressClass
	var classList netv1.IngressClassList
	if err := s.client.List(ctx, &classList); err != nil {
		s.logger.Errorf("failed to list networking v1 IngressClasses: %v", err)
	}
	for i := range classList.Items {
		if classList.Items[i].Spec.Controller != IngressClassKongController {
			continue
		}
		classes = append(classes, &classList.Items[i])
	}

	sort.SliceStable(classes, func(i, j int) bool {
		return strings.Compare(classes[i].Name, classes[j].Name) < 0
	})

	return classes
}

// ListIngressClassParametersV1Alpha1 returns the list of IngressClassParameters in the Ingress v1alpha1 store.
func (s Store) ListIngressClassParametersV1Alpha1(ctx context.Context) []*kongv1alpha1.IngressClassParameters {
	var icpList kongv1alpha1.IngressClassParametersList
	if err := s.client.List(ctx, &icpList); err != nil {
		return nil
	}

	icps := lo.Map(icpList.Items,
		func(p kongv1alpha1.IngressClassParameters, _ int) *kongv1alpha1.IngressClassParameters {
			return &p
		})

	sort.SliceStable(icps, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", icps[i].Namespace, icps[i].Name),
			fmt.Sprintf("%s/%s", icps[j].Namespace, icps[j].Name),
		) < 0
	})

	return icps
}

// ListIngressesV1beta1 returns the list of Ingresses in the Ingress v1beta1 store.
func (s Store) ListIngressesV1beta1(ctx context.Context) []*netv1beta1.Ingress {
	// filter ingress rules
	var ingresses []*netv1beta1.Ingress
	var ingressesV1Beta1 netv1beta1.IngressList
	if err := s.client.List(ctx, &ingressesV1Beta1); err != nil {
		s.logger.Errorf("failed to list networking v1 Ingresses: %v", err)
	}
	for _, item := range ingressesV1Beta1.Items {
		ing := item
		if !s.isValidIngressClass(&ing.ObjectMeta, annotations.IngressClassKey, s.ingressClassMatching) {
			continue
		}
		ingresses = append(ingresses, &ing)
	}

	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name),
		) < 0
	})
	return ingresses
}

// ListHTTPRoutes returns the list of HTTPRoutes in the HTTPRoute cache store.
func (s Store) ListHTTPRoutes(ctx context.Context) ([]*gatewayv1beta1.HTTPRoute, error) {
	var list gatewayv1beta1.HTTPRouteList
	if err := s.client.List(ctx, &list); err != nil {
		return nil, err
	}

	items := lo.Map(list.Items,
		func(p gatewayv1beta1.HTTPRoute, _ int) *gatewayv1beta1.HTTPRoute {
			return &p
		})

	sort.SliceStable(items, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", items[i].Namespace, items[i].Name),
			fmt.Sprintf("%s/%s", items[j].Namespace, items[j].Name),
		) < 0
	})

	return items, nil
}

// ListUDPRoutes returns the list of UDPRoutes in the UDPRoute cache store.
func (s Store) ListUDPRoutes(ctx context.Context) ([]*gatewayv1alpha2.UDPRoute, error) {
	var list gatewayv1alpha2.UDPRouteList
	if err := s.client.List(ctx, &list); err != nil {
		return nil, err
	}

	items := lo.Map(list.Items,
		func(p gatewayv1alpha2.UDPRoute, _ int) *gatewayv1alpha2.UDPRoute {
			return &p
		})

	sort.SliceStable(items, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", items[i].Namespace, items[i].Name),
			fmt.Sprintf("%s/%s", items[j].Namespace, items[j].Name),
		) < 0
	})

	return items, nil
}

// ListTCPRoutes returns the list of TCPRoutes in the TCPRoute cache store.
func (s Store) ListTCPRoutes(ctx context.Context) ([]*gatewayv1alpha2.TCPRoute, error) {
	var list gatewayv1alpha2.TCPRouteList
	if err := s.client.List(ctx, &list); err != nil {
		return nil, err
	}

	items := lo.Map(list.Items,
		func(p gatewayv1alpha2.TCPRoute, _ int) *gatewayv1alpha2.TCPRoute {
			return &p
		})

	sort.SliceStable(items, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", items[i].Namespace, items[i].Name),
			fmt.Sprintf("%s/%s", items[j].Namespace, items[j].Name),
		) < 0
	})

	return items, nil
}

// ListTLSRoutes returns the list of TLSRoutes in the TLSRoute cache store.
func (s Store) ListTLSRoutes(ctx context.Context) ([]*gatewayv1alpha2.TLSRoute, error) {
	var list gatewayv1alpha2.TLSRouteList
	if err := s.client.List(ctx, &list); err != nil {
		return nil, err
	}

	items := lo.Map(list.Items,
		func(p gatewayv1alpha2.TLSRoute, _ int) *gatewayv1alpha2.TLSRoute {
			return &p
		})

	sort.SliceStable(items, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", items[i].Namespace, items[i].Name),
			fmt.Sprintf("%s/%s", items[j].Namespace, items[j].Name),
		) < 0
	})

	return items, nil
}

// ListGRPCRoutes returns the list of GRPCRoutes in the GRPCRoute cache store.
func (s Store) ListGRPCRoutes(ctx context.Context) ([]*gatewayv1alpha2.GRPCRoute, error) {
	var list gatewayv1alpha2.GRPCRouteList
	if err := s.client.List(ctx, &list); err != nil {
		return nil, err
	}

	items := lo.Map(list.Items,
		func(p gatewayv1alpha2.GRPCRoute, _ int) *gatewayv1alpha2.GRPCRoute {
			return &p
		})

	sort.SliceStable(items, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", items[i].Namespace, items[i].Name),
			fmt.Sprintf("%s/%s", items[j].Namespace, items[j].Name),
		) < 0
	})

	return items, nil
}

// ListReferenceGrants returns the list of ReferenceGrants in the ReferenceGrant cache store.
func (s Store) ListReferenceGrants(ctx context.Context) ([]*gatewayv1beta1.ReferenceGrant, error) {
	var list gatewayv1beta1.ReferenceGrantList
	if err := s.client.List(ctx, &list); err != nil {
		return nil, err
	}

	items := lo.Map(list.Items,
		func(p gatewayv1beta1.ReferenceGrant, _ int) *gatewayv1beta1.ReferenceGrant {
			return &p
		})

	sort.SliceStable(items, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", items[i].Namespace, items[i].Name),
			fmt.Sprintf("%s/%s", items[j].Namespace, items[j].Name),
		) < 0
	})

	return items, nil
}

// ListGateways returns the list of Gateways in the Gateway cache store.
func (s Store) ListGateways(ctx context.Context) ([]*gatewayv1beta1.Gateway, error) {
	var list gatewayv1beta1.GatewayList
	if err := s.client.List(ctx, &list); err != nil {
		return nil, err
	}

	items := lo.Map(list.Items,
		func(p gatewayv1beta1.Gateway, _ int) *gatewayv1beta1.Gateway {
			return &p
		})

	sort.SliceStable(items, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", items[i].Namespace, items[i].Name),
			fmt.Sprintf("%s/%s", items[j].Namespace, items[j].Name),
		) < 0
	})

	return items, nil
}

// ListTCPIngresses returns the list of TCP Ingresses from
// configuration.konghq.com group.
func (s Store) ListTCPIngresses(ctx context.Context) ([]*kongv1beta1.TCPIngress, error) {
	var (
		ingresses   []*kongv1beta1.TCPIngress
		ingressList kongv1beta1.TCPIngressList
	)
	if err := s.client.List(ctx, &ingressList); err != nil {
		return nil, err
	}

	handlingClass := s.getIngressClassHandling(ctx)
	for i := range ingressList.Items {
		ingress := ingressList.Items[i]
		if s.isValidIngressClass(&ingress.ObjectMeta, annotations.IngressClassKey, handlingClass) {
			ingresses = append(ingresses, &ingress)
		}
	}

	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name),
		) < 0
	})
	return ingresses, nil
}

// ListUDPIngresses returns the list of UDP Ingresses.
func (s Store) ListUDPIngresses(ctx context.Context) ([]*kongv1beta1.UDPIngress, error) {
	var (
		ingresses   []*kongv1beta1.UDPIngress
		ingressList kongv1beta1.UDPIngressList
	)
	if err := s.client.List(ctx, &ingressList); err != nil {
		// older versions of the KIC do not support UDPIngress so short circuit to maintain support with them
		return nil, err
	}

	handlingClass := s.getIngressClassHandling(ctx)
	for i := range ingressList.Items {
		ingress := ingressList.Items[i]
		if s.isValidIngressClass(&ingress.ObjectMeta, annotations.IngressClassKey, handlingClass) {
			ingresses = append(ingresses, &ingress)
		}
	}

	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name),
		) < 0
	})
	return ingresses, nil
}

// ListKnativeIngresses returns the list of Knative Ingresses from
// ingresses.networking.internal.knative.dev group.
func (s Store) ListKnativeIngresses(ctx context.Context) ([]*knative.Ingress, error) {
	var (
		ingresses   []*knative.Ingress
		ingressList knative.IngressList
	)
	if err := s.client.List(ctx, &ingressList); err != nil {
		return nil, err
	}

	handlingClass := s.getIngressClassHandling(ctx)
	for i := range ingressList.Items {
		ingress := ingressList.Items[i]
		if s.isValidIngressClass(&ingress.ObjectMeta, annotations.KnativeIngressClassKey, handlingClass) ||
			s.isValidIngressClass(&ingress.ObjectMeta, annotations.KnativeIngressClassDeprecatedKey, handlingClass) {
			ingresses = append(ingresses, &ingress)
		}

	}

	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name),
		) < 0
	})
	return ingresses, nil
}

// GetEndpointsForService returns the internal endpoints for service 'namespace/name' inside k8s.
func (s Store) GetEndpointsForService(ctx context.Context, namespace, name string) (*corev1.Endpoints, error) {
	var (
		endpoints corev1.Endpoints
		key       = client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}
	)

	err := s.client.Get(ctx, key, &endpoints)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound{fmt.Sprintf("Endpoints for service %v not found", key)}
		}
		return nil, err
	}
	return &endpoints, nil
}

// GetKongPlugin returns the 'name' KongPlugin resource in namespace.
func (s Store) GetKongPlugin(ctx context.Context, namespace, name string) (*kongv1.KongPlugin, error) {
	var (
		plugin kongv1.KongPlugin
		key    = client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}
	)

	err := s.client.Get(ctx, key, &plugin)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound{fmt.Sprintf("KongPlugin %v not found", key)}
		}
		return nil, err
	}
	return &plugin, nil
}

// GetKongClusterPlugin returns the 'name' KongClusterPlugin resource.
func (s Store) GetKongClusterPlugin(ctx context.Context, name string) (*kongv1.KongClusterPlugin, error) {
	var (
		plugin kongv1.KongClusterPlugin
		key    = client.ObjectKey{
			Name: name,
		}
	)

	err := s.client.Get(ctx, key, &plugin)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound{fmt.Sprintf("KongClusterPlugin %v not found", key)}
		}
		return nil, err
	}
	return &plugin, nil
}

// GetKongIngress returns the 'name' KongIngress resource in namespace.
func (s Store) GetKongIngress(ctx context.Context, namespace, name string) (*kongv1.KongIngress, error) {
	var (
		ingress kongv1.KongIngress
		key     = client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}
	)

	err := s.client.Get(ctx, key, &ingress)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound{fmt.Sprintf("KongIngress %v not found", key)}
		}
		return nil, err
	}
	return &ingress, nil
}

// GetKongConsumer returns the 'name' KongConsumer resource in namespace.
func (s Store) GetKongConsumer(ctx context.Context, namespace, name string) (*kongv1.KongConsumer, error) {
	var (
		consumer kongv1.KongConsumer
		key      = client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}
	)

	err := s.client.Get(ctx, key, &consumer)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound{fmt.Sprintf("KongConsumer %v not found", key)}
		}
		return nil, err
	}
	return &consumer, nil
}

func (s Store) GetIngressClassName() string {
	return s.ingressClass
}

// GetIngressClassV1 returns the 'name' IngressClass resource.
func (s Store) GetIngressClassV1(ctx context.Context, name string) (*netv1.IngressClass, error) {
	var (
		ingressClass netv1.IngressClass
		key          = client.ObjectKey{
			Name: name,
		}
	)

	err := s.client.Get(ctx, key, &ingressClass)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound{fmt.Sprintf("IngressClass %v not found", name)}
		}
		return nil, err
	}
	return &ingressClass, nil
}

// GetIngressClassParametersV1Alpha1 returns IngressClassParameters for provided
// IngressClass.
func (s Store) GetIngressClassParametersV1Alpha1(ctx context.Context, ingressClass *netv1.IngressClass) (*kongv1alpha1.IngressClassParameters, error) {
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

	key := client.ObjectKey{
		Namespace: *ingressClass.Spec.Parameters.Namespace,
		Name:      ingressClass.Spec.Parameters.Name,
	}
	var icp kongv1alpha1.IngressClassParameters
	err := s.client.Get(ctx, key, &icp)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound{fmt.Sprintf("IngressClassParameters %v not found", key)}
		}
		return nil, err
	}
	return &icp, nil
}

// GetGateway returns gateway resource having specified namespace and name.
func (s Store) GetGateway(ctx context.Context, namespace string, name string) (*gatewayv1beta1.Gateway, error) {
	var (
		gateway gatewayv1beta1.Gateway
		key     = client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}
	)

	err := s.client.Get(ctx, key, &gateway)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound{fmt.Sprintf("Gateway %v not found", key)}
		}
		return nil, err
	}
	return &gateway, nil
}

// ListKongConsumers returns all KongConsumers filtered by the ingress.class
// annotation.
func (s Store) ListKongConsumers(ctx context.Context) []*kongv1.KongConsumer {
	var (
		consumers    []*kongv1.KongConsumer
		consumerList kongv1.KongConsumerList
	)
	if err := s.client.List(ctx, &consumerList); err != nil {
		return nil
	}

	handlingClass := s.getIngressClassHandling(ctx)
	for i := range consumerList.Items {
		plugin := consumerList.Items[i]
		if s.isValidIngressClass(&plugin.ObjectMeta, annotations.IngressClassKey, handlingClass) {
			consumers = append(consumers, &plugin)
		}
	}
	return consumers
}

// ListGlobalKongPlugins returns all KongPlugin resources
// filtered by the ingress.class annotation and with the
// label global:"true".
// Support for these global namespaced KongPlugins was removed in 0.10.0
// This function remains only to provide warnings to users with old configuration.
func (s Store) ListGlobalKongPlugins(ctx context.Context) ([]*kongv1.KongPlugin, error) {
	req, err := labels.NewRequirement("global", selection.Equals, []string{"true"})
	if err != nil {
		return nil, err
	}

	var (
		plugins    []*kongv1.KongPlugin
		pluginList kongv1.KongPluginList
	)
	err = s.client.List(ctx, &pluginList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*req),
	})
	if err != nil {
		return nil, err
	}

	handlingClass := s.getIngressClassHandling(ctx)
	for i := range pluginList.Items {
		plugin := pluginList.Items[i]
		if s.isValidIngressClass(&plugin.ObjectMeta, annotations.IngressClassKey, handlingClass) {
			plugins = append(plugins, &plugin)
		}
	}

	return plugins, nil
}

// ListGlobalKongClusterPlugins returns all KongClusterPlugin resources
// filtered by the ingress.class annotation and with the
// label global:"true".
func (s Store) ListGlobalKongClusterPlugins(ctx context.Context) ([]*kongv1.KongClusterPlugin, error) {
	req, err := labels.NewRequirement("global", selection.Equals, []string{"true"})
	if err != nil {
		return nil, err
	}

	var (
		plugins    []*kongv1.KongClusterPlugin
		pluginList kongv1.KongClusterPluginList
	)
	err = s.client.List(ctx, &pluginList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*req),
	})
	if err != nil {
		return nil, err
	}

	handlingClass := s.getIngressClassHandling(ctx)
	for i := range pluginList.Items {
		plugin := pluginList.Items[i]
		if s.isValidIngressClass(&plugin.ObjectMeta, annotations.IngressClassKey, handlingClass) {
			plugins = append(plugins, &plugin)
		}
	}

	return plugins, nil
}

// ListKongClusterPlugins lists all KongClusterPlugins that match expected ingress.class annotation.
func (s Store) ListKongClusterPlugins(ctx context.Context) []*kongv1.KongClusterPlugin {
	var (
		plugins    []*kongv1.KongClusterPlugin
		pluginList kongv1.KongClusterPluginList
	)
	if err := s.client.List(ctx, &pluginList); err != nil {
		return nil
	}

	handlingClass := s.getIngressClassHandling(ctx)
	for i := range pluginList.Items {
		plugin := pluginList.Items[i]
		if s.isValidIngressClass(&plugin.ObjectMeta, annotations.IngressClassKey, handlingClass) {
			plugins = append(plugins, &plugin)
		}
	}

	return plugins
}

// ListKongPlugins lists all KongPlugins.
func (s Store) ListKongPlugins(ctx context.Context) []*kongv1.KongPlugin {
	var pluginList kongv1.KongPluginList
	if err := s.client.List(ctx, &pluginList); err != nil {
		return nil
	}

	return lo.Map(pluginList.Items,
		func(p kongv1.KongPlugin, _ int) *kongv1.KongPlugin {
			return &p
		})
}

// ListCACerts returns all Secrets containing the label
// "konghq.com/ca-cert"="true".
func (s Store) ListCACerts(ctx context.Context) ([]*corev1.Secret, error) {
	req, err := labels.NewRequirement(caCertKey, selection.Equals, []string{"true"})
	if err != nil {
		return nil, err
	}

	var (
		secrets    []*corev1.Secret
		secretList corev1.SecretList
	)
	err = s.client.List(ctx, &secretList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*req),
	})
	if err != nil {
		return nil, err
	}

	handlingClass := s.getIngressClassHandling(ctx)
	for i := range secretList.Items {
		secret := secretList.Items[i]
		if s.isValidIngressClass(&secret.ObjectMeta, annotations.IngressClassKey, handlingClass) {
			secrets = append(secrets, &secret)
		}
	}

	return secrets, nil
}

func (s Store) networkingIngressV1Beta1(obj interface{}) *netv1beta1.Ingress {
	switch obj := obj.(type) {
	case *netv1beta1.Ingress:
		return obj

	default:
		s.logger.Errorf("cannot convert to networking v1beta1 Ingress: unsupported type: %v", reflect.TypeOf(obj))
		return nil
	}
}

// getIngressClassHandling returns annotations.ExactOrEmptyClassMatch if an IngressClass is the default class, or
// annotations.ExactClassMatch if the IngressClass is not default or does not exist.
func (s Store) getIngressClassHandling(ctx context.Context) annotations.ClassMatching {
	class, err := s.GetIngressClassV1(ctx, s.ingressClass)
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
	case netv1.SchemeGroupVersion.WithKind("IngressClass"):
		return &netv1.IngressClass{}, nil
	case netv1.SchemeGroupVersion.WithKind("Ingress"):
		return &netv1.Ingress{}, nil
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
	case gatewayv1beta1.SchemeGroupVersion.WithKind("HTTPRoutes"):
		return &gatewayv1beta1.HTTPRoute{}, nil
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
	case kongv1alpha1.SchemeGroupVersion.WithKind("IngressClassParameters"):
		return &kongv1alpha1.IngressClassParameters{}, nil
	// ----------------------------------------------------------------------------
	// Knative APIs
	// ----------------------------------------------------------------------------
	case knative.SchemeGroupVersion.WithKind("Ingress"):
		return &knative.Ingress{}, nil
	default:
		return nil, fmt.Errorf("%s is not a supported runtime.Object", gvk)
	}
}
