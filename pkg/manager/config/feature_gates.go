package config

import (
	"fmt"
)

const (
	// DocsURL provides a link to the documentation for feature gates in the KIC repository.
	DocsURL = "https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md"
)

type FeatureGates map[string]bool

// newFeatureGates creates FeatureGates from the given feature gate map, overriding the default settings.
func newFeatureGates(featureGates map[string]bool) (FeatureGates, error) {
	// generate a map of feature gates by string names to their controller enablement
	ctrlMap := FeatureGates(GetFeatureGatesDefaults())

	// override the default settings
	for feature, enabled := range featureGates {
		_, ok := ctrlMap[feature]
		if !ok {
			return ctrlMap, fmt.Errorf("%s is not a valid feature, please see the documentation: %s", feature, DocsURL)
		}
		ctrlMap[feature] = enabled
	}

	// KongCustomEntity requires FillIDs to be enabled, because custom entities requires stable IDs to fill in its "foreign" fields.
	if ctrlMap.Enabled(KongCustomEntityFeature) && !ctrlMap.Enabled(FillIDsFeature) {
		return nil, fmt.Errorf("%s is required if %s is enabled", FillIDsFeature, KongCustomEntityFeature)
	}

	return ctrlMap, nil
}

func (fg FeatureGates) Enabled(feature string) bool {
	return fg[feature]
}
