package mocks

import (
	"time"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

// MetricsRecorder is a mock implementation of metrics.Recorder.
type MetricsRecorder struct{}

func (m MetricsRecorder) RecordPushFailure(metrics.Protocol, time.Duration, string, int, error) {
}

func (m MetricsRecorder) RecordPushSuccess(metrics.Protocol, time.Duration, string) {
}

func (m MetricsRecorder) RecordFallbackPushSuccess(metrics.Protocol, time.Duration, string) {
}

func (m MetricsRecorder) RecordFallbackPushFailure(metrics.Protocol, time.Duration, string, int, error) {
}

func (m MetricsRecorder) RecordProcessedConfigSnapshotCacheHit() {
}

func (m MetricsRecorder) RecordProcessedConfigSnapshotCacheMiss() {
}

func (m MetricsRecorder) RecordTranslationFailure(time.Duration) {
}

func (m MetricsRecorder) RecordTranslationBrokenResources(int) {
}

func (m MetricsRecorder) RecordTranslationSuccess(time.Duration) {
}

func (m MetricsRecorder) RecordFallbackTranslationBrokenResources(int) {
}

func (m MetricsRecorder) RecordFallbackTranslationFailure(time.Duration) {
}

func (m MetricsRecorder) RecordFallbackTranslationSuccess(time.Duration) {
}

func (m MetricsRecorder) RecordFallbackCacheGenerationDuration(time.Duration, error) {
}
