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
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

func TestHTTPRouteReconciliation_DoesNotBlockSyncLoopWhenStatusQueueBufferIsExceeded(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithGatewayAPI)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gw, _ := deployGateway(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(gw.Namespace),
		WithGatewayFeatureEnabled,
		WithGatewayAPIControllers(),
		func(cfg *manager.Config) {
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
			Ports: builder.NewServicePort().
				WithName("http").
				WithProtocol(corev1.ProtocolTCP).
				WithPort(80).
				IntoSlice(),
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, &backendService))
	t.Cleanup(func() { _ = ctrlClient.Delete(ctx, &backendService) })

	httpRoute := gatewayapi.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: gw.Namespace,
			Name:      uuid.NewString(),
		},
		Spec: gatewayapi.HTTPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name: gatewayapi.ObjectName(gw.Name),
				}},
			},
			Rules: []gatewayapi.HTTPRouteRule{{
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

	require.Eventually(t, httpRouteGetsProgrammed(ctx, t, ctrlClient, httpRoute),
		time.Second*10, time.Millisecond*50)
}

func httpRouteGetsProgrammed(ctx context.Context, t *testing.T, cl client.Client, httpRoute gatewayapi.HTTPRoute) func() bool {
	return func() bool {
		if err := cl.Get(ctx, client.ObjectKeyFromObject(&httpRoute), &httpRoute); err != nil {
			t.Logf("Failed to get httpRoute: %v", err)
			return false
		}
		if len(httpRoute.Status.Parents) == 0 {
			t.Logf("no gateway parent in httpRoute status")
			return false
		}
		programmed, ok := lo.Find(httpRoute.Status.Parents[0].Conditions, func(c metav1.Condition) bool {
			return c.Type == string(gatewayapi.GatewayConditionProgrammed)
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
	}
}

func Test_WatchNamespaces(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 10 * time.Second
		tickTime = 100 * time.Millisecond
	)

	scheme := Scheme(t, WithGatewayAPI)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gw, _ := deployGateway(ctx, t, ctrlClient)
	hidden := CreateNamespace(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(gw.Namespace),
		WithGatewayFeatureEnabled,
		WithGatewayAPIControllers(),
		func(cfg *manager.Config) {
			// Enable status updates and change the queue's buffer size to 0 to
			// ensure that the status update notifications do not block the
			// sync loop despite the fact that the status update queue is full.
			// This is a regression test for https://github.com/Kong/kubernetes-ingress-controller/issues/4260.
			// The test will timeout if the sync loop blocks (effectively manager.Run does not return after canceling context).
			cfg.UpdateStatus = true
			cfg.UpdateStatusQueueBufferSize = 0
			// hidden is excluded
			cfg.WatchNamespaces = []string{gw.Namespace}
		})

	backendService := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: gw.Namespace,
			Name:      "backend-svc",
		},
		Spec: corev1.ServiceSpec{
			Ports: builder.NewServicePort().
				WithName("http").
				WithProtocol(corev1.ProtocolTCP).
				WithPort(80).
				IntoSlice(),
		},
	}

	hiddenService := backendService
	hiddenService.Namespace = hidden.Name

	require.NoError(t, ctrlClient.Create(ctx, &backendService))
	t.Cleanup(func() { _ = ctrlClient.Delete(ctx, &backendService) })

	require.NoError(t, ctrlClient.Create(ctx, &hiddenService))
	t.Cleanup(func() { _ = ctrlClient.Delete(ctx, &hiddenService) })

	httpRoute := gatewayapi.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: gw.Namespace,
			Name:      uuid.NewString(),
		},
		Spec: gatewayapi.HTTPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name: gatewayapi.ObjectName(gw.Name),
				}},
			},
			Rules: []gatewayapi.HTTPRouteRule{{
				Matches: builder.NewHTTPRouteMatch().
					WithPathPrefix("/path").
					ToSlice(),
				BackendRefs: builder.NewHTTPBackendRef(backendService.Name).
					WithPort(80).
					ToSlice(),
			}},
		},
	}

	hiddenRoute := httpRoute
	hiddenRoute.Namespace = hidden.Name

	require.NoError(t, ctrlClient.Create(ctx, &httpRoute))
	t.Cleanup(func() { _ = ctrlClient.Delete(ctx, &httpRoute) })

	require.NoError(t, ctrlClient.Create(ctx, &hiddenRoute))
	t.Cleanup(func() { _ = ctrlClient.Delete(ctx, &hiddenRoute) })

	require.Eventually(t, httpRouteGetsProgrammed(ctx, t, ctrlClient, httpRoute),
		waitTime, tickTime)

	require.Never(t, func() bool {
		if err := ctrlClient.Get(ctx, client.ObjectKeyFromObject(&hiddenRoute), &hiddenRoute); err != nil {
			return false
		}
		if len(hiddenRoute.Status.Parents) != 0 {
			t.Logf("gateway parent assigned to ignored namespace HTTPRoute")
			return true
		}
		return false
	}, waitTime, tickTime)
}
