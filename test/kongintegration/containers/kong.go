package containers

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

const (
	kongAdminPort       = "8001"
	kongProxyPort       = "8000"
	defaultRouterFlavor = "expressions"
)

type KongOpt func(*testcontainers.ContainerRequest)

func KongWithRouterFlavor(flavor string) KongOpt {
	return func(req *testcontainers.ContainerRequest) {
		req.Env["KONG_ROUTER_FLAVOR"] = flavor
	}
}

// Kong represents a docker container running Kong.
type Kong struct {
	container testcontainers.Container
}

// NewKong spawns a docker container running Kong (its image is determined by environment variables).
// It sets up a cleanup function that will terminate the container when the test finishes.
func NewKong(ctx context.Context, t *testing.T, opts ...KongOpt) Kong {
	req := testcontainers.ContainerRequest{
		Image:        kongImageUnderTest(),
		ExposedPorts: []string{kongAdminPort, kongProxyPort},
		Env: map[string]string{
			"KONG_DATABASE":      "off",
			"KONG_ADMIN_LISTEN":  fmt.Sprintf("0.0.0.0:%s", kongAdminPort),
			"KONG_PROXY_LISTEN":  fmt.Sprintf("0.0.0.0:%s", kongProxyPort),
			"KONG_ROUTER_FLAVOR": defaultRouterFlavor,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(kongAdminPort),
			wait.ForListeningPort(kongProxyPort),
		),
	}
	for _, opt := range opts {
		opt(&req)
	}
	kongC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	kong := Kong{
		container: kongC,
	}
	adminURL, err := url.Parse(kong.AdminURL(ctx, t))
	require.NoError(t, err)

	kongVersion, err := helpers.ValidateMinimalSupportedKongVersion(adminURL, consts.KongTestPassword) //nolint:contextcheck
	require.NoError(t, err)
	fmt.Printf("INFO: using Kong instance (version: %q) reachable at %s\n", kongVersion, adminURL)

	t.Cleanup(func() {
		assert.NoError(t, kongC.Terminate(ctx))
	})

	return kong
}

// AdminURL returns the admin API URL of the Kong container reachable from the host machine.
func (c Kong) AdminURL(ctx context.Context, t *testing.T) string {
	port, err := c.container.MappedPort(ctx, kongAdminPort)
	require.NoError(t, err)
	return fmt.Sprintf("http://localhost:%s", port.Port())
}

// ProxyURL returns the proxy URL of the Kong container reachable from the host machine.
func (c Kong) ProxyURL(ctx context.Context, t *testing.T) string {
	port, err := c.container.MappedPort(ctx, kongProxyPort)
	require.NoError(t, err)
	return fmt.Sprintf("http://localhost:%s", port.Port())
}

// kongImageUnderTest returns the Kong image to be used for integration tests. If both TEST_KONG_IMAGE and
// TEST_KONG_TAG are set, it will return the image and tag specified by them. Otherwise, it will return
// the default image (kong:latest).
func kongImageUnderTest() string {
	if testenv.KongImage() != "" && testenv.KongTag() != "" {
		return fmt.Sprintf("%s:%s", testenv.KongImage(), testenv.KongTag())
	}

	return "kong:latest"
}
