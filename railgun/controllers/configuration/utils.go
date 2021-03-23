package configuration

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/railgun/controllers"
)

// -----------------------------------------------------------------------------
// Utils - Controller Helper Functions
// -----------------------------------------------------------------------------

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

// IsAPIAvailable is a hack to short circuit controllers for APIs which aren't available on the cluster,
// enabling us to keep separate logic and logging for some legacy API versions.
func IsAPIAvailable(mgr ctrl.Manager, obj client.Object) (bool, error) {
	if err := mgr.GetAPIReader().Get(context.Background(), client.ObjectKey{Namespace: controllers.DefaultNamespace, Name: "non-existent"}, obj); err != nil {
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
