package testlabels

const (
	// Kind is the label key used to store the primary kind that's being tested.
	Kind         = "kind"
	KindUDPRoute = "UDPRoute"
)

const (
	// NetworkingFamily is the label key used to store the networking family of
	// resources that are being tests.
	//
	// Possible, values: "gatewaypi", "ingress".
	NetworkingFamily           = "networkingfamily"
	NetworkingFamilyGatewayAPI = "gatewayapi"
	NetworkingFamilyIngress    = "ingress"
)

const (
	// GatewayAPISupportLevel is the label key used to store the support level
	// of the Gateway API resources being tested.
	GatewayAPISupportLevel      = "gatewayapisupportlevel"
	GatewayAPISupportLevelAlpha = "alpha"
	GatewayAPISupportLevelBeta  = "beta"
)

const (
	// Example is the label key used to indicate whether the test is testing
	// example manifests.
	Example     = "example"
	ExampleTrue = "true"
)
