//go:build envtest

package envtest

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestSpecificGatewayNN(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 5 * time.Second
		tickTime = 500 * time.Millisecond
	)

	scheme := Scheme(t, WithGatewayAPI, WithKong)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		gw        = deployGateway(ctx, t, ctrlClient)
		nn        = client.ObjectKeyFromObject(&gw)
		gwIgnored = deployGateway(ctx, t, ctrlClient)
		nnIgnored = client.ObjectKeyFromObject(&gwIgnored)
	)

	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(gw.Namespace),
		WithGatewayFeatureEnabled,
		WithGatewayAPIControllers(),
		WithGatewayToReconcile(nn.String()),
	)

	createHTTPRoutes(ctx, t, ctrlClient, gw, 1)
	createHTTPRoutes(ctx, t, ctrlClient, gwIgnored, 1)

	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, nn, &gw)
		if err != nil {
			t.Logf("Failed to get gateway %s: %v", nn, err)
			return false
		}
		httpListener, ok := lo.Find(gw.Status.Listeners, func(listener gatewayapi.ListenerStatus) bool {
			return listener.Name == "http"
		})
		if !ok {
			t.Logf("Failed to find http listener status in gateway %s", nn)
			return false
		}
		if httpListener.AttachedRoutes != 1 {
			t.Logf("Expected %d routes to be attached to the http listener, got %d", 1, httpListener.AttachedRoutes)
			return false
		}

		return true
	}, waitTime, tickTime, "Failed to reconcile all HTTPRoutes")

	require.Never(t, func() bool {
		err := ctrlClient.Get(ctx, nnIgnored, &gwIgnored)
		if err != nil {
			t.Logf("Failed to get gateway %s: %v", nnIgnored, err)
			return true
		}

		// ignoredGW.Status.Listeners should be []
		if len(gwIgnored.Status.Listeners) != 0 {
			t.Logf("%s gateway should not be processed.", nnIgnored)
			return true
		}
		return false
	}, waitTime, tickTime, "Non configured gateway should not be processed")
}
