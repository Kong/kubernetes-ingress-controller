//+build integration_tests

package integration

import (
	"sync"

	ktfkind "github.com/kong/kubernetes-testing-framework/pkg/kind"
)

var (
	l            = sync.RWMutex{}
	proxyReadyCh = make(chan ktfkind.ProxyReadinessEvent)

	readinessEvent *ktfkind.ProxyReadinessEvent
)

// proxyReady is a threadsafe way to wait for the proxy to be ready
// and then receive the URLs where it can be reached.
func proxyReady() ktfkind.ProxyReadinessEvent {
	l.Lock()
	defer l.Unlock()

	if readinessEvent == nil {
		event := <-proxyReadyCh
		readinessEvent = &event
	}

	return *readinessEvent
}
