package translator

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func TestGenerateKongRoutesFromRouteRule_TCP(t *testing.T) {
	testCases := []struct {
		name      string
		route     *gatewayapi.TCPRoute
		gwPorts   []gatewayapi.PortNumber
		routeRule gatewayapi.TCPRouteRule
		expected  []kongstate.Route
	}{
		{
			name: "TCPRoute gets translated correctly to kong.Route",
			route: &gatewayapi.TCPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mytcproute-name",
					Namespace: "mynamespace",
				},
			},
			gwPorts: []gatewayapi.PortNumber{8080},
			routeRule: gatewayapi.TCPRouteRule{
				BackendRefs: []gatewayapi.BackendRef{
					{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Port: lo.ToPtr(gatewayapi.PortNumber(1234)),
						},
					},
				},
			},
			expected: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:      "mytcproute-name",
						Namespace: "mynamespace",
					},
					Route: kong.Route{
						Name: lo.ToPtr("tcproute.mynamespace.mytcproute-name.0.0"),
						Destinations: []*kong.CIDRPort{
							{
								Port: lo.ToPtr(8080),
							},
						},
						Protocols: []*string{
							lo.ToPtr("tcp"),
						},
						Tags: []*string{
							kong.String("k8s-name:mytcproute-name"),
							kong.String("k8s-namespace:mynamespace"),
						},
					},
				},
			},
		},
		{
			name: "TCPRoute with multiple backends and different Gateway port get translated correctly to kong.Route",
			route: &gatewayapi.TCPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mytcproute-name",
					Namespace: "mynamespace",
				},
			},
			routeRule: gatewayapi.TCPRouteRule{
				BackendRefs: []gatewayapi.BackendRef{
					{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Port: lo.ToPtr(gatewayapi.PortNumber(1234)),
						},
					},
					{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Port: lo.ToPtr(gatewayapi.PortNumber(5678)),
						},
					},
				},
			},
			gwPorts: []gatewayapi.PortNumber{8080, 8888},
			expected: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:      "mytcproute-name",
						Namespace: "mynamespace",
					},
					Route: kong.Route{
						Name: lo.ToPtr("tcproute.mynamespace.mytcproute-name.0.0"),
						Destinations: []*kong.CIDRPort{
							{
								Port: lo.ToPtr(8080),
							},
							{
								Port: lo.ToPtr(8888),
							},
						},
						Protocols: []*string{
							lo.ToPtr("tcp"),
						},
						Tags: []*string{
							kong.String("k8s-name:mytcproute-name"),
							kong.String("k8s-namespace:mynamespace"),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kongRoutes, err := generateKongRoutesFromRouteRule(tc.route, tc.gwPorts, 0, tc.routeRule)
			require.NoError(t, err)
			require.NotNil(t, kongRoutes)
			require.Equal(t, tc.expected, kongRoutes)
		})
	}
}

