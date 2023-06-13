package adminapi_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

func TestClientFactory_CreateAdminAPIClient(t *testing.T) {
	const (
		workspace  = "workspace"
		adminToken = "token"
	)

	testCases := []struct {
		name            string
		adminAPIReady   bool
		workspaceExists bool
		expectError     error
	}{
		{
			name:            "admin api is ready and workspace exists",
			adminAPIReady:   true,
			workspaceExists: true,
		},
		{
			name:            "admin api is ready and workspace doesn't exist",
			adminAPIReady:   true,
			workspaceExists: false,
		},
		{
			name:          "admin api is not ready",
			adminAPIReady: false,
			expectError:   adminapi.KongClientNotReadyError{},
		},
	}

	factory := adminapi.NewClientFactoryForWorkspace(workspace, adminapi.HTTPClientOpts{}, adminToken)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adminAPIServer := adminapi.NewMockAdminAPIServer(t, tc.adminAPIReady, tc.workspaceExists)
			adminAPI := httptest.NewServer(adminAPIServer)
			t.Cleanup(func() {
				adminAPI.Close()
			})

			client, err := factory.CreateAdminAPIClient(context.Background(), adminapi.DiscoveredAdminAPI{
				Address: adminAPI.URL,
				PodRef: k8stypes.NamespacedName{
					Namespace: "namespace",
					Name:      "name",
				},
			})

			if tc.expectError != nil {
				require.IsType(t, err, tc.expectError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, client)

			if !tc.workspaceExists {
				require.True(t, adminAPIServer.WasWorkspaceCreated(), "expected workspace to be created")
			}

			ref, ok := client.PodReference()
			require.True(t, ok, "expected pod reference to be attached to the client")
			require.Equal(t, k8stypes.NamespacedName{
				Namespace: "namespace",
				Name:      "name",
			}, ref)
		})
	}
}
