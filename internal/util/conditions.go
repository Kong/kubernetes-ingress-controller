package util

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// ConditionType can be any condition type, e.g. `gatewayv1.GatewayConditionProgrammed`.
type ConditionType string

// ConditionReason can be any condition reason, e.g. `gatewayv1.GatewayReasonProgrammed`.
type ConditionReason string

// CheckCondition tells if there's a condition matching the given type, reason, and status in conditions.
// It also makes sure that the condition's observed generation is no older than the resource's actual generation.
func CheckCondition(
	conditions []metav1.Condition,
	typ ConditionType,
	reason ConditionReason,
	status metav1.ConditionStatus,
	resourceGeneration int64,
) bool {
	for _, cond := range conditions {
		if cond.Type == string(typ) &&
			cond.Reason == string(reason) &&
			cond.Status == status &&
			cond.ObservedGeneration == resourceGeneration {
			return true
		}
	}
	return false
}
