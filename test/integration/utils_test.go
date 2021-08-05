//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

// -----------------------------------------------------------------------------
// Testing Timeouts
// -----------------------------------------------------------------------------

const (
	// waitTick is the default timeout tick interval for checking on ingress resources.
	waitTick = time.Second * 1

	// ingressWait is the default amount of time to wait for any particular ingress resource to be provisioned.
	ingressWait = time.Minute * 3

	// httpcTimeout is the default client timeout for HTTP clients used in tests.
	httpcTimeout = time.Second * 3

	// environmentCleanupTimeout is the amount of time that will be given by the test suite to the
	// testing environment to perform its cleanup when the test suite is shutting down.
	environmentCleanupTimeout = time.Minute * 3
)

// -----------------------------------------------------------------------------
// Testing Variables
// -----------------------------------------------------------------------------

var (
	// ctx the topical context of the test suite, can be used by test cases if they don't need
	// any special context as a function of the test
	ctx context.Context

	// cancel is the cancel function for the above global test context
	cancel context.CancelFunc

	// httpBinImage is the container image name we use for deploying the "httpbin" HTTP testing tool.
	// if you need a simple HTTP server for tests you're writing, use this and check the documentation.
	// See: https://github.com/postmanlabs/httpbin
	httpBinImage = "kennethreitz/httpbin"

	// ingressClass indicates the ingress class name which the tests will use for supported object reconciliation
	ingressClass = "kongtests"

	// elsewhere is the name of an alternative namespace
	elsewhere = "elsewhere"

	// controllerNamespace is the Kubernetes namespace where the controller is deployed
	controllerNamespace = "kong-system"

	// httpc is the default HTTP client to use for tests
	httpc = http.Client{Timeout: httpcTimeout}

	// watchNamespaces is a list of namespaces the controller watches
	watchNamespaces = strings.Join([]string{
		elsewhere,
		corev1.NamespaceDefault,
		testIngressEssentialsNamespace,
		testIngressClassNameSpecNamespace,
		testIngressHTTPSNamespace,
		testIngressHTTPSRedirectNamespace,
		testBulkIngressNamespace,
		testTCPIngressNamespace,
		testUDPIngressNamespace,
		testPluginsNamespace,
	}, ",")

	// env is the primary testing environment object which includes access to the Kubernetes cluster
	// and all the addons deployed in support of the tests.
	env environments.Environment

	// proxyURL provides access to the proxy endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyURL *url.URL

	// proxyAdminURL provides access to the Admin API endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyAdminURL *url.URL

	// proxyUDPURL provides access to the UDP API endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyUDPURL *url.URL

	// clusterVersion is a convenience var where the found version of the env.Cluster is stored.
	clusterVersion semver.Version
)

// -----------------------------------------------------------------------------
// Testing Variables - Environment Overrides
// -----------------------------------------------------------------------------

var (
	// dbmode indicates the database backend of the test cluster ("off" and "postgres" are supported)
	dbmode = os.Getenv("TEST_DATABASE_MODE")

	// clusterVersion indicates the version of Kubernetes to use for the tests (if the cluster was not provided by the caller)
	clusterVersionStr = os.Getenv("KONG_CLUSTER_VERSION")

	// existingCluster indicates whether or not the caller is providing their own cluster for running the tests.
	// These need to come in the format <TYPE>:<NAME> (e.g. "kind:<NAME>", "gke:<NAME>", e.t.c.).
	existingCluster = os.Getenv("KONG_TEST_CLUSTER")

	// keepTestCluster indicates whether the caller wants the cluster created by the test suite
	// to persist after the test for inspection. This has a nil effect when an existing cluster
	// is provided, as cleanup is not performed for existing clusters.
	keepTestCluster = os.Getenv("KONG_TEST_CLUSTER_PERSIST")

	// maxBatchSize indicates the maximum number of objects that should be POSTed per second during testing
	maxBatchSize = determineMaxBatchSize()
)

// -----------------------------------------------------------------------------
// Test Suite Exit Codes
// -----------------------------------------------------------------------------

const (
	// ExitCodeIncompatibleOptions is a POSIX compliant exit code for the test suite to
	// indicate that some combination of provided configurations were not compatible.
	ExitCodeIncompatibleOptions = 100

	// ExitCodeInvalidOptions is a POSIX compliant exit code for the test suite to indicate
	// that some of the provided runtime options were not valid and the tests could not run.
	ExitCodeInvalidOptions = 101

	// ExitCodeCantUseExistingCluster is a POSIX compliant exit code for the test suite to
	// indicate that an existing cluster provided for the tests was not usable.
	ExitCodeCantUseExistingCluster = 101

	// ExitCodeCantCreateCluster is a POSIX compliant exit code for the test suite to indicate
	// that a failure occurred when trying to create a Kubernetes cluster to run the tests.
	ExitCodeCantCreateCluster = 102

	// ExitCodeCleanupFailed is a POSIX compliant exit code for the test suite to indicate
	// that a failure occurred during cluster cleanup.
	ExitCodeCleanupFailed = 103

	// ExitCodeEnvSetupFailed is a generic exit code that can be used as a fallback for general
	// problems setting up the testing environment and/or cluster.
	ExitCodeEnvSetupFailed = 104
)

// -----------------------------------------------------------------------------
// Testing Utility Functions
// -----------------------------------------------------------------------------

// expect404WithNoRoute is used to check whether a given http response is (specifically) a Kong 404.
func expect404WithNoRoute(t *testing.T, proxyURL string, resp *http.Response) bool {
	if resp.StatusCode == http.StatusNotFound {
		// once the route is torn down and returning 404's, ensure that we got the expected response body back from Kong
		// Expected: {"message":"no Route matched with those values"}
		b := new(bytes.Buffer)
		_, err := b.ReadFrom(resp.Body)
		require.NoError(t, err)
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

// determineMaxBatchSize provides a size limit for the number of resources to POST in a single second during tests, and can be overridden with an ENV var if desired.
func determineMaxBatchSize() int {
	if v := os.Getenv("KONG_BULK_TESTING_BATCH_SIZE"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("Error: invalid batch size %s: %s", v, err))
		}
		return i
	}
	return 50
}

// -----------------------------------------------------------------------------
// Test.MAIN Utility Functions
// -----------------------------------------------------------------------------

// exitOnErrWithCode is a helper function meant for us in the test.Main to simplify failing and exiting
// the tests under unrecoverable error conditions. It will also attempt to perform any cluster
// cleaning necessary before exiting.
func exitOnErrWithCode(err error, exitCode int) {
	if err == nil {
		return
	}

	fmt.Println("WARNING: failure occurred, performing test cleanup")
	if env != nil && existingCluster == "" && keepTestCluster == "" {
		ctx, cancel := context.WithTimeout(context.Background(), environmentCleanupTimeout)
		defer cancel()

		fmt.Printf("INFO: cluster %s is being deleted\n", env.Cluster().Name())
		if cleanupErr := env.Cleanup(ctx); cleanupErr != nil {
			err = fmt.Errorf("cleanup failed after test failure occurred CLEANUP_FAILURE=(%s) TEST_FAILURE=(%s)", cleanupErr, err)
		}
	}

	fmt.Fprintf(os.Stderr, "Error: tests failed: %s\n", err)
	os.Exit(exitCode)
}

// exitOnErr is a wrapper around exitOnErrorWithCode that defaults to using the ExitCodeEnvSetupFailed
// exit code. This function is meant for convenience to wrap errors in setup that are hard to predict.
func exitOnErr(err error) {
	if err == nil {
		return
	}
	exitOnErrWithCode(err, ExitCodeEnvSetupFailed)
}
