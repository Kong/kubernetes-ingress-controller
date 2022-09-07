package parser

import (
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/types"
)

// refChecker is a wrapper type that facilitates checking whether a backenRef is allowed
// by a referenceGrantTo set.
type refChecker[T types.BackendRefT] struct {
	backendRef T
}

func newRefChecker[T types.BackendRefT](ref T) refChecker[T] {
	return refChecker[T]{
		backendRef: ref,
	}
}

// IsRefAllowedByGrant is a wrapper on top of isRefAllowedByGrant checks if backendRef (that RefChecker
// holds) is permitted by the provided namespace-indexed ReferenceGrantTo set: allowedRefs.
// allowedRefs is assumed to contain Tos that only match the backendRef's parent's From, as returned by
// getPermittedForReferenceGrantFrom.
func (rc refChecker[T]) IsRefAllowedByGrant(
	allowedRefs map[gatewayv1beta1.Namespace][]gatewayv1alpha2.ReferenceGrantTo,
) bool {
	switch br := (interface{})(rc.backendRef).(type) {
	case gatewayv1beta1.BackendRef:
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
	case gatewayv1alpha2.BackendRef:
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

	case gatewayv1alpha2.SecretObjectReference:
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

	case gatewayv1beta1.SecretObjectReference:
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
	allowed map[gatewayv1beta1.Namespace][]gatewayv1alpha2.ReferenceGrantTo,
) bool {
	if namespace == nil {
		// local references are always fine
		return true
	}
	for _, to := range allowed[gatewayv1beta1.Namespace(*namespace)] {
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
