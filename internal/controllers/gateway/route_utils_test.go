package gateway

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset/scheme"
)

func init() {
	if err := corev1.AddToScheme(scheme.Scheme); err != nil {
		fmt.Println("error while adding core1 scheme")
		os.Exit(1)
	}
	if err := gatewayv1beta1.Install(scheme.Scheme); err != nil {
		fmt.Println("error while adding gatewayv1beta1 scheme")
		os.Exit(1)
	}
}

func TestFilterHostnames(t *testing.T) {
	commonGateway := &gatewayv1beta1.Gateway{
		Spec: gatewayv1beta1.GatewaySpec{
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
			},
		},
	}

	testCases := []struct {
		name              string
		gateways          []supportedGatewayWithCondition
		httpRoute         *gatewayv1beta1.HTTPRoute
		expectedHTTPRoute *gatewayv1beta1.HTTPRoute
	}{
		{
			name: "listener 1 - specific",
			gateways: []supportedGatewayWithCondition{
				{
					gateway:      commonGateway,
					listenerName: "listener-1",
				},
			},
			httpRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("*.anotherwildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("*.nonmatchingwildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("very.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
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
			httpRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("non.matching.com"),
						util.StringToGatewayAPIHostnameV1Beta1("*.specific.com"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
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
			httpRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("non.matching.com"),
						util.StringToGatewayAPIHostnameV1Beta1("wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("foo.wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("bar.wildcard.io"),
						util.StringToGatewayAPIHostnameV1Beta1("foo.bar.wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
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
			httpRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("*.anotherwildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("*.anotherwildcard.io"),
					},
				},
			},
		},
		{
			name: "no match",
			gateways: []supportedGatewayWithCondition{
				{
					gateway: commonGateway,
				},
			},
			httpRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						util.StringToGatewayAPIHostnameV1Beta1("specific.but.wrong.com"),
						util.StringToGatewayAPIHostnameV1Beta1("wildcard.io"),
					},
				},
			},
			expectedHTTPRoute: &gatewayv1beta1.HTTPRoute{
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{},
				},
			},
		},
	}

	for _, tc := range testCases {
		filteredHTTPRoute := filterHostnames(tc.gateways, tc.httpRoute)
		assert.Equal(t, tc.expectedHTTPRoute.Spec, filteredHTTPRoute.Spec, tc.name)
	}
}

func addressOf[T any](v T) *T {
	return &v
}

