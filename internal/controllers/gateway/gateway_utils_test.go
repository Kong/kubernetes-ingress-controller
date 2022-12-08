package gateway

import (
	"testing"

	"github.com/stretchr/testify/require"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/address"
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
						Group: address.Of(gatewayv1beta1.Group("unknown.group.com")),
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
