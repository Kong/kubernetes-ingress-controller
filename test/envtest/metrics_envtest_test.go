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

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
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

	testCases := []struct {
		name                         string
		withPushError                bool
		fallbackConfigurationEnabled bool
		expectedMetrics              []string
		skippedMessage               string
	}{
		{
			name:                         "with push error and FallbackConfiguration enabled",
			skippedMessage:               "flaky, see https://github.com/Kong/kubernetes-ingress-controller/issues/6125",
			withPushError:                true,
			fallbackConfigurationEnabled: true,
			expectedMetrics: []string{
				metrics.MetricNameConfigPushCount,
				metrics.MetricNameConfigPushBrokenResources,
				metrics.MetricNameTranslationCount,
				metrics.MetricNameTranslationDuration,
				metrics.MetricNameTranslationBrokenResources,
				metrics.MetricNameConfigPushDuration,

				metrics.MetricNameFallbackTranslationBrokenResources,
				metrics.MetricNameFallbackTranslationCount,
				metrics.MetricNameFallbackConfigPushCount,
				metrics.MetricNameFallbackConfigPushSuccessTime,
				metrics.MetricNameFallbackConfigPushDuration,
				metrics.MetricNameFallbackConfigPushBrokenResources,
				metrics.MetricNameProcessedConfigSnapshotCacheHit,
				metrics.MetricNameProcessedConfigSnapshotCacheMiss,
			},
		},
		{
			name:          "without push error",
			withPushError: false,
			expectedMetrics: []string{
				metrics.MetricNameConfigPushCount,
				metrics.MetricNameConfigPushBrokenResources,
				metrics.MetricNameTranslationCount,
				metrics.MetricNameTranslationDuration,
				metrics.MetricNameTranslationBrokenResources,
				metrics.MetricNameConfigPushDuration,
				metrics.MetricNameConfigPushSuccessTime,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skippedMessage != "" {
				t.Skip(tc.skippedMessage)
			}
			ctx, cancel := context.WithTimeout(context.Background(), waitTime)
			defer cancel()

			var adminAPIOpts []mocks.AdminAPIHandlerOpt
			if tc.withPushError {
				adminAPIOpts = append(adminAPIOpts,
					mocks.WithConfigPostError([]byte(`{"flattened_errors": [{"errors": [{"messages": ["broken object"]}], "entity_tags": ["k8s-name:test-service","k8s-namespace:default","k8s-kind:Service","k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3","k8s-group:","k8s-version:v1"]}]}`)),
					mocks.WithConfigPostErrorOnlyOnFirstRequest(),
				)
			}
			cfg, _ := RunManager(ctx, t, envcfg,
				AdminAPIOptFns(adminAPIOpts...),
				func(cfg *manager.Config) {
					cfg.FeatureGates[featuregates.FallbackConfiguration] = tc.fallbackConfigurationEnabled
				},
			)

			wantMetrics := tc.expectedMetrics

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
		})
	}
}
