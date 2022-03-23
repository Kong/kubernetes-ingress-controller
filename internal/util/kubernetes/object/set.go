package object

import "sigs.k8s.io/controller-runtime/pkg/client"

type (
	gvk       string
	namespace string
	name      string
)

// Set is a de-duplicating list of Kubernetes objects used to test whether an
// object is a member of the set.
type Set struct {
	store map[gvk]map[namespace]map[name]struct{}
}

func (s *Set) Insert(objs ...client.Object) {
	for _, obj := range objs {
		if s.store == nil {
			s.store = make(map[gvk]map[namespace]map[name]struct{})
		}

		objGVK := obj.GetObjectKind().GroupVersionKind().String()
		objNS := obj.GetNamespace()
		objName := obj.GetName()

		if s.store[gvk(objGVK)] == nil {
			s.store[gvk(objGVK)] = make(map[namespace]map[name]struct{})
		}

		if s.store[gvk(objGVK)][namespace(objNS)] == nil {
			s.store[gvk(objGVK)][namespace(objNS)] = make(map[name]struct{})
		}

		s.store[gvk(objGVK)][namespace(objNS)][name(objName)] = struct{}{}
	}
}

func (s *Set) Has(obj client.Object) bool {
	gvkStr := obj.GetObjectKind().GroupVersionKind().String()
	gvkMap, ok := s.store[gvk(gvkStr)]
	if !ok {
		return false
	}

	namespaceStr := obj.GetNamespace()
	namespaceMap, ok := gvkMap[namespace(namespaceStr)]
	if !ok {
		return false
	}

	objName := obj.GetName()
	_, ok = namespaceMap[name(objName)]
	return ok
}
