/*
Copyright 2018 The Kong Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	versioned "github.com/kong/kubernetes-ingress-controller/pkg/client/configuration/clientset/versioned"
	internalinterfaces "github.com/kong/kubernetes-ingress-controller/pkg/client/configuration/informers/externalversions/internalinterfaces"
	v1 "github.com/kong/kubernetes-ingress-controller/pkg/client/configuration/listers/configuration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// KongCredentialInformer provides access to a shared informer and lister for
// KongCredentials.
type KongCredentialInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.KongCredentialLister
}

type kongCredentialInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewKongCredentialInformer constructs a new informer for KongCredential type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewKongCredentialInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredKongCredentialInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredKongCredentialInformer constructs a new informer for KongCredential type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredKongCredentialInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigurationV1().KongCredentials(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigurationV1().KongCredentials(namespace).Watch(context.TODO(), options)
			},
		},
		&configurationv1.KongCredential{},
		resyncPeriod,
		indexers,
	)
}

func (f *kongCredentialInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredKongCredentialInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *kongCredentialInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&configurationv1.KongCredential{}, f.defaultInformer)
}

func (f *kongCredentialInformer) Lister() v1.KongCredentialLister {
	return v1.NewKongCredentialLister(f.Informer().GetIndexer())
}
