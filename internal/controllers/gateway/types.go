//nolint:revive
package gateway

import (
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

type (
	BackendRef        = gatewayv1beta1.BackendRef
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
	PortNumber        = gatewayv1beta1.PortNumber
	ProtocolType      = gatewayv1beta1.ProtocolType
	RouteParentStatus = gatewayv1beta1.RouteParentStatus
	SectionName       = gatewayv1beta1.SectionName

	TCPRoute = gatewayv1alpha2.TCPRoute
	UDPRoute = gatewayv1alpha2.UDPRoute
	TLSRoute = gatewayv1alpha2.TLSRoute
)

const (
	HTTPProtocolType  = gatewayv1beta1.HTTPProtocolType
	HTTPSProtocolType = gatewayv1beta1.HTTPSProtocolType
	TLSProtocolType   = gatewayv1beta1.TLSProtocolType
	TCPProtocolType   = gatewayv1beta1.TCPProtocolType
	UDPProtocolType   = gatewayv1beta1.UDPProtocolType
)
