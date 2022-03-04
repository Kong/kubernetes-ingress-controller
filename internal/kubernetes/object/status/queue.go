package status

import (
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// ----------------------------------------------------------------------------
// Queue - Vars & Consts
// ----------------------------------------------------------------------------

// defaultBufferSize indicates the buffer size of the underlying channels that
// will be created for object kinds by default. This literally equates to the
// number of Kubernetes objects which can be in the queue at a single time.
const defaultBufferSize = 8192

// ----------------------------------------------------------------------------
// Queue - Public Types
// ----------------------------------------------------------------------------

// Queue provides a pub/sub queue with channels for individual Kubernetes
// objects, the purpose of which is to submit GenericEvents for those objects
// to trigger reconciliation when the object has been successfully configured in
// the dataplane so that its status can be updated (for instance, with IP
// address information in the case of Ingress resources).
type Queue struct {
	lock     sync.RWMutex
	channels map[string]chan event.GenericEvent
}

// NewQueue provides a new Queue object which can be used to
// publish status update events or subscribe to those events.
func NewQueue() *Queue {
	return &Queue{
		channels: make(map[string]chan event.GenericEvent),
	}
}

// ----------------------------------------------------------------------------
// Queue - Public Methods
// ----------------------------------------------------------------------------

// Publish emits a GenericEvent for the provided objects that indicates to
// subscribers that the status of that object needs to be updated.
func (q *Queue) Publish(obj client.Object) {
	ch := q.getChanForKind(obj.GetObjectKind().GroupVersionKind())
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
	return q.getChanForKind(gvk)
}

// ----------------------------------------------------------------------------
// Queue - Private Methods
// ----------------------------------------------------------------------------

func (q *Queue) getChanForKind(gvk schema.GroupVersionKind) chan event.GenericEvent {
	q.lock.Lock()
	defer q.lock.Unlock()
	ch, ok := q.channels[gvk.String()]
	if !ok { // if there's no channel built for this kind yet, make it
		ch = make(chan event.GenericEvent, defaultBufferSize)
		q.channels[gvk.String()] = ch
	}
	return ch
}
