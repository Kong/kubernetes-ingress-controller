package gateway

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/scheme"
)

func testIsRouteAttachedToReconciledGateway[routeT gatewayapi.RouteT](
	t *testing.T,
	cl client.Client,
	gatewayNN controllers.OptionalNamespacedName,
	route routeT,
	expectedResult bool,
) {
	logger := logr.Discard()

	result := IsRouteAttachedToReconciledGateway[routeT](cl, logger, gatewayNN, route)
	require.Equal(t, expectedResult, result)
}

func TestIsRouteAttachedToReconciledGateway(t *testing.T) {
	type httpRouteTestCase struct {
		name           string
		objects        []client.Object
		route          *gatewayapi.HTTPRoute
		gatewayNN      controllers.OptionalNamespacedName
		expectedResult bool
	}

	kongGWClass := &gatewayapi.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong",
		},
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: "konghq.com/kic-gateway-controller",
		},
	}

	anotherGWClass := &gatewayapi.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "another",
		},
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: "another",
		},
	}

	kongGW := &gatewayapi.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kong",
			Namespace: "default",
		},
		Spec: gatewayapi.GatewaySpec{
			GatewayClassName: "kong",
		},
	}

	anotherGW := &gatewayapi.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "another",
			Namespace: "default",
		},
		Spec: gatewayapi.GatewaySpec{
			GatewayClassName: "another",
		},
	}

	httpRouteTestCases := []httpRouteTestCase{
		{
			name: "single parent ref to gateway with expected class",
			objects: []client.Object{
				kongGW,
				kongGWClass,
			},
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kong-httproute",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Group:     lo.ToPtr(gatewayapi.V1Group),
								Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
								Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
								Name:      gatewayapi.ObjectName("kong"),
							},
						},
					},
				},
			},
			expectedResult: true,
		},
		{
			name: "single parent ref to gateway with another class",
			objects: []client.Object{
				anotherGW,
				anotherGWClass,
			},
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kong-httproute",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Group:     lo.ToPtr(gatewayapi.V1Group),
								Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
								Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
								Name:      gatewayapi.ObjectName("another"),
							},
						},
					},
				},
			},
			expectedResult: false,
		},
		{
			name: "single parent ref to specified gateway",
			objects: []client.Object{
				anotherGW,
				anotherGWClass,
			},
			gatewayNN: controllers.NewOptionalNamespacedName(mo.Some(
				k8stypes.NamespacedName{
					Namespace: "default",
					Name:      "another",
				},
			)),
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kong-httproute",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Group: lo.ToPtr(gatewayapi.V1Group),
								Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
								Name:  gatewayapi.ObjectName("another"),
							},
						},
					},
				},
			},
			expectedResult: true,
		},
		{
			name: "multiple parent refs with one pointing to reconciled gateway",
			objects: []client.Object{
				kongGW,
				kongGWClass,
			},
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kong-httproute",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Kind: lo.ToPtr(gatewayapi.Kind("Service")),
								Name: gatewayapi.ObjectName("kuma"),
							},
							{
								Group: lo.ToPtr(gatewayapi.V1Group),
								Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
								Name:  gatewayapi.ObjectName("kong"),
							},
						},
					},
				},
			},
			expectedResult: true,
		},
		{
			name: "parent ref pointing to non-exist gateway",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kong-httproute",
					Namespace: "default",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Group: lo.ToPtr(gatewayapi.V1Group),
								Kind:  lo.ToPtr(gatewayapi.Kind("Gateway")),
								Name:  gatewayapi.ObjectName("non-exist"),
							},
						},
					},
				},
			},
			expectedResult: false,
		},
	}

	for _, tc := range httpRouteTestCases {

		cl := fakeclient.NewClientBuilder().WithScheme(lo.Must(scheme.Get())).WithObjects(tc.objects...).Build()
		t.Run(tc.name, func(t *testing.T) {
			testIsRouteAttachedToReconciledGateway(
				t,
				cl,
				tc.gatewayNN,
				tc.route,
				tc.expectedResult,
			)
		},
		)
	}
}
