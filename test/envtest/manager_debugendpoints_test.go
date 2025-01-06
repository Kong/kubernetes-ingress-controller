//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

func TestDebugEndpoints(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 10 * time.Second
		tickTime = 10 * time.Millisecond
	)

	ctx, cancel := context.WithCancel(context.Background())
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
		name string
		port int
	}{
		{
			port: diagPort,
			name: "debug/pprof/",
		},
		{
			port: diagPort,
			name: "debug/config/successful",
		},
		{
			port: diagPort,
			name: "debug/config/failed",
		},
		{
			port: healthPort,
			name: "healthz",
		},
		{
			port: healthPort,
			name: "readyz",
		},
	}

	for _, u := range urls {
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
