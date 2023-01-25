package object

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	gvk string
)

type objectConfigurationStatus struct {
	generation int64
	succeeded  bool
}

type ConfigurationStatus string

const (
	ConfigurationStatusSucceeded ConfigurationStatus = "Succeeded"
	ConfigurationStatusFailed    ConfigurationStatus = "Failed"
	ConfigurationStatusUnknown   ConfigurationStatus = "Unknown"
)

// ConfigurationStatusSet is a de-duplicate set to store the configure status
// (succeeded, failed, unknown) of kubernetes objects.
type ConfigurationStatusSet struct {
	store map[gvk]map[types.NamespacedName]objectConfigurationStatus
}

func NewConfigurationStatusSet() *ConfigurationStatusSet {
	return &ConfigurationStatusSet{
		store: map[gvk]map[types.NamespacedName]objectConfigurationStatus{},
	}
}

func (s *ConfigurationStatusSet) Insert(obj client.Object, succeeded bool) {
	if s.store == nil {
		s.store = make(map[gvk]map[types.NamespacedName]objectConfigurationStatus)
	}

	objGVK := gvk(obj.GetObjectKind().GroupVersionKind().String())
	nsName := types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
	if s.store[objGVK] == nil {
		s.store[objGVK] = make(map[types.NamespacedName]objectConfigurationStatus)
	}
	s.store[objGVK][nsName] = objectConfigurationStatus{
		generation: obj.GetGeneration(),
		succeeded:  succeeded,
	}
}

func (s *ConfigurationStatusSet) Get(obj client.Object) ConfigurationStatus {
	if s.store == nil {
		return ConfigurationStatusUnknown
	}

	objGVK := gvk(obj.GetObjectKind().GroupVersionKind().String())
	nsName := types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}

	gvkMap, ok := s.store[objGVK]
	if !ok {
		return ConfigurationStatusUnknown
	}

	status, ok := gvkMap[nsName]
	if !ok {
		return ConfigurationStatusUnknown
	}

	// if the object generation is newer than the generation of current configuration,
	// the latest specification of the object may still not configured on Kong gateway, so "Unknown" is returned.
	if status.generation < obj.GetGeneration() {
		return ConfigurationStatusUnknown
	}

	if !status.succeeded {
		return ConfigurationStatusFailed
	}

	return ConfigurationStatusSucceeded
}
