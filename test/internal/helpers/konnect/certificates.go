package konnect

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	cpc "github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/controlplanesconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

// CreateClientCertificate creates a TLS client certificate and POSTs it to Konnect Control Plane configuration API
// so that KIC can use the certificates to authenticate against Konnect Admin API.
func CreateClientCertificate(ctx context.Context, t *testing.T, cpID string) (certPEM string, keyPEM string) {
	t.Helper()

	rgConfigClient, err := cpc.NewClientWithResponses(fmt.Sprintf(konnectControlPlanesConfigBaseURLFmt, cpID), cpc.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+accessToken())
			return nil
		}),
		cpc.WithHTTPClient(helpers.RetryableHTTPClient(helpers.DefaultHTTPClient())),
	)
	require.NoError(t, err)

	cert, key := certificate.MustGenerateSelfSignedCertPEMFormat()

	t.Log("creating client certificate in Konnect")
	resp, err := rgConfigClient.PostDpClientCertificatesWithResponse(ctx, cpc.PostDpClientCertificatesJSONRequestBody{
		Cert: string(cert),
	})
	require.NoError(t, err)
	require.Equalf(t, http.StatusCreated, resp.StatusCode(), "failed creating client certificate: %s", string(resp.Body))

	return string(cert), string(key)
}
