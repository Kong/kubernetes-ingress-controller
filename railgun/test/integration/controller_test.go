//+build integration_tests

package integration

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHealthEndpoint intentionally includes a magic value. localhost:10254 is the default
// controller health listen. This test should fail, and must be updated, if that default is changed.
func TestHealthEndpoint(t *testing.T) {
	_ = proxyReady()
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get("http://localhost:10254/healthz")
		if err != nil {
			t.Logf("WARNING: error while waiting for http://localhost:10254/healthz: %v", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			b.ReadFrom(resp.Body)
			return strings.Contains(b.String(), "ok")
		}
		return false
	}, ingressWait, waitTick)
}

// TestMetricsEndpoint intentionally includes a magic value. localhost:8080 is the default
// controller metrics listen. This test should fail, and must be updated, if that default is changed.
func TestMetricsEndpoint(t *testing.T) {
	_ = proxyReady()
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get("http://localhost:8080/metrics")
		if err != nil {
			t.Logf("WARNING: error while waiting for http://localhost:8080/metrics: %v", err)
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
