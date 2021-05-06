//+build integration_tests

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/pkg/manager"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	_ = proxyReady()
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("http://localhost:%v/healthz", manager.HealthzPort))
		if err != nil {
			t.Logf("WARNING: error while waiting for http://localhost:%v/healthz: %v", manager.HealthzPort, err)
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
	_ = proxyReady()
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("http://localhost:%v/metrics", manager.MetricsPort))
		if err != nil {
			t.Logf("WARNING: error while waiting for http://localhost:%v/metrics: %v", manager.MetricsPort, err)
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
