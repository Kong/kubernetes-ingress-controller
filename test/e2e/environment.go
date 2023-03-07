package e2e

import "os"

var (
	// clusterVersionStr indicates the Kubernetes cluster version to use when
	// generating a testing environment and allows the caller to provide a specific
	// version. If no version is provided the default version for the cluster
	// provisioner in the testing framework will be used.
	clusterVersionStr = os.Getenv("KONG_CLUSTER_VERSION")

	imageOverride         = os.Getenv("TEST_KONG_CONTROLLER_IMAGE_OVERRIDE")
	imageLoad             = os.Getenv("TEST_KONG_CONTROLLER_IMAGE_LOAD")
	kongImageOverride     = os.Getenv("TEST_KONG_IMAGE_OVERRIDE")
	kongImageLoad         = os.Getenv("TEST_KONG_IMAGE_LOAD")
	kongImagePullUsername = os.Getenv("TEST_KONG_PULL_USERNAME")
	kongImagePullPassword = os.Getenv("TEST_KONG_PULL_PASSWORD")

	// KONG_TEST_CLUSTER is to be filled when an already existing cluster should be used
	// in tests. It should be in a `<gke|kind>:<name>` format.
	// It takes precedence over KONG_TEST_CLUSTER_PROVIDER.
	existingCluster = os.Getenv("KONG_TEST_CLUSTER")

	// KONG_TEST_CLUSTER_PROVIDER is to be filled when a cluster of a given kind should
	// be created in tests. It can be either `gke` or `kind`.
	// It's not used when KONG_TEST_CLUSTER is set.
	clusterProvider = os.Getenv("KONG_TEST_CLUSTER_PROVIDER")
)
