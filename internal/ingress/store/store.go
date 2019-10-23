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
	"strings"

	"github.com/golang/glog"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/tools/cache"
)

// Storer is the interface that wraps the required methods to gather information
// about ingresses, services, secrets and ingress annotations.
type Storer interface {
	GetSecret(namespace, name string) (*apiv1.Secret, error)
	GetService(namespace, name string) (*apiv1.Service, error)
	GetEndpointsForService(namespace, name string) (*apiv1.Endpoints, error)
	GetKongIngress(namespace, name string) (*configurationv1.KongIngress, error)
	GetKongPlugin(namespace, name string) (*configurationv1.KongPlugin, error)
	GetKongConsumer(namespace, name string) (*configurationv1.KongConsumer, error)

	ListIngresses() []*networking.Ingress
	ListGlobalKongPlugins() ([]*configurationv1.KongPlugin, error)
	ListKongConsumers() []*configurationv1.KongConsumer
	ListKongCredentials() []*configurationv1.KongCredential
}

// Store implements Storer and can be used to list Ingress, Services
// and other resources from k8s APIserver. The backing stores should
// be synced and updated by the caller.
// It is ingressClass filter aware.
type Store struct {
	stores CacheStores

	isValidIngresClass func(objectMeta *metav1.ObjectMeta) bool
}

// CacheStores stores cache.Store for all Kinds of k8s objects that
// the Ingress Controller reads.
type CacheStores struct {
	Ingress  cache.Store
	Service  cache.Store
	Secret   cache.Store
	Endpoint cache.Store

	Plugin        cache.Store
	Consumer      cache.Store
	Credential    cache.Store
	Configuration cache.Store
}

// New creates a new object store to be used in the ingress controller
func New(cs CacheStores,
	isValidIngresClassFunc func(objectMeta *metav1.ObjectMeta) bool) Storer {
	return Store{
		stores:             cs,
		isValidIngresClass: isValidIngresClassFunc,
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
		return nil, fmt.Errorf("secret %v was not found", key)
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
		return nil, fmt.Errorf("service %v was not found", key)
	}
	return service.(*apiv1.Service), nil
}

// ListIngresses returns the list of Ingresses
func (s Store) ListIngresses() []*networking.Ingress {
	// filter ingress rules
	var ingresses []*networking.Ingress
	for _, item := range s.stores.Ingress.List() {
		ing := networkingIngressV1Beta1(item)
		if !s.isValidIngresClass(&ing.ObjectMeta) {
			continue
		}
		ingresses = append(ingresses, ing)
	}

	return ingresses
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
		return nil, fmt.Errorf("could not find endpoints for service %v", key)
	}
	return eps.(*apiv1.Endpoints), nil
}

// GetKongPlugin returns the 'name' KongPlugin resource in namespace.
func (s Store) GetKongPlugin(namespace, name string) (*configurationv1.KongPlugin, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	if strings.Contains(name, "/") {
		key = name
	}
	p, exists, err := s.stores.Plugin.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("plugin %v was not found", key)
	}
	return p.(*configurationv1.KongPlugin), nil
}

// GetKongIngress returns the 'name' KongIngress resource in namespace.
func (s Store) GetKongIngress(namespace, name string) (*configurationv1.KongIngress, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.stores.Configuration.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("KongIngress %v was not found", key)
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
		return nil, fmt.Errorf("consumer %v was not found", key)
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

var ingressConversionScheme *runtime.Scheme

func init() {
	ingressConversionScheme = runtime.NewScheme()
	extensions.AddToScheme(ingressConversionScheme)
	networking.AddToScheme(ingressConversionScheme)
}

func networkingIngressV1Beta1(obj interface{}) *networking.Ingress {
	networkingIngress, okNetworking := obj.(*networking.Ingress)
	if okNetworking {
		return networkingIngress
	}
	extensionsIngress, okExtension := obj.(*extensions.Ingress)
	if !okExtension {
		glog.Errorf("ingress resource can not be casted to extensions.Ingress" +
			" or networking.Ingress")
		return nil
	}
	networkingIngress = &networking.Ingress{}
	err := ingressConversionScheme.Convert(extensionsIngress, networkingIngress, nil)
	if err != nil {
		glog.Error("failed to convert extensions.Ingress "+
			"to networking.Ingress", err)
		return nil
	}
	return networkingIngress
}
