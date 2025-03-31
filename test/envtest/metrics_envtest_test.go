//go:build envtest

package envtest

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	prom "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func TestMetricsAreServed(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 200 * time.Second
		tickTime = 10 * time.Millisecond
	)

	allExpectedMetrics := []string{
		metrics.MetricNameConfigPushCount,
		metrics.MetricNameConfigPushBrokenResources,
		metrics.MetricNameTranslationCount,
		metrics.MetricNameTranslationDuration,
		metrics.MetricNameConfigPushSize,
		metrics.MetricNameTranslationBrokenResources,
		metrics.MetricNameConfigPushDuration,

		metrics.MetricNameFallbackTranslationBrokenResources,
		metrics.MetricNameFallbackTranslationCount,
		metrics.MetricNameFallbackTranslationDuration,
		metrics.MetricNameFallbackConfigPushSize,
		metrics.MetricNameFallbackConfigPushCount,
		metrics.MetricNameFallbackConfigPushSuccessTime,
		metrics.MetricNameFallbackConfigPushDuration,
		metrics.MetricNameFallbackConfigPushBrokenResources,
		metrics.MetricNameProcessedConfigSnapshotCacheHit,
		metrics.MetricNameProcessedConfigSnapshotCacheMiss,
	}

	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme)

	// We're going to make the first request to the admin API fail, so that the manager falls back to the last valid
	// configuration and we can test the fallback metrics.
	adminAPIOpts := []mocks.AdminAPIHandlerOpt{
		mocks.WithConfigPostError([]byte(`{"flattened_errors": [{"errors": [{"messages": ["broken object"]}], "entity_tags": ["k8s-name:test-service","k8s-namespace:default","k8s-kind:Service","k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3","k8s-group:","k8s-version:v1"]}]}`)),
		mocks.WithConfigPostErrorOnlyOnFirstRequest(),
	}
	addr := fmt.Sprintf("localhost:%d", helpers.GetFreePort(t))
	_ = RunManager(t.Context(), t, envcfg,
		AdminAPIOptFns(adminAPIOpts...),
		func(cfg *managercfg.Config) {
			cfg.FeatureGates[managercfg.FallbackConfigurationFeature] = true
		},
		WithMetricsAddr(addr),
	)

	metricsURL := fmt.Sprintf("http://%s/metrics", addr)
	t.Logf("waiting for metrics to be available at %q", metricsURL)

	assertMetric := func(t *assert.CollectT, reportedMetrics []*prom.MetricFamily, expectedMetricName string) {
		m, ok := lo.Find(reportedMetrics, func(m *prom.MetricFamily) bool {
			return m.GetName() == expectedMetricName
		})
		if !assert.True(t, ok, "expected metric %q not found in the response", expectedMetricName) {
			return
		}
		for _, m := range m.GetMetric() {
			containsInstanceID := lo.ContainsBy(m.GetLabel(), func(l *prom.LabelPair) bool {
				return l.GetName() == metrics.InstanceIDKey && l.GetValue() != ""
			})
			assert.True(t, containsInstanceID, "metric %q does not contain instance id label", expectedMetricName)
		}
	}
	httpClient := http.Client{
		Timeout: waitTime, // Set a timeout to avoid hanging forever.
	}
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		resp, err := httpClient.Get(metricsURL)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		reportedMetrics := extractAllMetrics(t, resp.Body)
		for _, expectedMetricName := range allExpectedMetrics {
			assertMetric(t, reportedMetrics, expectedMetricName)
		}
	}, waitTime, tickTime)
}

func extractAllMetrics(t *assert.CollectT, body io.Reader) []*prom.MetricFamily {
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(body)
	require.NoError(t, err)
	return lo.Values(mf)
}
