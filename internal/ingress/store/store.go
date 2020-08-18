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

	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	"github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/tools/cache"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
)

const (
	knativeIngressClassKey = "networking.knative.dev/ingress.class"
	caCertKey              = "konghq.com/ca-cert"
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
	GetSecret(namespace, name string) (*apiv1.Secret, error)
	GetService(namespace, name string) (*apiv1.Service, error)
	GetEndpointsForService(namespace, name string) (*apiv1.Endpoints, error)
	GetKongIngress(namespace, name string) (*configurationv1.KongIngress, error)
	GetKongPlugin(namespace, name string) (*configurationv1.KongPlugin, error)
	GetKongClusterPlugin(name string) (*configurationv1.KongClusterPlugin, error)
	GetKongConsumer(namespace, name string) (*configurationv1.KongConsumer, error)

	ListIngresses() []*networkingv1beta1.Ingress
	ListTCPIngresses() ([]*configurationv1beta1.TCPIngress, error)
	ListKnativeIngresses() ([]*knative.Ingress, error)
	ListGlobalKongPlugins() ([]*configurationv1.KongPlugin, error)
	ListGlobalKongClusterPlugins() ([]*configurationv1.KongClusterPlugin, error)
	ListKongConsumers() []*configurationv1.KongConsumer
	ListKongCredentials() []*configurationv1.KongCredential
	ListCACerts() ([]*apiv1.Secret, error)
}

// Store implements Storer and can be used to list Ingress, Services
// and other resources from k8s APIserver. The backing stores should
// be synced and updated by the caller.
// It is ingressClass filter aware.
type Store struct {
	stores CacheStores

	ingressClass string

	skipClasslessIngress bool

	isValidIngresClass func(objectMeta *metav1.ObjectMeta) bool

	logger logrus.FieldLogger
}

// CacheStores stores cache.Store for all Kinds of k8s objects that
// the Ingress Controller reads.
type CacheStores struct {
	Ingress    cache.Store
	TCPIngress cache.Store
	Service    cache.Store
	Secret     cache.Store
	Endpoint   cache.Store

	Plugin        cache.Store
	ClusterPlugin cache.Store
	Consumer      cache.Store
	Credential    cache.Store
	Configuration cache.Store

	KnativeIngress cache.Store
}

// New creates a new object store to be used in the ingress controller
func New(cs CacheStores, ingressClass string, skipClasslessIngress bool, logger logrus.FieldLogger) Storer {
	return Store{
		stores:               cs,
		ingressClass:         ingressClass,
		skipClasslessIngress: skipClasslessIngress,
		isValidIngresClass:   annotations.IngressClassValidatorFuncFromObjectMeta(ingressClass, skipClasslessIngress),
		logger:               logger,
	}
}

// GetSecret returns a Secret using the namespace and name as key
func (s Store) GetSecret(namespace, name string) (*apiv1.Secret, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	secret, exists, err := s.stores.Secret.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("Secret %v not found", key)}
	}
	return secret.(*apiv1.Secret), nil
}

// GetService returns a Service using the namespace and name as key
func (s Store) GetService(namespace, name string) (*apiv1.Service, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	service, exists, err := s.stores.Service.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("Service %v not found", key)}
	}
	return service.(*apiv1.Service), nil
}

// ListIngresses returns the list of Ingresses
func (s Store) ListIngresses() []*networkingv1beta1.Ingress {
	// filter ingress rules
	var ingresses []*networkingv1beta1.Ingress
	for _, item := range s.stores.Ingress.List() {
		ing := s.networkingIngressV1Beta1(item)
		if !s.isValidIngresClass(&ing.ObjectMeta) {
			continue
		}
		ingresses = append(ingresses, ing)
	}

	return ingresses
}

// ListTCPIngresses returns the list of TCP Ingresses from
// configuration.konghq.com group.
func (s Store) ListTCPIngresses() ([]*configurationv1beta1.TCPIngress, error) {
	var ingresses []*configurationv1beta1.TCPIngress
	err := cache.ListAll(s.stores.TCPIngress, labels.NewSelector(),
		func(ob interface{}) {
			ing, ok := ob.(*configurationv1beta1.TCPIngress)
			if ok && s.isValidIngresClass(&ing.ObjectMeta) {
				ingresses = append(ingresses, ing)
			}
		})
	if err != nil {
		return nil, err
	}
	return ingresses, nil
}

func (s Store) validKnativeIngressClass(objectMeta *metav1.ObjectMeta) bool {
	ingressAnnotationValue := objectMeta.GetAnnotations()[knativeIngressClassKey]
	if ingressAnnotationValue == "" &&
		s.ingressClass == annotations.DefaultIngressClass {
		return true
	}
	return ingressAnnotationValue == s.ingressClass
}

// ListKnativeIngresses returns the list of TCP Ingresses from
// configuration.konghq.com group.
func (s Store) ListKnativeIngresses() ([]*knative.Ingress, error) {
	var ingresses []*knative.Ingress
	if s.stores.KnativeIngress == nil {
		return ingresses, nil
	}
	err := cache.ListAll(s.stores.KnativeIngress, labels.NewSelector(),
		func(ob interface{}) {
			ing, ok := ob.(*knative.Ingress)
			if ok && s.validKnativeIngressClass(&ing.ObjectMeta) {
				ingresses = append(ingresses, ing)
			}
		})
	if err != nil {
		return nil, err
	}
	return ingresses, nil
}

