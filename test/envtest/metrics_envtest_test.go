//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

func TestMetricsAreServed(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 200 * time.Second
		tickTime = 10 * time.Millisecond
		maxDelay = 500 * time.Millisecond
	)

	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme)

	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()

	cfg, _ := RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
	)

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
			require.NoError(t,
				retry.Do(func() error {
					resp, err := http.Get(metricsURL)
					if err != nil {
						return fmt.Errorf("error %w checking %q", err, metricsURL)
					}

					defer resp.Body.Close()
					if http.StatusOK != resp.StatusCode {
						return fmt.Errorf("status code %v not as expected (200)", resp.StatusCode)
					}

					var parser expfmt.TextParser
					mf, err := parser.TextToMetricFamilies(resp.Body)
					if err != nil {
						return fmt.Errorf("error %w parsing %q", err, metricsURL)
					}

					if _, ok := mf[metric]; !ok {
						return fmt.Errorf("metric %q not present yet", metric)
					}
					return nil
				},
					retry.Context(ctx),
					retry.Delay(tickTime),
					retry.MaxDelay(maxDelay),
					retry.MaxJitter(maxDelay),
					retry.DelayType(retry.BackOffDelay),
					retry.Attempts(0), // We're using a context with timeout, so we don't need to limit the number of attempts.
					retry.LastErrorOnly(true),
					retry.OnRetry(func(_ uint, err error) {
						t.Logf("metric %s not present yet, err: %v", metric, err.Error())
					}),
				),
			)

			t.Logf("metric %q is present at /metrics", metric)
		})
	}
}
