package parser

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TranslationError represents an error occurring during translating Kubernetes objects into Kong ones.
// It can be associated with one or more Kubernetes objects.
type TranslationError struct {
	causingObjects []client.Object
	reason         string
}

// NewTranslationError creates a TranslationError with a reason that should be a human-readable explanation
// of the error reason, and an causingObjects slice that specifies what objects have caused the error.
func NewTranslationError(reason string, causingObjects ...client.Object) (TranslationError, error) {
	if reason == "" {
		reason = "unknown"
	}
	if len(causingObjects) < 1 {
		return TranslationError{}, fmt.Errorf("no causing objects specified, reason: %s", reason)
	}

	return TranslationError{
		causingObjects: causingObjects,
		reason:         reason,
	}, nil
}

// CausingObjects returns a slice of objects that have caused the translation error.
func (p TranslationError) CausingObjects() []client.Object {
	return p.causingObjects
}

// Reason returns a human-readable reason of the error.
func (p TranslationError) Reason() string {
	return p.reason
}

// TranslationErrorsCollector should be used to collect all translation errors that happen during the translation process.
type TranslationErrorsCollector struct {
	errors []TranslationError
	logger logrus.FieldLogger
}

func NewTranslationErrorsCollector(logger logrus.FieldLogger) (*TranslationErrorsCollector, error) {
	if logger == nil {
		return nil, errors.New("missing logger")
	}
	return &TranslationErrorsCollector{logger: logger}, nil
}

// PushTranslationError registers a translation error.
func (c *TranslationErrorsCollector) PushTranslationError(reason string, causingObjects ...client.Object) {
	translationErr, err := NewTranslationError(reason, causingObjects...)
	if err != nil {
		c.logger.Warningf("failed to create translation error: %w", err)
		return
	}

	c.errors = append(c.errors, translationErr)
}

// PopTranslationErrors returns all translation errors that occurred during the translation process and erases them
// in the collector. It makes the collector reusable during next translation runs.
func (c *TranslationErrorsCollector) PopTranslationErrors() []TranslationError {
	errs := c.errors
	c.errors = nil

	return errs
}
