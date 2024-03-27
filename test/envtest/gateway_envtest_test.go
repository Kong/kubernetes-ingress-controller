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
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// TestGatewayReconciliation_MoreThan100Routes verifies that if we create more
// than 100 HTTPRoutes, they all get reconciled and correctly attached to a
// Gateway's listener.
// It reproduces https://github.com/Kong/kubernetes-ingress-controller/issues/4456.
func TestGatewayReconciliation_MoreThan100Routes(t *testing.T) {
	scheme := Scheme(t, WithGatewayAPI, WithKong)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gw := deployGateway(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg, WithPublishService(gw.Namespace), WithGatewayFeatureEnabled)

	const numOfRoutes = 120
	createHTTPRoutes(ctx, t, ctrlClient, gw, numOfRoutes)

	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{Namespace: gw.Namespace, Name: gw.Name}, &gw)
		if err != nil {
			t.Logf("failed to get gateway %s/%s: %v", gw.Namespace, gw.Name, err)
			return false
		}
		httpListener, ok := lo.Find(gw.Status.Listeners, func(listener gatewayv1.ListenerStatus) bool {
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
	}, 3*time.Minute, time.Second, "failed to reconcile all HTTPRoutes")
}

// createHTTPRoutes creates a number of dummy HTTPRoutes for the given Gateway.
func createHTTPRoutes(
	ctx context.Context,
	t *testing.T,
	ctrlClient ctrlclient.Client,
	gw gatewayv1.Gateway,
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
		httpPort := gatewayv1.PortNumber(80)
		pathMatchPrefix := gatewayv1.PathMatchPathPrefix
		httpRoute := &gatewayv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: gw.Namespace,
			},
			Spec: gatewayv1.HTTPRouteSpec{
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{{
						Name: gatewayv1.ObjectName(gw.Name),
					}},
				},
				Rules: []gatewayv1.HTTPRouteRule{{
					Matches: []gatewayv1.HTTPRouteMatch{
						{
							Path: &gatewayv1.HTTPPathMatch{
								Type:  &pathMatchPrefix,
								Value: kong.String("/test-http-route"),
							},
						},
					},
					BackendRefs: []gatewayv1.HTTPBackendRef{{
						BackendRef: gatewayv1.BackendRef{
							BackendObjectReference: gatewayv1.BackendObjectReference{
								Name: gatewayv1.ObjectName("backend-svc"),
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
