package subtranslator

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

func TestBackendRefsToKongStateBackends(t *testing.T) {
	testcases := []struct {
		name        string
		route       client.Object
		backendRefs []gatewayapi.BackendRef
		allowed     map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo
		objects     store.FakeObjects
		expected    kongstate.ServiceBackends
	}{
		{
			name: "correct ReferenceGrant and an existing Service as backendRef returns a KongStateBackend with a Service",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						// BackendRefs are taken from the backendRefs argument.
						// This is done this way because the backendRefs are
						// extracted from the route's spec before.
						// Potentially this could be refactored to be extracted
						// in backendRefsToKongStateBackends using a type switch.
					}},
				},
			},
			backendRefs: []gatewayapi.BackendRef{
				builder.NewBackendRef("fake-service").WithPort(80).Build(),
			},
			allowed: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{
				gatewayapi.Namespace(corev1.NamespaceDefault): {
					{
						Group: "",
						Kind:  "Service",
						Name:  lo.ToPtr(gatewayapi.ObjectName("fake-service")),
					},
				},
			},
			objects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "fake-service",
							Namespace: corev1.NamespaceDefault,
						},
						Spec: corev1.ServiceSpec{
							Ports: []corev1.ServicePort{
								builder.NewServicePort().WithPort(80).Build(),
							},
						},
					},
				},
			},
			expected: func() kongstate.ServiceBackends {
				svcBackend, err := kongstate.NewServiceBackend(
					kongstate.ServiceBackendTypeKubernetesService,
					k8stypes.NamespacedName{Namespace: corev1.NamespaceDefault, Name: "fake-service"},
					kongstate.PortDef{
						Mode:   kongstate.PortModeByNumber,
						Number: 80,
					},
				)
				require.NoError(t, err)
				return kongstate.ServiceBackends{svcBackend}
			}(),
		},
		{
			name: "no ReferenceGrant and an existing Service as backendRef doesn't return a KongStateBackend with the Service",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
				},
			},
			allowed: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{
				gatewayapi.Namespace(corev1.NamespaceDefault): {
					{
						Group: "",
						Kind:  "ImaginaryKind",
						Name:  lo.ToPtr(gatewayapi.ObjectName("fake-service")),
					},
				},
			},
			objects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "fake-service",
							Namespace: corev1.NamespaceDefault,
						},
					},
				},
			},
			expected: kongstate.ServiceBackends{},
		},
		{
			name: "ReferenceGrant and a non existing Service as backendRef doesn't return a KongStateBackend with the Service",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
				},
			},
			backendRefs: []gatewayapi.BackendRef{
				builder.NewBackendRef("fake-service").WithPort(80).Build(),
			},
			allowed: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{
				gatewayapi.Namespace(corev1.NamespaceDefault): {
					{
						Group: "",
						Kind:  "Service",
						Name:  lo.ToPtr(gatewayapi.ObjectName("fake-service")),
					},
				},
			},
			expected: kongstate.ServiceBackends{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			fakestore, err := store.NewFakeStore(tc.objects)
			require.NoError(t, err)
			logger := logr.Discard()
			ret := backendRefsToKongStateBackends(logger, fakestore, tc.route, tc.backendRefs, tc.allowed)
			require.Equal(t, tc.expected, ret)
		})
	}
}

func commonRouteSpecMock(parentReferentName string) gatewayapi.CommonRouteSpec {
	return gatewayapi.CommonRouteSpec{
		ParentRefs: []gatewayapi.ParentReference{{
			Name: gatewayapi.ObjectName(parentReferentName),
		}},
	}
}
