//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

func TestDebugEndpoints(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 10 * time.Second
		tickTime = 10 * time.Millisecond
	)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	diagPort := helpers.GetFreePort(t)
	healthPort := helpers.GetFreePort(t)
	envcfg := Setup(t, scheme.Scheme)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithDiagnosticsServer(diagPort),
		WithHealthProbePort(healthPort),
		WithProfiling(),
	)

	urls := []struct {
		name         string
		port         int
		expectedCode int
	}{
		{
			port:         diagPort,
			name:         "debug/pprof/",
			expectedCode: http.StatusOK,
		},
		{
			port:         diagPort,
			name:         "debug/config/successful",
			expectedCode: http.StatusOK,
		},
		{
			port:         diagPort,
			name:         "debug/config/failed",
			expectedCode: http.StatusNoContent, // No failed push, so no content.
		},
		{
			port:         healthPort,
			name:         "healthz",
			expectedCode: http.StatusOK,
		},
		{
			port:         healthPort,
			name:         "readyz",
			expectedCode: http.StatusOK,
		},
	}

	for _, u := range urls {
		t.Run(u.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:%d/%s", u.port, u.name)
			eventuallHTTPGet(t, http.DefaultClient, url, u.expectedCode, waitTime, tickTime)
		})
	}
}

func eventuallHTTPGet(t *testing.T, h *http.Client, url string, expectedCode int, waitTime, tickTime time.Duration) {
	t.Helper()

	t.Logf("HTTP GET %s", url)
	assert.EventuallyWithT(t, func(t *assert.CollectT) {
		resp, err := h.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, expectedCode, resp.StatusCode)
	}, waitTime, tickTime)
}
