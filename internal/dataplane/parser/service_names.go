package parser

import (
	"fmt"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getUniqueKongServiceName provides a deterministic unique string name to be
// used as the name of a Kong Service given details of the Kubernetes object
// which that Service is generated for.
func getUniqueKongServiceNameForObject(obj client.Object) (serviceName string) {
	kind := strings.ToLower(obj.GetObjectKind().GroupVersionKind().Kind)
	return fmt.Sprintf("%s.%s.%s", kind, obj.GetNamespace(), obj.GetName())
}
