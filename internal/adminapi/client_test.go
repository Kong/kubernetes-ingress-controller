package adminapi_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func TestClientFactory_CreateAdminAPIClientAttachesPodReference(t *testing.T) {
	factory := adminapi.NewClientFactoryForWorkspace(logr.Discard(), "workspace", adminapi.ClientOpts{}, "")

	adminAPIHandler := mocks.NewAdminAPIHandler(t)
	adminAPIServer := httptest.NewServer(adminAPIHandler)
	t.Cleanup(func() { adminAPIServer.Close() })

	client, err := factory.CreateAdminAPIClient(context.Background(), adminapi.DiscoveredAdminAPI{
		Address: adminAPIServer.URL,
		PodRef: k8stypes.NamespacedName{
			Namespace: "namespace",
			Name:      "name",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	ref, ok := client.PodReference()
	require.True(t, ok, "expected pod reference to be attached to the client")
	require.Equal(t, k8stypes.NamespacedName{
		Namespace: "namespace",
		Name:      "name",
	}, ref)
}
