package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func Test_readyConditionExistsForObservedGeneration(t *testing.T) {
	t.Log("checking ready condition for currently ready gateway")
	currentlyReadyGateway := &gatewayv1beta1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 1,
		},
		Status: gatewayv1beta1.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayv1beta1.GatewayConditionReady),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1beta1.GatewayReasonReady),
			}},
		},
	}
	assert.True(t, isGatewayReady(currentlyReadyGateway))

	t.Log("checking ready condition for previously ready gateway that has since been updated")
	previouslyReadyGateway := &gatewayv1beta1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 2,
		},
		Status: gatewayv1beta1.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayv1beta1.GatewayConditionReady),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1beta1.GatewayReasonReady),
			}},
		},
	}
	assert.False(t, isGatewayReady(previouslyReadyGateway))

	t.Log("checking ready condition for a gateway which has never been ready")
	neverBeenReadyGateway := &gatewayv1beta1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 10,
		},
		Status: gatewayv1beta1.GatewayStatus{},
	}
	assert.False(t, isGatewayReady(neverBeenReadyGateway))
}

func Test_setGatewayCondtion(t *testing.T) {
	testCases := []struct {
		name            string
		gw              *gatewayv1beta1.Gateway
		condition       metav1.Condition
		conditionLength int
	}{
		{
			name: "no_such_condition_should_append_one",
			gw:   &gatewayv1beta1.Gateway{},
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
			gw: &gatewayv1beta1.Gateway{
				Status: gatewayv1beta1.GatewayStatus{
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
			gw: &gatewayv1beta1.Gateway{
				Status: gatewayv1beta1.GatewayStatus{
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

func Test_isGatewayMarkedAsScheduled(t *testing.T) {
	t.Log("verifying scheduled check for gateway object which has been scheduled")
	scheduledGateway := &gatewayv1beta1.Gateway{
		ObjectMeta: metav1.ObjectMeta{Generation: 1},
		Status: gatewayv1beta1.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayv1beta1.GatewayConditionScheduled),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1beta1.GatewayReasonScheduled),
			}},
		},
	}
	assert.True(t, isGatewayScheduled(scheduledGateway))

	t.Log("verifying scheduled check for gateway object which has not been scheduled")
	unscheduledGateway := &gatewayv1beta1.Gateway{}
	assert.False(t, isGatewayScheduled(unscheduledGateway))
}

func Test_getRefFromPublishService(t *testing.T) {
	t.Log("verifying refs for valid publish services")
	valid := "california/sanfrancisco"
	nsn, err := getRefFromPublishService(valid)
	assert.NoError(t, err)
	assert.Equal(t, "california", nsn.Namespace)
	assert.Equal(t, "sanfrancisco", nsn.Name)

	t.Log("verifying failure conditions for invalid publish services")
	invalid := "california.sanfrancisco"
	_, err = getRefFromPublishService(invalid)
	assert.Error(t, err)
	moreInvalid := "california/sanfrancisco/losangelos"
	_, err = getRefFromPublishService(moreInvalid)
	assert.Error(t, err)
}

func Test_pruneStatusConditions(t *testing.T) {
	t.Log("verifying that a gateway with minimal status conditions is not pruned")
	gateway := &gatewayv1beta1.Gateway{}
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

func Test_reconcileGatewaysIfClassMatches(t *testing.T) {
	t.Log("generating a gatewayclass to test reconciliation filters")
	gatewayClass := &gatewayv1beta1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "us",
		},
		Spec: gatewayv1beta1.GatewayClassSpec{
			ControllerName: ControllerName,
		},
	}

	t.Log("generating a list of matching controllers")
	matching := []gatewayv1beta1.Gateway{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sanfrancisco",
				Namespace: "california",
			},
			Spec: gatewayv1beta1.GatewaySpec{
				GatewayClassName: gatewayv1beta1.ObjectName(gatewayClass.Name),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sandiego",
				Namespace: "california",
			},
			Spec: gatewayv1beta1.GatewaySpec{
				GatewayClassName: gatewayv1beta1.ObjectName(gatewayClass.Name),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "losangelos",
				Namespace: "california",
			},
			Spec: gatewayv1beta1.GatewaySpec{
				GatewayClassName: gatewayv1beta1.ObjectName(gatewayClass.Name),
			},
		},
	}

	t.Log("generating a list of non-matching controllers")
	nonmatching := []gatewayv1beta1.Gateway{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hamburg",
				Namespace: "germany",
			},
			Spec: gatewayv1beta1.GatewaySpec{
				GatewayClassName: gatewayv1beta1.ObjectName("eu"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "paris",
				Namespace: "france",
			},
			Spec: gatewayv1beta1.GatewaySpec{
				GatewayClassName: gatewayv1beta1.ObjectName("eu"),
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
			NamespacedName: types.NamespacedName{
				Name:      "sanfrancisco",
				Namespace: "california",
			},
		},
		{
			NamespacedName: types.NamespacedName{
				Name:      "sandiego",
				Namespace: "california",
			},
		},
		{
			NamespacedName: types.NamespacedName{
				Name:      "losangelos",
				Namespace: "california",
			},
		},
	}
	assert.Equal(t, expected, reconcileGatewaysIfClassMatches(gatewayClass, append(matching, nonmatching...)))
	assert.Equal(t, expected, reconcileGatewaysIfClassMatches(gatewayClass, matching))
}

