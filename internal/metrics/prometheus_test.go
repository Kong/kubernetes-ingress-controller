package metrics

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	deckutils "github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"
	prom "github.com/prometheus/client_model/go"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
)

func TestNewGlobalCtrlRuntimeMetricsRecorder_DoesNotPanicWhenCalledTwice(t *testing.T) {
	require.NotPanics(t, func() {
		_ = NewGlobalCtrlRuntimeMetricsRecorder(uuid.New())
	})
	require.NotPanics(t, func() {
		_ = NewGlobalCtrlRuntimeMetricsRecorder(uuid.New())
	})
}

func TestGlobalCtrlRuntimeMetricsRecorder_AnyInstanceWritesToTheSameRegistry(t *testing.T) {
	m1 := NewGlobalCtrlRuntimeMetricsRecorder(uuid.New())
	m2 := NewGlobalCtrlRuntimeMetricsRecorder(uuid.New())

	const (
		firstDPHost  = "https://1.host"
		secondDPHost = "https://2.host"
	)
	m1.RecordPushSuccess(ProtocolDBLess, time.Millisecond, mo.Some(22), firstDPHost)
	m2.RecordPushSuccess(ProtocolDBLess, time.Millisecond, mo.Some(22), secondDPHost)

	// Verify that both instances write to the same registry (the controller-runtime's global registry).
	ctrlRuntimeGlobalRegistry := metrics.Registry
	metricFamilies, err := ctrlRuntimeGlobalRegistry.Gather()
	require.NoError(t, err)
	metricFamily, ok := lo.Find(metricFamilies, func(family *prom.MetricFamily) bool {
		return family.Name != nil && *family.Name == MetricNameConfigPushCount
	})
	require.True(t, ok, "expected to find %q metric", MetricNameConfigPushCount)

	assertMetricsContainHost := func(host string) {
		assert.True(t, lo.ContainsBy(metricFamily.Metric, func(metric *prom.Metric) bool {
			return lo.ContainsBy(metric.Label, func(label *prom.LabelPair) bool {
				return label.GetName() == DataplaneKey && label.GetValue() == host
			})
		}), "expected to find %q metric with dataplane label %q", MetricNameConfigPushCount, host)
	}
	assertMetricsContainHost(firstDPHost)
	assertMetricsContainHost(secondDPHost)
}

func TestRecordPush(t *testing.T) {
	mockSizeOfCfg := mo.Some(22)
	m := NewGlobalCtrlRuntimeMetricsRecorder(uuid.New())

	t.Run("recording push success works", func(t *testing.T) {
		require.NotPanics(t, func() {
			m.RecordPushSuccess(ProtocolDBLess, time.Millisecond, mockSizeOfCfg, "https://10.0.0.1:8080")
		})
	})
	t.Run("recording push failure works", func(t *testing.T) {
		require.NotPanics(t, func() {
			m.RecordPushFailure(ProtocolDBLess, time.Millisecond, mockSizeOfCfg, "https://10.0.0.1:8080", 5, fmt.Errorf("custom error"))
		})
	})
	// Verify that multiple call of NewGlobalCtrlRuntimeMetricsRecorder keeps all created metrics work.
	m2 := NewGlobalCtrlRuntimeMetricsRecorder(uuid.New())
	t.Run("recording push success works for old metrics", func(t *testing.T) {
		require.NotPanics(t, func() {
			m.RecordPushSuccess(ProtocolDBLess, time.Millisecond, mockSizeOfCfg, "https://10.0.0.1:8080")
		})
	})
	t.Run("recording push success works for new metrics", func(t *testing.T) {
		require.NotPanics(t, func() {
			m2.RecordPushSuccess(ProtocolDBLess, time.Millisecond, mockSizeOfCfg, "https://10.0.0.2:8080")
		})
	})
}

func TestRecordTranslation(t *testing.T) {
	m := NewGlobalCtrlRuntimeMetricsRecorder(uuid.New())
	t.Run("recording translation success works", func(t *testing.T) {
		require.NotPanics(t, func() {
			m.RecordTranslationSuccess(10 * time.Millisecond)
			m.RecordTranslationBrokenResources(0)
		})
	})
	t.Run("recording translation failure works", func(t *testing.T) {
		require.NotPanics(t, func() {
			m.RecordTranslationFailure(10 * time.Millisecond)
			m.RecordTranslationBrokenResources(9)
		})
	})
}

func TestPushFailureReason(t *testing.T) {
	apiConflictErr := kong.NewAPIError(http.StatusConflict, "conflict api error")
	networkErr := net.UnknownNetworkError("network error")
	genericError := errors.New("generic error")

	testCases := []struct {
		name           string
		err            error
		expectedReason string
	}{
		{
			name:           "generic_error",
			err:            genericError,
			expectedReason: FailureReasonOther,
		},
		{
			name:           "api_conflict_error",
			err:            apiConflictErr,
			expectedReason: FailureReasonConflict,
		},
		{
			name:           "api_conflict_error_wrapped",
			err:            fmt.Errorf("wrapped conflict api err: %w", apiConflictErr),
			expectedReason: FailureReasonConflict,
		},
		{
			name:           "deck_config_conflict_error_empty",
			err:            deckerrors.ConfigConflictError{},
			expectedReason: FailureReasonConflict,
		},
		{
			name:           "deck_config_conflict_error_with_generic_error",
			err:            deckerrors.ConfigConflictError{Err: genericError},
			expectedReason: FailureReasonConflict,
		},
		{
			name:           "deck_err_array_with_api_conflict_error",
			err:            deckutils.ErrArray{Errors: []error{apiConflictErr}},
			expectedReason: FailureReasonConflict,
		},
		{
			name:           "wrapped_deck_err_array_with_api_conflict_error",
			err:            fmt.Errorf("wrapped: %w", deckutils.ErrArray{Errors: []error{apiConflictErr}}),
			expectedReason: FailureReasonConflict,
		},
		{
			name:           "deck_err_array_with_generic_error",
			err:            deckutils.ErrArray{Errors: []error{genericError}},
			expectedReason: FailureReasonOther,
		},
		{
			name:           "deck_err_array_empty",
			err:            deckutils.ErrArray{Errors: []error{genericError}},
			expectedReason: FailureReasonOther,
		},
		{
			name:           "network_error",
			err:            networkErr,
			expectedReason: FailureReasonNetwork,
		},
		{
			name:           "network_error_wrapped_in_deck_config_conflict_error",
			err:            deckerrors.ConfigConflictError{Err: networkErr},
			expectedReason: FailureReasonNetwork,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			reason := pushFailureReason(tc.err)
			require.Equal(t, tc.expectedReason, reason)
		})
	}
}
