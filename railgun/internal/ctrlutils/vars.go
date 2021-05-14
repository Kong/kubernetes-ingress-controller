package ctrlutils

// -----------------------------------------------------------------------------
// General Controller Variables
// -----------------------------------------------------------------------------

var (
	// KongIngressFinalizer is the finalizer used to ensure Kong configuration cleanup for deleted resources.
	KongIngressFinalizer = "configuration.konghq.com/ingress"
)
