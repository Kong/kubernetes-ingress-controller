package kongstate

import (
	"regexp"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
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
					"correlation-id-config":    []byte(`{"header_name": "foo"}`),
					"correlation-id-generator": []byte(`"uuid"`),
					"correlation-id-invalid":   []byte(`"aaa`),
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
			name: "empty config and configPatch for whole object",
			args: args{
				plugin: kongv1.KongClusterPlugin{
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{}`),
					},
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
					"correlation-id-config":    []byte(`{"header_name": "foo"}`),
					"correlation-id-generator": []byte(`"uuid"`),
					"correlation-id-invalid":   []byte(`"aaa`),
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
			name: "empty config and configPatch for whole object",
			args: args{
				plugin: kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Protocols:  []kongv1.KongProtocol{"http"},
					PluginName: "correlation-id",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{}`),
					},
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

func TestGetKongIngressForServices(t *testing.T) {
	for _, tt := range []struct {
		name                string
		services            []*corev1.Service
		kongIngresses       []*kongv1.KongIngress
		expectedKongIngress *kongv1.KongIngress
		expectedError       error
	}{
		{
			name: "when no services are provided, no KongIngress will be provided",
			kongIngresses: []*kongv1.KongIngress{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongingress1",
					Namespace: corev1.NamespaceDefault,
				},
			}},
		},
		{
			name: "when none of the provided services have attached KongIngress resources, no KongIngress resources will be provided",
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			kongIngresses: []*kongv1.KongIngress{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongingress1",
					Namespace: corev1.NamespaceDefault,
				},
			}},
		},
		{
			name: "if at least one KongIngress resource is attached to a Service, it will be returned",
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress2",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			kongIngresses: []*kongv1.KongIngress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kongingress1",
						Namespace: corev1.NamespaceDefault,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kongingress2",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			expectedKongIngress: &kongv1.KongIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongingress2",
					Namespace: corev1.NamespaceDefault,
				},
			},
		},
		{
			name: "if multiple services have KongIngress resources this is accepted only if they're all attached to the same KongIngress",
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress2",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress2",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service3",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress2",
						},
					},
				},
			},
			kongIngresses: []*kongv1.KongIngress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kongingress1",
						Namespace: corev1.NamespaceDefault,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kongingress2",
						Namespace: corev1.NamespaceDefault,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kongingress3",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			expectedKongIngress: &kongv1.KongIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongingress2",
					Namespace: corev1.NamespaceDefault,
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			storer, err := store.NewFakeStore(store.FakeObjects{
				KongIngresses: tt.kongIngresses,
			})
			require.NoError(t, err)

			kongIngress, err := getKongIngressForServices(storer, tt.services)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedKongIngress, kongIngress)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
			}
		})
	}
}

func TestGetKongIngressFromObjectMeta(t *testing.T) {
	for _, tt := range []struct {
		name                string
		route               client.Object
		kongIngresses       []*kongv1.KongIngress
		expectedKongIngress *kongv1.KongIngress
		expectedError       error
	}{
		{
			name: "konghq.com/override annotation does not affect Gateway API's TCPRoute",
			route: &gatewayapi.TCPRoute{
				TypeMeta: gatewayapi.TCPRouteTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "muntaxabi-jugrofiyai-umumiy",
					Namespace: "behbudiy",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress1",
					},
				},
			},
			kongIngresses: []*kongv1.KongIngress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kongingress1",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			expectedKongIngress: nil,
		},
		{
			name: "konghq.com/override annotation does not affect Gateway API's UDPRoute",
			route: &gatewayapi.UDPRoute{
				TypeMeta: gatewayapi.UDPRouteTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "muntaxabi-jugrofiyai-umumiy",
					Namespace: "behbudiy",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress1",
					},
				},
			},
			kongIngresses: []*kongv1.KongIngress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kongingress1",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			expectedKongIngress: nil,
		},
		{
			name: "konghq.com/override annotation does not affect Gateway API's HTTPRoute",
			route: &gatewayapi.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "HTTPRoute",
					APIVersion: "gateway.networking.k8s.io/v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "muntaxabi-jugrofiyai-umumiy",
					Namespace: "behbudiy",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress1",
					},
				},
			},
			kongIngresses: []*kongv1.KongIngress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kongingress1",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			expectedKongIngress: nil,
		},
		{
			name: "konghq.com/override annotation does not affect Gateway API's TLSRoute",
			route: &gatewayapi.TLSRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "TLSRoute",
					APIVersion: "gateway.networking.k8s.io/v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "muntaxabi-jugrofiyai-umumiy",
					Namespace: "behbudiy",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress1",
					},
				},
			},
			kongIngresses: []*kongv1.KongIngress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kongingress1",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			expectedKongIngress: nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			storer, err := store.NewFakeStore(store.FakeObjects{
				KongIngresses: tt.kongIngresses,
			})
			require.NoError(t, err)

			obj := util.FromK8sObject(tt.route)
			kongIngress, err := getKongIngressFromObjectMeta(storer, obj)

			if tt.expectedError == nil {
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedKongIngress, kongIngress)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestPrettyPrintServiceList(t *testing.T) {
	for _, tt := range []struct {
		name     string
		services []*corev1.Service
		expected string
	}{
		{
			name:     "an empty list of services produces an empty string",
			expected: "",
		},
		{
			name: "a single service should just return the <namespace>/<name>",
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			expected: "default/test-service1",
		},
		{
			name: "multiple services should be comma deliniated",
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service3",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			expected: "default/test-service[0-9], default/test-service[0-9], default/test-service[0-9]",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			re := regexp.MustCompile(tt.expected)
			require.True(t, re.MatchString(prettyPrintServiceList(tt.services)))
		})
	}
}
