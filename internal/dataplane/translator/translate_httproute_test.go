package translator

import (
	"testing"

	"github.com/go-logr/zapr"
	"github.com/google/go-cmp/cmp"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

// httprouteGVK is the GVK for HTTPRoutes, needed in unit tests because
// we have to manually initialize objects that aren't retrieved from the
// Kubernetes API.
var httprouteGVK = schema.GroupVersionKind{
	Group:   "gateway.networking.k8s.io",
	Version: "v1beta1",
	Kind:    "HTTPRoute",
}

type testCaseIngressRulesFromHTTPRoutes struct {
	msg          string
	routes       []*gatewayapi.HTTPRoute
	storeObjects store.FakeObjects
	expected     func(routes []*gatewayapi.HTTPRoute) ingressRules
}

func TestValidateHTTPRoute(t *testing.T) {
	testCases := []struct {
		name             string
		httpRoute        *gatewayapi.HTTPRoute
		expressionRoutes bool
		expectedError    error
	}{
		{
			name: "valid HTTPRoute should pass the validation",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			},
			expressionRoutes: false,
			expectedError:    nil,
		},
		{
			name: "HTTPRoute with no rules should not pass the validation",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httproute-no-rule",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
				},
			},
			expressionRoutes: false,
			expectedError:    subtranslator.ErrRouteValidationNoRules,
		},
		{
			name: "HTTPRoute with query param match should pass validation with expression routes",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httproute-query-param-match",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
						Matches: builder.NewHTTPRouteMatch().WithQueryParam("foo", "bar").ToSlice(),
					}},
				},
			},
			expressionRoutes: true,
			expectedError:    nil,
		},
		{
			name: "HTTPRoute with query param match should not pass validation when expression routes disabled",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httproute-query-param-match",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
						Matches: builder.NewHTTPRouteMatch().WithQueryParam("foo", "bar").ToSlice(),
					}},
				},
			},
			expressionRoutes: false,
			expectedError:    subtranslator.ErrRouteValidationQueryParamMatchesUnsupported,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			featureFlags := FeatureFlags{
				ExpressionRoutes: tc.expressionRoutes,
			}
			err := validateHTTPRoute(tc.httpRoute, featureFlags)
			if tc.expectedError == nil {
				require.NoError(t, err, "should pass the validation")
			} else {
				require.ErrorIs(t, err, tc.expectedError, "should return expected error")
			}
		})
	}
}

