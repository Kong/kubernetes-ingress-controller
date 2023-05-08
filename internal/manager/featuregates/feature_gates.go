package featuregates

import (
	"fmt"

	"github.com/go-logr/logr"
)

// -----------------------------------------------------------------------------
// Feature Gates - Vars & Consts
// -----------------------------------------------------------------------------

const (
	// KnativeFeature is the name of the feature-gate for enabling/disabling Knative.
	KnativeFeature = "Knative"

	// GatewayFeature is the name of the feature-gate for enabling/disabling GatewayFeature APIs.
	GatewayFeature = "Gateway"

	// GatewayAlphaFeature is the name of the feature-gate for enabling or
	// disabling the Alpha maturity APIs and relevant features for Gateway API.
	GatewayAlphaFeature = "GatewayAlpha"

	// CombinedRoutesFeature is the name of the feature-gate for the newer object
	// translation logic that will combine routes for kong services when translating
	// objects like Ingress instead of creating a route per path.
	CombinedRoutesFeature = "CombinedRoutes"

	// ExpressionRoutesFeature is the name of the feature-gate for enabling KIC to translate
	// supported kubernetes objects into expression based routes in kong configrurations
	// when controlled kong gateway uses expression based router by configuring `router_flavor` to `expressions`.
	// Note: this feature is experimental and some resources and features may not be supported.
	// See: https://docs.konghq.com/gateway/latest/key-concepts/routes/expressions/ about expression based routes in Kong 3.0+.
	ExpressionRoutesFeature = "ExpressionRoutes"

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
		setupLog.Info("found configuration option for gated feature", "feature", feature, "enabled", enabled)
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
func GetFeatureGatesDefaults() map[string]bool {
	return map[string]bool{
		KnativeFeature:          false,
		GatewayFeature:          true,
		GatewayAlphaFeature:     false,
		CombinedRoutesFeature:   true,
		ExpressionRoutesFeature: false,
	}
}
