package util

import (
	"maps"

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

func FromK8sObject(obj client.Object) K8sObjectInfo {
	ret := K8sObjectInfo{
		Name:        obj.GetName(),
		Namespace:   obj.GetNamespace(),
		Annotations: maps.Clone(obj.GetAnnotations()),
	}
	if gvk := obj.GetObjectKind().GroupVersionKind(); gvk.String() != "" {
		ret.GroupVersionKind = gvk
	}
	return ret
}
