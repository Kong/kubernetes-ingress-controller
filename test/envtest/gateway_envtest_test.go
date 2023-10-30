//go:build envtest

package envtest

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// TestGatewayReconciliation_MoreThan100Routes verifies that if we create more
// than 100 HTTPRoutes, they all get reconciled and correctly attached to a
// Gateway's listener.
// It reproduces https://github.com/Kong/kubernetes-ingress-controller/issues/4456.
func TestGatewayReconciliation_MoreThan100Routes(t *testing.T) {
	t.Parallel()

	const (
		waitTime = time.Minute
		tickTime = 500 * time.Millisecond
	)

	scheme := Scheme(t, WithGatewayAPI, WithKong)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gw := deployGateway(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(gw.Namespace),
		WithGatewayFeatureEnabled,
		WithGatewayAPIControllers(),
	)

	const numOfRoutes = 120
	createHTTPRoutes(ctx, t, ctrlClient, gw, numOfRoutes)

	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{Namespace: gw.Namespace, Name: gw.Name}, &gw)
		if err != nil {
			t.Logf("failed to get gateway %s/%s: %v", gw.Namespace, gw.Name, err)
			return false
		}
		httpListener, ok := lo.Find(gw.Status.Listeners, func(listener gatewayapi.ListenerStatus) bool {
			return listener.Name == "http"
		})
		if !ok {
			t.Logf("failed to find http listener status in gateway %s/%s", gw.Namespace, gw.Name)
			return false
		}
		if httpListener.AttachedRoutes != numOfRoutes {
			t.Logf("expected %d routes to be attached to the http listener, got %d", numOfRoutes, httpListener.AttachedRoutes)
			return false
		}
		return true
	}, waitTime, tickTime, "failed to reconcile all HTTPRoutes")
}

// createHTTPRoutes creates a number of dummy HTTPRoutes for the given Gateway.
func createHTTPRoutes(
	ctx context.Context,
	t *testing.T,
	ctrlClient ctrlclient.Client,
	gw gatewayapi.Gateway,
	numOfRoutes int,
) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "backend-svc",
			Namespace: gw.Namespace,
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
	require.NoError(t, ctrlClient.Create(ctx, svc))
	t.Cleanup(func() { _ = ctrlClient.Delete(ctx, svc) })

	for i := 0; i < numOfRoutes; i++ {
		httpPort := gatewayapi.PortNumber(80)
		pathMatchPrefix := gatewayapi.PathMatchPathPrefix
		httpRoute := &gatewayapi.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: gw.Namespace,
			},
			Spec: gatewayapi.HTTPRouteSpec{
				CommonRouteSpec: gatewayapi.CommonRouteSpec{
					ParentRefs: []gatewayapi.ParentReference{{
						Name: gatewayapi.ObjectName(gw.Name),
					}},
				},
				Rules: []gatewayapi.HTTPRouteRule{{
					Matches: []gatewayapi.HTTPRouteMatch{
						{
							Path: &gatewayapi.HTTPPathMatch{
								Type:  &pathMatchPrefix,
								Value: kong.String("/test-http-route"),
							},
						},
					},
					BackendRefs: []gatewayapi.HTTPBackendRef{{
						BackendRef: gatewayapi.BackendRef{
							BackendObjectReference: gatewayapi.BackendObjectReference{
								Name: gatewayapi.ObjectName("backend-svc"),
								Port: &httpPort,
								Kind: util.StringToGatewayAPIKindPtr("Service"),
							},
						},
					}},
				}},
			},
		}

		err := ctrlClient.Create(ctx, httpRoute)
		require.NoError(t, err)
		t.Cleanup(func() { _ = ctrlClient.Delete(ctx, httpRoute) })
	}
}