func Test_getSupportedGatewayForRoute(t *testing.T) {
	t.Run("HTTPRoute", func(t *testing.T) {
		type expected struct {
			gateway      types.NamespacedName
			condition    metav1.Condition
			listenerName string
		}
		tests := []struct {
			name     string
			route    *HTTPRoute
			expected []expected
			objects  []client.Object
			wantErr  bool
		}{
			{
				name: "basic HTTPRoute gets accepted",
				route: &HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "basic-httproute",
						Namespace: "test-namespace",
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
							ParentRefs: []gatewayv1beta1.ParentReference{
								{
									Name: gatewayv1beta1.ObjectName("test-gateway"),
								},
							},
						},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								BackendRefs: []gatewayv1beta1.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
							},
						},
					},
				},
				objects: []client.Object{
					&Gateway{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "gateway.networking.k8s.io/v1beta1",
							Kind:       "Gateway",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-gateway",
							Namespace: "test-namespace",
							UID:       types.UID("ce7f0678-f59a-483c-80d1-243d3738d22c"),
						},
						Spec: gatewayv1beta1.GatewaySpec{
							GatewayClassName: "test-gatewayclass",
							Listeners: []gatewayv1beta1.Listener{
								{
									Name:     gatewayv1beta1.SectionName("http"),
									Protocol: gatewayv1beta1.HTTPProtocolType,
									Port:     gatewayv1beta1.PortNumber(80),
								},
							},
						},
						Status: gatewayv1beta1.GatewayStatus{
							Listeners: []gatewayv1beta1.ListenerStatus{
								{
									Name: gatewayv1beta1.SectionName("http"),
									Conditions: []metav1.Condition{
										{
											Type:   "Ready",
											Status: metav1.ConditionTrue,
										},
									},
								},
							},
						},
					},
					&GatewayClass{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-gatewayclass",
						},
						Spec: gatewayv1beta1.GatewayClassSpec{
							ControllerName: gatewayv1beta1.GatewayController("konghq.com/kic-gateway-controller"),
						},
					},
					&corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-namespace",
						},
					},
				},
				expected: []expected{
					{
						gateway: types.NamespacedName{
							Name:      "test-gateway",
							Namespace: "test-namespace",
						},
						listenerName: "",
						condition: metav1.Condition{
							Type:   string(gatewayv1beta1.RouteConditionAccepted),
							Status: metav1.ConditionTrue,
							Reason: string(gatewayv1beta1.RouteReasonAccepted),
						},
					},
				},
			},
			{
				name: "basic HTTPRoute specifying existing section name gets Accepted",
				route: &HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "basic-httproute",
						Namespace: "test-namespace",
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
							ParentRefs: []gatewayv1beta1.ParentReference{
								{
									Name:        gatewayv1beta1.ObjectName("test-gateway"),
									SectionName: addressOf(gatewayv1beta1.SectionName("http")),
								},
							},
						},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								BackendRefs: []gatewayv1beta1.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
							},
						},
					},
				},
				objects: []client.Object{
					&Gateway{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "gateway.networking.k8s.io/v1beta1",
							Kind:       "Gateway",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-gateway",
							Namespace: "test-namespace",
							UID:       types.UID("ce7f0678-f59a-483c-80d1-243d3738d22c"),
						},
						Spec: gatewayv1beta1.GatewaySpec{
							GatewayClassName: "test-gatewayclass",
							Listeners: []gatewayv1beta1.Listener{
								{
									Name:     gatewayv1beta1.SectionName("http"),
									Protocol: gatewayv1beta1.HTTPProtocolType,
									Port:     gatewayv1beta1.PortNumber(80),
								},
							},
						},
						Status: gatewayv1beta1.GatewayStatus{
							Listeners: []gatewayv1beta1.ListenerStatus{
								{
									Name: gatewayv1beta1.SectionName("http"),
									Conditions: []metav1.Condition{
										{
											Type:   "Ready",
											Status: metav1.ConditionTrue,
										},
									},
								},
							},
						},
					},
					&GatewayClass{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-gatewayclass",
						},
						Spec: gatewayv1beta1.GatewayClassSpec{
							ControllerName: gatewayv1beta1.GatewayController("konghq.com/kic-gateway-controller"),
						},
					},
					&corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-namespace",
						},
					},
				},
				expected: []expected{
					{
						gateway: types.NamespacedName{
							Name:      "test-gateway",
							Namespace: "test-namespace",
						},
						listenerName: "http",
						condition: metav1.Condition{
							Type:   string(gatewayv1beta1.RouteConditionAccepted),
							Status: metav1.ConditionTrue,
							Reason: string(gatewayv1beta1.RouteReasonAccepted),
						},
					},
				},
			},
			{
				name: "basic HTTPRoute specifying existing port gets Accepted",
				route: &HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "basic-httproute",
						Namespace: "test-namespace",
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
							ParentRefs: []gatewayv1beta1.ParentReference{
								{
									Name: gatewayv1beta1.ObjectName("test-gateway"),
									Port: addressOf(gatewayv1beta1.PortNumber(80)),
								},
							},
						},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								BackendRefs: []gatewayv1beta1.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
							},
						},
					},
				},
				objects: []client.Object{
					&Gateway{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "gateway.networking.k8s.io/v1beta1",
							Kind:       "Gateway",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-gateway",
							Namespace: "test-namespace",
							UID:       types.UID("ce7f0678-f59a-483c-80d1-243d3738d22c"),
						},
						Spec: gatewayv1beta1.GatewaySpec{
							GatewayClassName: "test-gatewayclass",
							Listeners: []gatewayv1beta1.Listener{
								{
									Name:     gatewayv1beta1.SectionName("http"),
									Protocol: gatewayv1beta1.HTTPProtocolType,
									Port:     gatewayv1beta1.PortNumber(80),
								},
							},
						},
						Status: gatewayv1beta1.GatewayStatus{
							Listeners: []gatewayv1beta1.ListenerStatus{
								{
									Name: gatewayv1beta1.SectionName("http"),
									Conditions: []metav1.Condition{
										{
											Type:   "Ready",
											Status: metav1.ConditionTrue,
										},
									},
								},
							},
						},
					},
					&GatewayClass{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-gatewayclass",
						},
						Spec: gatewayv1beta1.GatewayClassSpec{
							ControllerName: gatewayv1beta1.GatewayController("konghq.com/kic-gateway-controller"),
						},
					},
					&corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-namespace",
						},
					},
				},
				expected: []expected{
					{
						gateway: types.NamespacedName{
							Name:      "test-gateway",
							Namespace: "test-namespace",
						},
						listenerName: "",
						condition: metav1.Condition{
							Type:   string(gatewayv1beta1.RouteConditionAccepted),
							Status: metav1.ConditionTrue,
							Reason: string(gatewayv1beta1.RouteReasonAccepted),
						},
					},
				},
			},
			{
				name: "basic HTTPRoute specifying non-existing port does not get Accepted",
				route: &HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "basic-httproute",
						Namespace: "test-namespace",
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
							ParentRefs: []gatewayv1beta1.ParentReference{
								{
									Name: gatewayv1beta1.ObjectName("test-gateway"),
									Port: addressOf(gatewayv1beta1.PortNumber(80)),
								},
							},
						},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								BackendRefs: []gatewayv1beta1.HTTPBackendRef{
									builder.NewHTTPBackendRef("fake-service").WithPort(80).Build(),
								},
							},
						},
					},
				},
				objects: []client.Object{
					&Gateway{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "gateway.networking.k8s.io/v1beta1",
							Kind:       "Gateway",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-gateway",
							Namespace: "test-namespace",
							UID:       types.UID("ce7f0678-f59a-483c-80d1-243d3738d22c"),
						},
						Spec: gatewayv1beta1.GatewaySpec{
							GatewayClassName: "test-gatewayclass",
							Listeners: []gatewayv1beta1.Listener{
								{
									Name:     gatewayv1beta1.SectionName("http"),
									Protocol: gatewayv1beta1.HTTPProtocolType,
									Port:     gatewayv1beta1.PortNumber(81),
								},
							},
						},
						Status: gatewayv1beta1.GatewayStatus{
							Listeners: []gatewayv1beta1.ListenerStatus{
								{
									Name: gatewayv1beta1.SectionName("http"),
									Conditions: []metav1.Condition{
										{
											Type:   "Ready",
											Status: metav1.ConditionTrue,
										},
									},
								},
							},
						},
					},
					&GatewayClass{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-gatewayclass",
						},
						Spec: gatewayv1beta1.GatewayClassSpec{
							ControllerName: gatewayv1beta1.GatewayController("konghq.com/kic-gateway-controller"),
						},
					},
					&corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-namespace",
						},
					},
				},
				expected: []expected{
					{
						gateway: types.NamespacedName{
							Name:      "test-gateway",
							Namespace: "test-namespace",
						},
						listenerName: "",
						condition: metav1.Condition{
							Type:   string(gatewayv1beta1.RouteConditionAccepted),
							Status: metav1.ConditionFalse,
							Reason: string(RouteReasonNoMatchingListenerPort),
						},
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
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.Len(t, got, len(tt.expected))

					for i := range got {
						assert.Equalf(t, tt.expected[i].gateway.Namespace, got[i].gateway.Namespace, "gateway namespace #%d", i)
						assert.Equalf(t, tt.expected[i].gateway.Name, got[i].gateway.Name, "gateway name #%d", i)
						assert.Equalf(t, tt.expected[i].listenerName, got[i].listenerName, "listenerName #%d", i)
						assert.Equalf(t, tt.expected[i].condition, got[i].condition, "condition #%d", i)
					}
				}
			})
		}
	})
}
