//+build integration_tests

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"testing"

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

// expect404WithNoRoute is used to check whether a given http response is (specifically) a Kong 404.
func expect404WithNoRoute(t *testing.T, proxyURL string, resp *http.Response) bool {
	if resp.StatusCode == http.StatusNotFound {
		// once the route is torn down and returning 404's, ensure that we got the expected response body back from Kong
		// Expected: {"message":"no Route matched with those values"}
		b := new(bytes.Buffer)
		b.ReadFrom(resp.Body)
		body := struct {
			Message string `json:"message"`
		}{}
		if err := json.Unmarshal(b.Bytes(), &body); err != nil {
			t.Logf("WARNING: error decoding JSON from proxy while waiting for %s: %v", proxyURL, err)
			return false
		}
		return body.Message == "no Route matched with those values"
	}
	return false
}
