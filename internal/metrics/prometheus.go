package metrics

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/mo"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
)

// Recorder is an interface for recording metrics.
type Recorder interface {
	RecordPushFailure(p Protocol, duration time.Duration, size mo.Option[int], dataplane string, brokenResourcesCount int, err error)
	RecordPushSuccess(protocol Protocol, duration time.Duration, size mo.Option[int], target string)
	RecordFallbackPushSuccess(protocol Protocol, duration time.Duration, size mo.Option[int], target string)
	RecordFallbackPushFailure(protocol Protocol, duration time.Duration, size mo.Option[int], target string, failedResources int, err error)
	RecordProcessedConfigSnapshotCacheHit()
	RecordProcessedConfigSnapshotCacheMiss()
	RecordTranslationFailure(duration time.Duration)
	RecordTranslationBrokenResources(count int)
	RecordTranslationSuccess(duration time.Duration)
	RecordFallbackTranslationBrokenResources(count int)
	RecordFallbackTranslationFailure(duration time.Duration)
	RecordFallbackTranslationSuccess(duration time.Duration)
	RecordFallbackCacheGenerationDuration(since time.Duration, err error)
}

var _ Recorder = &GlobalCtrlRuntimeMetricsRecorder{}

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

	// InstanceIDKey defines the name of the metric label indicating which instance of the controller this time series is relevant for.
	InstanceIDKey string = "instance_id"
)

// Regular config push metrics names.
const (
	MetricNameConfigPushCount            = "ingress_controller_configuration_push_count"
	MetricNameConfigPushBrokenResources  = "ingress_controller_configuration_push_broken_resource_count"
	MetricNameConfigPushSuccessTime      = "ingress_controller_configuration_push_last_successful"
	MetricNameConfigPushSize             = "ingress_controller_configuration_push_size"
	MetricNameTranslationCount           = "ingress_controller_translation_count"
	MetricNameTranslationBrokenResources = "ingress_controller_translation_broken_resource_count"
	MetricNameTranslationDuration        = "ingress_controller_translation_duration_milliseconds"
	MetricNameConfigPushDuration         = "ingress_controller_configuration_push_duration_milliseconds"
)

// Fallback config push metrics names.
const (
	MetricNameFallbackTranslationCount           = "ingress_controller_fallback_translation_count"
	MetricNameFallbackTranslationBrokenResources = "ingress_controller_fallback_translation_broken_resource_count"
	MetricNameFallbackTranslationDuration        = "ingress_controller_fallback_translation_duration_milliseconds"
	MetricNameFallbackConfigPushSize             = "ingress_controller_fallback_configuration_push_size"
	MetricNameFallbackConfigPushCount            = "ingress_controller_fallback_configuration_push_count"
	MetricNameFallbackConfigPushSuccessTime      = "ingress_controller_fallback_configuration_push_last"
	MetricNameFallbackConfigPushDuration         = "ingress_controller_fallback_configuration_push_duration_milliseconds"
	MetricNameFallbackConfigPushBrokenResources  = "ingress_controller_fallback_configuration_push_broken_resource_count"
	MetricNameFallbackCacheGenerationDuration    = "ingress_controller_fallback_cache_generation_duration_milliseconds"
	MetricNameProcessedConfigSnapshotCacheHit    = "ingress_controller_processed_config_snapshot_cache_hit"
	MetricNameProcessedConfigSnapshotCacheMiss   = "ingress_controller_processed_config_snapshot_cache_miss"
)

