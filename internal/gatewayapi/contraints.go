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

// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/6000
// This is currently a "BackendRefT" but it's used as a generic target object.
// Needs to be renamed as such. It should be able to handle anything that satisfies
// https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client#Object
// ReferenceGrants deal with anything that has a Group, Kind, Namespace, and sometimes a Name.
// client.Object wraps two interfaces that have methods to surface all of those.

type BackendRefT interface {
	BackendRef |
		SecretObjectReference |
		PluginLabelReference
}

type PluginLabelReference struct {
	Namespace *string
	Name      string
}
