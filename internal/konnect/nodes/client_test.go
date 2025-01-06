package nodes_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/nodes"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
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

	c, err := nodes.NewClient(adminapi.KonnectConfig{Address: ts.URL})
	require.NoError(t, err)

	_, err = c.GetNode(context.Background(), "test-node-id")
	require.Error(t, err)

	err = c.DeleteNode(context.Background(), "test-node-id")
	require.NoError(t, err)
}
