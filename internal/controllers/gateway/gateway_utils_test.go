package gateway

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/scheme"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

func TestGetListenerSupportedRouteKinds(t *testing.T) {
	testCases := []struct {
		name                   string
		listener               gatewayapi.Listener
		expectedSupportedKinds []gatewayapi.RouteGroupKind
		resolvedRefsReason     gatewayapi.ListenerConditionReason
	}{
		{
			name: "only HTTP protocol specified",
			listener: gatewayapi.Listener{
				Protocol: gatewayapi.HTTPProtocolType,
			},
			expectedSupportedKinds: []gatewayapi.RouteGroupKind{
				builder.NewRouteGroupKind().HTTPRoute().Build(),
				builder.NewRouteGroupKind().GRPCRoute().Build(),
			},
			resolvedRefsReason: gatewayapi.ListenerReasonResolvedRefs,
		},
		{
			name: "only HTTPS protocol specified",
			listener: gatewayapi.Listener{
				Protocol: gatewayapi.HTTPSProtocolType,
			},
			expectedSupportedKinds: []gatewayapi.RouteGroupKind{
				builder.NewRouteGroupKind().HTTPRoute().Build(),
				builder.NewRouteGroupKind().GRPCRoute().Build(),
			},
			resolvedRefsReason: gatewayapi.ListenerReasonResolvedRefs,
		},
		{
			name: "only TCP protocol specified",
			listener: gatewayapi.Listener{
				Protocol: gatewayapi.TCPProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().TCPRoute().IntoSlice(),
			resolvedRefsReason:     gatewayapi.ListenerReasonResolvedRefs,
		},
		{
			name: "only UDP protocol specified",
			listener: gatewayapi.Listener{
				Protocol: gatewayapi.UDPProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().UDPRoute().IntoSlice(),
			resolvedRefsReason:     gatewayapi.ListenerReasonResolvedRefs,
		},
		{
			name: "only TLS protocol specified",
			listener: gatewayapi.Listener{
				Protocol: gatewayapi.TLSProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().TLSRoute().IntoSlice(),
			resolvedRefsReason:     gatewayapi.ListenerReasonResolvedRefs,
		},
		{
			name: "Kind not included in global gets discarded",
			listener: gatewayapi.Listener{
				Protocol: gatewayapi.HTTPProtocolType,
				AllowedRoutes: &gatewayapi.AllowedRoutes{
					Kinds: []gatewayapi.RouteGroupKind{
						{
							Group: lo.ToPtr(gatewayapi.Group("unknown.group.com")),
							Kind:  gatewayapi.Kind("UnknownKind"),
						},
						{
							Group: lo.ToPtr(gatewayapi.V1Group),
							Kind:  gatewayapi.Kind("HTTPRoute"),
						},
					},
				},
			},
			expectedSupportedKinds: []gatewayapi.RouteGroupKind{
				{
					Group: lo.ToPtr(gatewayapi.V1Group),
					Kind:  gatewayapi.Kind("HTTPRoute"),
				},
			},
			resolvedRefsReason: gatewayapi.ListenerReasonInvalidRouteKinds,
		},
		{
			name: "Kind included in global gets passed",
			listener: gatewayapi.Listener{
				Protocol: gatewayapi.HTTPProtocolType,
				AllowedRoutes: &gatewayapi.AllowedRoutes{
					Kinds: builder.NewRouteGroupKind().HTTPRoute().IntoSlice(),
				},
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().HTTPRoute().IntoSlice(),
			resolvedRefsReason:     gatewayapi.ListenerReasonResolvedRefs,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, reason := getListenerSupportedRouteKinds(tc.listener)
			require.Equal(t, tc.expectedSupportedKinds, got)
			require.Equal(t, tc.resolvedRefsReason, reason)
		})
	}
}

func TestGetListenerStatus(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	require.NoError(t, gatewayapi.InstallV1(scheme))

	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	testCases := []struct {
		name                     string
		gateway                  *gatewayapi.Gateway
		kongListens              []gatewayapi.Listener
		expectedListenerStatuses []gatewayapi.ListenerStatus
	}{
		{
			name: "only one listener",
			gateway: &gatewayapi.Gateway{
				TypeMeta: gatewayapi.V1GatewayTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name: "single-listener",
				},
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: "kong",
					Listeners: []gatewayapi.Listener{
						{
							Name:     "tcp-80",
							Port:     80,
							Protocol: gatewayapi.TCPProtocolType,
						},
					},
				},
			},
			kongListens: []gatewayapi.Listener{
				{
					Port:     80,
					Protocol: gatewayapi.TCPProtocolType,
				},
			},
			expectedListenerStatuses: []gatewayapi.ListenerStatus{
				{
					Name: gatewayapi.SectionName("tcp-80"),
					Conditions: []metav1.Condition{
						{
							Type:   string(gatewayapi.ListenerConditionAccepted),
							Status: metav1.ConditionTrue,
						},
					},
				},
			},
		},
		{
			name: "only one listener without a matching protocol or port",
			gateway: &gatewayapi.Gateway{
				TypeMeta: gatewayapi.V1GatewayTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name: "single-listener",
				},
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: "kong",
					Listeners: []gatewayapi.Listener{
						{
							Name:     "tcp-80",
							Port:     80,
							Protocol: gatewayapi.TLSProtocolType,
						},
					},
				},
			},
			kongListens: []gatewayapi.Listener{
				{
					Port:     80,
					Protocol: gatewayapi.TCPProtocolType,
				},
			},
			expectedListenerStatuses: []gatewayapi.ListenerStatus{
				{
					Name: gatewayapi.SectionName("tcp-80"),
					Conditions: []metav1.Condition{
						{
							Type:   string(gatewayapi.ListenerConditionAccepted),
							Status: metav1.ConditionFalse,
						},
					},
				},
			},
		},
		{
			name: "2 listeners, 1 with a matching protocol",
			gateway: &gatewayapi.Gateway{
				TypeMeta: gatewayapi.V1GatewayTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name: "single-listener",
				},
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: "kong",
					Listeners: []gatewayapi.Listener{
						{
							Name:     "tls-443",
							Port:     443,
							Protocol: gatewayapi.TLSProtocolType,
						},
						{
							Name:     "tcp-80",
							Port:     80,
							Protocol: gatewayapi.TCPProtocolType,
						},
					},
				},
			},
			kongListens: []gatewayapi.Listener{
				{
					Port:     80,
					Protocol: gatewayapi.TCPProtocolType,
				},
			},
			expectedListenerStatuses: []gatewayapi.ListenerStatus{
				{
					Name: gatewayapi.SectionName("tcp-80"),
					Conditions: []metav1.Condition{
						{
							Type:   string(gatewayapi.ListenerConditionAccepted),
							Status: metav1.ConditionTrue,
						},
					},
				},
				{
					Name: gatewayapi.SectionName("tls-443"),
					Conditions: []metav1.Condition{
						{
							Type:   string(gatewayapi.ListenerConditionAccepted),
							Status: metav1.ConditionFalse,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			statuses, err := getListenerStatus(ctx, tc.gateway, tc.kongListens, nil, client)
			require.NoError(t, err)
			require.Len(t, statuses, len(tc.expectedListenerStatuses), "should return expected number of listener statused")
			for _, expectedListenerStatus := range tc.expectedListenerStatuses {
				listenerStatus, ok := lo.Find(statuses, func(ls gatewayapi.ListenerStatus) bool {
					return ls.Name == expectedListenerStatus.Name
				})
				require.Truef(t, ok, "should find listener status of listener %s", expectedListenerStatus.Name)
				assertOnlyOneConditionForType(t, listenerStatus.Conditions)
				for _, expectedCondition := range expectedListenerStatus.Conditions {
					assert.Truef(t,
						lo.ContainsBy(listenerStatus.Conditions, func(c metav1.Condition) bool {
							return c.Type == expectedCondition.Type && c.Status == expectedCondition.Status
						}),
						"Condition %q should have Status %q: found listener conditions:\n%#v",
						expectedCondition.Type, expectedCondition.Status, listenerStatus.Conditions,
					)
				}
			}
		})
	}
}

