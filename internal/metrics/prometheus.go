package metrics

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckerrors"
)

type CtrlFuncMetrics struct {
	// ConfigPushCount is a Prometheus metric with semantics defined by its help string in NewCtrlFuncMetrics().
	ConfigPushCount *prometheus.CounterVec

	// TranslationCount is a Prometheus metric with semantics defined by its help string in NewCtrlFuncMetrics().
	TranslationCount *prometheus.CounterVec

	// ConfigPushDuration is a Prometheus metric with semantics defined by its help string in NewCtrlFuncMetrics().
	ConfigPushDuration *prometheus.HistogramVec
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
	MetricNameConfigPushCount    = "ingress_controller_configuration_push_count"
	MetricNameTranslationCount   = "ingress_controller_translation_count"
	MetricNameConfigPushDuration = "ingress_controller_configuration_push_duration_milliseconds"
)

var _once sync.Once

func NewCtrlFuncMetrics() *CtrlFuncMetrics {
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

	_once.Do(func() {
		metrics.Registry.MustRegister(controllerMetrics.ConfigPushCount, controllerMetrics.TranslationCount, controllerMetrics.ConfigPushDuration)
	})

	return controllerMetrics
}

// RecordPushSuccess records a successful configuration push.
func (c *CtrlFuncMetrics) RecordPushSuccess(p Protocol, d time.Duration, dataplane string) {
	dpOpt := withDataplane(dataplane)
	c.recordPushCount(p, dpOpt)
	c.recordPushDuration(p, d, dpOpt)
}

// RecordPushFailure records a failed configuration push.
func (c *CtrlFuncMetrics) RecordPushFailure(p Protocol, d time.Duration, dataplane string, err error) {
	dpOpt := withDataplane(dataplane)
	c.recordPushCount(p, dpOpt, withError(err))
	c.recordPushDuration(p, d, dpOpt, withFailure())
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
		SuccessKey:  SuccessTrue,
		ProtocolKey: string(p),
	}

	for _, opt := range opts {
		labels = opt(labels)
	}

	c.ConfigPushDuration.With(labels).Observe(float64(d.Milliseconds()))
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
