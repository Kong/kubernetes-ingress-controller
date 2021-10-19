package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	networking "k8s.io/api/networking/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/kongstate"
)

func TestMergeIngressRules(t *testing.T) {
	for _, tt := range []struct {
		name       string
		inputs     []ingressRules
		wantOutput *ingressRules
	}{
		{
			name: "empty list",
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{},
				ServiceNameToServices: map[string]kongstate.Service{},
			},
		},
		{
			name: "nil maps",
			inputs: []ingressRules{
				{}, {}, {},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{},
				ServiceNameToServices: map[string]kongstate.Service{},
			},
		},
		{
			name: "one input",
			inputs: []ingressRules{
				{
					SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}},
					ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}},
				ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}},
			},
		},
		{
			name: "three inputs",
			inputs: []ingressRules{
				{
					SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}},
					ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}},
				},
				{
					SecretNameToSNIs: map[string][]string{"g": {"h"}},
				},
				{
					ServiceNameToServices: map[string]kongstate.Service{"2": {Namespace: "carrot"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}, "g": {"h"}},
				ServiceNameToServices: map[string]kongstate.Service{"1": {Namespace: "potato"}, "2": {Namespace: "carrot"}},
			},
		},
		{
			name: "can merge SNI arrays",
			inputs: []ingressRules{
				{
					SecretNameToSNIs: map[string][]string{"a": {"b", "c"}},
				},
				{
					SecretNameToSNIs: map[string][]string{"a": {"d", "e"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{"a": {"b", "c", "d", "e"}},
				ServiceNameToServices: map[string]kongstate.Service{},
			},
		},
		{
			name: "overwrites services",
			inputs: []ingressRules{
				{
					ServiceNameToServices: map[string]kongstate.Service{"svc-name": {Namespace: "old"}},
				},
				{
					ServiceNameToServices: map[string]kongstate.Service{"svc-name": {Namespace: "new"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{},
				ServiceNameToServices: map[string]kongstate.Service{"svc-name": {Namespace: "new"}},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput := mergeIngressRules(tt.inputs...)
			assert.Equal(t, &gotOutput, tt.wantOutput)
		})
	}
}

func Test_addFromIngressV1beta1TLS(t *testing.T) {
	type args struct {
		tlsSections []networking.IngressTLS
		namespace   string
	}
	tests := []struct {
		name string
		args args
		want SecretNameToSNIs
	}{
		{
			args: args{
				tlsSections: []networking.IngressTLS{
					{
						Hosts: []string{
							"1.example.com",
							"2.example.com",
						},
						SecretName: "sooper-secret",
					},
					{
						Hosts: []string{
							"3.example.com",
							"4.example.com",
						},
						SecretName: "sooper-secret2",
					},
				},
				namespace: "foo",
			},
			want: SecretNameToSNIs{
				"foo/sooper-secret":  {"1.example.com", "2.example.com"},
				"foo/sooper-secret2": {"3.example.com", "4.example.com"},
			},
		},
		{
			args: args{
				tlsSections: []networking.IngressTLS{
					{
						Hosts: []string{
							"1.example.com",
						},
						SecretName: "sooper-secret",
					},
					{
						Hosts: []string{
							"3.example.com",
							"1.example.com",
							"4.example.com",
						},
						SecretName: "sooper-secret2",
					},
				},
				namespace: "foo",
			},
			want: SecretNameToSNIs{
				"foo/sooper-secret":  {"1.example.com"},
				"foo/sooper-secret2": {"3.example.com", "4.example.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newSecretNameToSNIs()
			m.addFromIngressV1beta1TLS(tt.args.tlsSections, tt.args.namespace)
			assert.Equal(t, m, tt.want)
		})
	}
}
