package envtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
	"github.com/kong/kubernetes-ingress-controller/v3/test/kongintegration/containers"
)

// runKongEnterprise runs a Kong EE container using the version from `test_dependencies.yaml`.
func runKongEnterprise(ctx context.Context, t *testing.T) containers.Kong {
	// Get the Kong Gateway version to use for the test from `test_dependencies.yaml` file.
	gatewayTag, err := testenv.GetDependencyVersion("envtests.kong-ee")
	require.NoError(t, err)

	// Prepare the container config modifier to set the Kong Gateway version.
	withEnvtestsVersion := func(request *testcontainers.ContainerRequest) {
		request.Image = fmt.Sprintf("kong/kong-gateway:%s", gatewayTag)
	}

	return containers.NewKong(ctx, t, withEnvtestsVersion)
}

// runKongGatewayWithoutStickySessionsSupport runs a Kong Gateway container that does
// not support sticky sessions, since every version >= 3.11.0 supports sticky sessions.
func runKongGatewayWithoutStickySessionsSupport(ctx context.Context, t *testing.T) containers.Kong {
	// Get the Kong Gateway version to use for the test from `test_dependencies.yaml` file.
	gatewayTag, err := testenv.GetDependencyVersion("envtests.kong-without-sticky-sessions")
	require.NoError(t, err)

	// Prepare the container config modifier to set the Kong Gateway version.
	withEnvtestsVersion := func(request *testcontainers.ContainerRequest) {
		request.Image = fmt.Sprintf("kong/kong:%s", gatewayTag)
	}

	return containers.NewKong(ctx, t, withEnvtestsVersion)
}
