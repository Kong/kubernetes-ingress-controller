package utils

import (
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
)

const defaultIngressClassAnnotation = "ingressclass.kubernetes.io/is-default-class"

// IsDefaultIngressClass returns whether an IngressClass is the default IngressClass.
func IsDefaultIngressClass(ingressClass *netv1.IngressClass) bool {
	return ingressClass.ObjectMeta.Annotations[defaultIngressClassAnnotation] == "true"
}

// MatchesIngressClass indicates whether or not an object belongs to a given ingress class.
func MatchesIngressClass(obj client.Object, controllerIngressClass string, isDefault bool) bool {
	objectIngressClass := obj.GetAnnotations()[annotations.IngressClassKey]
	if isDefault && IsIngressClassEmpty(obj) {
		return true
	}
	if ing, isV1Ingress := obj.(*netv1.Ingress); isV1Ingress {
		if ing.Spec.IngressClassName != nil && *ing.Spec.IngressClassName == controllerIngressClass {
			return true
		}
	}
	return objectIngressClass == controllerIngressClass
}

// GeneratePredicateFuncsForIngressClassFilter builds a controller-runtime reconciliation predicate function which filters out objects
// which have their ingress class set to the a value other than the controller class.
func GeneratePredicateFuncsForIngressClassFilter[T client.Object](name string) predicate.TypedFuncs[T] {
	preds := predicate.NewTypedPredicateFuncs(func(obj T) bool {
		// we assume true for isDefault here because the predicates have no client and cannot check if the class is
		// default. classless and are filtered out by Reconcile() if the configured class is not the default class
		return MatchesIngressClass(obj, name, true)
	})
	preds.UpdateFunc = func(e event.TypedUpdateEvent[T]) bool {
		return MatchesIngressClass(e.ObjectOld, name, true) || MatchesIngressClass(e.ObjectNew, name, true)
	}
	return preds
}

// IsIngressClassEmpty returns true if an object has no ingress class information or false otherwise.
func IsIngressClassEmpty(obj client.Object) bool {
	switch obj := obj.(type) {
	case *netv1.Ingress:
		// netv1.Ingress is the only kind with an explicit IngressClassName field. All other resources use annotations
		// the annotation is deprecated for netv1.Ingress, and the older Ingress versions are themselves deprecated
		// our CRDs use the annotation, but should probably transition to a field eventually to align with Ingress
		if _, ok := obj.GetAnnotations()[annotations.IngressClassKey]; !ok {
			return obj.Spec.IngressClassName == nil
		}
		return false
	default:
		if _, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
			return false
		}
		return true
	}
}

// CRDExists returns false if CRD does not exist.
func CRDExists(restMapper meta.RESTMapper, gvr schema.GroupVersionResource) bool {
	_, err := restMapper.KindsFor(gvr)
	return err == nil
}