func assertOnlyOneConditionForType(t *testing.T, conditions []metav1.Condition) {
	conditionsNum := lo.CountValuesBy(conditions, func(c metav1.Condition) string {
		return c.Type
	})
	for c, n := range conditionsNum {
		assert.Equalf(t, 1, n, "condition %s occurred %d times - expected 1 occurrence", c, n)
	}
}

func TestRouteAcceptedByGateways(t *testing.T) {
	testCases := []struct {
		name               string
		routeNamespace     string
		route              *gatewayapi.HTTPRoute
		gateways           []client.Object
		expectedGatewayNNs []k8stypes.NamespacedName
	}{
		{
			name:           "returns the gateway regardless of the route parent status conditions",
			routeNamespace: "default",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "route-1",
				},
				Status: gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ParentRef: gatewayapi.ParentReference{
									Name: "gateway-1",
								},
							},
						},
					},
				},
			},
			gateways: []client.Object{
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "gateway-1",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: builder.NewListener("http").WithPort(8080).HTTP().IntoSlice(),
					},
					Status: gatewayapi.GatewayStatus{
						Listeners: []gatewayapi.ListenerStatus{
							{
								Name: "http",
								SupportedKinds: []gatewayapi.RouteGroupKind{
									{
										Group: lo.ToPtr(gatewayapi.V1Group),
										Kind:  gatewayapi.Kind("HTTPRoute"),
									},
								},
							},
						},
					},
				},
			},
			expectedGatewayNNs: []k8stypes.NamespacedName{
				{
					Namespace: "default",
					Name:      "gateway-1",
				},
			},
		},
		{
			name:           "returns the gateway regardless of the route parent status conditions",
			routeNamespace: "default",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "route-1",
				},
				Status: gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ParentRef: gatewayapi.ParentReference{
									Name:  "gateway-1",
									Group: lo.ToPtr(gatewayapi.Group("wrong-group")),
								},
								Conditions: []metav1.Condition{
									{
										Status: metav1.ConditionTrue,
										Type:   string(gatewayapi.RouteConditionAccepted),
									},
								},
							},
							{
								ParentRef: gatewayapi.ParentReference{
									Name: "gateway-2",
									Kind: lo.ToPtr(gatewayapi.Kind("wrong-kind")),
								},
								Conditions: []metav1.Condition{
									{
										Status: metav1.ConditionTrue,
										Type:   string(gatewayapi.RouteConditionAccepted),
									},
								},
							},
							{
								ParentRef: gatewayapi.ParentReference{
									Name: "gateway-3",
								},
								Conditions: []metav1.Condition{
									{
										Status: metav1.ConditionTrue,
										Type:   string(gatewayapi.RouteConditionAccepted),
									},
								},
							},
							{
								ParentRef: gatewayapi.ParentReference{
									Name: "gateway-4",
								},
								Conditions: []metav1.Condition{
									{
										Status: metav1.ConditionFalse,
										Type:   string(gatewayapi.RouteConditionAccepted),
									},
								},
							},
						},
					},
				},
			},
			gateways: []client.Object{
				&gatewayapi.Gateway{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "wrong-group/v1",
						Kind:       "Gateway",
					},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "gateway-1",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: builder.NewListener("http").WithPort(8080).HTTP().IntoSlice(),
					},
					Status: gatewayapi.GatewayStatus{
						Listeners: []gatewayapi.ListenerStatus{
							{
								Name: "http",
								SupportedKinds: []gatewayapi.RouteGroupKind{
									{
										Group: lo.ToPtr(gatewayapi.V1Group),
										Kind:  gatewayapi.Kind("HTTPRoute"),
									},
								},
							},
						},
					},
				},
				&gatewayapi.Gateway{
					TypeMeta: metav1.TypeMeta{
						APIVersion: gatewayapi.GroupVersion.String(),
						Kind:       "wrong-kind",
					},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "gateway-2",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: builder.NewListener("http").WithPort(8080).HTTP().IntoSlice(),
					},
					Status: gatewayapi.GatewayStatus{
						Listeners: []gatewayapi.ListenerStatus{
							{
								Name: "http",
								SupportedKinds: []gatewayapi.RouteGroupKind{
									{
										Group: lo.ToPtr(gatewayapi.V1Group),
										Kind:  gatewayapi.Kind("HTTPRoute"),
									},
								},
							},
						},
					},
				},
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "gateway-3",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: builder.NewListener("http").WithPort(8080).HTTP().IntoSlice(),
					},
					Status: gatewayapi.GatewayStatus{
						Listeners: []gatewayapi.ListenerStatus{
							{
								Name: "http",
								SupportedKinds: []gatewayapi.RouteGroupKind{
									{
										Group: lo.ToPtr(gatewayapi.V1Group),
										Kind:  gatewayapi.Kind("HTTPRoute"),
									},
								},
							},
						},
					},
				},
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "gateway-4",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: builder.NewListener("http").WithPort(8080).HTTP().IntoSlice(),
					},
					Status: gatewayapi.GatewayStatus{
						Listeners: []gatewayapi.ListenerStatus{
							{
								Name: "http",
								SupportedKinds: []gatewayapi.RouteGroupKind{
									{
										Group: lo.ToPtr(gatewayapi.V1Group),
										Kind:  gatewayapi.Kind("HTTPRoute"),
									},
								},
							},
						},
					},
				},
			},
			expectedGatewayNNs: []k8stypes.NamespacedName{
				{
					Namespace: "default",
					Name:      "gateway-3",
				},
				{
					Namespace: "default",
					Name:      "gateway-4",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := scheme.Get()
			require.NoError(t, err)
			fakeClient := fake.NewClientBuilder().
				WithScheme(s).
				WithObjects(tc.gateways...).
				Build()

			gateways := routeAcceptedByGateways(
				context.Background(), fakeClient, logr.Discard(), tc.route,
			)
			assert.Equal(t, tc.expectedGatewayNNs, gateways)
		})
	}
}
