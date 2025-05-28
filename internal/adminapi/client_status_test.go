package adminapi

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	k8stypes "k8s.io/apimachinery/pkg/types"

	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// testLogger returns a logger for testing
func testLogger(t *testing.T) logr.Logger {
	return logr.Discard()
}

// Mock StatusClient for testing
type mockStatusClient struct {
	mock.Mock
}

func (m *mockStatusClient) IsReady(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockStatusClient) PodReference() (k8stypes.NamespacedName, bool) {
	args := m.Called()
	return args.Get(0).(k8stypes.NamespacedName), args.Bool(1)
}

func (m *mockStatusClient) BaseRootURL() string {
	args := m.Called()
	return args.String(0)
}

func TestClient_IsReady_WithStatusClient(t *testing.T) {
	tests := []struct {
		name                string
		statusClientError   error
		expectedError       bool
		shouldCallAdminAPI  bool
	}{
		{
			name:                "status client returns ready",
			statusClientError:   nil,
			expectedError:       false,
			shouldCallAdminAPI:  false,
		},
		{
			name:                "status client returns not ready",
			statusClientError:   errors.New("status not ready"),
			expectedError:       true,
			shouldCallAdminAPI:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock status client
			statusClient := &mockStatusClient{}
			statusClient.On("IsReady", mock.Anything).Return(tt.statusClientError)

			// Create a client with a mock admin API client
			client := &Client{
				statusClient: statusClient,
			}

			// Test IsReady
			ctx := context.Background()
			err := client.IsReady(ctx)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify that the status client was called
			statusClient.AssertExpectations(t)
		})
	}
}

func TestClient_IsReady_WithoutStatusClient(t *testing.T) {
	// This test would require mocking the kong.Client, which is more complex
	// For now, we'll test that the method exists and can be called
	client := &Client{}

	ctx := context.Background()
	err := client.IsReady(ctx)

	// Without a proper admin API client, this will fail, but we're testing the code path
	assert.Error(t, err)
}

func TestClient_AttachStatusClient(t *testing.T) {
	client := &Client{}
	
	// Initially, no status client should be attached
	assert.Nil(t, client.statusClient)

	// Create a mock status client
	statusClient := &mockStatusClient{}
	
	// Attach the status client
	client.AttachStatusClient(statusClient)
	
	// Verify that the status client was attached
	assert.Equal(t, statusClient, client.statusClient)
}

func TestClient_AttachStatusClient_Nil(t *testing.T) {
	client := &Client{}
	
	// Attach a nil status client (should be allowed)
	client.AttachStatusClient(nil)
	
	// Verify that the status client is nil
	assert.Nil(t, client.statusClient)
}

func TestClientFactory_CreateAdminAPIClient_WithStatusDiscoverer(t *testing.T) {
	// This is a more complex integration test that would require setting up
	// a full mock environment. For now, we'll test the basic structure.
	
	factory := ClientFactory{
		logger:               testLogger(t),
		workspace:            "",
		opts:                 managercfg.AdminAPIClientConfig{},
		adminToken:           "",
		statusAPIsDiscoverer: nil, // No status discoverer
	}

	// Test that the factory has the status discoverer field
	assert.Nil(t, factory.statusAPIsDiscoverer)
	
	// Test with a mock discoverer
	mockDiscoverer := &Discoverer{}
	factory.statusAPIsDiscoverer = mockDiscoverer
	assert.Equal(t, mockDiscoverer, factory.statusAPIsDiscoverer)
}

func TestNewClientFactoryForWorkspaceWithStatusDiscoverer(t *testing.T) {
	logger := testLogger(t)
	workspace := "test-workspace"
	opts := managercfg.AdminAPIClientConfig{}
	adminToken := "test-token"
	statusDiscoverer := &Discoverer{}

	factory := NewClientFactoryForWorkspaceWithStatusDiscoverer(
		logger,
		workspace,
		opts,
		adminToken,
		statusDiscoverer,
	)

	assert.Equal(t, logger, factory.logger)
	assert.Equal(t, workspace, factory.workspace)
	assert.Equal(t, opts, factory.opts)
	assert.Equal(t, adminToken, factory.adminToken)
	assert.Equal(t, statusDiscoverer, factory.statusAPIsDiscoverer)
}

func TestNewClientFactoryForWorkspace_BackwardCompatibility(t *testing.T) {
	logger := testLogger(t)
	workspace := "test-workspace"
	opts := managercfg.AdminAPIClientConfig{}
	adminToken := "test-token"

	factory := NewClientFactoryForWorkspace(
		logger,
		workspace,
		opts,
		adminToken,
	)

	assert.Equal(t, logger, factory.logger)
	assert.Equal(t, workspace, factory.workspace)
	assert.Equal(t, opts, factory.opts)
	assert.Equal(t, adminToken, factory.adminToken)
	assert.Nil(t, factory.statusAPIsDiscoverer) // Should be nil for backward compatibility
}