func Test_isGatewayControlledAndUnmanagedMode(t *testing.T) {
	var testControllerName gatewayv1beta1.GatewayController = "acme.io/gateway-controller"

	testCases := []struct {
		name           string
		GatewayClass   *gatewayv1beta1.GatewayClass
		expectedResult bool
	}{
		{
			name: "uncontrolled managed GatewayClass",
			GatewayClass: &gatewayv1beta1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "uncontrolled-managed",
				},
				Spec: gatewayv1beta1.GatewayClassSpec{
					ControllerName: testControllerName,
				},
			},
			expectedResult: false,
		},
		{
			name: "uncontrolled unmanaged GatewayClass",
			GatewayClass: &gatewayv1beta1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "uncontrolled-unmanaged",
					Annotations: map[string]string{
						annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
					},
				},
				Spec: gatewayv1beta1.GatewayClassSpec{
					ControllerName: testControllerName,
				},
			},
			expectedResult: false,
		},
		{
			name: "controlled managed GatewayClass",
			GatewayClass: &gatewayv1beta1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "controlled-managed",
				},
				Spec: gatewayv1beta1.GatewayClassSpec{
					ControllerName: ControllerName,
				},
			},
			expectedResult: false,
		},
		{
			name: "controlled unmanaged GatewayClass",
			GatewayClass: &gatewayv1beta1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "controlled-unmanaged",
					Annotations: map[string]string{
						annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
					},
				},
				Spec: gatewayv1beta1.GatewayClassSpec{
					ControllerName: ControllerName,
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

func Test_areAllowedRoutesConsistentByProtocol(t *testing.T) {
	same := gatewayv1alpha2.NamespacesFromSame
	all := gatewayv1alpha2.NamespacesFromAll
	selector := gatewayv1alpha2.NamespacesFromSelector

	inputs := []struct {
		expected bool
		message  string
		l        []gatewayv1alpha2.Listener
	}{
		{
			expected: true,
			message:  "empty",
			l:        []gatewayv1alpha2.Listener{},
		},
		{
			expected: true,
			message:  "no intersect",
			l: []gatewayv1alpha2.Listener{
				{
					Protocol: gatewayv1alpha2.UDPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &same,
						},
					},
				},
				{
					Protocol: gatewayv1alpha2.TCPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &all,
						},
					},
				},
			},
		},
		{
			expected: true,
			message:  "same allowed for each listener with same protocol",
			l: []gatewayv1alpha2.Listener{
				{
					Protocol: gatewayv1alpha2.UDPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &same,
						},
					},
				},
				{
					Protocol: gatewayv1alpha2.UDPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &same,
						},
					},
				},
				{
					Protocol: gatewayv1alpha2.TCPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &all,
						},
					},
				},
			},
		},
		{
			expected: false,
			message:  "different allowed for listeners with same protocol",
			l: []gatewayv1alpha2.Listener{
				{
					Protocol: gatewayv1alpha2.UDPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &same,
						},
					},
				},
				{
					Protocol: gatewayv1alpha2.UDPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &all,
						},
					},
				},
			},
		},
		{
			expected: true,
			message:  "same selector",
			l: []gatewayv1alpha2.Listener{
				{
					Protocol: gatewayv1alpha2.UDPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &selector,
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"key": "value"},
							},
						},
					},
				},
				{
					Protocol: gatewayv1alpha2.UDPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &selector,
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"key": "value"},
							},
						},
					},
				},
			},
		},
		{
			expected: false,
			message:  "different selector",
			l: []gatewayv1alpha2.Listener{
				{
					Protocol: gatewayv1alpha2.UDPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &selector,
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"key": "value"},
							},
						},
					},
				},
				{
					Protocol: gatewayv1alpha2.UDPProtocolType,
					AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
						Namespaces: &gatewayv1alpha2.RouteNamespaces{
							From: &selector,
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"key": "notvalue"},
							},
						},
					},
				},
			},
		},
	}
	for _, input := range inputs {
		assert.Equal(t, input.expected, areAllowedRoutesConsistentByProtocol(input.l), input.message)
	}
}

