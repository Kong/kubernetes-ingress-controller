package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func TestKongPluginFromK8SClusterPlugin(t *testing.T) {
	assert := assert.New(t)
	store, _ := store.NewFakeStore(store.FakeObjects{
		Secrets: []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "conf-secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"correlation-id-config":            []byte(`{"header_name": "foo"}`),
					"correlation-id-headername":        []byte(`"foo"`),
					"correlation-id-generator":         []byte(`"uuid"`),
					"correlation-id-invalid":           []byte(`"aaa`),
					"response-transformer-add-headers": []byte(`["h1:v1","h2:v2"]`),
				},
			},
		},
	})
	type args struct {
		plugin kongv1.KongClusterPlugin
	}
	tests := []struct {
		name    string
		args    args
		want    kong.Plugin
		wantErr bool
	}{
		{
			name: "basic configuration",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					Protocols:    []kongv1.KongProtocol{"http"},
					PluginName:   "correlation-id",
					InstanceName: "example",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("correlation-id"),
				Config: kong.Configuration{
					"header_name": "foo",
				},
				Protocols:    kong.StringSlice("http"),
				InstanceName: kong.String("example"),
			},
			wantErr: false,
		},
		{
			name: "secret configuration",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "correlation-id-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("correlation-id"),
				Config: kong.Configuration{
					"header_name": "foo",
				},
				Protocols: kong.StringSlice("http"),
			},
			wantErr: false,
		},
		{
			name: "missing secret configuration",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "correlation-id-config",
							Secret:    "missing",
							Namespace: "default",
						},
					},
				},
			},
			want:    kong.Plugin{},
			wantErr: true,
		},
		{
			name: "non-JSON configuration",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{{}`),
					},
				},
			},
			want:    kong.Plugin{},
			wantErr: true,
		},
		{
			name: "both Config and ConfigFrom set",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "correlation-id-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
			},
			want:    kong.Plugin{},
			wantErr: true,
		},
		{
			name: "Config and ConfigPatches set",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/generator",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Key:       "correlation-id-generator",
									Secret:    "conf-secret",
									Namespace: "default",
								},
							},
						},
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("correlation-id"),
				Config: kong.Configuration{
					"header_name": "foo",
					"generator":   "uuid",
				},
				Protocols: kong.StringSlice("http"),
			},
		},
		{
			name: "configPatch on subpath of non-exist path",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "response-transformer",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"replace":{"headers":["foo:bar"]}}`),
					},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/add/headers",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: "default",
									Key:       "response-transformer-add-headers",
									Secret:    "conf-secret",
								},
							},
						},
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("response-transformer"),
				Config: kong.Configuration{
					"replace": map[string]interface{}{
						"headers": []interface{}{
							"foo:bar",
						},
					},
					"add": map[string]interface{}{
						"headers": []interface{}{
							"h1:v1",
							"h2:v2",
						},
					},
				},
				Protocols: kong.StringSlice("http"),
			},
		},
		{
			name: "empty config and configPatch for particular paths",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config:     apiextensionsv1.JSON{},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/header_name",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: "default",
									Key:       "correlation-id-headername",
									Secret:    "conf-secret",
								},
							},
						},
						{
							Path: "/generator",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: "default",
									Key:       "correlation-id-generator",
									Secret:    "conf-secret",
								},
							},
						},
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("correlation-id"),
				Config: kong.Configuration{
					"header_name": "foo",
					"generator":   "uuid",
				},
				Protocols: kong.StringSlice("http"),
			},
		},
		{
			name: "empty config and configPatch for whole object",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config:     apiextensionsv1.JSON{},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: "default",
									Key:       "correlation-id-config",
									Secret:    "conf-secret",
								},
							},
						},
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("correlation-id"),
				Config: kong.Configuration{
					"header_name": "foo",
				},
				Protocols: kong.StringSlice("http"),
			},
		},
		{
			name: "missing secret in configPatches",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/generator",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: "default",
									Key:       "correlation-id-generator",
									Secret:    "missing-secret",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing key of secret in cofigPatches",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/generator",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: "default",
									Key:       "correlation-id-missing",
									Secret:    "conf-secret",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid value in configPatches",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/generator",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: "default",
									Key:       "correlation-id-invalid",
									Secret:    "conf-secret",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kongPluginFromK8SClusterPlugin(store, tt.args.plugin)
			if (err != nil) != tt.wantErr {
				t.Errorf("kongPluginFromK8SClusterPlugin error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(tt.want, got.Plugin)
			assert.NotEmpty(t, got.K8sParent)
		})
	}
}