// Metrics definitions for GlobalCtrlRuntimeMetricsRecorder.
var (
	configPushCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricNameConfigPushCount,
			Help: fmt.Sprintf(
				"Count of successful/failed configuration pushes to Kong. "+
					"`%s` describes the dataplane that was the target of configuration push. "+
					"`%s` describes the configuration protocol (`%s` or `%s`) in use. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`). "+
					"`%s` is populated in case of `%s=\"%s\"` and describes the reason of failure "+
					"(one of `%s`, `%s`, `%s`)."+
					"`%s` describes the instance of the controller that pushed the configuration.",
				DataplaneKey,
				ProtocolKey, ProtocolDBLess, ProtocolDeck,
				SuccessKey, SuccessFalse, SuccessTrue,
				FailureReasonKey, SuccessKey, SuccessFalse,
				FailureReasonConflict, FailureReasonNetwork, FailureReasonOther,
				InstanceIDKey,
			),
		},
		[]string{SuccessKey, ProtocolKey, FailureReasonKey, DataplaneKey, InstanceIDKey},
	)

	configPushBrokenResources = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameConfigPushBrokenResources,
			Help: fmt.Sprintf("The number of resources not accepted by Kong when attempting to push "+
				"configuration. `%s` describes the dataplane that was the target of the configuration push."+
				"`%s` describes the instance of the controller that pushed the configuration.",
				DataplaneKey,
				InstanceIDKey,
			),
		},
		[]string{DataplaneKey, InstanceIDKey},
	)

	translationCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricNameTranslationCount,
			Help: fmt.Sprintf(
				"Count of translations from Kubernetes state to Kong state. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`). "+
					"Unrecoveirable error in this case means KIC wasn't able to translate a Kubernetes object to Kong model."+
					"`%s` describes the instance of the controller that pushed the configuration.",
				SuccessKey, SuccessFalse, SuccessTrue,
				InstanceIDKey,
			),
		},
		[]string{SuccessKey, InstanceIDKey},
	)

	translationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: MetricNameTranslationDuration,
			Help: fmt.Sprintf(
				"Duration of  translating Kubernetes resources into Kong state. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`). "+
					"Unrecoverable error in this case means KIC wasn't able to translate a Kubernetes object to Kong model."+
					"`%s` describes the instance of the controller that pushed the configuration.",
				SuccessKey, SuccessFalse, SuccessTrue,
				InstanceIDKey,
			),
			Buckets: prometheus.ExponentialBucketsRange(1, float64(time.Minute.Milliseconds()), 30),
		},
		[]string{SuccessKey, InstanceIDKey},
	)

	translationBrokenResources = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameTranslationBrokenResources,
			Help: fmt.Sprintf("The number of resources that the controller cannot successfully translate to Kong "+
				"configuration."+
				"`%s` describes the instance of the controller that pushed the configuration.",
				InstanceIDKey,
			),
		},
		[]string{InstanceIDKey},
	)

	configPushDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: MetricNameConfigPushDuration,
			Help: fmt.Sprintf(
				"How long it took to push the configuration to Kong, in milliseconds. "+
					"`%s` describes the dataplane that was the target of configuration push. "+
					"`%s` describes the configuration protocol (`%s` or `%s`) in use. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`)."+
					"`%s` describes the instance of the controller that pushed the configuration.",
				DataplaneKey,
				ProtocolKey, ProtocolDBLess, ProtocolDeck,
				SuccessKey, SuccessFalse, SuccessTrue,
				InstanceIDKey,
			),
			Buckets: prometheus.ExponentialBuckets(100, 1.33, 30),
		},
		[]string{SuccessKey, ProtocolKey, DataplaneKey, InstanceIDKey},
	)

	configPushSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameConfigPushSize,
			Help: fmt.Sprintf(
				"The size of the configuration pushed to Kong, in bytes. "+
					"`%s` describes the dataplane that was the target of the configuration push. "+
					"`%s` describes the configuration protocol (metric is presented for `%s`, for `%s` it doesn't exist) in use. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`)."+
					"`%s` describes the instance of the controller that pushed the configuration.",
				DataplaneKey,
				ProtocolKey, ProtocolDBLess, ProtocolDeck,
				SuccessKey, SuccessFalse, SuccessTrue,
				InstanceIDKey,
			),
		},
		[]string{DataplaneKey, ProtocolKey, SuccessKey, InstanceIDKey},
	)

	configPushSuccessTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameConfigPushSuccessTime,
			Help: fmt.Sprintf("The time of the last successful configuration push. "+
				"`%s` describes the dataplane that was the target of the configuration push."+
				"`%s` describes the instance of the controller that pushed the configuration.",
				DataplaneKey,
				InstanceIDKey,
			),
		},
		[]string{DataplaneKey, InstanceIDKey},
	)

	fallbackTranslationCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricNameFallbackTranslationCount,
			Help: fmt.Sprintf("Count of translations from Kubernetes state to Kong state in fallback mode. "+
				"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`). "+
				"Unrecoverable error in this case means KIC wasn't able to translate a Kubernetes object to Kong model."+
				"`%s` describes the instance of the controller that pushed the configuration.",
				SuccessKey, SuccessFalse, SuccessTrue,
				InstanceIDKey,
			),
		},
		[]string{SuccessKey, InstanceIDKey},
	)

	fallbackTranslationBrokenResources = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameFallbackTranslationBrokenResources,
			Help: fmt.Sprintf("The number of resources that the controller cannot successfully translate to Kong "+
				"configuration in fallback mode."+
				"`%s` describes the instance of the controller that pushed the configuration.",
				InstanceIDKey,
			),
		},
		[]string{InstanceIDKey},
	)

	fallbackConfigPushSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameFallbackConfigPushSize,
			Help: fmt.Sprintf(
				"The size of the configuration pushed to Kong in fallback mode, in bytes. "+
					"`%s` describes the dataplane that was the target of the configuration push. "+
					"`%s` describes the configuration protocol (metric is presented for `%s`, for `%s` it doesn't exist) in use. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`)."+
					"`%s` describes the instance of the controller that pushed the configuration.",
				DataplaneKey,
				ProtocolKey, ProtocolDBLess, ProtocolDeck,
				SuccessKey, SuccessFalse, SuccessTrue,
				InstanceIDKey,
			),
		},
		[]string{DataplaneKey, ProtocolKey, SuccessKey, InstanceIDKey},
	)

	fallbackTranslationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: MetricNameFallbackTranslationDuration,
			Help: fmt.Sprintf(
				"Duration of translating Kubernetes resources into Kong state in fallback mode. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`). "+
					"Unrecoverable error in this case means KIC wasn't able to translate a Kubernetes object to Kong model."+
					"`%s` describes the instance of the controller that pushed the configuration.",
				SuccessKey, SuccessFalse, SuccessTrue,
				InstanceIDKey,
			),
			Buckets: prometheus.ExponentialBucketsRange(1, float64(time.Minute.Milliseconds()), 30),
		},
		[]string{SuccessKey, InstanceIDKey},
	)

	fallbackConfigPushCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricNameFallbackConfigPushCount,
			Help: fmt.Sprintf(
				"Count of successful/failed fallback configuration pushes to Kong. "+
					"`%s` describes the dataplane that was the target of configuration push. "+
					"`%s` describes the configuration protocol (`%s` or `%s`) in use. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`). "+
					"`%s` is populated in case of `%s=\"%s\"` and describes the reason of failure "+
					"(one of `%s`, `%s`, `%s`)."+
					"`%s` describes the instance of the controller that pushed the configuration.",
				DataplaneKey,
				ProtocolKey, ProtocolDBLess, ProtocolDeck,
				SuccessKey, SuccessFalse, SuccessTrue,
				FailureReasonKey, SuccessKey, SuccessFalse,
				FailureReasonConflict, FailureReasonNetwork, FailureReasonOther,
				InstanceIDKey,
			),
		},
		[]string{SuccessKey, ProtocolKey, FailureReasonKey, DataplaneKey, InstanceIDKey},
	)

	fallbackConfigPushSuccessTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameFallbackConfigPushSuccessTime,
			Help: fmt.Sprintf("The time of the last successful fallback configuration push. "+
				"`%s` describes the dataplane that was the target of the configuration push."+
				"`%s` describes the instance of the controller that pushed the configuration.",
				DataplaneKey,
				InstanceIDKey,
			),
		},
		[]string{DataplaneKey, InstanceIDKey},
	)

	fallbackConfigPushDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: MetricNameFallbackConfigPushDuration,
			Help: fmt.Sprintf(
				"How long it took to push the fallback configuration to Kong, in milliseconds. "+
					"`%s` describes the dataplane that was the target of configuration push. "+
					"`%s` describes the configuration protocol (`%s` or `%s`) in use. "+
					"`%s` describes whether there were unrecoverable errors (`%s`) or not (`%s`)."+
					"`%s` describes the instance of the controller that pushed the configuration.",
				DataplaneKey,
				ProtocolKey, ProtocolDBLess, ProtocolDeck,
				SuccessKey, SuccessFalse, SuccessTrue,
				InstanceIDKey,
			),
		},
		[]string{SuccessKey, ProtocolKey, DataplaneKey, InstanceIDKey},
	)

	fallbackConfigPushBrokenResources = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: MetricNameFallbackConfigPushBrokenResources,
			Help: fmt.Sprintf("The number of resources not accepted by Kong when attempting to push "+
				"fallback configuration. `%s` describes the dataplane that was the target of the configuration push."+
				"`%s` describes the instance of the controller that pushed the configuration.",
				DataplaneKey,
				InstanceIDKey,
			),
		},
		[]string{DataplaneKey, InstanceIDKey},
	)

	fallbackCacheGeneratingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: MetricNameFallbackCacheGenerationDuration,
			Help: fmt.Sprintf("How long it took to generate a fallback cache, in milliseconds. "+
				"`%s` describes whether the cache generation was successful (`%s`) or not (`%s`)."+
				"`%s` describes the instance of the controller that pushed the configuration.",
				SuccessKey, SuccessTrue, SuccessFalse,
				InstanceIDKey,
			),
		},
		[]string{SuccessKey, InstanceIDKey},
	)

	processedConfigSnapshotCacheHit = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricNameProcessedConfigSnapshotCacheHit,
			Help: fmt.Sprintf("The number of times the controller hit the processed config snapshot cache and skipped generating "+
				"a new one."+
				"`%s` describes the instance of the controller that pushed the configuration.",
				InstanceIDKey,
			),
		},
		[]string{InstanceIDKey},
	)

	processedConfigSnapshotCacheMiss = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricNameProcessedConfigSnapshotCacheMiss,
			Help: fmt.Sprintf("The number of times the controller missed the processed config snapshot cache and had to generate "+
				"a new one."+
				"`%s` describes the instance of the controller that pushed the configuration.",
				InstanceIDKey,
			),
		},
		[]string{InstanceIDKey},
	)
)

