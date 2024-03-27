//go:build envtest

package envtest

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
	gatewayclientv1beta1 "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned/typed/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers/conditions"
	"github.com/kong/kubernetes-ingress-controller/v2/test/mocks"
)

func TestGatewayWithGatewayClassReconciliation(t *testing.T) {
	t.Parallel()

	const (
		// unsupportedControllerName is the name of the controller used for those
		// gateways that are not supported by an actual controller (i.e., they won't be scheduled).
		unsupportedControllerName gatewayv1.GatewayController = "example.com/unsupported-gateway-controller"
	)

	scheme := Scheme(t, WithGatewayAPI)
	cfg := Setup(t, scheme)

	gatewayClient, err := gatewayclient.NewForConfig(cfg)
	require.NoError(t, err)

	client := NewControllerClient(t, scheme, cfg)

	testcases := []struct {
		Name         string
		GatewayClass gatewayv1.GatewayClass
		Gateway      gatewayv1.Gateway
		Test         func(
			ctx context.Context,
			t *testing.T,
			gwClient gatewayclientv1beta1.GatewayInterface,
			gwc gatewayv1.GatewayClass,
			gw gatewayv1.Gateway,
		)
	}{
		{
			Name: "unsupported gateway class",
			GatewayClass: gatewayv1.GatewayClass{
				Spec: gatewayv1.GatewayClassSpec{
					ControllerName: unsupportedControllerName,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "unsupported-gateway-class",
				},
			},
			Gateway: gatewayv1.Gateway{
				Spec: gatewayv1.GatewaySpec{
					GatewayClassName: gatewayv1.ObjectName("unsupported-gateway-class"),
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
				gwClient gatewayclientv1beta1.GatewayInterface,
				gwc gatewayv1.GatewayClass,
				gw gatewayv1.Gateway,
			) {
				t.Logf("deploying gateway class %s", gwc.Name)
				require.NoError(t, client.Create(ctx, &gwc))
				t.Cleanup(func() { _ = client.Delete(ctx, &gwc) })

				t.Logf("verifying that the unsupported Gateway %s does not get Accepted or Programmed by the controller", gw.Name)
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
			},
		},
		{
			Name: "managed gateway class",
			GatewayClass: gatewayv1.GatewayClass{
				Spec: gatewayv1.GatewayClassSpec{
					ControllerName: gateway.GetControllerName(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "managed-gateway-class",
				},
			},
			Gateway: gatewayv1.Gateway{
				Spec: gatewayv1.GatewaySpec{
					GatewayClassName: gatewayv1.ObjectName("managed-gateway-class"),
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
				gwClient gatewayclientv1beta1.GatewayInterface,
				gwc gatewayv1.GatewayClass,
				gw gatewayv1.Gateway,
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
				t.Cleanup(func() { _ = client.Delete(ctx, &gwc) })

				// Let's wait and check that the Gateway hasn't been reconciled by the operator.
				t.Log("verifying the Gateway is not reconciled as it is using a managed GatewayClass")
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
			},
		},
		{
			Name: "unmanaged gateway class",
			GatewayClass: gatewayv1.GatewayClass{
				Spec: gatewayv1.GatewayClassSpec{
					ControllerName: gateway.GetControllerName(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "unmanaged-gateway-class",
					Annotations: map[string]string{
						annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
					},
				},
			},
			Gateway: gatewayv1.Gateway{
				Spec: gatewayv1.GatewaySpec{
					GatewayClassName: gatewayv1.ObjectName("unmanaged-gateway-class"),
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
				gwClient gatewayclientv1beta1.GatewayInterface,
				gwc gatewayv1.GatewayClass,
				gw gatewayv1.Gateway,
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
				t.Cleanup(func() { _ = client.Delete(ctx, &gwc) })

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
						gateway, err := gwClient.Get(ctx, gw.Name, metav1.GetOptions{})
						require.NoError(t, err)
						t.Logf("expected to find an Accepted and Programmed conditions with Status True, got %#v", gateway.Status.Conditions)
						t.Fatalf("context got cancelled: %v", ctx.Err())
					case event := <-w.ResultChan():
						gateway, ok := event.Object.(*gatewayv1.Gateway)
						require.True(t, ok)

						if !conditions.Contain(gateway.Status.Conditions, conditions.WithType("Programmed"), conditions.WithStatus(metav1.ConditionTrue)) {
							continue
						}
						if !conditions.Contain(gateway.Status.Conditions, conditions.WithType("Accepted"), conditions.WithStatus(metav1.ConditionTrue)) {
							continue
						}
						break forLoop
					}
				}
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
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
			t.Cleanup(func() { _ = client.Delete(ctx, &tc.Gateway) })

			gwClient := gatewayClient.GatewayV1beta1().Gateways(ns.Name)
			tc.Test(ctx, t, gwClient, tc.GatewayClass, tc.Gateway)
		})
	}
}
