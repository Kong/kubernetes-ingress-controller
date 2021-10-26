package util

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// K8sObjectInfo describes a Kubernetes object.
type K8sObjectInfo struct {
	Name        string
	Namespace   string
	Annotations map[string]string
}

func deepCopy(m map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range m {
		result[k] = v
	}
	return result
}

func FromK8sObject(obj metav1.Object) K8sObjectInfo {
	return K8sObjectInfo{
		Name:        obj.GetName(),
		Namespace:   obj.GetNamespace(),
		Annotations: deepCopy(obj.GetAnnotations()),
	}
}
