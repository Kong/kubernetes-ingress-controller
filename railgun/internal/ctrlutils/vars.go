package ctrlutils

// -----------------------------------------------------------------------------
// General Controller Variables
// -----------------------------------------------------------------------------

var (
	// DefaultNamespace indicates the namespace that will be used by default
	// when no other is provided for the deployment or management of resources.
	DefaultNamespace = "kong-system"

	// KongIngressFinalizer is the finalizer used to ensure Kong configuration cleanup for deleted resources.
	KongIngressFinalizer = "configuration.konghq.com/ingress"
)
