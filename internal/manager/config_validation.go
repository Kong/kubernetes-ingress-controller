package manager

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	"github.com/samber/mo"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	cfgtypes "github.com/kong/kubernetes-ingress-controller/v3/internal/manager/config/types"
)

// https://github.com/kubernetes-sigs/gateway-api/blob/547122f7f55ac0464685552898c560658fb40073/apis/v1beta1/shared_types.go#L448-L463
var gatewayAPIControllerNameRegex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\/[A-Za-z0-9\/\-._~%!$&'()*+,;=:]+$`)

// *FromFlagValue functions are used to validate single flag values and set those in Config.
// They're meant to be used together with ValidatedValue[T] type.

func namespacedNameFromFlagValue(flagValue string) (OptionalNamespacedName, error) {
	parts := strings.SplitN(flagValue, "/", 3)
	if len(parts) != 2 {
		return OptionalNamespacedName{}, errors.New("the expected format is namespace/name")
	}
	if parts[0] == "" {
		return OptionalNamespacedName{}, errors.New("namespace cannot be empty")
	}
	if parts[1] == "" {
		return OptionalNamespacedName{}, errors.New("name cannot be empty")
	}

	return mo.Some(k8stypes.NamespacedName{
		Namespace: parts[0],
		Name:      parts[1],
	}), nil
}

func gatewayAPIControllerNameFromFlagValue(flagValue string) (string, error) {
	if !gatewayAPIControllerNameRegex.MatchString(flagValue) {
		return "", errors.New("the expected format is example.com/controller-name")
	}
	return flagValue, nil
}

func dnsStrategyFromFlagValue(flagValue string) (cfgtypes.DNSStrategy, error) {
	strategy := cfgtypes.DNSStrategy(flagValue)
	if err := strategy.Validate(); err != nil {
		return cfgtypes.DNSStrategy(""), err
	}
	return strategy, nil
}

// Validate validates the config. It should be used to validate the config variables' interdependencies.
// When a single variable is to be validated, *FromFlagValue function should be implemented.
func (c *Config) Validate() error {
	if c.flagSet != nil {
		if c.flagSet.Changed("kong-admin-svc") && c.flagSet.Changed("kong-admin-url") {
			return fmt.Errorf("can't set both --kong-admin-svc and --kong-admin-url")
		}
	}
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

	return nil
}

func (c *Config) validateKonnect() error {
	konnect := c.Konnect
	if !konnect.ConfigSynchronizationEnabled {
		return nil
	}

	if c.KongAdminSvc.IsAbsent() {
		return errors.New("--kong-admin-svc has to be set when using --konnect-sync-enabled")
	}
	if konnect.Address == "" {
		return errors.New("address not specified")
	}
	if konnect.ControlPlaneID == "" {
		return errors.New("control plane not specified")
	}
	if konnect.TLSClient.IsZero() {
		return fmt.Errorf("missing TLS client configuration")
	}
	if err := validateClientTLS(konnect.TLSClient); err != nil {
		return fmt.Errorf("TLS client config invalid: %w", err)
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
	if !c.FeatureGates[featuregates.FallbackConfiguration] && c.UseLastValidConfigForFallback {
		return fmt.Errorf(
			"--use-last-valid-config-for-fallback or CONTROLLER_USE_LAST_VALID_CONFIG_FOR_FALLBACK can only be used with %s feature gate enabled",
			featuregates.FallbackConfiguration,
		)
	}
	return nil
}

func validateClientTLS(clientTLS adminapi.TLSClientConfig) error {
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
