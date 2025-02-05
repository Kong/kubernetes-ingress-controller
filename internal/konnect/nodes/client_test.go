package nodes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/nodes"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/metadata"
)

type mockNodesServer struct {
	t *testing.T
}

func newMockNodesServer(t *testing.T) *mockNodesServer {
	return &mockNodesServer{
		t: t,
	}
}

func (m *mockNodesServer) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	require.Equal(m.t, metadata.UserAgent(), r.Header.Get("User-Agent"))
}

func TestNodesClientUserAgent(t *testing.T) {
	ts := httptest.NewServer(newMockNodesServer(t))
	t.Cleanup(ts.Close)

	c, err := nodes.NewClient(config.KonnectConfig{Address: ts.URL})
	require.NoError(t, err)

	_, err = c.GetNode(t.Context(), "test-node-id")
	require.Error(t, err)

	err = c.DeleteNode(t.Context(), "test-node-id")
	require.NoError(t, err)
}
