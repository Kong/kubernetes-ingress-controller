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

func (s Set) Insert(objs ...client.Object) Set {
	for _, obj := range objs {
		gvk, ok := s.store[gvk(obj.GetObjectKind().GroupVersionKind().String())]
		if !ok {
			gvk = make(map[namespace]map[name]struct{})
		}

		namespace, ok := gvk[namespace(obj.GetNamespace())]
		if !ok {
			namespace = make(map[name]struct{})
		}

		namespace[name(obj.GetName())] = struct{}{}
	}
	return s
}

func (s Set) Delete(objs ...client.Object) Set {
	for _, obj := range objs {
		gvk, ok := s.store[gvk(obj.GetObjectKind().GroupVersionKind().String())]
		if !ok {
			continue
		}

		namespace, ok := gvk[namespace(obj.GetNamespace())]
		if !ok {
			continue
		}

		delete(namespace, name(obj.GetName()))
	}
	return s
}

func (s Set) Has(obj client.Object) bool {
	gvk, ok := s.store[gvk(obj.GetObjectKind().GroupVersionKind().String())]
	if !ok {
		return false
	}

	namespace, ok := gvk[namespace(obj.GetNamespace())]
	if !ok {
		return false
	}

	_, ok = namespace[name(obj.GetName())]
	return ok
}
