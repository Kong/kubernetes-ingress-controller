package gateway

import (
	"testing"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
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
	assert.True(t, readyConditionExistsForObservedGeneration(currentlyReadyGateway))

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
	assert.False(t, readyConditionExistsForObservedGeneration(previouslyReadyGateway))

	t.Log("checking ready condition for a gateway which has never been ready")
	neverBeenReadyGateway := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 10,
		},
		Status: gatewayv1alpha2.GatewayStatus{},
	}
	assert.False(t, readyConditionExistsForObservedGeneration(neverBeenReadyGateway))
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
	assert.True(t, isGatewayMarkedAsScheduled(scheduledGateway))

	t.Log("verifying scheduled check for gateway object which has not been scheduled")
	unscheduledGateway := &gatewayv1alpha2.Gateway{}
	assert.False(t, isGatewayMarkedAsScheduled(unscheduledGateway))
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
	assert.False(t, isGatewayControlledAndUnmanagedMode(controlledGatewayClass, gatewayv1alpha2.Gateway{}))
	assert.False(t, isGatewayControlledAndUnmanagedMode(uncontrolledGatewayClass, gatewayv1alpha2.Gateway{}))
	assert.False(t, isGatewayControlledAndUnmanagedMode(uncontrolledGatewayClass, unmanagedGateway))
	assert.True(t, isGatewayControlledAndUnmanagedMode(controlledGatewayClass, unmanagedGateway))
}

