package translator

import (
	"strings"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

func TestIngressRulesFromTLSRoutesUsingExpressionRoutes(t *testing.T) {
	tlsRouteTypeMeta := metav1.TypeMeta{Kind: "TLSRoute", APIVersion: corev1.SchemeGroupVersion.String()}

	testCases := []struct {
		name                 string
		tcpRoutes            []*gatewayapi.TLSRoute
		services             []*corev1.Service
		expectedKongServices []kongstate.Service
		expectedKongRoutes   map[string][]kongstate.Route
		expectedFailures     []failures.ResourceFailure
	}{
		{
			name: "tlsroute with single rule",
			tcpRoutes: []*gatewayapi.TLSRoute{
				{
					TypeMeta: tlsRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "tlsroute-1",
					},
					Spec: gatewayapi.TLSRouteSpec{
						Hostnames: []gatewayapi.Hostname{
							"foo.com",
							"bar.com",
						},
						Rules: []gatewayapi.TLSRouteRule{
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
			// After https://github.com/Kong/kubernetes-ingress-controller/pull/5392
			// is merged the backendRef will be checked for existence in the store
			// so we need to add them here.
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "service1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "service2",
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("tlsroute.default.tlsroute-1.0"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(), builder.NewKongstateServiceBackend("service2").WithPortNumber(443).MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"tlsroute.default.tlsroute-1.0": {
					{
						Route: kong.Route{
							Name:         kong.String("tlsroute.default.tlsroute-1.0.0"),
							Expression:   kong.String(`(tls.sni == "foo.com") || (tls.sni == "bar.com")`),
							PreserveHost: kong.Bool(true),
							Protocols:    kong.StringSlice("tls"),
						},
						ExpressionRoutes: true,
					},
				},
			},
		},
		{
			name: "tlsroute with multiple rules",
			tcpRoutes: []*gatewayapi.TLSRoute{
				{
					TypeMeta: tlsRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "tlsroute-1",
					},
					Spec: gatewayapi.TLSRouteSpec{
						Hostnames: []gatewayapi.Hostname{
							"foo.com",
							"bar.com",
						},
						Rules: []gatewayapi.TLSRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
									builder.NewBackendRef("service2").WithPort(443).Build(),
								},
							},
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
			// After https://github.com/Kong/kubernetes-ingress-controller/pull/5392
			// is merged the backendRef will be checked for existence in the store
			// so we need to add them here.
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "service1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "service2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "service3",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "service4",
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("tlsroute.default.tlsroute-1.0"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(), builder.NewKongstateServiceBackend("service2").WithPortNumber(443).MustBuild(),
					},
				},
				{
					Service: kong.Service{
						Name: kong.String("tlsroute.default.tlsroute-1.1"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service3").WithPortNumber(8080).MustBuild(),
						builder.NewKongstateServiceBackend("service4").WithPortNumber(8443).MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"tlsroute.default.tlsroute-1.0": {
					{
						Route: kong.Route{
							Name:         kong.String("tlsroute.default.tlsroute-1.0.0"),
							Expression:   kong.String(`(tls.sni == "foo.com") || (tls.sni == "bar.com")`),
							PreserveHost: kong.Bool(true),
							Protocols:    kong.StringSlice("tls"),
						},
						ExpressionRoutes: true,
					},
				},
				"tlsroute.default.tlsroute-1.1": {
					{
						Route: kong.Route{
							Name:         kong.String("tlsroute.default.tlsroute-1.1.0"),
							Expression:   kong.String(`(tls.sni == "foo.com") || (tls.sni == "bar.com")`),
							PreserveHost: kong.Bool(true),
							Protocols:    kong.StringSlice("tls"),
						},
						ExpressionRoutes: true,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakestore, err := store.NewFakeStore(store.FakeObjects{
				TLSRoutes: tc.tcpRoutes,
				Services:  tc.services,
			})
			require.NoError(t, err)
			translator := mustNewTranslator(t, fakestore)
			translator.featureFlags.ExpressionRoutes = true

			failureCollector := failures.NewResourceFailuresCollector(zapr.NewLogger(zap.NewNop()))
			translator.failuresCollector = failureCollector

			result := translator.ingressRulesFromTLSRoutes()
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
