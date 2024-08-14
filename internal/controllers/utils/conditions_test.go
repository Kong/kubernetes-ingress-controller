package utils_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

func TestEnsureProgrammedCondition(t *testing.T) {
	const testObjectGeneration = 2
	var (
		expectedProgrammedConditionTrue = metav1.Condition{
			Type:               string(kongv1.ConditionProgrammed),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: testObjectGeneration,
			Reason:             string(kongv1.ReasonProgrammed),
			Message:            utils.ProgrammedConditionTrueMessage,
		}
		expectedProgrammedConditionFalse = metav1.Condition{
			Type:               string(kongv1.ConditionProgrammed),
			Status:             metav1.ConditionFalse,
			ObservedGeneration: testObjectGeneration,
			Reason:             string(kongv1.ReasonInvalid),
			Message:            utils.ProgrammedConditionFalseInvalidMessage,
		}
		expectedProgrammedConditionUnknown = metav1.Condition{
			Type:               string(kongv1.ConditionProgrammed),
			Status:             metav1.ConditionFalse,
			ObservedGeneration: testObjectGeneration,
			Reason:             string(kongv1.ReasonPending),
			Message:            utils.ProgrammedConditionFalsePendingMessage,
		}
	)

	testCases := []struct {
		name string

		configurationStatus object.ConfigurationStatus
		conditions          []metav1.Condition

		expectedUpdatedConditions []metav1.Condition
		expectedUpdateNeeded      bool
	}{
		{
			name:                      "condition already present with correct status and observed generation",
			configurationStatus:       object.ConfigurationStatusSucceeded,
			conditions:                []metav1.Condition{expectedProgrammedConditionTrue},
			expectedUpdatedConditions: []metav1.Condition{expectedProgrammedConditionTrue},
			expectedUpdateNeeded:      false,
		},
		{
			name:                "condition present with correct status but older observed generation",
			configurationStatus: object.ConfigurationStatusSucceeded,
			conditions: []metav1.Condition{
				func() metav1.Condition {
					cond := expectedProgrammedConditionTrue
					cond.ObservedGeneration = 1
					return cond
				}(),
			},
			expectedUpdatedConditions: []metav1.Condition{expectedProgrammedConditionTrue},
			expectedUpdateNeeded:      true,
		},
		{
			name:                "condition present with correct observed generation but different status",
			configurationStatus: object.ConfigurationStatusFailed,
			conditions: []metav1.Condition{
				func() metav1.Condition {
					cond := expectedProgrammedConditionFalse
					cond.Status = metav1.ConditionTrue
					return cond
				}(),
			},
			expectedUpdatedConditions: []metav1.Condition{expectedProgrammedConditionFalse},
			expectedUpdateNeeded:      true,
		},
		{
			name:                "condition present with correct observed generation but different reason",
			configurationStatus: object.ConfigurationStatusFailed,
			conditions: []metav1.Condition{
				func() metav1.Condition {
					cond := expectedProgrammedConditionFalse
					cond.Reason = string("SomeOtherReason")
					return cond
				}(),
			},
			expectedUpdatedConditions: []metav1.Condition{expectedProgrammedConditionFalse},
			expectedUpdateNeeded:      true,
		},
		{
			name:                      "Unknown status should not modify existing Programmed condition",
			configurationStatus:       object.ConfigurationStatusUnknown,
			conditions:                []metav1.Condition{expectedProgrammedConditionTrue},
			expectedUpdatedConditions: []metav1.Condition{expectedProgrammedConditionTrue},
			expectedUpdateNeeded:      false,
		},
		{
			name:                      "empty conditions",
			configurationStatus:       object.ConfigurationStatusSucceeded,
			conditions:                nil,
			expectedUpdatedConditions: []metav1.Condition{expectedProgrammedConditionTrue},
			expectedUpdateNeeded:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conditions, updateNeeded := utils.EnsureProgrammedCondition(tc.configurationStatus, testObjectGeneration, tc.conditions)
			assert.Equal(t, tc.expectedUpdateNeeded, updateNeeded)

			ignoreLastTransitionTime := cmpopts.IgnoreFields(metav1.Condition{}, "LastTransitionTime")
			diff := cmp.Diff(conditions, tc.expectedUpdatedConditions, ignoreLastTransitionTime)
			assert.Empty(t, diff, "conditions mismatch")
		})
	}
}
