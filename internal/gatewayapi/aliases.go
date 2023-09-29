package gatewayapi

import (
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// This file contains aliases for types and consts from the Gateway API.  Its purpose is to allow easy migration from
// one version of the Gateway API to another with minimal changes to the codebase.

type (
	AllowedRoutes             = gatewayv1beta1.AllowedRoutes
	BackendObjectReference    = gatewayv1beta1.BackendObjectReference
	BackendRef                = gatewayv1beta1.BackendRef
	CommonRouteSpec           = gatewayv1beta1.CommonRouteSpec
	Gateway                   = gatewayv1beta1.Gateway
	GatewayAddress            = gatewayv1beta1.GatewayAddress
	GatewayClass              = gatewayv1beta1.GatewayClass
	GatewayClassSpec          = gatewayv1beta1.GatewayClassSpec
	GatewayClassStatus        = gatewayv1beta1.GatewayClassStatus
	GatewayController         = gatewayv1beta1.GatewayController
	GatewayList               = gatewayv1beta1.GatewayList
	GatewaySpec               = gatewayv1beta1.GatewaySpec
	GatewayStatus             = gatewayv1beta1.GatewayStatus
	GatewayStatusAddress      = gatewayv1beta1.GatewayStatusAddress
	GatewayTLSConfig          = gatewayv1beta1.GatewayTLSConfig
	Group                     = gatewayv1beta1.Group
	HTTPBackendRef            = gatewayv1beta1.HTTPBackendRef
	HTTPHeader                = gatewayv1beta1.HTTPHeader
	HTTPHeaderFilter          = gatewayv1beta1.HTTPHeaderFilter
	HTTPHeaderMatch           = gatewayv1beta1.HTTPHeaderMatch
	HTTPHeaderName            = gatewayv1beta1.HTTPHeaderName
	HTTPMethod                = gatewayv1beta1.HTTPMethod
	HTTPPathMatch             = gatewayv1beta1.HTTPPathMatch
	HTTPQueryParamMatch       = gatewayv1beta1.HTTPQueryParamMatch
	HTTPRequestRedirectFilter = gatewayv1beta1.HTTPRequestRedirectFilter
	HTTPRoute                 = gatewayv1beta1.HTTPRoute
	HTTPRouteFilter           = gatewayv1beta1.HTTPRouteFilter
	HTTPRouteList             = gatewayv1beta1.HTTPRouteList
	HTTPRouteMatch            = gatewayv1beta1.HTTPRouteMatch
	HTTPRouteRule             = gatewayv1beta1.HTTPRouteRule
	HTTPRouteSpec             = gatewayv1beta1.HTTPRouteSpec
	HTTPRouteStatus           = gatewayv1beta1.HTTPRouteStatus
	Hostname                  = gatewayv1beta1.Hostname
	Kind                      = gatewayv1beta1.Kind
	Listener                  = gatewayv1beta1.Listener
	ListenerConditionReason   = gatewayv1beta1.ListenerConditionReason
	ListenerConditionType     = gatewayv1beta1.ListenerConditionType
	ListenerStatus            = gatewayv1beta1.ListenerStatus
	Namespace                 = gatewayv1beta1.Namespace
	ObjectName                = gatewayv1beta1.ObjectName
	ParentReference           = gatewayv1beta1.ParentReference
	PathMatchType             = gatewayv1beta1.PathMatchType
	PortNumber                = gatewayv1beta1.PortNumber
	PreciseHostname           = gatewayv1beta1.PreciseHostname
	ProtocolType              = gatewayv1beta1.ProtocolType
	ReferenceGrant            = gatewayv1beta1.ReferenceGrant
	ReferenceGrantFrom        = gatewayv1beta1.ReferenceGrantFrom
	ReferenceGrantList        = gatewayv1beta1.ReferenceGrantList
	ReferenceGrantSpec        = gatewayv1beta1.ReferenceGrantSpec
	ReferenceGrantTo          = gatewayv1beta1.ReferenceGrantTo
	RouteConditionReason      = gatewayv1beta1.RouteConditionReason
	RouteGroupKind            = gatewayv1beta1.RouteGroupKind
	RouteNamespaces           = gatewayv1beta1.RouteNamespaces
	RouteParentStatus         = gatewayv1beta1.RouteParentStatus
	RouteStatus               = gatewayv1beta1.RouteStatus
	SecretObjectReference     = gatewayv1beta1.SecretObjectReference
	SectionName               = gatewayv1beta1.SectionName

	GRPCBackendRef      = gatewayv1alpha2.GRPCBackendRef
	GRPCHeaderMatch     = gatewayv1alpha2.GRPCHeaderMatch
	GRPCHeaderName      = gatewayv1alpha2.GRPCHeaderName
	GRPCMethodMatch     = gatewayv1alpha2.GRPCMethodMatch
	GRPCMethodMatchType = gatewayv1alpha2.GRPCMethodMatchType
	GRPCRoute           = gatewayv1alpha2.GRPCRoute
	GRPCRouteList       = gatewayv1alpha2.GRPCRouteList
	GRPCRouteMatch      = gatewayv1alpha2.GRPCRouteMatch
	GRPCRouteRule       = gatewayv1alpha2.GRPCRouteRule
	GRPCRouteSpec       = gatewayv1alpha2.GRPCRouteSpec
	GRPCRouteStatus     = gatewayv1alpha2.GRPCRouteStatus
	TCPRoute            = gatewayv1alpha2.TCPRoute
	TCPRouteList        = gatewayv1alpha2.TCPRouteList
	TCPRouteRule        = gatewayv1alpha2.TCPRouteRule
	TCPRouteSpec        = gatewayv1alpha2.TCPRouteSpec
	TCPRouteStatus      = gatewayv1alpha2.TCPRouteStatus
	TLSRoute            = gatewayv1alpha2.TLSRoute
	TLSRouteList        = gatewayv1alpha2.TLSRouteList
	TLSRouteRule        = gatewayv1alpha2.TLSRouteRule
	TLSRouteSpec        = gatewayv1alpha2.TLSRouteSpec
	TLSRouteStatus      = gatewayv1alpha2.TLSRouteStatus
	UDPRoute            = gatewayv1alpha2.UDPRoute
	UDPRouteList        = gatewayv1alpha2.UDPRouteList
	UDPRouteRule        = gatewayv1alpha2.UDPRouteRule
	UDPRouteSpec        = gatewayv1alpha2.UDPRouteSpec
	UDPRouteStatus      = gatewayv1alpha2.UDPRouteStatus
)

const (
	FullPathHTTPPathModifier              = gatewayv1beta1.FullPathHTTPPathModifier
	GatewayClassConditionStatusAccepted   = gatewayv1beta1.GatewayClassConditionStatusAccepted
	GatewayClassReasonAccepted            = gatewayv1beta1.GatewayClassReasonAccepted
	GatewayConditionAccepted              = gatewayv1beta1.GatewayConditionAccepted
	GatewayConditionProgrammed            = gatewayv1beta1.GatewayConditionProgrammed
	GatewayReasonAccepted                 = gatewayv1beta1.GatewayReasonAccepted
	GatewayReasonPending                  = gatewayv1beta1.GatewayReasonPending
	GatewayReasonProgrammed               = gatewayv1beta1.GatewayReasonProgrammed
	HTTPMethodDelete                      = gatewayv1beta1.HTTPMethodDelete
	HTTPMethodGet                         = gatewayv1beta1.HTTPMethodGet
	HTTPProtocolType                      = gatewayv1beta1.HTTPProtocolType
	HTTPRouteFilterExtensionRef           = gatewayv1beta1.HTTPRouteFilterExtensionRef
	HTTPRouteFilterRequestHeaderModifier  = gatewayv1beta1.HTTPRouteFilterRequestHeaderModifier
	HTTPRouteFilterRequestMirror          = gatewayv1beta1.HTTPRouteFilterRequestMirror
	HTTPRouteFilterRequestRedirect        = gatewayv1beta1.HTTPRouteFilterRequestRedirect
	HTTPRouteFilterResponseHeaderModifier = gatewayv1beta1.HTTPRouteFilterResponseHeaderModifier
	HTTPRouteFilterURLRewrite             = gatewayv1beta1.HTTPRouteFilterURLRewrite
	HTTPSProtocolType                     = gatewayv1beta1.HTTPSProtocolType
	HeaderMatchExact                      = gatewayv1beta1.HeaderMatchExact
	HeaderMatchRegularExpression          = gatewayv1beta1.HeaderMatchRegularExpression
	HostnameAddressType                   = gatewayv1beta1.HostnameAddressType
	IPAddressType                         = gatewayv1beta1.IPAddressType
	ListenerConditionAccepted             = gatewayv1beta1.ListenerConditionAccepted
	ListenerConditionConflicted           = gatewayv1beta1.ListenerConditionConflicted
	ListenerConditionProgrammed           = gatewayv1beta1.ListenerConditionProgrammed
	ListenerConditionResolvedRefs         = gatewayv1beta1.ListenerConditionResolvedRefs
	ListenerReasonAccepted                = gatewayv1beta1.ListenerReasonAccepted
	ListenerReasonHostnameConflict        = gatewayv1beta1.ListenerReasonHostnameConflict
	ListenerReasonInvalid                 = gatewayv1beta1.ListenerReasonInvalid
	ListenerReasonInvalidCertificateRef   = gatewayv1beta1.ListenerReasonInvalidCertificateRef
	ListenerReasonInvalidRouteKinds       = gatewayv1beta1.ListenerReasonInvalidRouteKinds
	ListenerReasonNoConflicts             = gatewayv1beta1.ListenerReasonNoConflicts
	ListenerReasonPortUnavailable         = gatewayv1beta1.ListenerReasonPortUnavailable
	ListenerReasonProgrammed              = gatewayv1beta1.ListenerReasonProgrammed
	ListenerReasonProtocolConflict        = gatewayv1beta1.ListenerReasonProtocolConflict
	ListenerReasonRefNotPermitted         = gatewayv1beta1.ListenerReasonRefNotPermitted
	ListenerReasonResolvedRefs            = gatewayv1beta1.ListenerReasonResolvedRefs
	ListenerReasonUnsupportedProtocol     = gatewayv1beta1.ListenerReasonUnsupportedProtocol
	NamespacesFromAll                     = gatewayv1beta1.NamespacesFromAll
	NamespacesFromSame                    = gatewayv1beta1.NamespacesFromSame
	NamespacesFromSelector                = gatewayv1beta1.NamespacesFromSelector
	PathMatchExact                        = gatewayv1beta1.PathMatchExact
	PathMatchPathPrefix                   = gatewayv1beta1.PathMatchPathPrefix
	PathMatchRegularExpression            = gatewayv1beta1.PathMatchRegularExpression
	RouteConditionAccepted                = gatewayv1beta1.RouteConditionAccepted
	RouteConditionResolvedRefs            = gatewayv1beta1.RouteConditionResolvedRefs
	RouteReasonAccepted                   = gatewayv1beta1.RouteReasonAccepted
	RouteReasonBackendNotFound            = gatewayv1beta1.RouteReasonBackendNotFound
	RouteReasonInvalidKind                = gatewayv1beta1.RouteReasonInvalidKind
	RouteReasonNoMatchingListenerHostname = gatewayv1beta1.RouteReasonNoMatchingListenerHostname
	RouteReasonNoMatchingParent           = gatewayv1beta1.RouteReasonNoMatchingParent
	RouteReasonNotAllowedByListeners      = gatewayv1beta1.RouteReasonNotAllowedByListeners
	RouteReasonRefNotPermitted            = gatewayv1beta1.RouteReasonRefNotPermitted
	RouteReasonResolvedRefs               = gatewayv1beta1.RouteReasonResolvedRefs
	TCPProtocolType                       = gatewayv1beta1.TCPProtocolType
	TLSModePassthrough                    = gatewayv1beta1.TLSModePassthrough
	TLSModeTerminate                      = gatewayv1beta1.TLSModeTerminate
	TLSProtocolType                       = gatewayv1beta1.TLSProtocolType
	UDPProtocolType                       = gatewayv1beta1.UDPProtocolType

	GRPCMethodMatchExact             = gatewayv1alpha2.GRPCMethodMatchExact
	GRPCMethodMatchRegularExpression = gatewayv1alpha2.GRPCMethodMatchRegularExpression
)
