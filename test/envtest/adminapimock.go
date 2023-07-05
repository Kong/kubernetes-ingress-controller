package envtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const dblessConfig = `{
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

// StartAdminAPIServerMock starts a mock Kong Admin API server.
// Server's .Close() method will be called during test's cleanup.
func StartAdminAPIServerMock(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		t.Logf("Admin API config: %s %s %s", r.Method, r.URL, string(body))
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(dblessConfig)); err != nil {
			w.WriteHeader(500)
			return
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(dblessConfig)); err != nil {
			w.WriteHeader(500)
			return
		}
	})

	s := httptest.NewServer(mux)
	t.Cleanup(func() {
		s.Close()
	})
	return s
}
