//go:build integration_tests
// +build integration_tests

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/internal/manager"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

func TestMetricsEndpoint(t *testing.T) {
	t.Parallel()
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
