package envtest

import (
	"net/http/httptest"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

// StartAdminAPIServerMock starts a mock Kong Admin API server.
// It accepts a variadic list of options which can configure the test server.
//
// Server's .Close() method will be called during test's cleanup.
func StartAdminAPIServerMock(t *testing.T, opts ...mocks.AdminAPIHandlerOpt) *httptest.Server {
	t.Helper()

	handler := mocks.NewAdminAPIHandler(t, opts...)
	s := httptest.NewServer(handler)
	t.Cleanup(func() {
		s.Close()
	})
	return s
}
