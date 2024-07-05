package util

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// K8sObjectInfo describes a Kubernetes object.
type K8sObjectInfo struct {
	Name             string
	Namespace        string
	Annotations      map[string]string
	GroupVersionKind schema.GroupVersionKind
}

// FromK8sObject extracts information from a Kubernetes object.
// It performs a shallow copy of object annotations so any modifications after
// calling FromK8sObject will have an effect on the original object.
func FromK8sObject(obj client.Object) K8sObjectInfo {
	return K8sObjectInfo{
		Name:             obj.GetName(),
		Namespace:        obj.GetNamespace(),
		Annotations:      obj.GetAnnotations(),
		GroupVersionKind: obj.GetObjectKind().GroupVersionKind(),
	}
}
