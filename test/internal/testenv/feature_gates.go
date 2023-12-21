package testenv

import (
	"github.com/blang/semver/v4"

	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
)

func GetFeatureGates() string {
	// Due to the possibility of running tests between different versions,
	// it is necessary to adjust feature gates according to different KIC versions.
	//
	// Versions below 3.1.x cannot recognize the KongServiceFacade feature gate.
	// We only need to set `GatewayAlpha=true`.
	//
	// https://github.com/Kong/kubernetes-ingress-controller/issues/5373
	tag := ControllerTag()
	if tag != "" {
		// Currently, the latest version is 3.0.x, once 3.1.x is released,
		// we can remove this logic.
		if tag == "latest" {
			return "GatewayAlpha=true"
		}
		if v, err := semver.Make(tag); err == nil {
			minVersion, _ := semver.ParseRange("<3.1.x")
			if minVersion(v) {
				return "GatewayAlpha=true"
			}
		}
	}
	return consts.DefaultFeatureGates
}
