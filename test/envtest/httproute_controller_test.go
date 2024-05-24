//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

const (
	waitDuration = 5 * time.Second
	tickDuration = 100 * time.Millisecond
)

func TestHTTPRouteReconcilerProperlyReactsToReferenceGrant(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithGatewayAPI)
	cfg := Setup(t, scheme)
	client := NewControllerClient(t, scheme, cfg)

	reconciler := &gateway.HTTPRouteReconciler{
		Client:          client,
		DataplaneClient: mocks.Dataplane{},
	}

	// We use a deferred cancel to stop the manager and not wait for its timeout.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ns := CreateNamespace(ctx, t, client)
	nsRoute := CreateNamespace(ctx, t, client)

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "backend-1",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}
	require.NoError(t, client.Create(ctx, &svc))
	StartReconcilers(ctx, t, client.Scheme(), cfg, reconciler)

	gwc := gatewayapi.GatewayClass{
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				"konghq.com/gatewayclass-unmanaged": "placeholder",
			},
		},
	}
	require.NoError(t, client.Create(ctx, &gwc))
	t.Cleanup(func() { _ = client.Delete(ctx, &gwc) })

	gw := gatewayapi.Gateway{
		Spec: gatewayapi.GatewaySpec{
			GatewayClassName: gatewayapi.ObjectName(gwc.Name),
			Listeners: []gatewayapi.Listener{
				{
					Name:     gatewayapi.SectionName("http"),
					Port:     gatewayapi.PortNumber(80),
					Protocol: gatewayapi.HTTPProtocolType,
					AllowedRoutes: &gatewayapi.AllowedRoutes{
						Namespaces: &gatewayapi.RouteNamespaces{
							From: lo.ToPtr(gatewayapi.NamespacesFromAll),
						},
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      uuid.NewString(),
		},
	}
	require.NoError(t, client.Create(ctx, &gw))

	gwOld := gw.DeepCopy()
	gw.Status = gatewayapi.GatewayStatus{
		Addresses: []gatewayapi.GatewayStatusAddress{
			{
				Type:  lo.ToPtr(gatewayapi.IPAddressType),
				Value: "10.0.0.1",
			},
		},
		Conditions: []metav1.Condition{
			{
				Type:               "Programmed",
				Status:             metav1.ConditionTrue,
				Reason:             "Programmed",
				LastTransitionTime: metav1.Now(),
				ObservedGeneration: gw.Generation,
			},
			{
				Type:               "Accepted",
				Status:             metav1.ConditionTrue,
				Reason:             "Accepted",
				LastTransitionTime: metav1.Now(),
				ObservedGeneration: gw.Generation,
			},
		},
		Listeners: []gatewayapi.ListenerStatus{
			{
				Name: "http",
				Conditions: []metav1.Condition{
					{
						Type:               "Accepted",
						Status:             metav1.ConditionTrue,
						Reason:             "Accepted",
						LastTransitionTime: metav1.Now(),
					},
					{
						Type:               "Programmed",
						Status:             metav1.ConditionTrue,
						Reason:             "Programmed",
						LastTransitionTime: metav1.Now(),
					},
				},
				SupportedKinds: []gatewayapi.RouteGroupKind{
					{
						Group: lo.ToPtr(gatewayapi.Group(gatewayv1.GroupVersion.Group)),
						Kind:  "HTTPRoute",
					},
				},
			},
		},
	}
	require.NoError(t, client.Status().Patch(ctx, &gw, ctrlclient.MergeFrom(gwOld)))

	route := gatewayapi.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: "v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nsRoute.Name,
			Name:      uuid.NewString(),
		},
		Spec: gatewayapi.HTTPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name:      gatewayapi.ObjectName(gw.Name),
					Namespace: lo.ToPtr(gatewayapi.Namespace(ns.Name)),
				}},
			},
			Rules: []gatewayapi.HTTPRouteRule{{
				BackendRefs: builder.NewHTTPBackendRef("backend-1").WithNamespace(ns.Name).ToSlice(),
			}},
		},
	}
	require.NoError(t, client.Create(ctx, &route))

	nn := k8stypes.NamespacedName{
		Namespace: route.GetNamespace(),
		Name:      route.GetName(),
	}

	t.Logf("verifying that HTTPRoute has ResolvedRefs set to Status False and Reason RefNotPermitted")
	if !assert.Eventually(t,
		helpers.HTTPRouteEventuallyContainsConditions(ctx, t, client, nn,
			metav1.Condition{
				Type:   "ResolvedRefs",
				Status: "False",
				Reason: "RefNotPermitted",
			},
		),
		waitDuration, tickDuration,
	) {
		t.Fatal(printHTTPRoutesConditions(ctx, client, nn))
	}

	rg := gatewayapi.ReferenceGrant{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      uuid.NewString(),
		},
		Spec: gatewayapi.ReferenceGrantSpec{
			From: []gatewayapi.ReferenceGrantFrom{
				{
					Group:     gatewayapi.Group(gatewayv1.GroupVersion.Group),
					Kind:      "HTTPRoute",
					Namespace: gatewayapi.Namespace(nsRoute.Name),
				},
			},
			To: []gatewayapi.ReferenceGrantTo{
				{
					Group: "",
					Kind:  "Service",
				},
			},
		},
	}
	require.NoError(t, client.Create(ctx, &rg))
	t.Logf("verifying that HTTPRoute gets accepted by HTTPRouteReconciler after relevant ReferenceGrant gets created")
	if !assert.Eventually(t,
		helpers.HTTPRouteEventuallyContainsConditions(ctx, t, client, nn,
			metav1.Condition{
				Type:   "ResolvedRefs",
				Status: "True",
				Reason: "ResolvedRefs",
			},
			metav1.Condition{
				Type:   "Accepted",
				Status: "True",
				Reason: "Accepted",
			},
			// Programmed condition requires a bit more work with mocks.
			// It's set only when KubernetesObjectReports are enabled in the underlying
			// dataplane client and then it relies on what's returned by
			// dataplane client in KubernetesObjectConfigurationStatus().
			// This can be done but it's not the main focus of this test.
			// Related: https://github.com/Kong/kubernetes-ingress-controller/issues/3793
		),
		waitDuration, tickDuration,
	) {
		t.Fatal(printHTTPRoutesConditions(ctx, client, nn))
	}

	require.NoError(t, client.Delete(ctx, &rg))
	t.Logf("verifying that HTTPRoute gets its ResolvedRefs condition to Status False and Reason RefNotPermitted when relevant ReferenceGrant gets deleted")

	if !assert.Eventually(t,
		helpers.HTTPRouteEventuallyContainsConditions(ctx, t, client, nn,
			metav1.Condition{
				Type:   "ResolvedRefs",
				Status: "False",
				Reason: "RefNotPermitted",
			},
		),
		waitDuration, tickDuration,
	) {
		t.Fatal(printHTTPRoutesConditions(ctx, client, nn))
	}
}

