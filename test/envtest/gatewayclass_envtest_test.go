//go:build envtest

package envtest

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
	gatewayclientv1 "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned/typed/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/conditions"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func TestGatewayWithGatewayClassReconciliation(t *testing.T) {
	t.Parallel()

	const (
		// unsupportedControllerName is the name of the controller used for those
		// gateways that are not supported by an actual controller (i.e., they won't be scheduled).
		unsupportedControllerName gatewayapi.GatewayController = "example.com/unsupported-gateway-controller"

		waitTime = 3 * time.Second
		tickTime = 100 * time.Millisecond
	)

	scheme := Scheme(t, WithGatewayAPI)
	cfg := Setup(t, scheme)

	gatewayClient, err := gatewayclient.NewForConfig(cfg)
	require.NoError(t, err)

	client := NewControllerClient(t, scheme, cfg)

	testcases := []struct {
		Name         string
		GatewayClass gatewayapi.GatewayClass
		Gateway      gatewayapi.Gateway
		Test         func(
			ctx context.Context,
			t *testing.T,
			gwClient gatewayclientv1.GatewayInterface,
			gwc gatewayapi.GatewayClass,
			gw gatewayapi.Gateway,
		)
	}{
		{
			Name: "unsupported gateway class",
			GatewayClass: gatewayapi.GatewayClass{
				Spec: gatewayapi.GatewayClassSpec{
					ControllerName: unsupportedControllerName,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "unsupported-gateway-class",
				},
			},
			Gateway: gatewayapi.Gateway{
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: gatewayapi.ObjectName("unsupported-gateway-class"),
					Listeners: builder.NewListener("http").
						HTTP().
						WithPort(80).
						WithAllowedRoutes(builder.NewAllowedRoutesFromAllNamespaces()).
						IntoSlice(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
			},
			Test: func(
				ctx context.Context,
				t *testing.T,
				gwClient gatewayclientv1.GatewayInterface,
				gwc gatewayapi.GatewayClass,
				gw gatewayapi.Gateway,
			) {
				t.Logf("deploying gateway class %s", gwc.Name)
				require.NoError(t, client.Create(ctx, &gwc))
				t.Cleanup(func() { _ = client.Delete(context.Background(), &gwc) }) //nolint:contextcheck

				t.Logf("verifying that the unsupported Gateway %s does not get Accepted or Programmed by the controller", gw.Name)
				// NOTE: Ideally we wouldn't like to perform a busy wait loop here,
				// but rely on actual data like number of Reconciler calls.
				// However, this is currently not possible because the controllers
				// we have pass themselves as the Reconciler in manager.Options
				// hence wrapping the Reconciler() method is impossible with current
				// implementation.
				// Related issue: https://github.com/Kong/kubernetes-ingress-controller/issues/4190
				//
				// NOTE: we're not using a busy loop without a sleep because that could cause
				// the rate limiter to kick in and fail the test.
				require.Never(t, func() bool {
					gateway, err := gwClient.Get(ctx, gw.Name, metav1.GetOptions{})
					if err != nil {
						t.Logf("error getting Gateway %s: %v", gateway.Name, err)
						return true
					}
					if len(gateway.Status.Conditions) != 2 {
						t.Logf("expected 2 Status Conditions on the Gateway: Accepted and Programmed, got: %v", gateway.Status.Conditions)
						return true
					}

					if conditions.Contain(gateway.Status.Conditions, conditions.WithType("Programmed"), conditions.Not(conditions.WithStatus(metav1.ConditionUnknown))) {
						t.Logf("expected not to find a Programmed condition with Reason different than Unknown, got %#v", gateway.Status.Conditions)
						return true
					}
					if conditions.Contain(gateway.Status.Conditions, conditions.WithType("Accepted"), conditions.Not(conditions.WithStatus(metav1.ConditionUnknown))) {
						t.Logf("expected not to find a Accepted condition with Reason different than Unknown, got %#v", gateway.Status.Conditions)
						return true
					}
					return false
				}, waitTime, tickTime)
			},
		},
		{
			Name: "managed gateway class",
			GatewayClass: gatewayapi.GatewayClass{
				Spec: gatewayapi.GatewayClassSpec{
					ControllerName: gateway.GetControllerName(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "managed-gateway-class",
				},
			},
			Gateway: gatewayapi.Gateway{
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: gatewayapi.ObjectName("managed-gateway-class"),
					Listeners: builder.NewListener("http").
						HTTP().
						WithPort(80).
						WithAllowedRoutes(builder.NewAllowedRoutesFromAllNamespaces()).
						IntoSlice(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
			},
			Test: func(
				ctx context.Context,
				t *testing.T,
				gwClient gatewayclientv1.GatewayInterface,
				gwc gatewayapi.GatewayClass,
				gw gatewayapi.Gateway,
			) {
				t.Logf("verifying that the Gateway %s does not get scheduled by the controller due to missing its GatewayClass", gw.Name)
				// NOTE: Ideally we wouldn't like to perform a busy wait loop here,
				// but rely on actual data like number of Reconciler calls.
				// However, this is currently not possible because the controllers
				// we have pass themselves as the Reconciler in manager.Options
				// hence wrapping the Reconciler() method is impossible with current
				// implementation.
				// Related issue: https://github.com/Kong/kubernetes-ingress-controller/issues/4190
				for i := 0; i < 100; i++ {
					gateway, err := gwClient.Get(ctx, gw.Name, metav1.GetOptions{})
					require.NoError(t, err)
					require.Len(t, gateway.Status.Conditions, 2)

					if conditions.Contain(gateway.Status.Conditions, conditions.WithType("Programmed"), conditions.Not(conditions.WithStatus(metav1.ConditionUnknown))) {
						t.Fatalf("expected not to find a Programmed condition with Reason different than Unknown, got %#v", gateway.Status.Conditions)
					}
					if conditions.Contain(gateway.Status.Conditions, conditions.WithType("Accepted"), conditions.Not(conditions.WithStatus(metav1.ConditionUnknown))) {
						t.Fatalf("expected not to find a Accepted condition with Reason different than Unknown, got %#v", gateway.Status.Conditions)
					}
				}

				t.Logf("deploying gateway class %s", gwc.Name)
				require.NoError(t, client.Create(ctx, &gwc))
				t.Cleanup(func() { _ = client.Delete(context.Background(), &gwc) }) //nolint:contextcheck

				// Let's wait and check that the Gateway hasn't been reconciled by the operator.
				t.Log("verifying the Gateway is not reconciled as it is using a managed GatewayClass")

				// NOTE: we're not using a busy loop without a sleep because that could cause
				// the rate limiter to kick in and fail the test.
				require.Never(t, func() bool {
					gateway, err := gwClient.Get(ctx, gw.Name, metav1.GetOptions{})
					if err != nil {
						t.Logf("error getting Gateway %s: %v", gateway.Name, err)
						return true
					}
					if len(gateway.Status.Conditions) != 2 {
						t.Logf("expected 2 Status Conditions on the Gateway: Accepted and Programmed, got: %v", gateway.Status.Conditions)
						return true
					}

					if conditions.Contain(gateway.Status.Conditions, conditions.WithType("Programmed"), conditions.Not(conditions.WithStatus(metav1.ConditionUnknown))) {
						t.Logf("expected not to find a Programmed condition with Reason different than Unknown, got %#v", gateway.Status.Conditions)
						return true
					}
					if conditions.Contain(gateway.Status.Conditions, conditions.WithType("Accepted"), conditions.Not(conditions.WithStatus(metav1.ConditionUnknown))) {
						t.Logf("expected not to find a Accepted condition with Reason different than Unknown, got %#v", gateway.Status.Conditions)
						return true
					}
					return false
				}, waitTime, tickTime)
			},
		},
		{
			Name: "unmanaged gateway class",
			GatewayClass: gatewayapi.GatewayClass{
				Spec: gatewayapi.GatewayClassSpec{
					ControllerName: gateway.GetControllerName(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "unmanaged-gateway-class",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.GatewayClassUnmanagedKey: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
					},
				},
			},
			Gateway: gatewayapi.Gateway{
				Spec: gatewayapi.GatewaySpec{
					GatewayClassName: gatewayapi.ObjectName("unmanaged-gateway-class"),
					Listeners: builder.NewListener("http").
						HTTP().
						WithPort(80).
						WithAllowedRoutes(builder.NewAllowedRoutesFromAllNamespaces()).
						IntoSlice(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
			},
			Test: func(
				ctx context.Context,
				t *testing.T,
				gwClient gatewayclientv1.GatewayInterface,
				gwc gatewayapi.GatewayClass,
				gw gatewayapi.Gateway,
			) {
				t.Logf("verifying that the Gateway %s does not get scheduled by the controller due to missing its GatewayClass", gw.Name)
				// NOTE: Ideally we wouldn't like to perform a busy wait loop here,
				// but rely on actual data like number of Reconciler calls.
				// However, this is currently not possible because the controllers
				// we have pass themselves as the Reconciler in manager.Options
				// hence wrapping the Reconciler() method is impossible with current
				// implementation.
				// Related issue: https://github.com/Kong/kubernetes-ingress-controller/issues/4190
				//
				// NOTE: we're not using a busy loop without a sleep because that could cause
				// the rate limiter to kick in and fail the test.
				require.Never(t, func() bool {
					gateway, err := gwClient.Get(ctx, gw.Name, metav1.GetOptions{})
					if err != nil {
						t.Logf("error getting Gateway %s: %v", gateway.Name, err)
						return true
					}
					if len(gateway.Status.Conditions) != 2 {
						t.Logf("expected 2 Status Conditions on the Gateway: Accepted and Programmed, got: %v", gateway.Status.Conditions)
						return true
					}

					if conditions.Contain(gateway.Status.Conditions, conditions.WithType("Programmed"), conditions.Not(conditions.WithStatus(metav1.ConditionUnknown))) {
						t.Logf("expected not to find a Programmed condition with Reason different than Unknown, got %#v", gateway.Status.Conditions)
						return true
					}
					if conditions.Contain(gateway.Status.Conditions, conditions.WithType("Accepted"), conditions.Not(conditions.WithStatus(metav1.ConditionUnknown))) {
						t.Logf("expected not to find a Accepted condition with Reason different than Unknown, got %#v", gateway.Status.Conditions)
						return true
					}
					return false
				}, waitTime, tickTime)

				t.Logf("deploying gateway class %s", gwc.Name)
				require.NoError(t, client.Create(ctx, &gwc))
				t.Cleanup(func() { _ = client.Delete(context.Background(), &gwc) }) //nolint:contextcheck

				t.Logf("now that the GatewayClass exists, verifying that the Gateway %s gets Accepted and Programmed", gw.Name)

				w, err := gwClient.Watch(ctx, metav1.ListOptions{
					FieldSelector: "metadata.name=" + gw.Name,
					TypeMeta: metav1.TypeMeta{
						APIVersion: gatewayv1.GroupVersion.String(),
						Kind:       "Gateway",
					},
				})
				require.NoError(t, err)
				defer w.Stop()

			forLoop:
				for {
					select {
					case <-ctx.Done():
						t.Fatalf("context got cancelled: %v", ctx.Err())
					case event := <-w.ResultChan():
						gateway, ok := event.Object.(*gatewayapi.Gateway)
						require.True(t, ok, "expected to get a Gateway object, got %T", event.Object)

						if !conditions.Contain(gateway.Status.Conditions, conditions.WithType("Programmed"), conditions.WithStatus(metav1.ConditionTrue)) {
							t.Logf("Gateway %s still doesn't have the Programmed condition with Status True", gateway.Name)
							continue
						}
						if !conditions.Contain(gateway.Status.Conditions, conditions.WithType("Accepted"), conditions.WithStatus(metav1.ConditionTrue)) {
							t.Logf("Gateway %s still doesn't have the Accepted condition with Status True", gateway.Name)
							continue
						}
						break forLoop
					}
				}
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			var (
				ctx    context.Context
				cancel func()
			)

			// We use a deferred cancel to stop the manager and not wait for its timeout.
			if deadline, ok := t.Deadline(); ok {
				ctx, cancel = context.WithDeadline(context.Background(), deadline)
			} else {
				ctx, cancel = context.WithCancel(context.Background())
			}
			defer cancel()

			ns := CreateNamespace(ctx, t, client)

			svc := corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "publish-svc",
				},
				Spec: corev1.ServiceSpec{
					Ports: builder.NewServicePort().
						WithName("http").
						WithPort(80).
						WithProtocol(corev1.ProtocolTCP).
						WithTargetPort(intstr.FromInt(80)).
						IntoSlice(),
				},
			}
			require.NoError(t, client.Create(ctx, &svc))

			// GatewayReconciler runs GatewayClassReconciler, so we only need to
			// start the former.
			gwReconciler := &gateway.GatewayReconciler{
				Client: client,
				PublishServiceRef: k8stypes.NamespacedName{
					Namespace: ns.Name,
					Name:      svc.Name,
				},
				DataplaneClient:   mocks.Dataplane{},
				ReferenceIndexers: ctrlref.NewCacheIndexers(logr.Discard()),
			}
			StartReconcilers(ctx, t, client.Scheme(), cfg, gwReconciler)

			t.Logf("deploying gateway %s using %s gateway class", tc.Gateway.Name, tc.GatewayClass.Name)
			tc.Gateway.Namespace = ns.Name
			require.NoError(t, client.Create(ctx, &tc.Gateway))
			t.Cleanup(func() { _ = client.Delete(context.Background(), &tc.Gateway) })

			gwClient := gatewayClient.GatewayV1().Gateways(ns.Name)
			tc.Test(ctx, t, gwClient, tc.GatewayClass, tc.Gateway)
		})
	}
}
