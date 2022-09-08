package parser

import (
	"fmt"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var (
	// httprouteGVK is the GVK for HTTPRoutes, needed in unit tests because
	// we have to manually initialize objects that aren't retrieved from the
	// Kubernetes API.
	httprouteGVK = schema.GroupVersionKind{
		Group:   "gateway.networking.k8s.io",
		Version: "v1beta1",
		Kind:    "HTTPRoute",
	}

	pathMatchPrefix = gatewayv1beta1.PathMatchPathPrefix
	pathMatchRegex  = gatewayv1beta1.PathMatchRegularExpression
	pathMatchExact  = gatewayv1beta1.PathMatchExact
	queryMatchExact = gatewayv1beta1.QueryParamMatchExact
)

func Test_ingressRulesFromHTTPRoutes(t *testing.T) {
	fakestore, err := store.NewFakeStore(store.FakeObjects{})
	require.NoError(t, err)
	p := NewParser(logrus.New(), fakestore)
	httpPort := gatewayv1beta1.PortNumber(80)

	for _, tt := range []struct {
		msg      string
		routes   []*gatewayv1beta1.HTTPRoute
		expected ingressRules
		errs     []error
	}{
		{
			msg: "an empty list of HTTPRoutes should produce no ingress rules",
			expected: ingressRules{
				SecretNameToSNIs:      SecretNameToSNIs{},
				ServiceNameToServices: make(map[string]kongstate.Service),
			},
		},
		{
			msg: "an HTTPRoute rule with no matches can be routed if it has hostnames to match on",
			routes: []*gatewayv1beta1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: gatewayv1beta1.ObjectName("fake-gateway"),
						}},
					},
					Hostnames: []gatewayv1beta1.Hostname{
						"konghq.com",
						"www.konghq.com",
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName("fake-service"),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
							},
						}},
					}},
				},
			}},
			expected: ingressRules{
				SecretNameToSNIs: SecretNameToSNIs{},
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
						Backends: kongstate.ServiceBackends{{
							Name: "fake-service",
							PortDef: kongstate.PortDef{
								Mode:   kongstate.PortMode(1),
								Number: 80,
							},
						}},
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
							},
							Ingress: util.K8sObjectInfo{
								Name:        "basic-httproute",
								Namespace:   corev1.NamespaceDefault,
								Annotations: make(map[string]string),
								GroupVersionKind: schema.GroupVersionKind{
									Group:   "gateway.networking.k8s.io",
									Version: "v1beta1",
									Kind:    "HTTPRoute",
								},
							},
						}},
						Parent: &gatewayv1beta1.HTTPRoute{
							Spec: gatewayv1beta1.HTTPRouteSpec{
								CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
									ParentRefs: []gatewayv1beta1.ParentReference{
										{
											Name: gatewayv1beta1.ObjectName("fake-gateway"),
										},
									},
								},
								Hostnames: []gatewayv1beta1.Hostname{
									gatewayv1beta1.Hostname("konghq.com"),
									gatewayv1beta1.Hostname("www.konghq.com"),
								},
								Rules: []gatewayv1beta1.HTTPRouteRule{
									{
										BackendRefs: []gatewayv1beta1.HTTPBackendRef{
											{
												BackendRef: gatewayv1beta1.BackendRef{
													BackendObjectReference: gatewayv1beta1.BackendObjectReference{
														Name: gatewayv1beta1.ObjectName("fake-service"),
														Port: &httpPort,
														Kind: util.StringToGatewayAPIKindPtr("Service"),
													},
												},
											},
										},
									},
								},
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "basic-httproute",
								Namespace: "default",
							},
							TypeMeta: metav1.TypeMeta{
								Kind:       "HTTPRoute",
								APIVersion: "gateway.networking.k8s.io/v1beta1",
							},
						},
					},
				},
			},
		},
		{
			msg: "an HTTPRoute rule with no matches and no hostnames can't be routed",
			routes: []*gatewayv1beta1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: gatewayv1beta1.ObjectName("fake-gateway"),
						}},
					},
					// no hostnames present
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						// no match rules present
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName("fake-service"),
									Port: &httpPort,
								},
							},
						}},
					}},
				},
			}},
			expected: ingressRules{
				SecretNameToSNIs:      SecretNameToSNIs{},
				ServiceNameToServices: make(map[string]kongstate.Service),
			},
			errs: []error{
				fmt.Errorf("no match rules or hostnames specified"),
			},
		},
		{
			msg: "a single HTTPRoute with one match and one backendRef results in a single service",
			routes: []*gatewayv1beta1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: gatewayv1beta1.ObjectName("fake-gateway"),
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							Path: &gatewayv1beta1.HTTPPathMatch{
								Type:  &pathMatchPrefix,
								Value: kong.String("/httpbin"),
							},
						}},
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName("fake-service"),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
							},
						}},
					}},
				},
			}},
			expected: ingressRules{
				SecretNameToSNIs: SecretNameToSNIs{},
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
						Backends: kongstate.ServiceBackends{{
							Name: "fake-service",
							PortDef: kongstate.PortDef{
								Mode:   kongstate.PortMode(1),
								Number: 80,
							},
						}},
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
								StripPath: pointer.BoolPtr(false),
							},
							Ingress: util.K8sObjectInfo{
								Name:        "basic-httproute",
								Namespace:   corev1.NamespaceDefault,
								Annotations: make(map[string]string),
								GroupVersionKind: schema.GroupVersionKind{
									Group:   "gateway.networking.k8s.io",
									Version: "v1beta1",
									Kind:    "HTTPRoute",
								},
							},
						}},
						Parent: &gatewayv1beta1.HTTPRoute{
							Spec: gatewayv1beta1.HTTPRouteSpec{
								CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
									ParentRefs: []gatewayv1beta1.ParentReference{
										{
											Name: gatewayv1beta1.ObjectName("fake-gateway"),
										},
									},
								},
								Rules: []gatewayv1beta1.HTTPRouteRule{
									{
										Matches: []gatewayv1beta1.HTTPRouteMatch{
											{
												Path: &gatewayv1beta1.HTTPPathMatch{
													Type:  &pathMatchPrefix,
													Value: kong.String("/httpbin"),
												},
											},
										},
										BackendRefs: []gatewayv1beta1.HTTPBackendRef{
											{
												BackendRef: gatewayv1beta1.BackendRef{
													BackendObjectReference: gatewayv1beta1.BackendObjectReference{
														Name: gatewayv1beta1.ObjectName("fake-service"),
														Port: &httpPort,
														Kind: util.StringToGatewayAPIKindPtr("Service"),
													},
												},
											},
										},
									},
								},
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "basic-httproute",
								Namespace: "default",
							},
							TypeMeta: metav1.TypeMeta{
								Kind:       "HTTPRoute",
								APIVersion: "gateway.networking.k8s.io/v1beta1",
							},
						},
					},
				},
			},
		},
		{
			msg: "an HTTPRoute with no rules can't be routed",
			routes: []*gatewayv1beta1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: gatewayv1beta1.ObjectName("fake-gateway"),
						}},
					},
				},
			}},
			expected: ingressRules{
				SecretNameToSNIs:      SecretNameToSNIs{},
				ServiceNameToServices: make(map[string]kongstate.Service),
			},
			errs: []error{
				fmt.Errorf("no rules provided"),
			},
		},
		{
			msg: "an HTTPRoute with queryParam matches is not yet supported",
			routes: []*gatewayv1beta1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: gatewayv1beta1.ObjectName("fake-gateway"),
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							QueryParams: []gatewayv1beta1.HTTPQueryParamMatch{{
								Type:  &queryMatchExact,
								Name:  "username",
								Value: "kong",
							}},
						}},
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName("fake-service"),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
							},
						}},
					}},
				},
			}},
			expected: ingressRules{
				SecretNameToSNIs:      SecretNameToSNIs{},
				ServiceNameToServices: make(map[string]kongstate.Service),
			},
			errs: []error{
				fmt.Errorf("query param matches are not yet supported"),
			},
		},
		{
			msg: "an HTTPRoute with regex path matches is supported",
			routes: []*gatewayv1beta1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: gatewayv1beta1.ObjectName("fake-gateway"),
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							Path: &gatewayv1beta1.HTTPPathMatch{
								Type:  &pathMatchRegex,
								Value: kong.String("/httpbin$"),
							},
						}},
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName("fake-service"),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
							},
						}},
					}},
				},
			}},
			expected: ingressRules{
				SecretNameToSNIs: SecretNameToSNIs{},
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
						Backends: kongstate.ServiceBackends{{
							Name: "fake-service",
							PortDef: kongstate.PortDef{
								Mode:   kongstate.PortMode(1),
								Number: 80,
							},
						}},
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
								StripPath: pointer.BoolPtr(false),
							},
							Ingress: util.K8sObjectInfo{
								Name:        "basic-httproute",
								Namespace:   corev1.NamespaceDefault,
								Annotations: make(map[string]string),
								GroupVersionKind: schema.GroupVersionKind{
									Group:   "gateway.networking.k8s.io",
									Version: "v1beta1",
									Kind:    "HTTPRoute",
								},
							},
						}},
						Parent: &gatewayv1beta1.HTTPRoute{
							Spec: gatewayv1beta1.HTTPRouteSpec{
								CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
									ParentRefs: []gatewayv1beta1.ParentReference{
										{
											Name: gatewayv1beta1.ObjectName("fake-gateway"),
										},
									},
								},
								Rules: []gatewayv1beta1.HTTPRouteRule{
									{
										Matches: []gatewayv1beta1.HTTPRouteMatch{
											{
												Path: &gatewayv1beta1.HTTPPathMatch{
													Type:  &pathMatchRegex,
													Value: kong.String("/httpbin$"),
												},
											},
										},
										BackendRefs: []gatewayv1beta1.HTTPBackendRef{
											{
												BackendRef: gatewayv1beta1.BackendRef{
													BackendObjectReference: gatewayv1beta1.BackendObjectReference{
														Name: gatewayv1beta1.ObjectName("fake-service"),
														Port: &httpPort,
														Kind: util.StringToGatewayAPIKindPtr("Service"),
													},
												},
											},
										},
									},
								},
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "basic-httproute",
								Namespace: "default",
							},
							TypeMeta: metav1.TypeMeta{
								Kind:       "HTTPRoute",
								APIVersion: "gateway.networking.k8s.io/v1beta1",
							},
						},
					},
				},
			},
		},
		{
			msg: "an HTTPRoute with exact path matches translates to a terminated Kong regex route",
			routes: []*gatewayv1beta1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: gatewayv1beta1.ObjectName("fake-gateway"),
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							Path: &gatewayv1beta1.HTTPPathMatch{
								Type:  &pathMatchExact,
								Value: kong.String("/httpbin"),
							},
						}},
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName("fake-service"),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
							},
						}},
					}},
				},
			}},
			expected: ingressRules{
				SecretNameToSNIs: SecretNameToSNIs{},
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
						Backends: kongstate.ServiceBackends{{
							Name: "fake-service",
							PortDef: kongstate.PortDef{
								Mode:   kongstate.PortMode(1),
								Number: 80,
							},
						}},
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
								StripPath: pointer.BoolPtr(false),
							},
							Ingress: util.K8sObjectInfo{
								Name:        "basic-httproute",
								Namespace:   corev1.NamespaceDefault,
								Annotations: make(map[string]string),
								GroupVersionKind: schema.GroupVersionKind{
									Group:   "gateway.networking.k8s.io",
									Version: "v1beta1",
									Kind:    "HTTPRoute",
								},
							},
						}},
						Parent: &gatewayv1beta1.HTTPRoute{
							Spec: gatewayv1beta1.HTTPRouteSpec{
								CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
									ParentRefs: []gatewayv1beta1.ParentReference{
										{
											Name: gatewayv1beta1.ObjectName("fake-gateway"),
										},
									},
								},
								Rules: []gatewayv1beta1.HTTPRouteRule{
									{
										Matches: []gatewayv1beta1.HTTPRouteMatch{
											{
												Path: &gatewayv1beta1.HTTPPathMatch{
													Type:  &pathMatchExact,
													Value: kong.String("/httpbin"),
												},
											},
										},
										BackendRefs: []gatewayv1beta1.HTTPBackendRef{
											{
												BackendRef: gatewayv1beta1.BackendRef{
													BackendObjectReference: gatewayv1beta1.BackendObjectReference{
														Name: gatewayv1beta1.ObjectName("fake-service"),
														Port: &httpPort,
														Kind: util.StringToGatewayAPIKindPtr("Service"),
													},
												},
											},
										},
									},
								},
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "basic-httproute",
								Namespace: "default",
							},
							TypeMeta: metav1.TypeMeta{
								Kind:       "HTTPRoute",
								APIVersion: "gateway.networking.k8s.io/v1beta1",
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tt.msg, func(t *testing.T) {
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
			assert.Equal(t, tt.expected, ingressRules)

			// verify that we receive any and all expected errors
			assert.Equal(t, tt.errs, errs)
		})
	}
}