func init() {
	allMetrics := []prometheus.Collector{
		configPushCount,
		configPushBrokenResources,
		translationCount,
		translationDuration,
		translationBrokenResources,
		configPushDuration,
		configPushSize,
		configPushSuccessTime,
		fallbackTranslationCount,
		fallbackTranslationBrokenResources,
		fallbackTranslationDuration,
		fallbackConfigPushCount,
		fallbackConfigPushSuccessTime,
		fallbackConfigPushDuration,
		fallbackConfigPushSize,
		fallbackConfigPushBrokenResources,
		fallbackCacheGeneratingDuration,
		processedConfigSnapshotCacheHit,
		processedConfigSnapshotCacheMiss,
	}
	for _, m := range allMetrics {
		metrics.Registry.MustRegister(m)
	}
}

// InstanceID is an interface for a controller manager instance identifier.
type InstanceID interface {
	String() string
}

// GlobalCtrlRuntimeMetricsRecorder is a metrics recorder that uses a global Prometheus registry
// provided by the controller-runtime. Any instance of it will record metrics to the same registry.
//
// We want to expose KIC's custom metrics on the same endpoint as controller-runtime's built-in
// ones. Because of that, we have to use its global registry as CR doesn't allow injecting a custom one.
// Upstream issue regarding this: https://github.com/kubernetes-sigs/controller-runtime/issues/210.
type GlobalCtrlRuntimeMetricsRecorder struct {
	instanceID InstanceID
}