func TestGenerateKongRoutesFromRouteRule_UDP(t *testing.T) {
	testCases := []struct {
		name      string
		route     *gatewayapi.UDPRoute
		gwPorts   []gatewayapi.PortNumber
		routeRule gatewayapi.UDPRouteRule
		expected  []kongstate.Route
	}{
		{
			name: "UDPRoute gets translated correctly to kong.Route",
			route: &gatewayapi.UDPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "myudproute-name",
					Namespace: "mynamespace",
				},
			},
			routeRule: gatewayapi.UDPRouteRule{
				BackendRefs: []gatewayapi.BackendRef{
					{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Port: lo.ToPtr(gatewayapi.PortNumber(1234)),
						},
					},
				},
			},
			gwPorts: []gatewayapi.PortNumber{8080},
			expected: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:      "myudproute-name",
						Namespace: "mynamespace",
					},
					Route: kong.Route{
						Name: lo.ToPtr("udproute.mynamespace.myudproute-name.0.0"),
						Destinations: []*kong.CIDRPort{
							{
								Port: lo.ToPtr(8080),
							},
						},
						Protocols: []*string{
							lo.ToPtr("udp"),
						},
						Tags: []*string{
							kong.String("k8s-name:myudproute-name"),
							kong.String("k8s-namespace:mynamespace"),
						},
					},
				},
			},
		},
		{
			name: "UDPRoute with multiple backends and multiple Gateway ports gets translated correctly to kong.Route",
			route: &gatewayapi.UDPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "myudproute-name",
					Namespace: "mynamespace",
				},
			},
			routeRule: gatewayapi.UDPRouteRule{
				BackendRefs: []gatewayapi.BackendRef{
					{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Port: lo.ToPtr(gatewayapi.PortNumber(1234)),
						},
					},
					{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Port: lo.ToPtr(gatewayapi.PortNumber(5678)),
						},
					},
				},
			},
			gwPorts: []gatewayapi.PortNumber{8080, 8888},
			expected: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:      "myudproute-name",
						Namespace: "mynamespace",
					},
					Route: kong.Route{
						Name: lo.ToPtr("udproute.mynamespace.myudproute-name.0.0"),
						Destinations: []*kong.CIDRPort{
							{
								Port: lo.ToPtr(8080),
							},
							{
								Port: lo.ToPtr(8888),
							},
						},
						Protocols: []*string{
							lo.ToPtr("udp"),
						},
						Tags: []*string{
							kong.String("k8s-name:myudproute-name"),
							kong.String("k8s-namespace:mynamespace"),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kongRoutes, err := generateKongRoutesFromRouteRule(tc.route, tc.gwPorts, 0, tc.routeRule)
			require.NoError(t, err)
			require.NotNil(t, kongRoutes)
			require.Equal(t, tc.expected, kongRoutes)
		})
	}
}

func TestGenerateKongRoutesFromRouteRule_TLS(t *testing.T) {
	testcases := []struct {
		name      string
		route     *gatewayapi.TLSRoute
		routeRule gatewayapi.TLSRouteRule
		expected  []kongstate.Route
	}{
		{
			name: "TLSRoute gets translated correctly to kong.Route",
			route: &gatewayapi.TLSRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mytlsroute-name",
					Namespace: "mynamespace",
				},
				Spec: gatewayapi.TLSRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						"hostname.com",
						"hostname2.com",
					},
				},
			},
			routeRule: gatewayapi.TLSRouteRule{},
			expected: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:      "mytlsroute-name",
						Namespace: "mynamespace",
					},
					Route: kong.Route{
						Name: lo.ToPtr("tlsroute.mynamespace.mytlsroute-name.0.0"),
						SNIs: []*string{
							lo.ToPtr("hostname.com"),
							lo.ToPtr("hostname2.com"),
						},
						Protocols: []*string{
							lo.ToPtr("tls"),
						},
						Tags: []*string{
							kong.String("k8s-name:mytlsroute-name"),
							kong.String("k8s-namespace:mynamespace"),
						},
					},
				},
			},
		},
		{
			name: "TLSRoute without hostnames gets translated correctly to kong.Route without SNIs",
			route: &gatewayapi.TLSRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mytlsroute-name",
					Namespace: "mynamespace",
				},
				Spec: gatewayapi.TLSRouteSpec{},
			},
			routeRule: gatewayapi.TLSRouteRule{},
			expected: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:      "mytlsroute-name",
						Namespace: "mynamespace",
					},
					Route: kong.Route{
						Name: lo.ToPtr("tlsroute.mynamespace.mytlsroute-name.0.0"),
						SNIs: []*string{},
						Protocols: []*string{
							lo.ToPtr("tls"),
						},
						Tags: []*string{
							kong.String("k8s-name:mytlsroute-name"),
							kong.String("k8s-namespace:mynamespace"),
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// TLSRoute matches based on hostname with Gateway listener thus passing gwPorts is pointless.
			kongRoutes, err := generateKongRoutesFromRouteRule(tc.route, nil, 0, tc.routeRule)
			require.NoError(t, err)
			require.NotNil(t, kongRoutes)
			require.Equal(t, tc.expected, kongRoutes)
		})
	}
}
