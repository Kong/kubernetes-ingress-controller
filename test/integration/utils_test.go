//go:build integration_tests

package integration

import (
	"net/url"

	"github.com/kong/kubernetes-testing-framework/pkg/environments"
)

// -----------------------------------------------------------------------------
// Testing Variables
// -----------------------------------------------------------------------------

var (
	// env is the primary testing environment object which includes access to the Kubernetes cluster
	// and all the addons deployed in support of the tests.
	env environments.Environment

	// proxyHTTPURL provides access to the proxy endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyHTTPURL *url.URL

	// proxyHTTPSURL provides access to the proxy endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyHTTPSURL *url.URL

	// proxyAdminURL provides access to the Admin API endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyAdminURL *url.URL

	// proxyTCPURL provides access to the TCP API endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyTCPURL string

	// proxyTLSURL provides access to the TLS API endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyTLSURL string

	// proxyUDPURL provides access to the UDP API endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyUDPURL string

	// runInvalidConfigTests is set to true to run the test cases including invalid test cases.
	runInvalidConfigTests bool
)
