//go:build envtest

package envtest

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
)

func TestHTTPRouteReconciliation_DoesNotBlockSyncLoopWhenStatusQueueBufferIsExceeded(t *testing.T) {
	scheme := Scheme(t, WithGatewayAPI)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gw := deployGateway(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg, WithPublishService(gw.Namespace), WithGatewayFeatureEnabled, func(cfg *manager.Config) {
		// Enable status updates and change the queue's buffer size to 0 to
		// ensure that the status update notifications do not block the
		// sync loop despite the fact that the status update queue is full.
		// This is a regression test for https://github.com/Kong/kubernetes-ingress-controller/issues/4260.
		// The test will timeout if the sync loop blocks (effectively manager.Run does not return after canceling context).
		cfg.UpdateStatus = true
		cfg.UpdateStatusQueueBufferSize = 0
	})

	backendService := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: gw.Namespace,
			Name:      "backend-svc",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     80,
				},
			},
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, &backendService))
	t.Cleanup(func() { _ = ctrlClient.Delete(ctx, &backendService) })

	httpRoute := gatewayv1beta1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: gw.Namespace,
			Name:      uuid.NewString(),
		},
		Spec: gatewayv1beta1.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
				ParentRefs: []gatewayv1beta1.ParentReference{{
					Name: gatewayv1beta1.ObjectName(gw.Name),
				}},
			},
			Rules: []gatewayv1beta1.HTTPRouteRule{{
				Matches: builder.NewHTTPRouteMatch().
					WithPathPrefix("/path").
					ToSlice(),
				BackendRefs: builder.NewHTTPBackendRef(backendService.Name).
					WithPort(80).
					ToSlice(),
			}},
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, &httpRoute))
	t.Cleanup(func() { _ = ctrlClient.Delete(ctx, &httpRoute) })

	require.Eventually(t, func() bool {
		nn := k8stypes.NamespacedName{Name: httpRoute.Name, Namespace: httpRoute.Namespace}
		if err := ctrlClient.Get(ctx, nn, &httpRoute); err != nil {
			t.Logf("failed to get httpRoute: %v", err)
			return false
		}
		if len(httpRoute.Status.Parents) == 0 {
			t.Logf("no gateway parent in httpRoute status")
			return false
		}
		programmed, ok := lo.Find(httpRoute.Status.Parents[0].Conditions, func(c metav1.Condition) bool {
			return c.Type == string(gatewayv1beta1.GatewayConditionProgrammed)
		})
		if !ok {
			t.Logf("no programmed condition in httpRoute status")
			return false
		}
		if programmed.Status != metav1.ConditionTrue {
			t.Logf("programmed condition is not true")
			return false
		}
		return true
	}, time.Second*10, time.Millisecond*50)
	t.Cleanup(func() { _ = ctrlClient.Delete(ctx, &httpRoute) })
}

// deployGateway deploys a Gateway, GatewayClass, and publish Service for use in tests.
func deployGateway(ctx context.Context, t *testing.T, client client.Client) gatewayv1beta1.Gateway {
	ns := CreateNamespace(ctx, t, client)

	publishSvc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      PublishServiceName,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     8000,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
	require.NoError(t, client.Create(ctx, &publishSvc))
	t.Cleanup(func() { _ = client.Delete(ctx, &publishSvc) })

	gwc := gatewayv1beta1.GatewayClass{
		Spec: gatewayv1beta1.GatewayClassSpec{
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

	gw := gatewayv1beta1.Gateway{
		Spec: gatewayv1beta1.GatewaySpec{
			GatewayClassName: gatewayv1beta1.ObjectName(gwc.Name),
			Listeners: []gatewayv1beta1.Listener{
				{
					Name:     "http",
					Protocol: gatewayv1beta1.HTTPProtocolType,
					Port:     gatewayv1beta1.PortNumber(8000),
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      uuid.NewString(),
		},
	}
	require.NoError(t, client.Create(ctx, &gw))
	t.Cleanup(func() { _ = client.Delete(ctx, &gw) })

	return gw
}
