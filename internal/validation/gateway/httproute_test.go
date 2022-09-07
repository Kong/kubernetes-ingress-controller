package gateway

import (
	"fmt"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func TestValidateHTTPRoute(t *testing.T) {
	var (
		nonexistentListener = gatewayv1beta1.SectionName("listener-that-doesnt-exist")
		group               = gatewayv1beta1.Group("gateway.networking.k8s.io")
		defaultGWNamespace  = gatewayv1beta1.Namespace(corev1.NamespaceDefault)
		pathMatchRegex      = gatewayv1beta1.PathMatchRegularExpression
		headerMatchRegex    = gatewayv1beta1.HeaderMatchRegularExpression
		exampleGroup        = gatewayv1beta1.Group("example")
		podKind             = gatewayv1beta1.Kind("Pod")
	)

	for _, tt := range []struct {
		msg           string
		route         *gatewayv1beta1.HTTPRoute
		gateways      []*gatewayv1beta1.Gateway
		valid         bool
		validationMsg string
		err           error
	}{
		{
			msg: "if you provide errant gateways for validation, it fails validation",
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
			}, // no parentRefs
			gateways: []*gatewayv1beta1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1beta1.GatewaySpec{
					Listeners: []gatewayv1beta1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1beta1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
							Kinds: []gatewayv1beta1.RouteGroupKind{{
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
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name:        "testing-gateway",
							SectionName: &nonexistentListener,
						}},
					},
				},
			},
			gateways: []*gatewayv1beta1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1beta1.GatewaySpec{
					Listeners: []gatewayv1beta1.Listener{{
						Name:     "not-the-right-listener",
						Port:     80,
						Protocol: (gatewayv1beta1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
							Kinds: []gatewayv1beta1.RouteGroupKind{{
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
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
				},
			},
			gateways: []*gatewayv1beta1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1beta1.GatewaySpec{
					Listeners: []gatewayv1beta1.Listener{},
				},
			}},
			valid:         false,
			validationMsg: "couldn't find gateway listeners for httproute",
			err:           fmt.Errorf("no listeners could be found for gateway default/testing-gateway"),
		},
		{
			msg: "parentRefs which omit the namespace pass validation in the same namespace",
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
				},
			},
			gateways: []*gatewayv1beta1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1beta1.GatewaySpec{
					Listeners: []gatewayv1beta1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1beta1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
							Kinds: []gatewayv1beta1.RouteGroupKind{{
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
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
				},
			},
			gateways: []*gatewayv1beta1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1beta1.GatewaySpec{
					Listeners: []gatewayv1beta1.Listener{{
						Name:     "http-alternate",
						Port:     8000,
						Protocol: (gatewayv1beta1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
							Kinds: []gatewayv1beta1.RouteGroupKind{{
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
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							QueryParams: []gatewayv1beta1.HTTPQueryParamMatch{{
								Name:  "user-agent",
								Value: "netscape navigator",
							}},
						}},
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Namespace: &defaultGWNamespace,
								},
							},
						}},
					}},
				},
			},
			gateways: []*gatewayv1beta1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1beta1.GatewaySpec{
					Listeners: []gatewayv1beta1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1beta1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
							Kinds: []gatewayv1beta1.RouteGroupKind{{
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
			msg: "if an HTTPRoute is using regex path matching it fails validation due to lack of support",
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							Path: &gatewayv1beta1.HTTPPathMatch{
								Type:  &pathMatchRegex,
								Value: kong.String("^path/to/stuff/*$"),
							},
						}},
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Namespace: &defaultGWNamespace,
								},
							},
						}},
					}},
				},
			},
			gateways: []*gatewayv1beta1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1beta1.GatewaySpec{
					Listeners: []gatewayv1beta1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1beta1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
							Kinds: []gatewayv1beta1.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "httproute spec did not pass validation",
			err:           fmt.Errorf("regex path matching is not yet supported for httproute"),
		},
		{
			msg: "if an HTTPRoute is using regex header matching it fails validation due to lack of support",
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							Headers: []gatewayv1beta1.HTTPHeaderMatch{{
								Type:  &headerMatchRegex,
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Namespace: &defaultGWNamespace,
								},
							},
						}},
					}},
				},
			},
			gateways: []*gatewayv1beta1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1beta1.GatewaySpec{
					Listeners: []gatewayv1beta1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1beta1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
							Kinds: []gatewayv1beta1.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "httproute spec did not pass validation",
			err:           fmt.Errorf("regex header matching is not yet supported for httproute"),
		},
		{
			msg: "we don't support any group except core kubernetes for backendRefs",
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							Headers: []gatewayv1beta1.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{
							{
								BackendRef: gatewayv1beta1.BackendRef{
									BackendObjectReference: gatewayv1beta1.BackendObjectReference{
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
			gateways: []*gatewayv1beta1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1beta1.GatewaySpec{
					Listeners: []gatewayv1beta1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1beta1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
							Kinds: []gatewayv1beta1.RouteGroupKind{{
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
			route: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
						ParentRefs: []gatewayv1beta1.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{{
						Matches: []gatewayv1beta1.HTTPRouteMatch{{
							Headers: []gatewayv1beta1.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayv1beta1.HTTPBackendRef{
							{
								BackendRef: gatewayv1beta1.BackendRef{
									BackendObjectReference: gatewayv1beta1.BackendObjectReference{
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
			gateways: []*gatewayv1beta1.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayv1beta1.GatewaySpec{
					Listeners: []gatewayv1beta1.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayv1beta1.HTTPProtocolType),
						AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
							Kinds: []gatewayv1beta1.RouteGroupKind{{
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
		valid, validMsg, err := ValidateHTTPRoute(tt.route, tt.gateways...)
		assert.Equal(t, tt.valid, valid, tt.msg)
		assert.Equal(t, tt.validationMsg, validMsg, tt.msg)
		assert.Equal(t, tt.err, err, tt.msg)
	}
}
