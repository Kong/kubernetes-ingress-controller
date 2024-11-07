package konnect

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdkkonnectgo "github.com/Kong/sdk-konnect-go"
	sdkkonnectcomp "github.com/Kong/sdk-konnect-go/models/components"
	"github.com/Kong/sdk-konnect-go/retry"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/sdk"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
)

// CreateClientCertificate creates a TLS client certificate and POSTs it to Konnect Control Plane configuration API
// so that KIC can use the certificates to authenticate against Konnect Admin API.
func CreateClientCertificate(ctx context.Context, t *testing.T, cpID string) (certPEM string, keyPEM string) {
	t.Helper()

	sdk := sdk.New(accessToken(), serverURLOpt(),
		sdkkonnectgo.WithRetryConfig(retry.Config{
			Backoff: &retry.BackoffStrategy{
				InitialInterval: 100,
				MaxInterval:     2000,
				Exponent:        1.2,
				MaxElapsedTime:  10000,
			},
		}),
	)

	cert, key := certificate.MustGenerateSelfSignedCertPEMFormat()

	t.Log("creating client certificate in Konnect")
	resp, err := sdk.DPCertificates.CreateDataplaneCertificate(ctx, cpID, &sdkkonnectcomp.DataPlaneClientCertificateRequest{
		Cert: string(cert),
	})
	require.NoError(t, err)

	if !assert.Equal(t, http.StatusCreated, resp.GetStatusCode()) {
		body, err := io.ReadAll(resp.RawResponse.Body)
		if err != nil {
			body = []byte(err.Error())
		}
		require.Failf(t, "failed creating client certificate", "body %s", body)
		return "", ""
	}

	return string(cert), string(key)
}
