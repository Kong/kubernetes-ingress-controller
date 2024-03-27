package gateway

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset/scheme"
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
	commonGateway := &gatewayv1.Gateway{
		Spec: gatewayv1.GatewaySpec{
			Listeners: []Listener{
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
		httpRoute         *gatewayv1.HTTPRoute
		expectedHTTPRoute *gatewayv1.HTTPRoute
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
			httpRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("*.anotherwildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("*.nonmatchingwildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("very.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("very.specific.com"),
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
			httpRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("non.matching.com"),
						util.StringToGatewayAPIHostnameV1Beta1("*.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("very.specific.com"),
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
			httpRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("non.matching.com"),
						util.StringToGatewayAPIHostnameV1Beta1("wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("foo.wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("bar.wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("foo.bar.wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("foo.wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("bar.wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("foo.bar.wildcard.io"),
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
			httpRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("*.anotherwildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("*.anotherwildcard.io"),
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
			httpRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{},
				},
			},
			expectedHTTPRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{},
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
			httpRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("specific.but.wrong.com"),
						util.StringToGatewayAPIHostnameV1Beta1("wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{},
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
	gatewayClass := &GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-gatewayclass",
		},
		Spec: gatewayv1.GatewayClassSpec{
			ControllerName: gatewayv1.GatewayController("konghq.com/kic-gateway-controller"),
		},
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
	}

	routeConditionAccepted := func(status metav1.ConditionStatus, reason gatewayv1.RouteConditionReason) metav1.Condition {
		return metav1.Condition{
			Type:   string(gatewayv1.RouteConditionAccepted),
			Status: status,
			Reason: string(reason),
		}
	}

	type expected struct {
		condition    metav1.Condition
		listenerName string
	}

	goodGroup := gatewayv1.Group(gatewayv1alpha2.GroupName)
	goodKind := gatewayv1.Kind("Gateway")
	basicHTTPRoute := func() *HTTPRoute {
		return &HTTPRoute{
			TypeMeta: metav1.TypeMeta{
				Kind:       "HTTPRoute",
				APIVersion: gatewayv1.GroupVersion.Group + "/" + gatewayv1.GroupVersion.Version,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "basic-httproute",
				Namespace: "test-namespace",
			},
			Spec: gatewayv1.HTTPRouteSpec{
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{
						{
							Group: &goodGroup,
							Kind:  &goodKind,
							Name:  "test-gateway",
						},
					},
				},
				Rules: []gatewayv1.HTTPRouteRule{
					{
						BackendRefs: []gatewayv1.HTTPBackendRef{
							builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
						},
					},
				},
			},
		}
	}
	t.Run("HTTPRoute", func(t *testing.T) {
		gatewayWithHTTP80Ready := func() *Gateway {
			return &Gateway{
				TypeMeta: metav1.TypeMeta{
					APIVersion: gatewayv1.GroupVersion.String(),
					Kind:       "Gateway",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
					UID:       "ce7f0678-f59a-483c-80d1-243d3738d22c",
				},
				Spec: gatewayv1.GatewaySpec{
					GatewayClassName: "test-gatewayclass",
					Listeners:        builder.NewListener("http").WithPort(80).HTTP().IntoSlice(),
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: "http",
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
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
			route    *HTTPRoute
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
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
					},
				},
			},
			{
				name:  "basic HTTPRoute with TLS configuration gets accepted",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Spec.Listeners = builder.
							NewListener("http").WithPort(443).HTTPS().
							WithTLSConfig(&gatewayv1.GatewayTLSConfig{
								Mode: lo.ToPtr(gatewayv1.TLSModeTerminate),
							}).
							IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
					},
				},
			},
			{
				name: "basic HTTPRoute specifying existing section name gets Accepted",
				route: func() *HTTPRoute {
					r := basicHTTPRoute()
					r.Spec.ParentRefs[0].SectionName = lo.ToPtr(gatewayv1.SectionName("http"))
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
						condition:    routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
					},
				},
			},
			{
				name: "basic HTTPRoute specifying existing port gets Accepted",
				route: func() *HTTPRoute {
					r := basicHTTPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1.PortNumber(80))
					return r
				}(),
				objects: []client.Object{
					gatewayWithHTTP80Ready(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
					},
				},
			},
			{
				name: "basic HTTPRoute specifying non-existing port does not get Accepted",
				route: func() *HTTPRoute {
					r := basicHTTPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(PortNumber(80))
					return r
				}(),
				objects: []client.Object{
					func() *Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Spec.Listeners = builder.NewListener("http").WithPort(81).HTTP().IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, RouteReasonNoMatchingParent),
					},
				},
			},
			{
				name:  "basic HTTPRoute does not get accepted if it is not in the supported kinds",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Status.Listeners[0].SupportedKinds = nil
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayv1.RouteReasonNotAllowedByListeners),
					},
				},
			},
			{
				name:  "basic HTTPRoute does not get accepted if it is not permitted by allowed routes",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Spec.Listeners = builder.NewListener("http").
							WithPort(80).
							HTTP().
							WithAllowedRoutes(
								&gatewayv1.AllowedRoutes{
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
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayv1.RouteReasonNotAllowedByListeners),
					},
				},
			},
			{
				name:  "basic HTTPRoute does get accepted if allowed routes only specified Same namespace",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *Gateway {
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
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
					},
				},
			},
			{
				name: "HTTPRoute does not get accepted if Listener hostnames do not match route hostnames",
				route: func() *HTTPRoute {
					r := basicHTTPRoute()
					r.Spec.Hostnames = []gatewayv1.Hostname{"very.specific.com"}
					return r
				}(),
				objects: []client.Object{
					func() *Gateway {
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
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayv1.RouteReasonNoMatchingListenerHostname),
					},
				},
			},
			{
				name:  "HTTPRoute does not get accepted if Listener TLSConfig uses PassThrough",
				route: basicHTTPRoute(),
				objects: []client.Object{
					func() *Gateway {
						gw := gatewayWithHTTP80Ready()
						gw.Spec.Listeners = builder.
							NewListener("https").WithPort(443).HTTPS().
							WithTLSConfig(&gatewayv1.GatewayTLSConfig{
								Mode: lo.ToPtr(gatewayv1.TLSModePassthrough),
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
						condition: routeConditionAccepted(metav1.ConditionFalse, RouteReasonNoMatchingParent),
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

				got, err := getSupportedGatewayForRoute(context.Background(), fakeClient, tt.route)
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
		basicTCPRoute := func() *TCPRoute {
			return &TCPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "TCPRoute",
					APIVersion: gatewayv1alpha2.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-tcproute",
					Namespace: "test-namespace",
				},
				Spec: gatewayv1alpha2.TCPRouteSpec{
					CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
						ParentRefs: []gatewayv1alpha2.ParentReference{
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
		gatewayWithTCP80Ready := func() *Gateway {
			return &Gateway{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "gateway.networking.k8s.io/v1beta1",
					Kind:       "Gateway",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
					UID:       "ce7f0678-f59a-483c-80d1-243d3738d22c",
				},
				Spec: gatewayv1.GatewaySpec{
					GatewayClassName: "test-gatewayclass",
					Listeners:        builder.NewListener("tcp").WithPort(80).TCP().IntoSlice(),
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: "tcp",
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
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
			route    *TCPRoute
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
					condition: routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
				},
			},
			{
				name:  "basic TCPRoute does not get accepted because it is not in supported kinds",
				route: basicTCPRoute(),
				objects: []client.Object{
					func() *Gateway {
						gw := gatewayWithTCP80Ready()
						gw.Status.Listeners[0].SupportedKinds = nil
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionFalse, gatewayv1.RouteReasonNotAllowedByListeners),
				},
			},
			{
				name: "TCPRoute specifying existing port gets Accepted",
				route: func() *TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(80))
					r.Spec.Rules = []gatewayv1alpha2.TCPRouteRule{
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
					condition: routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
				},
			},
			{
				name: "TCPRoute specifying non existing port does not get Accepted",
				route: func() *TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(8000))
					r.Spec.Rules = []gatewayv1alpha2.TCPRouteRule{
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
					condition: routeConditionAccepted(metav1.ConditionFalse, RouteReasonNoMatchingParent),
				},
			},
			{
				name: "TCPRoute specifying in sectionName existing listener gets Accepted",
				route: func() *TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(80))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayv1alpha2.SectionName("tcp"))
					r.Spec.Rules = []gatewayv1alpha2.TCPRouteRule{
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
					condition:    routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
				},
			},
			{
				name: "TCPRoute specifying in sectionName non existing listener does not get Accepted",
				route: func() *TCPRoute {
					r := basicTCPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(80))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayv1alpha2.SectionName("unknown-listener"))
					r.Spec.Rules = []gatewayv1alpha2.TCPRouteRule{
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
					condition:    routeConditionAccepted(metav1.ConditionFalse, RouteReasonNoMatchingParent),
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

				got, err := getSupportedGatewayForRoute(context.Background(), fakeClient, tt.route)
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
		basicUDPRoute := func() *UDPRoute {
			return &UDPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "UDPRoute",
					APIVersion: gatewayv1alpha2.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-udproute",
					Namespace: "test-namespace",
				},
				Spec: gatewayv1alpha2.UDPRouteSpec{
					CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
						ParentRefs: []gatewayv1alpha2.ParentReference{
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
		gatewayWithUDP53Ready := func() *Gateway {
			return &Gateway{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "gateway.networking.k8s.io/v1beta1",
					Kind:       "Gateway",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
					UID:       "ce7f0678-f59a-483c-80d1-243d3738d22c",
				},
				Spec: gatewayv1.GatewaySpec{
					GatewayClassName: "test-gatewayclass",
					Listeners:        builder.NewListener("udp").WithPort(53).UDP().IntoSlice(),
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: "udp",
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
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
			route    *UDPRoute
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
					condition: routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
				},
			},
			{
				name:  "basic UDPRoute does not get accepted because it is not in supported kinds",
				route: basicUDPRoute(),
				objects: []client.Object{
					func() *Gateway {
						gw := gatewayWithUDP53Ready()
						gw.Status.Listeners[0].SupportedKinds = nil
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionFalse, gatewayv1.RouteReasonNotAllowedByListeners),
				},
			},
			{
				name: "UDPRoute specifying existing port gets Accepted",
				route: func() *UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(53))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
				},
			},
			{
				name: "UDPRoute specifying non existing port does not get Accepted",
				route: func() *UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(8000))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					condition: routeConditionAccepted(metav1.ConditionFalse, RouteReasonNoMatchingParent),
				},
			},
			{
				name: "UDPRoute specifying in sectionName existing listener gets Accepted",
				route: func() *UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(53))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayv1alpha2.SectionName("udp"))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "udp",
					condition:    routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
				},
			},
			{
				name: "UDPRoute specifying in sectionName non existing listener does not get Accepted",
				route: func() *UDPRoute {
					r := basicUDPRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(53))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayv1alpha2.SectionName("unknown-listener"))
					return r
				}(),
				objects: []client.Object{
					gatewayWithUDP53Ready(),
					gatewayClass,
					namespace,
				},
				expected: expected{
					listenerName: "unknown-listener",
					condition:    routeConditionAccepted(metav1.ConditionFalse, RouteReasonNoMatchingParent),
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

				got, err := getSupportedGatewayForRoute(context.Background(), fakeClient, tt.route)
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
		basicTLSRoute := func() *TLSRoute {
			return &TLSRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "TLSRoute",
					APIVersion: gatewayv1alpha2.GroupVersion.Group + "/" + gatewayv1alpha2.GroupVersion.Version,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basic-tlsroute",
					Namespace: "test-namespace",
				},
				Spec: gatewayv1alpha2.TLSRouteSpec{
					CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
						ParentRefs: []gatewayv1alpha2.ParentReference{
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
		gatewayWithTLS443PassthroughReady := func() *Gateway {
			return &Gateway{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "gateway.networking.k8s.io/v1beta1",
					Kind:       "Gateway",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
					UID:       "ce7f0678-f59a-483c-80d1-243d3738d22c",
				},
				Spec: gatewayv1.GatewaySpec{
					GatewayClassName: "test-gatewayclass",
					Listeners: builder.NewListener("tls").
						WithPort(443).
						TLS().
						WithTLSConfig(&gatewayv1.GatewayTLSConfig{
							Mode: lo.ToPtr(gatewayv1.TLSModePassthrough),
						}).IntoSlice(),
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: "tls",
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
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
			route    *TLSRoute
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
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
					},
				},
			},
			{
				name:  "basic TLSRoute does not get accepted because there is no listener with TLS in passthrough mode",
				route: basicTLSRoute(),
				objects: []client.Object{
					func() *Gateway {
						gw := gatewayWithTLS443PassthroughReady()
						gw.Spec.Listeners = builder.NewListener("tls").
							WithPort(443).
							TLS().
							WithTLSConfig(&gatewayv1.GatewayTLSConfig{
								Mode: lo.ToPtr(gatewayv1.TLSModeTerminate),
							}).IntoSlice()
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, RouteReasonNoMatchingParent),
					},
				},
			},
			{
				name:  "TLSRoute does not get accepted because it is not in supported kinds",
				route: basicTLSRoute(),
				objects: []client.Object{
					func() *Gateway {
						gw := gatewayWithTLS443PassthroughReady()
						gw.Status.Listeners[0].SupportedKinds = nil
						return gw
					}(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, gatewayv1.RouteReasonNotAllowedByListeners),
					},
				},
			},
			{
				name: "TLSRoute specifying existing port gets Accepted",
				route: func() *TLSRoute {
					r := basicTLSRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(443))
					return r
				}(),
				objects: []client.Object{
					gatewayWithTLS443PassthroughReady(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
					},
				},
			},
			{
				name: "TLSRoute specifying non existing port does not get Accepted",
				route: func() *TLSRoute {
					r := basicTLSRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(444))
					return r
				}(),
				objects: []client.Object{
					gatewayWithTLS443PassthroughReady(),
					gatewayClass,
					namespace,
				},
				expected: []expected{
					{
						condition: routeConditionAccepted(metav1.ConditionFalse, RouteReasonNoMatchingParent),
					},
				},
			},
			{
				name: "TLSRoute specifying in sectionName existing listener gets Accepted",
				route: func() *TLSRoute {
					r := basicTLSRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(443))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayv1alpha2.SectionName("tls"))
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
						condition:    routeConditionAccepted(metav1.ConditionTrue, gatewayv1.RouteReasonAccepted),
					},
				},
			},
			// unmatched sectionName.
			{
				name: "TLSRoute specifying in sectionName non existing listener does not get Accepted",
				route: func() *TLSRoute {
					r := basicTLSRoute()
					r.Spec.CommonRouteSpec.ParentRefs[0].Port = lo.ToPtr(gatewayv1alpha2.PortNumber(443))
					r.Spec.CommonRouteSpec.ParentRefs[0].SectionName = lo.ToPtr(gatewayv1alpha2.SectionName("unknown-listener"))
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
						condition:    routeConditionAccepted(metav1.ConditionFalse, RouteReasonNoMatchingParent),
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

				got, err := getSupportedGatewayForRoute(context.Background(), fakeClient, tt.route)
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
	badGroup := gatewayv1.Group("kechavakunduz.cholpon.uz")
	badKind := gatewayv1.Kind("razzoq")
	bustedParentHTTPRoute.Spec.ParentRefs = []gatewayv1.ParentReference{
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

		_, err := getSupportedGatewayForRoute(context.Background(), fakeClient, bustedParentHTTPRoute)
		require.Equal(t, fmt.Errorf("unsupported parent kind %s/%s", string(badGroup), string(badKind)), err)
	})
}

