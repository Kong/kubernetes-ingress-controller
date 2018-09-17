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
	"reflect"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	"github.com/eapache/channels"
	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	consumerv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/consumer/v1"
	credentialv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/credential/v1"
	pluginv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/plugin/v1"
	configurationclientv1 "github.com/kong/kubernetes-ingress-controller/internal/client/configuration/clientset/versioned"
	configurationinformer "github.com/kong/kubernetes-ingress-controller/internal/client/configuration/informers/externalversions"
	consumerclientv1 "github.com/kong/kubernetes-ingress-controller/internal/client/consumer/clientset/versioned"
	consumerinformer "github.com/kong/kubernetes-ingress-controller/internal/client/consumer/informers/externalversions"
	credentialclientv1 "github.com/kong/kubernetes-ingress-controller/internal/client/credential/clientset/versioned"
	credentialinformer "github.com/kong/kubernetes-ingress-controller/internal/client/credential/informers/externalversions"
	pluginclientv1 "github.com/kong/kubernetes-ingress-controller/internal/client/plugin/clientset/versioned"
	plugininformer "github.com/kong/kubernetes-ingress-controller/internal/client/plugin/informers/externalversions"
	"github.com/kong/kubernetes-ingress-controller/internal/file"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations/class"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations/parser"
	"github.com/kong/kubernetes-ingress-controller/internal/k8s"
)

// Storer is the interface that wraps the required methods to gather information
// about ingresses, services, secrets and ingress annotations.
type Storer interface {
	// GetSecret returns a Secret using the namespace and name as key
	GetSecret(key string) (*corev1.Secret, error)

	// GetService returns a Service using the namespace and name as key
	GetService(key string) (*corev1.Service, error)

	GetServiceEndpoints(svc *corev1.Service) (*corev1.Endpoints, error)

	// GetSecret returns an Ingress using the namespace and name as key
	GetIngress(key string) (*extensions.Ingress, error)

	// ListIngresses returns the list of Ingresses
	ListIngresses() []*extensions.Ingress

	// GetLocalSSLCert returns the local copy of a SSLCert
	GetLocalSSLCert(name string) (*ingress.SSLCert, error)

	// ListLocalSSLCerts returns the list of local SSLCerts
	ListLocalSSLCerts() []*ingress.SSLCert

	// Run initiates the synchronization of the controllers
	Run(stopCh chan struct{})

	GetKongPlugin(namespace, name string) (*pluginv1.KongPlugin, error)

	GetKongConsumer(namespace, name string) (*consumerv1.KongConsumer, error)

	GetKongIngress(namespace, name string) (*configurationv1.KongIngress, error)

	ListKongConsumers() []*consumerv1.KongConsumer

	ListKongCredentials() []*credentialv1.KongCredential

	ListGlobalKongPlugins() ([]*pluginv1.KongPlugin, error)
}

// EventType type of event associated with an informer
type EventType string

const (
	// CreateEvent event associated with new objects in an informer
	CreateEvent EventType = "CREATE"
	// UpdateEvent event associated with an object update in an informer
	UpdateEvent EventType = "UPDATE"
	// DeleteEvent event associated when an object is removed from an informer
	DeleteEvent EventType = "DELETE"
	// ConfigurationEvent event associated when a configuration object is created or updated
	ConfigurationEvent EventType = "CONFIGURATION"
)

// Event holds the context of an event
type Event struct {
	Type EventType
	Obj  interface{}
	Old  interface{}
}

// Lister returns the stores for ingresses, services, endpoints and secrets.
type Lister struct {
	Ingress  IngressLister
	Service  ServiceLister
	Endpoint EndpointLister
	Secret   SecretLister

	Kong struct {
		Plugin        cache.Store
		Consumer      cache.Store
		Credential    cache.Store
		Configuration cache.Store
	}
}

// Informer defines the required SharedIndexInformers that interact with the API server.
type Informer struct {
	Ingress       cache.SharedIndexInformer
	Endpoint      cache.SharedIndexInformer
	Service       cache.SharedIndexInformer
	Secret        cache.SharedIndexInformer
	Configuration cache.SharedIndexInformer

	Kong struct {
		Plugin        cache.SharedIndexInformer
		Consumer      cache.SharedIndexInformer
		Credential    cache.SharedIndexInformer
		Configuration cache.SharedIndexInformer
	}
}

