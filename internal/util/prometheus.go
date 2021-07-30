package util

import (
	"strconv"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type ControllerMetrics struct {
	// number of post /config to proxy successfully
	ConfigPass uint64
	// number of post /config to proxy failed
	ConfigFailure uint64
	// number of ingress analysis failure
	IngressParseFailure uint64
	// number of deck message failure
	DeckGenerationFailure uint64
	// duration of last successful confiuration sync
	ConfigureDuration uint64
}

func (ctrlMetrics *ControllerMetrics) IncCounter(counter *uint64) {
	atomic.AddUint64(counter, 1)
}

func (ctrlMetrics *ControllerMetrics) SetCounter(counter *uint64, val uint64) {
	atomic.StoreUint64(counter, val)
}

func (ctrlMetrics *ControllerMetrics) GeneratePrometheusCounter(name, help string, labels map[string]string) prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:        name,
			Help:        help,
			ConstLabels: labels,
		},
	)
}

func (ctrlMetrics *ControllerMetrics) GeneratePrometheusHistogram(name, help string, labels map[string]string) prometheus.Histogram {
	return prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        name,
			Help:        help,
			ConstLabels: labels,
		},
	)
}

func ControllerMetricsInit() *ControllerMetrics {
	controllerMetrics := &ControllerMetrics{
		ConfigPass:            0,
		ConfigFailure:         0,
		IngressParseFailure:   0,
		DeckGenerationFailure: 0,
		ConfigureDuration:     0,
	}

	configPassCounter := controllerMetrics.GeneratePrometheusCounter("sendconfig_sync_proxy_success", "number of post config proxy processed successfully.",
		map[string]string{"config_pass_total": strconv.FormatUint(controllerMetrics.ConfigPass, 10)})

	configFailureCounter := controllerMetrics.GeneratePrometheusCounter("sendconfig_sync_proxy_failure", "number of post config proxy processed successfully.",
		map[string]string{"config_failure_total": strconv.FormatUint(controllerMetrics.ConfigFailure, 10)})

	ingressParseFailure := controllerMetrics.GeneratePrometheusCounter("ingress_parse_failure", "number of ingress parse failed.",
		map[string]string{"ingress_parse_failure_total": strconv.FormatUint(controllerMetrics.IngressParseFailure, 10)})

	deckFailure := controllerMetrics.GeneratePrometheusCounter("deckgen_failure", "number of deck communication failed.",
		map[string]string{"deck_generation_failure_total": strconv.FormatUint(controllerMetrics.DeckGenerationFailure, 10)})

	configurationSyncDuration := controllerMetrics.GeneratePrometheusHistogram("sendconfig_sync_duration", "duration of last successful configuration.",
		map[string]string{"configuration_sync_duration": strconv.FormatUint(controllerMetrics.ConfigureDuration, 10)})

	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(configPassCounter, configFailureCounter, ingressParseFailure, deckFailure, configurationSyncDuration)
	return controllerMetrics
}
