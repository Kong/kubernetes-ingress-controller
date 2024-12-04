package main

// supportedTypes is a list of types that the cache stores support.
// To add a new type support, add a new entry to this list.
var supportedTypes = []cacheStoreSupportedType{
	// Core Kubernetes types
	{
		Type:       "Ingress",
		Package:    "netv1",
		StoreField: "IngressV1",
	},
	{
		Type:       "IngressClass",
		Package:    "netv1",
		StoreField: "IngressClassV1",
		KeyFunc:    clusterWideKeyFunc,
	},
	{
		Type:    "Service",
		Package: "corev1",
	},
	{
		Type:    "Secret",
		Package: "corev1",
	},
	{
		Type:    "EndpointSlice",
		Package: "discoveryv1",
	},
	// Gateway API types
	{
		Type:    "HTTPRoute",
		Package: "gatewayapi",
	},
	{
		Type:    "UDPRoute",
		Package: "gatewayapi",
	},
	{
		Type:    "TCPRoute",
		Package: "gatewayapi",
	},
	{
		Type:    "TLSRoute",
		Package: "gatewayapi",
	},
	{
		Type:    "GRPCRoute",
		Package: "gatewayapi",
	},
	{
		Type:    "ReferenceGrant",
		Package: "gatewayapi",
	},
	{
		Type:    "Gateway",
		Package: "gatewayapi",
	},
	{
		Type:    "BackendTLSPolicy",
		Package: "gatewayapi",
	},
	// Kong types
	{
		Type:       "KongPlugin",
		Package:    "kongv1",
		StoreField: "Plugin",
	},
	{
		Type:       "KongClusterPlugin",
		Package:    "kongv1",
		StoreField: "ClusterPlugin",
		KeyFunc:    clusterWideKeyFunc,
	},
	{
		Type:       "KongConsumer",
		Package:    "kongv1",
		StoreField: "Consumer",
	},
	{
		Type:       "KongConsumerGroup",
		Package:    "kongv1beta1",
		StoreField: "ConsumerGroup",
	},
	{
		Type:    "KongIngress",
		Package: "kongv1",
	},
	{
		Type:    "TCPIngress",
		Package: "kongv1beta1",
	},
	{
		Type:    "UDPIngress",
		Package: "kongv1beta1",
	},
	{
		Type:    "KongUpstreamPolicy",
		Package: "kongv1beta1",
	},
	{
		Type:       "IngressClassParameters",
		Package:    "kongv1alpha1",
		StoreField: "IngressClassParametersV1alpha1",
	},
	{
		Type:    "KongServiceFacade",
		Package: "incubatorv1alpha1",
	},
	{
		Type:    "KongVault",
		Package: "kongv1alpha1",
		KeyFunc: clusterWideKeyFunc,
	},
	{
		Type:    "KongCustomEntity",
		Package: "kongv1alpha1",
	},
}
