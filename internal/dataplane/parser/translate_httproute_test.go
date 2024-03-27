package parser

import (
	"strings"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
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
	msg      string
	routes   []*gatewayv1.HTTPRoute
	expected func(routes []*gatewayv1.HTTPRoute) ingressRules
	errs     []error
}

// common test cases  should work with legacy parser and combined routes parser.
func getIngressRulesFromHTTPRoutesCommonTestCases() []testCaseIngressRulesFromHTTPRoutes {
	return []testCaseIngressRulesFromHTTPRoutes{
		{
			msg: "an empty list of HTTPRoutes should produce no ingress rules",
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs:      newSecretNameToSNIs(),
					ServiceNameToParent:   map[string]client.Object{},
					ServiceNameToServices: make(map[string]kongstate.Service),
				}
			},
		},
		{
			msg: "an HTTPRoute rule with no matches can be routed if it has hostnames to match on",
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
					Hostnames: []gatewayv1.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
					Rules: []gatewayv1.HTTPRouteRule{{
						BackendRefs: []gatewayv1.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).Build(),
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
								Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "an HTTPRoute rule with no matches and no hostnames produces a catch-all rule",
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					// no hostnames present
					Rules: []gatewayv1.HTTPRouteRule{{
						// no match rules present
						BackendRefs: []gatewayv1.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).Build(),
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
								Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "a single HTTPRoute with one match and one backendRef results in a single service",
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayv1.HTTPRouteRule{{
						Matches: []gatewayv1.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin").Build(),
						},
						BackendRefs: []gatewayv1.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.0.0"),
									Paths: []*string{
										kong.String("/httpbin"),
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
								Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "an HTTPRoute with no rules can't be routed",
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs:      newSecretNameToSNIs(),
					ServiceNameToParent:   map[string]client.Object{},
					ServiceNameToServices: make(map[string]kongstate.Service),
				}
			},
			errs: []error{
				translators.ErrRouteValidationNoRules,
			},
		},
		{
			msg: "an HTTPRoute with queryParam matches is not yet supported",
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayv1.HTTPRouteRule{{
						Matches: []gatewayv1.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithQueryParam("username", "kong").Build(),
						},
						BackendRefs: []gatewayv1.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
				return ingressRules{
					SecretNameToSNIs:      newSecretNameToSNIs(),
					ServiceNameToParent:   map[string]client.Object{},
					ServiceNameToServices: make(map[string]kongstate.Service),
				}
			},
			errs: []error{
				translators.ErrRouteValidationQueryParamMatchesUnsupported,
			},
		},
		{
			msg: "an HTTPRoute with regex path matches is supported",
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayv1.HTTPRouteRule{{
						Matches: []gatewayv1.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathRegex("/httpbin$").Build(),
						},
						BackendRefs: []gatewayv1.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.0.0"),
									Paths: []*string{
										kong.String("/httpbin$"),
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
								Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
		{
			msg: "an HTTPRoute with exact path matches translates to a terminated Kong regex route",
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayv1.HTTPRouteRule{
						{
							Matches: []gatewayv1.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathExact("/httpbin").Build(),
							},
							BackendRefs: []gatewayv1.HTTPBackendRef{
								builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
							},
						},
					},
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.0.0"),
									Paths: []*string{
										kong.String("/httpbin$"),
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
								Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
	}
}

func getIngressRulesFromHTTPRoutesCombinedRoutesTestCases() []testCaseIngressRulesFromHTTPRoutes {
	return []testCaseIngressRulesFromHTTPRoutes{
		{
			msg: "a single HTTPRoute with multiple rules with equal backendRefs results in a single service",
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayv1.HTTPRouteRule{{
						Matches: []gatewayv1.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-1").Build(),
						},
						BackendRefs: []gatewayv1.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}, {
						Matches: []gatewayv1.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-2").Build(),
						},
						BackendRefs: []gatewayv1.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					}},
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								// only 1 route with two paths should be created
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.0"),
										Paths: []*string{
											kong.String("/httpbin-1"),
											kong.String("/httpbin-2"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
								},
							},
							Parent: routes[0],
						},
					},
				}
			},
		},

		{
			msg: "a single HTTPRoute with multiple rules with different backendRefs results in a multiple services",
			routes: []*gatewayv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "basic-httproute",
						Namespace: corev1.NamespaceDefault,
					},
					Spec: gatewayv1.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
						Rules: []gatewayv1.HTTPRouteRule{
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-1").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
							},
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-2").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(8080).Build(),
								},
							},
						},
					},
				},
			},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{ // only 1 route should be created for this service
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.0.0"),
									Paths: []*string{
										kong.String("/httpbin-1"),
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
								Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
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
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(8080).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{{
								Route: kong.Route{
									Name: kong.String("httproute.default.basic-httproute.1.0"),
									Paths: []*string{
										kong.String("/httpbin-2"),
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
								Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},

		{
			msg: "a single HTTPRoute with multiple rules and backendRefs generates consolidated routes",
			routes: []*gatewayv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "basic-httproute",
						Namespace: corev1.NamespaceDefault,
					},
					Spec: gatewayv1.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
						Rules: []gatewayv1.HTTPRouteRule{
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-1").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
									builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
								},
							},
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-2").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
									builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
								},
							},
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/httpbin-2").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("foo-v1").WithPort(8080).WithWeight(90).Build(),
									builder.NewHTTPBackendRef("foo-v3").WithPort(8080).WithWeight(10).Build(),
								},
							},
						},
					},
				},
			},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("foo-v1").WithPortNumber(80).WithWeight(90).Build(),
								builder.NewKongstateServiceBackend("foo-v2").WithPortNumber(8080).WithWeight(10).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.0"),
										Paths: []*string{
											kong.String("/httpbin-1"),
											kong.String("/httpbin-2"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
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
								builder.NewKongstateServiceBackend("foo-v1").WithPortNumber(8080).WithWeight(90).Build(),
								builder.NewKongstateServiceBackend("foo-v3").WithPortNumber(8080).WithWeight(10).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.2.0"),
										Paths: []*string{
											kong.String("/httpbin-2"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
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
			routes: []*gatewayv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "basic-httproute",
						Namespace: corev1.NamespaceDefault,
					},
					Spec: gatewayv1.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
						Rules: []gatewayv1.HTTPRouteRule{
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/path-0").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
								Filters: []gatewayv1.HTTPRouteFilter{
									{
										Type: gatewayv1.HTTPRouteFilterRequestHeaderModifier,
										RequestHeaderModifier: &gatewayv1.HTTPHeaderFilter{
											Add: []gatewayv1.HTTPHeader{
												{Name: "X-Test-Header-1", Value: "test-value-1"},
											},
										},
									},
								},
							},
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/path-1").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
								Filters: []gatewayv1.HTTPRouteFilter{
									{
										Type: gatewayv1.HTTPRouteFilterRequestHeaderModifier,
										RequestHeaderModifier: &gatewayv1.HTTPHeaderFilter{
											Add: []gatewayv1.HTTPHeader{
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
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								// two route  should be created, as the filters are different
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.0"),
										Paths: []*string{
											kong.String("/path-0"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
									Plugins: []kong.Plugin{
										{
											Name: kong.String("request-transformer"),
											Config: kong.Configuration{
												"append": map[string][]string{
													"headers": {"X-Test-Header-1:test-value-1"},
												},
											},
										},
									},
								},
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.1.0"),
										Paths: []*string{
											kong.String("/path-1"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
									Plugins: []kong.Plugin{
										{
											Name: kong.String("request-transformer"),
											Config: kong.Configuration{
												"append": map[string][]string{
													"headers": {"X-Test-Header-2:test-value-2"},
												},
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
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayv1.HTTPRouteRule{
						{
							Matches: []gatewayv1.HTTPRouteMatch{
								// Two matches eligible for consolidation into a single kong route
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-0").Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-1").Build(),
								// Other two matches eligible for consolidation, but not with the above two
								// as they have different methods
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-2").WithMethod(gatewayv1.HTTPMethodDelete).Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-3").WithMethod(gatewayv1.HTTPMethodDelete).Build(),
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
							BackendRefs: []gatewayv1.HTTPBackendRef{
								builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
								builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
							},
						},
					},
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("foo-v1").WithPortNumber(80).WithWeight(90).Build(),
								builder.NewKongstateServiceBackend("foo-v2").WithPortNumber(8080).WithWeight(10).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								// First two matches consolidated into a single route
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.0"),
										Paths: []*string{
											kong.String("/path-0"),
											kong.String("/path-1"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
								},
								// Second two matches consolidated into a single route
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.2"),
										Paths: []*string{
											kong.String("/path-2"),
											kong.String("/path-3"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
								},
								// Third two matches consolidated into a single route
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.4"),
										Paths: []*string{
											kong.String("/path-4"),
											kong.String("/path-5"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
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
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayv1.HTTPRouteRule{
						// Rule one has four matches, that can be consolidated into two kong routes
						{
							Matches: []gatewayv1.HTTPRouteMatch{
								// Two matches eligible for consolidation into a single kong route
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-0").Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-1").Build(),
								// Other two matches eligible for consolidation, but not with the above two
								// as they have different methods
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-2").WithMethod(gatewayv1.HTTPMethodDelete).Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-3").WithMethod(gatewayv1.HTTPMethodDelete).Build(),
							},
							BackendRefs: []gatewayv1.HTTPBackendRef{
								builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
								builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
							},
						},

						// Rule two:
						//	- shares the backend refs with rule one
						//	- has two matches, that can be consolidated with the first two matches of rule one
						{
							Matches: []gatewayv1.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-4").Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-5").Build(),
							},
							BackendRefs: []gatewayv1.HTTPBackendRef{
								builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
								builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
							},
						},

						// Rule three:
						//	- shares the backend refs with rule one
						//	- has two matches, that potentially could be consolidated with the first match of rule one
						//	- has a different filter than rule one, thus cannot be consolidated
						{
							Matches: []gatewayv1.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-6").Build(),
								builder.NewHTTPRouteMatch().WithPathPrefix("/path-7").Build(),
							},
							BackendRefs: []gatewayv1.HTTPBackendRef{
								builder.NewHTTPBackendRef("foo-v1").WithPort(80).WithWeight(90).Build(),
								builder.NewHTTPBackendRef("foo-v2").WithPort(8080).WithWeight(10).Build(),
							},
							Filters: []gatewayv1.HTTPRouteFilter{
								{
									Type: gatewayv1.HTTPRouteFilterRequestHeaderModifier,
									RequestHeaderModifier: &gatewayv1.HTTPHeaderFilter{
										Add: []gatewayv1.HTTPHeader{
											{Name: "X-Test-Header-1", Value: "test-value-1"},
										},
									},
								},
							},
						},
					},
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("foo-v1").WithPortNumber(80).WithWeight(90).Build(),
								builder.NewKongstateServiceBackend("foo-v2").WithPortNumber(8080).WithWeight(10).Build(),
							},
							Namespace: "default",
							Routes: []kongstate.Route{
								// First two matches from rule one and rule two consolidated into a single route
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.0"),
										Paths: []*string{
											kong.String("/path-0"),
											kong.String("/path-1"),
											kong.String("/path-4"),
											kong.String("/path-5"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
								},
								// Second two matches consolidated into a single route
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.0.2"),
										Paths: []*string{
											kong.String("/path-2"),
											kong.String("/path-3"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
								},

								// Matches from rule 3, that has different filter, are not consolidated
								{
									Route: kong.Route{
										Name: kong.String("httproute.default.basic-httproute.2.0"),
										Paths: []*string{
											kong.String("/path-6"),
											kong.String("/path-7"),
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
									Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
									Plugins: []kong.Plugin{
										{
											Name: kong.String("request-transformer"),
											Config: kong.Configuration{
												"append": map[string][]string{
													"headers": {"X-Test-Header-1:test-value-1"},
												},
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
	}
}

func TestIngressRulesFromHTTPRoutes(t *testing.T) {
	fakestore, err := store.NewFakeStore(store.FakeObjects{})
	require.NoError(t, err)

	testCases := getIngressRulesFromHTTPRoutesCommonTestCases()

	for _, tt := range testCases {
		t.Run(tt.msg, func(t *testing.T) {
			p := mustNewParser(t, fakestore)

			ingressRules := newIngressRules()

			var errs []error
			for _, httproute := range tt.routes {
				// initialize the HTTPRoute object
				httproute.SetGroupVersionKind(httprouteGVK)

				// generate the ingress rules
				err := p.ingressRulesFromHTTPRoute(&ingressRules, httproute)
				if err != nil {
					errs = append(errs, err)
				}
			}

			// verify that we receive the expected values
			expectedIngressRules := tt.expected(tt.routes)
			assert.Equal(t, expectedIngressRules, ingressRules)

			// verify that we receive any and all expected errors
			for i := range tt.errs {
				assert.ErrorIs(t, errs[i], tt.errs[i])
			}
		})
	}
}

func TestIngressRulesFromHTTPRoutesWithCombinedServiceRoutes(t *testing.T) {
	fakestore, err := store.NewFakeStore(store.FakeObjects{})
	require.NoError(t, err)

	testCases := getIngressRulesFromHTTPRoutesCommonTestCases()
	testCases = append(testCases, getIngressRulesFromHTTPRoutesCombinedRoutesTestCases()...)

	for _, tt := range testCases {
		t.Run(tt.msg, func(t *testing.T) {
			p := mustNewParser(t, fakestore)
			p.featureFlags.CombinedServiceRoutes = true

			ingressRules := newIngressRules()

			var errs []error
			for _, httproute := range tt.routes {
				// initialize the HTTPRoute object
				httproute.SetGroupVersionKind(httprouteGVK)

				// generate the ingress rules
				err := p.ingressRulesFromHTTPRoute(&ingressRules, httproute)
				if err != nil {
					errs = append(errs, err)
				}
			}

			// verify that we receive the expected values
			expectedIngressRules := tt.expected(tt.routes)
			assert.Equal(t, expectedIngressRules, ingressRules)

			// verify that we receive any and all expected errors
			for i := range tt.errs {
				assert.ErrorIs(t, errs[i], tt.errs[i])
			}
		})
	}
}

func TestGetHTTPRouteHostnamesAsSliceOfStringPointers(t *testing.T) {
	for _, tt := range []struct {
		msg      string
		input    *gatewayv1.HTTPRoute
		expected []*string
	}{
		{
			msg:      "an HTTPRoute with no hostnames produces no hostnames",
			input:    &gatewayv1.HTTPRoute{},
			expected: []*string{},
		},
		{
			msg: "an HTTPRoute with a single hostname produces a list with that one hostname",
			input: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
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
			input: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
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
	fakestore, err := store.NewFakeStore(store.FakeObjects{})
	require.NoError(t, err)
	parser := mustNewParser(t, fakestore)
	parser.featureFlags.RegexPathPrefix = true
	parserWithCombinedServiceRoutes := mustNewParser(t, fakestore)
	parserWithCombinedServiceRoutes.featureFlags.RegexPathPrefix = true
	parserWithCombinedServiceRoutes.featureFlags.CombinedServiceRoutes = true
	httpPort := gatewayv1.PortNumber(80)

	for _, tt := range []testCaseIngressRulesFromHTTPRoutes{
		{
			msg: "an HTTPRoute with regex path matches is supported",
			routes: []*gatewayv1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: commonRouteSpecMock("fake-gateway"),
					Rules: []gatewayv1.HTTPRouteRule{{
						Matches: []gatewayv1.HTTPRouteMatch{
							builder.NewHTTPRouteMatch().WithPathRegex("/httpbin$").Build(),
						},
						BackendRefs: []gatewayv1.HTTPBackendRef{{
							BackendRef: gatewayv1.BackendRef{
								BackendObjectReference: gatewayv1.BackendObjectReference{
									Name: gatewayv1.ObjectName("fake-service"),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
							},
						}},
					}},
				},
			}},
			expected: func(routes []*gatewayv1.HTTPRoute) ingressRules {
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
								builder.NewKongstateServiceBackend("fake-service").WithPortNumber(80).Build(),
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
								Ingress: k8sObjectInfoOfHTTPRoute(routes[0]),
							}},
							Parent: routes[0],
						},
					},
				}
			},
		},
	} {
		withParser := func(p *Parser) func(t *testing.T) {
			return func(t *testing.T) {
				ingressRules := newIngressRules()

				var errs []error
				for _, httproute := range tt.routes {
					// initialize the HTTPRoute object
					httproute.SetGroupVersionKind(httprouteGVK)

					// generate the ingress rules
					err := p.ingressRulesFromHTTPRoute(&ingressRules, httproute)
					if err != nil {
						errs = append(errs, err)
					}
				}

				// verify that we receive the expected values
				expectedIngressRules := tt.expected(tt.routes)
				assert.Equal(t, expectedIngressRules, ingressRules)

				// verify that we receive any and all expected errors
				assert.Equal(t, tt.errs, errs)
			}
		}

		t.Run(tt.msg+" using legacy parser", withParser(parser))
		t.Run(tt.msg+" using combined service routes parser", withParser(parserWithCombinedServiceRoutes))
	}
}

func TestIngressRulesFromHTTPRoutesUsingExpressionRoutes(t *testing.T) {
	fakestore, err := store.NewFakeStore(store.FakeObjects{})
	require.NoError(t, err)
	parser := mustNewParser(t, fakestore)
	parser.featureFlags.CombinedServiceRoutes = true
	parser.featureFlags.ExpressionRoutes = true
	httpRouteTypeMeta := metav1.TypeMeta{Kind: "HTTPRoute", APIVersion: gatewayv1.SchemeGroupVersion.String()}

	newResourceFailure := func(reason string, objects ...client.Object) failures.ResourceFailure {
		failure, _ := failures.NewResourceFailure(reason, objects...)
		return failure
	}

	testCases := []struct {
		name                 string
		httpRoutes           []*gatewayv1.HTTPRoute
		expectedKongServices []kongstate.Service
		expectedKongRoutes   map[string][]kongstate.Route
		expectedFailures     []failures.ResourceFailure
	}{
		{
			name: "single HTTPRoute with no hostname and multiple matches",
			httpRoutes: []*gatewayv1.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayv1.HTTPRouteSpec{
						Rules: []gatewayv1.HTTPRouteRule{
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
									builder.NewHTTPRouteMatch().WithPathExact("/v1/barr").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("service1").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("httproute.default.httproute-1._.0"),
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
				"httproute.default.httproute-1._.0": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1._.0.0"),
							Expression:   kong.String(`http.path == "/v1/foo"`),
							PreserveHost: kong.Bool(true),
						},
						Plugins:          []kong.Plugin{},
						ExpressionRoutes: true,
					},
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1._.0.1"),
							Expression:   kong.String(`http.path == "/v1/barr"`),
							PreserveHost: kong.Bool(true),
						},
						Plugins:          []kong.Plugin{},
						ExpressionRoutes: true,
					},
				},
			},
		},
		{
			name: "single HTTPRoute with multiple hostnames and rules",
			httpRoutes: []*gatewayv1.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayv1.HTTPRouteSpec{
						Hostnames: []gatewayv1.Hostname{
							"foo.com",
							"*.bar.com",
						},
						Rules: []gatewayv1.HTTPRouteRule{
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("service1").WithPort(80).Build(),
								},
							},
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/v1/barr").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("service2").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("httproute.default.httproute-1.foo.com.0"),
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
						Name: kong.String("httproute.default.httproute-1._.bar.com.0"),
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
						Name: kong.String("httproute.default.httproute-1.foo.com.1"),
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
						Name: kong.String("httproute.default.httproute-1._.bar.com.1"),
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
				"httproute.default.httproute-1.foo.com.0": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1.foo.com.0.0"),
							Expression:   kong.String(`(http.host == "foo.com") && (http.path == "/v1/foo")`),
							PreserveHost: kong.Bool(true),
						},
						Plugins:          []kong.Plugin{},
						ExpressionRoutes: true,
					},
				},
				"httproute.default.httproute-1._.bar.com.0": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1._.bar.com.0.0"),
							Expression:   kong.String(`(http.host =^ ".bar.com") && (http.path == "/v1/foo")`),
							PreserveHost: kong.Bool(true),
						},
						Plugins:          []kong.Plugin{},
						ExpressionRoutes: true,
					},
				},
				"httproute.default.httproute-1.foo.com.1": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1.foo.com.1.0"),
							Expression:   kong.String(`(http.host == "foo.com") && (http.path == "/v1/barr")`),
							PreserveHost: kong.Bool(true),
						},
						Plugins:          []kong.Plugin{},
						ExpressionRoutes: true,
					},
				},
				"httproute.default.httproute-1._.bar.com.1": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1._.bar.com.1.0"),
							Expression:   kong.String(`(http.host =^ ".bar.com") && (http.path == "/v1/barr")`),
							PreserveHost: kong.Bool(true),
						},
						Plugins:          []kong.Plugin{},
						ExpressionRoutes: true,
					},
				},
			},
		},
		{
			name: "single HTTPRoute with protocol and SNI annotations",
			httpRoutes: []*gatewayv1.HTTPRoute{
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
					Spec: gatewayv1.HTTPRouteSpec{
						Hostnames: []gatewayv1.Hostname{
							"foo.com",
						},
						Rules: []gatewayv1.HTTPRouteRule{
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("service1").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("httproute.default.httproute-1.foo.com.0"),
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
				"httproute.default.httproute-1.foo.com.0": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1.foo.com.0.0"),
							Expression:   kong.String(`(http.host == "foo.com") && (tls.sni == "foo.com") && (http.path == "/v1/foo")`),
							PreserveHost: kong.Bool(true),
						},
						Plugins:          []kong.Plugin{},
						ExpressionRoutes: true,
					},
				},
			},
		},
		{
			name: "multiple HTTPRoutes with translation failures",
			httpRoutes: []*gatewayv1.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-no-host-no-rule",
					},
					Spec: gatewayv1.HTTPRouteSpec{
						Hostnames: []gatewayv1.Hostname{"no-rule.example"},
					},
				},
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayv1.HTTPRouteSpec{
						Hostnames: []gatewayv1.Hostname{
							"foo.com",
						},
						Rules: []gatewayv1.HTTPRouteRule{
							{
								Matches: []gatewayv1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
								},
								BackendRefs: []gatewayv1.HTTPBackendRef{
									builder.NewHTTPBackendRef("service1").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedKongServices: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("httproute.default.httproute-1.foo.com.0"),
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
				"httproute.default.httproute-1.foo.com.0": {
					{
						Route: kong.Route{
							Name:         kong.String("httproute.default.httproute-1.foo.com.0.0"),
							Expression:   kong.String(`(http.host == "foo.com") && (http.path == "/v1/foo")`),
							PreserveHost: kong.Bool(true),
						},
						Plugins:          []kong.Plugin{},
						ExpressionRoutes: true,
					},
				},
			},
			expectedFailures: []failures.ResourceFailure{
				newResourceFailure(translators.ErrRouteValidationNoRules.Error(), &gatewayv1.HTTPRoute{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-no-host-no-rule",
					},
				}),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			failureCollector, err := failures.NewResourceFailuresCollector(logrus.New())
			require.NoError(t, err)
			parser.failuresCollector = failureCollector

			result := newIngressRules()
			parser.ingressRulesFromHTTPRoutesUsingExpressionRoutes(tc.httpRoutes, &result)
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

func TestIngressRulesFromSplitHTTPRouteMatchWithPriority(t *testing.T) {
	fakestore, err := store.NewFakeStore(store.FakeObjects{})
	require.NoError(t, err)
	parser := mustNewParser(t, fakestore)
	parser.featureFlags.CombinedServiceRoutes = true
	parser.featureFlags.ExpressionRoutes = true
	httpRouteTypeMeta := metav1.TypeMeta{Kind: "HTTPRoute", APIVersion: gatewayv1.SchemeGroupVersion.String()}

	testCases := []struct {
		name                string
		matchWithPriority   translators.SplitHTTPRouteMatchToKongRoutePriority
		expectedKongService kongstate.Service
		expectedKongRoute   kongstate.Route
		expectedError       error
	}{
		{
			name: "no hostname",
			matchWithPriority: translators.SplitHTTPRouteMatchToKongRoutePriority{
				Match: translators.SplitHTTPRouteMatch{
					Source: &gatewayv1.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "httproute-1",
						},
						Spec: gatewayv1.HTTPRouteSpec{
							Rules: []gatewayv1.HTTPRouteRule{
								{
									Matches: []gatewayv1.HTTPRouteMatch{
										builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
									},
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("service1").WithPort(80).Build(),
									},
								},
							},
						},
					},
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
				Priority: 1024,
			},
			expectedKongService: kongstate.Service{
				Service: kong.Service{
					Name: kong.String("httproute.default.httproute-1._.0"),
				},
				Backends: []kongstate.ServiceBackend{
					{
						Name:    "service1",
						PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
					},
				},
			},
			expectedKongRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("httproute.default.httproute-1._.0.0"),
					Expression:   kong.String(`http.path == "/v1/foo"`),
					PreserveHost: kong.Bool(true),
					StripPath:    kong.Bool(false),
					Priority:     kong.Uint64(1024),
				},
				Plugins:          []kong.Plugin{},
				ExpressionRoutes: true,
			},
		},
		{
			name: "precise hostname and filter",
			matchWithPriority: translators.SplitHTTPRouteMatchToKongRoutePriority{
				Match: translators.SplitHTTPRouteMatch{
					Source: &gatewayv1.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "httproute-1",
						},
						Spec: gatewayv1.HTTPRouteSpec{
							Hostnames: []gatewayv1.Hostname{
								"foo.com",
							},
							Rules: []gatewayv1.HTTPRouteRule{
								{
									Matches: []gatewayv1.HTTPRouteMatch{
										builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
										builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
									},
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("service1").WithPort(80).Build(),
									},
									Filters: []gatewayv1.HTTPRouteFilter{
										builder.NewHTTPRouteRequestRedirectFilter().
											WithRequestRedirectStatusCode(301).
											WithRequestRedirectHost("bar.com").
											Build(),
									},
								},
							},
						},
					},
					Hostname:   "foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
					RuleIndex:  0,
					MatchIndex: 1,
				},
				Priority: 1024,
			},
			expectedKongService: kongstate.Service{
				Service: kong.Service{
					Name: kong.String("httproute.default.httproute-1.foo.com.0"),
				},
				Backends: []kongstate.ServiceBackend{
					{
						Name:    "service1",
						PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
					},
				},
			},
			expectedKongRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("httproute.default.httproute-1.foo.com.0.1"),
					Expression:   kong.String(`(http.host == "foo.com") && (http.path == "/v1/foo")`),
					PreserveHost: kong.Bool(true),
					StripPath:    kong.Bool(false),
					Priority:     kong.Uint64(1024),
				},
				Plugins: []kong.Plugin{
					{
						Name: kong.String("request-termination"),
						Config: kong.Configuration{
							"status_code": kong.Int(301),
						},
					},
					{
						Name: kong.String("response-transformer"),
						Config: kong.Configuration{
							"add": map[string][]string{
								"headers": {"Location: http://bar.com:80/v1/foo"},
							},
						},
					},
				},
				ExpressionRoutes: true,
			},
		},
		{
			name: "wildcard hostname with multiple backends",
			matchWithPriority: translators.SplitHTTPRouteMatchToKongRoutePriority{
				Match: translators.SplitHTTPRouteMatch{
					Source: &gatewayv1.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "httproute-1",
						},
						Spec: gatewayv1.HTTPRouteSpec{
							Hostnames: []gatewayv1.Hostname{
								"*.foo.com",
							},
							Rules: []gatewayv1.HTTPRouteRule{
								{
									Matches: []gatewayv1.HTTPRouteMatch{
										builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
									},
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("service1").WithPort(80).WithWeight(10).Build(),
										builder.NewHTTPBackendRef("service2").WithPort(80).WithWeight(20).Build(),
									},
								},
							},
						},
					},
					Hostname:   "*.foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/v1/foo").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
				Priority: 1024,
			},
			expectedKongService: kongstate.Service{
				Service: kong.Service{
					Name: kong.String("httproute.default.httproute-1._.foo.com.0"),
				},
				Backends: []kongstate.ServiceBackend{
					{
						Name:    "service1",
						PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						Weight:  lo.ToPtr(int32(10)),
					},
					{
						Name:    "service2",
						PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
						Weight:  lo.ToPtr(int32(20)),
					},
				},
			},
			expectedKongRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("httproute.default.httproute-1._.foo.com.0.0"),
					Expression:   kong.String(`(http.host =^ ".foo.com") && (http.path == "/v1/foo")`),
					PreserveHost: kong.Bool(true),
					StripPath:    kong.Bool(false),
					Priority:     kong.Uint64(1024),
				},
				Plugins:          []kong.Plugin{},
				ExpressionRoutes: true,
			},
		},
		{
			name: "precise hostname and no match",
			matchWithPriority: translators.SplitHTTPRouteMatchToKongRoutePriority{
				Match: translators.SplitHTTPRouteMatch{
					Source: &gatewayv1.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "httproute-1",
						},
						Spec: gatewayv1.HTTPRouteSpec{
							Hostnames: []gatewayv1.Hostname{
								"a.foo.com",
							},
							Rules: []gatewayv1.HTTPRouteRule{
								{
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("service1").WithPort(80).Build(),
									},
								},
							},
						},
					},
					Hostname:   "a.foo.com",
					Match:      builder.NewHTTPRouteMatch().Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
				Priority: 1024,
			},
			expectedKongService: kongstate.Service{
				Service: kong.Service{
					Name: kong.String("httproute.default.httproute-1.a.foo.com.0"),
				},
				Backends: []kongstate.ServiceBackend{
					{
						Name:    "service1",
						PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
					},
				},
			},
			expectedKongRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("httproute.default.httproute-1.a.foo.com.0.0"),
					Expression:   kong.String(`http.host == "a.foo.com"`),
					PreserveHost: kong.Bool(true),
					StripPath:    kong.Bool(false),
					Priority:     kong.Uint64(1024),
				},
				Plugins:          []kong.Plugin{},
				ExpressionRoutes: true,
			},
		},
		{
			name: "no hostname and no match",
			matchWithPriority: translators.SplitHTTPRouteMatchToKongRoutePriority{
				Match: translators.SplitHTTPRouteMatch{
					Source: &gatewayv1.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "httproute-1",
						},
						Spec: gatewayv1.HTTPRouteSpec{
							Rules: []gatewayv1.HTTPRouteRule{
								{
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("service1").WithPort(80).Build(),
									},
								},
							},
						},
					},
					Hostname:   "",
					Match:      builder.NewHTTPRouteMatch().Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
				Priority: 1024,
			},
			expectedKongService: kongstate.Service{
				Service: kong.Service{
					Name: kong.String("httproute.default.httproute-1._.0"),
				},
				Backends: []kongstate.ServiceBackend{
					{
						Name:    "service1",
						PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(80)},
					},
				},
			},
			expectedKongRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("httproute.default.httproute-1._.0.0"),
					Expression:   kong.String(translators.CatchAllHTTPExpression),
					PreserveHost: kong.Bool(true),
					StripPath:    kong.Bool(false),
					Priority:     kong.Uint64(1024),
				},
				Plugins:          []kong.Plugin{},
				ExpressionRoutes: true,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			match := tc.matchWithPriority.Match
			tc.expectedKongRoute.Tags = util.GenerateTagsForObject(match.Source)
			tc.expectedKongRoute.Ingress = util.FromK8sObject(match.Source)

			res := newIngressRules()
			err := parser.ingressRulesFromSplitHTTPRouteMatchWithPriority(&res, tc.matchWithPriority)
			if tc.expectedError != nil {
				require.ErrorAs(t, err, tc.expectedError)
				return
			}
			require.NoError(t, err)
			kongService, ok := res.ServiceNameToServices[*tc.expectedKongService.Name]
			require.Truef(t, ok, "should find service %s", *tc.expectedKongService.Name)
			require.Equal(t, tc.expectedKongService.Backends, kongService.Backends)
			require.Len(t, kongService.Routes, 1)
			require.Equal(t, tc.expectedKongRoute, kongService.Routes[0])
		})
	}
}

func k8sObjectInfoOfHTTPRoute(route *gatewayv1.HTTPRoute) util.K8sObjectInfo {
	anotations := route.Annotations

	return util.K8sObjectInfo{
		Name:        route.Name,
		Namespace:   route.Namespace,
		Annotations: anotations,
		GroupVersionKind: schema.GroupVersionKind{
			Group:   "gateway.networking.k8s.io",
			Version: "v1beta1",
			Kind:    "HTTPRoute",
		},
	}
}

func commonRouteSpecMock(parentReferentName string) gatewayv1.CommonRouteSpec {
	return gatewayv1.CommonRouteSpec{
		ParentRefs: []gatewayv1.ParentReference{{
			Name: gatewayv1.ObjectName(parentReferentName),
		}},
	}
}