// Run initiates the synchronization of the controllers against the api server
func (c *Informer) Run(stopCh chan struct{}) {
	go c.Endpoint.Run(stopCh)
	go c.Service.Run(stopCh)
	go c.Secret.Run(stopCh)

	go c.Kong.Plugin.Run(stopCh)
	go c.Kong.Consumer.Run(stopCh)
	go c.Kong.Credential.Run(stopCh)
	go c.Kong.Configuration.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh,
		c.Endpoint.HasSynced,
		c.Service.HasSynced,
		c.Secret.HasSynced,
		c.Kong.Plugin.HasSynced,
		c.Kong.Consumer.HasSynced,
		c.Kong.Credential.HasSynced,
		c.Kong.Configuration.HasSynced,
	) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
	}

	// We need to wait before start syncing the ingress rules
	// because the rules requires content from other listers
	time.Sleep(1 * time.Second)
	go c.Ingress.Run(stopCh)
	if !cache.WaitForCacheSync(stopCh,
		c.Ingress.HasSynced,
	) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
	}
}

// k8sStore internal Storer implementation using informers and thread safe stores
type k8sStore struct {
	isOCSPCheckEnabled bool

	// informers contains the cache Informers
	informers *Informer

	// listers contains the cache.Store used in the ingress controller
	listers *Lister

	// sslStore local store of SSL certificates (certificates used in ingress)
	// this is required because the certificates must be present in the
	// container filesystem
	sslStore *SSLCertTracker

	// secretIngressMap contains information about which ingress references a
	// secret in the annotations.
	secretIngressMap ObjectRefMap

	filesystem file.Filesystem

	// updateCh
	updateCh *channels.RingChannel

	// mu mutex used to avoid simultaneous incovations to syncSecret
	mu *sync.Mutex
}

