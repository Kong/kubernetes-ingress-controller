package parser

import (
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// refChecker is a wrapper type that facilitates checking whether a backenRef is allowed
// by a referenceGrantTo set.
type refChecker[T gatewayapi.BackendRefT] struct {
	backendRef T
}

func newRefChecker[T gatewayapi.BackendRefT](ref T) refChecker[T] {
	return refChecker[T]{
		backendRef: ref,
	}
}

// IsRefAllowedByGrant is a wrapper on top of isRefAllowedByGrant checks if backendRef (that RefChecker
// holds) is permitted by the provided namespace-indexed ReferenceGrantTo set: allowedRefs.
// allowedRefs is assumed to contain Tos that only match the backendRef's parent's From, as returned by
// getPermittedForReferenceGrantFrom.
func (rc refChecker[T]) IsRefAllowedByGrant(
	allowedRefs map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo,
) bool {
	switch br := (interface{})(rc.backendRef).(type) {
	case gatewayapi.BackendRef:
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

	case gatewayapi.SecretObjectReference:
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
	allowed map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo,
) bool {
	if namespace == nil {
		// local references are always fine
		return true
	}
	for _, to := range allowed[gatewayapi.Namespace(*namespace)] {
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
