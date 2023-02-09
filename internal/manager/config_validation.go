package manager

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

// *FromFlagValue functions are used to validate single flag values and set those in Config.
// They're meant to be used together with ValidatedValue[T] type.

func namespacedNameFromFlagValue(flagValue string) (types.NamespacedName, error) {
	parts := strings.SplitN(flagValue, "/", 3)
	if len(parts) != 2 {
		return types.NamespacedName{}, errors.New("the expected format is namespace/name")
	}
	return types.NamespacedName{
		Namespace: parts[0],
		Name:      parts[1],
	}, nil
}

func gatewayAPIControllerNameFromFlagValue(flagValue string) (string, error) {
	// https://github.com/kubernetes-sigs/gateway-api/blob/547122f7f55ac0464685552898c560658fb40073/apis/v1beta1/shared_types.go#L448-L463
	re := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\/[A-Za-z0-9\/\-._~%!$&'()*+,;=:]+$`)
	if !re.Match([]byte(flagValue)) {
		return "", errors.New("the expected format is example.com/controller-name")
	}
	return flagValue, nil
}

// Validate validates the config. It should be used to validate the config variables' interdependencies.
// When a single variable is to be validated, *FromFlagValue function should be implemented.
func (c *Config) Validate() error {
	if c.flagSet != nil {
		if c.flagSet.Changed("kong-admin-svc") && c.flagSet.Changed("kong-admin-url") {
			return fmt.Errorf("can't set both --kong-admin-svc and --kong-admin-url")
		}
	}

	if err := c.validateKonnect(); err != nil {
		return fmt.Errorf("invalid konnect configuration: %w", err)
	}
	if err := c.validateKongAdminAPI(); err != nil {
		return fmt.Errorf("invalid kong admin api configuration: %w", err)
	}

	return nil
}

func (c *Config) validateKonnect() error {
	konnect := c.Konnect
	if !konnect.ConfigSynchronizationEnabled {
		return nil
	}

	if konnect.Address == "" {
		return errors.New("address not specified")
	}
	if konnect.RuntimeGroupID == "" {
		return errors.New("runtime group not specified")
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
