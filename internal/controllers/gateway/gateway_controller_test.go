package gateway

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func TestReadyConditionExistsForObservedGeneration(t *testing.T) {
	t.Log("checking programmed condition for currently ready gateway")
	currentlyProgrammedGateway := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 1,
		},
		Status: gatewayv1.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayv1.GatewayConditionProgrammed),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1.GatewayReasonProgrammed),
			}},
		},
	}
	assert.True(t, isGatewayProgrammed(currentlyProgrammedGateway))

	t.Log("checking programmed condition for previously programmed gateway that has since been updated")
	previouslyProgrammedGateway := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 2,
		},
		Status: gatewayv1.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayv1.GatewayConditionProgrammed),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1.GatewayReasonProgrammed),
			}},
		},
	}
	assert.False(t, isGatewayProgrammed(previouslyProgrammedGateway))

	t.Log("checking programmed condition for a gateway which has never been ready")
	neverBeenProgrammedGateway := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 10,
		},
		Status: gatewayv1.GatewayStatus{},
	}
	assert.False(t, isGatewayProgrammed(neverBeenProgrammedGateway))
}

func TestSetGatewayCondtion(t *testing.T) {
	testCases := []struct {
		name            string
		gw              *gatewayv1.Gateway
		condition       metav1.Condition
		conditionLength int
	}{
		{
			name: "no_such_condition_should_append_one",
			gw:   &gatewayv1.Gateway{},
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
			gw: &gatewayv1.Gateway{
				Status: gatewayv1.GatewayStatus{
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
			gw: &gatewayv1.Gateway{
				Status: gatewayv1.GatewayStatus{
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
	scheduledGateway := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{Generation: 1},
		Status: gatewayv1.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayv1.GatewayConditionAccepted),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1.GatewayReasonAccepted),
			}},
		},
	}
	assert.True(t, isGatewayScheduled(scheduledGateway))

	t.Log("verifying scheduled check for gateway object which has not been scheduled")
	unscheduledGateway := &gatewayv1.Gateway{}
	assert.False(t, isGatewayScheduled(unscheduledGateway))
}

func TestPruneStatusConditions(t *testing.T) {
	t.Log("verifying that a gateway with minimal status conditions is not pruned")
	gateway := &gatewayv1.Gateway{}
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
	gatewayClass := &gatewayv1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "us",
		},
		Spec: gatewayv1.GatewayClassSpec{
			ControllerName: GetControllerName(),
		},
	}

	t.Log("generating a list of matching controllers")
	matching := []gatewayv1.Gateway{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sanfrancisco",
				Namespace: "california",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: gatewayv1.ObjectName(gatewayClass.Name),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sandiego",
				Namespace: "california",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: gatewayv1.ObjectName(gatewayClass.Name),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "losangelos",
				Namespace: "california",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: gatewayv1.ObjectName(gatewayClass.Name),
			},
		},
	}

	t.Log("generating a list of non-matching controllers")
	nonmatching := []gatewayv1.Gateway{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hamburg",
				Namespace: "germany",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: gatewayv1.ObjectName("eu"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "paris",
				Namespace: "france",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: gatewayv1.ObjectName("eu"),
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

func TestIsGatewayControlledAndUnmanagedMode(t *testing.T) {
	var testControllerName gatewayv1.GatewayController = "acme.io/gateway-controller"

	testCases := []struct {
		name           string
		GatewayClass   *gatewayv1.GatewayClass
		expectedResult bool
	}{
		{
			name: "uncontrolled managed GatewayClass",
			GatewayClass: &gatewayv1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "uncontrolled-managed",
				},
				Spec: gatewayv1.GatewayClassSpec{
					ControllerName: testControllerName,
				},
			},
			expectedResult: false,
		},
		{
			name: "uncontrolled unmanaged GatewayClass",
			GatewayClass: &gatewayv1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "uncontrolled-unmanaged",
					Annotations: map[string]string{
						annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
					},
				},
				Spec: gatewayv1.GatewayClassSpec{
					ControllerName: testControllerName,
				},
			},
			expectedResult: false,
		},
		{
			name: "controlled managed GatewayClass",
			GatewayClass: &gatewayv1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "controlled-managed",
				},
				Spec: gatewayv1.GatewayClassSpec{
					ControllerName: GetControllerName(),
				},
			},
			expectedResult: false,
		},
		{
			name: "controlled unmanaged GatewayClass",
			GatewayClass: &gatewayv1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "controlled-unmanaged",
					Annotations: map[string]string{
						annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
					},
				},
				Spec: gatewayv1.GatewayClassSpec{
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
			assert.Equal(t, tc.expectedResult, isGatewayClassControlledAndUnmanaged(tc.GatewayClass))
		})
	}
}