func TestIngressRulesFromHTTPRoutes(t *testing.T) {
	testCases := []testCaseIngressRulesFromHTTPRoutes{
		{
			msg: "an empty list of HTTPRoutes should produce no ingress rules",
			expected: func(_ []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs:      newSecretNameToSNIs(),
					ServiceNameToParent:   map[string]client.Object{},
					ServiceNameToServices: make(map[string]kongstate.Service),
				}
			},
		},
		{
			msg: "an HTTPRoute rule with no matches can be routed if it has hostnames to match on",
			routes: []*gatewayapi.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "fake-service",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{ // only 1 service should be created
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created
								Route: kong.Route{
									Name:         kong.String("httproute.default.basic-httproute.0.0"),
									PreserveHost: kong.Bool(true),
									Protocols: []*string{
										kong.String("http"),
										kong.String("https"),
									},
									Hosts: []*string{
										kong.String("konghq.com"),
										kong.String("www.konghq.com"),
									},
									Tags: []*string{
										kong.String("k8s-name:basic-httproute"),
										kong.String("k8s-namespace:default"),
										kong.String("k8s-kind:HTTPRoute"),
										kong.String("k8s-group:gateway.networking.k8s.io"),
										kong.String("k8s-version:v1beta1"),
									},
								},
								Ingress: util.FromK8sObject(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "an HTTPRoute rule with no matches and no hostnames produces a catch-all rule",
			routes: []*gatewayapi.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					// no hostnames present
					Rules: []gatewayapi.HTTPRouteRule{{
						// no match rules present
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "fake-service",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{ // only 1 service should be created
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created
								Route: kong.Route{
									Name:         kong.String("httproute.default.basic-httproute.0.0"),
									PreserveHost: kong.Bool(true),
									Protocols: []*string{
										kong.String("http"),
										kong.String("https"),
									},
									Tags: []*string{
										kong.String("k8s-name:basic-httproute"),
										kong.String("k8s-namespace:default"),
										kong.String("k8s-kind:HTTPRoute"),
										kong.String("k8s-group:gateway.networking.k8s.io"),
										kong.String("k8s-version:v1beta1"),
									},
								},
								Ingress: util.FromK8sObject(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "a single HTTPRoute with one match and one backendRef results in a single service",
			routes: []*gatewayapi.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin").Build(),
						},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "fake-service",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{ // only 1 service should be created
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.0.0"),
									Paths: []*string{
										kong.String("~/httpbin$"),
										kong.String("/httpbin/"),
									},
									PreserveHost: kong.Bool(true),
									Protocols: []*string{
										kong.String("http"),
										kong.String("https"),
									},
									StripPath: lo.ToPtr(false),
									Tags: []*string{
										kong.String("k8s-name:basic-httproute"),
										kong.String("k8s-namespace:default"),
										kong.String("k8s-kind:HTTPRoute"),
										kong.String("k8s-group:gateway.networking.k8s.io"),
										kong.String("k8s-version:v1beta1"),
									},
								},
								Ingress: util.FromK8sObject(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "an HTTPRoute with regex path matches is supported",
			routes: []*gatewayapi.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathRegex("/httpbin$").Build(),
						},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "fake-service",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{ // only 1 service should be created
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.0.0"),
									Paths: []*string{
										kong.String("~/httpbin$"),
									},
									PreserveHost: kong.Bool(true),
									Protocols: []*string{
										kong.String("http"),
										kong.String("https"),
									},
									StripPath: lo.ToPtr(false),
									Tags: []*string{
										kong.String("k8s-name:basic-httproute"),
										kong.String("k8s-namespace:default"),
										kong.String("k8s-kind:HTTPRoute"),
										kong.String("k8s-group:gateway.networking.k8s.io"),
										kong.String("k8s-version:v1beta1"),
									},
								},
								Ingress: util.FromK8sObject(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "an HTTPRoute with exact path matches translates to a terminated Kong regex route",
			routes: []*gatewayapi.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathExact("/httpbin").Build(),
							},
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
							},
						},
					},
				},
			}},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "fake-service",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{ // only 1 service should be created
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.0.0"),
									Paths: []*string{
										kong.String("~/httpbin$"),
									},
									PreserveHost: kong.Bool(true),
									Protocols: []*string{
										kong.String("http"),
										kong.String("https"),
									},
									StripPath: lo.ToPtr(false),
									Tags: []*string{
										kong.String("k8s-name:basic-httproute"),
										kong.String("k8s-namespace:default"),
										kong.String("k8s-kind:HTTPRoute"),
										kong.String("k8s-group:gateway.networking.k8s.io"),
										kong.String("k8s-version:v1beta1"),
									},
								},
								Ingress: util.FromK8sObject(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "a single HTTPRoute with multiple rules with equal backendRefs results in a single service",
			routes: []*gatewayapi.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-1").Build(),
						},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}, {
						Matches: []gatewayapi.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-2").Build(),
						},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "fake-service",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{ // only 1 service should be created
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								// only 1 route with two paths should be created
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.0"),
										Paths: []*string{
											kong.String("~/httpbin-1$"),
											kong.String("/httpbin-1/"),
											kong.String("~/httpbin-2$"),
											kong.String("/httpbin-2/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},
							},
							Parent: routes[0],
						},
					},
				}
			},
		},

		{
			msg: "a single HTTPRoute with multiple rules with different backendRefs results in multiple services",
			routes: []*gatewayapi.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "basic-httproute",
						Namespace: corev1.NamespaceDefault,
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-1").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
							},
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-2").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(8080).Build(),
								},
							},
						},
					},
				},
			},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "fake-service",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
						"httproute.default.basic-httproute.1": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{ // 1 service per route should be created
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created for this service
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.0.0"),
									Paths: []*string{
										kong.String("~/httpbin-1$"),
										kong.String("/httpbin-1/"),
									},
									PreserveHost: kong.Bool(true),
									Protocols: []*string{
										kong.String("http"),
										kong.String("https"),
									},
									StripPath: lo.ToPtr(false),
									Tags: []*string{
										kong.String("k8s-name:basic-httproute"),
										kong.String("k8s-namespace:default"),
										kong.String("k8s-kind:HTTPRoute"),
										kong.String("k8s-group:gateway.networking.k8s.io"),
										kong.String("k8s-version:v1beta1"),
									},
								},
								Ingress: util.FromK8sObject(routes[0]),
							}},
							Parent: routes[0],
						},

						"httproute.default.basic-httproute.1": {
							Service: kong.Service{ // 1 service per route should be created
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.1"),
								Name:           kong.String("httproute.default.basic-httproute.1"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(8080).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.1.0"),
									Paths: []*string{
										kong.String("~/httpbin-2$"),
										kong.String("/httpbin-2/"),
									},
									PreserveHost: kong.Bool(true),
									Protocols: []*string{
										kong.String("http"),
										kong.String("https"),
									},
									StripPath: lo.ToPtr(false),
									Tags: []*string{
										kong.String("k8s-name:basic-httproute"),
										kong.String("k8s-namespace:default"),
										kong.String("k8s-kind:HTTPRoute"),
										kong.String("k8s-group:gateway.networking.k8s.io"),
										kong.String("k8s-version:v1beta1"),
									},
								},
								Ingress: util.FromK8sObject(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},

		{
			msg: "a single HTTPRoute with multiple rules and backendRefs generates consolidated routes",
			routes: []*gatewayapi.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "basic-httproute",
						Namespace: corev1.NamespaceDefault,
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-1").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
									builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
								},
							},
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-2").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
									builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
								},
							},
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-2").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("foo-v1").WithPort(8080).WithWeight(90).Build(),
									builder.NewHTTPBackendRef("foo-v3").WithPort(8080).WithWeight(10).Build(),
								},
							},
						},
					},
				},
			},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "foo-v1",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "foo-v2",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "foo-v3",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
						"httproute.default.basic-httproute.2": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("foo-v1").WithPortNumber(80).WithWeight(90).MustBuild(),
								builder.NewKongstateServiceBackend("foo-v2").WithPortNumber(8080).WithWeight(10).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.0"),
										Paths: []*string{
											kong.String("~/httpbin-1$"),
											kong.String("/httpbin-1/"),
											kong.String("~/httpbin-2$"),
											kong.String("/httpbin-2/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},
							},
							Parent: routes[0],
						},

						"httproute.default.basic-httproute.2": {
							Service: kong.Service{
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.2"),
								Name:           kong.String("httproute.default.basic-httproute.2"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("foo-v1").WithPortNumber(8080).WithWeight(90).MustBuild(),
								builder.NewKongstateServiceBackend("foo-v3").WithPortNumber(8080).WithWeight(10).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.2.0"),
										Paths: []*string{
											kong.String("~/httpbin-2$"),
											kong.String("/httpbin-2/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},
							},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "a single HTTPRoute with multiple rules with equal backendRefs and different filters results in a single service",
			routes: []*gatewayapi.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "basic-httproute",
						Namespace: corev1.NamespaceDefault,
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/path-0").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
								Filters: []gatewayapi.HTTPRouteFilter{
									{
										Type: gatewayapi.HTTPRouteFilterRequestHeaderModifier,
										RequestHeaderModifier: &gatewayapi.HTTPHeaderFilter{
											Add: []gatewayapi.HTTPHeader{
												{Name: "X-Test-Header-1", Value: "test-value-1"},
											},
										},
									},
								},
							},
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/path-1").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
								Filters: []gatewayapi.HTTPRouteFilter{
									{
										Type: gatewayapi.HTTPRouteFilterRequestHeaderModifier,
										RequestHeaderModifier: &gatewayapi.HTTPHeaderFilter{
											Add: []gatewayapi.HTTPHeader{
												{Name: "X-Test-Header-2", Value: "test-value-2"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "fake-service",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{ // only 1 service should be created
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								// two route  should be created, as the filters are different
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.0"),
										Paths: []*string{
											kong.String("~/path-0$"),
											kong.String("/path-0/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
									Plugins: []kong.Plugin{
										{
											Name: kong.String("request-transformer"),
											Config: kong.Configuration{
												"append": subtranslator.TransformerPluginConfig{
													Headers: []string{"X-Test-Header-1:test-value-1"},
												},
											},
											Tags: []*string{
												kong.String("k8s-name:basic-httproute"),
												kong.String("k8s-namespace:default"),
												kong.String("k8s-kind:HTTPRoute"),
												kong.String("k8s-group:gateway.networking.k8s.io"),
												kong.String("k8s-version:v1beta1"),
											},
										},
									},
								},
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.1.0"),
										Paths: []*string{
											kong.String("~/path-1$"),
											kong.String("/path-1/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
									Plugins: []kong.Plugin{
										{
											Name: kong.String("request-transformer"),
											Config: kong.Configuration{
												"append": subtranslator.TransformerPluginConfig{
													Headers: []string{"X-Test-Header-2:test-value-2"},
												},
											},
											Tags: []*string{
												kong.String("k8s-name:basic-httproute"),
												kong.String("k8s-namespace:default"),
												kong.String("k8s-kind:HTTPRoute"),
												kong.String("k8s-group:gateway.networking.k8s.io"),
												kong.String("k8s-version:v1beta1"),
											},
										},
									},
								},
							},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "a single HTTPRoute with single rule and multiple matches generates consolidated kong route paths",
			routes: []*gatewayapi.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								// Two matches eligible for consolidation into a single kong route
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-0").Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-1").Build(),
								// Other two matches eligible for consolidation, but not with the above two
								// as they have different methods
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-2").WithMethod(gatewayapi.HTTPMethodDelete).Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-3").WithMethod(gatewayapi.HTTPMethodDelete).Build(),
								// Other two matches eligible for consolidation, but not with the above two
								// as they have different headers
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-4").
									WithHeader("x-header-1", "x-value-1").
									WithHeader("x-header-2", "x-value-2").
									Build(),
								// Note the different header order
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-5").
									WithHeader("x-header-2", "x-value-2").
									WithHeader("x-header-1", "x-value-1").
									Build(),
							},
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
								builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
							},
						},
					},
				},
			}},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "foo-v1",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "foo-v2",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},

							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("foo-v1").WithPortNumber(80).WithWeight(90).MustBuild(),
								builder.NewKongstateServiceBackend("foo-v2").WithPortNumber(8080).WithWeight(10).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								// First two matches consolidated into a single route
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.0"),
										Paths: []*string{
											kong.String("~/path-0$"),
											kong.String("/path-0/"),
											kong.String("~/path-1$"),
											kong.String("/path-1/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},
								// Second two matches consolidated into a single route
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.2"),
										Paths: []*string{
											kong.String("~/path-2$"),
											kong.String("/path-2/"),
											kong.String("~/path-3$"),
											kong.String("/path-3/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Methods:   []*string{kong.String("DELETE")},
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},
								// Third two matches consolidated into a single route
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.4"),
										Paths: []*string{
											kong.String("~/path-4$"),
											kong.String("/path-4/"),
											kong.String("~/path-5$"),
											kong.String("/path-5/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Headers: map[string][]string{
											"x-header-1": {"x-value-1"},
											"x-header-2": {"x-value-2"},
										},
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},
							},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "a single HTTPRoute with multiple rules and matches generates consolidated kong route paths",
			routes: []*gatewayapi.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayapi.HTTPRouteRule{
						// Rule one has four matches, that can be consolidated into two kong routes
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								// Two matches eligible for consolidation into a single kong route
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-0").Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-1").Build(),
								// Other two matches eligible for consolidation, but not with the above two
								// as they have different methods
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-2").WithMethod(gatewayapi.HTTPMethodDelete).Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-3").WithMethod(gatewayapi.HTTPMethodDelete).Build(),
							},
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
								builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
							},
						},

						// Rule two:
						//	- shares the backend refs with rule one
						//	- has two matches, that can be consolidated with the first two matches of rule one
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-4").Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-5").Build(),
							},
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
								builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
							},
						},

						// Rule three:
						//	- shares the backend refs with rule one
						//	- has two matches, that potentially could be consolidated with the first match of rule one
						//	- has a different filter than rule one, thus cannot be consolidated
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-6").Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-7").Build(),
							},
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
								builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
							},
							Filters: []gatewayapi.HTTPRouteFilter{
								{
									Type: gatewayapi.HTTPRouteFilterRequestHeaderModifier,
									RequestHeaderModifier: &gatewayapi.HTTPHeaderFilter{
										Add: []gatewayapi.HTTPHeader{
											{Name: "X-Test-Header-1", Value: "test-value-1"},
										},
									},
								},
							},
						},
					},
				},
			}},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "foo-v1",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "foo-v2",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("foo-v1").WithPortNumber(80).WithWeight(90).MustBuild(),
								builder.NewKongstateServiceBackend("foo-v2").WithPortNumber(8080).WithWeight(10).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								// First two matches from rule one and rule two consolidated into a single route
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.0"),
										Paths: []*string{
											kong.String("~/path-0$"),
											kong.String("/path-0/"),
											kong.String("~/path-1$"),
											kong.String("/path-1/"),
											kong.String("~/path-4$"),
											kong.String("/path-4/"),
											kong.String("~/path-5$"),
											kong.String("/path-5/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},
								// Second two matches consolidated into a single route
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.2"),
										Paths: []*string{
											kong.String("~/path-2$"),
											kong.String("/path-2/"),
											kong.String("~/path-3$"),
											kong.String("/path-3/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Methods:   []*string{kong.String("DELETE")},
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},

								// Matches from rule 3, that has different filter, are not consolidated
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.2.0"),
										Paths: []*string{
											kong.String("~/path-6$"),
											kong.String("/path-6/"),
											kong.String("~/path-7$"),
											kong.String("/path-7/"),
										},
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:basic-httproute"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1beta1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
									Plugins: []kong.Plugin{
										{
											Name: kong.String("request-transformer"),
											Config: kong.Configuration{
												"append": subtranslator.TransformerPluginConfig{
													Headers: []string{"X-Test-Header-1:test-value-1"},
												},
											},
											Tags: []*string{
												kong.String("k8s-name:basic-httproute"),
												kong.String("k8s-namespace:default"),
												kong.String("k8s-kind:HTTPRoute"),
												kong.String("k8s-group:gateway.networking.k8s.io"),
												kong.String("k8s-version:v1beta1"),
											},
										},
									},
								},
							},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "a single HTTPRoute with timeouts will set the timeout in the service",
			routes: []*gatewayapi.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
						Timeouts: func() *gatewayapi.HTTPRouteTimeouts {
							timeout := gatewayapi.Duration("500ms")
							return &gatewayapi.HTTPRouteTimeouts{
								BackendRequest: &timeout,
							}
						}(),
					}},
				},
			}},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: corev1.NamespaceDefault,
							Name:      "fake-service",
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{ // only 1 service should be created
								ConnectTimeout: kong.Int(500),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(500),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(500),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created
								Route: kong.Route{
									Name:         kong.String("httproute.default.basic-httproute.0.0"),
									PreserveHost: kong.Bool(true),
									Protocols: []*string{
										kong.String("http"),
										kong.String("https"),
									},
									Hosts: []*string{
										kong.String("konghq.com"),
										kong.String("www.konghq.com"),
									},
									Tags: []*string{
										kong.String("k8s-name:basic-httproute"),
										kong.String("k8s-namespace:default"),
										kong.String("k8s-kind:HTTPRoute"),
										kong.String("k8s-group:gateway.networking.k8s.io"),
										kong.String("k8s-version:v1beta1"),
									},
								},
								Ingress: util.FromK8sObject(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.msg, func(t *testing.T) {
			fakestore, err := store.NewFakeStore(tt.storeObjects)
			require.NoError(t, err)

			p := mustNewTranslator(t, fakestore)

			ingressRules := newIngressRules()

			for _, httproute := range tt.routes {
				// initialize the HTTPRoute object
				httproute.SetGroupVersionKind(httprouteGVK)
			}

			p.ingressRulesFromHTTPRoutesWithCombinedService(tt.routes, &ingressRules)
			// verify that we receive the expected values
			expectedIngressRules := tt.expected(tt.routes)
			assert.Empty(t, cmp.Diff(expectedIngressRules, ingressRules, cmp.AllowUnexported(SecretNameToSNIs{}, kongstate.ServiceBackend{})))
		})
	}
}

