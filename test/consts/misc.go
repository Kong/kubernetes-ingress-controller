package consts

const (
	// KongTestPassword is used as a password only within the context of transient integration test runs
	// and is left static to help developers debug failures in those testing environments.
	KongTestPassword = "password"

	// KongTestWorkspace is used as a workspace only within the context of transient integration test runs
	// when Kong Enterprise is enabled and a database is used (DBmode != off) and is left static to help
	// developers debug failures in those testing environments.
	KongTestWorkspace = "notdefault"

	// IngressClass indicates the ingress class name which the tests will use for supported object reconciliation.
	IngressClass = "kongtests"

	// ControllerNamespace is the Kubernetes namespace where the controller is deployed.
	ControllerNamespace = "kong"
)
