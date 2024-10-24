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
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

var udpRouteTypeMeta = metav1.TypeMeta{Kind: "UDPRoute", APIVersion: gatewayv1alpha2.SchemeGroupVersion.String()}

func TestIngressRulesFromUDPRoutes(t *testing.T) {
	testCases := []struct {
		name                 string
		gateways             []*gatewayapi.Gateway
		udpRoutes            []*gatewayapi.UDPRoute
		services             []*corev1.Service
		expectedKongServices []kongstate.Service
		expectedKongRoutes   map[string][]kongstate.Route
		expectedFailures     []failures.ResourceFailure
	}{
		{
			name: "single UDPRoute with single rule, single backendref and Gateway in the same namespace",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "default",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
							builder.NewListener("udp80").WithPort(80).UDP().Build(),
							builder.NewListener("tcp80").WithPort(80).TCP().Build(),
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
			udpRoutes: []*gatewayapi.UDPRoute{
				{
					TypeMeta: udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "single-rule",
						Namespace: "default",
					},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name: "gateway-1",
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
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
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name:     kong.String("udproute.default.single-rule.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
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
			name: "multiple UDPRoute with single rule, different SectionName and Gateway in the same namespace",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "default",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
							builder.NewListener("udp80").WithPort(80).UDP().Build(),
							builder.NewListener("udp81").WithPort(81).UDP().Build(),
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
			udpRoutes: []*gatewayapi.UDPRoute{
				{
					TypeMeta: udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rule-1",
						Namespace: "default",
					},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:        "gateway-1",
									SectionName: lo.ToPtr(gatewayapi.SectionName("udp80")),
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta: udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rule-2",
						Namespace: "default",
					},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:        "gateway-1",
									SectionName: lo.ToPtr(gatewayapi.SectionName("udp81")),
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service2").WithPort(81).Build(),
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
						Name:     kong.String("udproute.default.rule-1.0"),
						Protocol: kong.String("udp"),
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
						Name:     kong.String("udproute.default.rule-2.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service2").WithPortNumber(81).MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"udproute.default.rule-1.0": {
					{
						Route: kong.Route{
							Name: kong.String("udproute.default.rule-1.0.0"),
							Destinations: []*kong.CIDRPort{
								{Port: kong.Int(80)},
							},
							Protocols: kong.StringSlice("udp"),
						},
					},
				},
				"udproute.default.rule-2.0": {
					{
						Route: kong.Route{
							Name: kong.String("udproute.default.rule-2.0.0"),
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
			name: "single UDPRoute with single rule and multiple backendRefs and Gateway in a different namespace",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "test-1",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
							builder.NewListener("udp80").WithPort(80).UDP().Build(),
							builder.NewListener("udp81").WithPort(81).UDP().Build(),
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
			udpRoutes: []*gatewayapi.UDPRoute{
				{
					TypeMeta: udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "multiple-backends",
						Namespace: "default",
					},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:      "gateway-1",
									Namespace: lo.ToPtr(gatewayapi.Namespace("test-1")),
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
									builder.NewBackendRef("service1").WithPort(81).Build(),
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
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name:     kong.String("udproute.default.multiple-backends.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(81).
							MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"udproute.default.multiple-backends.0": {
					{
						Route: kong.Route{
							Name: kong.String("udproute.default.multiple-backends.0.0"),
							Destinations: []*kong.CIDRPort{
								{Port: kong.Int(80)},
								{Port: kong.Int(81)},
							},
							Protocols: kong.StringSlice("udp"),
						},
					},
				},
			},
		},
		{
			name: "multiple UDPRoutes with translation errors",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "default",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
							builder.NewListener("udp80").WithPort(80).UDP().Build(),
							builder.NewListener("udp8080").WithPort(8080).UDP().Build(),
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
			udpRoutes: []*gatewayapi.UDPRoute{
				{
					TypeMeta:   udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{Name: "single-rule", Namespace: "default"},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:        "gateway-1",
									SectionName: lo.ToPtr(gatewayapi.SectionName("udp80")),
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta:   udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{Name: "single-rule-2", Namespace: "default"},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:        "gateway-1",
									SectionName: lo.ToPtr(gatewayapi.SectionName("udp8080")),
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service2").WithPort(8080).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta:   udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{Name: "no-rule", Namespace: "default"},
					Spec:       gatewayapi.UDPRouteSpec{},
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
						Name:     kong.String("udproute.default.single-rule.0"),
						Protocol: kong.String("udp"),
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
						Name:     kong.String("udproute.default.single-rule-2.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service2").WithPortNumber(8080).MustBuild(),
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
				"udproute.default.single-rule-2.0": {
					{
						Route: kong.Route{
							Name: kong.String("udproute.default.single-rule-2.0.0"),
							Destinations: []*kong.CIDRPort{
								{Port: kong.Int(8080)},
							},
							Protocols: kong.StringSlice("udp"),
						},
					},
				},
			},
			expectedFailures: []failures.ResourceFailure{
				newResourceFailure(
					t, subtranslator.ErrRouteValidationNoRules.Error(),
					&gatewayapi.UDPRoute{
						TypeMeta:   udpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "no-rule"},
					},
				),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakestore, err := store.NewFakeStore(store.FakeObjects{
				Gateways:  tc.gateways,
				UDPRoutes: tc.udpRoutes,
				Services:  tc.services,
			})
			require.NoError(t, err)
			translator := mustNewTranslator(t, fakestore)

			failureCollector := failures.NewResourceFailuresCollector(zapr.NewLogger(zap.NewNop()))
			translator.failuresCollector = failureCollector

			result := translator.ingressRulesFromUDPRoutes()
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
		gateways             []*gatewayapi.Gateway
		udpRoutes            []*gatewayapi.UDPRoute
		services             []*corev1.Service
		expectedKongServices []kongstate.Service
		expectedKongRoutes   map[string][]kongstate.Route
		expectedFailures     []failures.ResourceFailure
	}{
		{
			name: "UDPRoute with single rule, single backendref and Gateway in the same namespace",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "default",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
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
			udpRoutes: []*gatewayapi.UDPRoute{
				{
					TypeMeta: udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "single-rule",
						Namespace: "default",
					},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name: "gateway-1",
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
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
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name:     kong.String("udproute.default.single-rule.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"udproute.default.single-rule.0": {
					{
						Route: kong.Route{
							Name:       kong.String("udproute.default.single-rule.0.0"),
							Expression: kong.String("net.dst.port == 80"),
							Protocols:  kong.StringSlice("udp"),
						},
					},
				},
			},
		},
		{
			name: "UDPRoute with single rule, multiple backendrefs and Gateway in a different namespace",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "test-1",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
							builder.NewListener("udp80").WithPort(80).UDP().Build(),
							builder.NewListener("udp81").WithPort(81).UDP().Build(),
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
			udpRoutes: []*gatewayapi.UDPRoute{
				{
					TypeMeta: udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rule-1",
						Namespace: "default",
					},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:        "gateway-1",
									Namespace:   lo.ToPtr(gatewayapi.Namespace("test-1")),
									SectionName: lo.ToPtr(gatewayapi.SectionName("udp80")),
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(8080).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta: udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rule-2",
						Namespace: "default",
					},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:        "gateway-1",
									Namespace:   lo.ToPtr(gatewayapi.Namespace("test-1")),
									SectionName: lo.ToPtr(gatewayapi.SectionName("udp81")),
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service2").WithPort(8181).Build(),
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
						Name:     kong.String("udproute.default.rule-1.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service1").WithPortNumber(8080).MustBuild(),
					},
				},
				{
					Service: kong.Service{
						Name:     kong.String("udproute.default.rule-2.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service2").WithPortNumber(8181).MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"udproute.default.rule-1.0": {
					{
						Route: kong.Route{
							Name:       kong.String("udproute.default.rule-1.0.0"),
							Expression: kong.String("net.dst.port == 80"),
							Protocols:  kong.StringSlice("udp"),
						},
					},
				},
				"udproute.default.rule-2.0": {
					{
						Route: kong.Route{
							Name:       kong.String("udproute.default.rule-2.0.0"),
							Expression: kong.String("net.dst.port == 81"),
							Protocols:  kong.StringSlice("udp"),
						},
					},
				},
			},
		},
		{
			name: "single UDPRoute with single rule and multiple backendRefs",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "default",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
							builder.NewListener("udp80").WithPort(80).UDP().Build(),
							builder.NewListener("udp81").WithPort(81).UDP().Build(),
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
			udpRoutes: []*gatewayapi.UDPRoute{
				{
					TypeMeta: udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "multiple-backends",
						Namespace: "default",
					},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name: "gateway-1",
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
									builder.NewBackendRef("service1").WithPort(81).Build(),
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
						Name:     kong.String("udproute.default.multiple-backends.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
						builder.NewKongstateServiceBackend("service1").
							WithNamespace("default").
							WithPortNumber(81).
							MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"udproute.default.multiple-backends.0": {
					{
						Route: kong.Route{
							Name:       kong.String("udproute.default.multiple-backends.0.0"),
							Expression: kong.String("(net.dst.port == 80) || (net.dst.port == 81)"),
							Protocols:  kong.StringSlice("udp"),
						},
					},
				},
			},
		},
		{
			name: "multiple UDPRoutes with translation errors",
			gateways: []*gatewayapi.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gateway-1",
						Namespace: "default",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{
							builder.NewListener("udp80").WithPort(80).UDP().Build(),
							builder.NewListener("udp8080").WithPort(8080).UDP().Build(),
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
			udpRoutes: []*gatewayapi.UDPRoute{
				{
					TypeMeta:   udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{Name: "single-rule", Namespace: "default"},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:        "gateway-1",
									SectionName: lo.ToPtr(gatewayapi.SectionName("udp80")),
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service1").WithPort(80).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta:   udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{Name: "single-rule-2", Namespace: "default"},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{
								{
									Name:        "gateway-1",
									SectionName: lo.ToPtr(gatewayapi.SectionName("udp8080")),
								},
							},
						},
						Rules: []gatewayapi.UDPRouteRule{
							{
								BackendRefs: []gatewayapi.BackendRef{
									builder.NewBackendRef("service2").WithPort(8080).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta:   udpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{Name: "no-rule", Namespace: "default"},
					Spec:       gatewayapi.UDPRouteSpec{},
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
						Name:     kong.String("udproute.default.single-rule.0"),
						Protocol: kong.String("udp"),
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
						Name:     kong.String("udproute.default.single-rule-2.0"),
						Protocol: kong.String("udp"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service2").WithPortNumber(8080).MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"udproute.default.single-rule.0": {
					{
						Route: kong.Route{
							Name:       kong.String("udproute.default.single-rule.0.0"),
							Expression: kong.String("net.dst.port == 80"),
							Protocols:  kong.StringSlice("udp"),
						},
					},
				},
				"udproute.default.single-rule-2.0": {
					{
						Route: kong.Route{
							Name:       kong.String("udproute.default.single-rule-2.0.0"),
							Expression: kong.String("net.dst.port == 8080"),
							Protocols:  kong.StringSlice("udp"),
						},
					},
				},
			},
			expectedFailures: []failures.ResourceFailure{
				newResourceFailure(
					t, subtranslator.ErrRouteValidationNoRules.Error(),
					&gatewayapi.UDPRoute{
						TypeMeta:   udpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "no-rule"},
					},
				),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakestore, err := store.NewFakeStore(store.FakeObjects{
				Gateways:  tc.gateways,
				UDPRoutes: tc.udpRoutes,
				Services:  tc.services,
			})
			require.NoError(t, err)
			translator := mustNewTranslator(t, fakestore)
			translator.featureFlags.ExpressionRoutes = true

			failureCollector := failures.NewResourceFailuresCollector(zapr.NewLogger(zap.NewNop()))
			translator.failuresCollector = failureCollector

			result := translator.ingressRulesFromUDPRoutes()
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
