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
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
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
					"correlation-id-config": []byte(`{"header_name": "foo"}`),
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
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
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
				Protocols: kong.StringSlice("http"),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kongPluginFromK8SClusterPlugin(store, tt.args.plugin)
			if (err != nil) != tt.wantErr {
				t.Errorf("kongPluginFromK8SClusterPlugin error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(tt.want, got)
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
					"correlation-id-config": []byte(`{"header_name": "foo"}`),
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
					Protocols:  []configurationv1.KongProtocol{"http"},
					PluginName: "correlation-id",
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
				Protocols: kong.StringSlice("http"),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kongPluginFromK8SPlugin(store, tt.args.plugin)
			if (err != nil) != tt.wantErr {
				t.Errorf("kongPluginFromK8SPlugin error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(tt.want, got)
		})
	}
}

func Test_getKongIngressForServices(t *testing.T) {
	for _, tt := range []struct {
		name                string
		services            map[string]*corev1.Service
		kongIngresses       []*configurationv1.KongIngress
		expectedKongIngress *configurationv1.KongIngress
		expectedError       error
	}{
		{
			name: "when no services are provided, no KongIngress will be provided",
			kongIngresses: []*configurationv1.KongIngress{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongingress1",
					Namespace: corev1.NamespaceDefault,
				},
			}},
		},
		{
			name: "when none of the provided services have attached KongIngress resources, no KongIngress resources will be provided",
			services: map[string]*corev1.Service{
				"test-service1": {
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
					},
				},
				"test-service2": {
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			kongIngresses: []*configurationv1.KongIngress{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongingress1",
					Namespace: corev1.NamespaceDefault,
				},
			}},
		},
		{
			name: "if at least one KongIngress resource is attached to a Service, it will be returned",
			services: map[string]*corev1.Service{
				"test-service1": {
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress2",
						},
					},
				},
				"test-service2": {
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
					},
				},
			},
			kongIngresses: []*configurationv1.KongIngress{
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
			expectedKongIngress: &configurationv1.KongIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongingress2",
					Namespace: corev1.NamespaceDefault,
				},
			},
		},
		{
			name: "if multiple services have KongIngress resources this is accepted only if they're all attached to the same KongIngress",
			services: map[string]*corev1.Service{
				"test-service1": {
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress2",
						},
					},
				},
				"test-service2": {
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress2",
						},
					},
				},
				"test-service3": {
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service3",
						Namespace: corev1.NamespaceDefault,
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.ConfigurationKey: "test-kongingress2",
						},
					},
				},
			},
			kongIngresses: []*configurationv1.KongIngress{
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
			expectedKongIngress: &configurationv1.KongIngress{
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
				assert.Equal(t, tt.expectedKongIngress, kongIngress)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			}
		})
	}
}

func TestGetKongIngressFromObjectMeta(t *testing.T) {
	for _, tt := range []struct {
		name                string
		route               client.Object
		kongIngresses       []*configurationv1.KongIngress
		expectedKongIngress *configurationv1.KongIngress
		expectedError       error
	}{
		{
			name: "konghq.com/override annotation does not affect Gateway API's TCPRoute",
			route: &gatewayv1alpha2.TCPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "TCPRoute",
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
			kongIngresses: []*configurationv1.KongIngress{
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
			route: &gatewayv1alpha2.UDPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "UDPRoute",
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
			kongIngresses: []*configurationv1.KongIngress{
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
			route: &gatewayv1alpha2.HTTPRoute{
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
			kongIngresses: []*configurationv1.KongIngress{
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
			route: &gatewayv1alpha2.TLSRoute{
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
			kongIngresses: []*configurationv1.KongIngress{
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

func Test_prettyPrintServiceList(t *testing.T) {
	for _, tt := range []struct {
		name     string
		services map[string]*corev1.Service
		expected string
	}{
		{
			name:     "an empty list of services produces an empty string",
			expected: "",
		},
		{
			name: "a single service should just return the <namespace>/<name>",
			services: map[string]*corev1.Service{
				"test-service1": {
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
			services: map[string]*corev1.Service{
				"test-service1": {
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
					},
				},
				"test-service2": {
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service2",
						Namespace: corev1.NamespaceDefault,
					},
				},
				"test-service3": {
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
			assert.True(t, re.MatchString(PrettyPrintServiceList(tt.services)))
		})
	}
}
