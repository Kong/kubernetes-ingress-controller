//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset/scheme"
)

func TestMetricsAreServed(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 30 * time.Second
		tickTime = 10 * time.Millisecond
	)

	envcfg := Setup(t, scheme.Scheme)
	cfg := ConfigForEnvConfig(t, envcfg)
	cfg.EnableProfiling = false
	cfg.EnableConfigDumps = false

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		// NOTE: We're not running rootcmd.Run() or rootcmd.RunWithLogger() here
		// hecause that sets up signal handling and that in turn uses a mutex to ensure
		// only one signal handler is running at a time.
		// We could try to work around this but that code calls os.Exit(1) whenever
		// the root context is cancelled and that not what we want to test here.

		logger, err := manager.SetupLoggers(&cfg, io.Discard)
		require.NoError(t, err)

		err = manager.Run(ctx, &cfg, util.ConfigDumpDiagnostic{}, logger)
		require.NoError(t, err)
	}(ctx)

	wantMetrics := []string{
		metrics.MetricNameConfigPushCount,
		metrics.MetricNameConfigPushBrokenResources,
		metrics.MetricNameTranslationCount,
		metrics.MetricNameTranslationBrokenResources,
		metrics.MetricNameConfigPushDuration,
		metrics.MetricNameConfigPushSuccessTime,
	}

	metricsURL := fmt.Sprintf("http://%s/metrics", cfg.MetricsAddr)
	t.Logf("waiting for metrics to be available at %q", metricsURL)

	for _, metric := range wantMetrics {
		metric := metric
		t.Run(metric, func(t *testing.T) {
			require.Eventually(t, func() bool {
				resp, err := http.Get(metricsURL)
				if err != nil {
					t.Logf("err %v checking %q", err, metricsURL)
					return false
				}

				defer resp.Body.Close()
				if http.StatusOK != resp.StatusCode {
					t.Logf("status code %v", resp.StatusCode)
					return false
				}

				var parser expfmt.TextParser
				mf, err := parser.TextToMetricFamilies(resp.Body)
				if err != nil {
					t.Logf("err %v parsing %q", err, metricsURL)
					return false
				}

				if _, ok := mf[metric]; !ok {
					return false
				}
				return true
			}, waitTime, tickTime)
			t.Logf("metric %q is present at /metrics", metric)
		})
	}
}
