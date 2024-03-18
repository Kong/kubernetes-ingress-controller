package gateway

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset/scheme"
)

func init() {
	if err := corev1.AddToScheme(scheme.Scheme); err != nil {
		fmt.Println("error while adding core1 scheme")
		os.Exit(1)
	}
	if err := gatewayv1.Install(scheme.Scheme); err != nil {
		fmt.Println("error while adding gatewayv1 scheme")
		os.Exit(1)
	}
}

func TestFilterHostnames(t *testing.T) {
	commonGateway := &gatewayapi.Gateway{
		Spec: gatewayapi.GatewaySpec{
			Listeners: []gatewayapi.Listener{
				{
					Name:     "listener-1",
					Hostname: util.StringToGatewayAPIHostnamePtr("very.specific.com"),
				},
				{
					Name:     "listener-2",
					Hostname: util.StringToGatewayAPIHostnamePtr("*.wildcard.io"),
				},
				{
					Name:     "listener-3",
					Hostname: util.StringToGatewayAPIHostnamePtr("*.anotherwildcard.io"),
				},
				{
					Name: "listener-4",
				},
			},
		},
	}

	testCases := []struct {
		name              string
		gateways          []supportedGatewayWithCondition
		httpRoute         *gatewayapi.HTTPRoute
		expectedHTTPRoute *gatewayapi.HTTPRoute
		hasError          bool
		errString         string
	}{
		{
			name: "listener 1 - specific",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listener-1",
				},
			},
			httpRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						util.StringToGatewayAPIHostname("*.anotherwildcard.io"),
						util.StringToGatewayAPIHostname("*.nonmatchingwildcard.io"),
						util.StringToGatewayAPIHostname("very.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						util.StringToGatewayAPIHostname("very.specific.com"),
					},
				},
			},
		},
		{
			name: "listener 1 - wildcard",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listener-1",
				},
			},
			httpRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						util.StringToGatewayAPIHostname("non.matching.com"),
						util.StringToGatewayAPIHostname("*.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						util.StringToGatewayAPIHostname("very.specific.com"),
					},
				},
			},
		},
		{
			name: "listener 2",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listener-2",
				},
			},
			httpRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						util.StringToGatewayAPIHostname("non.matching.com"),
						util.StringToGatewayAPIHostname("wildcard.io"),
						util.StringToGatewayAPIHostname("foo.wildcard.io"),
						util.StringToGatewayAPIHostname("bar.wildcard.io"),
						util.StringToGatewayAPIHostname("foo.bar.wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						util.StringToGatewayAPIHostname("foo.wildcard.io"),
						util.StringToGatewayAPIHostname("bar.wildcard.io"),
						util.StringToGatewayAPIHostname("foo.bar.wildcard.io"),
					},
				},
			},
		},
		{
			name: "listener 3 - wildcard",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listener-3",
				},
			},
			httpRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						util.StringToGatewayAPIHostname("*.anotherwildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						util.StringToGatewayAPIHostname("*.anotherwildcard.io"),
					},
				},
			},
		},
		{
			name: "no listner specified - no hostname",
			gateways: []supportedGatewayWithCondition{
				{
					gateway: commonGateway,
				},
			},
			httpRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{},
				},
			},
			expectedHTTPRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{},
				},
			},
		},
		{
			name: "listener 1 - no match",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listner-1",
				},
			},
			httpRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						util.StringToGatewayAPIHostname("specific.but.wrong.com"),
						util.StringToGatewayAPIHostname("wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{},
				},
			},
			hasError:  true,
			errString: "no matching hostnames in listener",
		},
	}

	for _, tc := range testCases {
		filteredHTTPRoute, err := filterHostnames(tc.gateways, tc.httpRoute)
		if tc.hasError {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errString)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedHTTPRoute.Spec, filteredHTTPRoute.Spec, tc.name)
		}
	}
}

