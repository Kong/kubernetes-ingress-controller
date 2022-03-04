package utils

import (
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

const defaultIngressClassAnnotation = "ingressclass.kubernetes.io/is-default-class"

// HasAnnotation is a helper function to determine whether an object has a given annotation, and whether it's
// to the value provided.
func HasAnnotation(obj client.Object, key, expectedValue string) bool {
	foundValue, ok := obj.GetAnnotations()[key]
	return ok && foundValue == expectedValue
}

// IsDefaultIngressClass returns whether an IngressClass is the default IngressClass
func IsDefaultIngressClass(obj client.Object) bool {
	if ingressClass, ok := obj.(*netv1.IngressClass); ok {
		return ingressClass.ObjectMeta.Annotations[defaultIngressClassAnnotation] == "true"
	}
	return false
}

// MatchesIngressClassName indicates whether or not an object indicates that it's supported by the ingress class name provided.
func MatchesIngressClassName(obj client.Object, ingressClassName string, isDefault bool) bool {
	if ing, ok := obj.(*netv1.Ingress); ok {
		if ing.Spec.IngressClassName != nil && *ing.Spec.IngressClassName == ingressClassName {
			return true
		} else if ing.Spec.IngressClassName == nil && isDefault {
			_, standard := obj.GetAnnotations()[annotations.IngressClassKey]
			_, knative := obj.GetAnnotations()[annotations.KnativeIngressClassKey]
			if !standard && !knative {
				return true
			}
		}
	}

	if _, ok := obj.(*knative.Ingress); ok {
		return HasAnnotation(obj, annotations.KnativeIngressClassKey, ingressClassName)
	}

	return HasAnnotation(obj, annotations.IngressClassKey, ingressClassName)
}

// GeneratePredicateFuncsForIngressClassFilter builds a controller-runtime reconciliation predicate function which filters out objects
// which do not have the "kubernetes.io/ingress.class" annotation configured and set to the provided value or in their .spec.
func GeneratePredicateFuncsForIngressClassFilter(name string, specCheckEnabled, annotationCheckEnabled bool) predicate.Funcs {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		if annotationCheckEnabled && IsIngressClassAnnotationConfigured(obj, name) {
			return true
		}
		if specCheckEnabled {
			if IsIngressClassSpecConfigured(obj, name) {
				return true
			}
			if IsIngressClassSpecEmpty(obj) {
				// we include Ingresses with _no_ ingressClassName in case we're handling the default IngressClass,
				// and will filter them out if not in MatchesIngressClassName()
				return true
			}
		}
		return false
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		if annotationCheckEnabled && IsIngressClassAnnotationConfigured(e.ObjectOld, name) || IsIngressClassAnnotationConfigured(e.ObjectNew, name) {
			return true
		}
		if specCheckEnabled {
			if IsIngressClassSpecConfigured(e.ObjectOld, name) || IsIngressClassSpecConfigured(e.ObjectNew, name) {
				return true
			}
			if IsIngressClassSpecEmpty(e.ObjectOld) || IsIngressClassSpecEmpty(e.ObjectNew) {
				return true
			}
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

	if foundIngressClassName, ok := obj.GetAnnotations()[annotations.KnativeIngressClassKey]; ok {
		if foundIngressClassName == expectedIngressClassName {
			return true
		}
	}

	return false
}

// IsIngressClassSpecConfigured determines whether an object has IngressClassName field in its spec and whether the value
// matches the provide IngressClassName (and is therefore an object configured to be reconciled by that class).
func IsIngressClassSpecConfigured(obj client.Object, expectedIngressClassName string) bool {
	switch obj := obj.(type) {
	case *netv1.Ingress:
		return obj.Spec.IngressClassName != nil && *obj.Spec.IngressClassName == expectedIngressClassName
	case *netv1beta1.Ingress:
		return obj.Spec.IngressClassName != nil && *obj.Spec.IngressClassName == expectedIngressClassName
	case *extv1beta1.Ingress:
		return obj.Spec.IngressClassName != nil && *obj.Spec.IngressClassName == expectedIngressClassName
	}
	return false
}

// IsIngressClassSpecEmpty checks if a networking/v1 Ingress has no ingressClassName set
func IsIngressClassSpecEmpty(obj client.Object) bool {
	switch obj := obj.(type) {
	case *netv1.Ingress:
		return obj.Spec.IngressClassName == nil
	default:
		return false
	}
}

// CRDExists returns false if CRD does not exist
func CRDExists(client client.Client, gvr schema.GroupVersionResource) bool {
	_, err := client.RESTMapper().KindFor(gvr)
	return !meta.IsNoMatchError(err)
}
