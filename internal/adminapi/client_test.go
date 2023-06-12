package adminapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

type mockAdminAPIServer struct {
	mux *http.ServeMux
	t   *testing.T

	workspaceWasCreated bool
	m                   sync.RWMutex
}

func newMockAdminAPIServer(t *testing.T, ready, workspaceExists bool) *mockAdminAPIServer {
	srv := &mockAdminAPIServer{
		t: t,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if !ready {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
			_, _ = w.Write([]byte(`{}`))
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	mux.HandleFunc("/workspace/workspaces/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if workspaceExists {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	mux.HandleFunc("/workspaces", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if !workspaceExists {
				srv.m.Lock()
				defer srv.m.Unlock()
				srv.workspaceWasCreated = true
				_, _ = w.Write([]byte(`{"id": "workspace"}`))
				w.WriteHeader(http.StatusCreated)
			} else {
				t.Errorf("unexpected workspace creation")
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	srv.mux = mux

	return srv
}

func (m *mockAdminAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.t.Logf("mockAdminAPIServer received request: %s %s", r.Method, r.URL)
	m.mux.ServeHTTP(w, r)
}

func (m *mockAdminAPIServer) WasWorkspaceCreated() bool {
	m.m.RLock()
	defer m.m.RUnlock()
	return m.workspaceWasCreated
}

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
			adminAPIServer := newMockAdminAPIServer(t, tc.adminAPIReady, tc.workspaceExists)
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
