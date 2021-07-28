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

package v1beta1

import (
	v1beta1 "github.com/kong/kubernetes-ingress-controller/apis/configuration/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// TCPIngressLister helps list TCPIngresses.
// All objects returned here must be treated as read-only.
type TCPIngressLister interface {
	// List lists all TCPIngresses in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.TCPIngress, err error)
	// TCPIngresses returns an object that can list and get TCPIngresses.
	TCPIngresses(namespace string) TCPIngressNamespaceLister
	TCPIngressListerExpansion
}

// tCPIngressLister implements the TCPIngressLister interface.
type tCPIngressLister struct {
	indexer cache.Indexer
}

// NewTCPIngressLister returns a new TCPIngressLister.
func NewTCPIngressLister(indexer cache.Indexer) TCPIngressLister {
	return &tCPIngressLister{indexer: indexer}
}

// List lists all TCPIngresses in the indexer.
func (s *tCPIngressLister) List(selector labels.Selector) (ret []*v1beta1.TCPIngress, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.TCPIngress))
	})
	return ret, err
}

// TCPIngresses returns an object that can list and get TCPIngresses.
func (s *tCPIngressLister) TCPIngresses(namespace string) TCPIngressNamespaceLister {
	return tCPIngressNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TCPIngressNamespaceLister helps list and get TCPIngresses.
// All objects returned here must be treated as read-only.
type TCPIngressNamespaceLister interface {
	// List lists all TCPIngresses in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.TCPIngress, err error)
	// Get retrieves the TCPIngress from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.TCPIngress, error)
	TCPIngressNamespaceListerExpansion
}

// tCPIngressNamespaceLister implements the TCPIngressNamespaceLister
// interface.
type tCPIngressNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all TCPIngresses in the indexer for a given namespace.
func (s tCPIngressNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.TCPIngress, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.TCPIngress))
	})
	return ret, err
}

// Get retrieves the TCPIngress from the indexer for a given namespace and name.
func (s tCPIngressNamespaceLister) Get(name string) (*v1beta1.TCPIngress, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("tcpingress"), name)
	}
	return obj.(*v1beta1.TCPIngress), nil
}