func NewGlobalCtrlRuntimeMetricsRecorder(instanceID InstanceID) *GlobalCtrlRuntimeMetricsRecorder {
	return &GlobalCtrlRuntimeMetricsRecorder{instanceID: instanceID}
}

// RecordPushSuccess records a successful configuration push.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordPushSuccess(p Protocol, d time.Duration, size mo.Option[int], dataplane string) {
	dpOpt := withDataplane(dataplane)
	c.recordPushCount(p, dpOpt)
	c.recordPushDuration(p, d, dpOpt)
	c.recordPushSuccessTime(dpOpt)
	c.recordPushBrokenResources(0, dpOpt)
	c.recordConfigPushSize(p, size, dpOpt)
}

// RecordPushFailure records a failed configuration push.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordPushFailure(p Protocol, d time.Duration, size mo.Option[int], dataplane string, count int, err error) {
	dpOpt := withDataplane(dataplane)
	c.recordPushCount(p, dpOpt, withError(err))
	c.recordPushDuration(p, d, dpOpt, withFailure())
	c.recordPushBrokenResources(count, dpOpt)
	c.recordConfigPushSize(p, size, dpOpt, withFailure())
}

// RecordTranslationSuccess records a successful configuration translation.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordTranslationSuccess(duration time.Duration) {
	translationCount.With(prometheus.Labels{
		SuccessKey:    SuccessTrue,
		InstanceIDKey: c.instanceID.String(),
	}).Inc()
	translationDuration.With(prometheus.Labels{
		SuccessKey:    SuccessTrue,
		InstanceIDKey: c.instanceID.String(),
	}).Observe(float64(duration.Milliseconds()))
}

