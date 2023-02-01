//go:build integration_tests

package adminapi_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

func TestNewKongClientForKonnectRuntimeGroup(t *testing.T) {
	t.Skip("There's no infrastructure for Konnect tests yet")

	ctx := context.Background()
	const runtimeGroupID = "adf78c28-5763-4394-a9a4-a9436a1bea7d"

	c, err := adminapi.NewKongClientForKonnectRuntimeGroup(ctx, adminapi.KonnectConfig{
		ConfigSynchronizationEnabled: true,
		RuntimeGroupID:               runtimeGroupID,
		Address:                      "https://us.kic.api.konghq.tech",
		TLSClient: adminapi.TLSClientConfig{
			Cert: os.Getenv("KONG_TEST_KONNECT_TLS_CLIENT_CERT"),
			Key:  os.Getenv("KONG_TEST_KONNECT_TLS_CLIENT_KEY"),
		},
	})
	require.NoError(t, err)

	require.True(t, c.IsKonnect())
	require.Equal(t, runtimeGroupID, c.KonnectRuntimeGroup())

	_, err = c.AdminAPIClient().Services.ListAll(ctx)
	require.NoError(t, err)
}
