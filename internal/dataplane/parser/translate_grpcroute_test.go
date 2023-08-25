package parser

import (
	"strconv"
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
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
)

func TestIngressRulesFromGRPCRoutesUsingExpressionRoutes(t *testing.T) {
	fakestore, err := store.NewFakeStore(store.FakeObjects{})
	require.NoError(t, err)
	parser := mustNewParser(t, fakestore)
	parser.featureFlags.CombinedServiceRoutes = true
	parser.featureFlags.ExpressionRoutes = true
	grpcRouteTypeMeta := metav1.TypeMeta{Kind: "GRPCRoute", APIVersion: gatewayv1alpha2.SchemeGroupVersion.String()}

	testCases := []struct {
		name                 string
		grpcRoutes           []*gatewayv1alpha2.GRPCRoute
		expectedKongServices []kongstate.Service
		// service name -> routes
		expectedKongRoutes map[string][]kongstate.Route
		expectedFailures   []failures.ResourceFailure
	}{
		{
			name: "single GRPCRoute with multiple hostnames and multiple rules",
			grpcRoutes: []*gatewayv1alpha2.GRPCRoute{
				{
					TypeMeta: grpcRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grpcroute-1",
					},
					Spec: gatewayv1alpha2.GRPCRouteSpec{
						Hostnames: []gatewayv1alpha2.Hostname{
							"foo.com",
							"*.bar.com",
						},
						Rules: []gatewayv1alpha2.GRPCRouteRule{
							{
								Matches: []gatewayv1alpha2.GRPCRouteMatch{
									{
										Method: &gatewayv1alpha2.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
											Service: lo.ToPtr("v1"),
											Method:  lo.ToPtr("foo"),
										},
									},
								},
								BackendRefs: []gatewayv1alpha2.GRPCBackendRef{
									{
										BackendRef: builder.NewBackendRef("service1").WithPort(80).Build(),
									},
								},
							},
							{
								Matches: []gatewayv1alpha2.GRPCRouteMatch{
									{
										Method: &gatewayv1alpha2.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
											Service: lo.ToPtr("v1"),
											Method:  lo.ToPtr("foobar"),
										},
									},
								},
								BackendRefs: []gatewayv1alpha2.GRPCBackendRef{
									{
										BackendRef: builder.NewBackendRef("service2").WithPort(80).Build(),
									},
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("grpcroute.default.grpcroute-1.foo.com.0"),
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
						Name: kong.String("grpcroute.default.grpcroute-1.foo.com.1"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service2",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						},
					},
				},
				{
					Service: kong.Service{
						Name: kong.String("grpcroute.default.grpcroute-1._.bar.com.0"),
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
						Name: kong.String("grpcroute.default.grpcroute-1._.bar.com.1"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service2",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						},
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
			grpcRoutes: []*gatewayv1alpha2.GRPCRoute{
				{
					TypeMeta: grpcRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grpcroute-1",
					},
					Spec: gatewayv1alpha2.GRPCRouteSpec{
						Hostnames: []gatewayv1alpha2.Hostname{
							"foo.com",
						},
						Rules: []gatewayv1alpha2.GRPCRouteRule{
							{
								Matches: []gatewayv1alpha2.GRPCRouteMatch{
									{
										Method: &gatewayv1alpha2.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
											Service: lo.ToPtr("v1"),
											Method:  lo.ToPtr("foo"),
										},
									},
									{
										Method: &gatewayv1alpha2.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
											Service: lo.ToPtr("v1"),
											Method:  lo.ToPtr("foobar"),
										},
									},
								},
								BackendRefs: []gatewayv1alpha2.GRPCBackendRef{
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
					Spec: gatewayv1alpha2.GRPCRouteSpec{
						Rules: []gatewayv1alpha2.GRPCRouteRule{
							{
								Matches: []gatewayv1alpha2.GRPCRouteMatch{
									{
										Method: &gatewayv1alpha2.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
											Service: lo.ToPtr("v2"),
											Method:  lo.ToPtr("foo"),
										},
									},
								},
								BackendRefs: []gatewayv1alpha2.GRPCBackendRef{
									{
										BackendRef: builder.NewBackendRef("service2").WithPort(80).Build(),
									},
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("grpcroute.default.grpcroute-1.foo.com.0"),
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
						Name: kong.String("grpcroute.default.grpcroute-2._.0"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service2",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						},
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
			grpcRoutes: []*gatewayv1alpha2.GRPCRoute{
				{
					TypeMeta: grpcRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grpcroute-1",
					},
					Spec: gatewayv1alpha2.GRPCRouteSpec{
						Rules: []gatewayv1alpha2.GRPCRouteRule{
							{
								Matches: []gatewayv1alpha2.GRPCRouteMatch{
									{
										Method: &gatewayv1alpha2.GRPCMethodMatch{
											Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
											Service: lo.ToPtr("v1"),
											Method:  lo.ToPtr("foo"),
										},
									},
								},
								BackendRefs: []gatewayv1alpha2.GRPCBackendRef{
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
					Spec: gatewayv1alpha2.GRPCRouteSpec{},
				},
				{
					TypeMeta: grpcRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grpcroute-no-hostnames-no-matches",
					},
					Spec: gatewayv1alpha2.GRPCRouteSpec{
						Rules: []gatewayv1alpha2.GRPCRouteRule{
							{
								BackendRefs: []gatewayv1alpha2.GRPCBackendRef{
									{
										BackendRef: builder.NewBackendRef("service0").WithPort(80).Build(),
									},
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("grpcroute.default.grpcroute-1._.0"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service2",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						},
					},
				},
				{
					Service: kong.Service{
						Name: kong.String("grpcroute.default.grpcroute-no-hostnames-no-matches._.0"),
					},
					Backends: []kongstate.ServiceBackend{
						{
							Name:    "service0",
							PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						},
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
							Expression: kong.String(translators.CatchAllHTTPExpression),
						},
					},
				},
			},
			expectedFailures: []failures.ResourceFailure{
				newResourceFailure(t, translators.ErrRouteValidationNoRules.Error(),
					&gatewayv1alpha2.GRPCRoute{
						TypeMeta: grpcRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "grpcroute-no-rules",
						},
					}),
			},
		},
	}

	for i, tc := range testCases {
		indexStr := strconv.Itoa(i)
		tc := tc
		t.Run(indexStr+"-"+tc.name, func(t *testing.T) {
			failureCollector, err := failures.NewResourceFailuresCollector(logrus.New())
			require.NoError(t, err)
			parser.failuresCollector = failureCollector

			result := newIngressRules()
			parser.ingressRulesFromGRPCRoutesUsingExpressionRoutes(tc.grpcRoutes, &result)
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
			// check translation failures
			translationFailures := failureCollector.PopResourceFailures()
			require.Equal(t, len(tc.expectedFailures), len(translationFailures))
			for _, expectedTranslationFailure := range tc.expectedFailures {
				expectedFailureMessage := expectedTranslationFailure.Message()
				require.True(t, lo.ContainsBy(translationFailures, func(failure failures.ResourceFailure) bool {
					return strings.Contains(failure.Message(), expectedFailureMessage)
				}))
			}
		})

	}
}