func TestHTTPRouteReconciler_RemovesOutdatedParentStatuses(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithGatewayAPI)
	cfg := Setup(t, scheme)
	client := NewControllerClient(t, scheme, cfg)

	reconciler := &gateway.HTTPRouteReconciler{
		Client:          client,
		DataplaneClient: mocks.Dataplane{},
	}

	// We use a deferred cancel to stop the manager and not wait for its timeout.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ns := CreateNamespace(ctx, t, client)
	nsRoute := CreateNamespace(ctx, t, client)

	StartReconcilers(ctx, t, client.Scheme(), cfg, reconciler)

	gwc := gatewayapi.GatewayClass{
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "kong-gwclass-",
			Annotations: map[string]string{
				"konghq.com/gatewayclass-unmanaged": "placeholder",
			},
		},
	}
	require.NoError(t, client.Create(ctx, &gwc))
	t.Cleanup(func() { _ = client.Delete(ctx, &gwc) })
	gwcNonKong := gatewayapi.GatewayClass{
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: "acme.com/dummy-controller",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "non-kong-gwclass-",
		},
	}
	require.NoError(t, client.Create(ctx, &gwcNonKong))
	t.Cleanup(func() { _ = client.Delete(ctx, &gwcNonKong) })

	gw := gatewayapi.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    ns.Name,
			GenerateName: "gw-kong-",
		},
		Spec: gatewayapi.GatewaySpec{
			GatewayClassName: gatewayapi.ObjectName(gwc.Name),
			Listeners: []gatewayapi.Listener{
				{
					Name:          gatewayapi.SectionName("http"),
					Port:          gatewayapi.PortNumber(80),
					Protocol:      gatewayapi.HTTPProtocolType,
					AllowedRoutes: builder.NewAllowedRoutesFromAllNamespaces(),
				},
			},
		},
	}
	require.NoError(t, client.Create(ctx, &gw))

	gwNonKong := gatewayapi.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "gw-nonkong-",
			Namespace:    ns.Name,
		},
		Spec: gatewayapi.GatewaySpec{
			GatewayClassName: gatewayapi.ObjectName(gwcNonKong.Name),
			Listeners: []gatewayapi.Listener{
				{
					Name:          gatewayapi.SectionName("http"),
					Port:          gatewayapi.PortNumber(80),
					Protocol:      gatewayapi.HTTPProtocolType,
					AllowedRoutes: builder.NewAllowedRoutesFromAllNamespaces(),
				},
			},
		},
	}
	require.NoError(t, client.Create(ctx, &gwNonKong))

	route := gatewayapi.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: "v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    nsRoute.Name,
			GenerateName: "httproute-kong-",
		},
		Spec: gatewayapi.HTTPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name:      gatewayapi.ObjectName(gw.Name),
					Namespace: lo.ToPtr(gatewayapi.Namespace(ns.Name)),
				}},
			},
			Rules: []gatewayapi.HTTPRouteRule{{
				BackendRefs: builder.NewHTTPBackendRef("backend-1").WithNamespace(ns.Name).ToSlice(),
			}},
		},
	}
	require.NoError(t, client.Create(ctx, &route))
	// Status has to be updated separately.
	route.Status = gatewayapi.HTTPRouteStatus{
		RouteStatus: gatewayapi.RouteStatus{
			Parents: []gatewayapi.RouteParentStatus{
				{
					ParentRef: gatewayapi.ParentReference{
						Name: gatewayapi.ObjectName(gwNonKong.Name),
					},
					ControllerName: gateway.GetControllerName(),
				},
			},
		},
	}
	require.NoError(t, client.Status().Update(ctx, &route))

	routeNonKong := gatewayapi.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: "v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    nsRoute.Name,
			GenerateName: "httproute-nonkong-",
		},
		Spec: gatewayapi.HTTPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name:      gatewayapi.ObjectName(gwNonKong.Name),
					Namespace: lo.ToPtr(gatewayapi.Namespace(ns.Name)),
				}},
			},
			Rules: []gatewayapi.HTTPRouteRule{{
				BackendRefs: builder.NewHTTPBackendRef("backend-1").WithNamespace(ns.Name).ToSlice(),
			}},
		},
	}
	require.NoError(t, client.Create(ctx, &routeNonKong))
	// Status has to be updated separately.
	routeNonKong.Status = gatewayapi.HTTPRouteStatus{
		RouteStatus: gatewayapi.RouteStatus{
			Parents: []gatewayapi.RouteParentStatus{
				{
					ParentRef: gatewayapi.ParentReference{
						Name: gatewayapi.ObjectName(gwNonKong.Name),
					},
					ControllerName: gwcNonKong.Spec.ControllerName,
				},
			},
		},
	}
	require.NoError(t, client.Status().Update(ctx, &routeNonKong))

	t.Run("routes attached to Gateways that are reconciled by KIC should have other Gateway refs cleared from status", func(t *testing.T) {
		require.Eventually(t, func() bool {
			if err := client.Get(ctx, ctrlclient.ObjectKeyFromObject(&route), &route); err != nil {
				t.Logf("failed to get HTTPRoute %s: %v", ctrlclient.ObjectKeyFromObject(&route), err)
				return false
			}

			if staleStatusFound := lo.ContainsBy(route.Status.Parents, func(ps gatewayapi.RouteParentStatus) bool {
				return string(ps.ParentRef.Name) == gwNonKong.Name
			}); staleStatusFound {
				t.Logf("found stale status for parent Gateway %s that does not belong to KIC GatewayClass", gwNonKong.Name)
				return false
			}

			return true
		}, waitDuration, tickDuration, "expected stale status to be removed from HTTPRoute")
	})

	t.Run("routes that were attached to Gateways that are reconciled by KIC and now become detached should have KIC Gateway refs cleared from status", func(t *testing.T) {
		require.Eventually(t, func() bool {
			if err := client.Get(ctx, ctrlclient.ObjectKeyFromObject(&route), &route); err != nil {
				t.Logf("failed to get HTTPRoute %s: %v", ctrlclient.ObjectKeyFromObject(&route), err)
				return false
			}
			route.Spec.ParentRefs = nil
			if err := client.Status().Update(ctx, &route); err != nil {
				t.Logf("failed to update HTTPRoute %s: %v", ctrlclient.ObjectKeyFromObject(&route), err)
				return false
			}

			if err := client.Get(ctx, ctrlclient.ObjectKeyFromObject(&route), &route); err != nil {
				t.Logf("failed to get HTTPRoute %s: %v", ctrlclient.ObjectKeyFromObject(&route), err)
				return false
			}

			return len(route.Status.Parents) == 0
		}, waitDuration, tickDuration, "expected stale KIC Gateway parentRef to be removed from HTTPRoute status")
	})

	t.Run("routes attached to Gateways that are not reconciled by KIC should not have other Gateway refs cleared from status", func(t *testing.T) {
		require.Never(t, func() bool {
			if err := client.Get(ctx, ctrlclient.ObjectKeyFromObject(&routeNonKong), &routeNonKong); err != nil {
				t.Logf("failed to get HTTPRoute %s: %v", ctrlclient.ObjectKeyFromObject(&routeNonKong), err)
				return true
			}

			if staleStatusFound := lo.ContainsBy(routeNonKong.Status.Parents, func(ps gatewayapi.RouteParentStatus) bool {
				return string(ps.ParentRef.Name) == gwNonKong.Name
			}); !staleStatusFound {
				t.Logf("status for parent %s not found, it should not have been cleared", gwNonKong.Name)
				return true
			}

			return false
		}, waitDuration, tickDuration, "expected status to not be removed from HTTPRoute")
	})
}