// New creates a new object store to be used in the ingress controller
func New(checkOCSP bool,
	namespace, configmap, tcp, udp, defaultSSLCertificate string,
	resyncPeriod time.Duration,
	client clientset.Interface,
	clientConf *rest.Config,
	fs file.Filesystem,
	updateCh *channels.RingChannel) Storer {

	store := &k8sStore{
		isOCSPCheckEnabled: checkOCSP,
		informers:          &Informer{},
		listers:            &Lister{},
		sslStore:           NewSSLCertTracker(),
		filesystem:         fs,
		updateCh:           updateCh,
		mu:                 &sync.Mutex{},
		secretIngressMap:   NewObjectRefMap(),
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{
		Interface: client.CoreV1().Events(namespace),
	})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{
		Component: "kong-ingress-controller",
	})

	ingEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ing := obj.(*extensions.Ingress)
			if !class.IsValid(ing) {
				a, _ := parser.GetStringAnnotation(class.IngressKey, ing)
				glog.Infof("ignoring add for ingress %v based on annotation %v with value %v", ing.Name, class.IngressKey, a)
				return
			}

			recorder.Eventf(ing, corev1.EventTypeNormal, "CREATE", fmt.Sprintf("Ingress %s/%s", ing.Namespace, ing.Name))

			store.updateSecretIngressMap(ing)
			store.syncSecrets(ing)

			updateCh.In() <- Event{
				Type: CreateEvent,
				Obj:  obj,
			}
		},
		DeleteFunc: func(obj interface{}) {
			ing, ok := obj.(*extensions.Ingress)
			if !ok {
				// If we reached here it means the ingress was deleted but its final state is unrecorded.
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					glog.Errorf("couldn't get object from tombstone %#v", obj)
					return
				}
				ing, ok = tombstone.Obj.(*extensions.Ingress)
				if !ok {
					glog.Errorf("Tombstone contained object that is not an Ingress: %#v", obj)
					return
				}
			}
			if !class.IsValid(ing) {
				glog.Infof("ignoring delete for ingress %v based on annotation %v", ing.Name, class.IngressKey)
				return
			}
			recorder.Eventf(ing, corev1.EventTypeNormal, "DELETE", fmt.Sprintf("Ingress %s/%s", ing.Namespace, ing.Name))

			key := k8s.MetaNamespaceKey(ing)
			store.secretIngressMap.Delete(key)

			updateCh.In() <- Event{
				Type: DeleteEvent,
				Obj:  obj,
			}
		},
		UpdateFunc: func(old, cur interface{}) {
			oldIng := old.(*extensions.Ingress)
			curIng := cur.(*extensions.Ingress)
			validOld := class.IsValid(oldIng)
			validCur := class.IsValid(curIng)
			if !validOld && validCur {
				glog.Infof("creating ingress %v based on annotation %v", curIng.Name, class.IngressKey)
				recorder.Eventf(curIng, corev1.EventTypeNormal, "CREATE", fmt.Sprintf("Ingress %s/%s", curIng.Namespace, curIng.Name))
			} else if validOld && !validCur {
				glog.Infof("removing ingress %v based on annotation %v", curIng.Name, class.IngressKey)
				recorder.Eventf(curIng, corev1.EventTypeNormal, "DELETE", fmt.Sprintf("Ingress %s/%s", curIng.Namespace, curIng.Name))
			} else if validCur && !reflect.DeepEqual(old, cur) {
				recorder.Eventf(curIng, corev1.EventTypeNormal, "UPDATE", fmt.Sprintf("Ingress %s/%s", curIng.Namespace, curIng.Name))
			}

			store.updateSecretIngressMap(curIng)
			store.syncSecrets(curIng)

			updateCh.In() <- Event{
				Type: UpdateEvent,
				Obj:  cur,
			}
		},
	}

	secrEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			sec := obj.(*corev1.Secret)
			key := k8s.MetaNamespaceKey(sec)

			// find references in ingresses and update local ssl certs
			if ings := store.secretIngressMap.Reference(key); len(ings) > 0 {
				glog.Infof("secret %v was added and it is used in ingress annotations. Parsing...", key)
				for _, ingKey := range ings {
					ing, err := store.GetIngress(ingKey)
					if err != nil {
						glog.Errorf("could not find Ingress %v in local store", ingKey)
						continue
					}
					store.syncSecrets(ing)
				}
				updateCh.In() <- Event{
					Type: CreateEvent,
					Obj:  obj,
				}
			}
		},
		UpdateFunc: func(old, cur interface{}) {
			if !reflect.DeepEqual(old, cur) {
				sec := cur.(*corev1.Secret)
				key := k8s.MetaNamespaceKey(sec)

				// find references in ingresses and update local ssl certs
				if ings := store.secretIngressMap.Reference(key); len(ings) > 0 {
					glog.Infof("secret %v was updated and it is used in ingress annotations. Parsing...", key)
					for _, ingKey := range ings {
						ing, err := store.GetIngress(ingKey)
						if err != nil {
							glog.Errorf("could not find Ingress %v in local store", ingKey)
							continue
						}
						store.syncSecrets(ing)
					}
					updateCh.In() <- Event{
						Type: UpdateEvent,
						Obj:  cur,
					}
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			sec, ok := obj.(*corev1.Secret)
			if !ok {
				// If we reached here it means the secret was deleted but its final state is unrecorded.
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					glog.Errorf("couldn't get object from tombstone %#v", obj)
					return
				}
				sec, ok = tombstone.Obj.(*corev1.Secret)
				if !ok {
					glog.Errorf("Tombstone contained object that is not a Secret: %#v", obj)
					return
				}
			}

			store.sslStore.Delete(k8s.MetaNamespaceKey(sec))

			updateCh.In() <- Event{
				Type: DeleteEvent,
				Obj:  obj,
			}
		},
	}

	epEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			updateCh.In() <- Event{
				Type: CreateEvent,
				Obj:  obj,
			}
		},
		DeleteFunc: func(obj interface{}) {
			updateCh.In() <- Event{
				Type: DeleteEvent,
				Obj:  obj,
			}
		},
		UpdateFunc: func(old, cur interface{}) {
			oep := old.(*corev1.Endpoints)
			ocur := cur.(*corev1.Endpoints)
			if !reflect.DeepEqual(ocur.Subsets, oep.Subsets) {
				updateCh.In() <- Event{
					Type: UpdateEvent,
					Obj:  cur,
				}
			}
		},
	}

	serviceEventHandler := cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, cur interface{}) {
			updateCh.In() <- Event{
				Type: ConfigurationEvent,
				Obj:  cur,
			}
		},
		DeleteFunc: func(obj interface{}) {
			updateCh.In() <- Event{
				Type: DeleteEvent,
				Obj:  obj,
			}
		},
	}

	// create informers factory, enable and assign required informers
	infFactory := informers.NewFilteredSharedInformerFactory(client, resyncPeriod, namespace, func(*metav1.ListOptions) {})

	store.informers.Ingress = infFactory.Extensions().V1beta1().Ingresses().Informer()
	store.listers.Ingress.Store = store.informers.Ingress.GetStore()

	store.informers.Endpoint = infFactory.Core().V1().Endpoints().Informer()
	store.listers.Endpoint.Store = store.informers.Endpoint.GetStore()

	store.informers.Secret = infFactory.Core().V1().Secrets().Informer()
	store.listers.Secret.Store = store.informers.Secret.GetStore()

	store.informers.Service = infFactory.Core().V1().Services().Informer()
	store.listers.Service.Store = store.informers.Service.GetStore()

	store.informers.Ingress.AddEventHandler(ingEventHandler)
	store.informers.Endpoint.AddEventHandler(epEventHandler)
	store.informers.Secret.AddEventHandler(secrEventHandler)
	store.informers.Service.AddEventHandler(serviceEventHandler)

	crdEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			updateCh.In() <- Event{
				Type: ConfigurationEvent,
				Obj:  obj,
			}
		},
		DeleteFunc: func(obj interface{}) {
			updateCh.In() <- Event{
				Type: ConfigurationEvent,
				Obj:  obj,
			}
		},
		UpdateFunc: func(old, cur interface{}) {
			updateCh.In() <- Event{
				Type: ConfigurationEvent,
				Obj:  cur,
			}
		},
	}

	pluginClient, _ := pluginclientv1.NewForConfig(clientConf)
	pluginFactory := plugininformer.NewFilteredSharedInformerFactory(pluginClient, resyncPeriod, namespace, func(*metav1.ListOptions) {})

	store.informers.Kong.Plugin = pluginFactory.Configuration().V1().KongPlugins().Informer()
	store.listers.Kong.Plugin = store.informers.Kong.Plugin.GetStore()
	store.informers.Kong.Plugin.AddEventHandler(crdEventHandler)

	consumerClient, _ := consumerclientv1.NewForConfig(clientConf)
	consumerFactory := consumerinformer.NewFilteredSharedInformerFactory(consumerClient, resyncPeriod, namespace, func(*metav1.ListOptions) {})

	store.informers.Kong.Consumer = consumerFactory.Configuration().V1().KongConsumers().Informer()
	store.listers.Kong.Consumer = store.informers.Kong.Consumer.GetStore()
	store.informers.Kong.Consumer.AddEventHandler(crdEventHandler)

	credClient, _ := credentialclientv1.NewForConfig(clientConf)
	credentialFactory := credentialinformer.NewFilteredSharedInformerFactory(credClient, resyncPeriod, namespace, func(*metav1.ListOptions) {})

	store.informers.Kong.Credential = credentialFactory.Configuration().V1().KongCredentials().Informer()
	store.listers.Kong.Credential = store.informers.Kong.Credential.GetStore()
	store.informers.Kong.Credential.AddEventHandler(crdEventHandler)

	confClient, _ := configurationclientv1.NewForConfig(clientConf)
	configFactory := configurationinformer.NewFilteredSharedInformerFactory(confClient, resyncPeriod, namespace, func(*metav1.ListOptions) {})

	store.informers.Kong.Configuration = configFactory.Configuration().V1().KongIngresses().Informer()
	store.listers.Kong.Configuration = store.informers.Kong.Configuration.GetStore()
	store.informers.Kong.Configuration.AddEventHandler(crdEventHandler)

	return store
}

