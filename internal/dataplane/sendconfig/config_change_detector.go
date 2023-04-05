package sendconfig

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
)

const (
	// WellKnownInitialHash is the hash of an empty configuration.
	WellKnownInitialHash = "00000000000000000000000000000000"
)

type ConfigurationChangeDetector interface {
	// HasConfigurationChanged verifies whether configuration has changed by comparing
	// old and new config's SHAs.
	// In case the SHAs are equal, it still can return true if a client is considered
	// crashed or just booted up based on its status.
	// In case the status indicates an empty config and the desired config is also empty
	// this will return false to prevent continuously sending empty configuration to Gateway.
	HasConfigurationChanged(
		ctx context.Context,
		oldSHA, newSHA []byte,
		targetConfig *file.Content,
		client KonnectAwareClient,
		statusClient StatusClient,
	) (bool, error)
}

type KonnectAwareClient interface {
	IsKonnect() bool
}

type StatusClient interface {
	Status(context.Context) (*kong.Status, error)
}

type DefaultConfigurationChangeDetector struct {
	log logrus.FieldLogger
}

func NewDefaultClientConfigurationChangeDetector(log logrus.FieldLogger) *DefaultConfigurationChangeDetector {
	return &DefaultConfigurationChangeDetector{log: log}
}

func (d *DefaultConfigurationChangeDetector) HasConfigurationChanged(
	ctx context.Context,
	oldSHA, newSHA []byte,
	targetConfig *file.Content,
	client KonnectAwareClient,
	statusClient StatusClient,
) (bool, error) {
	if !bytes.Equal(oldSHA, newSHA) {
		return true, nil
	}

	// In case of Konnect, we skip further steps that are meant to detect Kong instances crash/reset
	// that are not relevant for Konnect.
	// We're sure that if oldSHA and newSHA are equal, we are safe to skip the update.
	if client.IsKonnect() {
		return false, nil
	}

	// Check if a Kong instance has no configuration yet (could mean it crashed, was rebooted, etc.).
	hasNoConfiguration, err := kongHasNoConfiguration(ctx, statusClient, d.log)
	if err != nil {
		return false, fmt.Errorf("failed to verify kong readiness: %w", err)
	}

	// Kong instance has no configuration, we should push despite the oldSHA and newSHA being equal...
	if hasNoConfiguration {
		// ... unless we're trying to push an empty config in such case skip.
		if cmp.Equal(targetConfig, &file.Content{},
			cmp.FilterPath(
				func(p cmp.Path) bool {
					path := p.String()
					return path == "FormatVersion" || path == "Info"
				},
				cmp.Ignore(),
			),
		) {
			return false, nil
		}

		return true, nil
	}

	return false, nil
}

// kongHasNoConfiguration checks Kong's status endpoint and read its config hash.
// If the config hash reported by Kong is the known empty hash, it's considered crashed.
// This allows providing configuration to Kong instances that have unexpectedly crashed and
// lost their configuration.
func kongHasNoConfiguration(ctx context.Context, client StatusClient, log logrus.FieldLogger) (bool, error) {
	status, err := client.Status(ctx)
	if err != nil {
		return false, err
	}

	if hasNoConfig := status.ConfigurationHash == WellKnownInitialHash; hasNoConfig {
		return true, nil
	}

	return false, nil
}
