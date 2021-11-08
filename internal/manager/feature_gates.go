package manager

import (
	"fmt"

	"github.com/go-logr/logr"
)

// featureGatesDocsURL provides a link to the documentation for feature gates in the KIC repository
const featureGatesDocsURL = "https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md"

// setupFeatureGates converts feature gates to controller enablement
func setupFeatureGates(setupLog logr.Logger, c *Config) (map[string]bool, error) {
	// generate a map of feature gates by string names to their controller enablement
	ctrlMap := getFeatureGatesDefaults()

	// override the default settings
	for feature, enabled := range c.FeatureGates {
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
		"Knative": false,
		"Gateway": false,
	}
}
