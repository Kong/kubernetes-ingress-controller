package ctrlutils

// -----------------------------------------------------------------------------
// General Controller Variables
// -----------------------------------------------------------------------------

// KongIngressFinalizer is the finalizer used to ensure Kong configuration cleanup for deleted resources.
const KongIngressFinalizer = "configuration.konghq.com/ingress"
const KnativeIngressFinalizer = "networking.internal.knative.dev/ingress"
