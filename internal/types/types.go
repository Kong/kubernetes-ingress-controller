package types

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

type HostnameT interface {
	gatewayv1beta1.Hostname | gatewayv1alpha2.Hostname | string
}

type ListenerT interface {
	gatewayv1beta1.Listener | gatewayv1alpha2.Listener
}

type RouteT interface {
	client.Object

	*gatewayv1beta1.HTTPRoute |
		*gatewayv1alpha2.UDPRoute |
		*gatewayv1alpha2.TCPRoute |
		*gatewayv1alpha2.TLSRoute
}

type ParentReferenceT interface {
	gatewayv1alpha2.ParentReference | gatewayv1beta1.ParentReference
}

type BackendRefT interface {
	gatewayv1alpha2.BackendRef |
		gatewayv1beta1.BackendRef |
		gatewayv1alpha2.SecretObjectReference |
		gatewayv1beta1.SecretObjectReference
}

type ParentStatusT interface {
	*gatewayv1alpha2.RouteParentStatus | *gatewayv1beta1.RouteParentStatus
}
