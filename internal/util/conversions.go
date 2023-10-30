package util

import (
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// -----------------------------------------------------------------------------
// Type conversion Utilities
// -----------------------------------------------------------------------------

// StringToGatewayAPIHostname converts a string to a gatewayapi.Hostname.
func StringToGatewayAPIHostname(hostname string) gatewayapi.Hostname {
	return (gatewayapi.Hostname)(hostname)
}

// StringToGatewayAPIHostnamePtr converts a string to a *gatewayapi.Hostname.
func StringToGatewayAPIHostnamePtr(hostname string) *gatewayapi.Hostname {
	return lo.ToPtr(gatewayapi.Hostname(hostname))
}

// StringToGatewayAPIHostnameV1Beta1Ptr converts a string to a *gatewayapi.Hostname.
func StringToGatewayAPIHostnameV1Beta1Ptr(hostname string) *gatewayapi.Hostname {
	return lo.ToPtr(gatewayapi.Hostname(hostname))
}

// StringToGatewayAPIKindV1Alpha2Ptr converts a string to a *gatewayapi.Kind.
func StringToGatewayAPIKindV1Alpha2Ptr(kind string) *gatewayapi.Kind {
	return lo.ToPtr(gatewayapi.Kind(kind))
}

// StringToGatewayAPIKindPtr converts a string to a *gatewayapi.Kind.
func StringToGatewayAPIKindPtr(kind string) *gatewayapi.Kind {
	return lo.ToPtr(gatewayapi.Kind(kind))
}
