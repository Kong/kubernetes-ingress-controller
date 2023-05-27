package envtest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// StartAdminAPIServerMock starts a mock Kong Admin API server.
// Server's .Close() method will be called during test's cleanup.
func StartAdminAPIServerMock(t *testing.T) *httptest.Server {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Root struct {
			Configuration map[string]any `json:"configuration"`
			Version       string         `json:"version"`
		}
		body, err := json.Marshal(Root{
			Version: "3.3.0",
			Configuration: map[string]any{
				"database":      "off",
				"router_flavor": "traditional",
			},
		})
		if err != nil {
			w.WriteHeader(500)
			return
		}

		if _, err := w.Write(body); err != nil {
			w.WriteHeader(500)
			return
		}
	})
	s := httptest.NewServer(handler)
	t.Cleanup(func() {
		s.Close()
	})
	return s
}
