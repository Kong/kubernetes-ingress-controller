package sendconfig

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
)

const (
	// WellKnownInitialHash is the hash of an empty configuration.
	WellKnownInitialHash = "00000000000000000000000000000000"
)

type ConfigurationChangeDetector interface {
	HasConfigurationChanged(
		ctx context.Context,
		oldSHA, newSHA []byte,
		targetConfig *file.Content,
		statusClient StatusClient,
	) (bool, error)
}

type StatusClient interface {
	Status(context.Context) (*kong.Status, error)
}

// KongGatewayConfigurationChangeDetector detects changes in Kong Gateway configuration. Besides comparing SHA hashes,
// it also checks if the Kong Gateway has no configuration yet or if the configuration to be pushed is empty.
type KongGatewayConfigurationChangeDetector struct {
	logger logr.Logger
}

func NewKongGatewayConfigurationChangeDetector(logger logr.Logger) *KongGatewayConfigurationChangeDetector {
	return &KongGatewayConfigurationChangeDetector{logger: logger}
}

// HasConfigurationChanged verifies whether configuration has changed by comparing
// old and new config's SHAs. In case the SHAs are equal, it still can return true if a client is considered  crashed or
// just booted up based on its status. In case the status indicates an empty config and the desired config is also empty
// this will return false to prevent continuously sending empty configuration to Gateway.
func (d *KongGatewayConfigurationChangeDetector) HasConfigurationChanged(
	ctx context.Context,
	oldSHA, newSHA []byte,
	targetConfig *file.Content,
	statusClient StatusClient,
) (bool, error) {
	if !bytes.Equal(oldSHA, newSHA) {
		return true, nil
	}

	// Check if a Kong instance has no configuration yet (could mean it crashed, was rebooted, etc.).
	hasNoConfiguration, err := kongHasNoConfiguration(ctx, statusClient)
	if err != nil {
		return false, fmt.Errorf("failed to verify kong readiness: %w", err)
	}

	// Kong instance has no configuration, we should push despite the oldSHA and newSHA being equal...
	if hasNoConfiguration {
		// ... unless we're trying to push an empty config in such case skip.
		if deckgen.IsContentEmpty(targetConfig) {
			return false, nil
		}

		return true, nil
	}

	return false, nil
}

// kongHasNoConfiguration checks Kong's status endpoint and read its config hash. If the config hash reported by Kong is
// the known empty hash, it's considered crashed. This allows providing configuration to Kong instances that have
// unexpectedly crashed and lost their configuration.
func kongHasNoConfiguration(ctx context.Context, client StatusClient) (bool, error) {
	status, err := client.Status(ctx)
	if err != nil {
		return false, err
	}

	if hasNoConfig := status.ConfigurationHash == WellKnownInitialHash; hasNoConfig {
		return true, nil
	}

	return false, nil
}

// KonnectConfigurationChangeDetector detects changes in Konnect configuration by comparing SHA hashes only.
type KonnectConfigurationChangeDetector struct{}

func NewKonnectConfigurationChangeDetector() *KonnectConfigurationChangeDetector {
	return &KonnectConfigurationChangeDetector{}
}

// HasConfigurationChanged verifies whether configuration has changed by comparing old and new config's SHAs.
func (d *KonnectConfigurationChangeDetector) HasConfigurationChanged(
	_ context.Context,
	oldSHA, newSHA []byte,
	_ *file.Content,
	_ StatusClient,
) (bool, error) {
	return !bytes.Equal(oldSHA, newSHA), nil
}