func TestKongPluginFromK8SPlugin(t *testing.T) {
	assert := assert.New(t)
	store, _ := store.NewFakeStore(store.FakeObjects{
		Secrets: []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "conf-secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"correlation-id-config":            []byte(`{"header_name": "foo"}`),
					"correlation-id-headername":        []byte(`"foo"`),
					"correlation-id-generator":         []byte(`"uuid"`),
					"correlation-id-invalid":           []byte(`"aaa`),
					"response-transformer-add-headers": []byte(`["h1:v1","h2:v2"]`),
				},
			},
		},
	})
	type args struct {
		plugin kongv1.KongPlugin
	}
	tests := []struct {
		name    string
		args    args
		want    kong.Plugin
		wantErr bool
	}{
		{
			name: "basic configuration",
			args: args{
				plugin: kongv1.KongPlugin{
					Protocols:    []kongv1.KongProtocol{"http"},
					PluginName:   "correlation-id",
					InstanceName: "example",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("correlation-id"),
				Config: kong.Configuration{
					"header_name": "foo",
				},
				Protocols:    kong.StringSlice("http"),
				InstanceName: kong.String("example"),
			},
			wantErr: false,
		},
		{
			name: "secret configuration",
			args: args{
				plugin: kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "correlation-id-config",
							Secret: "conf-secret",
						},
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("correlation-id"),
				Config: kong.Configuration{
					"header_name": "foo",
				},
				Protocols: kong.StringSlice("http"),
			},
			wantErr: false,
		},
		{
			name: "missing secret configuration",
			args: args{
				plugin: kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "correlation-id-config",
							Secret: "missing",
						},
					},
				},
			},
			want:    kong.Plugin{},
			wantErr: true,
		},
		{
			name: "non-JSON configuration",
			args: args{
				plugin: kongv1.KongPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{{}`),
					},
				},
			},
			want:    kong.Plugin{},
			wantErr: true,
		},
		{
			name: "both Config and ConfigFrom set",
			args: args{
				plugin: kongv1.KongPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "correlation-id-config",
							Secret: "conf-secret",
						},
					},
				},
			},
			want:    kong.Plugin{},
			wantErr: true,
		},
		{
			name: "config and configPatches set",
			args: args{
				plugin: kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/generator",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Key:    "correlation-id-generator",
									Secret: "conf-secret",
								},
							},
						},
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("correlation-id"),
				Config: kong.Configuration{
					"header_name": "foo",
					"generator":   "uuid",
				},
				Protocols: kong.StringSlice("http"),
			},
		},
		{
			name: "configPatch on subpath of non-exist path",
			args: args{
				plugin: kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "response-transformer",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"replace":{"headers":["foo:bar"]}}`),
					},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/add/headers",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Key:    "response-transformer-add-headers",
									Secret: "conf-secret",
								},
							},
						},
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("response-transformer"),
				Config: kong.Configuration{
					"replace": map[string]interface{}{
						"headers": []interface{}{
							"foo:bar",
						},
					},
					"add": map[string]interface{}{
						"headers": []interface{}{
							"h1:v1",
							"h2:v2",
						},
					},
				},
				Protocols: kong.StringSlice("http"),
			},
		},
		{
			name: "empty config and configPatch for particular paths",
			args: args{
				plugin: kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config:     apiextensionsv1.JSON{},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/header_name",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Key:    "correlation-id-headername",
									Secret: "conf-secret",
								},
							},
						},
						{
							Path: "/generator",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Key:    "correlation-id-generator",
									Secret: "conf-secret",
								},
							},
						},
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("correlation-id"),
				Config: kong.Configuration{
					"header_name": "foo",
					"generator":   "uuid",
				},
				Protocols: kong.StringSlice("http"),
			},
		},
		{
			name: "empty config and configPatch for whole object",
			args: args{
				plugin: kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config:     apiextensionsv1.JSON{},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Key:    "correlation-id-config",
									Secret: "conf-secret",
								},
							},
						},
					},
				},
			},
			want: kong.Plugin{
				Name: kong.String("correlation-id"),
				Config: kong.Configuration{
					"header_name": "foo",
				},
				Protocols: kong.StringSlice("http"),
			},
		},
		{
			name: "missing secret in configPatches",
			args: args{
				plugin: kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/generator",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Key:    "correlation-id-generator",
									Secret: "missing-secret",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing key of secret in configPatches",
			args: args{
				plugin: kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/generator",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Key:    "correlation-id-missing",
									Secret: "conf-secret",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid value in configPatches",
			args: args{
				plugin: kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/generator",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Key:    "correlation-id-invalid",
									Secret: "conf-secret",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kongPluginFromK8SPlugin(store, tt.args.plugin)
			if (err != nil) != tt.wantErr {
				t.Errorf("kongPluginFromK8SPlugin error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// don't care about tags in this test
			got.Tags = nil
			assert.Equal(tt.want, got.Plugin)
			assert.NotEmpty(t, got.K8sParent)
		})
	}
}
