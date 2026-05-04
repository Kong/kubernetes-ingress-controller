package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configurationv1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1"

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
		plugin configurationv1.KongClusterPlugin
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
				plugin: configurationv1.KongClusterPlugin{
					Protocols:    []configurationv1.KongProtocol{"http"},
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
				plugin: configurationv1.KongClusterPlugin{
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					ConfigFrom: &configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
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
				plugin: configurationv1.KongClusterPlugin{
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					ConfigFrom: &configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
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
				plugin: configurationv1.KongClusterPlugin{
					Protocols:  []configurationv1.KongProtocol{"http"},
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
				plugin: configurationv1.KongClusterPlugin{
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigFrom: &configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
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
				plugin: configurationv1.KongClusterPlugin{
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []configurationv1.NamespacedConfigPatch{
						{
							Path: "/generator",
							ValueFrom: configurationv1.NamespacedConfigSource{
								SecretValue: configurationv1.NamespacedSecretValueFromSource{
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
				plugin: configurationv1.KongClusterPlugin{
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "response-transformer",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"replace":{"headers":["foo:bar"]}}`),
					},
					ConfigPatches: []configurationv1.NamespacedConfigPatch{
						{
							Path: "/add/headers",
							ValueFrom: configurationv1.NamespacedConfigSource{
								SecretValue: configurationv1.NamespacedSecretValueFromSource{
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
					"replace": map[string]any{
						"headers": []any{
							"foo:bar",
						},
					},
					"add": map[string]any{
						"headers": []any{
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
				plugin: configurationv1.KongClusterPlugin{
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config:     apiextensionsv1.JSON{},
					ConfigPatches: []configurationv1.NamespacedConfigPatch{
						{
							Path: "/header_name",
							ValueFrom: configurationv1.NamespacedConfigSource{
								SecretValue: configurationv1.NamespacedSecretValueFromSource{
									Namespace: "default",
									Key:       "correlation-id-headername",
									Secret:    "conf-secret",
								},
							},
						},
						{
							Path: "/generator",
							ValueFrom: configurationv1.NamespacedConfigSource{
								SecretValue: configurationv1.NamespacedSecretValueFromSource{
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
				plugin: configurationv1.KongClusterPlugin{
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config:     apiextensionsv1.JSON{},
					ConfigPatches: []configurationv1.NamespacedConfigPatch{
						{
							Path: "",
							ValueFrom: configurationv1.NamespacedConfigSource{
								SecretValue: configurationv1.NamespacedSecretValueFromSource{
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
				plugin: configurationv1.KongClusterPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []configurationv1.NamespacedConfigPatch{
						{
							Path: "/generator",
							ValueFrom: configurationv1.NamespacedConfigSource{
								SecretValue: configurationv1.NamespacedSecretValueFromSource{
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
				plugin: configurationv1.KongClusterPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []configurationv1.NamespacedConfigPatch{
						{
							Path: "/generator",
							ValueFrom: configurationv1.NamespacedConfigSource{
								SecretValue: configurationv1.NamespacedSecretValueFromSource{
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
				plugin: configurationv1.KongClusterPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []configurationv1.NamespacedConfigPatch{
						{
							Path: "/generator",
							ValueFrom: configurationv1.NamespacedConfigSource{
								SecretValue: configurationv1.NamespacedSecretValueFromSource{
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
		plugin configurationv1.KongPlugin
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
				plugin: configurationv1.KongPlugin{
					Protocols:    []configurationv1.KongProtocol{"http"},
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
				plugin: configurationv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					ConfigFrom: &configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
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
				plugin: configurationv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					ConfigFrom: &configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
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
				plugin: configurationv1.KongPlugin{
					Protocols:  []configurationv1.KongProtocol{"http"},
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
				plugin: configurationv1.KongPlugin{
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigFrom: &configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
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
				plugin: configurationv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []configurationv1.ConfigPatch{
						{
							Path: "/generator",
							ValueFrom: configurationv1.ConfigSource{
								SecretValue: configurationv1.SecretValueFromSource{
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
				plugin: configurationv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "response-transformer",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"replace":{"headers":["foo:bar"]}}`),
					},
					ConfigPatches: []configurationv1.ConfigPatch{
						{
							Path: "/add/headers",
							ValueFrom: configurationv1.ConfigSource{
								SecretValue: configurationv1.SecretValueFromSource{
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
					"replace": map[string]any{
						"headers": []any{
							"foo:bar",
						},
					},
					"add": map[string]any{
						"headers": []any{
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
				plugin: configurationv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config:     apiextensionsv1.JSON{},
					ConfigPatches: []configurationv1.ConfigPatch{
						{
							Path: "/header_name",
							ValueFrom: configurationv1.ConfigSource{
								SecretValue: configurationv1.SecretValueFromSource{
									Key:    "correlation-id-headername",
									Secret: "conf-secret",
								},
							},
						},
						{
							Path: "/generator",
							ValueFrom: configurationv1.ConfigSource{
								SecretValue: configurationv1.SecretValueFromSource{
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
				plugin: configurationv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config:     apiextensionsv1.JSON{},
					ConfigPatches: []configurationv1.ConfigPatch{
						{
							Path: "",
							ValueFrom: configurationv1.ConfigSource{
								SecretValue: configurationv1.SecretValueFromSource{
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
				plugin: configurationv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []configurationv1.ConfigPatch{
						{
							Path: "/generator",
							ValueFrom: configurationv1.ConfigSource{
								SecretValue: configurationv1.SecretValueFromSource{
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
				plugin: configurationv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []configurationv1.ConfigPatch{
						{
							Path: "/generator",
							ValueFrom: configurationv1.ConfigSource{
								SecretValue: configurationv1.SecretValueFromSource{
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
				plugin: configurationv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"header_name": "foo"}`),
					},
					ConfigPatches: []configurationv1.ConfigPatch{
						{
							Path: "/generator",
							ValueFrom: configurationv1.ConfigSource{
								SecretValue: configurationv1.SecretValueFromSource{
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

func TestPluginSanitizedCopy(t *testing.T) {
	parent := &configurationv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "plugin-parent",
			Namespace: "default",
		},
	}

	tests := []struct {
		name      string
		in        Plugin
		want      Plugin
		postCheck func(t *testing.T, in Plugin, got Plugin)
	}{
		{
			name: "redacts all top level config values and preserves other fields",
			in: Plugin{
				Plugin: kong.Plugin{
					ID:           kong.String("plugin-id"),
					Name:         kong.String("rate-limiting"),
					InstanceName: kong.String("instance-name"),
					RunOn:        kong.String("first"),
					Enabled:      kong.Bool(false),
					Protocols:    kong.StringSlice("http", "https"),
					Tags:         []*string{kong.String("tag-1"), kong.String("tag-2")},
					Ordering: &kong.PluginOrdering{
						Before: map[string][]string{
							"access": {"key-auth"},
						},
					},
					Config: kong.Configuration{
						"string":  "secret",
						"number":  123,
						"boolean": true,
						"object": map[string]any{
							"nested": "secret",
						},
						"array": []any{"first", "second"},
						"null":  nil,
					},
				},
				K8sParent: parent,
			},
			want: Plugin{
				Plugin: kong.Plugin{
					ID:           kong.String("plugin-id"),
					Name:         kong.String("rate-limiting"),
					InstanceName: kong.String("instance-name"),
					RunOn:        kong.String("first"),
					Enabled:      kong.Bool(false),
					Protocols:    kong.StringSlice("http", "https"),
					Tags:         []*string{kong.String("tag-1"), kong.String("tag-2")},
					Ordering: &kong.PluginOrdering{
						Before: map[string][]string{
							"access": {"key-auth"},
						},
					},
					Config: kong.Configuration{
						"string":  "{REDACTED}",
						"number":  "{REDACTED}",
						"boolean": "{REDACTED}",
						"object":  "{REDACTED}",
						"array":   "{REDACTED}",
						"null":    "{REDACTED}",
					},
				},
				K8sParent: parent,
			},
			postCheck: func(t *testing.T, in Plugin, got Plugin) {
				assert.Equal(t, "secret", in.Config["string"])
				assert.Equal(t, 123, in.Config["number"])
				assert.Equal(t, true, in.Config["boolean"])
				assert.Equal(t, map[string]any{"nested": "secret"}, in.Config["object"])
				assert.Equal(t, []any{"first", "second"}, in.Config["array"])
				assert.Nil(t, in.Config["null"])
				assert.Same(t, in.K8sParent, got.K8sParent)
			},
		},
		{
			name: "keeps nil config nil",
			in: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("key-auth"),
				},
			},
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("key-auth"),
				},
			},
			postCheck: func(t *testing.T, in Plugin, got Plugin) {
				assert.Nil(t, in.Config)
				assert.Nil(t, got.Config)
			},
		},
		{
			name: "returns an isolated config map copy",
			in: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("request-transformer"),
					Config: kong.Configuration{
						"add": "sensitive",
					},
				},
			},
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("request-transformer"),
					Config: kong.Configuration{
						"add": "{REDACTED}",
					},
				},
			},
			postCheck: func(t *testing.T, in Plugin, got Plugin) {
				got.Config["add"] = "changed"
				got.Config["new"] = "new-value"
				assert.Equal(t, "sensitive", in.Config["add"])
				_, exists := in.Config["new"]
				assert.False(t, exists)
			},
		},
		{
			name: "keeps empty config empty",
			in: Plugin{
				Plugin: kong.Plugin{
					Name:   kong.String("cors"),
					Config: kong.Configuration{},
				},
			},
			want: Plugin{
				Plugin: kong.Plugin{
					Name:   kong.String("cors"),
					Config: kong.Configuration{},
				},
			},
			postCheck: func(t *testing.T, in Plugin, got Plugin) {
				assert.Empty(t, in.Config)
				assert.Empty(t, got.Config)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.SanitizedCopy()

			assert.Equal(t, tt.want, got)

			if tt.postCheck != nil {
				tt.postCheck(t, tt.in, got)
			}
		})
	}
}
