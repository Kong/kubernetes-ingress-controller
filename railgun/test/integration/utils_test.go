//+build integration_tests

package integration

import (
	"net/url"
	"sync"
)

var (
	l = sync.RWMutex{}
	u *url.URL
)

// proxyURL is a threadsafe way to wait for the proxy to be ready
// and then receive the URL where it can be reached.
func proxyURL() *url.URL {
	l.Lock()
	defer l.Unlock()

	if u == nil {
		u = <-proxyReady
	}

	return u
}
