package clientset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongfake "github.com/kong/kubernetes-configuration/pkg/clientset/fake"
)

func TestClientset(t *testing.T) {
	t.Run("it can retrieve a fake KongPlugin", func(t *testing.T) {
		cl := kongfake.NewSimpleClientset(&kongv1.KongPlugin{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-plugin",
				Namespace: "test-ns",
			},
		})

		plugin, err := cl.ConfigurationV1().KongPlugins("test-ns").
			Get(t.Context(), "test-plugin", metav1.GetOptions{})
		require.NoError(t, err)
		require.NotNil(t, plugin)
		require.Equal(t, "test-plugin", plugin.Name)
	})
}
