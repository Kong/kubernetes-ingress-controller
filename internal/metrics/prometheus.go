package metrics

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
)

// descriptions of these metrics are found below, where their help text is set in NewCtrlFuncMetrics()

type CtrlFuncMetrics struct {
	ConfigPushCount *prometheus.CounterVec

	ConfigPushBrokenResources *prometheus.GaugeVec

	TranslationCount *prometheus.CounterVec

	TranslationBrokenResources prometheus.Gauge

	ConfigPushDuration *prometheus.HistogramVec

	ConfigPushSuccessTime *prometheus.GaugeVec
}

const (
	// SuccessTrue indicates that the operation was successful.
	SuccessTrue string = "true"
	// SuccessFalse indicates that the operation was not successful.
	SuccessFalse string = "false"

	// SuccessKey defines the key of the metric label indicating success/failure of an operation.
	SuccessKey string = "success"
)

type Protocol string

const (
	// ProtocolDBLess indicates that configuration was sent to Kong using the DB-less protocol (POST /config).
	ProtocolDBLess Protocol = "db-less"
	// ProtocolDeck indicates that configuration was sent to Kong using the DB mode protocol (deck sync).
	ProtocolDeck Protocol = "deck"

	// ProtocolKey defines the key of the metric label indicating which protocol KIC used to configure Kong.
	ProtocolKey string = "protocol"
)

const (
	// FailureReasonConflict indicates that the config push failed due to configuration conflicts.
	FailureReasonConflict string = "conflict"

	// FailureReasonNetwork indicates that the config push failed due to network issues.
	FailureReasonNetwork string = "network"

	// FailureReasonOther indicates that the config push failed due to other reasons.
	FailureReasonOther string = "other"

	// FailureReasonKey defines the key of the metric label indicating failure reason.
	FailureReasonKey string = "failure_reason"
)

const (
	// DataplaneKey defines the name of the metric label indicating which dataplane this time series is relevant for.
	DataplaneKey string = "dataplane"
)

const (
	MetricNameConfigPushCount            = "ingress_controller_configuration_push_count"
	MetricNameConfigPushBrokenResources  = "ingress_controller_configuration_push_broken_resource_count"
	MetricNameConfigPushSuccessTime      = "ingress_controller_configuration_push_last_successful"
	MetricNameTranslationCount           = "ingress_controller_translation_count"
	MetricNameTranslationBrokenResources = "ingress_controller_translation_broken_resource_count"
	MetricNameConfigPushDuration         = "ingress_controller_configuration_push_duration_milliseconds"
)

var _lock sync.Mutex

func NewCtrlFuncMetrics() *CtrlFuncMetrics {
	_lock.Lock()
	defer _lock.Unlock()

	controllerMetrics := &CtrlFuncMetrics{}

	controllerMetrics.ConfigPushCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricNameConfigPushCount,
			Help: fmt.Sprintf(
				"Count of successful/failed configuration pushes to Kong. "+
					"`%s` describes the dataplane that was the target of configuration push. "+
					"`%s` describes the configuration protocol (`%s` or `%s`) in use. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`). "+
					"`%s` is populated in case of `%s=\"%s\"` and describes the reason of failure "+
					"(one of `%s`, `%s`, `%s`).",
				DataplaneKey,
				ProtocolKey, ProtocolDBLess, ProtocolDeck,
				SuccessKey, SuccessFalse, SuccessTrue,
				FailureReasonKey, SuccessKey, SuccessFalse,
				FailureReasonConflict, FailureReasonNetwork, FailureReasonOther,
			),
		},
		[]string{SuccessKey, ProtocolKey, FailureReasonKey, DataplaneKey},
	)

	controllerMetrics.ConfigPushBrokenResources = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameConfigPushBrokenResources,
			Help: fmt.Sprintf("The number of resources not accepted by Kong when attempting to push "+
				"configuration. `%s` describes the dataplane that was the target of the configuration push.",
				DataplaneKey,
			),
		},
		[]string{DataplaneKey},
	)

	controllerMetrics.TranslationCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricNameTranslationCount,
			Help: fmt.Sprintf(
				"Count of translations from Kubernetes state to Kong state. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`). "+
					"Unrecoverable error in this case means KIC wasn't able to translate a Kubernetes object to Kong model.",
				SuccessKey, SuccessFalse, SuccessTrue,
			),
		},
		[]string{SuccessKey},
	)

	controllerMetrics.TranslationBrokenResources = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: MetricNameTranslationBrokenResources,
			Help: fmt.Sprintf("The number of resources that the controller cannot successfully translate to Kong " +
				"configuration",
			),
		},
	)

	controllerMetrics.ConfigPushDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: MetricNameConfigPushDuration,
			Help: fmt.Sprintf(
				"How long it took to push the configuration to Kong, in milliseconds. "+
					"`%s` describes the dataplane that was the target of configuration push. "+
					"`%s` describes the configuration protocol (`%s` or `%s`) in use. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`).",
				DataplaneKey,
				ProtocolKey, ProtocolDBLess, ProtocolDeck,
				SuccessKey, SuccessFalse, SuccessTrue,
			),
			Buckets: prometheus.ExponentialBuckets(100, 1.33, 30),
		},
		[]string{SuccessKey, ProtocolKey, DataplaneKey},
	)

	controllerMetrics.ConfigPushSuccessTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameConfigPushSuccessTime,
			Help: fmt.Sprintf("The time of the last successful configuration push. "+
				"`%s` describes the dataplane that was the target of the configuration push.",
				DataplaneKey,
			),
		},
		[]string{DataplaneKey},
	)

	metrics.Registry.Unregister(controllerMetrics.ConfigPushCount)
	metrics.Registry.Unregister(controllerMetrics.ConfigPushBrokenResources)
	metrics.Registry.Unregister(controllerMetrics.TranslationCount)
	metrics.Registry.Unregister(controllerMetrics.TranslationBrokenResources)
	metrics.Registry.Unregister(controllerMetrics.ConfigPushDuration)
	metrics.Registry.Unregister(controllerMetrics.ConfigPushSuccessTime)

	metrics.Registry.MustRegister(
		controllerMetrics.ConfigPushCount,
		controllerMetrics.ConfigPushBrokenResources,
		controllerMetrics.TranslationCount,
		controllerMetrics.TranslationBrokenResources,
		controllerMetrics.ConfigPushDuration,
		controllerMetrics.ConfigPushSuccessTime,
	)

	return controllerMetrics
}

