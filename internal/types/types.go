package types

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type HostnameT interface {
	gatewayv1.Hostname | string
}

type RouteT interface {
	client.Object

	*gatewayv1.HTTPRoute |
		*gatewayv1alpha2.UDPRoute |
		*gatewayv1alpha2.TCPRoute |
		*gatewayv1alpha2.TLSRoute |
		*gatewayv1alpha2.GRPCRoute
}

type BackendRefT interface {
	gatewayv1.BackendRef |
		gatewayv1.SecretObjectReference
}
