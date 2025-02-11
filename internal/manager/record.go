package manager

import (
	"github.com/samber/lo"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/consts"
)

// newEventRecorderForInstance creates a new event recorder that will annotate all events with the given instance ID.
// It implements standard the record.EventRecorder interface.
func newEventRecorderForInstance(instanceID InstanceID, eventRecorder record.EventRecorder) *eventRecorderWithID {
	return &eventRecorderWithID{
		idAnnotation: map[string]string{
			consts.InstanceIDAnnotationKey: instanceID.String(),
		},
		eventRecorder: eventRecorder,
	}
}

var _ record.EventRecorder = &eventRecorderWithID{}

type eventRecorderWithID struct {
	idAnnotation  map[string]string
	eventRecorder record.EventRecorder
}

func (e *eventRecorderWithID) AnnotatedEventf(
	object runtime.Object, annotations map[string]string, eventtype string, reason string, messageFmt string, args ...interface{},
) {
	e.eventRecorder.AnnotatedEventf(object, lo.Assign(e.idAnnotation, annotations), eventtype, reason, messageFmt, args...)
}

func (e *eventRecorderWithID) Event(object runtime.Object, eventtype string, reason string, message string) {
	e.eventRecorder.AnnotatedEventf(object, e.idAnnotation, eventtype, reason, message)
}

func (e *eventRecorderWithID) Eventf(object runtime.Object, eventtype string, reason string, messageFmt string, args ...interface{}) {
	e.eventRecorder.AnnotatedEventf(object, e.idAnnotation, eventtype, reason, messageFmt, args...)
}
