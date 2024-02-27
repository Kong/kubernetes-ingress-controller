package gateway

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func TestReadyConditionExistsForObservedGeneration(t *testing.T) {
	t.Log("checking programmed condition for currently ready gateway")
	currentlyProgrammedGateway := &gatewayapi.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 1,
		},
		Status: gatewayapi.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayapi.GatewayConditionProgrammed),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayapi.GatewayReasonProgrammed),
			}},
		},
	}
	assert.True(t, isGatewayProgrammed(currentlyProgrammedGateway))

	t.Log("checking programmed condition for previously programmed gateway that has since been updated")
	previouslyProgrammedGateway := &gatewayapi.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 2,
		},
		Status: gatewayapi.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayapi.GatewayConditionProgrammed),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayapi.GatewayReasonProgrammed),
			}},
		},
	}
	assert.False(t, isGatewayProgrammed(previouslyProgrammedGateway))

	t.Log("checking programmed condition for a gateway which has never been ready")
	neverBeenProgrammedGateway := &gatewayapi.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 10,
		},
		Status: gatewayapi.GatewayStatus{},
	}
	assert.False(t, isGatewayProgrammed(neverBeenProgrammedGateway))
}

func TestSetGatewayCondtion(t *testing.T) {
	testCases := []struct {
		name            string
		gw              *gatewayapi.Gateway
		condition       metav1.Condition
		conditionLength int
	}{
		{
			name: "no_such_condition_should_append_one",
			gw:   &gatewayapi.Gateway{},
			condition: metav1.Condition{
				Type:               "fake1",
				Status:             metav1.ConditionTrue,
				Reason:             "fake1",
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
			},
			conditionLength: 1,
		},
		{
			name: "have_condition_with_type_should_replace",
			gw: &gatewayapi.Gateway{
				Status: gatewayapi.GatewayStatus{
					Conditions: []metav1.Condition{
						{
							Type:               "fake1",
							Status:             metav1.ConditionFalse,
							Reason:             "fake1",
							ObservedGeneration: 1,
							LastTransitionTime: metav1.Now(),
						},
					},
				},
			},
			condition: metav1.Condition{
				Type:               "fake1",
				Status:             metav1.ConditionTrue,
				Reason:             "fake1",
				ObservedGeneration: 2,
				LastTransitionTime: metav1.Now(),
			},
			conditionLength: 1,
		},
		{
			name: "multiple_conditions_with_type_should_preserve_one",
			gw: &gatewayapi.Gateway{
				Status: gatewayapi.GatewayStatus{
					Conditions: []metav1.Condition{
						{
							Type:               "fake1",
							Status:             metav1.ConditionFalse,
							Reason:             "fake1",
							ObservedGeneration: 1,
							LastTransitionTime: metav1.Now(),
						},
						{
							Type:               "fake1",
							Status:             metav1.ConditionTrue,
							Reason:             "fake2",
							ObservedGeneration: 2,
							LastTransitionTime: metav1.Now(),
						},
						{
							Type:               "fake2",
							Status:             metav1.ConditionTrue,
							Reason:             "fake2",
							ObservedGeneration: 2,
							LastTransitionTime: metav1.Now(),
						},
					},
				},
			},
			condition: metav1.Condition{
				Type:               "fake1",
				Status:             metav1.ConditionTrue,
				Reason:             "fake1",
				ObservedGeneration: 3,
				LastTransitionTime: metav1.Now(),
			},
			conditionLength: 2,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			setGatewayCondition(tc.gw, tc.condition)
			t.Logf("checking conditions of gateway after setting")
			assert.Len(t, tc.gw.Status.Conditions, tc.conditionLength)

			conditionNum := 0
			var observedCondition metav1.Condition
			for _, condition := range tc.gw.Status.Conditions {
				if condition.Type == tc.condition.Type {
					conditionNum++
					observedCondition = condition
				}
			}
			assert.Equal(t, 1, conditionNum)
			assert.EqualValues(t, tc.condition, observedCondition)
		})
	}
}

func TestIsGatewayMarkedAsAccepted(t *testing.T) {
	t.Log("verifying scheduled check for gateway object which has been accepted")
	scheduledGateway := &gatewayapi.Gateway{
		ObjectMeta: metav1.ObjectMeta{Generation: 1},
		Status: gatewayapi.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayapi.GatewayConditionAccepted),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayapi.GatewayReasonAccepted),
			}},
		},
	}
	assert.True(t, isGatewayAccepted(scheduledGateway))

	t.Log("verifying scheduled check for gateway object which has not been scheduled")
	unscheduledGateway := &gatewayapi.Gateway{}
	assert.False(t, isGatewayAccepted(unscheduledGateway))
}

