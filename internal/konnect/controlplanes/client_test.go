package controlplanes_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/controlplanes"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
)

type mockControlPlanesServer struct {
	t *testing.T
}

func newMockControlPlanesServer(t *testing.T) *mockControlPlanesServer {
	return &mockControlPlanesServer{
		t: t,
	}
}

func (m *mockControlPlanesServer) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	require.Equal(m.t, metadata.UserAgent(), r.Header.Get("User-Agent"))
}

func TestControlPlanesClientUserAgent(t *testing.T) {
	ts := httptest.NewServer(newMockControlPlanesServer(t))
	t.Cleanup(ts.Close)

	c, err := controlplanes.NewClient(ts.URL)
	require.NoError(t, err)

	r, err := c.GetControlPlane(context.Background(), uuid.New())
	require.NoError(t, err)
	r.Body.Close()

	r, err = c.DeleteControlPlane(context.Background(), uuid.New())
	require.NoError(t, err)
	r.Body.Close()
}
