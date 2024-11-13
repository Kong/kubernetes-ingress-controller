package kongstate

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

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

			kongIngress, err := getKongIngressFromObjectMeta(storer, tt.route)

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
