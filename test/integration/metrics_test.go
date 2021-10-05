//go:build integration_tests
// +build integration_tests

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/internal/manager"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

func TestMetricsEndpoint(t *testing.T) {
	t.Parallel()

	wantMetrics := []string{
		"send_configuration_count",
		"ingress_parse_count",
		"proxy_configuration_duration_milliseconds",
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

		var parser expfmt.TextParser
		mf, err := parser.TextToMetricFamilies(resp.Body)
		if err != nil {
			t.Logf("WARNING: error when decoding prometheus metrics: %v", err)
			return false
		}

		for _, wantMetric := range wantMetrics {
			if _, ok := mf[wantMetric]; !ok {
				t.Logf("WARNING: metric not found in /metrics: %q", wantMetric)
				return false
			}
		}

		return true // All metrics from wantMetrics have been found in /metrics.
	}, ingressWait, waitTick)
}
