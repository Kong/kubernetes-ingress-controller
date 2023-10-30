package containers

import (
	"context"
	"strconv"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/kong/kubernetes-ingress-controller/v3/test"
)

// HTTPBin represents a docker container running the `kong/httpbin`.
type HTTPBin struct {
	container testcontainers.Container
}

// NewHTTPBin spawns a docker container running the `kong/httpbin` which can be used for proxy routing testing.
func NewHTTPBin(ctx context.Context, t *testing.T) HTTPBin {
	port, err := nat.NewPort("", strconv.Itoa(test.HTTPBinPort))
	require.NoError(t, err)
	req := testcontainers.ContainerRequest{
		Image:        test.HTTPBinImage,
		ExposedPorts: []string{MappedLocalPort(t, port)},
		WaitingFor:   wait.ForListeningPort(port),
	}
	httpBinC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, httpBinC.Terminate(ctx))
	})

	return HTTPBin{
		container: httpBinC,
	}
}

// IP returns the IP address of the HTTPBin container reachable from the docker network.
// Can be used to configure Kong services.
func (h HTTPBin) IP(ctx context.Context, t *testing.T) string {
	ip, err := h.container.ContainerIP(ctx)
	require.NoError(t, err)
	return ip
}
