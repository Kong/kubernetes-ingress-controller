//+build integration_tests

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"

	"github.com/kong/kubernetes-ingress-controller/railgun/internal/manager"
)

func TestHealthEndpoint(t *testing.T) {
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

func TestMetricsEndpoint(t *testing.T) {
	if useLegacyKIC() {
		t.Skip("metrics endpoint test does not apply to legacy KIC")
	}
	assert.Eventually(t, func() bool {
		metricsURL := fmt.Sprintf("http://localhost:%v/metrics", manager.MetricsPort)
		resp, err := httpc.Get(metricsURL)
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", metricsURL, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return false
		}
		decoder := expfmt.SampleDecoder{
			Dec:  expfmt.NewDecoder(resp.Body, expfmt.FmtText),
			Opts: &expfmt.DecodeOptions{},
		}

		var v model.Vector
		if err := decoder.Decode(&v); err != nil {
			t.Logf("decoder failed: %v", err)
			return false
		}

		return len(v) > 0
	}, ingressWait, waitTick)
}

func TestProfilingEndpoint(t *testing.T) {
	if useLegacyKIC() {
		t.Skip("profiling endpoint behaves differently in legacy KIC")
	}
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
