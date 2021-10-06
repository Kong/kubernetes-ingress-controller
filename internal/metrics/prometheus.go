package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type CtrlFuncMetrics struct {
	// ConfigPushCount counts the events of sending configuration to Kong,
	// using metric fields to distinguish between DB-less or DB-mode syncs,
	// and to tell successes from failures.
	ConfigPushCount *prometheus.CounterVec

	// TranslationCount counts the events of converting resources from Kubernetes to a KongState.
	TranslationCount *prometheus.CounterVec

	// ConfigPushDuration records the duration of each successful configuration sync.
	ConfigPushDuration prometheus.Histogram
}

// Success indicates the results of a function/operation
type Success string

const (
	// SuccessTrue operation successfully
	SuccessTrue Success = "true"
	// SuccessFalse operation failed
	SuccessFalse Success = "false"

	// SuccessKey success label within metrics
	SuccessKey model.LabelName = "success"
)

type ConfigType string

const (

	// ConfigDBLess says post config to proxy
	ConfigDBLess ConfigType = "db-less"
	// ConfigDeck says generate deck
	ConfigDeck ConfigType = "deck"

	// TypeKey type label within metrics
	TypeKey model.LabelName = "type"
)

func NewCtrlFuncMetrics() *CtrlFuncMetrics {
	controllerMetrics := &CtrlFuncMetrics{}

	controllerMetrics.ConfigPushCount =
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingress_controller_configuration_push_count",
				Help: "Count of successful/failed configuration pushes to Kong. `" +
					TypeKey + "` describes the configuration protocol (" + ConfigDBLess + " or " +
					ConfigDeck + ") in use. `" +
					SuccessKey + "` describes whether there were unrecoverable errors (`" +
					SuccessFalse + "`) or not (`" + SuccessTrue + "`).",
			},
			[]string{"success", "type"},
		)

	controllerMetrics.TranslationCount =
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingress_controller_translation_count",
				Help: "Count of translations from Kubernetes state to Kong state. `" +
					SuccessKey + "` describes whether there were unrecoverable errors (`" +
					SuccessFalse + "`) or not (`" + SuccessTrue + "`).",
			},
			[]string{"success"},
		)

	controllerMetrics.ConfigPushDuration =
		prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name: "ingress_controller_configuration_push_duration_milliseconds",
				Help: "How long it took to push the configuration to Kong, in milliseconds. `" +
					TypeKey + "` describes the configuration protocol (" + ConfigDBLess + " or " +
					ConfigDeck + ") in use. `" +
					SuccessKey + "` describes whether there were unrecoverable errors (`" +
					SuccessFalse + "`) or not (`" + SuccessTrue + "`).",
				Buckets: prometheus.ExponentialBuckets(100, 1.33, 30),
			},
			[]string{"success", "type"},
		)

	metrics.Registry.MustRegister(controllerMetrics.ConfigPushCount, controllerMetrics.TranslationCount, controllerMetrics.ConfigPushDuration)

	return controllerMetrics
}