// GetSecret returns a Secret using the namespace and name as key
func (s k8sStore) GetSecret(key string) (*corev1.Secret, error) {
	return s.listers.Secret.ByKey(key)
}

// ListLocalSecrets returns the list of local Secrets
func (s k8sStore) ListLocalSecrets() []*ingress.SSLCert {
	var certs []*ingress.SSLCert
	for _, item := range s.sslStore.List() {
		if s, ok := item.(*ingress.SSLCert); ok {
			certs = append(certs, s)
		}
	}

	return certs
}

// GetService returns a Service using the namespace and name as key
func (s k8sStore) GetService(key string) (*corev1.Service, error) {
	return s.listers.Service.ByKey(key)
}

// GetSecret returns an Ingress using the namespace and name as key
func (s k8sStore) GetIngress(key string) (*extensions.Ingress, error) {
	return s.listers.Ingress.ByKey(key)
}

// ListIngresses returns the list of Ingresses
func (s k8sStore) ListIngresses() []*extensions.Ingress {
	// filter ingress rules
	var ingresses []*extensions.Ingress
	for _, item := range s.listers.Ingress.List() {
		ing := item.(*extensions.Ingress)
		if !class.IsValid(ing) {
			continue
		}

		ingresses = append(ingresses, ing)
	}

	return ingresses
}

