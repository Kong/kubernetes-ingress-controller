package ctrlutils

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	netv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
)

// classSpec indicates the fieldName for objects which support indicating their Ingress Class by spec
const classSpec = "IngressClassName"

// CleanupFinalizer removes an object finalizer from an object which is currently being deleted.
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

	if _, ok := obj.(*knative.Ingress); ok {
		return HasAnnotation(obj, annotations.KnativeIngressClassKey, ingressClassName)
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

// IsObjectSupported is a helper function to check if the object has any configuring that
// indicates it is supported for a given ingress class name.
func IsObjectSupported(obj client.Object, ingressClassName string) bool {
	// currently we short circuit on services and endpoints, storing all that are found.
	// See: https://github.com/Kong/kubernetes-ingress-controller/issues/1259
	if _, ok := obj.(*corev1.Service); ok {
		return true
	}
	if _, ok := obj.(*corev1.Endpoints); ok {
		return true
	}
	return IsIngressClassAnnotationConfigured(obj, ingressClassName) || IsIngressClassSpecConfigured(obj, ingressClassName)
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

	if foundIngressClassName, ok := obj.GetAnnotations()[annotations.KnativeIngressClassKey]; ok {
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

// CRDExists returns false if CRD does not exist
func CRDExists(client client.Client, gvr schema.GroupVersionResource) bool {
	_, err := client.RESTMapper().KindFor(gvr)
	if meta.IsNoMatchError(err) {
		return false
	}
	return true
}

// Convert2ClientObject is a convenience method to convert normal Kubernetes objects into
// controller-runtime's client.Object type for any of our supported APIs.
func Convert2ClientObject(obj interface{}) (client.Object, error) {
	var cobj client.Object
	switch obj := obj.(type) {
	// Kubernetes Core API Support
	case *netv1beta1.Ingress:
		cobj = obj
	case netv1beta1.Ingress:
		cobj = &obj
	case *netv1.Ingress:
		cobj = obj
	case netv1.Ingress:
		cobj = &obj
	case *corev1.Service:
		cobj = obj
	case corev1.Service:
		cobj = &obj
	case *corev1.Endpoints:
		cobj = obj
	case corev1.Endpoints:
		cobj = &obj
	// Kong API Support
	case *kongv1.KongPlugin:
		cobj = obj
	case kongv1.KongPlugin:
		cobj = &obj
	case *kongv1.KongClusterPlugin:
		cobj = obj
	case kongv1.KongClusterPlugin:
		cobj = &obj
	case *kongv1.KongConsumer:
		cobj = obj
	case kongv1.KongConsumer:
		cobj = &obj
	case *kongv1.KongIngress:
		cobj = obj
	case kongv1.KongIngress:
		cobj = &obj
	case *kongv1beta1.TCPIngress:
		cobj = obj
	case kongv1beta1.TCPIngress:
		cobj = &obj
	case *kongv1beta1.UDPIngress:
		cobj = obj
	case kongv1beta1.UDPIngress:
		cobj = &obj
	// 3rd Party API Support
	case *knative.Ingress:
		cobj = obj
	case knative.Ingress:
		cobj = &obj
	default:
		return nil, fmt.Errorf("unexpected object type found: %T", obj)
	}
	return cobj, nil
}
