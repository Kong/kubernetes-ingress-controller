package parser

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/stretchr/testify/require"
)

func TestGetPluginsAssociatedWithSecret(t *testing.T) {
	s, err := store.NewFakeStore(store.FakeObjects{
		Secrets: []*corev1.Secret{},
	})
	require.NoError(t, err)
	secretID := "8a3753e0-093b-43d9-9d39-27985c987d92"

	associatedPlugins := getPluginsAssociatedWithSecret(s, secretID)
	require.Empty(t, associatedPlugins)
}
