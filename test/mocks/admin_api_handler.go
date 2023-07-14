package mocks

import (
	"net/http"
	"sync/atomic"
	"testing"
)

const defaultDBLessStatusResponse = `{
	"version": "3.3.0",
	"configuration": {
		"database": "off",
		"router_flavor": "traditional",
		"role": "traditional",
		"proxy_listeners": [
			{
				"backlog=%d+": false,
				"ipv6only=on": false,
				"ipv6only=off": false,
				"ssl": false,
				"so_keepalive=off": false,
				"so_keepalive=%w*:%w*:%d*": false,
				"listener": "0.0.0.0:8000",
				"bind": false,
				"port": 8000,
				"deferred": false,
				"so_keepalive=on": false,
				"http2": false,
				"proxy_protocol": false,
				"ip": "0.0.0.0",
				"reuseport": false
			}
		]
	}
}`

// AdminAPIHandler is a mock implementation of the Admin API. It only implements the endpoints that are
// required for the tests.
type AdminAPIHandler struct {
	mux *http.ServeMux
	t   *testing.T

	// ready is a flag that indicates whether the server should return a 200 OK or a 503 Service Unavailable.
	// It's set to true by default.
	ready bool

	// workspaceExists makes `/workspace/workspaces/:id` return 200 when true, or 404 otherwise.
	workspaceExists bool

	// workspaceWasCreated is set to true when a workspace `POST /workspaces` was called.
	workspaceWasCreated atomic.Bool
}

type AdminAPIHandlerOpt func(srv *AdminAPIHandler)

func WithWorkspaceExists(exists bool) AdminAPIHandlerOpt {
	return func(srv *AdminAPIHandler) {
		srv.workspaceExists = exists
	}
}

func WithReady(ready bool) AdminAPIHandlerOpt {
	return func(srv *AdminAPIHandler) {
		srv.ready = ready
	}
}

func NewAdminAPIHandler(t *testing.T, opts ...AdminAPIHandlerOpt) *AdminAPIHandler {
	h := &AdminAPIHandler{
		t:     t,
		ready: true,
	}

	for _, opt := range opts {
		opt(h)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(defaultDBLessStatusResponse))
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if !h.ready {
				w.WriteHeader(http.StatusServiceUnavailable)
			} else {
				_, _ = w.Write([]byte(defaultDBLessStatusResponse))
			}
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	mux.HandleFunc("/workspace/workspaces/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if h.workspaceExists {
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
			if !h.workspaceExists {
				h.workspaceWasCreated.Store(true)
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id": "workspace"}`))
			} else {
				t.Errorf("unexpected workspace creation")
			}
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(`{"version": "3.3.0"}`))
			return
		}
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	h.mux = mux
	return h
}

func (m *AdminAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.t.Logf("AdminAPIHandler received request: %s %s", r.Method, r.URL)
	m.mux.ServeHTTP(w, r)
}

func (m *AdminAPIHandler) WasWorkspaceCreated() bool {
	return m.workspaceWasCreated.Load()
}