func (s k8sStore) GetServiceEndpoints(svc *corev1.Service) (*corev1.Endpoints, error) {
	return s.listers.Endpoint.GetServiceEndpoints(svc)
}

// Run initiates the synchronization of the controllers
// and the initial synchronization of the secrets.
func (s k8sStore) Run(stopCh chan struct{}) {
	// start informers
	s.informers.Run(stopCh)

	if s.isOCSPCheckEnabled {
		go wait.Until(s.checkSSLChainIssues, 60*time.Second, stopCh)
	}
}

// syncSecrets synchronizes data from all Secrets referenced by the given
// Ingress with the local store and file system.
func (s k8sStore) syncSecrets(ing *extensions.Ingress) {
	key := k8s.MetaNamespaceKey(ing)
	for _, secrKey := range s.secretIngressMap.ReferencedBy(key) {
		s.syncSecret(secrKey)
	}
}

// updateSecretIngressMap takes an Ingress and updates all Secret objects it
// references in secretIngressMap.
func (s *k8sStore) updateSecretIngressMap(ing *extensions.Ingress) {
	key := k8s.MetaNamespaceKey(ing)
	glog.V(3).Infof("updating references to secrets for ingress %v", key)

	// delete all existing references first
	s.secretIngressMap.Delete(key)

	var refSecrets []string

	for _, tls := range ing.Spec.TLS {
		secrName := tls.SecretName
		if secrName != "" {
			secrKey := fmt.Sprintf("%v/%v", ing.Namespace, secrName)
			refSecrets = append(refSecrets, secrKey)
		}
	}

	// populate map with all secret references
	s.secretIngressMap.Insert(key, refSecrets...)
}

// ListLocalSSLCerts returns the list of local SSLCerts
func (s k8sStore) ListLocalSSLCerts() []*ingress.SSLCert {
	var certs []*ingress.SSLCert
	for _, item := range s.sslStore.List() {
		if s, ok := item.(*ingress.SSLCert); ok {
			certs = append(certs, s)
		}
	}

	return certs
}

// GetLocalSSLCert returns the local copy of a SSLCert
func (s k8sStore) GetLocalSSLCert(key string) (*ingress.SSLCert, error) {
	return s.sslStore.ByKey(key)
}

func (s k8sStore) GetKongPlugin(namespace, name string) (*pluginv1.KongPlugin, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.listers.Kong.Plugin.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("plugin %v was not found", key)
	}
	return p.(*pluginv1.KongPlugin), nil
}

func (s k8sStore) GetKongIngress(namespace, name string) (*configurationv1.KongIngress, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.listers.Kong.Configuration.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("ingress configuration %v was not found", key)
	}
	return p.(*configurationv1.KongIngress), nil
}

func (s k8sStore) GetKongConsumer(namespace, name string) (*consumerv1.KongConsumer, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	p, exists, err := s.listers.Kong.Consumer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("plugin %v was not found", key)
	}
	return p.(*consumerv1.KongConsumer), nil
}

func (s k8sStore) ListKongConsumers() []*consumerv1.KongConsumer {
	var consumers []*consumerv1.KongConsumer
	for _, item := range s.listers.Kong.Consumer.List() {
		if c, ok := item.(*consumerv1.KongConsumer); ok {
			consumers = append(consumers, c)
		}
	}

	return consumers
}

func (s k8sStore) ListKongCredentials() []*credentialv1.KongCredential {
	var credentials []*credentialv1.KongCredential
	for _, item := range s.listers.Kong.Credential.List() {
		if c, ok := item.(*credentialv1.KongCredential); ok {
			credentials = append(credentials, c)
		}
	}

	return credentials
}

func (s k8sStore) ListGlobalKongPlugins() ([]*pluginv1.KongPlugin, error) {

	var plugins []*pluginv1.KongPlugin
	// var globalPlugins []*pluginv1.KongPlugin
	req, err := labels.NewRequirement("global", selection.Equals, []string{"true"})
	if err != nil {
		return nil, err
	}
	err = cache.ListAll(s.listers.Kong.Plugin,
		labels.NewSelector().Add(*req),
		func(ob interface{}) {
			if p, ok := ob.(*pluginv1.KongPlugin); ok {
				plugins = append(plugins, p)
			}
		})
	if err != nil {
		return nil, err
	}
	return plugins, nil
}
