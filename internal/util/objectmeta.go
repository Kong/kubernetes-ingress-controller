package util

import (
	"maps"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (k K8sObjectInfo) GetAnnotations() map[string]string {
	return k.Annotations
}

func (k K8sObjectInfo) GetNamespace() string {
	return k.Namespace
}

// FromK8sObject extracts information from a Kubernetes object.
// It performs a shallow copy of object annotations so any modifications after
// calling FromK8sObject will have an effect on the original object.
func FromK8sObject(obj client.Object) K8sObjectInfo {
	return K8sObjectInfo{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
		// We return a copy of annotations map here because translator functions may modify annotations
		// and that change would then be stored in store which is not desired.
		Annotations:      maps.Clone(obj.GetAnnotations()),
		GroupVersionKind: obj.GetObjectKind().GroupVersionKind(),
	}
}

// TypeMetaFromGVK returns typemeta from groupversionkind of a k8s object.
func TypeMetaFromGVK(gvk schema.GroupVersionKind) metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
	}
}
