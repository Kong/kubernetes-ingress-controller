package util

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type ControllerFunctionalPrometheusMetrics struct {
	// ConfigCounter number of post /config to proxy successfully
	ConfigCounter *prometheus.CounterVec

	// ParseCounter number of ingress analysis failure
	ParseCounter *prometheus.CounterVec

	// ConfigureDurationHistogram duration of last successful confiuration sync
	ConfigureDurationHistogram prometheus.Histogram
}

type Success string

const (
	// ConfigSuccessTrue post-config to proxy successfully
	ConfigSuccessTrue Success = "true"
	// ConfigSuccessFalse post-config to proxy failed
	ConfigSuccessFalse Success = "false"
	// IngressParseTrue says that ingress parsed successful
	IngressParseTrue Success = "true"
	// IngressParseFalse ingress parsed failed
	IngressParseFalse Success = "false"
)

type ConfigType string

const (
	// ConfigProxy says post config to proxy
	ConfigProxy ConfigType = "post-config"
	// ConfigDeck says generate deck
	ConfigDeck ConfigType = "deck"
)

func (ctrlMetrics *ControllerFunctionalPrometheusMetrics) NewPrometheusHistogram(name, help string) prometheus.Histogram {
	return prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    help,
			Buckets: prometheus.ExponentialBuckets(1, 10, 4),
		},
	)
}

func ControllerMetricsInit() *ControllerFunctionalPrometheusMetrics {
	controllerMetrics := &ControllerFunctionalPrometheusMetrics{}

	controllerMetrics.ConfigCounter =
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "send_configuration_count",
				Help: "number of post config proxy processed successfully.",
			},
			[]string{"success", "type"},
		)

	controllerMetrics.ParseCounter =
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingress_parse_count",
				Help: "number of ingress parse.",
			},
			[]string{"success"},
		)

	controllerMetrics.ConfigureDurationHistogram = controllerMetrics.NewPrometheusHistogram("proxy_configuration_duration_milliseconds", "duration of last successful configuration.")

	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		controllerMetrics.ConfigCounter,
		controllerMetrics.ParseCounter,
		controllerMetrics.ConfigureDurationHistogram,
	)

	return controllerMetrics
}
