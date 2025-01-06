package testlabels

const (
	// Kind is the label key used to store the primary kind that's being tested.
	Kind                   = "kind"
	KindUDPRoute           = "UDPRoute"
	KindTCPRoute           = "TCPRoute"
	KindHTTPRoute          = "HTTPRoute"
	KindGRPCRoute          = "GRPCRoute"
	KindBackendTLSPolicy   = "BackendTLSPolicy"
	KindIngress            = "Ingress"
	KindKongServiceFacade  = "KongServiceFacade"
	KindKongUDPIngress     = "UDPIngress"
	KindKongTCPIngress     = "TCPIngress"
	KindKongUpstreamPolicy = "KongUpstreamPolicy"
	KindKongLicense        = "KongLicense"
	KindKongCustomEntity   = "KongCustomEntity"
)

const (
	// NetworkingFamily is the label key used to store the networking family of
	// resources that are being tests.
	//
	// Possible, values: "gatewayapi", "ingress".
	NetworkingFamily           = "networkingfamily"
	NetworkingFamilyGatewayAPI = "gatewayapi"
	NetworkingFamilyIngress    = "ingress"
)

const (
	// Example is the label key used to indicate whether the test is testing
	// example manifests.
	Example     = "example"
	ExampleTrue = "true"
)
