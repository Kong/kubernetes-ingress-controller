package ctrlutils

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// classSpec indicates the fieldName for objects which support indicating their Ingress Class by spec
const classSpec = "IngressClassName"

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
func HasAnnotation(obj client.Object, key, expectedValue string) bool {
	foundValue, ok := obj.GetAnnotations()[key]
	return ok && foundValue == expectedValue
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

type objWithIngressClassNameSpec struct {
	Spec struct{ IngressClassName *string }
}

// GeneratePredicateFuncsForIngressClassFilter builds a controller-runtime reconcilation predicate function which filters out objects
// which do not have the "kubernetes.io/ingress.class" annotation configured and set to the provided value or in their .spec.
func GeneratePredicateFuncsForIngressClassFilter(name string, specCheckEnabled, annotationCheckEnabled bool) predicate.Funcs {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		if annotationCheckEnabled && IsIngressClassAnnotationConfigured(obj, name) {
			return true
		}
		if specCheckEnabled && IsIngressClassSpecConfigured(obj, name) {
			return true
		}
		return false
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		if annotationCheckEnabled && IsIngressClassAnnotationConfigured(e.ObjectOld, name) || IsIngressClassAnnotationConfigured(e.ObjectNew, name) {
			return true
		}
		if specCheckEnabled && IsIngressClassSpecConfigured(e.ObjectOld, name) || IsIngressClassSpecConfigured(e.ObjectNew, name) {
			return true
		}
		return false
	}
	return preds
}

// IsIngressClassAnnotationConfigured determines whether an object has an ingress.class annotation configured that
// matches the provide IngressClassName (and is therefore an object configured to be reconciled by that class).
//
// NOTE: keep in mind that the ingress.class annotation is deprecated and will be removed in a future release
//       of Kubernetes in favor of the .spec based implementation.
func IsIngressClassAnnotationConfigured(obj client.Object, expectedIngressClassName string) bool {
	if foundIngressClassName, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
		if foundIngressClassName == expectedIngressClassName {
			return true
		}
	}
	return false
}

// IsIngressClassAnnotationConfigured determines whether an object has IngressClassName field in its spec and whether the value
// matches the provide IngressClassName (and is therefore an object configured to be reconciled by that class).
func IsIngressClassSpecConfigured(obj client.Object, expectedIngressClassName string) bool {
	switch obj := obj.(type) {
	case *netv1.Ingress:
		return obj.Spec.IngressClassName != nil && *obj.Spec.IngressClassName == expectedIngressClassName
	}
	return false
}

// GetSemver attempts to provide a semver.Version for the provided semver string and produces an error if it can not be parsed.
func GetSemver(v string) (semver.Version, error) {
	// fix enterprise edition semver adding patch number
	// fix enterprise edition version with dash
	// fix bad version formats like 0.13.0preview1
	re := regexp.MustCompile(`(\d+\.\d+)(?:[\.-](\d+))?(?:\-?(.+)$|$)`)
	m := re.FindStringSubmatch(v)
	if len(m) != 4 {
		return semver.Version{}, fmt.Errorf("Unknown Kong version : '%v'", v)
	}
	if m[2] == "" {
		m[2] = "0"
	}
	if m[3] != "" {
		m[3] = "-" + strings.Replace(m[3], "enterprise-edition", "enterprise", 1)
		m[3] = strings.Replace(m[3], ".", "", -1)
	}
	v = fmt.Sprintf("%s.%s%s", m[1], m[2], m[3])
	return semver.Make(v)
}
