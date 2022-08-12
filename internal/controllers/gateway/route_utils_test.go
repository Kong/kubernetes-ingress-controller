package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func Test_filterHostnames(t *testing.T) {
	commonGateway := &gatewayv1alpha2.Gateway{
		Spec: gatewayv1alpha2.GatewaySpec{
			Listeners: []gatewayv1alpha2.Listener{
				{
					Name:     "listener-1",
					Hostname: util.StringToGatewayAPIHostnamePtr("very.specific.com"),
				},
				{
					Name:     "listener-2",
					Hostname: util.StringToGatewayAPIHostnamePtr("*.wildcard.io"),
				},
				{
					Name:     "listener-3",
					Hostname: util.StringToGatewayAPIHostnamePtr("*.anotherwildcard.io"),
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
						util.StringToGatewayAPIHostname("*.anotherwildcard.io"),
						util.StringToGatewayAPIHostname("*.nonmatchingwildcard.io"),
						util.StringToGatewayAPIHostname("very.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						util.StringToGatewayAPIHostname("very.specific.com"),
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
						util.StringToGatewayAPIHostname("non.matching.com"),
						util.StringToGatewayAPIHostname("*.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						util.StringToGatewayAPIHostname("very.specific.com"),
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
						util.StringToGatewayAPIHostname("non.matching.com"),
						util.StringToGatewayAPIHostname("wildcard.io"),
						util.StringToGatewayAPIHostname("foo.wildcard.io"),
						util.StringToGatewayAPIHostname("bar.wildcard.io"),
						util.StringToGatewayAPIHostname("foo.bar.wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						util.StringToGatewayAPIHostname("foo.wildcard.io"),
						util.StringToGatewayAPIHostname("bar.wildcard.io"),
						util.StringToGatewayAPIHostname("foo.bar.wildcard.io"),
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
						util.StringToGatewayAPIHostname("*.anotherwildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1alpha2.HTTPRoute{
				Spec: gatewayv1alpha2.HTTPRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						util.StringToGatewayAPIHostname("*.anotherwildcard.io"),
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
						util.StringToGatewayAPIHostname("specific.but.wrong.com"),
						util.StringToGatewayAPIHostname("wildcard.io"),
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
