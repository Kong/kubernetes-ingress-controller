// Package kongintegration contains integration tests that require a Kong instance to be running, but
// do not require a Kubernetes cluster nor full Kong Ingress Controller deployment.
//
// All tests that verify KIC's individual components compatibility with Kong should be placed here.
package kongintegration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

func init() {
	// Disable testcontainer's reaper (Ryuk) globally for this package. It's needed because Ryuk requires
	// privileged mode to run, which is not desired and could cause issues in CI.
	// Unfortunately, there is no way to disable it another way (e.g. via testcontainer's API), so we have
	// to use this hack.
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		panic(fmt.Errorf("failed to disable testcontainer's reaper (Ryuk): %w", err))
	}
}

// spawnDBLessKongContainer spawns a Kong container with DB-less mode enabled. It returns the admin API
// URL to be used to configure it.
// It sets up a cleanup function that will terminate the container when the test finishes.
func spawnDBLessKongContainer(ctx context.Context, t *testing.T) (adminURL string) {
	req := testcontainers.ContainerRequest{
		Image:        kongImageUnderTest(),
		ExposedPorts: []string{"8001"},
		Env: map[string]string{
			"KONG_DATABASE":     "off",
			"KONG_ADMIN_LISTEN": "0.0.0.0:8001",
		},
		WaitingFor: wait.ForListeningPort("8001"),
	}
	kongC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		err := kongC.Terminate(ctx)
		assert.NoError(t, err)
	})

	mappedAdminAPIPort, err := kongC.MappedPort(ctx, "8001")
	require.NoError(t, err)

	return fmt.Sprintf("http://localhost:%s", mappedAdminAPIPort.Port())
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
