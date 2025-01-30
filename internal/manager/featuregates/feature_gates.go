package featuregates

import (
	"fmt"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// -----------------------------------------------------------------------------
// Feature Gates - Vars & Consts
// -----------------------------------------------------------------------------

const (
	// GatewayAlphaFeature is the name of the feature-gate for enabling or
	// disabling the Alpha maturity APIs and relevant features for Gateway API.

	// FillIDsFeature is the name of the feature-gate that makes KIC fill in the ID fields of Kong entities (Services,
	// Routes, and Consumers). It ensures that IDs remain stable across restarts of the controller.

	// RewriteURIsFeature is the name of the feature-gate for enabling/disabling konghq.com/rewrite annotation.

	// KongServiceFacade is the name of the feature-gate for enabling KongServiceFacade CR reconciliation.

	// SanitizeKonnectConfigDumps is the name of the feature-gate that enables sanitization of Konnect config dumps.

	// FallbackConfiguration is the name of the feature-gate that enables generating fallback configuration in the case
	// of entity errors returned by the Kong Admin API.

	// KongCustomEntity is the name of the feature-gate for enabling KongCustomEntity CR reconciliation
	// for configuring custom Kong entities that KIC does not support yet.
	// Requires feature gate `FillIDs` to be enabled.

	// CombinedServicesFromDifferentHTTPRoutes is the name of the feature gate that enables combining rules sharing the same backendRefs
	// from different HTTPRoutes in the same namespace into one Kong gateway service to reduce total number of Kong gateway services.

	// DocsURL provides a link to the documentation for feature gates in the KIC repository.
	DocsURL = "https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md"
)

type FeatureGates map[string]bool

// New creates FeatureGates from the given feature gate map, overriding the default settings.
func New(setupLog logr.Logger, featureGates map[string]bool) (FeatureGates, error) {
	// generate a map of feature gates by string names to their controller enablement
	ctrlMap := GetFeatureGatesDefaults()

	// override the default settings
	for feature, enabled := range featureGates {
		setupLog.Info("Found configuration option for gated feature", "feature", feature, "enabled", enabled)
		_, ok := ctrlMap[feature]
		if !ok {
			return ctrlMap, fmt.Errorf("%s is not a valid feature, please see the documentation: %s", feature, DocsURL)
		}
		ctrlMap[feature] = enabled
	}

	// KongCustomEntity requires FillIDs to be enabled, because custom entities requires stable IDs to fill in its "foreign" fields.
	if ctrlMap.Enabled(config.KongCustomEntityFeature) && !ctrlMap.Enabled(config.FillIDsFeature) {
		return nil, fmt.Errorf("%s is required if %s is enabled", config.FillIDsFeature, config.KongCustomEntityFeature)
	}

	return ctrlMap, nil
}

func (fg FeatureGates) Enabled(feature string) bool {
	return fg[feature]
}

// GetFeatureGatesDefaults initializes a feature gate map given the currently
// supported feature gates options and derives defaults for them based on
// manager configuration options if present.
//
// NOTE: if you're adding a new feature gate, it needs to be added here.
func GetFeatureGatesDefaults() FeatureGates {
	return map[string]bool{
		config.GatewayAlphaFeature:                            false,
		config.FillIDsFeature:                                 true,
		config.RewriteURIsFeature:                             false,
		config.KongServiceFacadeFeature:                       false,
		config.SanitizeKonnectConfigDumpsFeature:              true,
		config.FallbackConfigurationFeature:                   false,
		config.KongCustomEntityFeature:                        true,
		config.CombinedServicesFromDifferentHTTPRoutesFeature: false,
	}
}
