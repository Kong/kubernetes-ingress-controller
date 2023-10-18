package envtest

import (
	"context"
	"sync"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers"
)

// StartReconcilers creates a controller manager and starts the provided reconciler
// as its runnable.
// It also adds a t.Cleanup which waits for the manager to exit so that the test
// can be self contained and logs from different tests' managers don't mix up.
func StartReconcilers(ctx context.Context, t *testing.T, scheme *runtime.Scheme, cfg *rest.Config, reconcilers ...controllers.Reconciler) {
	t.Helper()

	core, logs := observer.New(zap.InfoLevel)
	logger := zapr.NewLogger(zap.New(core))
	o := manager.Options{
		Logger: logger,
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
	}

	mgr, err := ctrl.NewManager(cfg, o)
	require.NoError(t, err)

	for _, r := range reconcilers {
		r.SetLogger(mgr.GetLogger())
		require.NoError(t, r.SetupWithManager(mgr))
	}

	// This wait group makes it so that we wait for manager to exit.
	// This way we get clean test logs not mixing between tests.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.NoError(t, mgr.Start(ctx))
	}()
	t.Cleanup(func() {
		wg.Wait()

		if t.Failed() {
			t.Logf("Test %s failed: dumping controller logs\n", t.Name())
			for _, log := range logs.All() {
				t.Logf("%s %s\n", log.Time, log.Message)
			}
		}
	})
}

// NewControllerClient returns a new controller-runtime Client for provided runtime.Scheme and rest.Config.
func NewControllerClient(t *testing.T, scheme *runtime.Scheme, cfg *rest.Config) ctrlclient.Client {
	client, err := ctrlclient.New(cfg, ctrlclient.Options{
		Scheme: scheme,
	})
	require.NoError(t, err)
	return client
}