func TestPruneStatusConditions(t *testing.T) {
	t.Log("verifying that a gateway with minimal status conditions is not pruned")
	gateway := &gatewayapi.Gateway{}
	for i := 0; i < 4; i++ {
		gateway.Status.Conditions = append(gateway.Status.Conditions, metav1.Condition{Type: "fake", ObservedGeneration: int64(i)})
	}
	assert.Len(t, pruneGatewayStatusConds(gateway).Status.Conditions, 4)
	assert.Len(t, gateway.Status.Conditions, 4)

	t.Log("verifying that a gateway with the maximum allowed number of conditions is note pruned")
	for i := 0; i < 4; i++ {
		gateway.Status.Conditions = append(gateway.Status.Conditions, metav1.Condition{Type: "fake", ObservedGeneration: int64(i) + 4})
	}
	assert.Len(t, pruneGatewayStatusConds(gateway).Status.Conditions, 8)
	assert.Len(t, gateway.Status.Conditions, 8)

	t.Log("verifying that a gateway with too many status conditions is pruned")
	for i := 0; i < 4; i++ {
		gateway.Status.Conditions = append(gateway.Status.Conditions, metav1.Condition{Type: "fake", ObservedGeneration: int64(i) + 8})
	}
	assert.Len(t, pruneGatewayStatusConds(gateway).Status.Conditions, 8)
	assert.Len(t, gateway.Status.Conditions, 8)

	t.Log("verifying that the more recent 8 conditions were retained after the pruning")
	assert.Equal(t, int64(4), gateway.Status.Conditions[0].ObservedGeneration)
	assert.Equal(t, int64(11), gateway.Status.Conditions[7].ObservedGeneration)
}

func TestReconcileGatewaysIfClassMatches(t *testing.T) {
	t.Log("generating a gatewayclass to test reconciliation filters")
	gatewayClass := &gatewayapi.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "us",
		},
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: GetControllerName(),
		},
	}

	t.Log("generating a list of matching controllers")
	matching := []gatewayapi.Gateway{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sanfrancisco",
				Namespace: "california",
			},
			Spec: gatewayapi.GatewaySpec{
				GatewayClassName: gatewayapi.ObjectName(gatewayClass.Name),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sandiego",
				Namespace: "california",
			},
			Spec: gatewayapi.GatewaySpec{
				GatewayClassName: gatewayapi.ObjectName(gatewayClass.Name),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "losangelos",
				Namespace: "california",
			},
			Spec: gatewayapi.GatewaySpec{
				GatewayClassName: gatewayapi.ObjectName(gatewayClass.Name),
			},
		},
	}

	t.Log("generating a list of non-matching controllers")
	nonmatching := []gatewayapi.Gateway{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hamburg",
				Namespace: "germany",
			},
			Spec: gatewayapi.GatewaySpec{
				GatewayClassName: gatewayapi.ObjectName("eu"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "paris",
				Namespace: "france",
			},
			Spec: gatewayapi.GatewaySpec{
				GatewayClassName: gatewayapi.ObjectName("eu"),
			},
		},
	}

	t.Log("verifying reconciliation counts")
	assert.Len(t, reconcileGatewaysIfClassMatches(gatewayClass, append(matching, nonmatching...)), len(matching))
	assert.Len(t, reconcileGatewaysIfClassMatches(gatewayClass, matching), len(matching))
	assert.Len(t, reconcileGatewaysIfClassMatches(gatewayClass, nonmatching), 0)
	assert.Len(t, reconcileGatewaysIfClassMatches(gatewayClass, nil), 0)

	t.Log("verifying reconciliation results")
	expected := []reconcile.Request{
		{
			NamespacedName: k8stypes.NamespacedName{
				Name:      "sanfrancisco",
				Namespace: "california",
			},
		},
		{
			NamespacedName: k8stypes.NamespacedName{
				Name:      "sandiego",
				Namespace: "california",
			},
		},
		{
			NamespacedName: k8stypes.NamespacedName{
				Name:      "losangelos",
				Namespace: "california",
			},
		},
	}
	assert.Equal(t, expected, reconcileGatewaysIfClassMatches(gatewayClass, append(matching, nonmatching...)))
	assert.Equal(t, expected, reconcileGatewaysIfClassMatches(gatewayClass, matching))
}

func TestIsGatewayControlled(t *testing.T) {
	var testControllerName gatewayapi.GatewayController = "acme.io/gateway-controller"

	testCases := []struct {
		name           string
		GatewayClass   *gatewayapi.GatewayClass
		expectedResult bool
	}{
		{
			name: "uncontrolled GatewayClass",
			GatewayClass: &gatewayapi.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "uncontrolled",
				},
				Spec: gatewayapi.GatewayClassSpec{
					ControllerName: testControllerName,
				},
			},
			expectedResult: false,
		},
		{
			name: "controlled GatewayClass",
			GatewayClass: &gatewayapi.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "controlled",
				},
				Spec: gatewayapi.GatewayClassSpec{
					ControllerName: GetControllerName(),
				},
			},
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expectedResult, isGatewayClassControlled(tc.GatewayClass))
		})
	}
}

