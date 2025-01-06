package utils

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
)

// ObjectWithControlPlaneRef is an interface that represents an object that has a control plane reference.
type ObjectWithControlPlaneRef interface {
	GetControlPlaneRef() *kongv1alpha1.ControlPlaneRef
}

// GenerateCPReferenceMatchesPredicate generates a predicate function that filters out objects that have a control plane
// reference set to a value other than 'kic'.
func GenerateCPReferenceMatchesPredicate[T ObjectWithControlPlaneRef]() predicate.Predicate {
	return predicate.NewPredicateFuncs(func(o client.Object) bool {
		c, ok := o.(T)
		if !ok {
			return false
		}
		if cpRef := c.GetControlPlaneRef(); cpRef != nil {
			// If the cpRef is set, reconcile the object only if it is set explicitly to 'kic'.
			return cpRef.Type == kongv1alpha1.ControlPlaneRefKIC
		}
		// If there's no cpRef set, we should reconcile it as by default it's 'kic'.
		return true
	})
}
