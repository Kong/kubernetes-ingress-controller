package util

import (
	"github.com/samber/lo"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// -----------------------------------------------------------------------------
// Type conversion Utilities
// -----------------------------------------------------------------------------

// StringToGatewayAPIHostname converts a string to a gatewayv1alpha2.Hostname.
func StringToGatewayAPIHostname(hostname string) gatewayv1alpha2.Hostname {
	return (gatewayv1alpha2.Hostname)(hostname)
}

// StringToGatewayAPIHostnameV1Beta1 converts a string to a gatewayv1beta1.Hostname.
func StringToGatewayAPIHostnameV1Beta1(hostname string) gatewayv1beta1.Hostname {
	return (gatewayv1beta1.Hostname)(hostname)
}

// StringToGatewayAPIHostnamePtr converts a string to a *gatewayv1beta1.Hostname.
func StringToGatewayAPIHostnamePtr(hostname string) *gatewayv1beta1.Hostname {
	return lo.ToPtr(gatewayv1beta1.Hostname(hostname))
}

// StringToGatewayAPIHostnameV1Beta1Ptr converts a string to a *gatewayv1beta1.Hostname.
func StringToGatewayAPIHostnameV1Beta1Ptr(hostname string) *gatewayv1beta1.Hostname {
	return lo.ToPtr(gatewayv1beta1.Hostname(hostname))
}

// StringToGatewayAPIKindV1Alpha2Ptr converts a string to a *gatewayv1alpha2.Kind.
func StringToGatewayAPIKindV1Alpha2Ptr(kind string) *gatewayv1alpha2.Kind {
	return lo.ToPtr(gatewayv1alpha2.Kind(kind))
}

// StringToGatewayAPIKindPtr converts a string to a *gatewayv1beta1.Kind.
func StringToGatewayAPIKindPtr(kind string) *gatewayv1beta1.Kind {
	return lo.ToPtr(gatewayv1beta1.Kind(kind))
}