func TestGetReferenceGrantConditionReason(t *testing.T) {
	testCases := []struct {
		name             string
		gatewayNamespace string
		certRef          gatewayv1.SecretObjectReference
		referenceGrants  []gatewayv1beta1.ReferenceGrant
		expectedReason   string
	}{
		{
			name:           "empty reference",
			certRef:        gatewayv1.SecretObjectReference{},
			expectedReason: string(gatewayv1.ListenerReasonResolvedRefs),
		},
		{
			name:             "no need for reference",
			gatewayNamespace: "test",
			certRef: gatewayv1.SecretObjectReference{
				Kind: util.StringToGatewayAPIKindPtr("Secret"),
				Name: "testSecret",
			},
			expectedReason: string(gatewayv1.ListenerReasonResolvedRefs),
		},
		{
			name:             "reference not granted - secret name not matching",
			gatewayNamespace: "test",
			certRef: gatewayv1.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: lo.ToPtr(Namespace("otherNamespace")),
			},
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "otherNamespace",
					},
					Spec: gatewayv1beta1.ReferenceGrantSpec{
						From: []gatewayv1beta1.ReferenceGrantFrom{
							{
								Group:     gatewayV1beta1Group,
								Kind:      "Gateway",
								Namespace: "test",
							},
						},
						To: []gatewayv1beta1.ReferenceGrantTo{
							{
								Group: "",
								Kind:  "Secret",
								Name:  lo.ToPtr(gatewayv1alpha2.ObjectName("anotherSecret")),
							},
						},
					},
				},
			},
			expectedReason: string(gatewayv1.ListenerReasonRefNotPermitted),
		},
		{
			name:             "reference not granted - no grants specified",
			gatewayNamespace: "test",
			certRef: gatewayv1.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: lo.ToPtr(Namespace("otherNamespace")),
			},
			expectedReason: string(gatewayv1.ListenerReasonRefNotPermitted),
		},
		{
			name:             "reference granted, secret name not specified",
			gatewayNamespace: "test",
			certRef: gatewayv1.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: lo.ToPtr(Namespace("otherNamespace")),
			},
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "otherNamespace",
					},
					Spec: gatewayv1beta1.ReferenceGrantSpec{
						From: []gatewayv1beta1.ReferenceGrantFrom{
							// useless entry, just to furtherly test the function
							{
								Group:     "otherGroup",
								Kind:      "otherKind",
								Namespace: "test",
							},
							// good entry
							{
								Group:     gatewayV1beta1Group,
								Kind:      "Gateway",
								Namespace: "test",
							},
						},
						To: []gatewayv1beta1.ReferenceGrantTo{
							{
								Group: "",
								Kind:  "Secret",
							},
						},
					},
				},
			},
			expectedReason: string(gatewayv1.ListenerReasonResolvedRefs),
		},
		{
			name:             "reference granted, secret name specified",
			gatewayNamespace: "test",
			certRef: gatewayv1.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: lo.ToPtr(Namespace("otherNamespace")),
			},
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "otherNamespace",
					},
					Spec: gatewayv1beta1.ReferenceGrantSpec{
						From: []gatewayv1beta1.ReferenceGrantFrom{
							{
								Group:     gatewayV1beta1Group,
								Kind:      "Gateway",
								Namespace: "test",
							},
						},
						To: []gatewayv1beta1.ReferenceGrantTo{
							{
								Group: "",
								Kind:  "Secret",
								Name:  lo.ToPtr(gatewayv1alpha2.ObjectName("testSecret")),
							},
						},
					},
				},
			},
			expectedReason: string(gatewayv1.ListenerReasonResolvedRefs),
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
