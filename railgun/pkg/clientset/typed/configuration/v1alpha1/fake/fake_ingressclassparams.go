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

	v1alpha1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeIngressClassParamses implements IngressClassParamsInterface
type FakeIngressClassParamses struct {
	Fake *FakeConfigurationV1alpha1
	ns   string
}

var ingressclassparamsesResource = schema.GroupVersionResource{Group: "configuration", Version: "v1alpha1", Resource: "ingressclassparamses"}

var ingressclassparamsesKind = schema.GroupVersionKind{Group: "configuration", Version: "v1alpha1", Kind: "IngressClassParams"}

// Get takes name of the ingressClassParams, and returns the corresponding ingressClassParams object, and an error if there is any.
func (c *FakeIngressClassParamses) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.IngressClassParams, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(ingressclassparamsesResource, c.ns, name), &v1alpha1.IngressClassParams{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IngressClassParams), err
}

// List takes label and field selectors, and returns the list of IngressClassParamses that match those selectors.
func (c *FakeIngressClassParamses) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.IngressClassParamsList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(ingressclassparamsesResource, ingressclassparamsesKind, c.ns, opts), &v1alpha1.IngressClassParamsList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.IngressClassParamsList{ListMeta: obj.(*v1alpha1.IngressClassParamsList).ListMeta}
	for _, item := range obj.(*v1alpha1.IngressClassParamsList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested ingressClassParamses.
func (c *FakeIngressClassParamses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(ingressclassparamsesResource, c.ns, opts))

}

// Create takes the representation of a ingressClassParams and creates it.  Returns the server's representation of the ingressClassParams, and an error, if there is any.
func (c *FakeIngressClassParamses) Create(ctx context.Context, ingressClassParams *v1alpha1.IngressClassParams, opts v1.CreateOptions) (result *v1alpha1.IngressClassParams, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(ingressclassparamsesResource, c.ns, ingressClassParams), &v1alpha1.IngressClassParams{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IngressClassParams), err
}

// Update takes the representation of a ingressClassParams and updates it. Returns the server's representation of the ingressClassParams, and an error, if there is any.
func (c *FakeIngressClassParamses) Update(ctx context.Context, ingressClassParams *v1alpha1.IngressClassParams, opts v1.UpdateOptions) (result *v1alpha1.IngressClassParams, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(ingressclassparamsesResource, c.ns, ingressClassParams), &v1alpha1.IngressClassParams{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IngressClassParams), err
}

// Delete takes name of the ingressClassParams and deletes it. Returns an error if one occurs.
func (c *FakeIngressClassParamses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(ingressclassparamsesResource, c.ns, name), &v1alpha1.IngressClassParams{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeIngressClassParamses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(ingressclassparamsesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.IngressClassParamsList{})
	return err
}

// Patch applies the patch and returns the patched ingressClassParams.
func (c *FakeIngressClassParamses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.IngressClassParams, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(ingressclassparamsesResource, c.ns, name, pt, data, subresources...), &v1alpha1.IngressClassParams{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IngressClassParams), err
}
