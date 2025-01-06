package translator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func TestGetPluginsAssociatedWithCACertSecret(t *testing.T) {
	kongPluginWithSecret := func(name, secretID string) *kongv1.KongPlugin {
		return &kongv1.KongPlugin{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Config: apiextensionsv1.JSON{
				Raw: []byte(fmt.Sprintf(`{"ca_certificates": ["%s"]}`, secretID)),
			},
		}
	}
	kongClusterPluginWithSecret := func(name, secretID string) *kongv1.KongClusterPlugin {
		return &kongv1.KongClusterPlugin{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
			},
			Config: apiextensionsv1.JSON{
				Raw: []byte(fmt.Sprintf(`{"ca_certificates": ["%s"]}`, secretID)),
			},
		}
	}

	const (
		secretID        = "8a3753e0-093b-43d9-9d39-27985c987d92"
		anotherSecretID = "99fa09c7-f849-4449-891e-19b9a0015763"
	)
	var (
		associatedPlugin           = kongPluginWithSecret("associated_plugin", secretID)
		nonAssociatedPlugin        = kongPluginWithSecret("non_associated_plugin", anotherSecretID)
		associatedClusterPlugin    = kongClusterPluginWithSecret("associated_cluster_plugin", secretID)
		nonAssociatedClusterPlugin = kongClusterPluginWithSecret("non_associated_cluster_plugin", anotherSecretID)
	)
	storer, err := store.NewFakeStore(store.FakeObjects{
		KongPlugins:        []*kongv1.KongPlugin{associatedPlugin, nonAssociatedPlugin},
		KongClusterPlugins: []*kongv1.KongClusterPlugin{associatedClusterPlugin, nonAssociatedClusterPlugin},
	})
	require.NoError(t, err)

	gotPlugins := getPluginsAssociatedWithCACertSecret(secretID, storer)
	expectedPlugins := []client.Object{associatedPlugin, associatedClusterPlugin}
	require.ElementsMatch(t, expectedPlugins, gotPlugins, "expected plugins do not match actual ones")
}