func TestEnsureParentsProgrammedCondition(t *testing.T) {
	createGateway := func(nn k8stypes.NamespacedName) *Gateway {
		return &Gateway{
			TypeMeta: metav1.TypeMeta{
				APIVersion: gatewayv1.GroupVersion.String(),
				Kind:       "Gateway",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      nn.Name,
				Namespace: nn.Namespace,
				UID:       k8stypes.UID(uuid.NewString()),
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: "test-gatewayclass",
				Listeners:        builder.NewListener("http").WithPort(80).HTTP().IntoSlice(),
			},
			Status: gatewayv1.GatewayStatus{
				Listeners: []gatewayv1.ListenerStatus{
					{
						Name: "http-1",
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayv1.ListenerConditionProgrammed),
								Status: metav1.ConditionTrue,
							},
						},
						SupportedKinds: supportedRouteGroupKinds,
					},
					{
						Name: "http-2",
						Conditions: []metav1.Condition{
							{
								Type:   string(gatewayv1.ListenerConditionProgrammed),
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
			httpRouteFunc  func() *HTTPRoute
			gatewayFunc    func() []supportedGatewayWithCondition
			expectedUpdate bool
			expectedStatus *gatewayv1.HTTPRouteStatus
		}{
			{
				name: "Programmed condition gets properly set to Status True when parent status is already set in route",
				httpRouteFunc: func() *HTTPRoute {
					return &HTTPRoute{
						TypeMeta: metav1.TypeMeta{
							Kind:       "HTTPRoute",
							APIVersion: gatewayv1.GroupVersion.Group + "/" + gatewayv1.GroupVersion.Version,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayv1.HTTPRouteSpec{
							CommonRouteSpec: gatewayv1.CommonRouteSpec{
								ParentRefs: []gatewayv1.ParentReference{
									{
										Group: lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
										Kind:  lo.ToPtr(gatewayv1.Kind("Gateway")),
										Name:  "test-gateway",
									},
								},
							},
							Rules: []gatewayv1.HTTPRouteRule{
								{
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
						Status: gatewayv1.HTTPRouteStatus{
							RouteStatus: gatewayv1.RouteStatus{
								Parents: []gatewayv1.RouteParentStatus{
									{
										ParentRef: gatewayv1.ParentReference{
											Kind:      lo.ToPtr(gatewayv1.Kind("Gateway")),
											Group:     lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
											Name:      gatewayv1.ObjectName(gatewayNN1.Name),
											Namespace: (*gatewayv1.Namespace)(&gatewayNN1.Namespace),
										},
										ControllerName: "konghq.com/kic-gateway-controller",
										Conditions: []metav1.Condition{
											{
												Type:               string(gatewayv1.GatewayConditionAccepted),
												Message:            "",
												ObservedGeneration: 1,
												Status:             metav1.ConditionTrue,
												Reason:             string(gatewayv1.RouteConditionAccepted),
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
				expectedStatus: &gatewayv1.HTTPRouteStatus{
					RouteStatus: gatewayv1.RouteStatus{
						Parents: []gatewayv1.RouteParentStatus{
							{
								ParentRef: gatewayv1.ParentReference{
									Kind:      lo.ToPtr(gatewayv1.Kind("Gateway")),
									Group:     lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
									Name:      gatewayv1.ObjectName(gatewayNN1.Name),
									Namespace: (*gatewayv1.Namespace)(&gatewayNN1.Namespace),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               string(gatewayv1.GatewayConditionAccepted),
										Message:            "",
										ObservedGeneration: 1,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayv1.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayv1.RouteConditionAccepted),
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
				httpRouteFunc: func() *HTTPRoute {
					return &HTTPRoute{
						TypeMeta: metav1.TypeMeta{
							Kind:       "HTTPRoute",
							APIVersion: gatewayv1.GroupVersion.Group + "/" + gatewayv1.GroupVersion.Version,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayv1.HTTPRouteSpec{
							CommonRouteSpec: gatewayv1.CommonRouteSpec{
								ParentRefs: []gatewayv1.ParentReference{
									{
										Group: lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
										Kind:  lo.ToPtr(gatewayv1.Kind("Gateway")),
										Name:  "test-gateway",
									},
								},
							},
							Rules: []gatewayv1.HTTPRouteRule{
								{
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
						Status: gatewayv1.HTTPRouteStatus{
							RouteStatus: gatewayv1.RouteStatus{
								Parents: []gatewayv1.RouteParentStatus{
									{
										ParentRef: gatewayv1.ParentReference{
											Kind:      lo.ToPtr(gatewayv1.Kind("Gateway")),
											Group:     lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
											Name:      gatewayv1.ObjectName(gatewayNN1.Name),
											Namespace: (*gatewayv1.Namespace)(&gatewayNN1.Namespace),
										},
										ControllerName: "konghq.com/kic-gateway-controller",
										Conditions: []metav1.Condition{
											{
												Type:               string(gatewayv1.GatewayConditionAccepted),
												Message:            "",
												ObservedGeneration: 1,
												Status:             metav1.ConditionTrue,
												Reason:             string(gatewayv1.RouteConditionAccepted),
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
				expectedStatus: &gatewayv1.HTTPRouteStatus{
					RouteStatus: gatewayv1.RouteStatus{
						Parents: []gatewayv1.RouteParentStatus{
							{
								ParentRef: gatewayv1.ParentReference{
									Kind:      lo.ToPtr(gatewayv1.Kind("Gateway")),
									Group:     lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
									Name:      gatewayv1.ObjectName(gatewayNN1.Name),
									Namespace: (*gatewayv1.Namespace)(&gatewayNN1.Namespace),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               string(gatewayv1.GatewayConditionAccepted),
										Message:            "",
										ObservedGeneration: 1,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayv1.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayv1.RouteConditionAccepted),
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
				httpRouteFunc: func() *HTTPRoute {
					return &HTTPRoute{
						TypeMeta: metav1.TypeMeta{
							Kind:       "HTTPRoute",
							APIVersion: gatewayv1.GroupVersion.Group + "/" + gatewayv1.GroupVersion.Version,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayv1.HTTPRouteSpec{
							CommonRouteSpec: gatewayv1.CommonRouteSpec{
								ParentRefs: []gatewayv1.ParentReference{
									{
										Group:       lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
										Kind:        lo.ToPtr(gatewayv1.Kind("Gateway")),
										Name:        "test-gateway",
										SectionName: lo.ToPtr(gatewayv1.SectionName("http-2")),
									},
								},
							},
							Rules: []gatewayv1.HTTPRouteRule{
								{
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
						Status: gatewayv1.HTTPRouteStatus{
							RouteStatus: gatewayv1.RouteStatus{
								Parents: []gatewayv1.RouteParentStatus{
									{
										ParentRef: gatewayv1.ParentReference{
											Kind:        lo.ToPtr(gatewayv1.Kind("Gateway")),
											Group:       lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
											Name:        gatewayv1.ObjectName(gatewayNN1.Name),
											Namespace:   (*gatewayv1.Namespace)(&gatewayNN1.Namespace),
											SectionName: lo.ToPtr(gatewayv1.SectionName("http-2")),
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
				expectedStatus: &gatewayv1.HTTPRouteStatus{
					RouteStatus: gatewayv1.RouteStatus{
						Parents: []gatewayv1.RouteParentStatus{
							{
								ParentRef: gatewayv1.ParentReference{
									Kind:        lo.ToPtr(gatewayv1.Kind("Gateway")),
									Group:       lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
									Name:        gatewayv1.ObjectName(gatewayNN1.Name),
									Namespace:   (*gatewayv1.Namespace)(&gatewayNN1.Namespace),
									SectionName: lo.ToPtr(gatewayv1.SectionName("http-2")),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayv1.RouteConditionAccepted),
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
				httpRouteFunc: func() *HTTPRoute {
					return &HTTPRoute{
						TypeMeta: metav1.TypeMeta{
							Kind:       "HTTPRoute",
							APIVersion: gatewayv1.GroupVersion.Group + "/" + gatewayv1.GroupVersion.Version,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayv1.HTTPRouteSpec{
							CommonRouteSpec: gatewayv1.CommonRouteSpec{
								ParentRefs: []gatewayv1.ParentReference{
									{
										Group:       lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
										Kind:        lo.ToPtr(gatewayv1.Kind("Gateway")),
										Name:        gatewayv1.ObjectName(gateway1.Name),
										SectionName: lo.ToPtr(gatewayv1.SectionName("http-2")),
									},
									{
										Group:       lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
										Kind:        lo.ToPtr(gatewayv1.Kind("Gateway")),
										Name:        gatewayv1.ObjectName(gateway2.Name),
										SectionName: lo.ToPtr(gatewayv1.SectionName("http-1")),
									},
								},
							},
							Rules: []gatewayv1.HTTPRouteRule{
								{
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
						Status: gatewayv1.HTTPRouteStatus{},
					}
				},
				expectedUpdate: true,
				expectedStatus: &gatewayv1.HTTPRouteStatus{
					RouteStatus: gatewayv1.RouteStatus{
						Parents: []gatewayv1.RouteParentStatus{
							{
								ParentRef: gatewayv1.ParentReference{
									Kind:        lo.ToPtr(gatewayv1.Kind("Gateway")),
									Group:       lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
									Name:        gatewayv1.ObjectName(gatewayNN1.Name),
									Namespace:   (*gatewayv1.Namespace)(&gatewayNN1.Namespace),
									SectionName: lo.ToPtr(gatewayv1.SectionName("http-2")),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayv1.RouteConditionAccepted),
										LastTransitionTime: metav1.Now(),
									},
								},
							},
							{
								ParentRef: gatewayv1.ParentReference{
									Kind:        lo.ToPtr(gatewayv1.Kind("Gateway")),
									Group:       lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
									Name:        gatewayv1.ObjectName(gatewayNN2.Name),
									Namespace:   (*gatewayv1.Namespace)(&gatewayNN2.Namespace),
									SectionName: lo.ToPtr(gatewayv1.SectionName("http-1")),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayv1.RouteConditionAccepted),
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
				httpRouteFunc: func() *HTTPRoute {
					return &HTTPRoute{
						TypeMeta: metav1.TypeMeta{
							Kind:       "HTTPRoute",
							APIVersion: gatewayv1.GroupVersion.Group + "/" + gatewayv1.GroupVersion.Version,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayv1.HTTPRouteSpec{
							CommonRouteSpec: gatewayv1.CommonRouteSpec{
								ParentRefs: []gatewayv1.ParentReference{
									{
										Group: lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
										Kind:  lo.ToPtr(gatewayv1.Kind("Gateway")),
										Name:  "test-gateway",
									},
								},
							},
							Rules: []gatewayv1.HTTPRouteRule{
								{
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
					}
				},
				expectedUpdate: true,
				expectedStatus: &gatewayv1.HTTPRouteStatus{
					RouteStatus: gatewayv1.RouteStatus{
						Parents: []gatewayv1.RouteParentStatus{
							{
								ParentRef: gatewayv1.ParentReference{
									Kind:      lo.ToPtr(gatewayv1.Kind("Gateway")),
									Group:     lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
									Name:      gatewayv1.ObjectName(gatewayNN1.Name),
									Namespace: (*gatewayv1.Namespace)(&gatewayNN1.Namespace),
								},
								ControllerName: "konghq.com/kic-gateway-controller",
								Conditions: []metav1.Condition{
									{
										Type:               ConditionTypeProgrammed,
										Message:            "",
										ObservedGeneration: 42,
										Status:             metav1.ConditionTrue,
										Reason:             string(gatewayv1.RouteConditionAccepted),
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
				httpRouteFunc: func() *HTTPRoute {
					return &HTTPRoute{
						TypeMeta: metav1.TypeMeta{
							Kind:       "HTTPRoute",
							APIVersion: gatewayv1.GroupVersion.Group + "/" + gatewayv1.GroupVersion.Version,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:       "basic-httproute",
							Namespace:  gatewayNN1.Namespace,
							Generation: 42,
						},
						Spec: gatewayv1.HTTPRouteSpec{
							CommonRouteSpec: gatewayv1.CommonRouteSpec{
								ParentRefs: []gatewayv1.ParentReference{
									{
										Group: lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
										Kind:  lo.ToPtr(gatewayv1.Kind("Gateway")),
										Name:  "test-gateway",
									},
								},
							},
							Rules: []gatewayv1.HTTPRouteRule{
								{
									BackendRefs: []gatewayv1.HTTPBackendRef{
										builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
									},
								},
							},
						},
						Status: gatewayv1.HTTPRouteStatus{
							RouteStatus: gatewayv1.RouteStatus{
								Parents: []gatewayv1.RouteParentStatus{
									{
										ParentRef: gatewayv1.ParentReference{
											Kind:      lo.ToPtr(gatewayv1.Kind("Gateway")),
											Group:     lo.ToPtr(gatewayv1.Group(gatewayv1.GroupName)),
											Name:      gatewayv1.ObjectName(gatewayNN1.Name),
											Namespace: (*gatewayv1.Namespace)(&gatewayNN1.Namespace),
										},
										ControllerName: "konghq.com/kic-gateway-controller",
										Conditions: []metav1.Condition{
											{
												Type:               ConditionTypeProgrammed,
												Message:            "",
												ObservedGeneration: 42,
												Status:             metav1.ConditionTrue,
												Reason:             string(gatewayv1.RouteConditionAccepted),
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
						Reason: string(gatewayv1.RouteConditionAccepted),
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
		gateway       gatewayv1.Gateway
		httpRoute     *gatewayv1.HTTPRoute
		expectedValue bool
	}{
		{
			name: "accepted, allowedRoutes from the same namespace",
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayv1.SectionName("listener-1")),
							},
						},
					},
				},
			},
			gateway: gatewayv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1alpha2.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayv1.HTTPProtocolType,
						},
					},
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: gatewayv1.SectionName("listener-1"),
							SupportedKinds: []gatewayv1.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayv1.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
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
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "other-namespace",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayv1.SectionName("listener-1")),
							},
						},
					},
				},
			},
			gateway: gatewayv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1alpha2.Listener{
						{
							Name: "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSelectorNamespace(&metav1.LabelSelector{
								MatchLabels: map[string]string{
									"konghq.com/allowed-namespace": "true",
								},
							}),
							Protocol: gatewayv1.HTTPProtocolType,
						},
					},
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: gatewayv1.SectionName("listener-1"),
							SupportedKinds: []gatewayv1.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayv1.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
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
			name: "not accepted, listener not ready",
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayv1.SectionName("listener-1")),
							},
						},
					},
				},
			},
			gateway: gatewayv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1alpha2.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayv1.HTTPProtocolType,
						},
					},
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: gatewayv1.SectionName("listener-1"),
							SupportedKinds: []gatewayv1.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayv1.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
									Status: metav1.ConditionFalse,
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
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayv1.SectionName("wrong-listener")),
							},
						},
					},
				},
			},
			gateway: gatewayv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1alpha2.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayv1.HTTPProtocolType,
						},
					},
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: gatewayv1.SectionName("listener-1"),
							SupportedKinds: []gatewayv1.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayv1.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
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
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayv1.SectionName("listener-1")),
								Port:        lo.ToPtr(gatewayv1.PortNumber(8080)),
							},
						},
					},
				},
			},
			gateway: gatewayv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1alpha2.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayv1.HTTPProtocolType,
							Port:          9090,
						},
					},
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: gatewayv1.SectionName("listener-1"),
							SupportedKinds: []gatewayv1.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayv1.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
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
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayv1.SectionName("listener-1")),
							},
						},
					},
				},
			},
			gateway: gatewayv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1alpha2.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayv1.UDPProtocolType,
						},
					},
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: gatewayv1.SectionName("listener-1"),
							SupportedKinds: []gatewayv1.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayv1.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
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
			httpRoute: &gatewayv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "default",
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name:        "gateway",
								SectionName: lo.ToPtr(gatewayv1.SectionName("listener-1")),
							},
						},
					},
					Hostnames: []gatewayv1.Hostname{
						"wrong.hostname.com",
					},
				},
			},
			gateway: gatewayv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway",
					Namespace: "default",
				},
				Spec: gatewayv1.GatewaySpec{
					Listeners: []gatewayv1alpha2.Listener{
						{
							Name:          "listener-1",
							AllowedRoutes: builder.NewAllowedRoutesFromSameNamespaces(),
							Protocol:      gatewayv1.HTTPProtocolType,
							Hostname:      lo.ToPtr(gatewayv1.Hostname("foo.bar.com")),
						},
					},
				},
				Status: gatewayv1.GatewayStatus{
					Listeners: []gatewayv1.ListenerStatus{
						{
							Name: gatewayv1.SectionName("listener-1"),
							SupportedKinds: []gatewayv1.RouteGroupKind{
								{
									Group: lo.ToPtr(gatewayv1.Group("gateway.networking.k8s.io")),
									Kind:  "HTTPRoute",
								},
							},
							Conditions: []metav1.Condition{
								{
									Type:   string(gatewayv1.ListenerConditionProgrammed),
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
