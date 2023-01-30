package helpers

import (
	"net/http"
	"time"
)

// DefaultHTTPClient returns a client that should be used by default in tests.
// All defaults that should be propagated to tests for use should be changed in here.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}
