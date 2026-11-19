package containers

import (
	"net/netip"
	"strconv"
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

// FreeHostPortBindings returns a host config modifier that binds each of the given container ports to a free
// port on 0.0.0.0 of the host. The container ports must also be listed in testcontainers.ContainerRequest ExposedPorts.
//
// This is a workaround for a bug in docker that causes IPv4 and IPv6 port mappings to overlap with
// different containers on the same host. See: https://github.com/moby/moby/issues/42442.
func FreeHostPortBindings(t *testing.T, containerPorts ...string) func(*container.HostConfig) {
	bindings := make(network.PortMap, len(containerPorts))
	for _, containerPort := range containerPorts {
		port, err := network.ParsePort(containerPort)
		require.NoError(t, err)
		bindings[port] = []network.PortBinding{{
			HostIP:   netip.IPv4Unspecified(),
			HostPort: strconv.Itoa(helpers.GetFreePort(t)),
		}}
	}
	return func(hostConfig *container.HostConfig) {
		hostConfig.PortBindings = bindings
	}
}
