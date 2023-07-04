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
	subscriptions map[string]chan event.GenericEvent

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
		subscriptions:          make(map[string]chan event.GenericEvent),
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
	ch, ok := q.getSubscriptionForKind(obj.GetObjectKind().GroupVersionKind())
	if !ok {
		// There's no subscriber for this object kind - nothing to do.
		return
	}
	ch <- event.GenericEvent{Object: obj}
}

// Subscribe provides a consumer channel where generic events for the provided
// object kind will be published when the status for the object needs to be
// updated.
//
// Note that more than one consumer can subscribe to a channel for a particular
// object kind, but that this only represents a single channel: events will not
// be duplicated and each subscriber will receive events on a first come first
// serve basis.
func (q *Queue) Subscribe(gvk schema.GroupVersionKind) chan event.GenericEvent {
	return q.getOrCreateSubscriptionForKind(gvk)
}

// getOrCreateSubscriptionForKind returns the subscription channel for the provided object GVK.
// If the channel does not exist, it will be created.
func (q *Queue) getOrCreateSubscriptionForKind(gvk schema.GroupVersionKind) chan event.GenericEvent {
	q.lock.Lock()
	defer q.lock.Unlock()
	ch, ok := q.subscriptions[gvk.String()]
	if !ok {
		// If there's no channel built for this kind yet, make it.
		ch = make(chan event.GenericEvent, q.subscriptionBufferSize)
		q.subscriptions[gvk.String()] = ch
	}
	return ch
}

// getSubscriptionForKind returns the subscription channel for the provided object GVK.
// The second return value indicates whether the channel exists.
func (q *Queue) getSubscriptionForKind(gvk schema.GroupVersionKind) (chan event.GenericEvent, bool) {
	q.lock.RLock()
	defer q.lock.RUnlock()
	ch, ok := q.subscriptions[gvk.String()]
	return ch, ok
}
