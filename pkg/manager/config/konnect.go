package config

import (
	"time"
)

type KonnectConfig struct {
	// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/3922
	// ConfigSynchronizationEnabled is the only toggle we had prior to the addition of the license agent.
	// We likely want to combine these into a single Konnect toggle or piggyback off other Konnect functionality.
	ConfigSynchronizationEnabled bool
	ControlPlaneID               string
	Address                      string
	UploadConfigPeriod           time.Duration
	ConfigSyncConcurrency        int
	RefreshNodePeriod            time.Duration
	TLSClient                    TLSClientConfig

	LicenseSynchronizationEnabled bool
	InitialLicensePollingPeriod   time.Duration
	LicensePollingPeriod          time.Duration
	LicenseStorageEnabled         bool
	ConsumersSyncDisabled         bool
}
