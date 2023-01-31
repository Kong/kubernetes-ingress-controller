package consts

const (
	// KongTestPassword is used as a password only within the context of transient integration test runs
	// and is left static to help developers debug failures in those testing environments.
	KongTestPassword = "password"

	// IngressClass indicates the ingress class name which the tests will use for supported object reconciliation.
	IngressClass = "kongtests"

	// ControllerNamespace is the Kubernetes namespace where the controller is deployed.
	ControllerNamespace = "kong-system"
)
