package sendconfig

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
)

const (
	// WellKnownInitialHash is the hash of an empty configuration.
	WellKnownInitialHash = "00000000000000000000000000000000"
)

type ConfigurationChangeDetector interface {
	// HasConfigurationChanged verifies whether configuration has changed by comparing old and new config's SHAs.
	// In case the SHAs are equal, it still can return true if a client is considered crashed based on its status.
	HasConfigurationChanged(ctx context.Context, oldSHA, newSHA []byte, client KonnectAwareClient, statusClient StatusClient) (bool, error)
}

type KonnectAwareClient interface {
	IsKonnect() bool
}

type StatusClient interface {
	Status(context.Context) (*kong.Status, error)
}

type DefaultConfigurationChangeDetector struct {
	latestReportedSHA []byte
	shaLock           sync.RWMutex
	log               logrus.FieldLogger
}

func NewDefaultClientConfigurationChangeDetector(log logrus.FieldLogger) *DefaultConfigurationChangeDetector {
	return &DefaultConfigurationChangeDetector{log: log}
}

func (d *DefaultConfigurationChangeDetector) HasConfigurationChanged(
	ctx context.Context,
	oldSHA, newSHA []byte,
	client KonnectAwareClient,
	statusClient StatusClient,
) (bool, error) {
	if !bytes.Equal(oldSHA, newSHA) {
		return true, nil
	}
	if !d.hasSHAUpdateAlreadyBeenReported(newSHA) {
		d.log.Debugf("sha %s has been reported", hex.EncodeToString(newSHA))
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
	// Kong instance has no configuration, we should push despite the oldSHA and newSHA being equal.
	if hasNoConfiguration {
		return true, nil
	}

	return false, nil
}

// hasSHAUpdateAlreadyBeenReported is a helper function to allow
// sendconfig internals to be aware of the last logged/reported
// update to the Kong Admin API. Given the most recent update SHA,
// it will return true/false whether or not that SHA has previously
// been reported (logged, e.t.c.) so that the caller can make
// decisions (such as staggering or stifling duplicate log lines).
func (d *DefaultConfigurationChangeDetector) hasSHAUpdateAlreadyBeenReported(latestUpdateSHA []byte) bool {
	d.shaLock.Lock()
	defer d.shaLock.Unlock()
	if bytes.Equal(d.latestReportedSHA, latestUpdateSHA) {
		return true
	}
	d.latestReportedSHA = latestUpdateSHA
	return false
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
		log.Debugf("starting to send configuration (hash: %s)", status.ConfigurationHash)
		return true, nil
	}

	return false, nil
}
