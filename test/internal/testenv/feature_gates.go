//go:build !knative

package testenv

import "github.com/kong/kubernetes-ingress-controller/v2/test/consts"

func getFeatureGates() string {
	return consts.DefaultFeatureGates
}
