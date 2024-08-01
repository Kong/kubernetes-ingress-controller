package helpers

import (
	"sync"
	"testing"

	"github.com/phayes/freeport"
)

var freePortLock = sync.Mutex{}

// GetFreePort asks the kernel for a free open port that is ready to use.
// On top of that, it also makes sure that the port hasn't been used in the current test run yet to reduce
// chances of a race condition in parallel tests.
func GetFreePort(t *testing.T) int {
	var (
		freePort    int
		retriesLeft = 100
	)
	freePortLock.Lock()
	defer freePortLock.Unlock()
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
		t.Fatal("no ports available")
	}
	return freePort
}

// userPorts keeps track of ports that were used in the current test run.
var usedPorts sync.Map