// RecordTranslationFailure records a failed configuration translation.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordTranslationFailure(duration time.Duration) {
	translationCount.With(prometheus.Labels{
		SuccessKey:    SuccessFalse,
		InstanceIDKey: c.instanceID.String(),
	}).Inc()
	translationDuration.With(prometheus.Labels{
		SuccessKey:    SuccessFalse,
		InstanceIDKey: c.instanceID.String(),
	}).Observe(float64(duration.Milliseconds()))
}

// RecordTranslationBrokenResources records the number of resources failing translation.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordTranslationBrokenResources(count int) {
	translationBrokenResources.With(prometheus.Labels{
		InstanceIDKey: c.instanceID.String(),
	}).Set(float64(count))
}

// RecordFallbackTranslationFailure records a failed fallback configuration translation.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordFallbackTranslationFailure(duration time.Duration) {
	fallbackTranslationCount.With(prometheus.Labels{
		SuccessKey:    SuccessFalse,
		InstanceIDKey: c.instanceID.String(),
	}).Inc()
	fallbackTranslationDuration.With(prometheus.Labels{
		SuccessKey:    SuccessFalse,
		InstanceIDKey: c.instanceID.String(),
	}).Observe(float64(duration))
}

// RecordFallbackTranslationSuccess records a failed fallback configuration translation.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordFallbackTranslationSuccess(duration time.Duration) {
	fallbackTranslationCount.With(prometheus.Labels{
		SuccessKey:    SuccessTrue,
		InstanceIDKey: c.instanceID.String(),
	}).Inc()
	fallbackTranslationDuration.With(prometheus.Labels{
		SuccessKey:    SuccessTrue,
		InstanceIDKey: c.instanceID.String(),
	}).Observe(float64(duration))
}

// RecordProcessedConfigSnapshotCacheHit records a hit on the processed config snapshot cache.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordProcessedConfigSnapshotCacheHit() {
	processedConfigSnapshotCacheHit.With(prometheus.Labels{
		InstanceIDKey: c.instanceID.String(),
	}).Inc()
}

// RecordProcessedConfigSnapshotCacheMiss records a miss on the processed config snapshot cache.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordProcessedConfigSnapshotCacheMiss() {
	processedConfigSnapshotCacheMiss.With(prometheus.Labels{
		InstanceIDKey: c.instanceID.String(),
	}).Inc()
}

// RecordFallbackTranslationBrokenResources records the number of fallback resources failing translation.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordFallbackTranslationBrokenResources(count int) {
	fallbackTranslationBrokenResources.With(prometheus.Labels{
		InstanceIDKey: c.instanceID.String(),
	}).Set(float64(count))
}

// RecordFallbackPushSuccess records a successful fallback configuration push.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordFallbackPushSuccess(
	p Protocol, duration time.Duration, size mo.Option[int], dataplane string,
) {
	dpOpt := withDataplane(dataplane)
	c.recordFallbackPushCount(p, dpOpt)
	c.recordFallbackPushDuration(p, duration, dpOpt)
	c.recordFallbackPushSuccessTime(dpOpt)
	c.recordFallbackPushBrokenResources(0, dpOpt)
	c.recordFallbackConfigPushSize(p, size, dpOpt)
}

// RecordFallbackPushFailure records a failed fallback configuration push.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordFallbackPushFailure(
	p Protocol, duration time.Duration, size mo.Option[int], dataplane string, brokenResourcesCount int, err error,
) {
	dpOpt := withDataplane(dataplane)
	c.recordFallbackPushDuration(p, duration, dpOpt, withFailure())
	c.recordFallbackPushCount(p, dpOpt, withError(err))
	c.recordFallbackPushBrokenResources(brokenResourcesCount, dpOpt)
	c.recordFallbackConfigPushSize(p, size, dpOpt, withFailure())
}

