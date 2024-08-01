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
	// Regular config push metrics.
	ConfigPushCount            *prometheus.CounterVec
	ConfigPushBrokenResources  *prometheus.GaugeVec
	TranslationCount           *prometheus.CounterVec
	TranslationBrokenResources prometheus.Gauge
	TranslationDuration        *prometheus.HistogramVec
	ConfigPushDuration         *prometheus.HistogramVec
	ConfigPushSuccessTime      *prometheus.GaugeVec

	// Fallback config push metrics.
	FallbackTranslationCount           *prometheus.CounterVec
	FallbackTranslationBrokenResources prometheus.Gauge
	FallbackConfigPushCount            *prometheus.CounterVec
	FallbackConfigPushSuccessTime      *prometheus.GaugeVec
	FallbackConfigPushBrokenResources  *prometheus.GaugeVec
	FallbackConfigPushDuration         *prometheus.HistogramVec
	FallbackCacheGeneratingDuration    *prometheus.HistogramVec
	ProcessedConfigSnapshotCacheHit    prometheus.Counter
	ProcessedConfigSnapshotCacheMiss   prometheus.Counter
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

// Regular config push metrics names.
const (
	MetricNameConfigPushCount            = "ingress_controller_configuration_push_count"
	MetricNameConfigPushBrokenResources  = "ingress_controller_configuration_push_broken_resource_count"
	MetricNameConfigPushSuccessTime      = "ingress_controller_configuration_push_last_successful"
	MetricNameTranslationCount           = "ingress_controller_translation_count"
	MetricNameTranslationBrokenResources = "ingress_controller_translation_broken_resource_count"
	MetricNameTranslationDuration        = "ingress_controller_translation_duration_milliseconds"
	MetricNameConfigPushDuration         = "ingress_controller_configuration_push_duration_milliseconds"
)

