package configsecret

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// keyFor provides the string key that should be used to store the YAML contents of an object in the
// Kong configuration secret.
func KeyFor(obj runtime.Object, nsn types.NamespacedName) string {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return strings.Join([]string{gvk.Group, gvk.Version, gvk.Kind, nsn.Namespace, nsn.Name}, KeyDelimiter)
}
