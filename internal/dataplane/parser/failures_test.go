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

var someValidTranslationFailureCausingObjects = []client.Object{&kongv1.KongIngress{}, &kongv1.KongPlugin{}}

func TestTranslationFailure(t *testing.T) {
	t.Run("is_created_and_returns_reason_and_causing_objects", func(t *testing.T) {
		transErr, err := parser.NewTranslationFailure(someValidTranslationFailureReason, someValidTranslationFailureCausingObjects...)
		require.NoError(t, err)

		assert.Equal(t, someValidTranslationFailureReason, transErr.Reason())
		assert.ElementsMatch(t, someValidTranslationFailureCausingObjects, transErr.CausingObjects())
	})

	t.Run("fallbacks_to_unknown_reason_when_empty", func(t *testing.T) {
		transErr, err := parser.NewTranslationFailure("", someValidTranslationFailureCausingObjects...)
		require.NoError(t, err)
		require.Equal(t, parser.TranslationFailureReasonUnknown, transErr.Reason())
	})

	t.Run("requires_at_least_one_causing_object", func(t *testing.T) {
		_, err := parser.NewTranslationFailure(someValidTranslationFailureReason, someValidTranslationFailureCausingObjects[0])
		require.NoError(t, err)

		_, err = parser.NewTranslationFailure(someValidTranslationFailureReason)
		require.Error(t, err)
	})
}

func TestTranslationFailuresCollector(t *testing.T) {
	testLogger, _ := test.NewNullLogger()

	t.Run("is_created_when_logger_valid", func(t *testing.T) {
		collector, err := parser.NewTranslationFailuresCollector(testLogger)
		require.NoError(t, err)
		require.NotNil(t, collector)
	})

	t.Run("requires_non_nil_logger", func(t *testing.T) {
		_, err := parser.NewTranslationFailuresCollector(nil)
		require.Error(t, err)
	})

	t.Run("pushes_and_pops_translation_failures", func(t *testing.T) {
		collector, err := parser.NewTranslationFailuresCollector(testLogger)
		require.NoError(t, err)

		collector.PushTranslationFailure(someValidTranslationFailureReason, someValidTranslationFailureCausingObjects...)
		collector.PushTranslationFailure(someValidTranslationFailureReason, someValidTranslationFailureCausingObjects...)

		collectedErrors := collector.PopTranslationFailures()
		require.Len(t, collectedErrors, 2)
		require.Empty(t, collector.PopTranslationFailures(), "second call should not return any failure")
	})

	t.Run("does_not_crash_but_logs_warning_when_no_causing_objects_passed", func(t *testing.T) {
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
