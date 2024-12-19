package util

import (
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// -----------------------------------------------------------------------------
// Type conversion Utilities
// -----------------------------------------------------------------------------

// StringToTypedPtr converts a string to pointer to a typed designated by the provided type parameter.
func StringToTypedPtr[
	TT *T,
	T ~string,
](s string) TT {
	ret := T(s)
	return &ret
}

// StringToGatewayAPIKindPtr converts a string to a *gatewayapi.Kind.
func StringToGatewayAPIKindPtr(kind string) *gatewayapi.Kind {
	return StringToTypedPtr[*gatewayapi.Kind](kind)
}
