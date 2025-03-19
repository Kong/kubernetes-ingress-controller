package config

import (
	"errors"
	"fmt"
)

// Validate validates the config. It should be used to validate the config variables' interdependencies.
func (c *Config) Validate() error {
	if c.KongAdminToken != "" && c.KongAdminTokenPath != "" {
		return errors.New("both admin token and admin token file specified, only one allowed")
	}

	if err := c.validateKonnect(); err != nil {
		return fmt.Errorf("invalid konnect configuration: %w", err)
	}
	if err := c.validateKongAdminAPI(); err != nil {
		return fmt.Errorf("invalid kong admin api configuration: %w", err)
	}
	if err := c.validateFallbackConfiguration(); err != nil {
		return fmt.Errorf("invalid fallback config settings: %w", err)
	}
	if err := c.validateGatewayDiscovery(); err != nil {
		return fmt.Errorf("invalid gateway discovery configuration: %w", err)
	}

	return nil
}

func (c *Config) validateKonnect() error {
	if !c.Konnect.ConfigSynchronizationEnabled {
		return nil
	}

	if c.KongAdminSvc.IsAbsent() {
		return errors.New("--kong-admin-svc has to be set when using --konnect-sync-enabled")
	}
	if c.Konnect.Address == "" {
		return errors.New("address not specified")
	}
	if c.Konnect.ControlPlaneID == "" {
		return errors.New("control plane not specified")
	}
	if c.Konnect.TLSClient.IsZero() {
		return fmt.Errorf("missing TLS client configuration")
	}
	if err := validateClientTLS(c.Konnect.TLSClient); err != nil {
		return fmt.Errorf("TLS client config invalid: %w", err)
	}
	if c.Konnect.UploadConfigPeriod < MinKonnectConfigUploadPeriod {
		return fmt.Errorf("cannot set upload config period to be smaller than %s", MinKonnectConfigUploadPeriod.String())
	}
	return nil
}

func (c *Config) validateKongAdminAPI() error {
	if err := validateClientTLS(c.KongAdminAPIConfig.TLSClient); err != nil {
		return fmt.Errorf("TLS client config invalid: %w", err)
	}
	return nil
}

func (c *Config) validateFallbackConfiguration() error {
	if !c.FeatureGates[FallbackConfigurationFeature] && c.UseLastValidConfigForFallback {
		return fmt.Errorf(
			"--use-last-valid-config-for-fallback or CONTROLLER_USE_LAST_VALID_CONFIG_FOR_FALLBACK can only be used with %s feature gate enabled",
			FallbackConfigurationFeature,
		)
	}
	return nil
}

func (c *Config) validateGatewayDiscovery() error {
	// Skip validation if gateway discovery is not enabled.
	if _, ok := c.KongAdminSvc.Get(); !ok {
		return nil
	}

	if c.GatewayDiscoveryReadinessCheckInterval < MinDataPlanesReadinessReconciliationInterval {
		return fmt.Errorf("Readiness check reconciliation interval cannot be less than %s",
			MinDataPlanesReadinessReconciliationInterval)
	}
	if c.GatewayDiscoveryReadinessCheckTimeout >= c.GatewayDiscoveryReadinessCheckInterval {
		return fmt.Errorf("Readiness check timeout must be less than readiness check recociliation interval")
	}
	return nil
}

func validateClientTLS(clientTLS TLSClientConfig) error {
	if clientTLS.Cert != "" && clientTLS.CertFile != "" {
		return errors.New("both client certificate and client certificate file specified, only one allowed")
	}
	if clientTLS.Key != "" && clientTLS.KeyFile != "" {
		return errors.New("both client key and client key file specified, only one allowed")
	}

	clientCertProvided := clientTLS.Cert != "" || clientTLS.CertFile != ""
	clientKeyProvided := clientTLS.Key != "" || clientTLS.KeyFile != ""

	if clientCertProvided && !clientKeyProvided {
		return errors.New("client certificate was provided, but the client key was not")
	}

	if clientKeyProvided && !clientCertProvided {
		return errors.New("client key was provided, but the client certificate was not")
	}

	return nil
}
