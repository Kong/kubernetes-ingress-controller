package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/event"

	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func TestQueue(t *testing.T) {
	t.Log("creating a status queue")
	q := NewQueue()
	assert.Len(t, q.channels, 0, "no channels should be created by default")

	t.Log("generating Kubernetes objects to emit events for the queue")
	ing1 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      "ingress-test-1",
		},
	}
	ing2 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      "ingress-test-1",
		},
	}
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      "tcpingress-test-1",
		},
	}
	udp := &kongv1beta1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      "udpingress-test-1",
		},
	}

	t.Log("initializing kubernetes objects (this would normally be done by api client)")
	ing1.SetGroupVersionKind(ingGVK)
	ing2.SetGroupVersionKind(ingGVK)
	tcp.SetGroupVersionKind(tcpGVK)
	udp.SetGroupVersionKind(udpGVK)

	t.Log("verifying that events can be subscribed to for new object kinds")
	ingCH := q.Subscribe(ing1.GroupVersionKind())
	assert.Len(t, q.channels, 1, "internally a single channel should be created for the object kind")
	t.Logf("%+v", q.channels)
	assert.Len(t, ingCH, 0, "the underlying channel should be empty")

	t.Log("verifying an object can have an event published for it")
	q.Publish(ing1)
	assert.Len(t, q.channels, 1, "a channel was already created for the consumer: no more should be created")
	assert.Len(t, ingCH, 1, "the underlying channel should now contain one event")

	t.Log("verifying a published event can be consumed by the consumer")
	assert.Equal(t, event.GenericEvent{Object: ing1}, <-ingCH)
	assert.Len(t, ingCH, 0, "the event should be consumed")

	t.Log("verifying publishing different named objects for kinds that have already been seen")
	q.Publish(ing2)
	assert.Len(t, q.channels, 1, "a channel was already created for the object kind: no more should be created")
	assert.Len(t, ingCH, 1, "the underlying channel should now contain one event")

	t.Log("verifying that objects of new kinds can be published into the queue")
	q.Publish(tcp)
	q.Publish(udp)
	tcpCH := q.Subscribe(tcp.GroupVersionKind())
	udpCH := q.Subscribe(udp.GroupVersionKind())
	assert.Len(t, q.channels, 3, "2 new channels should have been created for the two new object kinds")
	assert.Len(t, ingCH, 1, "the underlying channel should contain 1 event")
	assert.Len(t, tcpCH, 1, "the underlying channel should contain 1 event")
	assert.Len(t, udpCH, 1, "the underlying channel should contain 1 event")

	t.Log("verifying that multiple events can be submitted for the same object")
	q.Publish(ing1)
	q.Publish(ing2)
	q.Publish(ing2)
	q.Publish(tcp)
	q.Publish(tcp)
	q.Publish(tcp)
	q.Publish(tcp)
	q.Publish(udp)
	q.Publish(udp)
	q.Publish(udp)
	q.Publish(udp)
	q.Publish(udp)
	assert.Len(t, q.channels, 3)
	assert.Len(t, ingCH, 4)
	assert.Len(t, tcpCH, 5)
	assert.Len(t, udpCH, 6)

	t.Log("verifying that all objects can be consumed and the queue can be drained")
	for i := 0; i < 4; i++ {
		assert.Equal(t, ingGVK, (<-ingCH).Object.GetObjectKind().GroupVersionKind())
	}
	for i := 0; i < 5; i++ {
		assert.Equal(t, event.GenericEvent{Object: tcp}, <-tcpCH)
	}
	for i := 0; i < 6; i++ {
		assert.Equal(t, event.GenericEvent{Object: udp}, <-udpCH)
	}
	assert.Len(t, q.channels, 3)
	assert.Len(t, ingCH, 0)
	assert.Len(t, tcpCH, 0)
	assert.Len(t, udpCH, 0)
}

// the GVKs for objects need to be initialized manually in the unit testing
// case as this would normally be done by the API and client for real objects.
var (
	ingGVK = schema.GroupVersionKind{
		Group:   "networking.k8s.io",
		Version: "v1",
		Kind:    "Ingress",
	}
	tcpGVK = schema.GroupVersionKind{
		Group:   "configuration.konghq.com",
		Version: "v1beta1",
		Kind:    "TCPIngress",
	}
	udpGVK = schema.GroupVersionKind{
		Group:   "configuration.konghq.com",
		Version: "v1beta1",
		Kind:    "UDPIngress",
	}
)
