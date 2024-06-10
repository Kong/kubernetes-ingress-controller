package consts

const (
	// DefaultFeatureGates is the default feature gates setting that should be
	// provided if none are provided by the user. This generally includes features
	// that are innocuous, or otherwise don't actually get triggered unless the
	// user takes further action.
	DefaultFeatureGates = "GatewayAlpha=true,KongServiceFacade=true,RewriteURIs=true,FallbackConfiguration=true,KongCustomEntity=true"
)
