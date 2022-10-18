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
	store map[gvk]map[namespace]map[name]int64
}

func (s *Set) Insert(objs ...client.Object) {
	for _, obj := range objs {
		if s.store == nil {
			s.store = make(map[gvk]map[namespace]map[name]int64)
		}

		objGeneration := obj.GetGeneration()
		objGVK := obj.GetObjectKind().GroupVersionKind().String()
		objNS := obj.GetNamespace()
		objName := obj.GetName()

		if s.store[gvk(objGVK)] == nil {
			s.store[gvk(objGVK)] = make(map[namespace]map[name]int64)
		}

		if s.store[gvk(objGVK)][namespace(objNS)] == nil {
			s.store[gvk(objGVK)][namespace(objNS)] = make(map[name]int64)
		}

		s.store[gvk(objGVK)][namespace(objNS)][name(objName)] = objGeneration
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
	observedGeneration, ok := namespaceMap[name(objName)]
	if !ok {
		return false
	}

	// don't consider an object to be configured in case it doesn't match the exact generation we are looking at
	expectedGeneration := obj.GetGeneration()
	return observedGeneration == expectedGeneration
}
