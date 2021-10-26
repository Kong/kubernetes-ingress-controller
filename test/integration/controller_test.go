//go:build integration_tests
// +build integration_tests

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()
	assert.Eventually(t, func() bool {
		healthzURL := fmt.Sprintf("http://localhost:%v/healthz", manager.HealthzPort)
		resp, err := httpc.Get(healthzURL)
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", healthzURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, ingressWait, waitTick)
}

func TestReadyEndpoint(t *testing.T) {
	t.Parallel()
	assert.Eventually(t, func() bool {
		readyzURL := fmt.Sprintf("http://localhost:%v/readyz", manager.HealthzPort)
		resp, err := httpc.Get(readyzURL)
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", readyzURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, ingressWait, waitTick)
}

func TestProfilingEndpoint(t *testing.T) {
	t.Parallel()
	assert.Eventually(t, func() bool {
		profilingURL := fmt.Sprintf("http://localhost:%v/debug/pprof/", manager.DiagnosticsPort)
		resp, err := httpc.Get(profilingURL)
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", profilingURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, ingressWait, waitTick)
}

func TestConfigEndpoint(t *testing.T) {
	t.Parallel()
	assert.Eventually(t, func() bool {
		successURL := fmt.Sprintf("http://localhost:%v/debug/config/successful", manager.DiagnosticsPort)
		failURL := fmt.Sprintf("http://localhost:%v/debug/config/failed", manager.DiagnosticsPort)
		successResp, err := httpc.Get(successURL)
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", successURL, err)
			return false
		}
		defer successResp.Body.Close()
		failResp, err := httpc.Get(failURL)
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", failURL, err)
			return false
		}
		defer failResp.Body.Close()
		return successResp.StatusCode == http.StatusOK && failResp.StatusCode == http.StatusOK
	}, ingressWait, waitTick)
}
