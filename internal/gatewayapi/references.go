package gatewayapi

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RefChecker is a wrapper type that facilitates checking whether a backendRef is allowed
// by a referenceGrantTo set.
type RefChecker[T BackendRefT] struct {
	target     client.Object
	backendRef T
}

// NewRefCheckerForRoute returns a RefChecker for the provided route and backendRef.
func NewRefCheckerForRoute[T BackendRefT](route client.Object, ref T) RefChecker[T] {
	return RefChecker[T]{
		target:     route,
		backendRef: ref,
	}
}

func NewRefCheckerForKongPlugin[T BackendRefT](target client.Object, requester T) RefChecker[T] {
	return RefChecker[T]{
		target:     target,
		backendRef: requester,
	}
}

// IsRefAllowedByGrant is a wrapper on top of isRefAllowedByGrant checks if backendRef (that RefChecker
// holds) is permitted by the provided namespace-indexed ReferenceGrantTo set: allowedRefs.
// allowedRefs is assumed to contain Tos that only match the backendRef's parent's From, as returned by
// getPermittedForReferenceGrantFrom.
func (rc RefChecker[T]) IsRefAllowedByGrant(
	allowedRefs map[Namespace][]ReferenceGrantTo,
) bool {
	switch br := (interface{})(rc.backendRef).(type) {
	case BackendRef:
		// NOTE TRR this is a catch-all that the plugins technically won't need as-is because they have their own
		// inherent namespace check: if no namespace specified the ref is assumed allowed because it must be local
		if br.Namespace == nil {
			return true
		}

		// If the namespace is specified but is the same as the route's namespace, then the ref is allowed.
		if rc.target.GetNamespace() == string(*br.Namespace) {
			// NOTE TRR however the plugin stuff does not check the "did someone ask for their own namespace" case and assumes
			// references will only be foreign
			return true
		}

		return isRefAllowedByGrant(
			(*string)(br.Namespace),
			(string)(br.Name),
			(string)(*br.Group),
			(string)(*br.Kind),
			allowedRefs,
		)

	// TODO TRR how do these acually differ? we don't check if the secret ref is specified but same, but in practice
	// that's only for the "you're holding it wrong" usage
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

	// TODO TRR this isn't really in the GWAPI space, but for temporary consistency, it lives there
	case PluginLabelReference:
		if br.Namespace == nil {
			return true
		}

		return isRefAllowedByGrant(
			(br.Namespace),
			(br.Name),
			"configuration.konghq.com", // TODO TRR we have some const for this somewhere? alternately maybe it's fed in like the rest
			"KongPlugin",               // TODO TRR ditto
			allowedRefs,
		)
	}

	return false
}

// isRefAllowedByGrant checks if backendRef is permitted by the provided namespace-indexed ReferenceGrantTo set: allowed.
// allowed is assumed to contain Tos that only match the backendRef's parent's From, as returned by
// getPermittedForReferenceGrantFrom.
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
