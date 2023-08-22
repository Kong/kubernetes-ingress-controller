package parser

import (
	"strings"
	"testing"

	"github.com/blang/semver/v4"
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
)

var (
	udpRouteTypeMeta = metav1.TypeMeta{Kind: "UDPRoute", APIVersion: gatewayv1alpha2.SchemeGroupVersion.String()}
)

func TestIngressRulesFromUDPRoutes(t *testing.T) {
	testCases := []struct {
		name                 string
		udpRoutes            []*gatewayv1alpha2.UDPRoute
		expectedKongServices []kongstate.Service
		expectedKongRoutes   map[string][]kongstate.Route
		expectedFailures     []failures.ResourceFailure
	}{
		{
			name: "single UDPRoute with single rule",
			udpRoutes: []*gatewayv1alpha2.UDPRoute{
				{
					TypeMeta:   udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{Name: "single-rule", Namespace: "default"},
					Spec: gatewayv1alpha2.UDPRouteSpec{
						Rules: []gatewayv1alpha2.UDPRouteRule{
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
						Name:     kong.String("udproute.default.single-rule.0"),
						Protocol: kong.String("udp"),
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
				"udproute.default.single-rule.0": {
					{
						Route: kong.Route{
							Name: kong.String("udproute.default.single-rule.0.0"),
							Destinations: []*kong.CIDRPort{
								{Port: kong.Int(80)},
							},
							Protocols: kong.StringSlice("udp"),
						},
					},
				},
			},
		},
		{
			name: "single UDPRoute with multiple rules",
			udpRoutes: []*gatewayv1alpha2.UDPRoute{
				{
					TypeMeta:   udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{Name: "multiple-rules", Namespace: "default"},
					Spec: gatewayv1alpha2.UDPRouteSpec{
						Rules: []gatewayv1alpha2.UDPRouteRule{
							{
								BackendRefs: []gatewayv1alpha2.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
								},
							},
							{
								BackendRefs: []gatewayv1alpha2.BackendRef{
									builder.NewBackendRef("service2").WithPort(81).Build(),
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name:     kong.String("udproute.default.multiple-rules.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service1",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						},
					},
				},
				{
					Service: kong.Service{
						Name:     kong.String("udproute.default.multiple-rules.1"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service2",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(81)},
						},
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"udproute.default.multiple-rules.0": {
					{
						Route: kong.Route{
							Name: kong.String("udproute.default.multiple-rules.0.0"),
							Destinations: []*kong.CIDRPort{
								{Port: kong.Int(80)},
							},
							Protocols: kong.StringSlice("udp"),
						},
					},
				},
				"udproute.default.multiple-rules.1": {
					{
						Route: kong.Route{
							Name: kong.String("udproute.default.multiple-rules.1.0"),
							Destinations: []*kong.CIDRPort{
								{Port: kong.Int(81)},
							},
							Protocols: kong.StringSlice("udp"),
						},
					},
				},
			},
		},
		{
			name: "single UDPRoute with single rule and multiple backendRefs",
			udpRoutes: []*gatewayv1alpha2.UDPRoute{
				{
					TypeMeta:   udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{Name: "multiple-backends", Namespace: "default"},
					Spec: gatewayv1alpha2.UDPRouteSpec{
						Rules: []gatewayv1alpha2.UDPRouteRule{
							{
								BackendRefs: []gatewayv1alpha2.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
									builder.NewBackendRef("service1").WithPort(81).Build(),
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name:     kong.String("udproute.default.multiple-backends.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service1",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						},
						{
							Name:    "service1",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(81)},
						},
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"udproute.default.multiple-backends.0": {
					{
						Route: kong.Route{
							Name: kong.String("udproute.default.multiple-backends.0.0"),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fakestore, err := store.NewFakeStore(store.FakeObjects{
				UDPRoutes: tc.udpRoutes,
			})
			require.NoError(t, err)
			parser := mustNewParser(t, fakestore)
			parser.featureFlags.CombinedServiceRoutes = true

			failureCollector, err := failures.NewResourceFailuresCollector(logrus.New())
			require.NoError(t, err)
			parser.failuresCollector = failureCollector

			result := parser.ingressRulesFromUDPRoutes()
			// check services
			require.Len(t, result.ServiceNameToServices, len(tc.expectedKongServices),
				"should have expected number of services")
			for _, expectedKongService := range tc.expectedKongServices {
				kongService, ok := result.ServiceNameToServices[*expectedKongService.Name]
				require.Truef(t, ok, "should find service %s", *expectedKongService.Name)
				require.Equalf(t, expectedKongService.Backends, kongService.Backends,
					"service %s should have expected backends", *expectedKongService.Name)
				require.Equalf(t, "udp", *kongService.Protocol, "service %s should use UDP protocol", *expectedKongService.Name)
				// check for routes attached to the service
				expectedKongRoutes := tc.expectedKongRoutes[*kongService.Name]
				require.Lenf(t, kongService.Routes, len(expectedKongRoutes),
					"service %s should have expected number of routes", *expectedKongService.Name)

				kongRouteNameToRoute := lo.SliceToMap(kongService.Routes, func(r kongstate.Route) (string, kongstate.Route) {
					return *r.Name, r
				})
				for _, expectedRoute := range expectedKongRoutes {
					routeName := expectedRoute.Name
					r, ok := kongRouteNameToRoute[*routeName]
					require.Truef(t, ok, "should find route %s", *routeName)
					require.Equalf(t, expectedRoute.Destinations, r.Destinations, "route %s should have expected destinations", *routeName)
					require.Equalf(t, expectedRoute.Protocols, r.Protocols, "route %s should have expected protocols", *routeName)
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

func TestIngressRulesFromUDPRoutesUsingExpressionRoutes(t *testing.T) {
	testCases := []struct {
		name                 string
		udpRoutes            []*gatewayv1alpha2.UDPRoute
		kongVersion          semver.Version
		expectedKongServices []kongstate.Service
		expectedKongRoutes   map[string][]kongstate.Route
		expectedFailures     []failures.ResourceFailure
	}{}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fakestore, err := store.NewFakeStore(store.FakeObjects{
				UDPRoutes: tc.udpRoutes,
			})
			require.NoError(t, err)
			parser := mustNewParser(t, fakestore)
			parser.featureFlags.CombinedServiceRoutes = true
			parser.featureFlags.ExpressionRoutes = true
			parser.kongVersion = tc.kongVersion

			failureCollector, err := failures.NewResourceFailuresCollector(logrus.New())
			require.NoError(t, err)
			parser.failuresCollector = failureCollector

			result := parser.ingressRulesFromUDPRoutes()
			// check services
			require.Len(t, result.ServiceNameToServices, len(tc.expectedKongServices),
				"should have expected number of services")
			for _, expectedKongService := range tc.expectedKongServices {
				kongService, ok := result.ServiceNameToServices[*expectedKongService.Name]
				require.Truef(t, ok, "should find service %s", *expectedKongService.Name)
				require.Equalf(t, expectedKongService.Backends, kongService.Backends,
					"service %s should have expected backends", *expectedKongService.Name)
				require.Equalf(t, "udp", kongService.Protocol, "service %s should use UDP protocol", *expectedKongService.Name)
				// check for routes attached to the service
				expectedKongRoutes := tc.expectedKongRoutes[*kongService.Name]
				require.Lenf(t, kongService.Routes, len(expectedKongRoutes),
					"service %s should have expected number of routes", *expectedKongService.Name)

				kongRouteNameToRoute := lo.SliceToMap(kongService.Routes, func(r kongstate.Route) (string, kongstate.Route) {
					return *r.Name, r
				})
				for _, expectedRoute := range expectedKongRoutes {
					routeName := expectedRoute.Name
					r, ok := kongRouteNameToRoute[*routeName]
					require.Truef(t, ok, "should find route %s", *routeName)
					require.Equalf(t, expectedRoute.Expression, r.Expression, "route %s should have expected expression", *routeName)
					require.Equalf(t, expectedRoute.Protocols, r.Protocols, "route %s should have expected protocols", *routeName)
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
