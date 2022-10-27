package parser

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// TranslationFailureReasonUnknown is used when no specific reason is specified when creating a TranslationFailure.
	TranslationFailureReasonUnknown = "unknown"
)

// TranslationFailure represents an error occurring during translating Kubernetes objects into Kong ones.
// It can be associated with one or more Kubernetes objects.
type TranslationFailure struct {
	causingObjects []client.Object
	reason         string
}

// NewTranslationFailure creates a TranslationFailure with a reason that should be a human-readable explanation
// of the error reason, and a causingObjects slice that specifies what objects have caused the error.
func NewTranslationFailure(reason string, causingObjects ...client.Object) (TranslationFailure, error) {
	if reason == "" {
		reason = TranslationFailureReasonUnknown
	}
	if len(causingObjects) < 1 {
		return TranslationFailure{}, fmt.Errorf("no causing objects specified, reason: %s", reason)
	}

	return TranslationFailure{
		causingObjects: causingObjects,
		reason:         reason,
	}, nil
}

// CausingObjects returns a slice of objects that have caused the translation error.
func (p TranslationFailure) CausingObjects() []client.Object {
	return p.causingObjects
}

// Reason returns a human-readable reason of the failure.
func (p TranslationFailure) Reason() string {
	return p.reason
}

// TranslationFailuresCollector should be used to collect all translation failures that happen during the translation process.
type TranslationFailuresCollector struct {
	failures []TranslationFailure
	logger   logrus.FieldLogger
}

func NewTranslationFailuresCollector(logger logrus.FieldLogger) (*TranslationFailuresCollector, error) {
	if logger == nil {
		return nil, errors.New("missing logger")
	}
	return &TranslationFailuresCollector{logger: logger}, nil
}

// PushTranslationFailure registers a translation failure.
func (c *TranslationFailuresCollector) PushTranslationFailure(reason string, causingObjects ...client.Object) {
	translationErr, err := NewTranslationFailure(reason, causingObjects...)
	if err != nil {
		c.logger.Warningf("failed to create translation failure: %w", err)
		return
	}

	c.failures = append(c.failures, translationErr)
}

// PopTranslationFailures returns all translation failures that occurred during the translation process and erases them
// in the collector. It makes the collector reusable during next translation runs.
func (c *TranslationFailuresCollector) PopTranslationFailures() []TranslationFailure {
	errs := c.failures
	c.failures = nil

	return errs
}
