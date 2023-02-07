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
	}{
		{
			name: "only HTTP protocol specified",
			listener: Listener{
				Protocol: HTTPProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().HTTPRoute().IntoSlice(),
		},
		{
			name: "only HTTPS protocol specified",
			listener: Listener{
				Protocol: HTTPSProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().HTTPRoute().IntoSlice(),
		},
		{
			name: "only TCP protocol specified",
			listener: Listener{
				Protocol: TCPProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().TCPRoute().IntoSlice(),
		},
		{
			name: "only UDP protocol specified",
			listener: Listener{
				Protocol: UDPProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().UDPRoute().IntoSlice(),
		},
		{
			name: "only TLS protocol specified",
			listener: Listener{
				Protocol: TLSProtocolType,
			},
			expectedSupportedKinds: builder.NewRouteGroupKind().TLSRoute().IntoSlice(),
		},
		{
			name: "Kind not included in global gets discarded",
			listener: Listener{
				Protocol: HTTPProtocolType,
				AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
					Kinds: []gatewayv1beta1.RouteGroupKind{{
						Group: lo.ToPtr(gatewayv1beta1.Group("unknown.group.com")),
						Kind:  Kind("UnknownKind"),
					}},
				},
			},
			expectedSupportedKinds: nil,
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
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := getListenerSupportedRouteKinds(tc.listener)
			require.Equal(t, tc.expectedSupportedKinds, got)
		})
	}
}

func TestGetListenerStatus_no_duplicated_Detached_condition(t *testing.T) {
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
	assertOnlyOneConditionOfType(t, listenerStatus.Conditions, gatewayv1beta1.ListenerConditionAccepted)
}

func assertOnlyOneConditionOfType(t *testing.T, conditions []metav1.Condition, typ gatewayv1beta1.ListenerConditionType) {
	conditionNum := 0
	for _, condition := range conditions {
		if condition.Type == string(typ) {
			conditionNum++
		}
	}
	assert.Equal(t, 1, conditionNum)
}
