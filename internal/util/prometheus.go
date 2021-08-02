package util

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ControllerFunctionalPrometheusMetrics struct {
	// ConfigCounter number of post /config to proxy successfully
	ConfigCounter *prometheus.CounterVec

	// ParseCounter number of ingress analysis failure
	ParseCounter *prometheus.CounterVec

	// ConfigureDurationHistogram duration of last successful confiuration sync
	ConfigureDurationHistogram prometheus.Histogram
}

// Success indicates the results of a function/operation
type Success string

const (
	// SuccessTrue operation successfully
	SuccessTrue Success = "true"
	// SuccessFalse operation failed
	SuccessFalse Success = "false"
)

type ConfigType string

const (
	// ConfigProxy says post config to proxy
	ConfigProxy ConfigType = "post-config"
	// ConfigDeck says generate deck
	ConfigDeck ConfigType = "deck"
)

func ControllerMetricsInit() *ControllerFunctionalPrometheusMetrics {
	controllerMetrics := &ControllerFunctionalPrometheusMetrics{}

	reg := prometheus.NewRegistry()

	controllerMetrics.ConfigCounter =
		promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "send_configuration_count",
				Help: "number of post config proxy processed successfully.",
			},
			[]string{"success", "type"},
		)

	controllerMetrics.ParseCounter =
		promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingress_parse_count",
				Help: "number of ingress parse.",
			},
			[]string{"success"},
		)

	controllerMetrics.ConfigureDurationHistogram =
		promauto.With(reg).NewHistogram(
			prometheus.HistogramOpts{
				Name:    "proxy_configuration_duration_milliseconds",
				Help:    "duration of last successful configuration.",
				Buckets: prometheus.ExponentialBuckets(1, 1.2, 20),
			},
		)

	return controllerMetrics
}
