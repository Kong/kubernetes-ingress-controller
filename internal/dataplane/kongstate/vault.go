package kongstate

import (
	"github.com/kong/go-kong/kong"

	configurationv1alpha1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1alpha1"
)

type Vault struct {
	kong.Vault

	K8sKongVault *configurationv1alpha1.KongVault
}
