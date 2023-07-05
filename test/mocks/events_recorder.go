package mocks

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
)

// EventRecorder is a mock implementation of the k8s.io/client-go/tools/record.EventRecorder interface.
type EventRecorder struct {
	events []string
	l      sync.RWMutex
}

func NewEventRecorder() *EventRecorder {
	return &EventRecorder{}
}

func (r *EventRecorder) Event(_ runtime.Object, eventtype, reason, message string) {
	r.writeEvent(eventtype, reason, "%s", message)
}

func (r *EventRecorder) Eventf(_ runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	r.writeEvent(eventtype, reason, messageFmt, args...)
}

func (r *EventRecorder) AnnotatedEventf(_ runtime.Object, _ map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	r.writeEvent(eventtype, reason, messageFmt, args...)
}

func (r *EventRecorder) Events() []string {
	r.l.RLock()
	defer r.l.RUnlock()
	copied := make([]string, len(r.events))
	copy(copied, r.events)
	return copied
}

func (r *EventRecorder) writeEvent(eventtype, reason, messageFmt string, args ...interface{}) {
	r.l.Lock()
	defer r.l.Unlock()
	r.events = append(r.events, fmt.Sprintf(eventtype+" "+reason+" "+messageFmt, args...))
}