// RecordPushSuccess records a successful configuration push.
func (c *CtrlFuncMetrics) RecordPushSuccess(p Protocol, d time.Duration, dataplane string) {
	dpOpt := withDataplane(dataplane)
	c.recordPushCount(p, dpOpt)
	c.recordPushDuration(p, d, dpOpt)
	c.recordPushSuccessTime(dpOpt)
	c.recordPushBrokenResources(0, dpOpt)
}

// RecordPushFailure records a failed configuration push.
func (c *CtrlFuncMetrics) RecordPushFailure(p Protocol, d time.Duration, dataplane string, count int, err error) {
	dpOpt := withDataplane(dataplane)
	c.recordPushCount(p, dpOpt, withError(err))
	c.recordPushDuration(p, d, dpOpt, withFailure())
	c.recordPushBrokenResources(count, dpOpt)
}

// RecordTranslationSuccess records a successful configuration translation.
func (c *CtrlFuncMetrics) RecordTranslationSuccess() {
	c.TranslationCount.With(prometheus.Labels{
		SuccessKey: SuccessTrue,
	}).Inc()
}

// RecordTranslationFailure records a failed configuration translation.
func (c *CtrlFuncMetrics) RecordTranslationFailure() {
	c.TranslationCount.With(prometheus.Labels{
		SuccessKey: SuccessFalse,
	}).Inc()
}

// RecordTranslationBrokenResources records the number of resources failing translation.
func (c *CtrlFuncMetrics) RecordTranslationBrokenResources(count int) {
	c.TranslationBrokenResources.Set(float64(count))
}

type recordOption func(prometheus.Labels) prometheus.Labels

func withError(err error) recordOption {
	return func(l prometheus.Labels) prometheus.Labels {
		l[FailureReasonKey] = pushFailureReason(err)
		l[SuccessKey] = SuccessFalse
		return l
	}
}

func withFailure() recordOption {
	return func(l prometheus.Labels) prometheus.Labels {
		l[SuccessKey] = SuccessFalse
		return l
	}
}

func withDataplane(dataplane string) recordOption {
	return func(l prometheus.Labels) prometheus.Labels {
		l[DataplaneKey] = dataplane
		return l
	}
}

func (c *CtrlFuncMetrics) recordPushCount(p Protocol, opts ...recordOption) {
	labels := prometheus.Labels{
		// although this is hardcoded to true here, the withError or withFailure opt function will flip it to false
		SuccessKey:       SuccessTrue,
		ProtocolKey:      string(p),
		FailureReasonKey: "",
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	c.ConfigPushCount.With(labels).Inc()
}

func (c *CtrlFuncMetrics) recordPushDuration(p Protocol, d time.Duration, opts ...recordOption) {
	labels := prometheus.Labels{
		// although this is hardcoded to true here, the withError or withFailure opt function will flip it to false
		SuccessKey:  SuccessTrue,
		ProtocolKey: string(p),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	c.ConfigPushDuration.With(labels).Observe(float64(d.Milliseconds()))
}

func (c *CtrlFuncMetrics) recordPushBrokenResources(count int, opts ...recordOption) {
	labels := prometheus.Labels{}

	for _, opt := range opts {
		labels = opt(labels)
	}

	c.ConfigPushBrokenResources.With(labels).Set(float64(count))
}

func (c *CtrlFuncMetrics) recordPushSuccessTime(opts ...recordOption) {
	labels := prometheus.Labels{}

	for _, opt := range opts {
		labels = opt(labels)
	}

	c.ConfigPushSuccessTime.With(labels).SetToCurrentTime()
}

// pushFailureReason extracts config push failure reason from an error returned
// from sendconfig's onUpdateInMemoryMode or onUpdateDBMode.
func pushFailureReason(err error) string {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return FailureReasonNetwork
	}

	if deckerrors.IsConflictErr(err) {
		return FailureReasonConflict
	}

	return FailureReasonOther
}
