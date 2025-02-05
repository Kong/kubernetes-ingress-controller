package mocks

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/scheme"
)

// EventRecorder is a mock implementation of the k8s.io/client-go/tools/record.EventRecorder interface.
type EventRecorder struct {
	events []string
	l      sync.RWMutex
}

func NewEventRecorder() *EventRecorder {
	return &EventRecorder{}
}

func (r *EventRecorder) Event(o runtime.Object, eventtype, reason, message string) {
	r.writeEvent(o, eventtype, reason, "%s", message)
}

func (r *EventRecorder) Eventf(o runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	r.writeEvent(o, eventtype, reason, messageFmt, args...)
}

func (r *EventRecorder) AnnotatedEventf(o runtime.Object, _ map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	r.writeEvent(o, eventtype, reason, messageFmt, args...)
}

func (r *EventRecorder) Events() []string {
	r.l.RLock()
	defer r.l.RUnlock()
	copied := make([]string, len(r.events))
	copy(copied, r.events)
	return copied
}

func (r *EventRecorder) writeEvent(o runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	r.l.Lock()
	defer r.l.Unlock()

	_ = util.PopulateTypeMeta(o, scheme.Get())
	fmtString := fmt.Sprintf("%s: %s %s %s", o.GetObjectKind().GroupVersionKind().Kind, eventtype, reason, messageFmt)
	r.events = append(r.events, fmt.Sprintf(fmtString, args...))
}