func Test_areAddressesEqual(t *testing.T) {
	ipAddressType := gatewayv1alpha2.IPAddressType
	hostnameAddressType := gatewayv1alpha2.HostnameAddressType
	namedAddressType := gatewayv1alpha2.NamedAddressType
	extraIPAddressType := gatewayv1alpha2.IPAddressType

	inputs := []struct {
		expected bool
		message  string
		l1       []gatewayv1alpha2.GatewayAddress
		l2       []gatewayv1alpha2.GatewayAddress
	}{
		{
			true,
			"two empty lists of addresses are equal",
			[]gatewayv1alpha2.GatewayAddress{},
			[]gatewayv1alpha2.GatewayAddress{},
		},
		{
			false,
			"empty list does not match with populated list",
			[]gatewayv1alpha2.GatewayAddress{},
			[]gatewayv1alpha2.GatewayAddress{{Value: "127.0.0.1"}},
		},
		{
			true,
			"two lists with partially populated structs match",
			[]gatewayv1alpha2.GatewayAddress{{Value: "127.0.0.1"}},
			[]gatewayv1alpha2.GatewayAddress{{Value: "127.0.0.1"}},
		},
		{
			false,
			"similar lists where only one has type do not match",
			[]gatewayv1alpha2.GatewayAddress{{Value: "127.0.0.1"}},
			[]gatewayv1alpha2.GatewayAddress{{
				Type:  &ipAddressType,
				Value: "127.0.0.1",
			}},
		},
		{
			false,
			"blantantly unmatching lists with the same count don't match",
			[]gatewayv1alpha2.GatewayAddress{{
				Type:  &namedAddressType,
				Value: "kong/proxy",
			}},
			[]gatewayv1alpha2.GatewayAddress{{
				Type:  &ipAddressType,
				Value: "127.0.0.1",
			}},
		},
		{
			true,
			"two identical lists of addresses with full attributes using hostname type match",
			[]gatewayv1alpha2.GatewayAddress{{
				Type:  &hostnameAddressType,
				Value: "konghq.com",
			}},
			[]gatewayv1alpha2.GatewayAddress{{
				Type:  &hostnameAddressType,
				Value: "konghq.com",
			}},
		},
		{
			true,
			"identical lists with more than one element match",
			[]gatewayv1alpha2.GatewayAddress{
				{
					Type:  &ipAddressType,
					Value: "192.168.1.1",
				},
				{
					Type:  &ipAddressType,
					Value: "192.168.1.2",
				},
			},
			[]gatewayv1alpha2.GatewayAddress{
				{
					Type:  &ipAddressType,
					Value: "192.168.1.1",
				},
				{
					Type:  &ipAddressType,
					Value: "192.168.1.2",
				},
			},
		},
		{
			true,
			"pointers with the same underlying value match even if the address is different",
			[]gatewayv1alpha2.GatewayAddress{
				{
					Type:  &ipAddressType,
					Value: "192.168.1.1",
				},
				{
					Type:  &ipAddressType,
					Value: "192.168.1.2",
				},
			},
			[]gatewayv1alpha2.GatewayAddress{
				{
					Type:  &ipAddressType,
					Value: "192.168.1.1",
				},
				{
					Type:  &extraIPAddressType,
					Value: "192.168.1.2",
				},
			},
		},
		{
			false,
			"two lists with the same contents but different ordering of those contents do not match",
			[]gatewayv1alpha2.GatewayAddress{
				{
					Type:  &ipAddressType,
					Value: "192.168.1.1",
				},
				{
					Type:  &ipAddressType,
					Value: "192.168.1.2",
				},
			},
			[]gatewayv1alpha2.GatewayAddress{
				{
					Type:  &ipAddressType,
					Value: "192.168.1.2",
				},
				{
					Type:  &ipAddressType,
					Value: "192.168.1.1",
				},
			},
		},
		{
			true,
			"large identical lists with various options match",
			[]gatewayv1alpha2.GatewayAddress{
				{Type: &ipAddressType, Value: "192.168.1.1"},
				{Type: &ipAddressType, Value: "192.168.1.2"},
				{Type: &ipAddressType, Value: "192.168.1.3"},
				{Type: &ipAddressType, Value: "192.168.1.4"},
				{Type: &ipAddressType, Value: "192.168.1.5"},
				{Type: &ipAddressType, Value: "192.168.1.6"},
				{Type: &ipAddressType, Value: "192.168.1.7"},
				{Type: &ipAddressType, Value: "192.168.1.8"},
				{Type: &ipAddressType, Value: "192.168.1.9"},
				{Type: &ipAddressType, Value: "192.168.1.10"},
				{Type: &hostnameAddressType, Value: "konghq.com"},
				{Type: &namedAddressType, Value: "kong/proxy"},
			},
			[]gatewayv1alpha2.GatewayAddress{
				{Type: &ipAddressType, Value: "192.168.1.1"},
				{Type: &ipAddressType, Value: "192.168.1.2"},
				{Type: &ipAddressType, Value: "192.168.1.3"},
				{Type: &ipAddressType, Value: "192.168.1.4"},
				{Type: &ipAddressType, Value: "192.168.1.5"},
				{Type: &ipAddressType, Value: "192.168.1.6"},
				{Type: &ipAddressType, Value: "192.168.1.7"},
				{Type: &ipAddressType, Value: "192.168.1.8"},
				{Type: &ipAddressType, Value: "192.168.1.9"},
				{Type: &ipAddressType, Value: "192.168.1.10"},
				{Type: &hostnameAddressType, Value: "konghq.com"},
				{Type: &namedAddressType, Value: "kong/proxy"},
			},
		},
		{
			false,
			"large lists with one difference don't match",
			[]gatewayv1alpha2.GatewayAddress{
				{Type: &ipAddressType, Value: "192.168.1.1"},
				{Type: &ipAddressType, Value: "192.168.1.2"},
				{Type: &ipAddressType, Value: "192.168.1.3"},
				{Type: &ipAddressType, Value: "192.168.1.4"},
				{Type: &ipAddressType, Value: "192.168.1.5"},
				{Type: &ipAddressType, Value: "192.168.1.6"},
				{Type: &ipAddressType, Value: "192.168.1.7"},
				{Type: &ipAddressType, Value: "192.168.1.8"},
				{Type: &ipAddressType, Value: "192.168.1.9"},
				{Type: &ipAddressType, Value: "192.168.1.10"},
				{Type: &hostnameAddressType, Value: "konghq.com"},
				{Type: &namedAddressType, Value: "kong/proxy"},
			},
			[]gatewayv1alpha2.GatewayAddress{
				{Type: &ipAddressType, Value: "192.168.1.1"},
				{Type: &ipAddressType, Value: "192.168.1.2"},
				{Type: &ipAddressType, Value: "192.168.1.3"},
				{Type: &ipAddressType, Value: "192.168.1.4"},
				{Type: &ipAddressType, Value: "192.168.1.5"},
				{Type: &ipAddressType, Value: "192.168.1.6"},
				{Type: &ipAddressType, Value: "192.168.1.7"},
				{Type: &ipAddressType, Value: "192.168.1.8"},
				{Type: &ipAddressType, Value: "192.168.1.9"},
				{Type: &ipAddressType, Value: "192.168.1.10"},
				{Type: &hostnameAddressType, Value: "konghq.com"},
				{Type: &namedAddressType, Value: "kong/admin"},
			},
		},
	}

	for _, input := range inputs {
		assert.Equal(t, input.expected, areAddressesEqual(input.l1, input.l2), input.message, input.l1, input.l2)
	}
}

