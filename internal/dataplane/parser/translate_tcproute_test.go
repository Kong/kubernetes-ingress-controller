package parser

import (
	"strings"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
)

func TestIngressRulesFromTCPRoutesUsingExpressionRoutes(t *testing.T) {
	tcpRouteTypeMeta := metav1.TypeMeta{Kind: "TCPRoute", APIVersion: gatewayv1alpha2.SchemeGroupVersion.String()}

	testCases := []struct {
		name                 string
		gateways             []*gatewayapi.Gateway
		tcpRoutes            []*gatewayapi.TCPRoute
		expectedKongServices []kongstate.Service
		expectedKongRoutes   map[string][]kongstate.Route
		expectedFailures     []failures.ResourceFailure
	}{
		{
			name: "TCPRoute with single rule, single backendref and Gateway in the same namespace",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "default",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
							builder.NewListener("tcp80").WithPort(80).TCP().Build(),
							builder.NewListener("udp80").WithPort(80).UDP().Build(),
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "test-1",
					},
				},
			},
			tcpRoutes: []*gatewayapi.TCPRoute{
				{
					TypeMeta: tcpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tcproute-1",
						Namespace: "default",
					},
					Spec: gatewayapi.TCPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name: "gateway-1",
								},
							},
						},
						Rules: []gatewayapi.TCPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(8080).Build(),
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
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(8080)},
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
			name: "TCPRoute with single rule, multiple backendrefs and Gateway in the different namespace",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "test-1",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
							builder.NewListener("tcp80").WithPort(80).TCP().Build(),
							builder.NewListener("tcp443").WithPort(443).TCP().Build(),
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "default",
					},
				},
			},
			tcpRoutes: []*gatewayapi.TCPRoute{
				{
					TypeMeta: tcpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tcproute-1",
						Namespace: "default",
					},
					Spec: gatewayapi.TCPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:      "gateway-1",
									Namespace: lo.ToPtr(gatewayapi.Namespace("test-1")),
								},
							},
						},
						Rules: []gatewayapi.TCPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
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
			name: "TCPRoute with multiple rules and multiple Gateways with multiple listeners and sectionName for Gateway specified",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "default",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
							{
								Name:     "tcp80",
								Port:     80,
								Protocol: gatewayapi.TCPProtocolType,
							},
							{
								Name:     "tcp443",
								Port:     443,
								Protocol: gatewayapi.TCPProtocolType,
							},
							{
								Name:     "tcp8080",
								Port:     8080,
								Protocol: gatewayapi.TCPProtocolType,
							},
							{
								Name:     "tcp8443",
								Port:     8443,
								Protocol: gatewayapi.TCPProtocolType,
							},
						},
					},
				},
			},
			tcpRoutes: []*gatewayapi.TCPRoute{
				{
					TypeMeta: tcpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tcproute-1",
						Namespace: "default",
					},
					Spec: gatewayapi.TCPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:        "gateway-1",
									SectionName: lo.ToPtr(gatewayapi.SectionName("tcp80")),
								},
								{
									Name:        "gateway-1",
									SectionName: lo.ToPtr(gatewayapi.SectionName("tcp443")),
								},
							},
						},

						Rules: []gatewayapi.TCPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
									builder.NewBackendRef("service2").WithPort(443).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta: tcpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tcproute-2",
						Namespace: "default",
					},
					Spec: gatewayapi.TCPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:        "gateway-1",
									SectionName: lo.ToPtr(gatewayapi.SectionName("tcp8080")),
								},
								{
									Name:        "gateway-1",
									SectionName: lo.ToPtr(gatewayapi.SectionName("tcp8443")),
								},
							},
						},

						Rules: []gatewayapi.TCPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
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
						Name: kong.String("tcproute.default.tcproute-2.0"),
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
				"tcproute.default.tcproute-2.0": {
					{
						Route: kong.Route{
							Name:         kong.String("tcproute.default.tcproute-2.0.0"),
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
			fakeStore, err := store.NewFakeStore(store.FakeObjects{
				Gateways:  tc.gateways,
				TCPRoutes: tc.tcpRoutes,
			})
			require.NoError(t, err)
			parser := mustNewParser(t, fakeStore)
			parser.featureFlags.ExpressionRoutes = true

			failureCollector := failures.NewResourceFailuresCollector(zapr.NewLogger(zap.NewNop()))
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
					require.Equalf(t, expectedRoute.Expression, r.Expression, "route %s should have expected expression", *routeName)
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
