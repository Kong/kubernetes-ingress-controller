package configfetcher

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/dump"
	"github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

type LastValidConfigFetcher interface {
	// TryFetchingValidConfigFromGateways tries to fetch a valid configuration from all gateways and persists it if found.
	TryFetchingValidConfigFromGateways(ctx context.Context, logger logr.Logger, gatewayClients []*adminapi.Client) error

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
	// workspace is the workspace name used in generating deterministic IDs. Only used when fillIDs = true.
	workspace string
	// licenseGetter is an optional license provider.
	licenseGetter license.Getter
}

func NewDefaultKongLastGoodConfigFetcher(fillIDs bool, workspace string) *DefaultKongLastGoodConfigFetcher {
	return &DefaultKongLastGoodConfigFetcher{
		config:    dump.Config{},
		fillIDs:   fillIDs,
		workspace: workspace,
	}
}

// InjectLicenseGetter adds a license getter to the config fetcher.
func (cf *DefaultKongLastGoodConfigFetcher) InjectLicenseGetter(licenseGetter license.Getter) {
	cf.licenseGetter = licenseGetter
}

func (cf *DefaultKongLastGoodConfigFetcher) LastValidConfig() (*kongstate.KongState, bool) {
	if cf.lastValidState != nil {
		// TODO the translator version of this also has a condition on
		// `t.featureFlags.EnterpriseEdition == true`. This isn't a typical feature
		// flag, and is instead set based on whether we see we're talking to a
		// kong-gateway image. It's not used for anything other than the license.
		// Do we actually need this? While the OSS image does not recognize license
		// entities, you'll only have license entities if you created a KongLicense
		// or set up Konnect credentials that can pull one. In either of those
		// cases you arguably _should_ be using the kong-gateway image, so we
		// shouldn't really exclude it if you're not :thinking:
		if cf.licenseGetter != nil {
			optionalLicense := cf.licenseGetter.GetLicense()
			if l, ok := optionalLicense.Get(); ok {
				// TODO cf.lastValidState is a pointer, and a pointer to a rather large
				// struct we probably don't want to copy. this assignment without
				// copying it does actually _modify_ the last good config, and could
				// technically transform it from valid to invalid in a contrived
				// scenario where the config had a valid license and the environment
				// provided an invalid one.
				//
				// The environment arguably shouldn't provide an invalid one, but we
				// have seen Konnect provide invalid licenses (though we think this was
				// a fluke due to a change in their license management code), and you
				// could potentially overwrite your previously valid KongLicense with
				// garbage.
				//
				// On the other hand, any license in the last valid config would simply
				// be whatever license you had when the config was last valid. That
				// license could have since expired, and attempting to use it would
				// fail. Using the latest available license is the better option in
				// that case.
				//
				// We can reasonably expect Konnect to provide valid licenses (Konnect
				// needs to fix something if they aren't) and a broken KongLicense is
				// something you should be able to fix easily: you should only have one
				// of them managed by the superuser, so there shouldn't be any
				// confusion around where it's coming from.
				cf.lastValidState.Licenses = []kongstate.License{{License: l}}
			}
		}
		return cf.lastValidState, true
	}
	return nil, false
}

func (cf *DefaultKongLastGoodConfigFetcher) StoreLastValidConfig(s *kongstate.KongState) {
	cf.lastValidState = s
}

func (cf *DefaultKongLastGoodConfigFetcher) TryFetchingValidConfigFromGateways(
	ctx context.Context,
	logger logr.Logger,
	gatewayClients []*adminapi.Client,
) error {
	logger.V(util.DebugLevel).Info("Fetching last good configuration from gateway clients", "count", len(gatewayClients))

	var (
		goodKongState *kongstate.KongState
		errs          error
		clientUsed    *adminapi.Client
	)
	for _, client := range gatewayClients {
		logger.V(util.DebugLevel).Info("Fetching configuration", "url", client.BaseRootURL())
		rs, err := cf.getKongRawState(ctx, client.AdminAPIClient())
		if err != nil {
			errs = errors.Join(errs, err)
		}
		if rs == nil {
			errs = errors.Join(errs, fmt.Errorf("failed to fetch configuration from %q, got nil kong raw state", client.BaseRootURL()))
			continue
		}
		status, err := cf.getKongStatus(ctx, client.AdminAPIClient())
		if err != nil {
			errs = errors.Join(errs, err)
		}
		if status == nil {
			continue
		}

		if status.ConfigurationHash != sendconfig.WellKnownInitialHash {
			// Get the first good one as the one to be used.
			clientUsed = client
			ks := KongRawStateToKongState(rs)
			goodKongState = ks
			break
		}
	}
	if goodKongState != nil {
		if cf.fillIDs {
			goodKongState.FillIDs(logger, cf.workspace)
		}
		cf.lastValidState = goodKongState
		logger.V(util.DebugLevel).Info("Last good configuration fetched from Kong node", "url", clientUsed.BaseRootURL())
	}
	return errs
}

func (cf *DefaultKongLastGoodConfigFetcher) getKongRawState(ctx context.Context, client *kong.Client) (*utils.KongRawState, error) {
	return dump.Get(ctx, client, cf.config)
}

func (cf *DefaultKongLastGoodConfigFetcher) getKongStatus(ctx context.Context, client *kong.Client) (*kong.Status, error) {
	return client.Status(ctx)
}
