package gatewayapi

import (
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	// TLSVerifyDepthKey is the key used to store the tls verify depth.
	// This is used in the BackendTLSPolicy options.
	TLSVerifyDepthKey AnnotationKey = "tls-verify-depth"
)

const (
	V1Group = Group(gatewayv1.GroupName)
)

var V1GroupVersion = gatewayv1.GroupVersion.Version