func Test_getHTTPRouteHostnamesAsSliceOfStringPointers(t *testing.T) {
	for _, tt := range []struct {
		msg      string
		input    *gatewayv1beta1.HTTPRoute
		expected []*string
	}{
		{
			msg:      "an HTTPRoute with no hostnames produces no hostnames",
			input:    &gatewayv1beta1.HTTPRoute{},
			expected: []*string{},
		},
		{
			msg: "an HTTPRoute with a single hostname produces a list with that one hostname",
			input: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
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
			input: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
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

func Test_ingressRulesFromHTTPRoutes_RegexPrefix(t *testing.T) {
	fakestore, err := store.NewFakeStore(store.FakeObjects{})
	require.NoError(t, err)
	p := NewParser(logrus.New(), fakestore)
	p.EnableRegexPathPrefix()
	httpPort := gatewayv1beta1.PortNumber(80)

	for _, tt := range []struct {
		msg      string
		routes   []*gatewayv1beta1.HTTPRoute
		expected ingressRules
		errs     []error
	}{
		{
			msg: "an HTTPRoute with regex path matches is supported",
			routes: []*gatewayv1beta1.HTTPRoute{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: gatewayv1beta1.ObjectName("fake-gateway"),
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							Path: &gatewayv1beta1.HTTPPathMatch{
								Type:  &pathMatchRegex,
								Value: kong.String("/httpbin$"),
							},
						}},
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName("fake-service"),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
							},
						}},
					}},
				},
			}},
			expected: ingressRules{
				SecretNameToSNIs: SecretNameToSNIs{},
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
						Backends: kongstate.ServiceBackends{{
							Name: "fake-service",
							PortDef: kongstate.PortDef{
								Mode:   kongstate.PortMode(1),
								Number: 80,
							},
						}},
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
								StripPath: pointer.BoolPtr(false),
							},
							Ingress: util.K8sObjectInfo{
								Name:        "basic-httproute",
								Namespace:   corev1.NamespaceDefault,
								Annotations: make(map[string]string),
								GroupVersionKind: schema.GroupVersionKind{
									Group:   "gateway.networking.k8s.io",
									Version: "v1beta1",
									Kind:    "HTTPRoute",
								},
							},
						}},
						Parent: &gatewayv1beta1.HTTPRoute{
							Spec: gatewayv1beta1.HTTPRouteSpec{
								CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
									ParentRefs: []gatewayv1beta1.ParentReference{
										{
											Name: gatewayv1beta1.ObjectName("fake-gateway"),
										},
									},
								},
								Rules: []gatewayv1beta1.HTTPRouteRule{
									{
										Matches: []gatewayv1beta1.HTTPRouteMatch{
											{
												Path: &gatewayv1beta1.HTTPPathMatch{
													Type:  &pathMatchRegex,
													Value: kong.String("/httpbin$"),
												},
											},
										},
										BackendRefs: []gatewayv1beta1.HTTPBackendRef{
											{
												BackendRef: gatewayv1beta1.BackendRef{
													BackendObjectReference: gatewayv1beta1.BackendObjectReference{
														Name: gatewayv1beta1.ObjectName("fake-service"),
														Port: &httpPort,
														Kind: util.StringToGatewayAPIKindPtr("Service"),
													},
												},
											},
										},
									},
								},
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "basic-httproute",
								Namespace: "default",
							},
							TypeMeta: metav1.TypeMeta{
								Kind:       "HTTPRoute",
								APIVersion: "gateway.networking.k8s.io/v1beta1",
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tt.msg, func(t *testing.T) {
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
			assert.Equal(t, tt.expected, ingressRules)

			// verify that we receive any and all expected errors
			assert.Equal(t, tt.errs, errs)
		})
	}
}
