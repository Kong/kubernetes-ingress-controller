package failures

import (
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
)

const someValidResourceFailureReason = "some valid message"

func TestResourceFailure(t *testing.T) {
	t.Run("is created and returns message and causing objects", func(t *testing.T) {
		transErr, err := NewResourceFailure(someValidResourceFailureReason, someResourceFailureCausingObjects()...)
		require.NoError(t, err)

		assert.Equal(t, someValidResourceFailureReason, transErr.Message())
		assert.ElementsMatch(t, someResourceFailureCausingObjects(), transErr.CausingObjects())
	})

	t.Run("fallbacks to unknown message when empty", func(t *testing.T) {
		transErr, err := NewResourceFailure("", someResourceFailureCausingObjects()...)
		require.NoError(t, err)
		require.Equal(t, ResourceFailureReasonUnknown, transErr.Message())
	})

	t.Run("requires at least one causing object", func(t *testing.T) {
		_, err := NewResourceFailure(someValidResourceFailureReason, validCausingObject())
		require.NoError(t, err)

		_, err = NewResourceFailure(someValidResourceFailureReason)
		require.Error(t, err)
	})

	t.Run("requires valid objects", func(t *testing.T) {
		_, err := NewResourceFailure(someValidResourceFailureReason, nil)
		assert.Error(t, err, "expected a nil object to be rejected")

		emptyGVK := validCausingObject()
		emptyGVK.APIVersion = ""
		emptyGVK.Kind = ""
		_, err = NewResourceFailure(someValidResourceFailureReason, emptyGVK)
		assert.Error(t, err, "expected an empty GVK object to be rejected")

		noName := validCausingObject()
		noName.Name = ""
		_, err = NewResourceFailure(someValidResourceFailureReason, noName)
		assert.Error(t, err, "expected an empty name object to be rejected")

		noNamespace := validCausingObject()
		noNamespace.Namespace = ""
		_, err = NewResourceFailure(someValidResourceFailureReason, noNamespace)
		assert.NoError(t, err, "expected an empty namespace object to also be accepted")
	})
}

func TestResourceFailuresCollector(t *testing.T) {
	t.Run("is created when logger valid", func(t *testing.T) {
		logger := zapr.NewLogger(zap.NewNop())

		collector := NewResourceFailuresCollector(logger)
		require.NotNil(t, collector)
	})

	t.Run("pushes, logs and pops resource failures", func(t *testing.T) {
		core, logs := observer.New(zap.InfoLevel)
		logger := zapr.NewLogger(zap.New(core))

		collector := NewResourceFailuresCollector(logger)

		collector.PushResourceFailure(someValidResourceFailureReason, someResourceFailureCausingObjects()...)
		collector.PushResourceFailure(someValidResourceFailureReason, someResourceFailureCausingObjects()...)

		numberOfCausingObjects := len(someResourceFailureCausingObjects())
		require.Equal(t, logs.Len(), numberOfCausingObjects*2, "expecting one log entry per causing object")
		assertErrorLogs(t, logs)

		collectedErrors := collector.PopResourceFailures()
		require.Len(t, collectedErrors, 2)
		require.Empty(t, collector.PopResourceFailures(), "second call should not return any failure")
	})

	t.Run("does not crash but logs error when no causing objects passed", func(t *testing.T) {
		core, logs := observer.New(zap.DebugLevel)
		logger := zapr.NewLogger(zap.New(core))

		collector := NewResourceFailuresCollector(logger)

		collector.PushResourceFailure(someValidResourceFailureReason)

		require.NotZero(t, logs.Len())
		lastLog := logs.All()[logs.Len()-1]
		require.NotNil(t, lastLog)
		require.Equal(t, zap.ErrorLevel, lastLog.Level)
		require.Len(t, collector.PopResourceFailures(), 0, "no failures expected - causing objects missing")
	})
}

func assertErrorLogs(t *testing.T, logs *observer.ObservedLogs) {
	for i := range logs.All() {
		assert.Equalf(t, zapcore.ErrorLevel, logs.All()[i].Entry.Level, "%d-nth log entry expected to have ErrorLevel", i)
	}
}

func someResourceFailureCausingObjects() []client.Object {
	return []client.Object{validCausingObject(), validCausingObject()}
}

func validCausingObject() *kongv1.KongPlugin {
	return &kongv1.KongPlugin{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KongPlugin",
			APIVersion: kongv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "plugin-name",
			Namespace: "default",
		},
	}
}
