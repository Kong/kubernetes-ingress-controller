package kongstate

import (
	"github.com/kong/go-kong/kong"

	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
)

type Vault struct {
	kong.Vault

	K8sKongVault *kongv1alpha1.KongVault
}
