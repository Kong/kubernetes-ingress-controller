//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"

	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

// -----------------------------------------------------------------------------
// Testing Variables
// -----------------------------------------------------------------------------

var (
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

	// runInvalidConfigTests is set to true to run the test cases including invalid test cases.
	runInvalidConfigTests bool
)

// -----------------------------------------------------------------------------
// Test.MAIN Utility Functions
// -----------------------------------------------------------------------------

// exitOnErrWithCode is a helper function meant for us in the test.Main to simplify failing and exiting
// the tests under unrecoverable error conditions. It will also attempt to perform any cluster
// cleaning necessary before exiting.
func exitOnErrWithCode(ctx context.Context, err error, exitCode int) {
	if err == nil {
		return
	}

	fmt.Println("WARNING: failure occurred, performing test cleanup")
	if env != nil && testenv.ExistingClusterName() == "" && testenv.KeepTestCluster() == "" {
		ctx, cancel := context.WithTimeout(ctx, environmentCleanupTimeout)
		defer cancel()

		fmt.Printf("INFO: cluster %s is being deleted\n", env.Cluster().Name())
		if cleanupErr := env.Cleanup(ctx); cleanupErr != nil {
			err = fmt.Errorf("cleanup failed after test failure occurred CLEANUP_FAILURE=(%s): %w", cleanupErr, err)
		}
	}

	fmt.Fprintf(os.Stderr, "Error: tests failed: %s\n", err)
	os.Exit(exitCode)
}

// exitOnErr is a wrapper around exitOnErrorWithCode that defaults to using the ExitCodeEnvSetupFailed
// exit code. This function is meant for convenience to wrap errors in setup that are hard to predict.
func exitOnErr(ctx context.Context, err error) {
	if err == nil {
		return
	}
	exitOnErrWithCode(ctx, err, consts.ExitCodeEnvSetupFailed)
}
