package controlplanesconfig_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/controlplanesconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
)

type mockControlPlanesConfigServer struct {
	t *testing.T
}

func newMockControlPlanesConfigServer(t *testing.T) *mockControlPlanesConfigServer {
	return &mockControlPlanesConfigServer{
		t: t,
	}
}

func (m *mockControlPlanesConfigServer) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	require.Equal(m.t, metadata.UserAgent(), r.Header.Get("User-Agent"))
}

func TestControlPlanesConfigClientUserAgent(t *testing.T) {
	ts := httptest.NewServer(newMockControlPlanesConfigServer(t))
	t.Cleanup(ts.Close)

	c, err := controlplanesconfig.NewClient(ts.URL)
	require.NoError(t, err)

	r, err := c.GetNodes(context.Background(), &controlplanesconfig.GetNodesParams{})
	require.NoError(t, err)
	r.Body.Close()

	r, err = c.DeleteCoreEntities(context.Background(), "test")
	require.NoError(t, err)
	r.Body.Close()
}
