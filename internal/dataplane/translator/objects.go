package translator

import "sigs.k8s.io/controller-runtime/pkg/client"

// ObjectsCollector collects objects for later use.
// Its methods are safe to call with a nil receiver.
type ObjectsCollector struct {
	objects []client.Object
}

func NewObjectsCollector() *ObjectsCollector {
	return &ObjectsCollector{}
}

// Add adds an object to the collector. Noop if the receiver is nil.
func (p *ObjectsCollector) Add(obj client.Object) {
	if p == nil {
		return
	}

	p.objects = append(p.objects, obj)
}

// Pop returns the objects collected so far and resets the collector.
// Returns nil if the receiver is nil.
func (p *ObjectsCollector) Pop() []client.Object {
	if p == nil {
		return nil
	}

	objs := p.objects
	p.objects = nil
	return objs
}
