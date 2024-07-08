package gatewayapi

import (
	"reflect"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RefChecker is a wrapper type that facilitates checking whether a backendRef is allowed
// by a referenceGrantTo set.
type RefChecker[T BackendRefT] struct {
	target     client.Object
	backendRef T
	log        logr.Logger
}

// NewRefCheckerForRoute returns a RefChecker for the provided route and backendRef.
func NewRefCheckerForRoute[T BackendRefT](log logr.Logger, route client.Object, ref T) RefChecker[T] {
	return RefChecker[T]{
		target:     route,
		backendRef: ref,
	}
}

func NewRefCheckerForKongPlugin[T BackendRefT](log logr.Logger, target client.Object, requester T) RefChecker[T] {
	return RefChecker[T]{
		target:     target,
		backendRef: requester,
	}
}

// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/6000 this has separate cases for different types,
// but doesn't do anything meaningfully different for them (it only fills in some default info that should be available
// from the involved objects' methods). We want a generic utility that handles relationship checks for any
// client.Object. Any behavior particular to specific combinations of GVK->GVK relationships should be handled in the
// code that implements those relationships.

// IsRefAllowedByGrant is a wrapper on top of isRefAllowedByGrant checks if backendRef (that RefChecker
// holds) is permitted by the provided namespace-indexed ReferenceGrantTo set: allowedRefs.
// allowedRefs is assumed to contain Tos that only match the backendRef's parent's From, as returned by
// GetPermittedForReferenceGrantFrom.
func (rc RefChecker[T]) IsRefAllowedByGrant(
	allowedRefs map[Namespace][]ReferenceGrantTo,
) bool {
	switch br := (interface{})(rc.backendRef).(type) {
	case BackendRef:
		// NOTE https://github.com/Kong/kubernetes-ingress-controller/issues/6000
		// This is a catch-all that the plugins technically won't need as-is because they have their own
		// inherent namespace check: if no namespace specified the ref is assumed allowed because it must be local.
		if br.Namespace == nil {
			return true
		}

		// If the namespace is specified but is the same as the route's namespace, then the ref is allowed.
		if rc.target.GetNamespace() == string(*br.Namespace) {
			return true
		}

		return isRefAllowedByGrant(
			(*string)(br.Namespace),
			(string)(br.Name),
			(string)(*br.Group),
			(string)(*br.Kind),
			allowedRefs,
		)

	case SecretObjectReference:
		if br.Namespace == nil {
			return true
		}

		return isRefAllowedByGrant(
			(*string)(br.Namespace),
			(string)(br.Name),
			(string)(*br.Group),
			(string)(*br.Kind),
			allowedRefs,
		)

	case PluginLabelReference:
		if br.Namespace == nil {
			return true
		}

		return isRefAllowedByGrant(
			(br.Namespace),
			(br.Name),
			"configuration.konghq.com", // TODO https://github.com/Kong/kubernetes-ingress-controller/issues/6000
			"KongPlugin",               // TODO These magic strings should become unnecessary once we work with client.Object
			allowedRefs,
		)
	}

	return false
}

// isRefAllowedByGrant checks if backendRef is permitted by the provided namespace-indexed ReferenceGrantTo set: allowed.
// allowed is assumed to contain Tos that only match the backendRef's parent's From, as returned by
// GetPermittedForReferenceGrantFrom.
func isRefAllowedByGrant(
	namespace *string,
	name string,
	group string,
	kind string,
	allowed map[Namespace][]ReferenceGrantTo,
) bool {
	if namespace == nil {
		// local references are always fine
		return true
	}
	for _, to := range allowed[Namespace(*namespace)] {
		if string(to.Group) == group && string(to.Kind) == kind {
			if to.Name != nil {
				if string(*to.Name) == name {
					return true
				}
			} else {
				// if no referent name specified, matching group/kind is sufficient
				return true
			}
		}
	}

	return false
}

// GetPermittedForReferenceGrantFrom takes a ReferenceGrant From (a namespace, group, and kind) and returns a map
// from a namespace to a slice of ReferenceGrant Tos. When a To is included in the slice, the key namespace has a
// ReferenceGrant with those Tos and the input From.
func GetPermittedForReferenceGrantFrom(
	log logr.Logger,
	from ReferenceGrantFrom,
	grants []*ReferenceGrant,
) map[Namespace][]ReferenceGrantTo {
	allowed := make(map[Namespace][]ReferenceGrantTo)
	// loop over all From values in all grants. if we find a match, add all Tos to the list of Tos allowed for the
	// grant namespace. this technically could add duplicate copies of the Tos if there are duplicate Froms (it makes
	// no sense to add them, but it's allowed), but duplicate Tos are harmless (we only care about having at least one
	// matching To when checking if a ReferenceGrant allows a reference)
	for _, grant := range grants {
		for _, otherFrom := range grant.Spec.From {
			if reflect.DeepEqual(from, otherFrom) {
				allowed[Namespace(grant.ObjectMeta.Namespace)] = append(allowed[Namespace(grant.ObjectMeta.Namespace)], grant.Spec.To...)
			}
		}
	}

	return allowed
}
