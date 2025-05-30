package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"

	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// testLogger returns a logger for testing.
func testLogger(_ *testing.T) logr.Logger {
	return logr.Discard()
}

func TestStatusClient_IsReady(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError bool
		errorContains string
	}{
		{
			name:          "status endpoint returns 200",
			statusCode:    http.StatusOK,
			expectedError: false,
		},
		{
			name:          "status endpoint returns 500",
			statusCode:    http.StatusInternalServerError,
			expectedError: true,
			errorContains: "status endpoint returned 500",
		},
		{
			name:          "status endpoint returns 404",
			statusCode:    http.StatusNotFound,
			expectedError: true,
			errorContains: "status endpoint returned 404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that responds with the specified status code
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/status", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			// Create a status client
			discoveredAPI := DiscoveredAdminAPI{
				Address: server.URL,
				PodRef: k8stypes.NamespacedName{
					Name:      "test-pod",
					Namespace: "test-namespace",
				},
			}

			client, err := NewStatusClient(discoveredAPI, managercfg.AdminAPIClientConfig{})
			require.NoError(t, err)

			// Test IsReady
			ctx := t.Context()
			err = client.IsReady(ctx)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStatusClient_PodReference(t *testing.T) {
	podRef := k8stypes.NamespacedName{
		Name:      "test-pod",
		Namespace: "test-namespace",
	}

	discoveredAPI := DiscoveredAdminAPI{
		Address: "https://example.com:8100",
		PodRef:  podRef,
	}

	client, err := NewStatusClient(discoveredAPI, managercfg.AdminAPIClientConfig{})
	require.NoError(t, err)

	gotPodRef, ok := client.PodReference()
	assert.True(t, ok)
	assert.Equal(t, podRef, gotPodRef)
}

func TestStatusClient_BaseRootURL(t *testing.T) {
	expectedURL := "https://example.com:8100"
	discoveredAPI := DiscoveredAdminAPI{
		Address: expectedURL,
		PodRef: k8stypes.NamespacedName{
			Name:      "test-pod",
			Namespace: "test-namespace",
		},
	}

	client, err := NewStatusClient(discoveredAPI, managercfg.AdminAPIClientConfig{})
	require.NoError(t, err)

	assert.Equal(t, expectedURL, client.BaseRootURL())
}

func TestStatusClientFactory_CreateStatusClient(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError bool
		errorContains string
	}{
		{
			name:          "successful status client creation",
			statusCode:    http.StatusOK,
			expectedError: false,
		},
		{
			name:          "status client creation fails when status check fails",
			statusCode:    http.StatusInternalServerError,
			expectedError: true,
			errorContains: "status client not ready",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			// Create factory
			factory := NewStatusClientFactory(
				testLogger(t),
				managercfg.AdminAPIClientConfig{},
			)

			discoveredAPI := DiscoveredAdminAPI{
				Address: server.URL,
				PodRef: k8stypes.NamespacedName{
					Name:      "test-pod",
					Namespace: "test-namespace",
				},
			}

			// Test CreateStatusClient
			ctx := t.Context()
			client, err := factory.CreateStatusClient(ctx, discoveredAPI)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, client)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}