// RecordFallbackCacheGenerationDuration records the duration of a fallback cache generation.
func (c *GlobalCtrlRuntimeMetricsRecorder) RecordFallbackCacheGenerationDuration(d time.Duration, err error) {
	labels := prometheus.Labels{
		SuccessKey:    SuccessTrue,
		InstanceIDKey: c.instanceID.String(),
	}
	if err != nil {
		labels[SuccessKey] = SuccessFalse
	}
	fallbackCacheGeneratingDuration.With(labels).Observe(float64(d.Milliseconds()))
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

func (c *GlobalCtrlRuntimeMetricsRecorder) recordPushCount(p Protocol, opts ...recordOption) {
	labels := prometheus.Labels{
		// although this is hardcoded to true here, the withError or withFailure opt function will flip it to false
		SuccessKey:       SuccessTrue,
		ProtocolKey:      string(p),
		FailureReasonKey: "",
		InstanceIDKey:    c.instanceID.String(),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	configPushCount.With(labels).Inc()
}

func (c *GlobalCtrlRuntimeMetricsRecorder) recordPushDuration(p Protocol, d time.Duration, opts ...recordOption) {
	labels := prometheus.Labels{
		// although this is hardcoded to true here, the withError or withFailure opt function will flip it to false
		SuccessKey:    SuccessTrue,
		ProtocolKey:   string(p),
		InstanceIDKey: c.instanceID.String(),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	configPushDuration.With(labels).Observe(float64(d.Milliseconds()))
}

func (c *GlobalCtrlRuntimeMetricsRecorder) recordPushBrokenResources(count int, opts ...recordOption) {
	labels := prometheus.Labels{
		InstanceIDKey: c.instanceID.String(),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	configPushBrokenResources.With(labels).Set(float64(count))
}

func (c *GlobalCtrlRuntimeMetricsRecorder) recordPushSuccessTime(opts ...recordOption) {
	labels := prometheus.Labels{
		InstanceIDKey: c.instanceID.String(),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	configPushSuccessTime.With(labels).SetToCurrentTime()
}

func (c *GlobalCtrlRuntimeMetricsRecorder) recordFallbackPushCount(p Protocol, opts ...recordOption) {
	labels := prometheus.Labels{
		// Although this is hardcoded to true here, the withError or withFailure opt function will flip it to false.
		SuccessKey:       SuccessTrue,
		ProtocolKey:      string(p),
		FailureReasonKey: "",
		InstanceIDKey:    c.instanceID.String(),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	fallbackConfigPushCount.With(labels).Inc()
}

func (c *GlobalCtrlRuntimeMetricsRecorder) recordFallbackConfigPushSize(p Protocol, size mo.Option[int], opts ...recordOption) {
	// When size is missing do not report this metric at all.
	value, ok := size.Get()
	if !ok {
		return
	}
	// Although this is hardcoded to true here, the withError or withFailure opt function will flip it to false.
	labels := prometheus.Labels{
		SuccessKey:    SuccessTrue,
		ProtocolKey:   string(p),
		InstanceIDKey: c.instanceID.String(),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	fallbackConfigPushSize.With(labels).Set(float64(value))
}

func (c *GlobalCtrlRuntimeMetricsRecorder) recordConfigPushSize(p Protocol, size mo.Option[int], opts ...recordOption) {
	// When size is missing do not report this metric at all.
	value, ok := size.Get()
	if !ok {
		return
	}
	// Although this is hardcoded to true here, the withError or withFailure opt function will flip it to false.
	labels := prometheus.Labels{
		SuccessKey:    SuccessTrue,
		ProtocolKey:   string(p),
		InstanceIDKey: c.instanceID.String(),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	configPushSize.With(labels).Set(float64(value))
}

func (c *GlobalCtrlRuntimeMetricsRecorder) recordFallbackPushSuccessTime(opts ...recordOption) {
	labels := prometheus.Labels{
		InstanceIDKey: c.instanceID.String(),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	fallbackConfigPushSuccessTime.With(labels).SetToCurrentTime()
}

func (c *GlobalCtrlRuntimeMetricsRecorder) recordFallbackPushDuration(p Protocol, d time.Duration, opts ...recordOption) {
	labels := prometheus.Labels{
		// Although this is hardcoded to true here, the withError or withFailure opt function will flip it to false.
		SuccessKey:    SuccessTrue,
		ProtocolKey:   string(p),
		InstanceIDKey: c.instanceID.String(),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	fallbackConfigPushDuration.With(labels).Observe(float64(d.Milliseconds()))
}

func (c *GlobalCtrlRuntimeMetricsRecorder) recordFallbackPushBrokenResources(brokenObjectsCount int, opts ...recordOption) {
	labels := prometheus.Labels{
		InstanceIDKey: c.instanceID.String(),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	fallbackConfigPushBrokenResources.With(labels).Set(float64(brokenObjectsCount))
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
