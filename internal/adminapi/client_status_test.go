package adminapi

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	k8stypes "k8s.io/apimachinery/pkg/types"

	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

func TestClient_AttachStatusClient(t *testing.T) {
	client := &Client{}
	
	// Initially, no status client should be attached
	assert.Nil(t, client.statusClient)

	// Create a real status client for testing
	discoveredAPI := DiscoveredAdminAPI{
		Address: "https://example.com:8100",
		PodRef: k8stypes.NamespacedName{
			Name:      "test-pod",
			Namespace: "test-namespace",
		},
	}
	
	statusClient, err := NewStatusClient(discoveredAPI, managercfg.AdminAPIClientConfig{})
	assert.NoError(t, err)
	
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

func TestNewClientFactoryForWorkspaceWithStatusDiscoverer(t *testing.T) {
	logger := logr.Discard()
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
	logger := logr.Discard()
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

func TestClientFactory_HasStatusDiscoverer(t *testing.T) {
	factory := ClientFactory{
		logger:               logr.Discard(),
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