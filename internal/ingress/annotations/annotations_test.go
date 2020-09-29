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
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIngressClassValidatorFunc(t *testing.T) {
	tests := []struct {
		ingress       string        // the class set on the Ingress resource
		classMatching ClassMatching // the "user" classless ingress flag value, translated to its match strategy
		controller    string        // the class set on the controller
		isValid       bool          // the expected verdict
	}{
		{"", ExactOrEmptyClassMatch, "", true},
		{"", ExactOrEmptyClassMatch, DefaultIngressClass, true},
		{"", ExactClassMatch, DefaultIngressClass, false},
		{DefaultIngressClass, ExactOrEmptyClassMatch, DefaultIngressClass, true},
		{DefaultIngressClass, ExactClassMatch, DefaultIngressClass, true},
		{"custom", ExactOrEmptyClassMatch, "custom", true},
		{"", ExactOrEmptyClassMatch, "killer", true},
		{"custom", ExactOrEmptyClassMatch, DefaultIngressClass, false},
		{"custom", ExactClassMatch, DefaultIngressClass, false},
		{"", ExactOrEmptyClassMatch, "custom", true},
		{"", ExactClassMatch, "kozel", false},
		{"kozel", ExactOrEmptyClassMatch, "kozel", true},
		{"kozel", ExactClassMatch, "kozel", true},
		{"", ExactOrEmptyClassMatch, "killer", true},
		{"custom", ExactOrEmptyClassMatch, "kozel", false},
		{"custom", ExactClassMatch, "kozel", false},
	}

	ing := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: v1.NamespaceDefault,
		},
	}

	ingv1 := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
			Annotations: map[string]string{
				IngressClassKey: DefaultIngressClass,
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: nil,
		},
	}

	data := map[string]string{}
	ing.SetAnnotations(data)
	for _, test := range tests {
		ing.Annotations[IngressClassKey] = test.ingress
		ingv1.Spec.IngressClassName = &test.ingress
		// TODO: unclear if we truly use IngressClassValidatorFunc anymore
		// IngressClassValidatorFuncFromObjectMeta appears to effectively supersede it, and is what we use in store
		// IngressClassValidatorFunc appears to be a test-only relic at this point
		f := IngressClassValidatorFunc(test.controller)
		fmeta := IngressClassValidatorFuncFromObjectMeta(test.controller)
		fv1 := IngressClassValidatorFuncFromV1Ingress(test.controller)

		result := f(&ing.ObjectMeta, test.classMatching)
		if result != test.isValid {
			t.Errorf("test %v - expected %v but %v was returned", test, test.isValid, result)
		}
		resultMeta := fmeta(&ing.ObjectMeta, test.classMatching)
		if resultMeta != test.isValid {
			t.Errorf("meta test %v - expected %v but %v was returned", test, test.isValid, resultMeta)
		}
		resultV1 := fv1(ingv1, test.classMatching)
		if resultV1 != test.isValid {
			t.Errorf("v1 test %v - expected %v but %v was returned", test, test.isValid, resultV1)
		}
	}
}

func TestExtractPath(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/path": "/foo",
				},
			},
			want: "/foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractPath(tt.args.anns); got != tt.want {
				t.Errorf("ExtractPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractKongPluginsFromAnnotations(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/plugins": "kp-rl, kp-cors",
				},
			},
			want: []string{"kp-rl", "kp-cors"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractKongPluginsFromAnnotations(tt.args.anns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractKongPluginsFromAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractConfigurationName(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/override": "foo",
				},
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractConfigurationName(tt.args.anns); got != tt.want {
				t.Errorf("ExtractConfigurationName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractProtocolName(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/protocol": "foo",
				},
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractProtocolName(tt.args.anns); got != tt.want {
				t.Errorf("ExtractProtocolName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractProtocolNames(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/protocols": "foo,bar",
				},
			},
			want: []string{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractProtocolNames(tt.args.anns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractProtocolNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractClientCertificate(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/client-cert": "foo",
				},
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractClientCertificate(tt.args.anns); got != tt.want {
				t.Errorf("ExtractClientCertificate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasServiceUpstreamAnnotation(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "basic sanity",
			args: args{
				anns: map[string]string{
					"ingress.kubernetes.io/service-upstream": "true",
				},
			},
			want: true,
		},
		{
			name: "garbage value",
			args: args{
				anns: map[string]string{
					"ingress.kubernetes.io/service-upstream": "42",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasServiceUpstreamAnnotation(tt.args.anns); got != tt.want {
				t.Errorf("HasServiceUpstreamAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractStripPath(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/strip-path": "true",
				},
			},
			want: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractStripPath(tt.args.anns); got != tt.want {
				t.Errorf("ExtractStripPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractHTTPSRedirectStatusCode(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/https-redirect-status-code": "302",
				},
			},
			want: "302",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractHTTPSRedirectStatusCode(tt.args.anns); got != tt.want {
				t.Errorf("ExtractHTTPSRedirectStatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasForceSSLRedirectAnnotation(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "basic sanity",
			args: args{
				anns: map[string]string{
					"ingress.kubernetes.io/force-ssl-redirect": "true",
				},
			},
			want: true,
		},
		{
			name: "garbage value",
			args: args{
				anns: map[string]string{
					"ingress.kubernetes.io/force-ssl-redirect": "xyz",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasForceSSLRedirectAnnotation(tt.args.anns); got != tt.want {
				t.Errorf("HasForceSSLRedirectAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractPreserveHost(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/preserve-host": "true",
				},
			},
			want: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractPreserveHost(tt.args.anns); got != tt.want {
				t.Errorf("ExtractPreserveHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractRegexPriority(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/regex-priority": "10",
				},
			},
			want: "10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractRegexPriority(tt.args.anns); got != tt.want {
				t.Errorf("ExtractRegexPriority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractHostHeader(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/host-header": "example.net",
				},
			},
			want: "example.net",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractHostHeader(tt.args.anns); got != tt.want {
				t.Errorf("ExtractHostHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractMethods(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			want: []string{},
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/methods": "POST,GET",
				},
			},
			want: []string{"POST", "GET"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractMethods(tt.args.anns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractMethods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractSNIs(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			want: []string{},
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/snis": "hrodna.kong.example,katowice.kong.example",
				},
			},
			want: []string{"hrodna.kong.example", "katowice.kong.example"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			if got, _ = ExtractSNIs(tt.args.anns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractSNIs() = %v, want %v", got, tt.want)
			}
		})
	}
}
