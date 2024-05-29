package gatewayapi

import (
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

var (
	InstallV1    = gatewayv1.Install
	GroupVersion = gatewayv1.GroupVersion
)

// This file contains aliases for types and consts from the Gateway API.  Its purpose is to allow easy migration from
// one version of the Gateway API to another with minimal changes to the codebase.

type (
	AllowedRoutes             = gatewayv1.AllowedRoutes
	BackendObjectReference    = gatewayv1.BackendObjectReference
	BackendRef                = gatewayv1.BackendRef
	CommonRouteSpec           = gatewayv1.CommonRouteSpec
	Duration                  = gatewayv1.Duration
	Gateway                   = gatewayv1.Gateway
	GatewayAddress            = gatewayv1.GatewayAddress
	GatewayClass              = gatewayv1.GatewayClass
	GatewayClassList          = gatewayv1.GatewayClassList
	GatewayClassSpec          = gatewayv1.GatewayClassSpec
	GatewayClassStatus        = gatewayv1.GatewayClassStatus
	GatewayController         = gatewayv1.GatewayController
	GatewayList               = gatewayv1.GatewayList
	GatewaySpec               = gatewayv1.GatewaySpec
	GatewayStatus             = gatewayv1.GatewayStatus
	GatewayStatusAddress      = gatewayv1.GatewayStatusAddress
	GatewayTLSConfig          = gatewayv1.GatewayTLSConfig
	Group                     = gatewayv1.Group
	HTTPBackendRef            = gatewayv1.HTTPBackendRef
	HTTPHeader                = gatewayv1.HTTPHeader
	HTTPHeaderFilter          = gatewayv1.HTTPHeaderFilter
	HTTPHeaderMatch           = gatewayv1.HTTPHeaderMatch
	HTTPHeaderName            = gatewayv1.HTTPHeaderName
	HTTPMethod                = gatewayv1.HTTPMethod
	HTTPPathMatch             = gatewayv1.HTTPPathMatch
	HTTPQueryParamMatch       = gatewayv1.HTTPQueryParamMatch
	HTTPRequestRedirectFilter = gatewayv1.HTTPRequestRedirectFilter
	HTTPRoute                 = gatewayv1.HTTPRoute
	HTTPRouteFilter           = gatewayv1.HTTPRouteFilter
	HTTPRouteFilterType       = gatewayv1.HTTPRouteFilterType
	HTTPURLRewriteFilter      = gatewayv1.HTTPURLRewriteFilter
	HTTPPathModifier          = gatewayv1.HTTPPathModifier
	HTTPRouteList             = gatewayv1.HTTPRouteList
	HTTPRouteMatch            = gatewayv1.HTTPRouteMatch
	HTTPRouteRule             = gatewayv1.HTTPRouteRule
	HTTPRouteTimeouts         = gatewayv1.HTTPRouteTimeouts
	HTTPRequestMirrorFilter   = gatewayv1.HTTPRequestMirrorFilter
	LocalObjectReference      = gatewayv1.LocalObjectReference
	HTTPRouteSpec             = gatewayv1.HTTPRouteSpec
	HTTPRouteStatus           = gatewayv1.HTTPRouteStatus
	Hostname                  = gatewayv1.Hostname
	Kind                      = gatewayv1.Kind
	Listener                  = gatewayv1.Listener
	ListenerConditionReason   = gatewayv1.ListenerConditionReason
	ListenerConditionType     = gatewayv1.ListenerConditionType
	ListenerStatus            = gatewayv1.ListenerStatus
	Namespace                 = gatewayv1.Namespace
	ObjectName                = gatewayv1.ObjectName
	ParentReference           = gatewayv1.ParentReference
	PathMatchType             = gatewayv1.PathMatchType
	PortNumber                = gatewayv1.PortNumber
	PreciseHostname           = gatewayv1.PreciseHostname
	ProtocolType              = gatewayv1.ProtocolType
	ReferenceGrant            = gatewayv1beta1.ReferenceGrant
	ReferenceGrantFrom        = gatewayv1beta1.ReferenceGrantFrom
	ReferenceGrantList        = gatewayv1beta1.ReferenceGrantList
	ReferenceGrantSpec        = gatewayv1beta1.ReferenceGrantSpec
	ReferenceGrantTo          = gatewayv1beta1.ReferenceGrantTo
	RouteConditionReason      = gatewayv1.RouteConditionReason
	RouteGroupKind            = gatewayv1.RouteGroupKind
	RouteNamespaces           = gatewayv1.RouteNamespaces
	RouteParentStatus         = gatewayv1.RouteParentStatus
	RouteStatus               = gatewayv1.RouteStatus
	SecretObjectReference     = gatewayv1.SecretObjectReference
	SectionName               = gatewayv1.SectionName
	GRPCBackendRef            = gatewayv1.GRPCBackendRef
	GRPCHeaderMatch           = gatewayv1.GRPCHeaderMatch
	GRPCHeaderName            = gatewayv1.GRPCHeaderName
	GRPCMethodMatch           = gatewayv1.GRPCMethodMatch
	GRPCMethodMatchType       = gatewayv1.GRPCMethodMatchType
	GRPCRoute                 = gatewayv1.GRPCRoute
	GRPCRouteList             = gatewayv1.GRPCRouteList
	GRPCRouteMatch            = gatewayv1.GRPCRouteMatch
	GRPCRouteRule             = gatewayv1.GRPCRouteRule
	GRPCRouteSpec             = gatewayv1.GRPCRouteSpec
	GRPCRouteStatus           = gatewayv1.GRPCRouteStatus

	PolicyAncestorStatus = gatewayv1alpha2.PolicyAncestorStatus
	PolicyStatus         = gatewayv1alpha2.PolicyStatus
	TCPRoute             = gatewayv1alpha2.TCPRoute
	TCPRouteList         = gatewayv1alpha2.TCPRouteList
	TCPRouteRule         = gatewayv1alpha2.TCPRouteRule
	TCPRouteSpec         = gatewayv1alpha2.TCPRouteSpec
	TCPRouteStatus       = gatewayv1alpha2.TCPRouteStatus
	TLSRoute             = gatewayv1alpha2.TLSRoute
	TLSRouteList         = gatewayv1alpha2.TLSRouteList
	TLSRouteRule         = gatewayv1alpha2.TLSRouteRule
	TLSRouteSpec         = gatewayv1alpha2.TLSRouteSpec
	TLSRouteStatus       = gatewayv1alpha2.TLSRouteStatus
	UDPRoute             = gatewayv1alpha2.UDPRoute
	UDPRouteList         = gatewayv1alpha2.UDPRouteList
	UDPRouteRule         = gatewayv1alpha2.UDPRouteRule
	UDPRouteSpec         = gatewayv1alpha2.UDPRouteSpec
	UDPRouteStatus       = gatewayv1alpha2.UDPRouteStatus
)

const (
	FullPathHTTPPathModifier              = gatewayv1.FullPathHTTPPathModifier
	PrefixMatchHTTPPathModifier           = gatewayv1.PrefixMatchHTTPPathModifier
	GatewayClassConditionStatusAccepted   = gatewayv1.GatewayClassConditionStatusAccepted
	GatewayClassReasonAccepted            = gatewayv1.GatewayClassReasonAccepted
	GatewayConditionAccepted              = gatewayv1.GatewayConditionAccepted
	GatewayConditionProgrammed            = gatewayv1.GatewayConditionProgrammed
	GatewayReasonAccepted                 = gatewayv1.GatewayReasonAccepted
	GatewayReasonPending                  = gatewayv1.GatewayReasonPending
	GatewayReasonProgrammed               = gatewayv1.GatewayReasonProgrammed
	HTTPMethodDelete                      = gatewayv1.HTTPMethodDelete
	HTTPMethodGet                         = gatewayv1.HTTPMethodGet
	HTTPProtocolType                      = gatewayv1.HTTPProtocolType
	HTTPRouteFilterExtensionRef           = gatewayv1.HTTPRouteFilterExtensionRef
	HTTPRouteFilterRequestHeaderModifier  = gatewayv1.HTTPRouteFilterRequestHeaderModifier
	HTTPRouteFilterRequestMirror          = gatewayv1.HTTPRouteFilterRequestMirror
	HTTPRouteFilterRequestRedirect        = gatewayv1.HTTPRouteFilterRequestRedirect
	HTTPRouteFilterResponseHeaderModifier = gatewayv1.HTTPRouteFilterResponseHeaderModifier
	HTTPRouteFilterURLRewrite             = gatewayv1.HTTPRouteFilterURLRewrite
	HTTPSProtocolType                     = gatewayv1.HTTPSProtocolType
	HeaderMatchExact                      = gatewayv1.HeaderMatchExact
	HeaderMatchRegularExpression          = gatewayv1.HeaderMatchRegularExpression
	HostnameAddressType                   = gatewayv1.HostnameAddressType
	IPAddressType                         = gatewayv1.IPAddressType
	ListenerConditionAccepted             = gatewayv1.ListenerConditionAccepted
	ListenerConditionConflicted           = gatewayv1.ListenerConditionConflicted
	ListenerConditionProgrammed           = gatewayv1.ListenerConditionProgrammed
	ListenerConditionResolvedRefs         = gatewayv1.ListenerConditionResolvedRefs
	ListenerReasonAccepted                = gatewayv1.ListenerReasonAccepted
	ListenerReasonHostnameConflict        = gatewayv1.ListenerReasonHostnameConflict
	ListenerReasonInvalid                 = gatewayv1.ListenerReasonInvalid
	ListenerReasonInvalidCertificateRef   = gatewayv1.ListenerReasonInvalidCertificateRef
	ListenerReasonInvalidRouteKinds       = gatewayv1.ListenerReasonInvalidRouteKinds
	ListenerReasonNoConflicts             = gatewayv1.ListenerReasonNoConflicts
	ListenerReasonPortUnavailable         = gatewayv1.ListenerReasonPortUnavailable
	ListenerReasonProgrammed              = gatewayv1.ListenerReasonProgrammed
	ListenerReasonProtocolConflict        = gatewayv1.ListenerReasonProtocolConflict
	ListenerReasonRefNotPermitted         = gatewayv1.ListenerReasonRefNotPermitted
	ListenerReasonResolvedRefs            = gatewayv1.ListenerReasonResolvedRefs
	ListenerReasonUnsupportedProtocol     = gatewayv1.ListenerReasonUnsupportedProtocol
	NamespacesFromAll                     = gatewayv1.NamespacesFromAll
	NamespacesFromSame                    = gatewayv1.NamespacesFromSame
	NamespacesFromSelector                = gatewayv1.NamespacesFromSelector
	PathMatchExact                        = gatewayv1.PathMatchExact
	PathMatchPathPrefix                   = gatewayv1.PathMatchPathPrefix
	PathMatchRegularExpression            = gatewayv1.PathMatchRegularExpression
	QueryParamMatchExact                  = gatewayv1.QueryParamMatchExact
	QueryParamMatchRegularExpression      = gatewayv1.QueryParamMatchRegularExpression
	RouteConditionAccepted                = gatewayv1.RouteConditionAccepted
	RouteConditionResolvedRefs            = gatewayv1.RouteConditionResolvedRefs
	RouteReasonAccepted                   = gatewayv1.RouteReasonAccepted
	RouteReasonBackendNotFound            = gatewayv1.RouteReasonBackendNotFound
	RouteReasonInvalidKind                = gatewayv1.RouteReasonInvalidKind
	RouteReasonNoMatchingListenerHostname = gatewayv1.RouteReasonNoMatchingListenerHostname
	RouteReasonNoMatchingParent           = gatewayv1.RouteReasonNoMatchingParent
	RouteReasonNotAllowedByListeners      = gatewayv1.RouteReasonNotAllowedByListeners
	RouteReasonRefNotPermitted            = gatewayv1.RouteReasonRefNotPermitted
	RouteReasonResolvedRefs               = gatewayv1.RouteReasonResolvedRefs
	TCPProtocolType                       = gatewayv1.TCPProtocolType
	TLSModePassthrough                    = gatewayv1.TLSModePassthrough
	TLSModeTerminate                      = gatewayv1.TLSModeTerminate
	TLSProtocolType                       = gatewayv1.TLSProtocolType
	UDPProtocolType                       = gatewayv1.UDPProtocolType

	GRPCMethodMatchExact             = gatewayv1.GRPCMethodMatchExact
	GRPCMethodMatchRegularExpression = gatewayv1.GRPCMethodMatchRegularExpression

	PolicyConditionAccepted = gatewayv1alpha2.PolicyConditionAccepted
	PolicyReasonAccepted    = gatewayv1alpha2.PolicyReasonAccepted
	PolicyReasonConflicted  = gatewayv1alpha2.PolicyReasonConflicted
)
