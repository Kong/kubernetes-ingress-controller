package controllers

// -----------------------------------------------------------------------------
// Environment Variables
// -----------------------------------------------------------------------------

var (
	// CtrlNamespaceEnv provides the name of the environment variable where controllers which
	// manage the Kong configuration should use to find (or create) the configuration secret.
	CtrlNamespaceEnv = "KONG_CONFIGURATION_NAMESPACE"

	// ExternalCtrlEnv is an environment variable used to indicate whether the controller is running
	// outside the cluster where the proxy instances are running. If unset it is assumed that the proxy
	// instances can be reached via their Pod IP address. Only accepts "true" to enable.
	ExternalCtrlEnv = "KONG_EXTERNAL_CONTROLLER"
)

// -----------------------------------------------------------------------------
// General Controller Variables
// -----------------------------------------------------------------------------

var (
	// ProxyInstanceLabel is a label used for controllers (such as the secret configuration
	// controller) to identify which pods are running the Kong proxy which needs to be configured.
	ProxyInstanceLabel = "konghq.com/proxy-instance"
)