func Test_areListenersEqual(t *testing.T) {
	genericHostname1 := gatewayv1alpha2.Hostname("konghq.com")
	genericHostname2 := gatewayv1alpha2.Hostname("docs.konghq.com")
	extraGenericHostname2 := gatewayv1alpha2.Hostname("docs.konghq.com")

	inputs := []struct {
		expected bool
		message  string
		l1       []gatewayv1alpha2.Listener
		l2       []gatewayv1alpha2.Listener
	}{
		{
			true,
			"two empty lists of addresses are equal",
			[]gatewayv1alpha2.Listener{},
			[]gatewayv1alpha2.Listener{},
		},
		{
			false,
			"empty list does not match with populated list",
			[]gatewayv1alpha2.Listener{},
			[]gatewayv1alpha2.Listener{{
				Name: gatewayv1alpha2.SectionName("http"),
			}},
		},
		{
			true,
			"two lists with partially populated structs match",
			[]gatewayv1alpha2.Listener{{
				Name: gatewayv1alpha2.SectionName("http"),
			}},
			[]gatewayv1alpha2.Listener{{
				Name: gatewayv1alpha2.SectionName("http"),
			}},
		},
		{
			false,
			"similar lists where only one has type do not match",
			[]gatewayv1alpha2.Listener{{
				Name: gatewayv1alpha2.SectionName("http"),
			}},
			[]gatewayv1alpha2.Listener{{
				Name:     gatewayv1alpha2.SectionName("http"),
				Hostname: &genericHostname1,
			}},
		},
		{
			false,
			"blantantly unmatching lists with the same count don't match",
			[]gatewayv1alpha2.Listener{{
				Name:     gatewayv1alpha2.SectionName("udp"),
				Hostname: &genericHostname2,
			}},
			[]gatewayv1alpha2.Listener{{
				Name:     gatewayv1alpha2.SectionName("http"),
				Hostname: &genericHostname1,
			}},
		},
		{
			true,
			"two identical lists of addresses with several attributes match",
			[]gatewayv1alpha2.Listener{{
				Name:     gatewayv1alpha2.SectionName("http"),
				Hostname: &genericHostname1,
				Port:     gatewayv1alpha2.PortNumber(80),
				Protocol: gatewayv1alpha2.HTTPProtocolType,
			}},
			[]gatewayv1alpha2.Listener{{
				Name:     gatewayv1alpha2.SectionName("http"),
				Hostname: &genericHostname1,
				Port:     gatewayv1alpha2.PortNumber(80),
				Protocol: gatewayv1alpha2.HTTPProtocolType,
			}},
		},
		{
			true,
			"identical lists with more than one element match",
			[]gatewayv1alpha2.Listener{
				{
					Name:     gatewayv1alpha2.SectionName("http"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(80),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("https"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(443),
					Protocol: gatewayv1alpha2.HTTPSProtocolType,
				},
			},
			[]gatewayv1alpha2.Listener{
				{
					Name:     gatewayv1alpha2.SectionName("http"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(80),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("https"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(443),
					Protocol: gatewayv1alpha2.HTTPSProtocolType,
				},
			},
		},
		{
			true,
			"pointers with the same underlying value match even if the address is different",
			[]gatewayv1alpha2.Listener{{
				Name:     gatewayv1alpha2.SectionName("http"),
				Hostname: &genericHostname2,
				Port:     gatewayv1alpha2.PortNumber(80),
				Protocol: gatewayv1alpha2.HTTPProtocolType,
			}},
			[]gatewayv1alpha2.Listener{{
				Name:     gatewayv1alpha2.SectionName("http"),
				Hostname: &extraGenericHostname2,
				Port:     gatewayv1alpha2.PortNumber(80),
				Protocol: gatewayv1alpha2.HTTPProtocolType,
			}},
		},
		{
			false,
			"two lists with the same contents but different ordering of those contents do not match",
			[]gatewayv1alpha2.Listener{
				{
					Name:     gatewayv1alpha2.SectionName("http"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(80),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("https"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(443),
					Protocol: gatewayv1alpha2.HTTPSProtocolType,
				},
			},
			[]gatewayv1alpha2.Listener{
				{
					Name:     gatewayv1alpha2.SectionName("https"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(443),
					Protocol: gatewayv1alpha2.HTTPSProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("http"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(80),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
			},
		},
		{
			true,
			"large identical lists with various options match",
			[]gatewayv1alpha2.Listener{
				{
					Name:     gatewayv1alpha2.SectionName("https"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(443),
					Protocol: gatewayv1alpha2.HTTPSProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("http"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(80),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("httpalt"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(8080),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("udp"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(53),
					Protocol: gatewayv1alpha2.UDPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("tcp"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(9000),
					Protocol: gatewayv1alpha2.TCPProtocolType,
				},
			},
			[]gatewayv1alpha2.Listener{
				{
					Name:     gatewayv1alpha2.SectionName("https"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(443),
					Protocol: gatewayv1alpha2.HTTPSProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("http"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(80),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("httpalt"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(8080),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("udp"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(53),
					Protocol: gatewayv1alpha2.UDPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("tcp"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(9000),
					Protocol: gatewayv1alpha2.TCPProtocolType,
				},
			},
		},
		{
			false,
			"large lists with one difference don't match",
			[]gatewayv1alpha2.Listener{
				{
					Name:     gatewayv1alpha2.SectionName("https"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(443),
					Protocol: gatewayv1alpha2.HTTPSProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("http"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(80),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("httpalt"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(8080),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("udp"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(53),
					Protocol: gatewayv1alpha2.UDPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("tcp"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(9000),
					Protocol: gatewayv1alpha2.TCPProtocolType,
				},
			},
			[]gatewayv1alpha2.Listener{
				{
					Name:     gatewayv1alpha2.SectionName("https"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(443),
					Protocol: gatewayv1alpha2.HTTPSProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("http"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(80),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("httpalt"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(8080),
					Protocol: gatewayv1alpha2.HTTPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("udp"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(53),
					Protocol: gatewayv1alpha2.UDPProtocolType,
				},
				{
					Name:     gatewayv1alpha2.SectionName("udp"),
					Hostname: &genericHostname1,
					Port:     gatewayv1alpha2.PortNumber(9000),
					Protocol: gatewayv1alpha2.UDPProtocolType,
				},
			},
		},
	}

	for _, input := range inputs {
		assert.Equal(t, input.expected, areListenersEqual(input.l1, input.l2), input.message, input.l1, input.l2)
	}
}
