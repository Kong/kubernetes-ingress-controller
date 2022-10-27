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

const someValidTranslationErrorReason = "some valid reason"

var someValidTranslationErrorCausingObjects = []client.Object{&kongv1.KongIngress{}, &kongv1.KongPlugin{}}

func TestTranslationError(t *testing.T) {
	t.Run("is_created_and_returns_reason_and_causing_objects", func(t *testing.T) {
		transErr, err := parser.NewTranslationError(someValidTranslationErrorReason, someValidTranslationErrorCausingObjects...)
		require.NoError(t, err)

		assert.Equal(t, someValidTranslationErrorReason, transErr.Reason())
		assert.ElementsMatch(t, someValidTranslationErrorCausingObjects, transErr.CausingObjects())
	})

	t.Run("fallbacks_to_unknown_reason_when_empty", func(t *testing.T) {
		transErr, err := parser.NewTranslationError("", someValidTranslationErrorCausingObjects...)
		require.NoError(t, err)
		require.Equal(t, "unknown", transErr.Reason())
	})

	t.Run("requires_at_least_one_causing_object", func(t *testing.T) {
		_, err := parser.NewTranslationError(someValidTranslationErrorReason, someValidTranslationErrorCausingObjects[0])
		require.NoError(t, err)

		_, err = parser.NewTranslationError(someValidTranslationErrorReason)
		require.Error(t, err)
	})
}

func TestTranslationErrorsCollector(t *testing.T) {
	testLogger, _ := test.NewNullLogger()

	t.Run("is_created_when_logger_valid", func(t *testing.T) {
		collector, err := parser.NewTranslationErrorsCollector(testLogger)
		require.NoError(t, err)
		require.NotNil(t, collector)
	})

	t.Run("requires_non_nil_logger", func(t *testing.T) {
		_, err := parser.NewTranslationErrorsCollector(nil)
		require.Error(t, err)
	})

	t.Run("pushes_and_pops_translation_errors", func(t *testing.T) {
		collector, err := parser.NewTranslationErrorsCollector(testLogger)
		require.NoError(t, err)

		collector.PushTranslationError(someValidTranslationErrorReason, someValidTranslationErrorCausingObjects...)
		collector.PushTranslationError(someValidTranslationErrorReason, someValidTranslationErrorCausingObjects...)

		collectedErrors := collector.PopTranslationErrors()
		require.Len(t, collectedErrors, 2)
		require.Empty(t, collector.PopTranslationErrors(), "second call should not return any error")
	})

	t.Run("does_not_crash_but_logs_warning_when_no_causing_objects_passed", func(t *testing.T) {
		logger, loggerHook := test.NewNullLogger()
		collector, err := parser.NewTranslationErrorsCollector(logger)
		require.NoError(t, err)

		collector.PushTranslationError(someValidTranslationErrorReason)

		lastLog := loggerHook.LastEntry()
		require.NotNil(t, lastLog)
		require.Equal(t, logrus.WarnLevel, lastLog.Level)
		require.Len(t, collector.PopTranslationErrors(), 0, "no errors expected - causing objects missing")
	})
}
