package testenv

import "github.com/kong/kubernetes-ingress-controller/v3/test/consts"

func getFeatureGates() string {
	return consts.DefaultFeatureGates
}
