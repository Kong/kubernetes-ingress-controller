package manager

import (
	"fmt"

	"github.com/go-logr/logr"
)

// -----------------------------------------------------------------------------
// Feature Gates - Vars & Consts
// -----------------------------------------------------------------------------

const (
	// knativeFeature is the name of the feature-gate for enabling/disabling Knative.
	knativeFeature = "Knative"

	// gatewayFeature is the name of the feature-gate for enabling/disabling Gateway APIs.
	gatewayFeature = "Gateway"

	// gatewayAlphaFeature is the name of the feature-gate for enabling or
	// disabling the Alpha maturity APIs and relevant features for Gateway API.
	gatewayAlphaFeature = "GatewayAlpha"

	// combinedRoutesFeature is the name of the feature-gate for the newer object
	// translation logic that will combine routes for kong services when translating
	// objects like Ingress instead of creating a route per path.
	combinedRoutesFeature = "CombinedRoutes"

	// featureGatesDocsURL provides a link to the documentation for feature gates in the KIC repository.
	featureGatesDocsURL = "https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md"
)

// setupFeatureGates converts feature gates to controller enablement.
func setupFeatureGates(setupLog logr.Logger, featureGates map[string]bool) (map[string]bool, error) {
	// generate a map of feature gates by string names to their controller enablement
	ctrlMap := getFeatureGatesDefaults()

	// override the default settings
	for feature, enabled := range featureGates {
		setupLog.Info("found configuration option for gated feature", "feature", feature, "enabled", enabled)
		_, ok := ctrlMap[feature]
		if !ok {
			return ctrlMap, fmt.Errorf("%s is not a valid feature, please see the documentation: %s", feature, featureGatesDocsURL)
		}
		ctrlMap[feature] = enabled
	}

	return ctrlMap, nil
}

// getFeatureGatesDefaults initializes a feature gate map given the currently
// supported feature gates options and derives defaults for them based on
// manager configuration options if present.
//
// NOTE: if you're adding a new feature gate, it needs to be added here.
func getFeatureGatesDefaults() map[string]bool {
	return map[string]bool{
		knativeFeature:        false,
		gatewayFeature:        true,
		gatewayAlphaFeature:   false,
		combinedRoutesFeature: true,
	}
}
