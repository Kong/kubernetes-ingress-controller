package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func Test_filterHostnames(t *testing.T) {
	commonGateway := &gatewayv1beta1.Gateway{
		Spec: gatewayv1beta1.GatewaySpec{
			Listeners: []Listener{
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
		httpRoute         *gatewayv1beta1.HTTPRoute
		expectedHTTPRoute *gatewayv1beta1.HTTPRoute
	}{
		{
			name: "listener 1 - specific",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listener-1",
				},
			},
			httpRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("*.anotherwildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("*.nonmatchingwildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("very.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("very.specific.com"),
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
			httpRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("non.matching.com"),
						util.StringToGatewayAPIHostnameV1Beta1("*.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("very.specific.com"),
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
			httpRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("non.matching.com"),
						util.StringToGatewayAPIHostnameV1Beta1("wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("foo.wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("bar.wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("foo.bar.wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("foo.wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("bar.wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("foo.bar.wildcard.io"),
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
			httpRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("*.anotherwildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("*.anotherwildcard.io"),
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
			httpRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("specific.but.wrong.com"),
						util.StringToGatewayAPIHostnameV1Beta1("wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{},
				},
			},
		},
	}

	for _, tc := range testCases {
		filteredHTTPRoute := filterHostnames(tc.gateways, tc.httpRoute)
		assert.Equal(t, tc.expectedHTTPRoute.Spec, filteredHTTPRoute.Spec, tc.name)
	}
}
