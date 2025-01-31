package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/samber/mo"
	k8stypes "k8s.io/apimachinery/pkg/types"

	cfgtypes "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// https://github.com/kubernetes-sigs/gateway-api/blob/547122f7f55ac0464685552898c560658fb40073/apis/v1beta1/shared_types.go#L448-L463
var gatewayAPIControllerNameRegex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\/[A-Za-z0-9\/\-._~%!$&'()*+,;=:]+$`)

// Validate validates the config. It should be used to validate the config variables' interdependencies.
// When a single variable is to be validated, *FromFlagValue function should be implemented.
func (c *CLIConfig) Validate() error {
	if c.flagSet != nil {
		if c.flagSet.Changed("kong-admin-svc") && c.flagSet.Changed("kong-admin-url") {
			return fmt.Errorf("can't set both --kong-admin-svc and --kong-admin-url")
		}
	}
	return c.Config.Validate()
}

// *FromFlagValue functions are used to validate single flag values and set those in CLIConfig.
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

func metricsAccessFilterFromFlagValue(flagValue string) (cfgtypes.MetricsAccessFilter, error) {
	switch flagValue {
	case string(cfgtypes.MetricsAccessFilterOff), string(cfgtypes.MetricsAccessFilterRBAC):
		return cfgtypes.MetricsAccessFilter(flagValue), nil
	default:
		return "", fmt.Errorf("unsupported metrics filter %s", flagValue)
	}
}

func gatewayAPIControllerNameFromFlagValue(flagValue string) (string, error) {
	if !gatewayAPIControllerNameRegex.MatchString(flagValue) {
		return "", errors.New("the expected format is example.com/controller-name")
	}
	return flagValue, nil
}
