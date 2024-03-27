//nolint:revive
package gateway

import (
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type (
	BackendRef        = gatewayv1.BackendRef
	Gateway           = gatewayv1.Gateway
	GatewayAddress    = gatewayv1.GatewayAddress
	GatewayClass      = gatewayv1.GatewayClass
	Group             = gatewayv1.Group
	Hostname          = gatewayv1.Hostname
	HTTPRoute         = gatewayv1.HTTPRoute
	Kind              = gatewayv1.Kind
	Listener          = gatewayv1.Listener
	ListenerStatus    = gatewayv1.ListenerStatus
	Namespace         = gatewayv1.Namespace
	ObjectName        = gatewayv1.ObjectName
	ParentReference   = gatewayv1.ParentReference
	PortNumber        = gatewayv1.PortNumber
	ProtocolType      = gatewayv1.ProtocolType
	RouteParentStatus = gatewayv1.RouteParentStatus
	SectionName       = gatewayv1.SectionName

	TCPRoute  = gatewayv1alpha2.TCPRoute
	UDPRoute  = gatewayv1alpha2.UDPRoute
	TLSRoute  = gatewayv1alpha2.TLSRoute
	GRPCRoute = gatewayv1alpha2.GRPCRoute
)

const (
	HTTPProtocolType  = gatewayv1.HTTPProtocolType
	HTTPSProtocolType = gatewayv1.HTTPSProtocolType
	TLSProtocolType   = gatewayv1.TLSProtocolType
	TCPProtocolType   = gatewayv1.TCPProtocolType
	UDPProtocolType   = gatewayv1.UDPProtocolType
)
