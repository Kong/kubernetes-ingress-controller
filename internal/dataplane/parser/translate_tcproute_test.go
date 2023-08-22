package parser

import (
	"strings"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
)

func TestIngressRulesFromTCPRoutesUsingExpressionRoutes(t *testing.T) {
	tcpRouteTypeMeta := metav1.TypeMeta{Kind: "TCPRoute", APIVersion: gatewayv1alpha2.SchemeGroupVersion.String()}

	testCases := []struct {
		name                 string
		tcpRoutes            []*gatewayv1alpha2.TCPRoute
		expectedKongServices []kongstate.Service
		expectedKongRoutes   map[string][]kongstate.Route
		expectedFailures     []failures.ResourceFailure
	}{
		{
			name: "tcproute with single rule and single backendref",
			tcpRoutes: []*gatewayv1alpha2.TCPRoute{
				{
					TypeMeta: tcpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "tcproute-1",
					},
					Spec: gatewayv1alpha2.TCPRouteSpec{
						Rules: []gatewayv1alpha2.TCPRouteRule{
							{
								BackendRefs: []gatewayv1alpha2.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("tcproute.default.tcproute-1.0"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service1",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						},
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"tcproute.default.tcproute-1.0": {
					{
						Route: kong.Route{
							Name:         kong.String("tcproute.default.tcproute-1.0.0"),
							Expression:   kong.String(`net.dst.port == 80`),
							PreserveHost: kong.Bool(true),
							Protocols:    kong.StringSlice("tcp"),
						},
						ExpressionRoutes: true,
					},
				},
			},
		},
		{
			name: "tcproute with single rule and multiple backendrefs",
			tcpRoutes: []*gatewayv1alpha2.TCPRoute{
				{
					TypeMeta: tcpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "tcproute-1",
					},
					Spec: gatewayv1alpha2.TCPRouteSpec{
						Rules: []gatewayv1alpha2.TCPRouteRule{
							{
								BackendRefs: []gatewayv1alpha2.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
									builder.NewBackendRef("service2").WithPort(443).Build(),
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("tcproute.default.tcproute-1.0"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service1",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						},
						{
							Name:    "service2",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(443)},
						},
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"tcproute.default.tcproute-1.0": {
					{
						Route: kong.Route{
							Name:         kong.String("tcproute.default.tcproute-1.0.0"),
							Expression:   kong.String(`(net.dst.port == 80) || (net.dst.port == 443)`),
							PreserveHost: kong.Bool(true),
							Protocols:    kong.StringSlice("tcp"),
						},
						ExpressionRoutes: true,
					},
				},
			},
		},
		{
			name: "tcproute with multiple rules",
			tcpRoutes: []*gatewayv1alpha2.TCPRoute{
				{
					TypeMeta: tcpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "tcproute-1",
					},
					Spec: gatewayv1alpha2.TCPRouteSpec{
						Rules: []gatewayv1alpha2.TCPRouteRule{
							{
								BackendRefs: []gatewayv1alpha2.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
									builder.NewBackendRef("service2").WithPort(443).Build(),
								},
							},
							{
								BackendRefs: []gatewayv1alpha2.BackendRef{
									builder.NewBackendRef("service3").WithPort(8080).Build(),
									builder.NewBackendRef("service4").WithPort(8443).Build(),
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("tcproute.default.tcproute-1.0"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service1",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						},
						{
							Name:    "service2",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(443)},
						},
					},
				},
				{
					Service: kong.Service{
						Name: kong.String("tcproute.default.tcproute-1.1"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service3",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(8080)},
						},
						{
							Name:    "service4",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(8443)},
						},
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"tcproute.default.tcproute-1.0": {
					{
						Route: kong.Route{
							Name:         kong.String("tcproute.default.tcproute-1.0.0"),
							Expression:   kong.String(`(net.dst.port == 80) || (net.dst.port == 443)`),
							PreserveHost: kong.Bool(true),
							Protocols:    kong.StringSlice("tcp"),
						},
						ExpressionRoutes: true,
					},
				},
				"tcproute.default.tcproute-1.1": {
					{
						Route: kong.Route{
							Name:         kong.String("tcproute.default.tcproute-1.1.0"),
							Expression:   kong.String(`(net.dst.port == 8080) || (net.dst.port == 8443)`),
							PreserveHost: kong.Bool(true),
							Protocols:    kong.StringSlice("tcp"),
						},
						ExpressionRoutes: true,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fakestore, err := store.NewFakeStore(store.FakeObjects{TCPRoutes: tc.tcpRoutes})
			require.NoError(t, err)
			parser := mustNewParser(t, fakestore)
			parser.featureFlags.ExpressionRoutes = true
			parser.kongVersion = versions.ExpressionRouterL4Cutoff

			failureCollector, err := failures.NewResourceFailuresCollector(logrus.New())
			require.NoError(t, err)
			parser.failuresCollector = failureCollector

			result := parser.ingressRulesFromTCPRoutes()
			// check services
			require.Len(t, result.ServiceNameToServices, len(tc.expectedKongServices),
				"should have expected number of services")
			for _, expectedKongService := range tc.expectedKongServices {
				kongService, ok := result.ServiceNameToServices[*expectedKongService.Name]
				require.Truef(t, ok, "should find service %s", expectedKongService.Name)
				require.Equal(t, expectedKongService.Backends, kongService.Backends)
				// check routes
				expectedKongRoutes := tc.expectedKongRoutes[*kongService.Name]
				require.Len(t, kongService.Routes, len(expectedKongRoutes))

				kongRouteNameToRoute := lo.SliceToMap(kongService.Routes, func(r kongstate.Route) (string, kongstate.Route) {
					return *r.Name, r
				})
				for _, expectedRoute := range expectedKongRoutes {
					routeName := expectedRoute.Name
					r, ok := kongRouteNameToRoute[*routeName]
					require.Truef(t, ok, "should find route %s", *routeName)
					require.Equal(t, expectedRoute.Expression, r.Expression)
					require.Equal(t, expectedRoute.Protocols, r.Protocols)
				}
			}
			// check translation failures
			translationFailures := failureCollector.PopResourceFailures()
			require.Len(t, translationFailures, len(tc.expectedFailures))
			for _, expectedTranslationFailure := range tc.expectedFailures {
				expectedFailureMessage := expectedTranslationFailure.Message()
				require.True(t, lo.ContainsBy(translationFailures, func(failure failures.ResourceFailure) bool {
					return strings.Contains(failure.Message(), expectedFailureMessage)
				}))
			}
		})
	}
}
