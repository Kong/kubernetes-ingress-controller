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
	"bytes"
	"fmt"
	"strings"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/tools/cache"
)

// Storer is the interface that wraps the required methods to gather information
// about ingresses, services, secrets and ingress annotations.
type Storer interface {
	GetSecret(key string) (*apiv1.Secret, error)
	GetCertFromSecret(string) (*utils.RawSSLCert, error)
	GetService(key string) (*apiv1.Service, error)
	GetServiceEndpoints(svc *apiv1.Service) (*apiv1.Endpoints, error)
	ListIngresses() []*extensions.Ingress
	GetKongIngress(namespace, name string) (*configurationv1.KongIngress, error)
	GetKongConsumer(namespace, name string) (*configurationv1.KongConsumer, error)
	ListKongConsumers() []*configurationv1.KongConsumer
	GetKongPlugin(namespace, name string) (*configurationv1.KongPlugin, error)
	ListGlobalKongPlugins() ([]*configurationv1.KongPlugin, error)
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
func (s Store) GetSecret(key string) (*apiv1.Secret, error) {
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
func (s Store) GetService(key string) (*apiv1.Service, error) {
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
func (s Store) ListIngresses() []*extensions.Ingress {
	// filter ingress rules
	var ingresses []*extensions.Ingress
	for _, item := range s.stores.Ingress.List() {
		ing := item.(*extensions.Ingress)
		if !s.isValidIngresClass(&ing.ObjectMeta) {
			continue
		}

		ingresses = append(ingresses, ing)
	}

	return ingresses
}

// GetServiceEndpoints returns the internal endpoints for svc inside the
// current k8s cluster.
func (s Store) GetServiceEndpoints(svc *apiv1.Service) (*apiv1.Endpoints, error) {
	key := fmt.Sprintf("%v/%v", svc.Namespace, svc.Name)
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
		return nil, nil
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

// GetCertFromSecret gets an SSL cert from k8s secret.
func (s Store) GetCertFromSecret(secretName string) (*utils.RawSSLCert, error) {
	secret, err := s.GetSecret(secretName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving secret %v: %v", secretName, err)
	}
	cert, okcert := secret.Data[apiv1.TLSCertKey]
	key, okkey := secret.Data[apiv1.TLSPrivateKeyKey]

	if !okcert || !okkey {
		return nil, fmt.Errorf("no keypair could be found in %v", secretName)
	}

	cert = []byte(strings.TrimSpace(bytes.NewBuffer(cert).String()))
	key = []byte(strings.TrimSpace(bytes.NewBuffer(key).String()))

	return &utils.RawSSLCert{
		Cert: cert,
		Key:  key,
	}, nil
}
