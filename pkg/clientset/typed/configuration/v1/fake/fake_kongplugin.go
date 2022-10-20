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

package fake

import (
	"context"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeKongPlugins implements KongPluginInterface
type FakeKongPlugins struct {
	Fake *FakeConfigurationV1
	ns   string
}

var kongpluginsResource = schema.GroupVersionResource{Group: "configuration", Version: "v1", Resource: "kongplugins"}

var kongpluginsKind = schema.GroupVersionKind{Group: "configuration", Version: "v1", Kind: "KongPlugin"}

// Get takes name of the kongPlugin, and returns the corresponding kongPlugin object, and an error if there is any.
func (c *FakeKongPlugins) Get(ctx context.Context, name string, options v1.GetOptions) (result *configurationv1.KongPlugin, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(kongpluginsResource, c.ns, name), &configurationv1.KongPlugin{})

	if obj == nil {
		return nil, err
	}
	return obj.(*configurationv1.KongPlugin), err
}

// List takes label and field selectors, and returns the list of KongPlugins that match those selectors.
func (c *FakeKongPlugins) List(ctx context.Context, opts v1.ListOptions) (result *configurationv1.KongPluginList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(kongpluginsResource, kongpluginsKind, c.ns, opts), &configurationv1.KongPluginList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &configurationv1.KongPluginList{ListMeta: obj.(*configurationv1.KongPluginList).ListMeta}
	for _, item := range obj.(*configurationv1.KongPluginList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested kongPlugins.
func (c *FakeKongPlugins) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(kongpluginsResource, c.ns, opts))

}

// Create takes the representation of a kongPlugin and creates it.  Returns the server's representation of the kongPlugin, and an error, if there is any.
func (c *FakeKongPlugins) Create(ctx context.Context, kongPlugin *configurationv1.KongPlugin, opts v1.CreateOptions) (result *configurationv1.KongPlugin, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(kongpluginsResource, c.ns, kongPlugin), &configurationv1.KongPlugin{})

	if obj == nil {
		return nil, err
	}
	return obj.(*configurationv1.KongPlugin), err
}

// Update takes the representation of a kongPlugin and updates it. Returns the server's representation of the kongPlugin, and an error, if there is any.
func (c *FakeKongPlugins) Update(ctx context.Context, kongPlugin *configurationv1.KongPlugin, opts v1.UpdateOptions) (result *configurationv1.KongPlugin, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(kongpluginsResource, c.ns, kongPlugin), &configurationv1.KongPlugin{})

	if obj == nil {
		return nil, err
	}
	return obj.(*configurationv1.KongPlugin), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeKongPlugins) UpdateStatus(ctx context.Context, kongPlugin *configurationv1.KongPlugin, opts v1.UpdateOptions) (*configurationv1.KongPlugin, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(kongpluginsResource, "status", c.ns, kongPlugin), &configurationv1.KongPlugin{})

	if obj == nil {
		return nil, err
	}
	return obj.(*configurationv1.KongPlugin), err
}

// Delete takes name of the kongPlugin and deletes it. Returns an error if one occurs.
func (c *FakeKongPlugins) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(kongpluginsResource, c.ns, name, opts), &configurationv1.KongPlugin{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeKongPlugins) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(kongpluginsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &configurationv1.KongPluginList{})
	return err
}

// Patch applies the patch and returns the patched kongPlugin.
func (c *FakeKongPlugins) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *configurationv1.KongPlugin, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(kongpluginsResource, c.ns, name, pt, data, subresources...), &configurationv1.KongPlugin{})

	if obj == nil {
		return nil, err
	}
	return obj.(*configurationv1.KongPlugin), err
}
