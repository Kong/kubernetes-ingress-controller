package controlplanes_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	sdkkonnectgo "github.com/Kong/sdk-konnect-go"
	sdkkonnectcomp "github.com/Kong/sdk-konnect-go/models/components"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/sdk"
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

func (m *mockControlPlanesServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		require.Equal(m.t, metadata.UserAgent(), r.Header.Get("User-Agent"))
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		require.Equal(m.t, metadata.UserAgent(), r.Header.Get("User-Agent"))

	}
}

func TestControlPlanesClientUserAgent(t *testing.T) {
	ts := httptest.NewServer(newMockControlPlanesServer(t))
	t.Cleanup(ts.Close)

	ctx := context.Background()
	sdk := sdk.New("kpat_xxx", sdkkonnectgo.WithServerURL(ts.URL))

	_, err := sdk.ControlPlanes.CreateControlPlane(ctx, sdkkonnectcomp.CreateControlPlaneRequest{
		Name: "test",
	})
	// NOTE: just check the user agent and do not attempt to mock out the entire response.
	require.ErrorContains(t, err, "unknown content-type received: : Status 201")

	_, err = sdk.ControlPlanes.DeleteControlPlane(ctx, "id")
	// NOTE: just check the user agent and do not attempt to mock out the entire response.
	require.ErrorContains(t, err, "unknown status code returned: Status 200")
}
