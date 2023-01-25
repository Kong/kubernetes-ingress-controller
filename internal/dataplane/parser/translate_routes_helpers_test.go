package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func TestGenerateKongRoutesFromRouteRule_TCP(t *testing.T) {
	testcases := []struct {
		name      string
		route     *gatewayv1alpha2.TCPRoute
		routeRule gatewayv1alpha2.TCPRouteRule
		expected  []kongstate.Route
	}{
		{
			name: "TCPRoute gets translated correctly to kong.Route",
			route: &gatewayv1alpha2.TCPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mytcproute-name",
					Namespace: "mynamespace",
				},
			},
			routeRule: gatewayv1alpha2.TCPRouteRule{
				BackendRefs: []gatewayv1alpha2.BackendRef{
					{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Port: lo.ToPtr(gatewayv1alpha2.PortNumber(1234)),
						},
					},
				},
			},
			expected: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:        "mytcproute-name",
						Namespace:   "mynamespace",
						Annotations: map[string]string{},
					},
					Route: kong.Route{
						Name: lo.ToPtr("tcproute.mynamespace.mytcproute-name.0.0"),
						Destinations: []*kong.CIDRPort{
							{
								Port: lo.ToPtr(1234),
							},
						},
						Protocols: []*string{
							lo.ToPtr("tcp"),
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			kongRoutes, err := generateKongRoutesFromRouteRule(tc.route, 0, tc.routeRule)
			require.NoError(t, err)
			require.NotNil(t, kongRoutes)
			if !cmp.Equal(tc.expected, kongRoutes) {
				t.Logf("actual []kongstate.Route differs from expected\n%s", cmp.Diff(tc.expected, kongRoutes))
				t.Fail()
			}
		})
	}
}

func TestGenerateKongRoutesFromRouteRule_UDP(t *testing.T) {
	testcases := []struct {
		name      string
		route     *gatewayv1alpha2.UDPRoute
		routeRule gatewayv1alpha2.UDPRouteRule
		expected  []kongstate.Route
	}{
		{
			name: "UDPRoute gets translated correctly to kong.Route",
			route: &gatewayv1alpha2.UDPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "myudproute-name",
					Namespace: "mynamespace",
				},
			},
			routeRule: gatewayv1alpha2.UDPRouteRule{
				BackendRefs: []gatewayv1alpha2.BackendRef{
					{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Port: lo.ToPtr(gatewayv1alpha2.PortNumber(1234)),
						},
					},
				},
			},
			expected: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:        "myudproute-name",
						Namespace:   "mynamespace",
						Annotations: map[string]string{},
					},
					Route: kong.Route{
						Name: lo.ToPtr("udproute.mynamespace.myudproute-name.0.0"),
						Destinations: []*kong.CIDRPort{
							{
								Port: lo.ToPtr(1234),
							},
						},
						Protocols: []*string{
							lo.ToPtr("udp"),
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			kongRoutes, err := generateKongRoutesFromRouteRule(tc.route, 0, tc.routeRule)
			require.NoError(t, err)
			require.NotNil(t, kongRoutes)
			if !cmp.Equal(tc.expected, kongRoutes) {
				t.Logf("actual []kongstate.Route differs from expected\n%s", cmp.Diff(tc.expected, kongRoutes))
				t.Fail()
			}
		})
	}
}

func TestGenerateKongRoutesFromRouteRule_TLS(t *testing.T) {
	testcases := []struct {
		name      string
		route     *gatewayv1alpha2.TLSRoute
		routeRule gatewayv1alpha2.TLSRouteRule
		expected  []kongstate.Route
	}{
		{
			name: "TLSRoute gets translated correctly to kong.Route",
			route: &gatewayv1alpha2.TLSRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mytlsroute-name",
					Namespace: "mynamespace",
				},
				Spec: gatewayv1alpha2.TLSRouteSpec{
					Hostnames: []gatewayv1alpha2.Hostname{
						"hostname.com",
						"hostname2.com",
					},
				},
			},
			routeRule: gatewayv1alpha2.TLSRouteRule{},
			expected: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:        "mytlsroute-name",
						Namespace:   "mynamespace",
						Annotations: map[string]string{},
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
					},
				},
			},
		},
		{
			name: "TLSRoute without hostnames gets translated correctly to kong.Route without SNIs",
			route: &gatewayv1alpha2.TLSRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mytlsroute-name",
					Namespace: "mynamespace",
				},
				Spec: gatewayv1alpha2.TLSRouteSpec{},
			},
			routeRule: gatewayv1alpha2.TLSRouteRule{},
			expected: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:        "mytlsroute-name",
						Namespace:   "mynamespace",
						Annotations: map[string]string{},
					},
					Route: kong.Route{
						Name: lo.ToPtr("tlsroute.mynamespace.mytlsroute-name.0.0"),
						SNIs: []*string{},
						Protocols: []*string{
							lo.ToPtr("tls"),
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			kongRoutes, err := generateKongRoutesFromRouteRule(tc.route, 0, tc.routeRule)
			require.NoError(t, err)
			require.NotNil(t, kongRoutes)
			if !cmp.Equal(tc.expected, kongRoutes) {
				t.Logf("actual []kongstate.Route differs from expected\n%s", cmp.Diff(tc.expected, kongRoutes))
				t.Fail()
			}
		})
	}
}
