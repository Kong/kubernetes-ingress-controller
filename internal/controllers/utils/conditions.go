package utils

import (
	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	"github.com/samber/lo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object"
)

const (
	// ProgrammedConditionTrueMessage is the message for the programmed condition when it is True.
	ProgrammedConditionTrueMessage = "Object was successfully configured in Kong."

	// ProgrammedConditionFalseInvalidMessage is the message for the programmed condition when it is False with reason Invalid.
	ProgrammedConditionFalseInvalidMessage = "Object failed to be configured in Kong - see its attached Events for more information."

	// ProgrammedConditionFalsePendingMessage is the message for the programmed condition when it is False with reason Pending.
	ProgrammedConditionFalsePendingMessage = "Object is pending configuration in Kong."
)

type ProgrammedConditionOption func(object.ConfigurationStatus, *metav1.Condition)

// WithUnknownMessage sets the message of the desired Programmed condition to the given message if the
// configuration status is Unknown.
func WithUnknownMessage(message string) ProgrammedConditionOption {
	return func(status object.ConfigurationStatus, condition *metav1.Condition) {
		if status == object.ConfigurationStatusUnknown {
			condition.Message = message
		}
	}
}

// EnsureProgrammedCondition ensures that the programmed condition is present in the conditions slice with the
// status reflecting the current configuration status of the object.
// If the condition is already present with the correct status, the conditions slice is returned unmodified and false is
// returned as the second return value. If the condition is not present or has the wrong status, the conditions slice is
// returned with the condition updated and true is returned.
func EnsureProgrammedCondition(
	configurationStatus object.ConfigurationStatus,
	objectGeneration int64,
	conditions []metav1.Condition,
	options ...ProgrammedConditionOption,
) (
	updatedConditions []metav1.Condition,
	updateNeeded bool,
) {
	var (
		status  metav1.ConditionStatus
		reason  kongv1.ConditionReason
		message string
	)
	switch configurationStatus {
	case object.ConfigurationStatusSucceeded:
		status = metav1.ConditionTrue
		reason = kongv1.ReasonProgrammed
		message = ProgrammedConditionTrueMessage
	case object.ConfigurationStatusFailed:
		status = metav1.ConditionFalse
		reason = kongv1.ReasonInvalid
		message = ProgrammedConditionFalseInvalidMessage
	case object.ConfigurationStatusUnknown:
		status = metav1.ConditionFalse
		reason = kongv1.ReasonPending
		message = ProgrammedConditionFalsePendingMessage
	}

	desiredCondition := metav1.Condition{
		Type:               string(kongv1.ConditionProgrammed),
		Status:             status,
		ObservedGeneration: objectGeneration,
		LastTransitionTime: metav1.Now(),
		Reason:             string(reason),
		Message:            message,
	}
	for _, opt := range options {
		opt(configurationStatus, &desiredCondition)
	}

	hasMatchingCondition := util.CheckCondition(
		conditions,
		util.ConditionType(desiredCondition.Type),
		util.ConditionReason(desiredCondition.Reason),
		desiredCondition.Status,
		desiredCondition.ObservedGeneration,
	)

	if hasMatchingCondition {
		return conditions, false
	}

	_, idx, ok := lo.FindIndexOf(conditions, func(c metav1.Condition) bool { return c.Type == string(kongv1.ConditionProgrammed) })
	if !ok {
		conditions = append(conditions, desiredCondition)
	} else {
		// Do not update existing "Programmed" condition to Unknown to prevent races on updating status when new instance starts.
		if configurationStatus == object.ConfigurationStatusUnknown {
			return conditions, false
		}
		conditions[idx] = desiredCondition
	}

	return conditions, true
}
