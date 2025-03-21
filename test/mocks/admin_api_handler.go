package mocks

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
)

const (
	defaultDBLessRootResponse = `{
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
	defaultDBLessStatusResponseWithoutConfigurationHash = `{
	"memory": {
	  "workers_lua_vms": [
		{
		  "http_allocated_gc": "43.99 MiB",
		  "pid": 1260
		},
		{
		  "http_allocated_gc": "43.98 MiB",
		  "pid": 1261
		}
	  ],
	  "lua_shared_dicts": {
		"kong_secrets": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"prometheus_metrics": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong_locks": {
		  "allocated_slabs": "0.06 MiB",
		  "capacity": "8.00 MiB"
		},
		"kong_healthchecks": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong_cluster_events": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong_rate_limiting_counters": {
		  "allocated_slabs": "0.08 MiB",
		  "capacity": "12.00 MiB"
		},
		"kong_core_db_cache": {
		  "allocated_slabs": "0.76 MiB",
		  "capacity": "128.00 MiB"
		},
		"kong_core_db_cache_miss": {
		  "allocated_slabs": "0.08 MiB",
		  "capacity": "12.00 MiB"
		},
		"kong_db_cache": {
		  "allocated_slabs": "0.76 MiB",
		  "capacity": "128.00 MiB"
		},
		"kong_db_cache_miss": {
		  "allocated_slabs": "0.08 MiB",
		  "capacity": "12.00 MiB"
		}
	  }
	},
	"server": {
	  "connections_reading": 0,
	  "total_requests": 615,
	  "connections_writing": 3,
	  "connections_handled": 615,
	  "connections_waiting": 0,
	  "connections_accepted": 615,
	  "connections_active": 3
	}
}`
)

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

	// configurationHash specifies the configuration hash of mocked Kong instance
	// return in /status response.
	configurationHash string
}

type AdminAPIHandlerOpt func(h *AdminAPIHandler)

func WithConfigurationHash(hash string) AdminAPIHandlerOpt {
	return func(h *AdminAPIHandler) {
		h.configurationHash = hash
	}
}

func WithWorkspaceExists(exists bool) AdminAPIHandlerOpt {
	return func(h *AdminAPIHandler) {
		h.workspaceExists = exists
	}
}

func WithReady(ready bool) AdminAPIHandlerOpt {
	return func(h *AdminAPIHandler) {
		h.ready = ready
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
			_, _ = w.Write([]byte(defaultDBLessRootResponse))
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if !h.ready {
				w.WriteHeader(http.StatusServiceUnavailable)
			} else {
				if h.configurationHash != "" {
					_, _ = w.Write([]byte(formatDBLessStatusResponseWithConfigurationHash(h.configurationHash)))
				} else {
					_, _ = w.Write([]byte(defaultDBLessStatusResponseWithoutConfigurationHash))
				}
			}
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	// The handler for the status call in a specific workspace.
	mux.HandleFunc("/workspace/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if !h.ready {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			if !h.workspaceExists && !h.workspaceWasCreated.Load() {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if h.configurationHash != "" {
				_, _ = w.Write([]byte(formatDBLessStatusResponseWithConfigurationHash(h.configurationHash)))
			} else {
				_, _ = w.Write([]byte(defaultDBLessStatusResponseWithoutConfigurationHash))
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
	// this path gets spammed by the readiness checker and shouldn't be of interest for logging
	if r.URL.Path != "/status" {
		m.t.Logf("AdminAPIHandler received request: %s %s", r.Method, r.URL)
	}
	m.mux.ServeHTTP(w, r)
}

func (m *AdminAPIHandler) WasWorkspaceCreated() bool {
	return m.workspaceWasCreated.Load()
}

func formatDBLessStatusResponseWithConfigurationHash(hash string) string {
	const defaultDBLessStatusResponseWithConfigurationHash = `{
	"configuration_hash": "%s",
	"memory": {
	  "workers_lua_vms": [
		{
		  "http_allocated_gc": "43.99 MiB",
		  "pid": 1260
		},
		{
		  "http_allocated_gc": "43.98 MiB",
		  "pid": 1261
		}
	  ],
	  "lua_shared_dicts": {
		"kong_secrets": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"prometheus_metrics": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong_locks": {
		  "allocated_slabs": "0.06 MiB",
		  "capacity": "8.00 MiB"
		},
		"kong_healthchecks": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong_cluster_events": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong_rate_limiting_counters": {
		  "allocated_slabs": "0.08 MiB",
		  "capacity": "12.00 MiB"
		},
		"kong_core_db_cache": {
		  "allocated_slabs": "0.76 MiB",
		  "capacity": "128.00 MiB"
		},
		"kong_core_db_cache_miss": {
		  "allocated_slabs": "0.08 MiB",
		  "capacity": "12.00 MiB"
		},
		"kong_db_cache": {
		  "allocated_slabs": "0.76 MiB",
		  "capacity": "128.00 MiB"
		},
		"kong_db_cache_miss": {
		  "allocated_slabs": "0.08 MiB",
		  "capacity": "12.00 MiB"
		}
	  }
	},
	"server": {
	  "connections_reading": 0,
	  "total_requests": 615,
	  "connections_writing": 3,
	  "connections_handled": 615,
	  "connections_waiting": 0,
	  "connections_accepted": 615,
	  "connections_active": 3
	}
}`

	return fmt.Sprintf(defaultDBLessStatusResponseWithConfigurationHash, hash)
}
