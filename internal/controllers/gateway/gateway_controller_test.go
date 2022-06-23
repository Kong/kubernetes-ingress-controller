package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

func Test_readyConditionExistsForObservedGeneration(t *testing.T) {
	t.Log("checking ready condition for currently ready gateway")
	currentlyReadyGateway := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 1,
		},
		Status: gatewayv1alpha2.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayv1alpha2.GatewayConditionReady),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1alpha2.GatewayReasonReady),
			}},
		},
	}
	assert.True(t, isGatewayReady(currentlyReadyGateway))

	t.Log("checking ready condition for previously ready gateway that has since been updated")
	previouslyReadyGateway := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 2,
		},
		Status: gatewayv1alpha2.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayv1alpha2.GatewayConditionReady),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1alpha2.GatewayReasonReady),
			}},
		},
	}
	assert.False(t, isGatewayReady(previouslyReadyGateway))

	t.Log("checking ready condition for a gateway which has never been ready")
	neverBeenReadyGateway := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 10,
		},
		Status: gatewayv1alpha2.GatewayStatus{},
	}
	assert.False(t, isGatewayReady(neverBeenReadyGateway))
}

func Test_isGatewayMarkedAsScheduled(t *testing.T) {
	t.Log("verifying scheduled check for gateway object which has been scheduled")
	scheduledGateway := &gatewayv1alpha2.Gateway{
		Status: gatewayv1alpha2.GatewayStatus{
			Conditions: []metav1.Condition{{
				Type:               string(gatewayv1alpha2.GatewayConditionScheduled),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1alpha2.GatewayReasonScheduled),
			}},
		},
	}
	assert.True(t, isGatewayScheduled(scheduledGateway))

	t.Log("verifying scheduled check for gateway object which has not been scheduled")
	unscheduledGateway := &gatewayv1alpha2.Gateway{}
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
	gateway := &gatewayv1alpha2.Gateway{}
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
	gatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "us",
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: ControllerName,
		},
	}

	t.Log("generating a list of matching controllers")
	matching := []gatewayv1alpha2.Gateway{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sanfrancisco",
				Namespace: "california",
			},
			Spec: gatewayv1alpha2.GatewaySpec{
				GatewayClassName: gatewayv1alpha2.ObjectName(gatewayClass.Name),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sandiego",
				Namespace: "california",
			},
			Spec: gatewayv1alpha2.GatewaySpec{
				GatewayClassName: gatewayv1alpha2.ObjectName(gatewayClass.Name),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "losangelos",
				Namespace: "california",
			},
			Spec: gatewayv1alpha2.GatewaySpec{
				GatewayClassName: gatewayv1alpha2.ObjectName(gatewayClass.Name),
			},
		},
	}

	t.Log("generating a list of non-matching controllers")
	nonmatching := []gatewayv1alpha2.Gateway{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hamburg",
				Namespace: "germany",
			},
			Spec: gatewayv1alpha2.GatewaySpec{
				GatewayClassName: gatewayv1alpha2.ObjectName("eu"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "paris",
				Namespace: "france",
			},
			Spec: gatewayv1alpha2.GatewaySpec{
				GatewayClassName: gatewayv1alpha2.ObjectName("eu"),
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
	t.Log("generating a gatewayclass controlled by this controller implementation")
	controlledGatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "us",
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: ControllerName,
		},
	}

	t.Log("generating a gatewayclass not controlled by this implementation")
	uncontrolledGatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "eu",
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: "acme.io/gateway-controller",
		},
	}

	t.Log("creating an unmanaged mode enabled gateway")
	unmanagedGateway := gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.GatewayUnmanagedAnnotation: "true",
			},
		},
	}

	t.Log("verifying the results for several gateways")
	assert.False(t, isGatewayInClassAndUnmanaged(controlledGatewayClass, gatewayv1alpha2.Gateway{}))
	assert.False(t, isGatewayInClassAndUnmanaged(uncontrolledGatewayClass, gatewayv1alpha2.Gateway{}))
	assert.False(t, isGatewayInClassAndUnmanaged(uncontrolledGatewayClass, unmanagedGateway))
	assert.True(t, isGatewayInClassAndUnmanaged(controlledGatewayClass, unmanagedGateway))
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
