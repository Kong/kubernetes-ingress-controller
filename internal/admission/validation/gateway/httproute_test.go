package gateway

import (
	"context"
	"fmt"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
)

func TestValidateHTTPRoute(t *testing.T) {
	var (
		nonexistentListener = gatewayv1.SectionName("listener-that-doesnt-exist")
		group               = gatewayv1.Group("gateway.networking.k8s.io")
		defaultGWNamespace  = gatewayv1.Namespace(corev1.NamespaceDefault)
		exampleGroup        = gatewayv1.Group("example")
		podKind             = gatewayv1.Kind("Pod")
	)

	for _, tt := range []struct {
		msg           string
		route         *gatewayv1.HTTPRoute
		gateways      []*gatewayv1.Gateway
		valid         bool
		validationMsg string
		err           error
	}{
		{
			msg: "if you provide errant gateways for validation, it fails validation",
			route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
			}, // no parentRefs
			gateways: []*gatewayv1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1.AllowedRoutes{
							Kinds: []gatewayv1.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "couldn't determine parentRefs for httproute",
			err:           fmt.Errorf("no parentRef matched gateway default/testing-gateway"),
		},
		{
			msg: "if you use sectionname to attach to a non-existent gateway listener, it fails validation",
			route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Name:        "testing-gateway",
							SectionName: &nonexistentListener,
						}},
					},
				},
			},
			gateways: []*gatewayv1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1.Listener{{
						Name:     "not-the-right-listener",
						Port:     80,
						Protocol: (gatewayv1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1.AllowedRoutes{
							Kinds: []gatewayv1.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "couldn't find gateway listeners for httproute",
			err:           fmt.Errorf("sectionname referenced listener listener-that-doesnt-exist was not found on gateway default/testing-gateway"),
		},
		{
			msg: "if the provided gateway has NO listeners, the HTTPRoute fails validation",
			route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
				},
			},
			gateways: []*gatewayv1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1.Listener{},
				},
			}},
			valid:         false,
			validationMsg: "couldn't find gateway listeners for httproute",
			err:           fmt.Errorf("no listeners could be found for gateway default/testing-gateway"),
		},
		{
			msg: "parentRefs which omit the namespace pass validation in the same namespace",
			route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
				},
			},
			gateways: []*gatewayv1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1.AllowedRoutes{
							Kinds: []gatewayv1.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid: true,
		},
		{
			msg: "if the gateway listener doesn't support HTTPRoute, validation fails",
			route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
				},
			},
			gateways: []*gatewayv1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1.Listener{{
						Name:     "http-alternate",
						Port:     8000,
						Protocol: (gatewayv1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1.AllowedRoutes{
							Kinds: []gatewayv1.RouteGroupKind{{
								Group: &group,
								Kind:  "TCPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "httproute linked gateway listeners did not pass validation",
			err:           fmt.Errorf("HTTPRoute not supported by listener http-alternate"),
		},
		{
			msg: "if an HTTPRoute is using queryparams matching it fails validation due to lack of support",
			route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayv1.HTTPRouteRule{{
						Matches: []gatewayv1.HTTPRouteMatch{{
							QueryParams: []gatewayv1.HTTPQueryParamMatch{{
								Name:  "user-agent",
								Value: "netscape navigator",
							}},
						}},
						BackendRefs: []gatewayv1.HTTPBackendRef{{
							BackendRef: gatewayv1.BackendRef{
								BackendObjectReference: gatewayv1.BackendObjectReference{
									Namespace: &defaultGWNamespace,
								},
							},
						}},
					}},
				},
			},
			gateways: []*gatewayv1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1.AllowedRoutes{
							Kinds: []gatewayv1.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "httproute spec did not pass validation",
			err:           fmt.Errorf("queryparam matching is not yet supported for httproute"),
		},
		{
			msg: "we don't support any group except core kubernetes for backendRefs",
			route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayv1.HTTPRouteRule{{
						Matches: []gatewayv1.HTTPRouteMatch{{
							Headers: []gatewayv1.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayv1.HTTPBackendRef{
							{
								BackendRef: gatewayv1.BackendRef{
									BackendObjectReference: gatewayv1.BackendObjectReference{
										Group:     &exampleGroup,
										Kind:      &podKind,
										Namespace: &defaultGWNamespace,
										Name:      "service1",
									},
								},
							},
						},
					}},
				},
			},
			gateways: []*gatewayv1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1.AllowedRoutes{
							Kinds: []gatewayv1.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "httproute spec did not pass validation",
			err:           fmt.Errorf("example is not a supported group for httproute backendRefs, only core is supported"),
		},
		{
			msg: "we don't support any core kind except Service for backendRefs",
			route: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayv1.HTTPRouteRule{{
						Matches: []gatewayv1.HTTPRouteMatch{{
							Headers: []gatewayv1.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayv1.HTTPBackendRef{
							{
								BackendRef: gatewayv1.BackendRef{
									BackendObjectReference: gatewayv1.BackendObjectReference{
										Kind:      &podKind,
										Namespace: &defaultGWNamespace,
										Name:      "service1",
									},
								},
							},
						},
					}},
				},
			},
			gateways: []*gatewayv1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1.AllowedRoutes{
							Kinds: []gatewayv1.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "httproute spec did not pass validation",
			err:           fmt.Errorf("Pod is not a supported kind for httproute backendRefs, only Service is supported"),
		},
	} {
		// Passed Kong version and routesValidator are irrelevant for the above test cases.
		valid, validMsg, err := ValidateHTTPRoute(
			context.Background(), mockRoutesValidator{}, parser.FeatureFlags{}, semver.MustParse("3.0.0"), tt.route, tt.gateways...,
		)
		assert.Equal(t, tt.valid, valid, tt.msg)
		assert.Equal(t, tt.validationMsg, validMsg, tt.msg)
		assert.Equal(t, tt.err, err, tt.msg)
	}
}

type mockRoutesValidator struct{}

func (mockRoutesValidator) Validate(_ context.Context, _ *kong.Route) (bool, string, error) {
	return true, "", nil
}
