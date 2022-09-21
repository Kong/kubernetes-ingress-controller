package util_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func TestCheckCondition(t *testing.T) {
	expectedType := util.ConditionType(gatewayv1alpha2.ListenerConditionReady)
	expectedReason := util.ConditionReason(gatewayv1alpha2.GatewayReasonScheduled)
	expectedStatus := metav1.ConditionTrue
	generation := int64(1)
	givenConditions := []metav1.Condition{
		{
			Type:               string(expectedType),
			Reason:             string(expectedReason),
			Status:             expectedStatus,
			ObservedGeneration: generation,
		},
	}

	otherType := util.ConditionType(gatewayv1alpha2.ListenerConditionConflicted)
	otherReason := util.ConditionReason(gatewayv1alpha2.GatewayReasonReady)
	otherStatus := metav1.ConditionFalse

	testCases := []struct {
		name           string
		givenType      util.ConditionType
		givenReason    util.ConditionReason
		givenStatus    metav1.ConditionStatus
		expectedResult bool
	}{
		{
			name:           "all_as_expected_should_give_true",
			givenType:      expectedType,
			givenReason:    expectedReason,
			givenStatus:    expectedStatus,
			expectedResult: true,
		},
		{
			name:           "all_but_type_as_expected_should_give_false",
			givenType:      otherType,
			givenReason:    expectedReason,
			givenStatus:    expectedStatus,
			expectedResult: false,
		},
		{
			name:           "all_but_reason_as_expected_should_give_false",
			givenType:      expectedType,
			givenReason:    otherReason,
			givenStatus:    expectedStatus,
			expectedResult: false,
		},
		{
			name:           "all_but_status_as_expected_should_give_false",
			givenType:      expectedType,
			givenReason:    expectedReason,
			givenStatus:    otherStatus,
			expectedResult: false,
		},
		{
			name:           "only_type_as_expected_should_give_false",
			givenType:      expectedType,
			givenReason:    otherReason,
			givenStatus:    otherStatus,
			expectedResult: false,
		},
		{
			name:           "only_reason_as_expected_should_give_false",
			givenType:      otherType,
			givenReason:    expectedReason,
			givenStatus:    otherStatus,
			expectedResult: false,
		},
		{
			name:           "only_status_as_expected_should_give_false",
			givenType:      otherType,
			givenReason:    otherReason,
			givenStatus:    expectedStatus,
			expectedResult: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ok := util.CheckCondition(
				givenConditions,
				testCase.givenType,
				testCase.givenReason,
				testCase.givenStatus,
				generation,
			)
			require.Equal(t, testCase.expectedResult, ok)
		})
	}
}

func TestCheckCondition_observed_generations_lower_than_actual_are_ignored(t *testing.T) {
	expectedType := util.ConditionType(gatewayv1alpha2.ListenerConditionReady)
	expectedReason := util.ConditionReason(gatewayv1alpha2.GatewayReasonScheduled)
	expectedStatus := metav1.ConditionTrue
	givenConditions := []metav1.Condition{
		{
			Type:               string(expectedType),
			Reason:             string(expectedReason),
			Status:             expectedStatus,
			ObservedGeneration: 1,
		},
	}
	generationHigherThanObserved := int64(2)

	ok := util.CheckCondition(
		givenConditions,
		expectedType,
		expectedReason,
		expectedStatus,
		generationHigherThanObserved,
	)
	require.False(t, ok, "expected to not match any condition due to low observed generation")
}
