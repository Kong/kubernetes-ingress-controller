package inputs

import "sigs.k8s.io/controller-runtime/pkg/client"

// hasFinalizer is a helper function to check whether a client.Object
// already has a specific finalizer set.
func hasFinalizer(obj client.Object, finalizer string) bool {
	hasFinalizer := false
	for _, finalizer := range obj.GetFinalizers() {
		if finalizer == finalizer {
			hasFinalizer = true
		}
	}
	return hasFinalizer
}
