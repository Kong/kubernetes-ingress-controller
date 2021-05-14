package ctrlutils

// -----------------------------------------------------------------------------
// General Controller Variables
// -----------------------------------------------------------------------------

var (
	// DefaultNamespace indicates the namespace that will be used by default
	// when no other is provided for the deployment or management of resources.
	DefaultNamespace = "kong-system"

	// ConfigSecretName indicates the name of the Secret object where Ingress controllers will upload
	// ingress objects for eventual parsing and configuration in the Kong Proxy APIs.
	ConfigSecretName = "kong-config"

	// ProxyInstanceLabel is a label used for controllers (such as the secret configuration
	// controller) to identify which pods are running the Kong proxy which needs to be configured.
	ProxyInstanceLabel = "konghq.com/proxy-instance"

	// KongIngressFinalizer is the finalizer used to ensure Kong configuration cleanup for deleted resources.
	KongIngressFinalizer = "configuration.konghq.com/ingress"
)
