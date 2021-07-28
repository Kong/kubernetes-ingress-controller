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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/kong/kubernetes-ingress-controller/apis/configuration/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// KongPluginLister helps list KongPlugins.
// All objects returned here must be treated as read-only.
type KongPluginLister interface {
	// List lists all KongPlugins in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.KongPlugin, err error)
	// KongPlugins returns an object that can list and get KongPlugins.
	KongPlugins(namespace string) KongPluginNamespaceLister
	KongPluginListerExpansion
}

// kongPluginLister implements the KongPluginLister interface.
type kongPluginLister struct {
	indexer cache.Indexer
}

// NewKongPluginLister returns a new KongPluginLister.
func NewKongPluginLister(indexer cache.Indexer) KongPluginLister {
	return &kongPluginLister{indexer: indexer}
}

// List lists all KongPlugins in the indexer.
func (s *kongPluginLister) List(selector labels.Selector) (ret []*v1.KongPlugin, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.KongPlugin))
	})
	return ret, err
}

// KongPlugins returns an object that can list and get KongPlugins.
func (s *kongPluginLister) KongPlugins(namespace string) KongPluginNamespaceLister {
	return kongPluginNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// KongPluginNamespaceLister helps list and get KongPlugins.
// All objects returned here must be treated as read-only.
type KongPluginNamespaceLister interface {
	// List lists all KongPlugins in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.KongPlugin, err error)
	// Get retrieves the KongPlugin from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.KongPlugin, error)
	KongPluginNamespaceListerExpansion
}

// kongPluginNamespaceLister implements the KongPluginNamespaceLister
// interface.
type kongPluginNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all KongPlugins in the indexer for a given namespace.
func (s kongPluginNamespaceLister) List(selector labels.Selector) (ret []*v1.KongPlugin, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.KongPlugin))
	})
	return ret, err
}

// Get retrieves the KongPlugin from the indexer for a given namespace and name.
func (s kongPluginNamespaceLister) Get(name string) (*v1.KongPlugin, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("kongplugin"), name)
	}
	return obj.(*v1.KongPlugin), nil
}
