package gateway

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

func init() {
	if err := gatewayv1.Install(scheme.Scheme); err != nil {
		panic(err)
	}
}

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
			expectedSupportedKinds: builder.NewRouteGroupKind().HTTPRoute().IntoSlice(),
			resolvedRefsReason:     gatewayapi.ListenerReasonResolvedRefs,
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
							Group: &gatewayV1Group,
							Kind:  gatewayapi.Kind("HTTPRoute"),
						},
					},
				},
			},
			expectedSupportedKinds: []gatewayapi.RouteGroupKind{
				{
					Group: &gatewayV1Group,
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

func TestGetListenerStatus_no_duplicated_condition(t *testing.T) {
	ctx := context.Background()
	client := fake.NewClientBuilder().Build()

	statuses, err := getListenerStatus(ctx, &gatewayapi.Gateway{
		Spec: gatewayapi.GatewaySpec{
			GatewayClassName: "kong",
			Listeners: []gatewayapi.Listener{
				{
					Port:     80,
					Protocol: "TCP",
				},
			},
		},
	}, nil, nil, client)
	require.NoError(t, err)
	require.Len(t, statuses, 1, "only one listener status expected as only one listener was defined")
	listenerStatus := statuses[0]
	assertOnlyOneConditionForType(t, listenerStatus.Conditions)
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
		name           string
		routeNamespace string
		parentStatuses []gatewayapi.RouteParentStatus
		gateways       []k8stypes.NamespacedName
	}{
		{
			name:           "no parentStatus with accepted condition",
			routeNamespace: "default",
			parentStatuses: []gatewayapi.RouteParentStatus{
				{
					ParentRef: gatewayapi.ParentReference{
						Name: "gateway-1",
					},
				},
			},
			gateways: []k8stypes.NamespacedName{
				// Gateways should be included even when route is not accepted.
				{
					Namespace: "default",
					Name:      "gateway-1",
				},
			},
		},
		{
			name:           "a subset of parentStatus with correct params",
			routeNamespace: "default",
			parentStatuses: []gatewayapi.RouteParentStatus{
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
			gateways: []k8stypes.NamespacedName{
				{
					Namespace: "default",
					Name:      "gateway-3",
				},
				// Gateways should be included even when route is not accepted.
				{
					Namespace: "default",
					Name:      "gateway-4",
				},
			},
		},
		{
			name:           "all parentStatuses",
			routeNamespace: "default",
			parentStatuses: []gatewayapi.RouteParentStatus{
				{
					ParentRef: gatewayapi.ParentReference{
						Name: "gateway-1",
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
						Name:      "gateway-2",
						Namespace: lo.ToPtr(gatewayapi.Namespace("namespace-2")),
					},
					Conditions: []metav1.Condition{
						{
							Status: metav1.ConditionTrue,
							Type:   string(gatewayapi.RouteConditionAccepted),
						},
					},
				},
			},
			gateways: []k8stypes.NamespacedName{
				{
					Namespace: "default",
					Name:      "gateway-1",
				},
				{
					Namespace: "namespace-2",
					Name:      "gateway-2",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gateways := routeAcceptedByGateways(tc.routeNamespace, tc.parentStatuses)
			assert.Equal(t, tc.gateways, gateways)
		})
	}
}
