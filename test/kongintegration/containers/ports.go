package containers

import (
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"

	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

// MappedLocalPort returns a port mapping for a container port that can be used to access the
// container from the host. The returned string is in the format expected by testcontainers.ContainerRequest ExposedPorts.
//
// This is a workaround for a bug in docker that causes IPv4 and IPv6 port mappings to overlap with
// different containers on the same host. See: https://github.com/moby/moby/issues/42442.
func MappedLocalPort(t *testing.T, containerPort nat.Port) string {
	// Return the port mapping in the format expected by testcontainers.ContainerRequest ExposedPorts:
	// <host>:<host-port>:<container-port>.
	return fmt.Sprintf("0.0.0.0:%d:%s", helpers.GetFreePort(t), containerPort.Port())
}
