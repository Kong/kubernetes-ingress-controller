package translator

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// refChecker is a wrapper type that facilitates checking whether a backenRef is allowed
// by a referenceGrantTo set.
type refChecker[T gatewayapi.BackendRefT] struct {
	route      client.Object
	backendRef T
}

// newRefCheckerForRoute returns a refChecker for the provided route and backendRef.
func newRefCheckerForRoute[T gatewayapi.BackendRefT](route client.Object, ref T) refChecker[T] {
	return refChecker[T]{
		route:      route,
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

		// If the namespace is specified but is the same as the route's namespace, then the ref is allowed.
		if rc.route.GetNamespace() == string(*br.Namespace) {
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
