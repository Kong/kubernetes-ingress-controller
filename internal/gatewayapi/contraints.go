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
		SecretObjectReference |
		// TODO TRR the plugin isn't really a backend in this sense, so that's a bit messy
		// Do we need to distinguish between backend and other types of reference? at what level?
		PluginLabelReference
}

type PluginLabelReference struct {
	Namespace *string
	Name      string
}
