package inputs

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

// isAPIAvailable is a hack to short circuit controllers for APIs which aren't available on the cluster,
// enabling us to keep separate logic and logging for some legacy API versions.
func isAPIAvailable(mgr ctrl.Manager, obj client.Object) (bool, error) {
	if err := mgr.GetAPIReader().Get(context.Background(), client.ObjectKey{Namespace: "kube-system", Name: "non-existent"}, obj); err != nil {
		if strings.Contains(err.Error(), "no matches for kind") {
			return false, nil
		}
		if errors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	}

	return true, nil
}
