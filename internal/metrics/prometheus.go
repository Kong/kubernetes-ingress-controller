package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
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
	ConfigPushDuration *prometheus.HistogramVec
}

const (
	// SuccessTrue operation successfully
	SuccessTrue string = "true"
	// SuccessFalse operation failed
	SuccessFalse string = "false"

	// SuccessKey success label within metrics
	SuccessKey string = "success"
)

const (
	// ConfigDBLess says post config to proxy
	ConfigDBLess string = "db-less"
	// ConfigDeck says generate deck
	ConfigDeck string = "deck"

	// TypeKey type label within metrics
	TypeKey string = "type"
)

const (
	MetricNameConfigPushCount    = "ingress_controller_configuration_push_count"
	MetricNameTranslationCount   = "ingress_controller_translation_count"
	MetricNameConfigPushDuration = "ingress_controller_configuration_push_duration_milliseconds"
)

func NewCtrlFuncMetrics() *CtrlFuncMetrics {
	controllerMetrics := &CtrlFuncMetrics{}

	controllerMetrics.ConfigPushCount =
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: MetricNameConfigPushCount,
				Help: "Count of successful/failed configuration pushes to Kong. `" +
					TypeKey + "` describes the configuration protocol (" + ConfigDBLess + " or " +
					ConfigDeck + ") in use. `" +
					SuccessKey + "` describes whether there were unrecoverable errors (`" +
					SuccessFalse + "`) or not (`" + SuccessTrue + "`).",
			},
			[]string{SuccessKey, TypeKey},
		)

	controllerMetrics.TranslationCount =
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: MetricNameTranslationCount,
				Help: "Count of translations from Kubernetes state to Kong state. `" +
					SuccessKey + "` describes whether there were unrecoverable errors (`" +
					SuccessFalse + "`) or not (`" + SuccessTrue + "`).",
			},
			[]string{SuccessKey},
		)

	controllerMetrics.ConfigPushDuration =
		prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: MetricNameConfigPushDuration,
				Help: "How long it took to push the configuration to Kong, in milliseconds. `" +
					TypeKey + "` describes the configuration protocol (" + ConfigDBLess + " or " +
					ConfigDeck + ") in use. `" +
					SuccessKey + "` describes whether there were unrecoverable errors (`" +
					SuccessFalse + "`) or not (`" + SuccessTrue + "`).",
				Buckets: prometheus.ExponentialBuckets(100, 1.33, 30),
			},
			[]string{SuccessKey, TypeKey},
		)

	metrics.Registry.MustRegister(controllerMetrics.ConfigPushCount, controllerMetrics.TranslationCount, controllerMetrics.ConfigPushDuration)

	return controllerMetrics
}
