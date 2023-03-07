package gateway

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
)

func TestGetListenerSupportedRouteKinds(t *testing.T) {
	testCases := []struct {
		name                   string
		listener               Listener
		expectedSupportedKinds []gatewayv1beta1.RouteGroupKind
		resolvedRefsReason     gatewayv1beta1.ListenerConditionReason
	}{
		{
			name: "only HTTP protocol specified",
			listener: Listener{
				Protocol: HTTPProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().HTTPRoute().IntoSlice(),
			resolvedRefsReason:     gatewayv1beta1.ListenerReasonResolvedRefs,
		},
		{
			name: "only HTTPS protocol specified",
			listener: Listener{
				Protocol: HTTPSProtocolType,
			},
			expectedSupportedKinds: []gatewayv1beta1.RouteGroupKind{
				builder.NewRouteGroupKind().HTTPRoute().Build(),
				builder.NewRouteGroupKind().GRPCRoute().Build(),
			},
			resolvedRefsReason: gatewayv1beta1.ListenerReasonResolvedRefs,
		},
		{
			name: "only TCP protocol specified",
			listener: Listener{
				Protocol: TCPProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().TCPRoute().IntoSlice(),
			resolvedRefsReason:     gatewayv1beta1.ListenerReasonResolvedRefs,
		},
		{
			name: "only UDP protocol specified",
			listener: Listener{
				Protocol: UDPProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().UDPRoute().IntoSlice(),
			resolvedRefsReason:     gatewayv1beta1.ListenerReasonResolvedRefs,
		},
		{
			name: "only TLS protocol specified",
			listener: Listener{
				Protocol: TLSProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().TLSRoute().IntoSlice(),
			resolvedRefsReason:     gatewayv1beta1.ListenerReasonResolvedRefs,
		},
		{
			name: "Kind not included in global gets discarded",
			listener: Listener{
				Protocol: HTTPProtocolType,
				AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
					Kinds: []gatewayv1beta1.RouteGroupKind{
						{
							Group: lo.ToPtr(gatewayv1beta1.Group("unknown.group.com")),
							Kind:  Kind("UnknownKind"),
						},
						{
							Group: &gatewayV1beta1Group,
							Kind:  Kind("HTTPRoute"),
						},
					},
				},
			},
			expectedSupportedKinds: []gatewayv1beta1.RouteGroupKind{
				{
					Group: &gatewayV1beta1Group,
					Kind:  Kind("HTTPRoute"),
				},
			},
			resolvedRefsReason: gatewayv1beta1.ListenerReasonInvalidRouteKinds,
		},
		{
			name: "Kind included in global gets passed",
			listener: Listener{
				Protocol: HTTPProtocolType,
				AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
					Kinds: builder.NewRouteGroupKind().HTTPRoute().IntoSlice(),
				},
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().HTTPRoute().IntoSlice(),
			resolvedRefsReason:     gatewayv1beta1.ListenerReasonResolvedRefs,
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

	statuses, err := getListenerStatus(ctx, &Gateway{
		Spec: gatewayv1beta1.GatewaySpec{
			GatewayClassName: "kong",
			Listeners: []gatewayv1beta1.Listener{
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