func TestGetSupportedGatewayForRoute(t *testing.T) {
	gatewayClass := &gatewayapi.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-gatewayclass",
		},
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: gatewayapi.GatewayController("konghq.com/kic-gateway-controller"),
		},
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
	}

	routeConditionAccepted := func(status metav1.ConditionStatus, reason gatewayapi.RouteConditionReason) metav1.Condition {
		return metav1.Condition{
			Type:   string(gatewayapi.RouteConditionAccepted),
			Status: status,
			Reason: string(reason),
		}
	}

	type expected struct {
		condition    metav1.Condition
		listenerName string
	}

	goodGroup := gatewayapi.Group(gatewayv1.GroupName)
	goodKind := gatewayapi.Kind("Gateway")
	basicHTTPRoute := func() *gatewayapi.HTTPRoute {
		return &gatewayapi.HTTPRoute{
			TypeMeta: gatewayapi.V1GatewayTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name:      "basic-httproute",
				Namespace: "test-namespace",
			},
			Spec: gatewayapi.HTTPRouteSpec{
				CommonRouteSpec: gatewayapi.CommonRouteSpec{
					ParentRefs: []gatewayapi.ParentReference{
						{
							Group: &goodGroup,
							Kind:  &goodKind,
							Name:  "test-gateway",
						},
					},
				},
				Rules: []gatewayapi.HTTPRouteRule{
					{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					},
				},
			},
		}
	}
	t.Run("HTTPRoute", func(t *testing.T) {
		gatewayWithHTTP80Ready := func() *gatewayapi.Gateway {
			return &gatewayapi.Gateway{
				TypeMeta: gatewayapi.V1GatewayTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
					UID:       "ce7f0678-f59a-483c-80d1-243d3738d22c",
				},
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: "test-gatewayclass",
					Listeners:        builder.NewListener("http").WithPort(80).HTTP().IntoSlice(),
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: "http",
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
							SupportedKinds: supportedRouteGroupKinds,
						},
					},
				},
			}
		}

		tests := []struct {
			name     string
			route    *gatewayapi.HTTPRoute
			expected []expected
			objects  []client.Object
		}{
			{
				name:  "basic HTTPRoute gets accepted",
				route: basicHTTPRoute(),
				objects: []client.Object{
					gatewayWithHTTP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
					},
				},
			},
			{
				name:  "basic HTTPRoute with TLS configuration gets accepted",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Spec.Listeners = builder.
							NewListener("http").WithPort(443).HTTPS().
							WithTLSConfig(&gatewayapi.GatewayTLSConfig{
								Mode: lo.ToPtr(gatewayapi.TLSModeTerminate),
							}).
							IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
					},
				},
			},
			{
				name: "basic HTTPRoute specifying existing section name gets Accepted",
				route: func() *gatewayapi.HTTPRoute {
					r := basicHTTPRoute()
					r.Spec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("http"))
					return r
				}(),
				objects: []client.Object{
					gatewayWithHTTP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						listenerName: "http",
						condition:    routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
					},
				},
			},
			{
				name: "basic HTTPRoute specifying existing port gets Accepted",
				route: func() *gatewayapi.HTTPRoute {
					r := basicHTTPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(80))
					return r
				}(),
				objects: []client.Object{
					gatewayWithHTTP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
					},
				},
			},
			{
				name: "basic HTTPRoute specifying non-existing port does not get Accepted",
				route: func() *gatewayapi.HTTPRoute {
					r := basicHTTPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(80))
					return r
				}(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Spec.Listeners = builder.NewListener("http").WithPort(81).HTTP().IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
					},
				},
			},
			{
				name:  "basic HTTPRoute does not get accepted if it is not in the supported kinds",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Status.Listeners[0].SupportedKinds = nil
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNotAllowedByListeners),
					},
				},
			},
			{
				name:  "basic HTTPRoute does not get accepted if it is not permitted by allowed routes",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Spec.Listeners = builder.NewListener("http").
							WithPort(80).
							HTTP().
							WithAllowedRoutes(
								&gatewayapi.AllowedRoutes{
									Kinds: builder.NewRouteGroupKind().TCPRoute().IntoSlice(),
								},
							).
							IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNotAllowedByListeners),
					},
				},
			},
			{
				name:  "basic HTTPRoute does get accepted if allowed routes only specified Same namespace",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Spec.Listeners = builder.NewListener("http").
							WithPort(80).
							HTTP().
							WithAllowedRoutes(builder.NewAllowedRoutesFromSameNamespaces()).
							IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
					},
				},
			},
			{
				name: "HTTPRoute does not get accepted if Listener hostnames do not match route hostnames",
				route: func() *gatewayapi.HTTPRoute {
					r := basicHTTPRoute()
					r.Spec.Hostnames = []gatewayapi.Hostname{"very.specific.com"}
					return r
				}(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Spec.Listeners = builder.NewListener("http").
							WithPort(80).
							HTTP().
							WithAllowedRoutes(builder.NewAllowedRoutesFromSameNamespaces()).
							WithHostname("hostname.com").
							IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingListenerHostname),
					},
				},
			},
			{
				name:  "HTTPRoute does not get accepted if Listener TLSConfig uses PassThrough",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Spec.Listeners = builder.
							NewListener("https").WithPort(443).HTTPS().
							WithTLSConfig(&gatewayapi.GatewayTLSConfig{
								Mode: lo.ToPtr(gatewayapi.TLSModePassthrough),
							}).
							IntoSlice()
						gw.Status.Listeners[0].Name = "https"
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
					},
				},
			},
			{
				name:  "HTTPRoute does not get accepted if Listener doesn't match route's protocol",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithHTTP80Ready()
						// Using matching Listener name, but wrong protocol.
						gw.Spec.Listeners = builder.NewListener("http").WithPort(80).UDP().IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
					},
				},
			},
			{
				name:  "HTTPRoute does not get accepted if Listener doesn't match route's section name",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithHTTP80Ready()
						// Using Listener's name not matching Route's section name.
						gw.Spec.Listeners = builder.NewListener("not-http").WithPort(80).HTTP().IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNotAllowedByListeners),
					},
				},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				fakeClient := fakeclient.
					NewClientBuilder().
					WithScheme(scheme.Scheme).
					WithObjects(tt.objects...).
					Build()

				got, err := getSupportedGatewayForRoute(context.Background(), logr.Discard(), fakeClient, tt.route, controllers.OptionalNamespacedName{})
				require.NoError(t, err)
				require.Len(t, got, len(tt.expected))

				for i := range got {
					assert.Equalf(t, "test-namespace", got[i].gateway.Namespace, "gateway namespace #%d", i)
					assert.Equalf(t, "test-gateway", got[i].gateway.Name, "gateway name #%d", i)
					assert.Equalf(t, tt.expected[i].listenerName, got[i].listenerName, "listenerName #%d", i)
					assert.Equalf(t, tt.expected[i].condition, got[i].condition, "condition #%d", i)
				}
			})
		}
	})

	t.Run("TCPRoute", func(t *testing.T) {
		basicTCPRoute := func() *gatewayapi.TCPRoute {
			return &gatewayapi.TCPRoute{
				TypeMeta: gatewayapi.TCPRouteTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-tcproute",
					Namespace: "test-namespace",
				},
				Spec: gatewayapi.TCPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Group: &goodGroup,
								Kind:  &goodKind,
								Name:  "test-gateway",
							},
						},
					},
				},
			}
		}
		gatewayWithTCP80Ready := func() *gatewayapi.Gateway {
			return &gatewayapi.Gateway{
				TypeMeta: gatewayapi.V1GatewayTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
					UID:       "ce7f0678-f59a-483c-80d1-243d3738d22c",
				},
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: "test-gatewayclass",
					Listeners:        builder.NewListener("tcp").WithPort(80).TCP().IntoSlice(),
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: "tcp",
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
							SupportedKinds: builder.NewRouteGroupKind().TCPRoute().IntoSlice(),
						},
					},
				},
			}
		}

		type expected struct {
			condition    metav1.Condition
			listenerName string
		}
		tests := []struct {
			name     string
			route    *gatewayapi.TCPRoute
			expected expected
			objects  []client.Object
			wantErr  bool
		}{
			{
				name:  "basic TCPRoute does get accepted because it is in supported kinds",
				route: basicTCPRoute(),
				objects: []client.Object{
					gatewayWithTCP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
				},
			},
			{
				name:  "basic TCPRoute does not get accepted because it is not in supported kinds",
				route: basicTCPRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithTCP80Ready()
						gw.Status.Listeners[0].SupportedKinds = nil
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNotAllowedByListeners),
				},
			},
			{
				name: "TCPRoute specifying existing port gets Accepted",
				route: func() *gatewayapi.TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(80))
					r.Spec.Rules = []gatewayapi.TCPRouteRule{
						{
							BackendRefs: builder.NewBackendRef("fake-service").WithPort(80).ToSlice(),
						},
					}
					return r
				}(),
				objects: []client.Object{
					gatewayWithTCP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
				},
			},
			{
				name: "TCPRoute specifying non existing port does not get Accepted",
				route: func() *gatewayapi.TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(8000))
					r.Spec.Rules = []gatewayapi.TCPRouteRule{
						{
							BackendRefs: builder.NewBackendRef("fake-service").WithPort(80).ToSlice(),
						},
					}
					return r
				}(),
				objects: []client.Object{
					gatewayWithTCP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
				},
			},
			{
				name: "TCPRoute specifying in sectionName existing listener gets Accepted",
				route: func() *gatewayapi.TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(80))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("tcp"))
					r.Spec.Rules = []gatewayapi.TCPRouteRule{
						{
							BackendRefs: builder.NewBackendRef("fake-service").WithPort(80).ToSlice(),
						},
					}
					return r
				}(),
				objects: []client.Object{
					gatewayWithTCP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "tcp",
					condition:    routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
				},
			},
			{
				name: "TCPRoute specifying in sectionName non existing listener does not get Accepted",
				route: func() *gatewayapi.TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(80))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("unknown-listener"))
					r.Spec.Rules = []gatewayapi.TCPRouteRule{
						{
							BackendRefs: builder.NewBackendRef("fake-service").WithPort(80).ToSlice(),
						},
					}
					return r
				}(),
				objects: []client.Object{
					gatewayWithTCP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "unknown-listener",
					condition:    routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
				},
			},
			{
				name: "TCPRoute specifying in sectionName existing listener with a matching port gets Accepted",
				route: func() *gatewayapi.TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(80))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("tcp"))
					r.Spec.Rules = []gatewayapi.TCPRouteRule{
						{
							BackendRefs: builder.NewBackendRef("fake-service").WithPort(80).ToSlice(),
						},
					}
					return r
				}(),
				objects: []client.Object{
					gatewayWithTCP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "tcp",
					condition:    routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
				},
			},
			{
				name: "TCPRoute specifying in sectionName existing listener with a non-matching port does not gets Accepted",
				route: func() *gatewayapi.TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(8080))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("tcp"))
					r.Spec.Rules = []gatewayapi.TCPRouteRule{
						{
							BackendRefs: builder.NewBackendRef("fake-service").WithPort(80).ToSlice(),
						},
					}
					return r
				}(),
				objects: []client.Object{
					gatewayWithTCP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "tcp",
					condition:    routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
				},
			},
			{
				name: "TCPRoute specifying in sectionName non existing listener with an existing port does not gets Accepted",
				route: func() *gatewayapi.TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(80))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("unknown-listener"))
					r.Spec.Rules = []gatewayapi.TCPRouteRule{
						{
							BackendRefs: builder.NewBackendRef("fake-service").WithPort(80).ToSlice(),
						},
					}
					return r
				}(),
				objects: []client.Object{
					gatewayWithTCP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "unknown-listener",
					condition:    routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
				},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				fakeClient := fakeclient.
					NewClientBuilder().
					WithScheme(scheme.Scheme).
					WithObjects(tt.objects...).
					Build()

				got, err := getSupportedGatewayForRoute(context.Background(), logr.Discard(), fakeClient, tt.route, controllers.OptionalNamespacedName{})
				require.NoError(t, err)
				require.Len(t, got, 1)
				match := got[0]

				assert.Equal(t, "test-namespace", match.gateway.Namespace)
				assert.Equal(t, "test-gateway", match.gateway.Name)
				assert.Equal(t, tt.expected.listenerName, match.listenerName)
				assert.Equal(t, tt.expected.condition, match.condition)
			})
		}
	})

	t.Run("UDPRoute", func(t *testing.T) {
		basicUDPRoute := func() *gatewayapi.UDPRoute {
			return &gatewayapi.UDPRoute{
				TypeMeta: gatewayapi.UDPRouteTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-udproute",
					Namespace: "test-namespace",
				},
				Spec: gatewayapi.UDPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Group: &goodGroup,
								Kind:  &goodKind,
								Name:  "test-gateway",
							},
						},
					},
				},
			}
		}
		gatewayWithUDP53Ready := func() *gatewayapi.Gateway {
			return &gatewayapi.Gateway{
				TypeMeta: gatewayapi.V1GatewayTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
					UID:       "ce7f0678-f59a-483c-80d1-243d3738d22c",
				},
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: "test-gatewayclass",
					Listeners:        builder.NewListener("udp").WithPort(53).UDP().IntoSlice(),
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: "udp",
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
							SupportedKinds: builder.NewRouteGroupKind().UDPRoute().IntoSlice(),
						},
					},
				},
			}
		}

		type expected struct {
			condition    metav1.Condition
			listenerName string
		}
		tests := []struct {
			name     string
			route    *gatewayapi.UDPRoute
			expected expected
			objects  []client.Object
			wantErr  bool
		}{
			{
				name:  "basic UDPRoute does get accepted because it is in supported kinds",
				route: basicUDPRoute(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
				},
			},
			{
				name:  "basic UDPRoute does not get accepted because it is not in supported kinds",
				route: basicUDPRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithUDP53Ready()
						gw.Status.Listeners[0].SupportedKinds = nil
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNotAllowedByListeners),
				},
			},
			{
				name: "UDPRoute specifying existing port gets Accepted",
				route: func() *gatewayapi.UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(53))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
				},
			},
			{
				name: "UDPRoute specifying non existing port does not get Accepted",
				route: func() *gatewayapi.UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(8000))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
				},
			},
			{
				name: "UDPRoute specifying in sectionName existing listener gets Accepted",
				route: func() *gatewayapi.UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(53))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("udp"))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "udp",
					condition:    routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
				},
			},
			{
				name: "UDPRoute specifying in sectionName non existing listener does not get Accepted",
				route: func() *gatewayapi.UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(53))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("unknown-listener"))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "unknown-listener",
					condition:    routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
				},
			},
			{
				name: "UDPRoute specifying in sectionName existing listener with a matching port gets Accepted",
				route: func() *gatewayapi.UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(53))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("udp"))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "udp",
					condition:    routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
				},
			},
			{
				name: "UDPRoute specifying in sectionName existing listener with a non-matching port does not get Accepted",
				route: func() *gatewayapi.UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(533))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("udp"))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "udp",
					condition:    routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
				},
			},
			{
				name: "UDPRoute specifying in sectionName non existing listener with an existing port does not get Accepted",
				route: func() *gatewayapi.UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(53))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("unknown-listener"))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "unknown-listener",
					condition:    routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
				},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				fakeClient := fakeclient.
					NewClientBuilder().
					WithScheme(scheme.Scheme).
					WithObjects(tt.objects...).
					Build()

				got, err := getSupportedGatewayForRoute(context.Background(), logr.Discard(), fakeClient, tt.route, controllers.OptionalNamespacedName{})
				require.NoError(t, err)
				require.Len(t, got, 1)
				match := got[0]

				assert.Equal(t, "test-namespace", match.gateway.Namespace)
				assert.Equal(t, "test-gateway", match.gateway.Name)
				assert.Equal(t, tt.expected.listenerName, match.listenerName)
				assert.Equal(t, tt.expected.condition, match.condition)
			})
		}
	})

	t.Run("TLSRoute", func(t *testing.T) {
		basicTLSRoute := func() *gatewayapi.TLSRoute {
			return &gatewayapi.TLSRoute{
				TypeMeta: gatewayapi.TLSRouteTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-tlsroute",
					Namespace: "test-namespace",
				},
				Spec: gatewayapi.TLSRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Group: &goodGroup,
								Kind:  &goodKind,
								Name:  "test-gateway",
							},
						},
					},
				},
			}
		}
		gatewayWithTLS443PassthroughReady := func() *gatewayapi.Gateway {
			return &gatewayapi.Gateway{
				TypeMeta: gatewayapi.V1GatewayTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
					UID:       "ce7f0678-f59a-483c-80d1-243d3738d22c",
				},
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: "test-gatewayclass",
					Listeners: builder.NewListener("tls").
						WithPort(443).
						TLS().
						WithTLSConfig(&gatewayapi.GatewayTLSConfig{
							Mode: lo.ToPtr(gatewayapi.TLSModePassthrough),
						}).IntoSlice(),
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: "tls",
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
							SupportedKinds: builder.NewRouteGroupKind().TLSRoute().IntoSlice(),
						},
					},
				},
			}
		}

		type expected struct {
			condition    metav1.Condition
			listenerName string
		}
		tests := []struct {
			name     string
			route    *gatewayapi.TLSRoute
			expected []expected
			objects  []client.Object
			wantErr  bool
		}{
			{
				name:  "basic TLSRoute does get accepted because it is in supported kinds and there is a listener with TLS in passthrough mode",
				route: basicTLSRoute(),
				objects: []client.Object{
					gatewayWithTLS443PassthroughReady(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
					},
				},
			},
			{
				name:  "basic TLSRoute does not get accepted because there is no listener with TLS in passthrough mode",
				route: basicTLSRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithTLS443PassthroughReady()
						gw.Spec.Listeners = builder.NewListener("tls").
							WithPort(443).
							TLS().
							WithTLSConfig(&gatewayapi.GatewayTLSConfig{
								Mode: lo.ToPtr(gatewayapi.TLSModeTerminate),
							}).IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
					},
				},
			},
			{
				name:  "TLSRoute does not get accepted because it is not in supported kinds",
				route: basicTLSRoute(),
				objects: []client.Object{
					func() *gatewayapi.Gateway {
						gw := gatewayWithTLS443PassthroughReady()
						gw.Status.Listeners[0].SupportedKinds = nil
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNotAllowedByListeners),
					},
				},
			},
			{
				name: "TLSRoute specifying existing port gets Accepted",
				route: func() *gatewayapi.TLSRoute {
					r := basicTLSRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(443))
					return r
				}(),
				objects: []client.Object{
					gatewayWithTLS443PassthroughReady(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
					},
				},
			},
			{
				name: "TLSRoute specifying non existing port does not get Accepted",
				route: func() *gatewayapi.TLSRoute {
					r := basicTLSRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(444))
					return r
				}(),
				objects: []client.Object{
					gatewayWithTLS443PassthroughReady(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
					},
				},
			},
			{
				name: "TLSRoute specifying in sectionName existing listener gets Accepted",
				route: func() *gatewayapi.TLSRoute {
					r := basicTLSRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(443))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("tls"))
					return r
				}(),
				objects: []client.Object{
					gatewayWithTLS443PassthroughReady(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						listenerName: "tls",
						condition:    routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
					},
				},
			},
			// unmatched sectionName.
			{
				name: "TLSRoute specifying in sectionName non existing listener does not get Accepted",
				route: func() *gatewayapi.TLSRoute {
					r := basicTLSRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayapi.PortNumber(443))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayapi.SectionName("unknown-listener"))
					return r
				}(),
				objects: []client.Object{
					gatewayWithTLS443PassthroughReady(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						listenerName: "unknown-listener",
						condition:    routeConditionAccepted(metav1.ConditionFalse, gatewayapi.RouteReasonNoMatchingParent),
					},
				},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				fakeClient := fakeclient.
					NewClientBuilder().
					WithScheme(scheme.Scheme).
					WithObjects(tt.objects...).
					Build()

				got, err := getSupportedGatewayForRoute(context.Background(), logr.Discard(), fakeClient, tt.route, controllers.OptionalNamespacedName{})
				require.NoError(t, err)
				require.Len(t, got, len(tt.expected))

				for i := range got {
					assert.Equalf(t, "test-namespace", got[i].gateway.Namespace, "gateway namespace #%d", i)
					assert.Equalf(t, "test-gateway", got[i].gateway.Name, "gateway name #%d", i)
					assert.Equalf(t, tt.expected[i].listenerName, got[i].listenerName, "listenerName #%d", i)
					assert.Equalf(t, tt.expected[i].condition, got[i].condition, "condition #%d", i)
				}
			})
		}
	})

	bustedParentHTTPRoute := basicHTTPRoute()
	badGroup := gatewayapi.Group("kechavakunduz.cholpon.uz")
	badKind := gatewayapi.Kind("razzoq")
	bustedParentHTTPRoute.Spec.ParentRefs = []gatewayapi.ParentReference{
		{
			Name:  "not-a-gateway",
			Kind:  &badKind,
			Group: &badGroup,
		},
	}
	t.Run("invalid parentRef kind rejected", func(t *testing.T) {
		fakeClient := fakeclient.
			NewClientBuilder().
			WithScheme(scheme.Scheme).
			Build()

		_, err := getSupportedGatewayForRoute(context.Background(), logr.Discard(), fakeClient, bustedParentHTTPRoute, controllers.OptionalNamespacedName{})
		require.Equal(t, fmt.Errorf("unsupported parent kind %s/%s", string(badGroup), string(badKind)), err)
	})

	t.Run("single Gateway", func(t *testing.T) {
		namedGateway := func(name string) *gatewayapi.Gateway {
			return &gatewayapi.Gateway{
				TypeMeta: gatewayapi.V1GatewayTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: "test-namespace",
					UID:       "ce7f0678-f59a-483c-80d1-243d3738d22c",
				},
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: "test-gatewayclass",
					Listeners:        builder.NewListener("http").WithPort(80).HTTP().IntoSlice(),
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: "http",
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
							SupportedKinds: supportedRouteGroupKinds,
						},
					},
				},
			}
		}

		basicHTTPRoute := func(gateway string) *gatewayapi.HTTPRoute {
			return &gatewayapi.HTTPRoute{
				TypeMeta: gatewayapi.V1GatewayTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-httproute",
					Namespace: "test-namespace",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Group: &goodGroup,
								Kind:  &goodKind,
								Name:  gatewayapi.ObjectName(gateway),
							},
						},
					},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
							},
						},
					},
				},
			}
		}

		tests := []struct {
			name        string
			route       *gatewayapi.HTTPRoute
			expected    []expected
			expectedErr error
			objects     []client.Object
		}{
			{
				name:  "HTTPRoute with bound Gateway parent is accepted",
				route: basicHTTPRoute("good-gateway"),
				objects: []client.Object{
					namedGateway("good-gateway"),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayapi.RouteReasonAccepted),
					},
				},
			},
			{
				name:  "HTTPRoute with other parent Gateway finds no matching Gateways",
				route: basicHTTPRoute("bad-gateway"),
				objects: []client.Object{
					namedGateway("bad-gateway"),
					gatewayClass,
					namespace,
				},
				expected:    []expected{},
				expectedErr: fmt.Errorf("no supported Gateway found for route"),
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				fakeClient := fakeclient.
					NewClientBuilder().
					WithScheme(scheme.Scheme).
					WithObjects(tt.objects...).
					Build()

				got, err := getSupportedGatewayForRoute(
					context.Background(),
					logr.Discard(),
					fakeClient,
					tt.route,
					mo.Some(k8stypes.NamespacedName{
						Namespace: namespace.Name,
						Name:      "good-gateway",
					}),
				)
				if tt.expectedErr != nil {
					require.Equal(t, tt.expectedErr, err)
				} else {
					require.NoError(t, err)
				}
				require.Len(t, got, len(tt.expected))

				for i := range got {
					assert.Equalf(t, "test-namespace", got[i].gateway.Namespace, "gateway namespace #%d", i)
					assert.Equalf(t, "good-gateway", got[i].gateway.Name, "gateway name #%d", i)
					assert.Equalf(t, tt.expected[i].listenerName, got[i].listenerName, "listenerName #%d", i)
					assert.Equalf(t, tt.expected[i].condition, got[i].condition, "condition #%d", i)
				}
			})
		}
	})
}

