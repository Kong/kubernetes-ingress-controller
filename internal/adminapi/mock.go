package adminapi

import (
	"net/http"
	"sync"
	"testing"
)

// MockAdminAPIServer is a mock implementation of the Admin API server. It only implements the endpoints that are
// required for the tests.
type MockAdminAPIServer struct {
	mux *http.ServeMux
	t   *testing.T

	workspaceWasCreated bool
	m                   sync.RWMutex
}

func NewMockAdminAPIServer(t *testing.T, ready, workspaceExists bool) *MockAdminAPIServer {
	srv := &MockAdminAPIServer{
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

func (m *MockAdminAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.t.Logf("MockAdminAPIServer received request: %s %s", r.Method, r.URL)
	m.mux.ServeHTTP(w, r)
}

func (m *MockAdminAPIServer) WasWorkspaceCreated() bool {
	m.m.RLock()
	defer m.m.RUnlock()
	return m.workspaceWasCreated
}
