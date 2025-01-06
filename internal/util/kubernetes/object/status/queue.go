package status

import (
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// DefaultBufferSize indicates the buffer size of the underlying channels that
// will be created for object kinds by default. This literally equates to the
// number of Kubernetes objects which can be in the queue at a single time.
const DefaultBufferSize = 8192

// Queue provides a pub/sub queue with channels for individual Kubernetes
// objects, the purpose of which is to submit GenericEvents for those objects
// to trigger reconciliation when the object has been successfully configured in
// the dataplane so that its status can be updated (for instance, with IP
// address information in the case of Ingress resources).
type Queue struct {
	// subscriptionBufferSize indicates the buffer size of the underlying channels
	// that will be created for object kinds' subscriptions.
	subscriptionBufferSize int

	// subscriptions indexed by the string representation of the object GVK.
	subscriptions map[string][]chan event.GenericEvent

	// lock protects the subscriptions map.
	lock sync.RWMutex
}

// QueueOption provides a functional option for configuring a Queue object.
type QueueOption func(*Queue)

// WithBufferSize sets the buffer size of the underlying channels that will be
// created for object kinds.
func WithBufferSize(size int) QueueOption {
	return func(q *Queue) {
		q.subscriptionBufferSize = size
	}
}

// NewQueue provides a new Queue object which can be used to
// publish status update events or subscribe to those events.
func NewQueue(opts ...QueueOption) *Queue {
	q := &Queue{
		subscriptionBufferSize: DefaultBufferSize,
		subscriptions:          make(map[string][]chan event.GenericEvent),
	}
	for _, opt := range opts {
		opt(q)
	}

	return q
}

// Publish emits a GenericEvent for the provided objects that indicates to
// subscribers that the status of that object needs to be updated.
// It's a no-op if there are no subscriptions for the object kind.
func (q *Queue) Publish(obj client.Object) {
	// Publish the event to all subscribers.
	for _, ch := range q.getSubscriptionsForKind(obj.GetObjectKind().GroupVersionKind()) {
		ch <- event.GenericEvent{Object: obj}
	}
}

// Subscribe provides a consumer channel where generic events for the provided
// object kind will be published when the status for the object needs to be
// updated.
//
// Please note every call to Subscribe will create a new subscription channel.
func (q *Queue) Subscribe(gvk schema.GroupVersionKind) chan event.GenericEvent {
	return q.createSubscriptionForKind(gvk)
}

// createSubscriptionForKind creates a new subscription channel for the provided object GVK and returns it.
func (q *Queue) createSubscriptionForKind(gvk schema.GroupVersionKind) chan event.GenericEvent {
	newSubscriptionCh := make(chan event.GenericEvent, q.subscriptionBufferSize)
	q.lock.Lock()
	defer q.lock.Unlock()
	q.subscriptions[gvk.String()] = append(q.subscriptions[gvk.String()], newSubscriptionCh)
	return newSubscriptionCh
}

// getSubscriptionsForKind returns all the subscription channels for the provided object GVK.
// It may return nil if there are no subscriptions for the object kind.
func (q *Queue) getSubscriptionsForKind(gvk schema.GroupVersionKind) []chan event.GenericEvent {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.subscriptions[gvk.String()]
}