func printHTTPRoutesConditions(ctx context.Context, client ctrlclient.Client, nn k8stypes.NamespacedName) string {
	var route gatewayapi.HTTPRoute
	err := client.Get(ctx, ctrlclient.ObjectKey{Namespace: nn.Namespace, Name: nn.Name}, &route)
	if err != nil {
		return fmt.Sprintf("Failed to get HTTPRoute %s/%s when trying to print its conditions", nn.Namespace, nn.Name)
	}

	if len(route.Status.Parents) == 0 {
		return fmt.Sprintf("HTTPRoute %s/%s has no parents in Status", nn.Namespace, nn.Name)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("HTTPRoute %s/%s has the following Parents in Status:", nn.Namespace, nn.Name))
	for _, p := range route.Status.Parents {
		if p.ParentRef.Namespace != nil {
			_, _ = sb.WriteString(fmt.Sprintf("\nParent %s/%s: ", *p.ParentRef.Namespace, string(p.ParentRef.Name)))
		} else {
			_, _ = sb.WriteString(fmt.Sprintf("\nParent %s: ", string(p.ParentRef.Name)))
		}
		for _, c := range p.Conditions {
			s := fmt.Sprintf(
				"\n\tcondition: Type:%s, Status:%s, Reason:%s, ObservedGeneration:%d",
				c.Type, c.Status, c.Reason, c.ObservedGeneration,
			)
			_, _ = sb.WriteString(s)
		}
		_ = sb.WriteByte('\n')
	}
	return sb.String()
}
