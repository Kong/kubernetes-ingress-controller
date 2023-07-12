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

	v1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeKongConsumerGroups implements KongConsumerGroupInterface
type FakeKongConsumerGroups struct {
	Fake *FakeConfigurationV1beta1
	ns   string
}

var kongconsumergroupsResource = v1beta1.SchemeGroupVersion.WithResource("kongconsumergroups")

var kongconsumergroupsKind = v1beta1.SchemeGroupVersion.WithKind("KongConsumerGroup")

// Get takes name of the kongConsumerGroup, and returns the corresponding kongConsumerGroup object, and an error if there is any.
func (c *FakeKongConsumerGroups) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.KongConsumerGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(kongconsumergroupsResource, c.ns, name), &v1beta1.KongConsumerGroup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.KongConsumerGroup), err
}

// List takes label and field selectors, and returns the list of KongConsumerGroups that match those selectors.
func (c *FakeKongConsumerGroups) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.KongConsumerGroupList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(kongconsumergroupsResource, kongconsumergroupsKind, c.ns, opts), &v1beta1.KongConsumerGroupList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.KongConsumerGroupList{ListMeta: obj.(*v1beta1.KongConsumerGroupList).ListMeta}
	for _, item := range obj.(*v1beta1.KongConsumerGroupList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested kongConsumerGroups.
func (c *FakeKongConsumerGroups) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(kongconsumergroupsResource, c.ns, opts))

}

// Create takes the representation of a kongConsumerGroup and creates it.  Returns the server's representation of the kongConsumerGroup, and an error, if there is any.
func (c *FakeKongConsumerGroups) Create(ctx context.Context, kongConsumerGroup *v1beta1.KongConsumerGroup, opts v1.CreateOptions) (result *v1beta1.KongConsumerGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(kongconsumergroupsResource, c.ns, kongConsumerGroup), &v1beta1.KongConsumerGroup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.KongConsumerGroup), err
}

// Update takes the representation of a kongConsumerGroup and updates it. Returns the server's representation of the kongConsumerGroup, and an error, if there is any.
func (c *FakeKongConsumerGroups) Update(ctx context.Context, kongConsumerGroup *v1beta1.KongConsumerGroup, opts v1.UpdateOptions) (result *v1beta1.KongConsumerGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(kongconsumergroupsResource, c.ns, kongConsumerGroup), &v1beta1.KongConsumerGroup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.KongConsumerGroup), err
}

// Delete takes name of the kongConsumerGroup and deletes it. Returns an error if one occurs.
func (c *FakeKongConsumerGroups) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(kongconsumergroupsResource, c.ns, name, opts), &v1beta1.KongConsumerGroup{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeKongConsumerGroups) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(kongconsumergroupsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.KongConsumerGroupList{})
	return err
}

// Patch applies the patch and returns the patched kongConsumerGroup.
func (c *FakeKongConsumerGroups) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.KongConsumerGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(kongconsumergroupsResource, c.ns, name, pt, data, subresources...), &v1beta1.KongConsumerGroup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.KongConsumerGroup), err
}
