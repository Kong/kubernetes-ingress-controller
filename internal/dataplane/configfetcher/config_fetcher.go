package configfetcher

import (
	"context"
	"errors"

	"github.com/kong/deck/dump"
	"github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/sirupsen/logrus"
)

type LastValidConfigFetcher interface {
	// TryFetchingValidConfigFromGateways tries to fetch a valid configuration from all gateways and persists it if found.
	TryFetchingValidConfigFromGateways(ctx context.Context, logger logrus.FieldLogger, gatewayClients []*adminapi.Client) error

	// LastValidConfig returns the last valid config and true if there's one available. Otherwise, second return value is false.
	LastValidConfig() (*kongstate.KongState, bool)

	// StoreLastValidConfig stores a given configuration as the last valid config. Should be used when the configuration was successfully accepted by a gateway.
	StoreLastValidConfig(s *kongstate.KongState)
}

type DefaultKongLastGoodConfigFetcher struct {
	config         dump.Config
	lastValidState *kongstate.KongState
	// fillIDs enables the last valid kongState to be filled in the IDs fields of Kong entities
	// - Services, Routes, and Consumers - based on their names. It ensures that IDs remain
	// stable across restarts of the controller.
	fillIDs bool
}

func NewDefaultKongLastGoodConfigFetcher(fillIDs bool) *DefaultKongLastGoodConfigFetcher {
	return &DefaultKongLastGoodConfigFetcher{
		config:  dump.Config{},
		fillIDs: fillIDs,
	}
}

func (cf *DefaultKongLastGoodConfigFetcher) LastValidConfig() (*kongstate.KongState, bool) {
	if cf.lastValidState != nil {
		return cf.lastValidState, true
	}
	return nil, false
}

func (cf *DefaultKongLastGoodConfigFetcher) StoreLastValidConfig(s *kongstate.KongState) {
	cf.lastValidState = s
}

func (cf *DefaultKongLastGoodConfigFetcher) TryFetchingValidConfigFromGateways(
	ctx context.Context,
	logger logrus.FieldLogger,
	gatewayClients []*adminapi.Client,
) error {
	logger.Debugf("fetching last good configuration from %d gateway clients", len(gatewayClients))

	var goodKongState *kongstate.KongState
	var errs error
	for _, client := range gatewayClients {
		rs, err := cf.getKongRawState(ctx, client.AdminAPIClient())
		if err != nil {
			errs = errors.Join(errs, err)
		}
		status, err := cf.getKongStatus(ctx, client.AdminAPIClient())
		if err != nil {
			errs = errors.Join(errs, err)
		}
		ks := KongRawStateToKongState(rs)
		if status.ConfigurationHash != sendconfig.WellKnownInitialHash {
			// Get the first good one as the one to be used.
			goodKongState = ks
			break
		}
	}
	if goodKongState != nil {
		if cf.fillIDs {
			goodKongState.FillIDs(logger)
		}
		cf.lastValidState = goodKongState
		logger.Debug("last good configuration fetched from a Kong node")
	}
	return errs
}

func (cf *DefaultKongLastGoodConfigFetcher) getKongRawState(ctx context.Context, client *kong.Client) (*utils.KongRawState, error) {
	return dump.Get(ctx, client, cf.config)
}

func (cf *DefaultKongLastGoodConfigFetcher) getKongStatus(ctx context.Context, client *kong.Client) (*kong.Status, error) {
	return client.Status(ctx)
}
