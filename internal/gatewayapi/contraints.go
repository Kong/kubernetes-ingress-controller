package gatewayapi

import "sigs.k8s.io/controller-runtime/pkg/client"

type HostnameT interface {
	Hostname | string
}

type RouteT interface {
	client.Object

	*HTTPRoute |
		*UDPRoute |
		*TCPRoute |
		*TLSRoute |
		*GRPCRoute
}

type BackendRefT interface {
	BackendRef |
		SecretObjectReference
}
