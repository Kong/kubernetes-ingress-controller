//+build integration_tests

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/railgun/manager"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	_ = proxyReady()
	assert.Eventually(t, func() bool {
		healthzURL := fmt.Sprintf("http://localhost:%v/healthz", manager.HealthzPort)
		resp, err := httpc.Get(healthzURL)
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", healthzURL, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return true
		}
		return false
	}, ingressWait, waitTick)
}

func TestMetricsEndpoint(t *testing.T) {
	if useLegacyKIC() {
		// The metrics endpoint was intentionally changed for 2.0. Skip if legacy.
		return
	}
	_ = proxyReady()
	assert.Eventually(t, func() bool {
		metricsURL := fmt.Sprintf("http://localhost:%v/metrics", manager.MetricsPort)
		resp, err := httpc.Get(metricsURL)
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", metricsURL, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			b.ReadFrom(resp.Body)
			return strings.Contains(b.String(), "controller_runtime_active_workers")
		}
		return false
	}, ingressWait, waitTick)
}
