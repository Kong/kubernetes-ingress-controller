package envtest

import (
	"context"
	"sync"
	"testing"

	"github.com/bombsimon/logrusr/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers"
)

// StartReconciler creates a controller manager and starts the provided reconciler
// as its runnable.
// It also adds a t.Cleanup which waits for the maanger to exit so that the test
// can be self contained and logs from different tests' managers don't mix up.
func StartReconciler(ctx context.Context, t *testing.T, scheme *runtime.Scheme, cfg *rest.Config, reconciler controllers.Reconciler) {
	o := manager.Options{
		Logger:             logrusr.New(logrus.New()),
		Scheme:             scheme,
		MetricsBindAddress: "0",
	}

	mgr, err := ctrl.NewManager(cfg, o)
	require.NoError(t, err)

	reconciler.SetLogger(mgr.GetLogger())

	require.NoError(t, reconciler.SetupWithManager(mgr))

	// This wait group makes it so that we wait for manager to exit.
	// This way we get clean test logs not mixing between tests.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.NoError(t, mgr.Start(ctx))
	}()
	t.Cleanup(func() { wg.Wait() })
}
