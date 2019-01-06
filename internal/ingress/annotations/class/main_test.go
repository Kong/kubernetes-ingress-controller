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

package class

import (
	"testing"

	consumerv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/consumer/v1"
	credentialv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/credential/v1"
	pluginv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/plugin/v1"
	api "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsValidClass(t *testing.T) {
	dc := DefaultClass
	ic := IngressClass
	// restore original values after the tests
	defer func() {
		DefaultClass = dc
		IngressClass = ic
	}()

	tests := []struct {
		ingress    string
		controller string
		defClass   string
		isValid    bool
	}{
		{"", "", "nginx", true},
		{"", "nginx", "nginx", true},
		{"nginx", "nginx", "nginx", true},
		{"custom", "custom", "nginx", true},
		{"", "killer", "nginx", false},
		{"custom", "nginx", "nginx", false},
	}

	ing := &extensions.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
	}

	data := map[string]string{}
	ing.SetAnnotations(data)
	for _, test := range tests {
		ing.Annotations[IngressKey] = test.ingress

		IngressClass = test.controller
		DefaultClass = test.defClass

		b := IsValid(&ing.ObjectMeta)
		if b != test.isValid {
			t.Errorf("test %v - expected %v but %v was returned", test, test.isValid, b)
		}
	}
}

func TestCandAddResource(t *testing.T) {
	dc := DefaultClass
	ic := IngressClass
	// restore original values after the tests
	defer func() {
		DefaultClass = dc
		IngressClass = ic
	}()

	tests := []struct {
		ingress    string
		controller string
		defClass   string
		isValid    bool
	}{
		{"", "", "nginx", true},
		{"", "nginx", "nginx", true},
		{"nginx", "nginx", "nginx", true},
		{"custom", "custom", "nginx", true},
		{"", "killer", "nginx", false},
		{"custom", "nginx", "nginx", false},
	}

	crd := &pluginv1.KongPlugin{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
	}

	data := map[string]string{}
	crd.SetAnnotations(data)
	for _, test := range tests {
		crd.Annotations[IngressKey] = test.ingress

		IngressClass = test.controller
		DefaultClass = test.defClass

		b := CanAddResource(crd)
		if b != test.isValid {
			t.Errorf("test %v - expected %v but %v was returned", test, test.isValid, b)
		}
	}
}

func TestCandDeleteResource(t *testing.T) {
	dc := DefaultClass
	ic := IngressClass
	// restore original values after the tests
	defer func() {
		DefaultClass = dc
		IngressClass = ic
	}()

	tests := []struct {
		ingress    string
		controller string
		defClass   string
		isValid    bool
	}{
		{"", "", "nginx", true},
		{"", "nginx", "nginx", true},
		{"nginx", "nginx", "nginx", true},
		{"custom", "custom", "nginx", true},
		{"", "killer", "nginx", false},
		{"custom", "nginx", "nginx", false},
	}

	crd := &credentialv1.KongCredential{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
	}

	data := map[string]string{}
	crd.SetAnnotations(data)
	for _, test := range tests {
		crd.Annotations[IngressKey] = test.ingress

		IngressClass = test.controller
		DefaultClass = test.defClass

		b := CanDeleteResource(crd)
		if b != test.isValid {
			t.Errorf("test %v - expected %v but %v was returned", test, test.isValid, b)
		}
	}
}

func TestCandUpdateResource(t *testing.T) {
	dc := DefaultClass
	ic := IngressClass
	// restore original values after the tests
	defer func() {
		DefaultClass = dc
		IngressClass = ic
	}()

	tests := []struct {
		ingress    string
		controller string
		defClass   string
		isValid    bool
	}{
		{"", "", "nginx", true},
		{"", "nginx", "nginx", true},
		{"nginx", "nginx", "nginx", true},
		{"custom", "custom", "nginx", true},
		{"", "killer", "nginx", false},
		{"custom", "nginx", "nginx", false},
	}

	crd := &consumerv1.KongConsumer{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
	}

	data := map[string]string{}
	crd.SetAnnotations(data)
	for _, test := range tests {
		crd.Annotations[IngressKey] = test.ingress

		IngressClass = test.controller
		DefaultClass = test.defClass

		b, _, _ := CanUpdateResource(crd)
		if b != test.isValid {
			t.Errorf("test %v - expected %v but %v was returned", test, test.isValid, b)
		}
	}
}
