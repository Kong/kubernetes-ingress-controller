package containers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
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

func KongWithDBMode(networkName string) KongOpt {
	return func(req *testcontainers.ContainerRequest) {
		req.Networks = []string{networkName}

		req.Env["KONG_DATABASE"] = "postgres"
		req.Env["KONG_PG_DATABASE"] = postgresDatabase
		req.Env["KONG_PG_USER"] = postgresUser
		req.Env["KONG_PG_PASSWORD"] = postgresPassword
		req.Env["KONG_PG_HOST"] = postgresContainerNetworkAlias
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
		Image: kongImageUnderTest(),
		ExposedPorts: []string{
			MappedLocalPort(t, kongAdminPort),
			MappedLocalPort(t, kongProxyPort),
		},
		Env: map[string]string{
			"KONG_DATABASE":      "off",
			"KONG_ADMIN_LISTEN":  fmt.Sprintf("0.0.0.0:%s", kongAdminPort),
			"KONG_PROXY_LISTEN":  fmt.Sprintf("0.0.0.0:%s", kongProxyPort),
			"KONG_ROUTER_FLAVOR": defaultRouterFlavor,
			"KONG_LICENSE_DATA":  os.Getenv("KONG_LICENSE_DATA"),
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
	t.Cleanup(func() { //nolint:contextcheck
		// If the container is already terminated, we don't need to terminate it again.
		if kongC.IsRunning() {
			assert.NoError(t, kongC.Terminate(context.Background()))
		}
	})

	adminURL, err := url.Parse(kong.AdminURL(ctx, t))
	require.NoError(t, err)

	const (
		tickTime = 100 * time.Millisecond
		waitTime = time.Minute
	)
	versionCtx, cancel := context.WithTimeout(ctx, waitTime)
	defer cancel()

	require.NoError(t,
		retry.Do(
			func() error {
				reqCtx, cancel := context.WithTimeout(ctx, test.RequestTimeout)
				defer cancel()
				kongVersion, err := helpers.ValidateMinimalSupportedKongVersion(reqCtx, adminURL, consts.KongTestPassword)
				if err != nil {
					return err
				}

				t.Logf("using Kong instance (version: %q) reachable at %s", kongVersion, adminURL)
				return nil
			},
			retry.Context(versionCtx),
			retry.Attempts(0),
			retry.Delay(tickTime),
			retry.DelayType(retry.FixedDelay),
			retry.LastErrorOnly(true),
			retry.OnRetry(func(_ uint, err error) {
				t.Logf("failed validating Kong version: %v", err)
			}),
			retry.RetryIf(func(err error) bool {
				return !errors.As(err, &helpers.TooOldKongGatewayError{})
			}),
		),
	)

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

func (c Kong) Terminate(ctx context.Context) error {
	return c.container.Terminate(ctx)
}

// kongImageUnderTest returns the Kong image to be used for integration tests. If both TEST_KONG_IMAGE and
// TEST_KONG_TAG are set, it will return the image and tag specified by them. Otherwise, it will return
// the default image (kong:latest or kong/kong-gateway if EE tests enabled).
func kongImageUnderTest() string {
	if testenv.KongImage() != "" && testenv.KongTag() != "" {
		return fmt.Sprintf("%s:%s", testenv.KongImage(), testenv.KongTag())
	}

	if testenv.KongEnterpriseEnabled() {
		return "kong/kong-gateway:latest"
	}
	return "kong:latest"
}