func Test_getReferenceGrantConditionReason(t *testing.T) {
	testCases := []struct {
		name             string
		gatewayNamespace string
		certRef          gatewayv1beta1.SecretObjectReference
		referenceGrants  []gatewayv1alpha2.ReferenceGrant
		expectedReason   string
	}{
		{
			name:             "no need for reference",
			gatewayNamespace: "test",
			certRef: gatewayv1beta1.SecretObjectReference{
				Kind: util.StringToGatewayAPIKindPtr("Secret"),
				Name: "testSecret",
			},
			expectedReason: string(gatewayv1alpha2.ListenerReasonResolvedRefs),
		},
		{
			name:             "reference not granted",
			gatewayNamespace: "test",
			certRef: gatewayv1beta1.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: (*Namespace)(pointer.StringPtr("otherNamespace")),
			},
			referenceGrants: []gatewayv1alpha2.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "otherNamespace",
					},
					Spec: gatewayv1alpha2.ReferenceGrantSpec{
						From: []gatewayv1alpha2.ReferenceGrantFrom{
							{
								Group:     (gatewayv1alpha2.Group)(gatewayV1beta1Group),
								Kind:      "Gateway",
								Namespace: "test",
							},
						},
						To: []gatewayv1alpha2.ReferenceGrantTo{
							{
								Group: "",
								Kind:  "Secret",
								Name:  (*gatewayv1alpha2.ObjectName)(pointer.StringPtr("anotherSecret")),
							},
						},
					},
				},
			},
			expectedReason: string(gatewayv1alpha2.ListenerReasonInvalidCertificateRef),
		},
		{
			name:             "reference granted, secret name not specified",
			gatewayNamespace: "test",
			certRef: gatewayv1beta1.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: (*Namespace)(pointer.StringPtr("otherNamespace")),
			},
			referenceGrants: []gatewayv1alpha2.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "otherNamespace",
					},
					Spec: gatewayv1alpha2.ReferenceGrantSpec{
						From: []gatewayv1alpha2.ReferenceGrantFrom{
							// useless entry, just to furtherly test the function
							{
								Group:     "otherGroup",
								Kind:      "otherKind",
								Namespace: "test",
							},
							// good entry
							{
								Group:     (gatewayv1alpha2.Group)(gatewayV1beta1Group),
								Kind:      "Gateway",
								Namespace: "test",
							},
						},
						To: []gatewayv1alpha2.ReferenceGrantTo{
							{
								Group: "",
								Kind:  "Secret",
							},
						},
					},
				},
			},
			expectedReason: string(gatewayv1alpha2.ListenerReasonResolvedRefs),
		},
		{
			name:             "reference granted, secret name specified",
			gatewayNamespace: "test",
			certRef: gatewayv1beta1.SecretObjectReference{
				Kind:      util.StringToGatewayAPIKindPtr("Secret"),
				Name:      "testSecret",
				Namespace: (*Namespace)(pointer.StringPtr("otherNamespace")),
			},
			referenceGrants: []gatewayv1alpha2.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "otherNamespace",
					},
					Spec: gatewayv1alpha2.ReferenceGrantSpec{
						From: []gatewayv1alpha2.ReferenceGrantFrom{
							{
								Group:     (gatewayv1alpha2.Group)(gatewayV1beta1Group),
								Kind:      "Gateway",
								Namespace: "test",
							},
						},
						To: []gatewayv1alpha2.ReferenceGrantTo{
							{
								Group: "",
								Kind:  "Secret",
								Name:  (*gatewayv1alpha2.ObjectName)(pointer.StringPtr("testSecret")),
							},
						},
					},
				},
			},
			expectedReason: string(gatewayv1alpha2.ListenerReasonResolvedRefs),
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
