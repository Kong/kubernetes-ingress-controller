package ctrlutils

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CleanupFinalizer ensures that a deleted resource is no longer present in the object cache.
func CleanupFinalizer(ctx context.Context, c client.Client, log logr.Logger, nsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
	if HasFinalizer(obj, KongIngressFinalizer) {
		log.Info("kong ingress finalizer needs to be removed from a resource which is deleting", "ingress", obj.GetName(), "finalizer", KongIngressFinalizer)
		finalizers := []string{}
		for _, finalizer := range obj.GetFinalizers() {
			if finalizer != KongIngressFinalizer {
				finalizers = append(finalizers, finalizer)
			}
		}
		obj.SetFinalizers(finalizers)
		if err := c.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("the kong ingress finalizer was removed from an a resource which is deleting", "ingress", obj.GetName(), "finalizer", KongIngressFinalizer)
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// HasFinalizer is a helper function to check whether a client.Object
// already has a specific finalizer set.
func HasFinalizer(obj client.Object, finalizer string) bool {
	hasFinalizer := false
	for _, foundFinalizer := range obj.GetFinalizers() {
		if foundFinalizer == finalizer {
			hasFinalizer = true
		}
	}
	return hasFinalizer
}

// IsAPIAvailable is a hack to short circuit controllers for APIs which aren't available on the cluster,
// enabling us to keep separate logic and logging for some legacy API versions.
func IsAPIAvailable(mgr ctrl.Manager, obj client.Object) (bool, error) {
	if err := mgr.GetAPIReader().Get(context.Background(), client.ObjectKey{Namespace: DefaultNamespace, Name: "non-existent"}, obj); err != nil {
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

// HasAnnotation is a helper function to determine whether an object has a given annotation, and whether it's
// to the value provided.
func HasAnnotation(obj client.Object, key, val string) bool {
	if v, ok := obj.GetAnnotations()[key]; ok {
		if v == val {
			return true
		}
	}
	return false
}

// MatchesIngressClassName indicates whether or not an object indicates that it's supported by the ingress class name provided.
func MatchesIngressClassName(obj client.Object, ingressClassName string) bool {
	if ing, ok := obj.(*netv1.Ingress); ok {
		if ing.Spec.IngressClassName != nil && *ing.Spec.IngressClassName == ingressClassName {
			return true
		}
	}
	return HasAnnotation(obj, annotations.IngressClassKey, ingressClassName)
}
