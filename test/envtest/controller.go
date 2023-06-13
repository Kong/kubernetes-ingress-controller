package envtest

import (
	"bytes"
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

// StartReconcilers creates a controller manager and starts the provided reconciler
// as its runnable.
// It also adds a t.Cleanup which waits for the maanger to exit so that the test
// can be self contained and logs from different tests' managers don't mix up.
func StartReconcilers(ctx context.Context, t *testing.T, scheme *runtime.Scheme, cfg *rest.Config, reconcilers ...controllers.Reconciler) {
	t.Helper()

	var b bytes.Buffer
	log := logrus.New()
	log.Out = &b
	o := manager.Options{
		Logger:             logrusr.New(log),
		Scheme:             scheme,
		MetricsBindAddress: "0",
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
			t.Logf("Test %s failed: dumping controller logs\n%s", t.Name(), b.String())
		}
	})
}
