package store

import "reflect"

// namespacedKeyFunc returns a key for a namespaced object.
func namespacedKeyFunc(obj interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	name := v.FieldByName("Name")
	namespace := v.FieldByName("Namespace")
	return namespace.String() + "/" + name.String(), nil
}

// clusterWideKeyFunc returns a key for a cluster-wide object.
func clusterWideKeyFunc(obj interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	return v.FieldByName("Name").String(), nil
}
