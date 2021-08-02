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
)

type ConfigType string

const (
	// ConfigProxy says post config to proxy
	ConfigProxy ConfigType = "post-config"
	// ConfigDeck says generate deck
	ConfigDeck ConfigType = "deck"
)

func ControllerMetricsInit() *CtrlFuncMetrics {
	controllerMetrics := &CtrlFuncMetrics{}

	reg := prometheus.NewRegistry()

	controllerMetrics.ConfigCounter =
		promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "send_configuration_count",
				Help: "Counts the events of sending configuration to Kong, using metric fields to distinguish between DB-less or DB-mode syncs, and to tell successes from failures..",
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
				Buckets: prometheus.ExponentialBuckets(1, 1.2, 20),
			},
		)

	return controllerMetrics
}
