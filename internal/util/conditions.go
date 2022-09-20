package util

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ConditionType string

type ConditionReason string

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
			cond.ObservedGeneration >= resourceGeneration {
			return true
		}
	}
	return false
}