// Fallback config push metrics names.
const (
	MetricNameFallbackTranslationCount           = "ingress_controller_fallback_translation_count"
	MetricNameFallbackTranslationBrokenResources = "ingress_controller_fallback_translation_broken_resource_count"
	MetricNameFallbackConfigPushCount            = "ingress_controller_fallback_configuration_push_count"
	MetricNameFallbackConfigPushSuccessTime      = "ingress_controller_fallback_configuration_push_last"
	MetricNameFallbackConfigPushDuration         = "ingress_controller_fallback_configuration_push_duration_milliseconds"
	MetricNameFallbackConfigPushBrokenResources  = "ingress_controller_fallback_configuration_push_broken_resource_count"
	MetricNameFallbackCacheGenerationDuration    = "ingress_controller_fallback_cache_generation_duration_milliseconds"
	MetricNameProcessedConfigSnapshotCacheHit    = "ingress_controller_processed_config_snapshot_cache_hit"
	MetricNameProcessedConfigSnapshotCacheMiss   = "ingress_controller_processed_config_snapshot_cache_miss"
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

	// TODO: add new metric for fallback config generation or add a "fallback" label?
	controllerMetrics.TranslationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: MetricNameTranslationDuration,
			Help: fmt.Sprintf(
				"Duration of translations from Kubernetes state to Kong state."+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`). "+
					"Unrecoverable error in this case means KIC wasn't able to translate a Kubernetes object to Kong model.",
				SuccessKey, SuccessFalse, SuccessTrue,
			),
			// TODO: A simple configuration translation with 1 ingress and 1 service costs 100~200 microseconds.
			// Should we use a smaller starting value of the buckets?
			Buckets: prometheus.ExponentialBuckets(1, 2, 20),
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

	controllerMetrics.FallbackTranslationCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricNameFallbackTranslationCount,
			Help: fmt.Sprintf("Count of translations from Kubernetes state to Kong state in fallback mode. "+
				"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`). "+
				"Unrecoverable error in this case means KIC wasn't able to translate a Kubernetes object to Kong model.",
				SuccessKey, SuccessFalse, SuccessTrue,
			),
		},
		[]string{SuccessKey},
	)

	controllerMetrics.FallbackTranslationBrokenResources = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: MetricNameFallbackTranslationBrokenResources,
			Help: fmt.Sprintf("The number of resources that the controller cannot successfully translate to Kong " +
				"configuration in fallback mode.",
			),
		},
	)

	controllerMetrics.FallbackConfigPushCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricNameFallbackConfigPushCount,
			Help: fmt.Sprintf(
				"Count of successful/failed fallback configuration pushes to Kong. "+
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

	controllerMetrics.FallbackConfigPushSuccessTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameFallbackConfigPushSuccessTime,
			Help: fmt.Sprintf("The time of the last successful fallback configuration push. "+
				"`%s` describes the dataplane that was the target of the configuration push.",
				DataplaneKey,
			),
		},
		[]string{DataplaneKey},
	)

	controllerMetrics.FallbackConfigPushDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: MetricNameFallbackConfigPushDuration,
			Help: fmt.Sprintf(
				"How long it took to push the fallback configuration to Kong, in milliseconds. "+
					"`%s` describes the dataplane that was the target of configuration push. "+
					"`%s` describes the configuration protocol (`%s` or `%s`) in use. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`).",
				DataplaneKey,
				ProtocolKey, ProtocolDBLess, ProtocolDeck,
				SuccessKey, SuccessFalse, SuccessTrue,
			),
		},
		[]string{SuccessKey, ProtocolKey, DataplaneKey},
	)

	controllerMetrics.FallbackConfigPushBrokenResources = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameFallbackConfigPushBrokenResources,
			Help: fmt.Sprintf("The number of resources not accepted by Kong when attempting to push "+
				"fallback configuration. `%s` describes the dataplane that was the target of the configuration push.",
				DataplaneKey,
			),
		},
		[]string{DataplaneKey},
	)

	controllerMetrics.FallbackCacheGeneratingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: MetricNameFallbackCacheGenerationDuration,
			Help: fmt.Sprintf("How long it took to generate a fallback cache, in milliseconds. "+
				"`%s` describes whether the cache generation was successful (`%s`) or not (`%s`).",
				SuccessKey, SuccessTrue, SuccessFalse,
			),
		},
		[]string{SuccessKey},
	)

	controllerMetrics.ProcessedConfigSnapshotCacheHit = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: MetricNameProcessedConfigSnapshotCacheHit,
			Help: "The number of times the controller hit the processed config snapshot cache and skipped generating " +
				"a new one.",
		},
	)

	controllerMetrics.ProcessedConfigSnapshotCacheMiss = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: MetricNameProcessedConfigSnapshotCacheMiss,
			Help: "The number of times the controller missed the processed config snapshot cache and had to generate " +
				"a new one.",
		},
	)

	allMetrics := []prometheus.Collector{
		controllerMetrics.ConfigPushCount,
		controllerMetrics.ConfigPushBrokenResources,
		controllerMetrics.TranslationCount,
		controllerMetrics.TranslationDuration,
		controllerMetrics.TranslationBrokenResources,
		controllerMetrics.ConfigPushDuration,
		controllerMetrics.ConfigPushSuccessTime,
		controllerMetrics.FallbackTranslationBrokenResources,
		controllerMetrics.FallbackTranslationCount,
		controllerMetrics.FallbackConfigPushCount,
		controllerMetrics.FallbackConfigPushSuccessTime,
		controllerMetrics.FallbackConfigPushDuration,
		controllerMetrics.FallbackConfigPushBrokenResources,
		controllerMetrics.FallbackCacheGeneratingDuration,
		controllerMetrics.ProcessedConfigSnapshotCacheHit,
		controllerMetrics.ProcessedConfigSnapshotCacheMiss,
	}
	for _, m := range allMetrics {
		metrics.Registry.Unregister(m)
		metrics.Registry.MustRegister(m)
	}

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
func (c *CtrlFuncMetrics) RecordTranslationSuccess(duration time.Duration) {
	c.TranslationCount.With(prometheus.Labels{
		SuccessKey: SuccessTrue,
	}).Inc()
	c.TranslationDuration.With(prometheus.Labels{
		SuccessKey: SuccessTrue,
	}).Observe(float64(duration.Milliseconds()))
}

// RecordTranslationFailure records a failed configuration translation.
func (c *CtrlFuncMetrics) RecordTranslationFailure(duration time.Duration) {
	c.TranslationCount.With(prometheus.Labels{
		SuccessKey: SuccessFalse,
	}).Inc()
	c.TranslationDuration.With(prometheus.Labels{
		SuccessKey: SuccessFalse,
	}).Observe(float64(duration.Milliseconds()))
}

// RecordTranslationBrokenResources records the number of resources failing translation.
func (c *CtrlFuncMetrics) RecordTranslationBrokenResources(count int) {
	c.TranslationBrokenResources.Set(float64(count))
}

