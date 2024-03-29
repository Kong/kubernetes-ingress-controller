/*
Copyright 2021 Kong, Inc.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
	scheme "github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// KongServiceFacadesGetter has a method to return a KongServiceFacadeInterface.
// A group's client should implement this interface.
type KongServiceFacadesGetter interface {
	KongServiceFacades(namespace string) KongServiceFacadeInterface
}

// KongServiceFacadeInterface has methods to work with KongServiceFacade resources.
type KongServiceFacadeInterface interface {
	Create(ctx context.Context, kongServiceFacade *v1alpha1.KongServiceFacade, opts v1.CreateOptions) (*v1alpha1.KongServiceFacade, error)
	Update(ctx context.Context, kongServiceFacade *v1alpha1.KongServiceFacade, opts v1.UpdateOptions) (*v1alpha1.KongServiceFacade, error)
	UpdateStatus(ctx context.Context, kongServiceFacade *v1alpha1.KongServiceFacade, opts v1.UpdateOptions) (*v1alpha1.KongServiceFacade, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.KongServiceFacade, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.KongServiceFacadeList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.KongServiceFacade, err error)
	KongServiceFacadeExpansion
}

// kongServiceFacades implements KongServiceFacadeInterface
type kongServiceFacades struct {
	client rest.Interface
	ns     string
}

// newKongServiceFacades returns a KongServiceFacades
func newKongServiceFacades(c *IncubatorV1alpha1Client, namespace string) *kongServiceFacades {
	return &kongServiceFacades{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the kongServiceFacade, and returns the corresponding kongServiceFacade object, and an error if there is any.
func (c *kongServiceFacades) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.KongServiceFacade, err error) {
	result = &v1alpha1.KongServiceFacade{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kongservicefacades").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of KongServiceFacades that match those selectors.
func (c *kongServiceFacades) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.KongServiceFacadeList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.KongServiceFacadeList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kongservicefacades").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested kongServiceFacades.
func (c *kongServiceFacades) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("kongservicefacades").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a kongServiceFacade and creates it.  Returns the server's representation of the kongServiceFacade, and an error, if there is any.
func (c *kongServiceFacades) Create(ctx context.Context, kongServiceFacade *v1alpha1.KongServiceFacade, opts v1.CreateOptions) (result *v1alpha1.KongServiceFacade, err error) {
	result = &v1alpha1.KongServiceFacade{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("kongservicefacades").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kongServiceFacade).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a kongServiceFacade and updates it. Returns the server's representation of the kongServiceFacade, and an error, if there is any.
func (c *kongServiceFacades) Update(ctx context.Context, kongServiceFacade *v1alpha1.KongServiceFacade, opts v1.UpdateOptions) (result *v1alpha1.KongServiceFacade, err error) {
	result = &v1alpha1.KongServiceFacade{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("kongservicefacades").
		Name(kongServiceFacade.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kongServiceFacade).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *kongServiceFacades) UpdateStatus(ctx context.Context, kongServiceFacade *v1alpha1.KongServiceFacade, opts v1.UpdateOptions) (result *v1alpha1.KongServiceFacade, err error) {
	result = &v1alpha1.KongServiceFacade{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("kongservicefacades").
		Name(kongServiceFacade.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kongServiceFacade).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the kongServiceFacade and deletes it. Returns an error if one occurs.
func (c *kongServiceFacades) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kongservicefacades").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *kongServiceFacades) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kongservicefacades").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched kongServiceFacade.
func (c *kongServiceFacades) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.KongServiceFacade, err error) {
	result = &v1alpha1.KongServiceFacade{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("kongservicefacades").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