func TestIngressRulesFromHTTPRoutesCombinedServicesAcrossHTTPRoutes(t *testing.T) {
	var (
		serviceTypeMeta = metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		}
		httpRouteTypeMeta = metav1.TypeMeta{
			APIVersion: string(gatewayapi.V1Group) + "/" + gatewayapi.V1GroupVersion,
			Kind:       "HTTPRoute",
		}

		fakeService = &corev1.Service{
			TypeMeta: serviceTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-service",
				Namespace: "default",
			},
		}
	)

	// test cases when combined services from different HTTPRoute is enabled.
	testCasesCombinedServicesAcrossHTTPRoutes := []testCaseIngressRulesFromHTTPRoutes{
		{
			msg: "an empty list of HTTPRoutes should produce no ingress rules",
			expected: func(_ []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs:      newSecretNameToSNIs(),
					ServiceNameToParent:   map[string]client.Object{},
					ServiceNameToServices: make(map[string]kongstate.Service),
				}
			},
		},
		{
			msg: "a single HTTPRoute with multiple rules sharing the same backend should be translated to one service",
			routes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-1",
						Namespace: "default",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
								},
							},
						},
					},
				},
			},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{fakeService},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.svc.default.fake-service.80": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.svc.default.fake-service.80": {
							Service: kong.Service{
								Host:           kong.String("httproute.default.svc.default.fake-service.80"),
								Name:           kong.String("httproute.default.svc.default.fake-service.80"),
								ConnectTimeout: kong.Int(60000),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: []kongstate.ServiceBackend{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Parent:    routes[0],
							Routes: []kongstate.Route{
								// Only 1 route should be created.
								{
									Route: kong.Route{
										Name:         kong.String("httproute.default.httproute-1.0.0"),
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										Paths: kong.StringSlice(
											"~/foo$",
											"~/bar$",
										),
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:httproute-1"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},
							},
						},
					},
				}
			},
		},
		{
			msg: "rules with the same backends in multiple HTTPRoutes from same namespace should be translated to one service",
			routes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-1",
						Namespace: "default",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-2",
						Namespace: "default",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
								},
							},
						},
					},
				},
			},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{fakeService},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.svc.default.fake-service.80": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.svc.default.fake-service.80": {
							Service: kong.Service{
								Host:           kong.String("httproute.default.svc.default.fake-service.80"),
								Name:           kong.String("httproute.default.svc.default.fake-service.80"),
								ConnectTimeout: kong.Int(60000),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: []kongstate.ServiceBackend{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Parent:    routes[0],
							Routes: []kongstate.Route{
								// 2 routes should be created for the 2 HTTPRoutes.
								{
									Route: kong.Route{
										Name:         kong.String("httproute.default.httproute-1.0.0"),
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										Paths: kong.StringSlice(
											"~/foo$",
										),
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:httproute-1"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},
								{
									Route: kong.Route{
										Name:         kong.String("httproute.default.httproute-2.0.0"),
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										Paths: kong.StringSlice(
											"~/bar$",
										),
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:httproute-2"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1"),
										},
									},
									Ingress: util.FromK8sObject(routes[1]),
								},
							},
						},
					},
				}
			},
		},
		{
			msg: "rules with the same backends in HTTPRoutes from different namespaces should be translated into multiple services",
			routes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-1",
						Namespace: "default",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-2",
						Namespace: "ns-alter",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithNamespace("default").WithPort(80).Build(),
								},
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
								},
							},
						},
					},
				},
			},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{fakeService},
				ReferenceGrants: []*gatewayapi.ReferenceGrant{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "grant-from-ns-alter-to-services",
							Namespace: "default",
						},
						Spec: gatewayapi.ReferenceGrantSpec{
							From: []gatewayapi.ReferenceGrantFrom{
								{
									Group:     gatewayapi.V1Group,
									Kind:      "HTTPRoute",
									Namespace: "ns-alter",
								},
							},
							To: []gatewayapi.ReferenceGrantTo{
								{
									Group: gatewayapi.Group(corev1.GroupName),
									Kind:  gatewayapi.Kind("Service"),
								},
							},
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.svc.default.fake-service.80":  routes[0],
						"httproute.ns-alter.svc.default.fake-service.80": routes[1],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.svc.default.fake-service.80": {
							Service: kong.Service{
								Host:           kong.String("httproute.default.svc.default.fake-service.80"),
								Name:           kong.String("httproute.default.svc.default.fake-service.80"),
								ConnectTimeout: kong.Int(60000),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: []kongstate.ServiceBackend{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Parent:    routes[0],
							Routes: []kongstate.Route{
								{
									Route: kong.Route{
										Name:         kong.String("httproute.default.httproute-1.0.0"),
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										Paths: kong.StringSlice(
											"~/foo$",
										),
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:httproute-1"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
								},
							},
						},
						"httproute.ns-alter.svc.default.fake-service.80": {
							Service: kong.Service{
								Host:           kong.String("httproute.ns-alter.svc.default.fake-service.80"),
								Name:           kong.String("httproute.ns-alter.svc.default.fake-service.80"),
								ConnectTimeout: kong.Int(60000),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: []kongstate.ServiceBackend{
								builder.NewKongstateServiceBackend("fake-service").WithNamespace("default").WithPortNumber(80).MustBuild(),
							},
							Namespace: "ns-alter",
							Parent:    routes[1],
							Routes: []kongstate.Route{
								{
									Route: kong.Route{
										Name:         kong.String("httproute.ns-alter.httproute-2.0.0"),
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										Paths: kong.StringSlice(
											"~/bar$",
										),
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:httproute-2"),
											kong.String("k8s-namespace:ns-alter"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1"),
										},
									},
									Ingress: util.FromK8sObject(routes[1]),
								},
							},
						},
					},
				}
			},
		},
		{
			msg: "rules with no backends but RequestRedirect filter should be translated to one service",
			routes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-1",
						Namespace: "default",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Filters: []gatewayapi.HTTPRouteFilter{
									builder.NewHTTPRouteRequestRedirectFilter().WithRequestRedirectHost("konghq.com").WithRequestRedirectStatusCode(301).Build(),
								},
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/kong").Build(),
								},
							},
							{
								Filters: []gatewayapi.HTTPRouteFilter{
									builder.NewHTTPRouteRequestRedirectFilter().WithRequestRedirectHost("kumahq.com").WithRequestRedirectStatusCode(301).Build(),
								},
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/kuma").Build(),
								},
							},
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.svc._": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.svc._": {
							Service: kong.Service{
								Host:           kong.String("httproute.default.svc._"),
								Name:           kong.String("httproute.default.svc._"),
								ConnectTimeout: kong.Int(60000),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends:  []kongstate.ServiceBackend{},
							Namespace: "default",
							Parent:    routes[0],
							Routes: []kongstate.Route{
								// 2 routes where one for each rule should be created because they have different filters.
								{
									Route: kong.Route{
										Name:         kong.String("httproute.default.httproute-1.0.0"),
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										Paths: kong.StringSlice(
											"~/kong$",
										),
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:httproute-1"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
									Plugins: []kong.Plugin{
										{
											Name: kong.String("request-termination"),
											Config: kong.Configuration{
												"status_code": kong.Int(301),
											},
											Tags: []*string{
												kong.String("k8s-name:httproute-1"),
												kong.String("k8s-namespace:default"),
												kong.String("k8s-kind:HTTPRoute"),
												kong.String("k8s-group:gateway.networking.k8s.io"),
												kong.String("k8s-version:v1"),
											},
										},
										{
											Name: kong.String("response-transformer"),
											Config: kong.Configuration{
												"add": subtranslator.TransformerPluginConfig{
													Headers: []string{"Location: http://konghq.com:80/kong"},
												},
											},
											Tags: []*string{
												kong.String("k8s-name:httproute-1"),
												kong.String("k8s-namespace:default"),
												kong.String("k8s-kind:HTTPRoute"),
												kong.String("k8s-group:gateway.networking.k8s.io"),
												kong.String("k8s-version:v1"),
											},
										},
									},
								},
								{
									Route: kong.Route{
										Name:         kong.String("httproute.default.httproute-1.1.0"),
										PreserveHost: kong.Bool(true),
										Protocols: []*string{
											kong.String("http"),
											kong.String("https"),
										},
										Paths: kong.StringSlice(
											"~/kuma$",
										),
										StripPath: lo.ToPtr(false),
										Tags: []*string{
											kong.String("k8s-name:httproute-1"),
											kong.String("k8s-namespace:default"),
											kong.String("k8s-kind:HTTPRoute"),
											kong.String("k8s-group:gateway.networking.k8s.io"),
											kong.String("k8s-version:v1"),
										},
									},
									Ingress: util.FromK8sObject(routes[0]),
									Plugins: []kong.Plugin{
										{
											Name: kong.String("request-termination"),
											Config: kong.Configuration{
												"status_code": kong.Int(301),
											},
											Tags: []*string{
												kong.String("k8s-name:httproute-1"),
												kong.String("k8s-namespace:default"),
												kong.String("k8s-kind:HTTPRoute"),
												kong.String("k8s-group:gateway.networking.k8s.io"),
												kong.String("k8s-version:v1"),
											},
										},
										{
											Name: kong.String("response-transformer"),
											Config: kong.Configuration{
												"add": subtranslator.TransformerPluginConfig{
													Headers: []string{"Location: http://kumahq.com:80/kuma"},
												},
											},
											Tags: []*string{
												kong.String("k8s-name:httproute-1"),
												kong.String("k8s-namespace:default"),
												kong.String("k8s-kind:HTTPRoute"),
												kong.String("k8s-group:gateway.networking.k8s.io"),
												kong.String("k8s-version:v1"),
											},
										},
									},
								},
							},
						},
					},
				}
			},
		},
	}
	for _, tc := range testCasesCombinedServicesAcrossHTTPRoutes {
		t.Run(tc.msg, func(t *testing.T) {
			fakestore, err := store.NewFakeStore(tc.storeObjects)
			require.NoError(t, err)

			p := mustNewTranslator(t, fakestore)
			// Set feature flag CombinedServicesFromDifferentHTTPRoutes to true to enable combining services from different HTTPRoutes.
			p.featureFlags.CombinedServicesFromDifferentHTTPRoutes = true

			ingressRules := newIngressRules()
			p.ingressRulesFromHTTPRoutesWithCombinedService(tc.routes, &ingressRules)
			// verify that we receive the expected values
			expectedIngressRules := tc.expected(tc.routes)
			assert.Empty(t, cmp.Diff(expectedIngressRules, ingressRules, cmp.AllowUnexported(SecretNameToSNIs{}, kongstate.ServiceBackend{})))
		})
	}
}

