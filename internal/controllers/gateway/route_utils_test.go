package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

func Test_filterHostnames(t *testing.T) {
	commonGateway := &gatewayv1alpha2.Gateway{
		Spec: gatewayv1alpha2.GatewaySpec{
			Listeners: []gatewayv1alpha2.Listener{
				{
					Name:     "listener-1",
					Hostname: (*gatewayv1alpha2.Hostname)(pointer.StringPtr("very.specific.com")),
				},
				{
					Name:     "listener-2",
					Hostname: (*gatewayv1alpha2.Hostname)(pointer.StringPtr("*.wildcard.io")),
				},
				{
					Name:     "listener-3",
					Hostname: (*gatewayv1alpha2.Hostname)(pointer.StringPtr("*.anotherwildcard.io")),
				},
			},
		},
	}

	testCases := []struct {
		name              string
		gateways          []supportedGatewayWithCondition
		httpRoute         *gatewayv1alpha2.HTTPRoute
		expectedHTTPRoute *gatewayv1alpha2.HTTPRoute
	}{
		{
			name: "listener 1 - specific",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listener-1",
				},
			},
			httpRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						(gatewayv1alpha2.Hostname)("*.anotherwildcard.io"),
						(gatewayv1alpha2.Hostname)("*.nonmatchingwildcard.io"),
						(gatewayv1alpha2.Hostname)("very.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						(gatewayv1alpha2.Hostname)("very.specific.com"),
					},
				},
			},
		},
		{
			name: "listener 1 - wildcard",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listener-1",
				},
			},
			httpRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						(gatewayv1alpha2.Hostname)("non.matching.com"),
						(gatewayv1alpha2.Hostname)("*.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						(gatewayv1alpha2.Hostname)("very.specific.com"),
					},
				},
			},
		},
		{
			name: "listener 2",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listener-2",
				},
			},
			httpRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						(gatewayv1alpha2.Hostname)("non.matching.com"),
						(gatewayv1alpha2.Hostname)("wildcard.io"),
						(gatewayv1alpha2.Hostname)("foo.wildcard.io"),
						(gatewayv1alpha2.Hostname)("bar.wildcard.io"),
						(gatewayv1alpha2.Hostname)("foo.bar.wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						(gatewayv1alpha2.Hostname)("foo.wildcard.io"),
						(gatewayv1alpha2.Hostname)("bar.wildcard.io"),
						(gatewayv1alpha2.Hostname)("foo.bar.wildcard.io"),
					},
				},
			},
		},
		{
			name: "listener 3 - wildcard",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listener-3",
				},
			},
			httpRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						(gatewayv1alpha2.Hostname)("*.anotherwildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						(gatewayv1alpha2.Hostname)("*.anotherwildcard.io"),
					},
				},
			},
		},
		{
			name: "no match",
			gateways: []supportedGatewayWithCondition{
				{
					gateway: commonGateway,
				},
			},
			httpRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						(gatewayv1alpha2.Hostname)("specific.but.wrong.com"),
						(gatewayv1alpha2.Hostname)("wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{},
				},
			},
		},
	}

	for _, tc := range testCases {
		filteredHTTPRoute := filterHostnames(tc.gateways, tc.httpRoute)
		assert.Equal(t, tc.expectedHTTPRoute.Spec, filteredHTTPRoute.Spec, tc.name)
	}
}
