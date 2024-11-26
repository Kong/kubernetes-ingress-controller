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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
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
		{"custom", IgnoreClassMatch, "custom", true},
	}

	ing := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: corev1.NamespaceDefault,
		},
	}

	ingv1 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
			Annotations: map[string]string{
				IngressClassKey: DefaultIngressClass,
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: nil,
		},
	}

	data := map[string]string{}
	ing.SetAnnotations(data)
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		ing.Annotations[IngressClassKey] = test.ingress
		ingv1.Spec.IngressClassName = &test.ingress
		// TODO: unclear if we truly use IngressClassValidatorFunc anymore
		// IngressClassValidatorFuncFromObjectMeta appears to effectively supersede it, and is what we use in store
		// IngressClassValidatorFunc appears to be a test-only relic at this point
		fmeta := IngressClassValidatorFuncFromObjectMeta(test.controller)
		fv1 := IngressClassValidatorFuncFromV1Ingress(test.controller)

		resultMeta := fmeta(&ing.ObjectMeta, IngressClassKey, test.classMatching)
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
			want: nil,
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
			want: nil,
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

func TestExtractRequestBuffering(t *testing.T) {
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
					"konghq.com/request-buffering": "true",
				},
			},
			want: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ExtractRequestBuffering(tt.args.anns)
			if tt.want == "" {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
			if got != tt.want {
				t.Errorf("ExtractRequestBuffering() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractResponseBuffering(t *testing.T) {
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
					"konghq.com/response-buffering": "true",
				},
			},
			want: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ExtractResponseBuffering(tt.args.anns)
			if tt.want == "" {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
			if got != tt.want {
				t.Errorf("ExtractResponseBuffering() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractHostAliases(t *testing.T) {
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
			want: nil,
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/host-aliases": "foo.kong.com,bar.kong.com",
				},
			},
			want: []string{"foo.kong.com", "bar.kong.com"},
		},
		{
			name: "misconfigured",
			args: args{
				anns: map[string]string{
					"konghq.com/host-aliases": "",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			if got, _ = ExtractHostAliases(tt.args.anns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractHostAliases() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractConnectTimeout(t *testing.T) {
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
					"konghq.com/connect-timeout": "3000",
				},
			},
			want: "3000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ExtractConnectTimeout(tt.args.anns)
			if tt.want == "" {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExtractWriteTimeout(t *testing.T) {
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
					"konghq.com/write-timeout": "3000",
				},
			},
			want: "3000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ExtractWriteTimeout(tt.args.anns)
			if tt.want == "" {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExtractReadTimeout(t *testing.T) {
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
					"konghq.com/read-timeout": "3000",
				},
			},
			want: "3000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ExtractReadTimeout(tt.args.anns)
			if tt.want == "" {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExtractRetries(t *testing.T) {
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
					"konghq.com/retries": "3000",
				},
			},
			want: "3000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ExtractRetries(tt.args.anns)
			if tt.want == "" {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
			if got != tt.want {
				t.Errorf("ExtractRetries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractHeaders(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "empty",
			want: map[string][]string{},
		},
		{
			name: "empty with custom separator",
			args: args{
				anns: map[string]string{
					"konghq.com/headers-separator": ";",
				},
			},
			want: map[string][]string{},
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/headers.foo": "foo",
				},
			},
			want: map[string][]string{"foo": {"foo"}},
		},
		{
			name: "no separator",
			args: args{
				anns: map[string]string{
					"konghq.com/headersfoo": "foo",
				},
			},
			want: map[string][]string{},
		},
		{
			name: "separator with no header results in empty header value",
			args: args{
				anns: map[string]string{
					"konghq.com/headers.foo": "foo,",
				},
			},
			want: map[string][]string{"foo": {"foo", ""}},
		},
		{
			name: "no header name",
			args: args{
				anns: map[string]string{
					"konghq.com/headers.": "foo",
				},
			},
			want: map[string][]string{},
		},
		{
			name: "multiple header, multiple values, trailing spaces",
			args: args{
				anns: map[string]string{
					"konghq.com/headers.x-example":    "foo, bar, baz  ",
					"konghq.com/headers.x-additional": "foo",
				},
			},
			want: map[string][]string{
				"x-example":    {"foo", "bar", "baz"},
				"x-additional": {"foo"},
			},
		},
		{
			name: "multiple header, multiple values, custom separator",
			args: args{
				anns: map[string]string{
					"konghq.com/headers-separator":    ";",
					"konghq.com/headers.x-example":    "foo, bar;baz",
					"konghq.com/headers.x-additional": "foo",
				},
			},
			want: map[string][]string{
				"x-example":    {"foo, bar", "baz"},
				"x-additional": {"foo"},
			},
		},
		{
			name: "multiple header, multiple values, custom separator, leading & trailing spaces",
			args: args{
				anns: map[string]string{
					"konghq.com/headers-separator":    ";",
					"konghq.com/headers.x-example":    " foo, bar;cat,dog ;   baz ",
					"konghq.com/headers.x-additional": "foo;",
				},
			},
			want: map[string][]string{
				"x-example":    {"foo, bar", "cat,dog", "baz"},
				"x-additional": {"foo", ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ExtractHeaders(tt.args.anns)
			if len(tt.want) == 0 {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
			for key, val := range tt.want {
				actual, ok := got[key]
				assert.True(t, ok)
				assert.Equal(t, val, actual)
			}
		})
	}
}

func TestExtractPathHandling(t *testing.T) {
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
					"konghq.com/path-handling": "v1",
				},
			},
			want: "v1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ExtractPathHandling(tt.args.anns)
			if tt.want == "" {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
			if got != tt.want {
				t.Errorf("ExtractPathHandling() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractRewriteURI(t *testing.T) {
	type args struct {
		anns map[string]string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		exist bool
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "non-empty",
			args: args{
				anns: map[string]string{
					"konghq.com/rewrite": "/foo/$1",
				},
			},
			want:  "/foo/$1",
			exist: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, exist := ExtractRewriteURI(tt.args.anns)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.exist, exist)
		})
	}
}

func TestExtractTLSVerify(t *testing.T) {
	_, ok := ExtractTLSVerify(nil)
	assert.False(t, ok)

	_, ok = ExtractTLSVerify(map[string]string{})
	assert.False(t, ok)

	v, ok := ExtractTLSVerify(map[string]string{AnnotationPrefix + TLSVerifyKey: "true"})
	assert.True(t, ok)
	assert.Equal(t, true, v)

	v, ok = ExtractTLSVerify(map[string]string{AnnotationPrefix + TLSVerifyKey: "false"})
	assert.True(t, ok)
	assert.Equal(t, false, v)
}

func TestExtractTLSVerifyDepth(t *testing.T) {
	_, ok := ExtractTLSVerifyDepth(nil)
	assert.False(t, ok)

	_, ok = ExtractTLSVerifyDepth(map[string]string{})
	assert.False(t, ok)

	_, ok = ExtractTLSVerifyDepth(map[string]string{AnnotationPrefix + TLSVerifyDepthKey: "non-integer"})
	assert.False(t, ok)

	v, ok := ExtractTLSVerifyDepth(map[string]string{AnnotationPrefix + TLSVerifyDepthKey: "1"})
	assert.True(t, ok)
	assert.Equal(t, 1, v)
}

func TestExtractCACertificates(t *testing.T) {
	v := ExtractCACertificates(nil)
	assert.Empty(t, v)

	v = ExtractCACertificates(map[string]string{})
	assert.Empty(t, v)

	v = ExtractCACertificates(map[string]string{AnnotationPrefix + CACertificatesKey: "foo,bar"})
	assert.Len(t, v, 2)
	assert.Equal(t, "foo", v[0])
	assert.Equal(t, "bar", v[1])
}
