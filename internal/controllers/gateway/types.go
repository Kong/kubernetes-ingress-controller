//nolint:revive
package gateway

import (
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

type (
	Gateway           = gatewayv1beta1.Gateway
	GatewayAddress    = gatewayv1beta1.GatewayAddress
	GatewayClass      = gatewayv1beta1.GatewayClass
	Group             = gatewayv1beta1.Group
	Hostname          = gatewayv1beta1.Hostname
	HTTPRoute         = gatewayv1beta1.HTTPRoute
	Kind              = gatewayv1beta1.Kind
	Listener          = gatewayv1beta1.Listener
	ListenerStatus    = gatewayv1beta1.ListenerStatus
	Namespace         = gatewayv1beta1.Namespace
	ObjectName        = gatewayv1beta1.ObjectName
	ParentReference   = gatewayv1beta1.ParentReference
	RouteParentStatus = gatewayv1beta1.RouteParentStatus
	PortNumber        = gatewayv1beta1.PortNumber
	ProtocolType      = gatewayv1beta1.ProtocolType
	SectionName       = gatewayv1beta1.SectionName
)

const (
	HTTPProtocolType  = gatewayv1beta1.HTTPProtocolType
	HTTPSProtocolType = gatewayv1beta1.HTTPSProtocolType
	TLSProtocolType   = gatewayv1beta1.TLSProtocolType
	TCPProtocolType   = gatewayv1beta1.TCPProtocolType
	UDPProtocolType   = gatewayv1beta1.UDPProtocolType
)
