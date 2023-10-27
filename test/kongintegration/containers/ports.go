package containers

import (
	"fmt"
	"sync"

	"github.com/docker/go-connections/nat"
	"github.com/phayes/freeport"
)

// MappedLocalPort returns a port mapping for a container port that can be used to access the
// container from the host. The returned string is in the format expected by testcontainers.ContainerRequest ExposedPorts.
//
// This is a workaround for a bug in docker that causes IPv4 and IPv6 port mappings to overlap with
// different containers on the same host. See: https://github.com/moby/moby/issues/42442.
func MappedLocalPort(containerPort nat.Port) string {
	var (
		freePort    int
		retriesLeft = 100
	)
	for {
		// Get a random free port, but do not use it yet...
		var err error
		freePort, err = freeport.GetFreePort()
		if err != nil {
			continue
		}

		// ... First, check if the port has been used in this test run already to reduce chances of a race condition.
		_, wasUsed := usedPorts.LoadOrStore(freePort, true)

		// The port hasn't been used in this test run - we can use it. It was stored in usedPorts, so it will not be
		// used again during this test run.
		if !wasUsed {
			break
		}

		// Otherwise, the port was used in this test run. We need to get another one.
		freePort = 0
		retriesLeft--
		if retriesLeft == 0 {
			break
		}
	}
	if freePort == 0 {
		panic("no ports available")
	}

	// Return the port mapping in the format expected by testcontainers.ContainerRequest ExposedPorts:
	// <host>:<host-port>:<container-port>.
	return fmt.Sprintf("0.0.0.0:%d:%s", freePort, containerPort.Port())
}

// userPorts keeps track of ports that were used in the current test run.
var usedPorts sync.Map
