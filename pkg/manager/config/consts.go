package config

import "time"

const (
	// LeaderElectionEnabled is a constant that represents a value that should be used to enable leader election.
	LeaderElectionEnabled = "enabled"
	// LeaderElectionDisabled is a constant that represents a value that should be used to disable leader election.
	LeaderElectionDisabled = "disabled"
)

const (
	// DefaultDataPlanesReadinessReconciliationInterval is the interval at which the manager will run DPs readiness reconciliation loop.
	// It's the same as the default interval of a Kubernetes container's readiness probe.
	DefaultDataPlanesReadinessReconciliationInterval = 10 * time.Second
	// MinDataPlanesReadinessReconciliationInterval is the minimum interval of DPs readiness reconciliation loop.
	MinDataPlanesReadinessReconciliationInterval = 3 * time.Second
	// DefaultDataPlanesReadinessCheckTimeout is the default timeout of readiness check.
	// When a readiness check request did not get response within the timeout, the gateway instance will turn into `Pending` status.
	DefaultDataPlanesReadinessCheckTimeout = 5 * time.Second
)

const (
	// MinKonnectConfigUploadPeriod is the minimum period between operations to upload Kong configuration to Konnect.
	MinKonnectConfigUploadPeriod = 10 * time.Second
	// DefaultKonnectConfigUploadPeriod is the default period between operations to upload Kong configuration to Konnect.
	DefaultKonnectConfigUploadPeriod = 30 * time.Second
)
