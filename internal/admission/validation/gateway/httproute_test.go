package gateway

import (
	"context"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/scheme"
)

func TestValidateHTTPRoute(t *testing.T) {
	var (
		group              = gatewayapi.Group("gateway.networking.k8s.io")
		defaultGWNamespace = gatewayapi.Namespace(corev1.NamespaceDefault)
		exampleGroup       = gatewayapi.Group("example")
		podKind            = gatewayapi.Kind("Pod")
		gatewayClassName   = gatewayapi.ObjectName("kong")
		gatewayClass       = &gatewayapi.GatewayClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(gatewayClassName),
			},
			Spec: gatewayapi.GatewayClassSpec{
				ControllerName: gatewaycontroller.GetControllerName(),
			},
		}
	)

	for _, tt := range []struct {
		msg           string
		route         *gatewayapi.HTTPRoute
		cachedObjects []client.Object
		valid         bool
		validationMsg string
		err           error
	}{
		{
			msg: "route with no parentRef is accepted with no validations",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
			}, // no parentRefs
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid: true,
		},
		{
			msg: "route with parentRef to non-gateway object is accepted with no vlaidation",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
								Namespace: lo.ToPtr(gatewayapi.Namespace(corev1.NamespaceDefault)),
								Name:      gatewayapi.ObjectName("kuma-cp"),
							},
						},
					},
				},
			}, // parentRef to a Service
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid: true,
		},
		{
			msg: "parentRefs which omit the namespace pass validation in the same namespace",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
				},
			},
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid: true,
		},
		{
			msg: "if an HTTPRoute is using queryparams matching it fails validation due to only supporting expression router",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{{
							QueryParams: []gatewayapi.HTTPQueryParamMatch{{
								Name:  "user-agent",
								Value: "netscape navigator",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Namespace: &defaultGWNamespace,
								},
							},
						}},
					}},
				},
			},
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid:         false,
			validationMsg: "HTTPRoute spec did not pass validation: rules[0].matches[0]: queryparam matching is supported with expression router only",
		},
		{
			msg: "we don't support any group except core kubernetes for backendRefs",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{{
							Headers: []gatewayapi.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
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
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid:         false,
			validationMsg: "HTTPRoute spec did not pass validation: rules[0].backendRefs[0]: example is not a supported group for httproute backendRefs, only core is supported",
		},
		{
			msg: "we don't support any core kind except Service for backendRefs",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{{
							Headers: []gatewayapi.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
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
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid:         false,
			validationMsg: "HTTPRoute spec did not pass validation: rules[0].backendRefs[0]: Pod is not a supported kind for httproute backendRefs, only Service is supported",
		},
		{
			msg: "we do not support RequestMirror filter",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{{
							Headers: []gatewayapi.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "service1",
									},
								},
							},
						},
						Filters: []gatewayapi.HTTPRouteFilter{
							{
								Type: gatewayapi.HTTPRouteFilterRequestMirror,
							},
						},
					}},
				},
			},
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid:         false,
			validationMsg: "HTTPRoute spec did not pass validation: rules[0].filters[0]: filter type RequestMirror is unsupported",
		},
		{
			msg: "we only support setting the timeout to the same value",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{{
							Headers: []gatewayapi.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "service1",
									},
								},
							},
						},
						Timeouts: &gatewayapi.HTTPRouteTimeouts{
							BackendRequest: lo.ToPtr(gatewayapi.Duration("1s")),
						},
					}, {
						Matches: []gatewayapi.HTTPRouteMatch{{
							Headers: []gatewayapi.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "service1",
									},
								},
							},
						},
						Timeouts: &gatewayapi.HTTPRouteTimeouts{
							BackendRequest: lo.ToPtr(gatewayapi.Duration("1s")),
						},
					}},
				},
			},
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid: true,
		},
		{
			msg: "we don't support setting the timeout to different value",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{{
							Headers: []gatewayapi.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "service1",
									},
								},
							},
						},
						Timeouts: &gatewayapi.HTTPRouteTimeouts{
							BackendRequest: lo.ToPtr(gatewayapi.Duration("1s")),
						},
					}, {
						Matches: []gatewayapi.HTTPRouteMatch{{
							Headers: []gatewayapi.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "service1",
									},
								},
							},
						},
						Timeouts: &gatewayapi.HTTPRouteTimeouts{
							BackendRequest: lo.ToPtr(gatewayapi.Duration("5s")),
						},
					}},
				},
			},
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid:         false,
			validationMsg: "HTTPRoute spec did not pass validation: timeout is set for one of the rules, but a different value is set in another rule",
		},
		{
			msg: "we do not support filters in backendRefs",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{{
							Headers: []gatewayapi.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "service1",
									},
								},
								Filters: []gatewayapi.HTTPRouteFilter{
									{
										Type: gatewayapi.HTTPRouteFilterRequestHeaderModifier,
									},
								},
							},
						},
					}},
				},
			},
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid:         false,
			validationMsg: "HTTPRoute spec did not pass validation: rules[0].backendRefs[0]: filters in backendRef is unsupported",
		},
		{
			msg: "invalid protocols",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.ProtocolsKey: "ohno",
					},
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{},
				},
			},
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid:         false,
			validationMsg: "HTTPRoute has invalid Kong annotations: invalid konghq.com/protocols value: ohno",
		},
		{
			msg: "HTTPRoute URLRewrite ReplaceFullPath",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Filters: []gatewayapi.HTTPRouteFilter{
								{
									Type: gatewayapi.HTTPRouteFilterURLRewrite,
									URLRewrite: &gatewayapi.HTTPURLRewriteFilter{
										Path: &gatewayapi.HTTPPathModifier{
											Type:            gatewayapi.FullPathHTTPPathModifier,
											ReplaceFullPath: lo.ToPtr("/new-path"),
										},
									},
								},
							},
						},
					},
				},
			},
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid: true,
		},
		{
			msg: "HTTPRoute URLRewrite Hostname",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Filters: []gatewayapi.HTTPRouteFilter{
								{
									Type: gatewayapi.HTTPRouteFilterURLRewrite,
									URLRewrite: &gatewayapi.HTTPURLRewriteFilter{
										Hostname: lo.ToPtr(gatewayapi.PreciseHostname("example.com")),
									},
								},
							},
						},
					},
				},
			},
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid:         false,
			validationMsg: "HTTPRoute spec did not pass validation: rules[0].filters[0]: filter type URLRewrite (with hostname replace) is unsupported",
		},
		{
			msg: "HTTPRoute URLRewrite ReplacePrefixMatch",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Filters: []gatewayapi.HTTPRouteFilter{
								{
									Type: gatewayapi.HTTPRouteFilterURLRewrite,
									URLRewrite: &gatewayapi.HTTPURLRewriteFilter{
										Path: &gatewayapi.HTTPPathModifier{
											Type:               gatewayapi.PrefixMatchHTTPPathModifier,
											ReplacePrefixMatch: lo.ToPtr("/new"),
										},
									},
								},
							},
						},
					},
				},
			},
			cachedObjects: []client.Object{
				gatewayClass,
				&gatewayapi.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: corev1.NamespaceDefault,
						Name:      "testing-gateway",
					},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayClassName,
						Listeners: []gatewayapi.Listener{{
							Name:     "http",
							Port:     80,
							Protocol: (gatewayapi.HTTPProtocolType),
							AllowedRoutes: &gatewayapi.AllowedRoutes{
								Kinds: []gatewayapi.RouteGroupKind{{
									Group: &group,
									Kind:  "HTTPRoute",
								}},
							},
						}},
					},
				},
			},
			valid: true,
		},
	} {
		t.Run(tt.msg, func(t *testing.T) {
			fakeClient := fakeclient.
				NewClientBuilder().
				WithScheme(lo.Must(scheme.Get())).
				WithObjects(tt.cachedObjects...).
				Build()

			// Passed routesValidator is irrelevant for the above test cases.
			valid, validMsg, err := ValidateHTTPRoute(
				context.Background(), mockRoutesValidator{}, translator.FeatureFlags{}, tt.route, fakeClient,
			)
			assert.Equal(t, tt.valid, valid, tt.msg)
			assert.Equal(t, tt.validationMsg, validMsg, tt.msg)
			assert.Equal(t, tt.err, err, tt.msg)
		})
	}
}

type mockRoutesValidator struct{}

func (mockRoutesValidator) Validate(_ context.Context, _ *kong.Route) (bool, string, error) {
	return true, "", nil
}