func TestEnsureParentsProgrammedCondition(t *testing.T) {
	createGateway := func(nn k8stypes.NamespacedName) *gatewayapi.Gateway {
		return &gatewayapi.Gateway{
			TypeMeta: gatewayapi.V1GatewayTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name:      nn.Name,
				Namespace: nn.Namespace,
				UID:       k8stypes.UID(uuid.NewString()),
			},
			Spec: gatewayapi.GatewaySpec{
				GatewayClassName: "test-gatewayclass",
				Listeners:        builder.NewListener("http").WithPort(80).HTTP().IntoSlice(),
			},
			Status: gatewayapi.GatewayStatus{
				Listeners: []gatewayapi.ListenerStatus{
					{
						Name: "http-1",
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayapi.ListenerConditionProgrammed),
								Status: metav1.ConditionTrue,
							},
						},
						SupportedKinds: supportedRouteGroupKinds,
					},
					{
						Name: "http-2",
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayapi.ListenerConditionProgrammed),
								Status: metav1.ConditionTrue,
							},
						},
						SupportedKinds: supportedRouteGroupKinds,
					},
				},
			},
		}
	}

	t.Run("HTTPRoute", func(t *testing.T) {
		gatewayNN1 := k8stypes.NamespacedName{
			Namespace: "test-namespace",
			Name:      "test-gateway",
		}
		gateway1 := createGateway(gatewayNN1)
		gatewayNN2 := k8stypes.NamespacedName{
			Namespace: "test-namespace",
			Name:      "test-gateway-2",
		}
		gateway2 := createGateway(gatewayNN2)

		tests := []struct {
			name           string
			httpRouteFunc  func() *gatewayapi.HTTPRoute
			gatewayFunc    func() []supportedGatewayWithCondition
			expectedUpdate bool
			expectedStatus *gatewayapi.HTTPRouteStatus
		}{
			{
				name: "Programmed condition gets properly set to Status True when parent status is already set in route",
				httpRouteFunc: func() *gatewayapi.HTTPRoute {
					return &gatewayapi.HTTPRoute{
						TypeMeta: gatewayapi.V1HTTPRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayapi.HTTPRouteSpec{
							CommonRouteSpec: gatewayapi.CommonRouteSpec{
								ParentRefs: []gatewayapi.ParentReference{
									{
										Group: lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
										Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:  "test-gateway",
									},
								},
							},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									BackendRefs: []gatewayapi.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
						Status: gatewayapi.HTTPRouteStatus{
							RouteStatus: gatewayapi.RouteStatus{
								Parents: []gatewayapi.RouteParentStatus{
									{
										ParentRef: gatewayapi.ParentReference{
											Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
											Group:     lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
											Name:      gatewayapi.ObjectName(gatewayNN1.Name),
											Namespace: (*gatewayapi.Namespace)(&gatewayNN1.Namespace),
										},
										ControllerName: "konghq.com/kic-gateway-controller",
										Conditions: []metav1.Condition{
											{
												Type:               string(gatewayapi.GatewayConditionAccepted),
												Message:            "",
												ObservedGeneration: 1,
												Status:             metav1.ConditionTrue,
												Reason:             string(gatewayapi.RouteConditionAccepted),
												LastTransitionTime: metav1.Now(),
											},
											{
												Type:               ConditionTypeProgrammed,
												Message:            "",
												ObservedGeneration: 1,
												Status:             metav1.ConditionUnknown,
												Reason:             string(metav1.ConditionUnknown),
												LastTransitionTime: metav1.Now(),
											},
										},
									},
								},
							},
						},
					}
				},
				expectedUpdate: true,
				expectedStatus: &gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ParentRef: gatewayapi.ParentReference{
									Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
									Group:     lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
									Name:      gatewayapi.ObjectName(gatewayNN1.Name),
									Namespace: (*gatewayapi.Namespace)(&gatewayNN1.Namespace),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               string(gatewayapi.GatewayConditionAccepted),
										Message:            "",
										ObservedGeneration: 1,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayapi.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayapi.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
								},
							},
						},
					},
				},
				gatewayFunc: func() []supportedGatewayWithCondition {
					return []supportedGatewayWithCondition{
						{gateway: gateway1},
					}
				},
			},
			{
				name: "Programmed condition gets properly set to Status True when Programmed condition is not present in route's parent status",
				httpRouteFunc: func() *gatewayapi.HTTPRoute {
					return &gatewayapi.HTTPRoute{
						TypeMeta: gatewayapi.V1HTTPRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayapi.HTTPRouteSpec{
							CommonRouteSpec: gatewayapi.CommonRouteSpec{
								ParentRefs: []gatewayapi.ParentReference{
									{
										Group: lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
										Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:  "test-gateway",
									},
								},
							},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									BackendRefs: []gatewayapi.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
						Status: gatewayapi.HTTPRouteStatus{
							RouteStatus: gatewayapi.RouteStatus{
								Parents: []gatewayapi.RouteParentStatus{
									{
										ParentRef: gatewayapi.ParentReference{
											Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
											Group:     lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
											Name:      gatewayapi.ObjectName(gatewayNN1.Name),
											Namespace: (*gatewayapi.Namespace)(&gatewayNN1.Namespace),
										},
										ControllerName: "konghq.com/kic-gateway-controller",
										Conditions: []metav1.Condition{
											{
												Type:               string(gatewayapi.GatewayConditionAccepted),
												Message:            "",
												ObservedGeneration: 1,
												Status:             metav1.ConditionTrue,
												Reason:             string(gatewayapi.RouteConditionAccepted),
												LastTransitionTime: metav1.Now(),
											},
										},
									},
								},
							},
						},
					}
				},
				expectedUpdate: true,
				expectedStatus: &gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ParentRef: gatewayapi.ParentReference{
									Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
									Group:     lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
									Name:      gatewayapi.ObjectName(gatewayNN1.Name),
									Namespace: (*gatewayapi.Namespace)(&gatewayNN1.Namespace),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               string(gatewayapi.GatewayConditionAccepted),
										Message:            "",
										ObservedGeneration: 1,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayapi.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayapi.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
								},
							},
						},
					},
				},
				gatewayFunc: func() []supportedGatewayWithCondition {
					return []supportedGatewayWithCondition{
						{gateway: gateway1},
					}
				},
			},
			{
				name: "Programmed condition gets properly set to Status True when Programmed condition is not present in route's parent status and Parent Section is specified",
				httpRouteFunc: func() *gatewayapi.HTTPRoute {
					return &gatewayapi.HTTPRoute{
						TypeMeta: gatewayapi.V1HTTPRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayapi.HTTPRouteSpec{
							CommonRouteSpec: gatewayapi.CommonRouteSpec{
								ParentRefs: []gatewayapi.ParentReference{
									{
										Group:       lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
										Kind:        lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:        "test-gateway",
										SectionName: lo.ToPtr(gatewayapi.SectionName("http-2")),
									},
								},
							},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									BackendRefs: []gatewayapi.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
						Status: gatewayapi.HTTPRouteStatus{
							RouteStatus: gatewayapi.RouteStatus{
								Parents: []gatewayapi.RouteParentStatus{
									{
										ParentRef: gatewayapi.ParentReference{
											Kind:        lo.ToPtr(gatewayapi.Kind("Gateway")),
											Group:       lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
											Name:        gatewayapi.ObjectName(gatewayNN1.Name),
											Namespace:   (*gatewayapi.Namespace)(&gatewayNN1.Namespace),
											SectionName: lo.ToPtr(gatewayapi.SectionName("http-2")),
										},
										ControllerName: "konghq.com/kic-gateway-controller",
										Conditions:     []metav1.Condition{},
									},
								},
							},
						},
					}
				},
				expectedUpdate: true,
				expectedStatus: &gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ParentRef: gatewayapi.ParentReference{
									Kind:        lo.ToPtr(gatewayapi.Kind("Gateway")),
									Group:       lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
									Name:        gatewayapi.ObjectName(gatewayNN1.Name),
									Namespace:   (*gatewayapi.Namespace)(&gatewayNN1.Namespace),
									SectionName: lo.ToPtr(gatewayapi.SectionName("http-2")),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayapi.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
								},
							},
						},
					},
				},
				gatewayFunc: func() []supportedGatewayWithCondition {
					return []supportedGatewayWithCondition{
						{
							gateway:      gateway1,
							listenerName: "http-2",
						},
					}
				},
			},
			{
				name: "Programmed condition gets properly set to Status True when route's parent status is not set and Parent Section is specified with 2 gateways both with section name specified",
				httpRouteFunc: func() *gatewayapi.HTTPRoute {
					return &gatewayapi.HTTPRoute{
						TypeMeta: gatewayapi.V1HTTPRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayapi.HTTPRouteSpec{
							CommonRouteSpec: gatewayapi.CommonRouteSpec{
								ParentRefs: []gatewayapi.ParentReference{
									{
										Group:       lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
										Kind:        lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:        gatewayapi.ObjectName(gateway1.Name),
										SectionName: lo.ToPtr(gatewayapi.SectionName("http-2")),
									},
									{
										Group:       lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
										Kind:        lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:        gatewayapi.ObjectName(gateway2.Name),
										SectionName: lo.ToPtr(gatewayapi.SectionName("http-1")),
									},
								},
							},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									BackendRefs: []gatewayapi.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
						Status: gatewayapi.HTTPRouteStatus{},
					}
				},
				expectedUpdate: true,
				expectedStatus: &gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ParentRef: gatewayapi.ParentReference{
									Kind:        lo.ToPtr(gatewayapi.Kind("Gateway")),
									Group:       lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
									Name:        gatewayapi.ObjectName(gatewayNN1.Name),
									Namespace:   (*gatewayapi.Namespace)(&gatewayNN1.Namespace),
									SectionName: lo.ToPtr(gatewayapi.SectionName("http-2")),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayapi.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
								},
							},
							{
								ParentRef: gatewayapi.ParentReference{
									Kind:        lo.ToPtr(gatewayapi.Kind("Gateway")),
									Group:       lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
									Name:        gatewayapi.ObjectName(gatewayNN2.Name),
									Namespace:   (*gatewayapi.Namespace)(&gatewayNN2.Namespace),
									SectionName: lo.ToPtr(gatewayapi.SectionName("http-1")),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayapi.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
								},
							},
						},
					},
				},
				gatewayFunc: func() []supportedGatewayWithCondition {
					return []supportedGatewayWithCondition{
						{
							gateway:      gateway1,
							listenerName: "http-2",
						},
						{
							gateway:      gateway2,
							listenerName: "http-1",
						},
					}
				},
			},
			{
				name: "Programmed condition gets properly added to route's parents status when no status for that parent is present yet",
				httpRouteFunc: func() *gatewayapi.HTTPRoute {
					return &gatewayapi.HTTPRoute{
						TypeMeta: gatewayapi.V1HTTPRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayapi.HTTPRouteSpec{
							CommonRouteSpec: gatewayapi.CommonRouteSpec{
								ParentRefs: []gatewayapi.ParentReference{
									{
										Group: lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
										Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:  "test-gateway",
									},
								},
							},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									BackendRefs: []gatewayapi.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
					}
				},
				expectedUpdate: true,
				expectedStatus: &gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ParentRef: gatewayapi.ParentReference{
									Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
									Group:     lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
									Name:      gatewayapi.ObjectName(gatewayNN1.Name),
									Namespace: (*gatewayapi.Namespace)(&gatewayNN1.Namespace),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayapi.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
								},
							},
						},
					},
				},
				gatewayFunc: func() []supportedGatewayWithCondition {
					return []supportedGatewayWithCondition{
						{gateway: gateway1},
					}
				},
			},
			{
				name: "no update is being done when an expected Programmed condition is already in place",
				httpRouteFunc: func() *gatewayapi.HTTPRoute {
					return &gatewayapi.HTTPRoute{
						TypeMeta: gatewayapi.V1HTTPRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayapi.HTTPRouteSpec{
							CommonRouteSpec: gatewayapi.CommonRouteSpec{
								ParentRefs: []gatewayapi.ParentReference{
									{
										Group: lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
										Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
										Name:  "test-gateway",
									},
								},
							},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									BackendRefs: []gatewayapi.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
						Status: gatewayapi.HTTPRouteStatus{
							RouteStatus: gatewayapi.RouteStatus{
								Parents: []gatewayapi.RouteParentStatus{
									{
										ParentRef: gatewayapi.ParentReference{
											Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
											Group:     lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
											Name:      gatewayapi.ObjectName(gatewayNN1.Name),
											Namespace: (*gatewayapi.Namespace)(&gatewayNN1.Namespace),
										},
										ControllerName: "konghq.com/kic-gateway-controller",
										Conditions: []metav1.Condition{
											{
												Type:               ConditionTypeProgrammed,
												Message:            "",
												ObservedGeneration: 42,
												Status:             metav1.ConditionTrue,
												Reason:             string(gatewayapi.RouteConditionAccepted),
												LastTransitionTime: metav1.Now(),
											},
										},
									},
								},
							},
						},
					}
				},
				expectedUpdate: false,
				gatewayFunc: func() []supportedGatewayWithCondition {
					return []supportedGatewayWithCondition{
						{gateway: gateway1},
					}
				},
			},
		}

		for _, tc := range tests {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				var (
					ctx       = context.Background()
					httproute = tc.httpRouteFunc()
					gateways  = tc.gatewayFunc()
				)

				fakeClient := fakeclient.
					NewClientBuilder().
					WithScheme(scheme.Scheme).
					WithObjects(httproute).
					WithStatusSubresource(httproute).
					Build()

				updated, err := ensureParentsProgrammedCondition(ctx, fakeClient.Status(), httproute, httproute.Status.Parents, gateways,
					metav1.Condition{
						Status: metav1.ConditionTrue,
						Reason: string(gatewayapi.RouteConditionAccepted),
					},
				)
				require.NoError(t, err)
				if tc.expectedUpdate {
					require.True(t, updated)
					require.NoError(t, fakeClient.Get(ctx, client.ObjectKeyFromObject(httproute), httproute))

					diff := cmp.Diff(*tc.expectedStatus, httproute.Status,
						// Don't compare the time since metav1.Now() is used.
						cmp.FilterPath(
							func(p cmp.Path) bool { return p.String() == "RouteStatus.Parents.Conditions.LastTransitionTime" },
							cmp.Ignore(),
						),
					)
					if diff != "" {
						t.Errorf("HTTPRoute Status not as expected:\n%s", diff)
					}
				} else {
					require.False(t, updated)
				}
			})
		}
	})
}

