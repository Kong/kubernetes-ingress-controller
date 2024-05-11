package gatewayapi

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

var V1GatewayTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1.GroupVersion.String(),
	Kind:       "Gateway",
}

var V1beta1GatewayTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1beta1.GroupVersion.String(),
	Kind:       "Gateway",
}

var V1GatewayClassTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1.GroupVersion.String(),
	Kind:       "GatewayClass",
}

var V1beta1GatewayClassTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1beta1.GroupVersion.String(),
	Kind:       "GatewayClass",
}

var V1HTTPRouteTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1.GroupVersion.String(),
	Kind:       "HTTPRoute",
}

var V1beta1HTTPRouteTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1beta1.GroupVersion.String(),
	Kind:       "HTTPRoute",
}

var ReferenceGrantTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1beta1.GroupVersion.String(),
	Kind:       "ReferenceGrant",
}

var GRPCRouteTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1.GroupVersion.String(),
	Kind:       "GRPCRoute",
}

var TCPRouteTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1alpha2.GroupVersion.String(),
	Kind:       "TCPRoute",
}

var TLSRouteTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1alpha2.GroupVersion.String(),
	Kind:       "TLSRoute",
}

var UDPRouteTypeMeta = metav1.TypeMeta{
	APIVersion: gatewayv1alpha2.GroupVersion.String(),
	Kind:       "UDPRoute",
}

var (
	V1GatewayGVResource = metav1.GroupVersionResource{
		Group:    gatewayv1.GroupVersion.Group,
		Version:  gatewayv1.GroupVersion.Version,
		Resource: "gateways",
	}
	V1HTTPRouteGVResource = metav1.GroupVersionResource{
		Group:    gatewayv1.GroupVersion.Group,
		Version:  gatewayv1.GroupVersion.Version,
		Resource: "httproutes",
	}
	V1beta1GatewayGVResource = metav1.GroupVersionResource{
		Group:    gatewayv1beta1.GroupVersion.Group,
		Version:  gatewayv1beta1.GroupVersion.Version,
		Resource: "gateways",
	}
	V1beta1HTTPRouteGVResource = metav1.GroupVersionResource{
		Group:    gatewayv1beta1.GroupVersion.Group,
		Version:  gatewayv1beta1.GroupVersion.Version,
		Resource: "httproutes",
	}
)
