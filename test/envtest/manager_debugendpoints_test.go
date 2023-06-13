//go:build envtest
// +build envtest

package envtest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/cmd/rootcmd"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset/scheme"
)

func TestDebugEndpoints(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 10 * time.Second
		tickTime = 10 * time.Millisecond
	)

	envcfg := Setup(t, scheme.Scheme)
	cfg := ConfigForEnvConfig(t, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		// NOTE: We're not running rootcmd.Run() or rootcmd.RunWithLogger() here
		// hecause that sets up signal handling and that in turn uses a mutex to ensure
		// only one signal handler is running at a time.
		// We could try to work around this but that code calls os.Exit(1) whenever
		// the root context is cancelled and that not what we want to test here.

		deprecatedLogger, logger, err := manager.SetupLoggers(&cfg, io.Discard)
		require.NoError(t, err)
		diag, err := rootcmd.StartDiagnosticsServer(ctx, manager.DiagnosticsPort, &cfg, logger)
		require.NoError(t, err)

		err = manager.Run(ctx, &cfg, diag.ConfigDumps, deprecatedLogger)
		require.NoError(t, err)
	}(ctx)

	urls := []struct {
		name string
		port int
	}{
		{
			port: manager.DiagnosticsPort,
			name: "debug/pprof/",
		},
		{
			port: manager.DiagnosticsPort,
			name: "debug/config/successful",
		},
		{
			port: manager.DiagnosticsPort,
			name: "debug/config/failed",
		},
		{
			port: manager.HealthzPort,
			name: "healthz",
		},
		{
			port: manager.HealthzPort,
			name: "readyz",
		},
	}

	for _, u := range urls {
		u := u
		t.Run(u.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:%d/%s", u.port, u.name)
			eventuallHTTPGet(t, http.DefaultClient, url, waitTime, tickTime)
		})
	}
}

func eventuallHTTPGet(t *testing.T, h *http.Client, url string, waitTime, tickTime time.Duration) {
	t.Helper()

	t.Logf("HTTP GET %s", url)
	assert.EventuallyWithT(t, func(t *assert.CollectT) {
		resp, err := h.Get(url)
		if !assert.NoError(t, err) {
			return
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}, waitTime, tickTime)
}
