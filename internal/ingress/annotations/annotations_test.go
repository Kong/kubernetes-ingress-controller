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

package annotations

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExtractKongPluginsFromAnnotations(t *testing.T) {
	data := map[string]string{
		"plugins.konghq.com": "kp-rl, kp-cors",
	}

	ka := ExtractKongPluginsFromAnnotations(data)
	if len(ka) != 2 {
		t.Errorf("expected two keys but %v returned", len(ka))
	}
	if ka[0] != "kp-rl" {
		t.Errorf("expected first element to be 'kp-rl'")
	}
}

func TestExtractConfigurationName(t *testing.T) {
	data := map[string]string{
		"configuration.konghq.com": "demo",
	}

	cn := ExtractConfigurationName(data)
	if cn != "demo" {
		t.Errorf("expected demo as configuration name but got %v", cn)
	}
}

func TestExtractProtocolName(t *testing.T) {
	data := map[string]string{
		"configuration.konghq.com/protocol": "grpc",
	}

	pn := ExtractProtocolName(data)
	if pn != "grpc" {
		t.Errorf("expected grpc as configuration name but got %v", pn)
	}
}

func TestExtractProtocolNames(t *testing.T) {
	data := map[string]string{
		"configuration.konghq.com/protocols": "grpc,grpcs",
	}

	s := []string{"grpc", "grpcs"}

	pns := ExtractProtocolNames(data)
	if !reflect.DeepEqual(pns, s) {
		t.Errorf("expected grpc,grpcs as configuration name but got %v", pns)
	}
}

func TestExtractClientCert(t *testing.T) {
	data := map[string]string{
		"configuration.konghq.com/client-cert": "secret1",
	}

	secret := ExtractClientCertificate(data)
	if secret != "secret1" {
		t.Errorf("expected secret as secret1 but got %v", secret)
	}
}

func TestIngrssClassValidatorFunc(t *testing.T) {
	tests := []struct {
		ingress    string
		controller string
		isValid    bool
	}{
		{"", "", true},
		{"", "kong", true},
		{"kong", "kong", true},
		{"custom", "custom", true},
		{"", "killer", false},
		{"custom", "kong", false},
	}

	ing := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: v1.NamespaceDefault,
		},
	}

	data := map[string]string{}
	ing.SetAnnotations(data)
	for _, test := range tests {
		ing.Annotations[ingressClassKey] = test.ingress
		f := IngressClassValidatorFunc(test.controller)
		b := f(&ing.ObjectMeta)
		if b != test.isValid {
			t.Errorf("test %v - expected %v but %v was returned", test, test.isValid, b)
		}
	}
}
