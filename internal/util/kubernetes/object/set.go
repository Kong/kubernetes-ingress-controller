package object

import "sigs.k8s.io/controller-runtime/pkg/client"

type (
	gvk        string
	namespace  string
	name       string
	generation int64
)

type object struct {
	generation int64
	succeeded  bool
}

type ConfiguredStatus string

const (
	ConfiguredStatusSucceeded ConfiguredStatus = "Succeeded"
	ConfiguredStatusFailed    ConfiguredStatus = "Failed"
	ConfiguredStatusUnknown   ConfiguredStatus = "Unknown"
)

// Set is a de-duplicating list of Kubernetes objects used to test whether an
// object is a member of the set.
type Set struct {
	store map[gvk]map[namespace]map[name]object
}

func (s *Set) InsertSucceeded(objs ...client.Object) {
	s.insert(true, objs...)
}

func (s *Set) InsertFailed(objs ...client.Object) {
	s.insert(false, objs...)
}

func (s *Set) insert(succeeded bool, objs ...client.Object) {
	for _, obj := range objs {
		if s.store == nil {
			s.store = make(map[gvk]map[namespace]map[name]object)
		}

		objGVK := obj.GetObjectKind().GroupVersionKind().String()
		objNS := obj.GetNamespace()
		objName := obj.GetName()
		objGeneration := obj.GetGeneration()

		if s.store[gvk(objGVK)] == nil {
			s.store[gvk(objGVK)] = make(map[namespace]map[name]object)
		}

		if s.store[gvk(objGVK)][namespace(objNS)] == nil {
			s.store[gvk(objGVK)][namespace(objNS)] = make(map[name]object)
		}

		s.store[gvk(objGVK)][namespace(objNS)][name(objName)] = object{
			generation: objGeneration,
			succeeded:  succeeded,
		}
	}
}

func (s *Set) Get(obj client.Object) ConfiguredStatus {
	gvkStr := obj.GetObjectKind().GroupVersionKind().String()
	gvkMap, ok := s.store[gvk(gvkStr)]
	if !ok {
		return ConfiguredStatusUnknown
	}

	namespaceStr := obj.GetNamespace()
	namespaceMap, ok := gvkMap[namespace(namespaceStr)]
	if !ok {
		return ConfiguredStatusUnknown
	}

	objName := obj.GetName()
	storedObj, ok := namespaceMap[name(objName)]
	if !ok {
		return ConfiguredStatusUnknown
	}

	if storedObj.generation < obj.GetGeneration() {
		return ConfiguredStatusUnknown
	}

	if !storedObj.succeeded {
		return ConfiguredStatusFailed
	}

	return ConfiguredStatusSucceeded
}
