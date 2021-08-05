package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CtrlFuncMetrics struct {
	// ConfigCounter counts the events of sending configuration to Kong,
	// using metric fields to distinguish between DB-less or DB-mode syncs,
	// and to tell successes from failures.
	ConfigCounter *prometheus.CounterVec

	// ParseCounter counts the events of converting resources from Kubernetes to a KongState.
	ParseCounter *prometheus.CounterVec

	// ConfigureDurationHistogram records the duration of each successful configuration sync.
	ConfigureDurationHistogram prometheus.Histogram
}

// Success indicates the results of a function/operation
type Success string

const (
	// SuccessTrue operation successfully
	SuccessTrue Success = "true"
	// SuccessFalse operation failed
	SuccessFalse Success = "false"

	// SuccessKey success label within metrics
	SuccessKey Success = "success"
)

type ConfigType string

const (

	// ConfigProxy says post config to proxy
	ConfigProxy ConfigType = "post-config"
	// ConfigDeck says generate deck
	ConfigDeck ConfigType = "deck"

	// TypeKey type label within metrics
	TypeKey ConfigType = "type"
)

func ControllerMetricsInit() *CtrlFuncMetrics {
	controllerMetrics := &CtrlFuncMetrics{}

	reg := prometheus.NewRegistry()

	controllerMetrics.ConfigCounter =
		promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "send_configuration_count",
				Help: "Counts the success/failure events of converting kubernetes resources to a KongState, including conversion.",
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
				Help:    "Duration of last successful configuration.",
				Buckets: prometheus.ExponentialBuckets(100, 1.33, 30),
			},
		)

	return controllerMetrics
}
