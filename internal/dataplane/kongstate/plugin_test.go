package kongstate

import (
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
)

func TestKongPluginFromK8SClusterPlugin(t *testing.T) {
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

	tests := []struct {
		name    string
		plugin  kongv1.KongClusterPlugin
		want    Plugin
		wantErr bool
	}{
		{
			name: "basic configuration",
			plugin: kongv1.KongClusterPlugin{
				Protocols:    []kongv1.KongProtocol{"http"},
				PluginName:   "correlation-id",
				InstanceName: "example",
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"header_name": "foo"}`),
				},
			},
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("correlation-id"),
					Config: kong.Configuration{
						"header_name": "foo",
					},
					Protocols:    kong.StringSlice("http"),
					InstanceName: kong.String("example"),
				},
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{JSONPaths: []string{}},
			},
			wantErr: false,
		},
		{
			name: "secret configuration",
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
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("correlation-id"),
					Config: kong.Configuration{
						"header_name": "foo",
					},
					Protocols: kong.StringSlice("http"),
				},
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
					WholeConfigIsSensitive: true,
					JSONPaths:              []string{},
				},
			},
			wantErr: false,
		},
		{
			name: "missing secret configuration",
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
			want:    Plugin{},
			wantErr: true,
		},
		{
			name: "non-JSON configuration",
			plugin: kongv1.KongClusterPlugin{
				Protocols:  []kongv1.KongProtocol{"http"},
				PluginName: "correlation-id",
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{{}`),
				},
			},
			want:    Plugin{},
			wantErr: true,
		},
		{
			name: "both Config and ConfigFrom set",
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
			want:    Plugin{},
			wantErr: true,
		},
		{
			name: "Config and ConfigPatches set",
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
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("correlation-id"),
					Config: kong.Configuration{
						"header_name": "foo",
						"generator":   "uuid",
					},
					Protocols: kong.StringSlice("http"),
				},
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
					JSONPaths: []string{"/generator"},
				},
			},
		},
		{
			name: "configPatch on subpath of non-exist path",
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
			want: Plugin{
				Plugin: kong.Plugin{
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
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
					JSONPaths: []string{"/add/headers"},
				},
			},
		},
		{
			name: "empty config and configPatch for particular paths",
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
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("correlation-id"),
					Config: kong.Configuration{
						"header_name": "foo",
						"generator":   "uuid",
					},
					Protocols: kong.StringSlice("http"),
				},
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
					JSONPaths: []string{"/header_name", "/generator"},
				},
			},
		},
		{
			name: "empty config and configPatch for whole object",
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
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("correlation-id"),
					Config: kong.Configuration{
						"header_name": "foo",
					},
					Protocols: kong.StringSlice("http"),
				},
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
					JSONPaths: []string{""},
				},
			},
		},
		{
			name: "missing secret in configPatches",
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
			wantErr: true,
		},
		{
			name: "missing key of secret in cofigPatches",
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
			wantErr: true,
		},
		{
			name: "invalid value in configPatches",
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kongPluginFromK8SClusterPlugin(store, tt.plugin)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			tt.want.K8sParent = tt.plugin.DeepCopy()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKongPluginFromK8SPlugin(t *testing.T) {
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
	tests := []struct {
		name    string
		plugin  kongv1.KongPlugin
		want    Plugin
		wantErr bool
	}{
		{
			name: "basic configuration",
			plugin: kongv1.KongPlugin{
				Protocols:    []kongv1.KongProtocol{"http"},
				PluginName:   "correlation-id",
				InstanceName: "example",
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"header_name": "foo"}`),
				},
			},
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("correlation-id"),
					Config: kong.Configuration{
						"header_name": "foo",
					},
					Protocols:    kong.StringSlice("http"),
					InstanceName: kong.String("example"),
				},
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{JSONPaths: []string{}},
			},
			wantErr: false,
		},
		{
			name: "secret configuration",
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
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("correlation-id"),
					Config: kong.Configuration{
						"header_name": "foo",
					},
					Protocols: kong.StringSlice("http"),
				},
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
					WholeConfigIsSensitive: true,
					JSONPaths:              []string{},
				},
			},
			wantErr: false,
		},
		{
			name: "missing secret configuration",
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
			wantErr: true,
		},
		{
			name: "non-JSON configuration",
			plugin: kongv1.KongPlugin{
				Protocols:  []kongv1.KongProtocol{"http"},
				PluginName: "correlation-id",
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{{}`),
				},
			},
			wantErr: true,
		},
		{
			name: "both Config and ConfigFrom set",
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
			wantErr: true,
		},
		{
			name: "config and configPatches set",
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
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("correlation-id"),
					Config: kong.Configuration{
						"header_name": "foo",
						"generator":   "uuid",
					},
					Protocols: kong.StringSlice("http"),
				},
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
					JSONPaths: []string{"/generator"},
				},
			},
		},
		{
			name: "configPatch on subpath of non-exist path",
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
			want: Plugin{
				Plugin: kong.Plugin{
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
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
					JSONPaths: []string{"/add/headers"},
				},
			},
		},
		{
			name: "empty config and configPatch for particular paths",
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
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("correlation-id"),
					Config: kong.Configuration{
						"header_name": "foo",
						"generator":   "uuid",
					},
					Protocols: kong.StringSlice("http"),
				},
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
					JSONPaths: []string{"/header_name", "/generator"},
				},
			},
		},
		{
			name: "empty config and configPatch for whole object",
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
			want: Plugin{
				Plugin: kong.Plugin{
					Name: kong.String("correlation-id"),
					Config: kong.Configuration{
						"header_name": "foo",
					},
					Protocols: kong.StringSlice("http"),
				},
				SensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
					JSONPaths: []string{""},
				},
			},
		},
		{
			name: "missing secret in configPatches",
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
			wantErr: true,
		},
		{
			name: "missing key of secret in configPatches",
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
			wantErr: true,
		},
		{
			name: "invalid value in configPatches",
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kongPluginFromK8SPlugin(store, tt.plugin)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			// don't care about tags in this test
			got.Tags = nil
			tt.want.K8sParent = tt.plugin.DeepCopy()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPlugin_SanitizedCopy(t *testing.T) {
	// this needs a static random seed because some auths generate random values
	uuid.SetRand(rand.New(rand.NewSource(1))) //nolint:gosec
	testCases := []struct {
		name                    string
		config                  kong.Configuration
		sensitiveFieldsMeta     PluginSensitiveFieldsMetadata
		expectedSanitizedConfig kong.Configuration
	}{
		{
			name: "sensitive fields are redacted with JSONPaths",
			config: kong.Configuration{
				"secret": "secret-value",
				"object": map[string]interface{}{
					"secretObjectField": "secret-object-field-value",
				},
			},
			sensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
				JSONPaths: []string{
					"/secret",
					"/object/secretObjectField",
				},
			},
			expectedSanitizedConfig: kong.Configuration{
				"secret": "{vault://redacted-value}",
				"object": map[string]interface{}{
					"secretObjectField": "{vault://redacted-value}",
				},
			},
		},
		{
			name: "invalid JSONPath doesn't panic and redacts whole config as fallback",
			config: kong.Configuration{
				"secret": "secret-value",
				"object": map[string]interface{}{
					"secretObjectField": "secret-object-field-value",
				},
			},
			sensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
				JSONPaths: []string{
					"/not-existing-path",
				},
			},
			expectedSanitizedConfig: kong.Configuration{
				"secret": "{vault://redacted-value}",
				"object": "{vault://redacted-value}",
			},
		},
		{
			name: "whole config to sanitize",
			config: kong.Configuration{
				"secret": "secret-value",
				"object": map[string]interface{}{
					"secretObjectField": "secret-object-field-value",
				},
			},
			sensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
				WholeConfigIsSensitive: true,
			},
			expectedSanitizedConfig: kong.Configuration{
				"secret": "{vault://redacted-value}",
				"object": "{vault://redacted-value}",
			},
		},
		{
			name: "single empty JSON path - whole config is redacted",
			config: kong.Configuration{
				"secret": "secret-value",
				"object": map[string]interface{}{
					"secretObjectField": "secret-object-field-value",
				},
			},
			sensitiveFieldsMeta: PluginSensitiveFieldsMetadata{
				JSONPaths: []string{""},
			},
			expectedSanitizedConfig: kong.Configuration{
				"secret": "{vault://redacted-value}",
				"object": "{vault://redacted-value}",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Plugin{
				Plugin: kong.Plugin{
					Config: tc.config,
				},
				SensitiveFieldsMeta: tc.sensitiveFieldsMeta,
			}
			sanitized := p.SanitizedCopy()
			assert.Equal(t, tc.expectedSanitizedConfig, sanitized.Config)
		})
	}
}