func TestIsGatewayUnmanaged(t *testing.T) {
	testCases := []struct {
		name                    string
		GatewayClassAnnotations map[string]string
		expectedResult          bool
	}{
		{
			name: "unmanaged GatewayClass",
			GatewayClassAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.GatewayClassUnmanagedKey: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
			},
			expectedResult: true,
		},
		{
			name:                    "managed GatewayClass",
			GatewayClassAnnotations: map[string]string{},
			expectedResult:          false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expectedResult, isGatewayClassUnmanaged(tc.GatewayClassAnnotations))
		})
	}
}

func TestGetReferenceGrantConditionReason(t *testing.T) {
	testCases := []struct {
		name             string
		gatewayNamespace string
		certRef          gatewayapi.SecretObjectReference
		referenceGrants  []gatewayapi.ReferenceGrant
		expectedReason   string
	}{
		{
			name:           "empty reference",
			certRef:        gatewayapi.SecretObjectReference{},
			expectedReason: string(gatewayapi.ListenerReasonResolvedRefs),
		},
		{
			name:             "no need for reference",
			gatewayNamespace: "test",
			certRef: gatewayapi.SecretObjectReference{
				Kind: util.StringToGatewayAPIKindPtr("Secret"),
				Name: "testSecret",
			},
			expectedReason: string(gatewayapi.ListenerReasonResolvedRefs),
		},
		{
			name:             "reference not granted - secret name not matching",
			gatewayNamespace: "test",
			certRef: gatewayapi.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: lo.ToPtr(gatewayapi.Namespace("otherNamespace")),
			},
			referenceGrants: []gatewayapi.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "otherNamespace",
					},
					Spec: gatewayapi.ReferenceGrantSpec{
						From: []gatewayapi.ReferenceGrantFrom{
							{
								Group:     gatewayapi.V1Group,
								Kind:      "Gateway",
								Namespace: "test",
							},
						},
						To: []gatewayapi.ReferenceGrantTo{
							{
								Group: "",
								Kind:  "Secret",
								Name:  lo.ToPtr(gatewayapi.ObjectName("anotherSecret")),
							},
						},
					},
				},
			},
			expectedReason: string(gatewayapi.ListenerReasonRefNotPermitted),
		},
		{
			name:             "reference not granted - no grants specified",
			gatewayNamespace: "test",
			certRef: gatewayapi.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: lo.ToPtr(gatewayapi.Namespace("otherNamespace")),
			},
			expectedReason: string(gatewayapi.ListenerReasonRefNotPermitted),
		},
		{
			name:             "reference granted, secret name not specified",
			gatewayNamespace: "test",
			certRef: gatewayapi.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: lo.ToPtr(gatewayapi.Namespace("otherNamespace")),
			},
			referenceGrants: []gatewayapi.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "otherNamespace",
					},
					Spec: gatewayapi.ReferenceGrantSpec{
						From: []gatewayapi.ReferenceGrantFrom{
							// useless entry, just to furtherly test the function
							{
								Group:     "otherGroup",
								Kind:      "otherKind",
								Namespace: "test",
							},
							// good entry
							{
								Group:     gatewayapi.V1Group,
								Kind:      "Gateway",
								Namespace: "test",
							},
						},
						To: []gatewayapi.ReferenceGrantTo{
							{
								Group: "",
								Kind:  "Secret",
							},
						},
					},
				},
			},
			expectedReason: string(gatewayapi.ListenerReasonResolvedRefs),
		},
		{
			name:             "reference granted, secret name specified",
			gatewayNamespace: "test",
			certRef: gatewayapi.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: lo.ToPtr(gatewayapi.Namespace("otherNamespace")),
			},
			referenceGrants: []gatewayapi.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "otherNamespace",
					},
					Spec: gatewayapi.ReferenceGrantSpec{
						From: []gatewayapi.ReferenceGrantFrom{
							{
								Group:     gatewayapi.V1Group,
								Kind:      "Gateway",
								Namespace: "test",
							},
						},
						To: []gatewayapi.ReferenceGrantTo{
							{
								Group: "",
								Kind:  "Secret",
								Name:  lo.ToPtr(gatewayapi.ObjectName("testSecret")),
							},
						},
					},
				},
			},
			expectedReason: string(gatewayapi.ListenerReasonResolvedRefs),
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedReason, getReferenceGrantConditionReason(
			tc.gatewayNamespace,
			tc.certRef,
			tc.referenceGrants,
		))
	}
}
