package featuregates

import (
	"fmt"

	"github.com/go-logr/logr"
)

// -----------------------------------------------------------------------------
// Feature Gates - Vars & Consts
// -----------------------------------------------------------------------------

const (
	// GatewayAlphaFeature is the name of the feature-gate for enabling or
	// disabling the Alpha maturity APIs and relevant features for Gateway API.
	GatewayAlphaFeature = "GatewayAlpha"

	// FillIDsFeature is the name of the feature-gate that makes KIC fill in the ID fields of Kong entities (Services,
	// Routes, and Consumers). It ensures that IDs remain stable across restarts of the controller.
	FillIDsFeature = "FillIDs"

	// RewriteURIsFeature is the name of the feature-gate for enabling/disabling konghq.com/rewrite annotation.
	RewriteURIsFeature = "RewriteURIs"

	// ServiceFacade is the name of the feature-gate for enabling KongServiceFacade CR reconciliation.
	ServiceFacade = "ServiceFacade"

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
		GatewayAlphaFeature: false,
		FillIDsFeature:      true,
		RewriteURIsFeature:  false,
		ServiceFacade:       false,
	}
}