// GetEndpointsForService returns the internal endpoints for service
// 'namespace/name' inside k8s.
func (s Store) GetEndpointsForService(namespace, name string) (*apiv1.Endpoints, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	eps, exists, err := s.stores.Endpoint.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("Endpoints for service %v not found", key)}
	}
	return eps.(*apiv1.Endpoints), nil
}

// GetKongPlugin returns the 'name' KongPlugin resource in namespace.
func (s Store) GetKongPlugin(namespace, name string) (*configurationv1.KongPlugin, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.stores.Plugin.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("KongPlugin %v not found", key)}
	}
	return p.(*configurationv1.KongPlugin), nil
}

// GetKongClusterPlugin returns the 'name' KongClusterPlugin resource.
func (s Store) GetKongClusterPlugin(name string) (*configurationv1.KongClusterPlugin, error) {
	p, exists, err := s.stores.ClusterPlugin.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("KongClusterPluign %v not found", name)}
	}
	return p.(*configurationv1.KongClusterPlugin), nil
}

// GetKongIngress returns the 'name' KongIngress resource in namespace.
func (s Store) GetKongIngress(namespace, name string) (*configurationv1.KongIngress, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.stores.Configuration.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("KongIngress %v not found", name)}
	}
	return p.(*configurationv1.KongIngress), nil
}

// GetKongConsumer returns the 'name' KongConsumer resource in namespace.
func (s Store) GetKongConsumer(namespace, name string) (*configurationv1.KongConsumer, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.stores.Consumer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("KongConsumer %v not found", key)}
	}
	return p.(*configurationv1.KongConsumer), nil
}

// ListKongConsumers returns all KongConsumers filtered by the ingress.class
// annotation.
func (s Store) ListKongConsumers() []*configurationv1.KongConsumer {
	var consumers []*configurationv1.KongConsumer
	for _, item := range s.stores.Consumer.List() {
		c, ok := item.(*configurationv1.KongConsumer)
		if ok && s.isValidIngresClass(&c.ObjectMeta) {
			consumers = append(consumers, c)
		}
	}

	return consumers
}

// ListKongCredentials returns all KongCredential filtered by the ingress.class
// annotation.
func (s Store) ListKongCredentials() []*configurationv1.KongCredential {
	var credentials []*configurationv1.KongCredential
	for _, item := range s.stores.Credential.List() {
		c, ok := item.(*configurationv1.KongCredential)
		if ok && s.isValidIngresClass(&c.ObjectMeta) {
			credentials = append(credentials, c)
		}
	}

	return credentials
}

// ListGlobalKongPlugins returns all KongPlugin resources
// filtered by the ingress.class annotation and with the
// label global:"true".
// Support for these global namespaced KongPlugins was removed in 0.10.0
// This function remains only to provide warnings to users with old configuration
func (s Store) ListGlobalKongPlugins() ([]*configurationv1.KongPlugin, error) {

	var plugins []*configurationv1.KongPlugin
	// var globalPlugins []*configurationv1.KongPlugin
	req, err := labels.NewRequirement("global", selection.Equals, []string{"true"})
	if err != nil {
		return nil, err
	}
	err = cache.ListAll(s.stores.Plugin,
		labels.NewSelector().Add(*req),
		func(ob interface{}) {
			p, ok := ob.(*configurationv1.KongPlugin)
			if ok && s.isValidIngresClass(&p.ObjectMeta) {
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
func (s Store) ListGlobalKongClusterPlugins() ([]*configurationv1.KongClusterPlugin, error) {

	var plugins []*configurationv1.KongClusterPlugin
	req, err := labels.NewRequirement("global", selection.Equals, []string{"true"})
	if err != nil {
		return nil, err
	}
	err = cache.ListAll(s.stores.ClusterPlugin,
		labels.NewSelector().Add(*req),
		func(ob interface{}) {
			p, ok := ob.(*configurationv1.KongClusterPlugin)
			if ok && s.isValidIngresClass(&p.ObjectMeta) {
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
func (s Store) ListCACerts() ([]*apiv1.Secret, error) {
	var secrets []*apiv1.Secret
	req, err := labels.NewRequirement(caCertKey,
		selection.Equals, []string{"true"})
	if err != nil {
		return nil, err
	}
	err = cache.ListAll(s.stores.Secret,
		labels.NewSelector().Add(*req),
		func(ob interface{}) {
			p, ok := ob.(*apiv1.Secret)
			if ok {
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

func toNetworkingIngressV1Beta1(obj *extensions.Ingress) (*networkingv1beta1.Ingress, error) {
	js, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize object of type %v: %v", reflect.TypeOf(obj), err)
	}
	var out networkingv1beta1.Ingress
	if err := json.Unmarshal(js, &out); err != nil {
		return nil, fmt.Errorf("failed to deserialize json: %v", err)
	}
	out.APIVersion = networkingv1beta1.SchemeGroupVersion.String()
	return &out, nil
}
