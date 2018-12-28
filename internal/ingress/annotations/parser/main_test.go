/*
Copyright 2016 The Kubernetes Authors.

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

package parser

import (
	"testing"

	consumerv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/consumer/v1"
	credentialv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/credential/v1"
	pluginv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/plugin/v1"
	api "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildIngress() *extensions.Ingress {
	return &extensions.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
		Spec: extensions.IngressSpec{},
	}
}

func buildPlugin() *pluginv1.KongPlugin {
	return &pluginv1.KongPlugin{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
	}
}

func buildCredential() *credentialv1.KongCredential {
	return &credentialv1.KongCredential{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
	}
}

func buildConsumer() *consumerv1.KongConsumer {
	return &consumerv1.KongConsumer{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
	}
}
func TestGetBoolAnnotation(t *testing.T) {
	ing := buildIngress()

	_, err := GetBoolAnnotation("", nil)
	if err == nil {
		t.Errorf("expected error but retuned nil")
	}

	tests := []struct {
		name   string
		field  string
		value  string
		exp    bool
		expErr bool
	}{
		{"valid - false", "bool", "false", false, false},
		{"valid - true", "bool", "true", true, false},
	}

	data := map[string]string{}
	ing.SetAnnotations(data)

	for _, test := range tests {
		data[GetAnnotationWithPrefix(test.field)] = test.value

		u, err := GetBoolAnnotation(test.field, ing)
		if test.expErr {
			if err == nil {
				t.Errorf("%v: expected error but retuned nil", test.name)
			}
			continue
		}
		if u != test.exp {
			t.Errorf("%v: expected \"%v\" but \"%v\" was returned", test.name, test.exp, u)
		}

		delete(data, test.field)
	}
}

func TestGetStringAnnotation(t *testing.T) {
	ing := buildIngress()

	_, err := GetStringAnnotation("", nil)
	if err == nil {
		t.Errorf("expected error but retuned nil")
	}

	tests := []struct {
		name   string
		field  string
		value  string
		exp    string
		expErr bool
	}{
		{"valid - A", "string", "A", "A", false},
		{"valid - B", "string", "B", "B", false},
	}

	data := map[string]string{}
	ing.SetAnnotations(data)

	for _, test := range tests {
		data[GetAnnotationWithPrefix(test.field)] = test.value

		s, err := GetStringAnnotation(test.field, ing)
		if test.expErr {
			if err == nil {
				t.Errorf("%v: expected error but retuned nil", test.name)
			}
			continue
		}
		if s != test.exp {
			t.Errorf("%v: expected \"%v\" but \"%v\" was returned", test.name, test.exp, s)
		}

		delete(data, test.field)
	}
}

func TestGetStringAnnotationPlugin(t *testing.T) {
	res := buildPlugin()

	_, err := GetStringAnnotationPlugin("", nil)
	if err == nil {
		t.Errorf("expected error but retuned nil")
	}

	tests := []struct {
		name   string
		field  string
		value  string
		exp    string
		expErr bool
	}{
		{"valid - A", "string", "A", "A", false},
		{"valid - B", "string", "B", "B", false},
	}

	data := map[string]string{}
	res.SetAnnotations(data)

	for _, test := range tests {
		data[GetAnnotationWithPrefix(test.field)] = test.value

		s, err := GetStringAnnotationPlugin(test.field, res)
		if test.expErr {
			if err == nil {
				t.Errorf("%v: expected error but retuned nil", test.name)
			}
			continue
		}
		if s != test.exp {
			t.Errorf("%v: expected \"%v\" but \"%v\" was returned", test.name, test.exp, s)
		}

		delete(data, test.field)
	}
}

func TestGetStringAnnotationCredential(t *testing.T) {
	res := buildCredential()

	_, err := GetStringAnnotationCredential("", nil)
	if err == nil {
		t.Errorf("expected error but retuned nil")
	}

	tests := []struct {
		name   string
		field  string
		value  string
		exp    string
		expErr bool
	}{
		{"valid - A", "string", "A", "A", false},
		{"valid - B", "string", "B", "B", false},
	}

	data := map[string]string{}
	res.SetAnnotations(data)

	for _, test := range tests {
		data[GetAnnotationWithPrefix(test.field)] = test.value

		s, err := GetStringAnnotationCredential(test.field, res)
		if test.expErr {
			if err == nil {
				t.Errorf("%v: expected error but retuned nil", test.name)
			}
			continue
		}
		if s != test.exp {
			t.Errorf("%v: expected \"%v\" but \"%v\" was returned", test.name, test.exp, s)
		}

		delete(data, test.field)
	}
}

func TestGetStringAnnotationConsumer(t *testing.T) {
	res := buildConsumer()

	_, err := GetStringAnnotationConsumer("", nil)
	if err == nil {
		t.Errorf("expected error but retuned nil")
	}

	tests := []struct {
		name   string
		field  string
		value  string
		exp    string
		expErr bool
	}{
		{"valid - A", "string", "A", "A", false},
		{"valid - B", "string", "B", "B", false},
	}

	data := map[string]string{}
	res.SetAnnotations(data)

	for _, test := range tests {
		data[GetAnnotationWithPrefix(test.field)] = test.value

		s, err := GetStringAnnotationConsumer(test.field, res)
		if test.expErr {
			if err == nil {
				t.Errorf("%v: expected error but retuned nil", test.name)
			}
			continue
		}
		if s != test.exp {
			t.Errorf("%v: expected \"%v\" but \"%v\" was returned", test.name, test.exp, s)
		}

		delete(data, test.field)
	}
}

func TestGetIntAnnotation(t *testing.T) {
	ing := buildIngress()

	_, err := GetIntAnnotation("", nil)
	if err == nil {
		t.Errorf("expected error but retuned nil")
	}

	tests := []struct {
		name   string
		field  string
		value  string
		exp    int
		expErr bool
	}{
		{"valid - A", "string", "1", 1, false},
		{"valid - B", "string", "2", 2, false},
	}

	data := map[string]string{}
	ing.SetAnnotations(data)

	for _, test := range tests {
		data[GetAnnotationWithPrefix(test.field)] = test.value

		s, err := GetIntAnnotation(test.field, ing)
		if test.expErr {
			if err == nil {
				t.Errorf("%v: expected error but retuned nil", test.name)
			}
			continue
		}
		if s != test.exp {
			t.Errorf("%v: expected \"%v\" but \"%v\" was returned", test.name, test.exp, s)
		}

		delete(data, test.field)
	}
}
