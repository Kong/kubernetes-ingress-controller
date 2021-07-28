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

package v1

import (
	"context"
	"time"

	v1 "github.com/kong/kubernetes-ingress-controller/apis/configuration/v1"
	scheme "github.com/kong/kubernetes-ingress-controller/pkg/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// KongPluginsGetter has a method to return a KongPluginInterface.
// A group's client should implement this interface.
type KongPluginsGetter interface {
	KongPlugins(namespace string) KongPluginInterface
}

// KongPluginInterface has methods to work with KongPlugin resources.
type KongPluginInterface interface {
	Create(ctx context.Context, kongPlugin *v1.KongPlugin, opts metav1.CreateOptions) (*v1.KongPlugin, error)
	Update(ctx context.Context, kongPlugin *v1.KongPlugin, opts metav1.UpdateOptions) (*v1.KongPlugin, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.KongPlugin, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.KongPluginList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.KongPlugin, err error)
	KongPluginExpansion
}

// kongPlugins implements KongPluginInterface
type kongPlugins struct {
	client rest.Interface
	ns     string
}

// newKongPlugins returns a KongPlugins
func newKongPlugins(c *ConfigurationV1Client, namespace string) *kongPlugins {
	return &kongPlugins{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the kongPlugin, and returns the corresponding kongPlugin object, and an error if there is any.
func (c *kongPlugins) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.KongPlugin, err error) {
	result = &v1.KongPlugin{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kongplugins").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of KongPlugins that match those selectors.
func (c *kongPlugins) List(ctx context.Context, opts metav1.ListOptions) (result *v1.KongPluginList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.KongPluginList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kongplugins").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested kongPlugins.
func (c *kongPlugins) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("kongplugins").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a kongPlugin and creates it.  Returns the server's representation of the kongPlugin, and an error, if there is any.
func (c *kongPlugins) Create(ctx context.Context, kongPlugin *v1.KongPlugin, opts metav1.CreateOptions) (result *v1.KongPlugin, err error) {
	result = &v1.KongPlugin{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("kongplugins").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kongPlugin).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a kongPlugin and updates it. Returns the server's representation of the kongPlugin, and an error, if there is any.
func (c *kongPlugins) Update(ctx context.Context, kongPlugin *v1.KongPlugin, opts metav1.UpdateOptions) (result *v1.KongPlugin, err error) {
	result = &v1.KongPlugin{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("kongplugins").
		Name(kongPlugin.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kongPlugin).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the kongPlugin and deletes it. Returns an error if one occurs.
func (c *kongPlugins) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kongplugins").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *kongPlugins) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kongplugins").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched kongPlugin.
func (c *kongPlugins) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.KongPlugin, err error) {
	result = &v1.KongPlugin{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("kongplugins").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
