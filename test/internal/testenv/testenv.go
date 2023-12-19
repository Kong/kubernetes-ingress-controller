package testenv

import (
	"fmt"
	"os"
	"time"
)

// -----------------------------------------------------------------------------
// Testing Helpers for Environment Variables Overrides
// -----------------------------------------------------------------------------

type DBModeDatabase string

const (
	DBModeOff      DBModeDatabase = "off"
	DBModePostgres DBModeDatabase = "postgres"
)

// DBMode indicates the database backend of the test cluster ("off" and "postgres" are supported).
func DBMode() DBModeDatabase {
	switch dbmode := os.Getenv("TEST_DATABASE_MODE"); dbmode {
	case "", "off":
		return DBModeOff
	case "postgres":
		return DBModePostgres
	default:
		// TODO
		os.Exit(1)
		return ""
	}
}

// KongImage is the Kong image to use in lieu of the default.
func KongImage() string {
	return os.Getenv("TEST_KONG_IMAGE")
}

// KongTag is the Kong image tag to use in tests.
func KongTag() string {
	return os.Getenv("TEST_KONG_TAG")
}

// KongImageTag is the combined Kong image and tag if both are set, or empty string if not.
func KongImageTag() string {
	if KongImage() != "" && KongTag() != "" {
		return fmt.Sprintf("%s:%s", KongImage(), KongTag())
	}
	return ""
}

// KongEffectiveVersion is the effective semver of kong gateway.
// When testing against "nightly" image of kong gateway, we need to set the effective version for parsing semver in chart templates.
func KongEffectiveVersion() string {
	return os.Getenv("TEST_KONG_EFFECTIVE_VERSION")
}

// KongRouterFlavor returns router mode of Kong in tests. Currently supports:
// - `traditional`
// - `traditional_compatible`.
// - `expressions` (experimental, only for testing expression route related tests).
func KongRouterFlavor() string {
	rf := os.Getenv("TEST_KONG_ROUTER_FLAVOR")
	if rf != "" && rf != "traditional" && rf != "traditional_compatible" && rf != "expressions" {
		// TODO
		os.Exit(1)
	}

	return rf
}

// KongPullUsername is the Docker username to use for the Kong image pull secret.
func KongPullUsername() string {
	return os.Getenv("TEST_KONG_PULL_USERNAME")
}

// KongPullPassword is the Docker password to use for the Kong image pull secret.
func KongPullPassword() string {
	return os.Getenv("TEST_KONG_PULL_PASSWORD")
}

// KongEnterpriseEnabled enables Enterprise-specific tests when set to "true".
func KongEnterpriseEnabled() bool {
	return os.Getenv("TEST_KONG_ENTERPRISE") == "true"
}

// ClusterVersion indicates the version of Kubernetes to use for the tests
// (if the cluster was not provided by the caller).
func ClusterVersion() string {
	return os.Getenv("KONG_CLUSTER_VERSION")
}

// ControllerFeatureGates contains the feature gates that should be enabled
// for test runs in the controller.
// If none specified, we fall back to default values.
func ControllerFeatureGates() string {
	featureGates := os.Getenv("KONG_CONTROLLER_FEATURE_GATES")
	if featureGates == "" {
		featureGates = getFeatureGates()
	}
	return featureGates
}

// ExpressionRoutesEnabled indicates whether or not to enable expression routes
// for the Kong Gateway and the controller.
// If none specified, we fall back to default value - traditional_compatible.
func ExpressionRoutesEnabled() bool {
	return os.Getenv("KONG_TEST_EXPRESSION_ROUTES") == "true"
}

// -----------------------------------------------------------------------------
// Environment variables related helpers
// -----------------------------------------------------------------------------

// KeepTestCluster indicates whether the caller wants the cluster created by the test suite
// to persist after the test for inspection. This has a nil effect when an existing cluster
// is provided, as cleanup is not performed for existing clusters.
func KeepTestCluster() string {
	return os.Getenv("KONG_TEST_CLUSTER_PERSIST")
}

// ExistingClusterName indicates whether or not the caller is providing their own cluster for running the tests.
// These need to come in the format <TYPE>:<NAME> (e.g. "kind:<NAME>", "gke:<NAME>", e.t.c.).
func ExistingClusterName() string {
	return os.Getenv("KONG_TEST_CLUSTER")
}

// WaitForClusterDelete indicates whether or not to wait for cluster deletion to complete.
func WaitForClusterDelete() bool {
	return os.Getenv("KONG_TEST_CLUSTER_WAIT_FOR_DELETE") == "true"
}

// EnvironmentReadyTimeout returns the amount of time that will be given to wait for the environment
// ready, including all the dependencies (kong, metallb, etc)
// used here to make up a context to pass into environments.WaitForReady to trigger cleanup when timed out.
func EnvironmentReadyTimeout() time.Duration {
	const DefaultEnvironmentReadyTimeout = time.Minute * 10
	timeout, err := time.ParseDuration(os.Getenv("KONG_TEST_ENVIRONMENT_READY_TIMEOUT"))
	if err != nil {
		timeout = DefaultEnvironmentReadyTimeout
	}
	return timeout
}

// IsCI indicates whether or not the tests are running in a CI environment.
func IsCI() bool {
	// It's a common convention that e.g. GitHub, GitLab, and other CI providers
	// set the CI environment variable.
	return os.Getenv("CI") == "true"
}
