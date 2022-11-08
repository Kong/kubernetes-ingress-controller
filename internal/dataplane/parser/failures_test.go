package parser_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

const someValidTranslationFailureReason = "some valid reason"

func TestTranslationFailure(t *testing.T) {
	t.Run("is created and returns reason and causing objects", func(t *testing.T) {
		transErr, err := parser.NewTranslationFailure(someValidTranslationFailureReason, someValidTranslationFailureCausingObjects()...)
		require.NoError(t, err)

		assert.Equal(t, someValidTranslationFailureReason, transErr.Reason())
		assert.ElementsMatch(t, someValidTranslationFailureCausingObjects(), transErr.CausingObjects())
	})

	t.Run("fallbacks to unknown reason when empty", func(t *testing.T) {
		transErr, err := parser.NewTranslationFailure("", someValidTranslationFailureCausingObjects()...)
		require.NoError(t, err)
		require.Equal(t, parser.TranslationFailureReasonUnknown, transErr.Reason())
	})

	t.Run("requires at least one causing object", func(t *testing.T) {
		_, err := parser.NewTranslationFailure(someValidTranslationFailureReason, validCausingObject())
		require.NoError(t, err)

		_, err = parser.NewTranslationFailure(someValidTranslationFailureReason)
		require.Error(t, err)
	})

	t.Run("requires valid objects", func(t *testing.T) {
		_, err := parser.NewTranslationFailure(someValidTranslationFailureReason, nil)
		assert.Error(t, err, "expected a nil object to be rejected")

		emptyGVK := validCausingObject()
		emptyGVK.APIVersion = ""
		emptyGVK.Kind = ""
		_, err = parser.NewTranslationFailure(someValidTranslationFailureReason, emptyGVK)
		assert.Error(t, err, "expected an empty GVK object to be rejected")

		noName := validCausingObject()
		noName.Name = ""
		_, err = parser.NewTranslationFailure(someValidTranslationFailureReason, noName)
		assert.Error(t, err, "expected an empty name object to be rejected")

		noNamespace := validCausingObject()
		noNamespace.Namespace = ""
		_, err = parser.NewTranslationFailure(someValidTranslationFailureReason, noNamespace)
		assert.Error(t, err, "expected an empty namespace object to be rejected")
	})
}

func TestTranslationFailuresCollector(t *testing.T) {
	testLogger, _ := test.NewNullLogger()

	t.Run("is created when logger valid", func(t *testing.T) {
		collector, err := parser.NewTranslationFailuresCollector(testLogger)
		require.NoError(t, err)
		require.NotNil(t, collector)
	})

	t.Run("requires non nil logger", func(t *testing.T) {
		_, err := parser.NewTranslationFailuresCollector(nil)
		require.Error(t, err)
	})

	t.Run("pushes, logs and pops translation failures", func(t *testing.T) {
		logger, loggerHook := test.NewNullLogger()
		collector, err := parser.NewTranslationFailuresCollector(logger)
		require.NoError(t, err)

		collector.PushTranslationFailure(someValidTranslationFailureReason, someValidTranslationFailureCausingObjects()...)
		collector.PushTranslationFailure(someValidTranslationFailureReason, someValidTranslationFailureCausingObjects()...)

		numberOfCausingObjects := len(someValidTranslationFailureCausingObjects())
		require.Len(t, loggerHook.AllEntries(), numberOfCausingObjects*2, "expecting one log entry per causing object")
		assertErrorLogs(t, loggerHook)

		collectedErrors := collector.PopTranslationFailures()
		require.Len(t, collectedErrors, 2)
		require.Empty(t, collector.PopTranslationFailures(), "second call should not return any failure")
	})

	t.Run("does not crash but logs warning when no causing objects passed", func(t *testing.T) {
		logger, loggerHook := test.NewNullLogger()
		collector, err := parser.NewTranslationFailuresCollector(logger)
		require.NoError(t, err)

		collector.PushTranslationFailure(someValidTranslationFailureReason)

		lastLog := loggerHook.LastEntry()
		require.NotNil(t, lastLog)
		require.Equal(t, logrus.WarnLevel, lastLog.Level)
		require.Len(t, collector.PopTranslationFailures(), 0, "no failures expected - causing objects missing")
	})
}

func assertErrorLogs(t *testing.T, logHook *test.Hook) {
	for i := range logHook.AllEntries() {
		assert.Equalf(t, logrus.ErrorLevel, logHook.AllEntries()[i].Level, "%d-nth log entry expected to have ErrorLevel", i)
	}
}

func someValidTranslationFailureCausingObjects() []client.Object {
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
