package config

import (
	"fmt"
	"maps"
	"slices"
)

const (
	// DocsURL provides a link to the documentation for feature gates in the KIC repository.
	DocsURL = "https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md"
)

type FeatureGates map[string]bool

// NewFeatureGates creates FeatureGates from the given feature gate map, overriding the default settings.
func NewFeatureGates(featureGates map[string]bool) (FeatureGates, error) {
	// Generate a map of feature gates by string names to their controller enablement
	ctrlMap := GetFeatureGatesDefaults()

	// Override the default settings.
	for _, fgName := range slices.Sorted(maps.Keys(featureGates)) {
		if _, ok := ctrlMap[fgName]; !ok {
			return nil, fmt.Errorf("%s is not a valid feature, please see the documentation: %s", fgName, DocsURL)
		}
		ctrlMap[fgName] = featureGates[fgName]
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
