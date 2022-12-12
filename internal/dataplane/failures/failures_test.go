package failures

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
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
		assert.Error(t, err, "expected an empty namespace object to be rejected")
	})
}

func TestResourceFailuresCollector(t *testing.T) {
	testLogger, _ := test.NewNullLogger()

	t.Run("is created when logger valid", func(t *testing.T) {
		collector, err := NewResourceFailuresCollector(testLogger)
		require.NoError(t, err)
		require.NotNil(t, collector)
	})

	t.Run("requires non nil logger", func(t *testing.T) {
		_, err := NewResourceFailuresCollector(nil)
		require.Error(t, err)
	})

	t.Run("pushes, logs and pops resource failures", func(t *testing.T) {
		logger, loggerHook := test.NewNullLogger()
		collector, err := NewResourceFailuresCollector(logger)
		require.NoError(t, err)

		collector.PushResourceFailure(someValidResourceFailureReason, someResourceFailureCausingObjects()...)
		collector.PushResourceFailure(someValidResourceFailureReason, someResourceFailureCausingObjects()...)

		numberOfCausingObjects := len(someResourceFailureCausingObjects())
		require.Len(t, loggerHook.AllEntries(), numberOfCausingObjects*2, "expecting one log entry per causing object")
		assertErrorLogs(t, loggerHook)

		collectedErrors := collector.PopResourceFailures()
		require.Len(t, collectedErrors, 2)
		require.Empty(t, collector.PopResourceFailures(), "second call should not return any failure")
	})

	t.Run("does not crash but logs warning when no causing objects passed", func(t *testing.T) {
		logger, loggerHook := test.NewNullLogger()
		collector, err := NewResourceFailuresCollector(logger)
		require.NoError(t, err)

		collector.PushResourceFailure(someValidResourceFailureReason)

		lastLog := loggerHook.LastEntry()
		require.NotNil(t, lastLog)
		require.Equal(t, logrus.WarnLevel, lastLog.Level)
		require.Len(t, collector.PopResourceFailures(), 0, "no failures expected - causing objects missing")
	})
}

func assertErrorLogs(t *testing.T, logHook *test.Hook) {
	for i := range logHook.AllEntries() {
		assert.Equalf(t, logrus.ErrorLevel, logHook.AllEntries()[i].Level, "%d-nth log entry expected to have ErrorLevel", i)
	}
}

func someResourceFailureCausingObjects() []client.Object {
	return []client.Object{validCausingObject(), validCausingObject()}
}

func validCausingObject() *v1.KongPlugin {
	return &v1.KongPlugin{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KongPlugin",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "plugin-name",
			Namespace: "default",
		},
	}
}
