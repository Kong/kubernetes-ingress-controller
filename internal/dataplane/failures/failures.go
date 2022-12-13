package failures

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// ResourceFailureReasonUnknown is used when no specific message is specified when creating a ResourceFailure.
	ResourceFailureReasonUnknown = "unknown"
)

// ResourceFailure represents an error encountered when processing one or more Kubernetes resources into Kong
// configuration.
type ResourceFailure struct {
	causingObjects []client.Object
	message        string
}

// NewResourceFailure creates a ResourceFailure with a message that should be a human-readable explanation
// of the error message, and a causingObjects slice that specifies what objects have caused the error.
func NewResourceFailure(reason string, causingObjects ...client.Object) (ResourceFailure, error) {
	if reason == "" {
		reason = ResourceFailureReasonUnknown
	}
	if len(causingObjects) < 1 {
		return ResourceFailure{}, fmt.Errorf("no causing objects specified, message: %s", reason)
	}

	for _, obj := range causingObjects {
		if obj == nil {
			return ResourceFailure{}, errors.New("one of causing objects is nil")
		}
		gvk := obj.GetObjectKind().GroupVersionKind()
		if gvk.Empty() {
			return ResourceFailure{}, errors.New("one of causing objects has an empty GVK")
		}
		if obj.GetName() == "" {
			return ResourceFailure{}, fmt.Errorf("one of causing objects (%s) has no name", gvk.String())
		}
		if obj.GetNamespace() == "" {
			return ResourceFailure{}, fmt.Errorf("one of causing objects (%s) has no namespace", gvk.String())
		}
	}

	return ResourceFailure{
		causingObjects: causingObjects,
		message:        reason,
	}, nil
}

// CausingObjects returns a slice of objects involved in a resource processing failure.
func (p ResourceFailure) CausingObjects() []client.Object {
	return p.causingObjects
}

// Message returns a human-readable message describing the cause of the failure.
func (p ResourceFailure) Message() string {
	return p.message
}

// ResourceFailuresCollector collects resource failures across different stages of resource processing.
type ResourceFailuresCollector struct {
	failures []ResourceFailure
	logger   logrus.FieldLogger
}

func NewResourceFailuresCollector(logger logrus.FieldLogger) (*ResourceFailuresCollector, error) {
	if logger == nil {
		return nil, errors.New("missing logger")
	}
	return &ResourceFailuresCollector{logger: logger}, nil
}

// PushResourceFailure adds a resource processing failure to the collector and logs it.
func (c *ResourceFailuresCollector) PushResourceFailure(reason string, causingObjects ...client.Object) {
	resourceFailure, err := NewResourceFailure(reason, causingObjects...)
	if err != nil {
		c.logger.WithField("resource_failure_reason", reason).Warningf("failed to create resource failure: %s", err)
		return
	}

	c.failures = append(c.failures, resourceFailure)
	c.logResourceFailure(reason, causingObjects...)
}

// logResourceFailure logs an error with a resource processing failure message for each causing object.
func (c *ResourceFailuresCollector) logResourceFailure(reason string, causingObjects ...client.Object) {
	for _, obj := range causingObjects {
		c.logger.WithFields(logrus.Fields{
			"name":      obj.GetName(),
			"namespace": obj.GetNamespace(),
			"GVK":       obj.GetObjectKind().GroupVersionKind().String(),
		}).Errorf("resource processing failed: %s", reason)
	}
}

// PopResourceFailures returns all resource processing failures stored in the collector and clears the collector's
// stored failures. The collector can then be reused for the next iteration of the process it tracks.
func (c *ResourceFailuresCollector) PopResourceFailures() []ResourceFailure {
	errs := c.failures
	c.failures = nil

	return errs
}
