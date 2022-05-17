package utils

import (
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

const defaultIngressClassAnnotation = "ingressclass.kubernetes.io/is-default-class"

// IsDefaultIngressClass returns whether an IngressClass is the default IngressClass
func IsDefaultIngressClass(obj client.Object) bool {
	if ingressClass, ok := obj.(*netv1.IngressClass); ok {
		return ingressClass.ObjectMeta.Annotations[defaultIngressClassAnnotation] == "true"
	}
	return false
}

// MatchesIngressClass indicates whether or not an object belongs to a given ingress class
func MatchesIngressClass(obj client.Object, controllerIngressClass string, isDefault bool) bool {
	objectIngressClass := obj.GetAnnotations()[annotations.IngressClassKey]
	objectKnativeClass := obj.GetAnnotations()[annotations.KnativeIngressClassKey]
	objectKnativeClassAlt := obj.GetAnnotations()[annotations.KnativeIngressClassDeprecatedKey]
	if isDefault && IsIngressClassEmpty(obj) {
		return true
	}
	if ing, isV1Ingress := obj.(*netv1.Ingress); isV1Ingress {
		if ing.Spec.IngressClassName != nil && *ing.Spec.IngressClassName == controllerIngressClass {
			return true
		}
	}

	switch controllerIngressClass {
	case objectIngressClass:
		return true
	case objectKnativeClass:
		return true
	case objectKnativeClassAlt:
		return true
	}

	return false
}

// GeneratePredicateFuncsForIngressClassFilter builds a controller-runtime reconciliation predicate function which filters out objects
// which have their ingress class set to the a value other than the controller class
func GeneratePredicateFuncsForIngressClassFilter(name string) predicate.Funcs {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		// we assume true for isDefault here because the predicates have no client and cannot check if the class is
		// default. classless and are filtered out by Reconcile() if the configured class is not the default class
		return MatchesIngressClass(obj, name, true)
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		return MatchesIngressClass(e.ObjectOld, name, true) || MatchesIngressClass(e.ObjectNew, name, true)
	}
	return preds
}

// IsIngressClassEmpty returns true if an object has no ingress class information or false otherwise
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
		if _, ok := obj.GetAnnotations()[annotations.KnativeIngressClassKey]; ok {
			return false
		}
		if _, ok := obj.GetAnnotations()[annotations.KnativeIngressClassDeprecatedKey]; ok {
			return false
		}
		return true
	}
}

// CRDExists returns false if CRD does not exist
func CRDExists(client client.Client, gvr schema.GroupVersionResource) bool {
	_, err := client.RESTMapper().KindFor(gvr)
	return !meta.IsNoMatchError(err)
}
