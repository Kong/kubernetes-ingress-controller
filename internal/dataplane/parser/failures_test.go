package parser_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

const someValidTranslationFailureReason = "some valid reason"

func someValidTranslationFailureCausingObjects() []client.Object {
	return []client.Object{&kongv1.KongIngress{}, &kongv1.KongPlugin{}}
}

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
		_, err := parser.NewTranslationFailure(someValidTranslationFailureReason, someValidTranslationFailureCausingObjects()[0])
		require.NoError(t, err)

		_, err = parser.NewTranslationFailure(someValidTranslationFailureReason)
		require.Error(t, err)
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

	t.Run("pushes and pops translation failures", func(t *testing.T) {
		collector, err := parser.NewTranslationFailuresCollector(testLogger)
		require.NoError(t, err)

		collector.PushTranslationFailure(someValidTranslationFailureReason, someValidTranslationFailureCausingObjects()...)
		collector.PushTranslationFailure(someValidTranslationFailureReason, someValidTranslationFailureCausingObjects()...)

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