func TestIsRouteAcceptedByListener(t *testing.T) {
	testCases := []struct {
		name          string
		gateway       gatewayapi.Gateway
		httpRoute     *gatewayapi.HTTPRoute
		expectedValue bool
	}{
		{
			name: "accepted, allowedRoutes from the same namespace",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayapi.SectionName("listener-1")),
							},
						},
					},
				},
			},
			gateway: gatewayapi.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayapi.HTTPProtocolType,
						},
					},
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: gatewayapi.SectionName("listener-1"),
							SupportedKinds: []gatewayapi.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
						},
					},
				},
			},
			expectedValue: true,
		},
		{
			name: "accepted, allowedRoutes from selected namespace",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "other-namespace",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayapi.SectionName("listener-1")),
							},
						},
					},
				},
			},
			gateway: gatewayapi.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{
						{
							Name: "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSelectorNamespace(&metav1.LabelSelector{
								MatchLabels: map[string]string{
									"konghq.com/allowed-namespace": "true",
								},
							}),
							Protocol: gatewayapi.HTTPProtocolType,
						},
					},
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: gatewayapi.SectionName("listener-1"),
							SupportedKinds: []gatewayapi.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
						},
					},
				},
			},
			expectedValue: true,
		},
		{
			name: "route accepted but listener not programmed",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayapi.SectionName("listener-1")),
							},
						},
					},
				},
			},
			gateway: gatewayapi.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayapi.HTTPProtocolType,
						},
					},
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: gatewayapi.SectionName("listener-1"),
							SupportedKinds: []gatewayapi.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionFalse,
								},
							},
						},
					},
				},
			},
			expectedValue: true,
		},
		{
			name: "not accepted, not in listener's supportedKinds",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayapi.SectionName("listener-1")),
							},
						},
					},
				},
			},
			gateway: gatewayapi.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayapi.HTTPProtocolType,
						},
					},
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: gatewayapi.SectionName("listener-1"),
							SupportedKinds: []gatewayapi.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  "GRPCRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
						},
					},
				},
			},
			expectedValue: false,
		},
		{
			name: "not accepted, wrong sectionName",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayapi.SectionName("wrong-listener")),
							},
						},
					},
				},
			},
			gateway: gatewayapi.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayapi.HTTPProtocolType,
						},
					},
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: gatewayapi.SectionName("listener-1"),
							SupportedKinds: []gatewayapi.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
						},
					},
				},
			},
			expectedValue: false,
		},
		{
			name: "not accepted, wrong port",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayapi.SectionName("listener-1")),
								Port:        lo.ToPtr(gatewayapi.PortNumber(8080)),
							},
						},
					},
				},
			},
			gateway: gatewayapi.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayapi.HTTPProtocolType,
							Port:          9090,
						},
					},
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: gatewayapi.SectionName("listener-1"),
							SupportedKinds: []gatewayapi.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
						},
					},
				},
			},
			expectedValue: false,
		},
		{
			name: "not accepted, wrong protocol",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayapi.SectionName("listener-1")),
							},
						},
					},
				},
			},
			gateway: gatewayapi.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayapi.UDPProtocolType,
						},
					},
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: gatewayapi.SectionName("listener-1"),
							SupportedKinds: []gatewayapi.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
						},
					},
				},
			},
			expectedValue: false,
		},
		{
			name: "not accepted, wrong hostnames",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayapi.SectionName("listener-1")),
							},
						},
					},
					Hostnames: []gatewayapi.Hostname{
						"wrong.hostname.com",
					},
				},
			},
			gateway: gatewayapi.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayapi.HTTPProtocolType,
							Hostname:      lo.ToPtr(gatewayapi.Hostname("foo.bar.com")),
						},
					},
				},
				Status: gatewayapi.GatewayStatus{
					Listeners: []gatewayapi.ListenerStatus{
						{
							Name: gatewayapi.SectionName("listener-1"),
							SupportedKinds: []gatewayapi.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayapi.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayapi.ListenerConditionProgrammed),
									Status: metav1.ConditionTrue,
								},
							},
						},
					},
				},
			},
			expectedValue: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ctx       = context.Background()
				namespace = &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "other-namespace",
						Labels: map[string]string{
							"konghq.com/allowed-namespace": "true",
						},
					},
				}
				fakeClient = fakeclient.
						NewClientBuilder().
						WithScheme(scheme.Scheme).
						WithObjects(namespace).
						Build()
			)

			ok, err := isRouteAcceptedByListener(ctx, fakeClient, tc.httpRoute, tc.gateway, 0, tc.httpRoute.Spec.ParentRefs[0])
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedValue, ok)
		})
	}
}
