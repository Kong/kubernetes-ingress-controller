package translator

import (
	"strconv"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

func TestIngressRulesFromGRPCRoutesUsingExpressionRoutes(t *testing.T) {
	grpcRouteTypeMeta := metav1.TypeMeta{Kind: "GRPCRoute", APIVersion: gatewayv1.SchemeGroupVersion.String()}

	testCases := []struct {
		name                 string
		grpcRoutes           []*gatewayapi.GRPCRoute
		expectedKongServices []kongstate.Service
		services             []*corev1.Service
		// service name -> routes
		expectedKongRoutes map[string][]kongstate.Route
	}{
		{
			name: "single GRPCRoute with multiple hostnames and multiple rules",
			grpcRoutes: []*gatewayapi.GRPCRoute{
				{
					TypeMeta: grpcRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grpcroute-1",
					},
					Spec: gatewayapi.GRPCRouteSpec{
						Hostnames: []gatewayapi.Hostname{
							"foo.com",
							"*.bar.com",
						},
						Rules: []gatewayapi.GRPCRouteRule{
							{
								Matches: []gatewayapi.GRPCRouteMatch{
									{
										Method: &gatewayapi.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
											Service: lo.ToPtr("v1"),
											Method:  lo.ToPtr("foo"),
										},
									},
								},
								BackendRefs: []gatewayapi.GRPCBackendRef{
									{
										BackendRef: builder.NewBackendRef("service1").WithPort(80).Build(),
									},
								},
							},
							{
								Matches: []gatewayapi.GRPCRouteMatch{
									{
										Method: &gatewayapi.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
											Service: lo.ToPtr("v1"),
											Method:  lo.ToPtr("foobar"),
										},
									},
								},
								BackendRefs: []gatewayapi.GRPCBackendRef{
									{
										BackendRef: builder.NewBackendRef("service2").WithPort(80).Build(),
									},
								},
							},
						},
					},
				},
			},
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
						Name: kong.String("grpcroute.default.grpcroute-1.foo.com.0"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
				{
					Service: kong.Service{
						Name: kong.String("grpcroute.default.grpcroute-1.foo.com.1"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service2").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
				{
					Service: kong.Service{
						Name: kong.String("grpcroute.default.grpcroute-1._.bar.com.0"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
				{
					Service: kong.Service{
						Name: kong.String("grpcroute.default.grpcroute-1._.bar.com.1"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service2").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"grpcroute.default.grpcroute-1.foo.com.0": {
					{
						Route: kong.Route{
							Name:       kong.String("grpcroute.default.grpcroute-1.foo.com.0.0"),
							Expression: kong.String(`(http.path == "/v1/foo") && (http.host == "foo.com")`),
						},
					},
				},
				"grpcroute.default.grpcroute-1.foo.com.1": {
					{
						Route: kong.Route{
							Name:       kong.String("grpcroute.default.grpcroute-1.foo.com.1.0"),
							Expression: kong.String(`(http.path == "/v1/foobar") && (http.host == "foo.com")`),
						},
					},
				},
				"grpcroute.default.grpcroute-1._.bar.com.0": {
					{
						Route: kong.Route{
							Name:       kong.String("grpcroute.default.grpcroute-1._.bar.com.0.0"),
							Expression: kong.String(`(http.path == "/v1/foo") && (http.host =^ ".bar.com")`),
						},
					},
				},
				"grpcroute.default.grpcroute-1._.bar.com.1": {
					{
						Route: kong.Route{
							Name:       kong.String("grpcroute.default.grpcroute-1._.bar.com.1.0"),
							Expression: kong.String(`(http.path == "/v1/foobar") && (http.host =^ ".bar.com")`),
						},
					},
				},
			},
		},
		{
			name: "multiple GRPCRoutes with multiple matches",
			grpcRoutes: []*gatewayapi.GRPCRoute{
				{
					TypeMeta: grpcRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grpcroute-1",
					},
					Spec: gatewayapi.GRPCRouteSpec{
						Hostnames: []gatewayapi.Hostname{
							"foo.com",
						},
						Rules: []gatewayapi.GRPCRouteRule{
							{
								Matches: []gatewayapi.GRPCRouteMatch{
									{
										Method: &gatewayapi.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
											Service: lo.ToPtr("v1"),
											Method:  lo.ToPtr("foo"),
										},
									},
									{
										Method: &gatewayapi.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
											Service: lo.ToPtr("v1"),
											Method:  lo.ToPtr("foobar"),
										},
									},
								},
								BackendRefs: []gatewayapi.GRPCBackendRef{
									{
										BackendRef: builder.NewBackendRef("service1").WithPort(80).Build(),
									},
								},
							},
						},
					},
				},
				{
					TypeMeta: grpcRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grpcroute-2",
					},
					Spec: gatewayapi.GRPCRouteSpec{
						Rules: []gatewayapi.GRPCRouteRule{
							{
								Matches: []gatewayapi.GRPCRouteMatch{
									{
										Method: &gatewayapi.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
											Service: lo.ToPtr("v2"),
											Method:  lo.ToPtr("foo"),
										},
									},
								},
								BackendRefs: []gatewayapi.GRPCBackendRef{
									{
										BackendRef: builder.NewBackendRef("service2").WithPort(80).Build(),
									},
								},
							},
						},
					},
				},
			},
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
						Name: kong.String("grpcroute.default.grpcroute-1.foo.com.0"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
				{
					Service: kong.Service{
						Name: kong.String("grpcroute.default.grpcroute-2._.0"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service2").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"grpcroute.default.grpcroute-1.foo.com.0": {
					{
						Route: kong.Route{
							Name:       kong.String("grpcroute.default.grpcroute-1.foo.com.0.0"),
							Expression: kong.String(`(http.path == "/v1/foo") && (http.host == "foo.com")`),
						},
					},
					{
						Route: kong.Route{
							Name:       kong.String("grpcroute.default.grpcroute-1.foo.com.0.1"),
							Expression: kong.String(`(http.path == "/v1/foobar") && (http.host == "foo.com")`),
						},
					},
				},
				"grpcroute.default.grpcroute-2._.0": {
					{
						Route: kong.Route{
							Name:       kong.String("grpcroute.default.grpcroute-2._.0.0"),
							Expression: kong.String(`http.path == "/v2/foo"`),
						},
					},
				},
			},
		},
		{
			name: "multiple GRPCRoutes with translation error",
			grpcRoutes: []*gatewayapi.GRPCRoute{
				{
					TypeMeta: grpcRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grpcroute-1",
					},
					Spec: gatewayapi.GRPCRouteSpec{
						Rules: []gatewayapi.GRPCRouteRule{
							{
								Matches: []gatewayapi.GRPCRouteMatch{
									{
										Method: &gatewayapi.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
											Service: lo.ToPtr("v1"),
											Method:  lo.ToPtr("foo"),
										},
									},
								},
								BackendRefs: []gatewayapi.GRPCBackendRef{
									{
										BackendRef: builder.NewBackendRef("service2").WithPort(80).Build(),
									},
								},
							},
						},
					},
				},
				{
					TypeMeta: grpcRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grpcroute-no-rules",
					},
					Spec: gatewayapi.GRPCRouteSpec{},
				},
				{
					TypeMeta: grpcRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grpcroute-no-hostnames-no-matches",
					},
					Spec: gatewayapi.GRPCRouteSpec{
						Rules: []gatewayapi.GRPCRouteRule{
							{
								BackendRefs: []gatewayapi.GRPCBackendRef{
									{
										BackendRef: builder.NewBackendRef("service0").WithPort(80).Build(),
									},
								},
							},
						},
					},
				},
			},
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "service0",
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
						Name: kong.String("grpcroute.default.grpcroute-1._.0"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service2").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
				{
					Service: kong.Service{
						Name: kong.String("grpcroute.default.grpcroute-no-hostnames-no-matches._.0"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service0").WithPortNumber(80).MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"grpcroute.default.grpcroute-1._.0": {
					{
						Route: kong.Route{
							Name:       kong.String("grpcroute.default.grpcroute-1._.0.0"),
							Expression: kong.String(`http.path == "/v1/foo"`),
						},
					},
				},
				"grpcroute.default.grpcroute-no-hostnames-no-matches._.0": {
					{
						Route: kong.Route{
							Name:       kong.String("grpcroute.default.grpcroute-no-hostnames-no-matches._.0.0"),
							Expression: kong.String(subtranslator.CatchAllHTTPExpression),
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		indexStr := strconv.Itoa(i)
		tc := tc
		t.Run(indexStr+"-"+tc.name, func(t *testing.T) {
			failureCollector := failures.NewResourceFailuresCollector(zapr.NewLogger(zap.NewNop()))

			fakestore, err := store.NewFakeStore(store.FakeObjects{
				GRPCRoutes: tc.grpcRoutes,
				Services:   tc.services,
			})
			require.NoError(t, err)
			translator := mustNewTranslator(t, fakestore)
			translator.featureFlags.ExpressionRoutes = true
			translator.failuresCollector = failureCollector

			result := newIngressRules()
			translator.ingressRulesFromGRPCRoutesUsingExpressionRoutes(tc.grpcRoutes, &result)
			// check services
			require.Equal(t, len(tc.expectedKongServices), len(result.ServiceNameToServices),
				"should have expected number of services")
			for _, expectedKongService := range tc.expectedKongServices {
				kongService, ok := result.ServiceNameToServices[*expectedKongService.Name]
				require.Truef(t, ok, "should find service %s", expectedKongService.Name)
				require.Equal(t, expectedKongService.Backends, kongService.Backends)
				// check routes
				expectedKongRoutes := tc.expectedKongRoutes[*kongService.Name]
				require.Equal(t, len(expectedKongRoutes), len(kongService.Routes))

				kongRouteNameToRoute := lo.SliceToMap(kongService.Routes, func(r kongstate.Route) (string, kongstate.Route) {
					return *r.Name, r
				})
				for _, expectedRoute := range expectedKongRoutes {
					routeName := expectedRoute.Name
					r, ok := kongRouteNameToRoute[*routeName]
					require.Truef(t, ok, "should find route %s", *routeName)
					require.Equal(t, *expectedRoute.Expression, *r.Expression)
				}
			}
			// Check translation failures.
			translationFailures := failureCollector.PopResourceFailures()
			require.Empty(t, translationFailures)
		})

	}
}
