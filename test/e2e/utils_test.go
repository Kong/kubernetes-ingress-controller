//+build e2e_tests

package e2e

import (
	"net/http"
	"os"
	"time"
)

var (
	// clusterVersionStr indicates the Kubernetes cluster version to use when
	// generating a testing environment and allows the caller to provide a specific
	// version. If no version is provided the default version for the cluster
	// provisioner in the testing framework will be used.
	clusterVersionStr = os.Getenv("KONG_CLUSTER_VERSION")

	// httpc is a standard HTTP client for tests to use that has a low default
	// timeout instead of the longer default provided by the http stdlib.
	httpc = http.Client{Timeout: time.Second * 10}
)
