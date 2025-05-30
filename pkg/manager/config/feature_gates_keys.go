package config

const (
	// GatewayAlphaFeature is the name of the feature-gate for enabling or
	// disabling the Alpha maturity APIs and relevant features for Gateway API.
	GatewayAlphaFeature = "GatewayAlpha"

	// FillIDsFeature is the name of the feature-gate that makes KIC fill in the ID fields of Kong entities (Services,
	// Routes, and Consumers). It ensures that IDs remain stable across restarts of the controller.
	FillIDsFeature = "FillIDs"

	// RewriteURIsFeature is the name of the feature-gate for enabling/disabling konghq.com/rewrite annotation.
	RewriteURIsFeature = "RewriteURIs"

	// KongServiceFacadeFeature is the name of the feature-gate for enabling KongServiceFacade CR reconciliation.
	KongServiceFacadeFeature = "KongServiceFacade"

	// SanitizeKonnectConfigDumpsFeature is the name of the feature-gate that enables sanitization of Konnect config dumps.
	SanitizeKonnectConfigDumpsFeature = "SanitizeKonnectConfigDumps"

	// FallbackConfigurationFeature is the name of the feature-gate that enables generating fallback configuration in the case
	// of entity errors returned by the Kong Admin API.
	FallbackConfigurationFeature = "FallbackConfiguration"

	// KongCustomEntityFeature is the name of the feature-gate for enabling KongCustomEntity CR reconciliation
	// for configuring custom Kong entities that KIC does not support yet.
	// Requires feature gate `FillIDs` to be enabled.
	KongCustomEntityFeature = "KongCustomEntity"

	// CombinedServicesFromDifferentHTTPRoutesFeature is the name of the feature gate that enables combining rules sharing the same backendRefs
	// from different HTTPRoutes in the same namespace into one Kong gateway service to reduce total number of Kong gateway services.
	CombinedServicesFromDifferentHTTPRoutesFeature = "CombinedServicesFromDifferentHTTPRoutes"

	// StickySessionsTerminatingEndpointsFeature is the name of the feature gate that enables keeping terminating endpoints
	// in Kong upstreams with weight=0 for sticky sessions. This allows existing sessions to continue while preventing
	// new sessions from being assigned to terminating pods.
	StickySessionsTerminatingEndpointsFeature = "StickySessionsTerminatingEndpoints"
)

// GetFeatureGatesDefaults returns the default values for all feature gates.
//
// NOTE: if you're adding a new feature gate, it needs to be added here.
func GetFeatureGatesDefaults() FeatureGates {
	return map[string]bool{
		GatewayAlphaFeature:                            false,
		FillIDsFeature:                                 true,
		RewriteURIsFeature:                             false,
		KongServiceFacadeFeature:                       false,
		SanitizeKonnectConfigDumpsFeature:              true,
		FallbackConfigurationFeature:                   false,
		KongCustomEntityFeature:                        true,
		CombinedServicesFromDifferentHTTPRoutesFeature: false,
		StickySessionsTerminatingEndpointsFeature:      false,
	}
}