// RecordFallbackTranslationFailure records a failed fallback configuration translation.
func (c *CtrlFuncMetrics) RecordFallbackTranslationFailure() {
	c.FallbackTranslationCount.With(prometheus.Labels{
		SuccessKey: SuccessFalse,
	}).Inc()
}

// RecordFallbackTranslationSuccess records a failed fallback configuration translation.
func (c *CtrlFuncMetrics) RecordFallbackTranslationSuccess() {
	c.FallbackTranslationCount.With(prometheus.Labels{
		SuccessKey: SuccessTrue,
	}).Inc()
}

// RecordProcessedConfigSnapshotCacheHit records a hit on the processed config snapshot cache.
func (c *CtrlFuncMetrics) RecordProcessedConfigSnapshotCacheHit() {
	c.ProcessedConfigSnapshotCacheHit.Inc()
}

// RecordProcessedConfigSnapshotCacheMiss records a miss on the processed config snapshot cache.
func (c *CtrlFuncMetrics) RecordProcessedConfigSnapshotCacheMiss() {
	c.ProcessedConfigSnapshotCacheMiss.Inc()
}

// RecordFallbackTranslationBrokenResources records the number of fallback resources failing translation.
func (c *CtrlFuncMetrics) RecordFallbackTranslationBrokenResources(count int) {
	c.FallbackTranslationBrokenResources.Set(float64(count))
}

// RecordFallbackPushSuccess records a successful fallback configuration push.
func (c *CtrlFuncMetrics) RecordFallbackPushSuccess(p Protocol, duration time.Duration, dataplane string) {
	dpOpt := withDataplane(dataplane)
	c.recordFallbackPushCount(p, dpOpt)
	c.recordFallbackPushDuration(p, duration, dpOpt)
	c.recordFallbackPushSuccessTime(dpOpt)
	c.recordFallbackPushBrokenResources(0, dpOpt)
}

// RecordFallbackPushFailure records a failed fallback configuration push.
func (c *CtrlFuncMetrics) RecordFallbackPushFailure(p Protocol, duration time.Duration, dataplane string, brokenResourcesCount int, err error) {
	dpOpt := withDataplane(dataplane)
	c.recordFallbackPushDuration(p, duration, dpOpt, withFailure())
	c.recordFallbackPushCount(p, dpOpt, withError(err))
	c.recordFallbackPushBrokenResources(brokenResourcesCount, dpOpt)
}

// RecordFallbackCacheGenerationDuration records the duration of a fallback cache generation.
func (c *CtrlFuncMetrics) RecordFallbackCacheGenerationDuration(d time.Duration, err error) {
	labels := prometheus.Labels{
		SuccessKey: SuccessTrue,
	}
	if err != nil {
		labels[SuccessKey] = SuccessFalse
	}
	c.FallbackCacheGeneratingDuration.With(labels).Observe(float64(d.Milliseconds()))
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

func (c *CtrlFuncMetrics) recordFallbackPushCount(p Protocol, opts ...recordOption) {
	labels := prometheus.Labels{
		// Although this is hardcoded to true here, the withError or withFailure opt function will flip it to false.
		SuccessKey:       SuccessTrue,
		ProtocolKey:      string(p),
		FailureReasonKey: "",
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	c.FallbackConfigPushCount.With(labels).Inc()
}

func (c *CtrlFuncMetrics) recordFallbackPushSuccessTime(opts ...recordOption) {
	labels := prometheus.Labels{}

	for _, opt := range opts {
		labels = opt(labels)
	}

	c.FallbackConfigPushSuccessTime.With(labels).SetToCurrentTime()
}

func (c *CtrlFuncMetrics) recordFallbackPushDuration(p Protocol, d time.Duration, opts ...recordOption) {
	labels := prometheus.Labels{
		// Although this is hardcoded to true here, the withError or withFailure opt function will flip it to false.
		SuccessKey:  SuccessTrue,
		ProtocolKey: string(p),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	c.FallbackConfigPushDuration.With(labels).Observe(float64(d.Milliseconds()))
}

func (c *CtrlFuncMetrics) recordFallbackPushBrokenResources(brokenObjectsCount int, opts ...recordOption) {
	labels := prometheus.Labels{}

	for _, opt := range opts {
		labels = opt(labels)
	}

	c.FallbackConfigPushBrokenResources.With(labels).Set(float64(brokenObjectsCount))
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
