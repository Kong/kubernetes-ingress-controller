package gateway

import (
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// handleUpdateError is a helper function to handle errors that occur when updating
// an object or its status.
// If the error is a conflict, it will return a requeue result.
// Otherwise, it will return the provided error.
func handleUpdateError(err error, log logr.Logger, obj client.Object) (ctrl.Result, error) {
	if apierrors.IsConflict(err) {
		debug(log, obj, "Conflict found when updating, retrying")
		return ctrl.Result{Requeue: true}, nil
	}
	return ctrl.Result{}, err
}