func TestGetHTTPRouteHostnamesAsSliceOfStringPointers(t *testing.T) {
	for _, tt := range []struct {
		msg      string
		input    *gatewayapi.HTTPRoute
		expected []*string
	}{
		{
			msg:      "an HTTPRoute with no hostnames produces no hostnames",
			input:    &gatewayapi.HTTPRoute{},
			expected: []*string{},
		},
		{
			msg: "an HTTPRoute with a single hostname produces a list with that one hostname",
			input: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
					},
				},
			},
			expected: []*string{
				kong.String("konghq.com"),
			},
		},
		{
			msg: "an HTTPRoute with multiple hostnames produces a list with the same hostnames",
			input: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						"konghq.com",
						"www.konghq.com",
						"docs.konghq.com",
					},
				},
			},
			expected: []*string{
				kong.String("konghq.com"),
				kong.String("www.konghq.com"),
				kong.String("docs.konghq.com"),
			},
		},
	} {
		t.Run(tt.msg, func(t *testing.T) {
			assert.Equal(t, tt.expected, getHTTPRouteHostnamesAsSliceOfStringPointers(tt.input))
		})
	}
}

func TestIngressRulesFromHTTPRoutes_RegexPrefix(t *testing.T) {
	for _, tt := range []testCaseIngressRulesFromHTTPRoutes{
		{
			msg: "an HTTPRoute with regex path matches is supported",
			routes: []*gatewayapi.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathRegex("/httpbin$").Build(),
						},
						BackendRefs: []gatewayapi.HTTPBackendRef{{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Name: gatewayapi.ObjectName("fake-service"),
									Port: lo.ToPtr(gatewayapi.PortNumber(80)),
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
							},
						}},
					}},
				},
			}},
			storeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "fake-service",
							Namespace: corev1.NamespaceDefault,
						},
					},
				},
			},
			expected: func(routes []*gatewayapi.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs: newSecretNameToSNIs(),
					ServiceNameToParent: map[string]client.Object{
						"httproute.default.basic-httproute.0": routes[0],
					},
					ServiceNameToServices: map[string]kongstate.Service{
						"httproute.default.basic-httproute.0": {
							Service: kong.Service{ // only 1 service should be created
								ConnectTimeout: kong.Int(60000),
								Host:           kong.String("httproute.default.basic-httproute.0"),
								Name:           kong.String("httproute.default.basic-httproute.0"),
								Protocol:       kong.String("http"),
								ReadTimeout:    kong.Int(60000),
								Retries:        kong.Int(5),
								WriteTimeout:   kong.Int(60000),
							},
							Backends: kongstate.ServiceBackends{
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).MustBuild(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.0.0"),
									Paths: []*string{
										kong.String("~/httpbin$"),
									},
									PreserveHost: kong.Bool(true),
									Protocols: []*string{
										kong.String("http"),
										kong.String("https"),
									},
									StripPath: lo.ToPtr(false),
									Tags: []*string{
										kong.String("k8s-name:basic-httproute"),
										kong.String("k8s-namespace:default"),
										kong.String("k8s-kind:HTTPRoute"),
										kong.String("k8s-group:gateway.networking.k8s.io"),
										kong.String("k8s-version:v1beta1"),
									},
								},
								Ingress: util.FromK8sObject(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
	} {
		withTranslator := func(tran *Translator) func(t *testing.T) {
			return func(t *testing.T) {
				ingressRules := newIngressRules()

				for _, httproute := range tt.routes {
					// initialize the HTTPRoute object
					httproute.SetGroupVersionKind(httprouteGVK)
				}
				tran.ingressRulesFromHTTPRoutesWithCombinedService(tt.routes, &ingressRules)

				// verify that we receive the expected values
				expectedIngressRules := tt.expected(tt.routes)
				assert.Equal(t, expectedIngressRules, ingressRules)
			}
		}

		fakestore, err := store.NewFakeStore(tt.storeObjects)
		require.NoError(t, err)
		translator := mustNewTranslator(t, fakestore)

		t.Run(tt.msg+" using combined service routes translator", withTranslator(translator))
	}
}

func TestIngressRulesFromHTTPRoutesUsingExpressionRoutes(t *testing.T) {
	httpRouteTypeMeta := metav1.TypeMeta{Kind: "HTTPRoute", APIVersion: gatewayv1beta1.GroupVersion.String()}

	testCases := []struct {
		name                 string
		httpRoutes           []*gatewayapi.HTTPRoute
		expectedKongServices []kongstate.Service
		expectedKongRoutes   map[string][]kongstate.Route
		fakeObjects          store.FakeObjects
	}{
		{
			name: "single HTTPRoute with no hostname and multiple matches",
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
									builder.NewHTTPRouteMatch().WithPathExact("/v1/barr").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			fakeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "service-1",
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-1.80"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service-1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"httproute.default.svc.default.service-1.80": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1._.0.0"),
							Expression:   kong.String(`http.path == "/v1/foo"`),
							PreserveHost: kong.Bool(true),
						},
						ExpressionRoutes: true,
					},
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1._.0.1"),
							Expression:   kong.String(`http.path == "/v1/barr"`),
							PreserveHost: kong.Bool(true),
						},
						ExpressionRoutes: true,
					},
				},
			},
		},
		{
			name: "single HTTPRoute with multiple hostnames and rules",
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Hostnames: []gatewayapi.Hostname{
							"foo.com",
							"*.bar.com",
						},
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
							},
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/v1/barr").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-2").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			fakeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "service-1",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "service-2",
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-1.80"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service-1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
				{
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-2.80"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service-2").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"httproute.default.svc.default.service-1.80": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1.foo.com.0.0"),
							Expression:   kong.String(`(http.host == "foo.com") && (http.path == "/v1/foo")`),
							PreserveHost: kong.Bool(true),
						},
						ExpressionRoutes: true,
					},
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1._.bar.com.0.0"),
							Expression:   kong.String(`(http.host =^ ".bar.com") && (http.path == "/v1/foo")`),
							PreserveHost: kong.Bool(true),
						},
						ExpressionRoutes: true,
					},
				},
				"httproute.default.svc.default.service-2.80": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1.foo.com.1.0"),
							Expression:   kong.String(`(http.host == "foo.com") && (http.path == "/v1/barr")`),
							PreserveHost: kong.Bool(true),
						},
						ExpressionRoutes: true,
					},
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1._.bar.com.1.0"),
							Expression:   kong.String(`(http.host =^ ".bar.com") && (http.path == "/v1/barr")`),
							PreserveHost: kong.Bool(true),
						},
						ExpressionRoutes: true,
					},
				},
			},
		},
		{
			name: "single HTTPRoute with protocol and SNI annotations",
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
						Annotations: map[string]string{
							"konghq.com/protocols": "https",
							"konghq.com/snis":      "foo.com",
						},
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Hostnames: []gatewayapi.Hostname{
							"foo.com",
						},
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			fakeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "service-1",
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-1.80"),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service-1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"httproute.default.svc.default.service-1.80": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1.foo.com.0.0"),
							Expression:   kong.String(`(http.host == "foo.com") && (tls.sni == "foo.com") && (http.path == "/v1/foo")`),
							PreserveHost: kong.Bool(true),
						},
						ExpressionRoutes: true,
					},
				},
			},
		},
		{
			name: "single HTTPRoute with backendTimeout configuration",
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
									builder.NewHTTPRouteMatch().WithPathExact("/v1/barr").Build(),
								},
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
								Timeouts: func() *gatewayapi.HTTPRouteTimeouts {
									timeout := gatewayapi.Duration("500ms")
									return &gatewayapi.HTTPRouteTimeouts{
										BackendRequest: &timeout,
									}
								}(),
							},
						},
					},
				},
			},
			fakeObjects: store.FakeObjects{
				Services: []*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "service-1",
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name:           kong.String("httproute.default.svc.default.service-1.80"),
						ConnectTimeout: kong.Int(500),
						ReadTimeout:    kong.Int(500),
						WriteTimeout:   kong.Int(500),
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("service-1").
							WithNamespace("default").
							WithPortNumber(80).
							MustBuild(),
					},
				},
			},
			expectedKongRoutes: map[string][]kongstate.Route{
				"httproute.default.svc.default.service-1.80": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1._.0.0"),
							Expression:   kong.String(`http.path == "/v1/foo"`),
							PreserveHost: kong.Bool(true),
						},
						ExpressionRoutes: true,
					},
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1._.0.1"),
							Expression:   kong.String(`http.path == "/v1/barr"`),
							PreserveHost: kong.Bool(true),
						},
						ExpressionRoutes: true,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakestore, err := store.NewFakeStore(tc.fakeObjects)
			require.NoError(t, err)
			translator := mustNewTranslator(t, fakestore)
			translator.featureFlags.ExpressionRoutes = true
			translator.failuresCollector = failures.NewResourceFailuresCollector(zapr.NewLogger(zap.NewNop()))

			result := newIngressRules()
			translator.ingressRulesFromHTTPRoutesWithCombinedService(tc.httpRoutes, &result)
			// check services
			require.Equal(t, len(tc.expectedKongServices), len(result.ServiceNameToServices),
				"should have expected number of services")
			for _, expectedKongService := range tc.expectedKongServices {
				kongService, ok := result.ServiceNameToServices[*expectedKongService.Name]
				require.Truef(t, ok, "should find service %s", *expectedKongService.Name)
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
		})
	}
}

func commonRouteSpecMock(parentReferentName string) gatewayapi.CommonRouteSpec {
	return gatewayapi.CommonRouteSpec{
		ParentRefs: []gatewayapi.ParentReference{{
			Name: gatewayapi.ObjectName(parentReferentName),
		}},
	}
}
