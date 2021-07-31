package util

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type ControllerFunctionalPrometheusMetrics struct {
	// number of post /config to proxy successfully
	ConfigCounter *prometheus.CounterVec

	// number of ingress analysis failure
	ParseCounter *prometheus.CounterVec

	// duration of last successful confiuration sync
	ConfigureDurationHistogram prometheus.Histogram
}

type Success string

const (
	// EnablementStatusDisabled says that the resource it controls is disabled.
	ConfigSuccessTrue Success = "true"
	// EnablementStatusEnabled says that the resource it controls is enabled.
	ConfigSuccessFalse Success = "false"

	IngressParseTrue  Success = "true"
	IngressParseFalse Success = "false"
)

type ConfigType string

const (
	ConfigProxy ConfigType = "post-config"
	ConfigDeck  ConfigType = "deck"
)

func (ctrlMetrics *ControllerFunctionalPrometheusMetrics) NewPrometheusCounter(name, help string, labels ...string) *prometheus.CounterVec {
	return 
	)
}

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
			[]string{"sucess", "type"},
		)

	controllerMetrics.ParseCounter =
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingress_parse_count",
				Help: "number of ingress parse.",
			},
			[]string{"sucess"},
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
