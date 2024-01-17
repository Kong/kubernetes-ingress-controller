//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestSpecificGatewayNN(t *testing.T) {
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
	ignoredGW := deployGateway(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(gw.Namespace),
		WithGatewayFeatureEnabled,
		WithGatewayAPIControllers(),
		WithGatewayToReconcile(fmt.Sprintf("%s/%s", gw.Namespace, gw.Name)),
	)

	createHTTPRoutes(ctx, t, ctrlClient, gw, 1)
	createHTTPRoutes(ctx, t, ctrlClient, ignoredGW, 1)

	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{Namespace: gw.Namespace, Name: gw.Name}, &gw)
		if err != nil {
			t.Logf("Failed to get gateway %s/%s: %v", gw.Namespace, gw.Name, err)
			return false
		}
		httpListener, ok := lo.Find(gw.Status.Listeners, func(listener gatewayapi.ListenerStatus) bool {
			return listener.Name == "http"
		})
		if !ok {
			t.Logf("Failed to find http listener status in gateway %s/%s", gw.Namespace, gw.Name)
			return false
		}
		if httpListener.AttachedRoutes != 1 {
			t.Logf("Expected %d routes to be attached to the http listener, got %d", 1, httpListener.AttachedRoutes)
			return false
		}

		err = ctrlClient.Get(ctx, k8stypes.NamespacedName{Namespace: ignoredGW.Namespace, Name: ignoredGW.Name}, &ignoredGW)
		if err != nil {
			t.Logf("Failed to get gateway %s/%s: %v", ignoredGW.Namespace, ignoredGW.Name, err)
			return false
		}

		// ignoredGW.Status.Listeners should be []
		if len(ignoredGW.Status.Listeners) != 0 {
			t.Logf("Expected %s/%s gateway should not be processed.", ignoredGW.Namespace, ignoredGW.Name)
			return false
		}

		return true
	}, waitTime, tickTime, "failed to reconcile all HTTPRoutes")
}
