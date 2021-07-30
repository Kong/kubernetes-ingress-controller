package util

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type ControllerFunctionalPrometheusMetrics struct {
	// number of post /config to proxy successfully
	ConfigPassCounter prometheus.Counter
	// number of post /config to proxy failed
	ConfigFailureCounter prometheus.Counter
	// number of ingress analysis failure
	ParseFailureCounter prometheus.Counter

	// duration of last successful confiuration sync
	ConfigureDurationHistogram prometheus.Histogram
}

func (ctrlMetrics *ControllerFunctionalPrometheusMetrics) GeneratePrometheusCounter(name, help string) prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
	)
}

func (ctrlMetrics *ControllerFunctionalPrometheusMetrics) GeneratePrometheusHistogram(name, help string) prometheus.Histogram {
	return prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    help,
			Buckets: prometheus.LinearBuckets(20, 5, 5),
		},
	)
}

func ControllerMetricsInit() *ControllerFunctionalPrometheusMetrics {
	controllerMetrics := &ControllerFunctionalPrometheusMetrics{}

	controllerMetrics.ConfigPassCounter = controllerMetrics.GeneratePrometheusCounter("sendconfig_sync_proxy_success", "number of post config proxy processed successfully.")
	controllerMetrics.ConfigFailureCounter = controllerMetrics.GeneratePrometheusCounter("sendconfig_sync_proxy_failure", "number of post config proxy processed successfully.")
	controllerMetrics.ParseFailureCounter = controllerMetrics.GeneratePrometheusCounter("ingress_parse_failure", "number of ingress parse failed.")
	controllerMetrics.ConfigureDurationHistogram = controllerMetrics.GeneratePrometheusHistogram("sendconfig_sync_duration", "duration of last successful configuration.")

	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		controllerMetrics.ConfigPassCounter,
		controllerMetrics.ConfigFailureCounter,
		controllerMetrics.ParseFailureCounter,
		controllerMetrics.ConfigureDurationHistogram)

	return controllerMetrics
}
