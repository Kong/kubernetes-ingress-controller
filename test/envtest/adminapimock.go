package envtest

import (
	"net/http/httptest"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/v2/test/mocks"
)

// StartAdminAPIServerMock starts a mock Kong Admin API server.
// Server's .Close() method will be called during test's cleanup.
func StartAdminAPIServerMock(t *testing.T) *httptest.Server {
	t.Helper()

	handler := mocks.NewAdminAPIHandler(t)
	s := httptest.NewServer(handler)
	t.Cleanup(func() {
		s.Close()
	})
	return s
}
