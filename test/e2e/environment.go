package e2e

import "os"

var (
	// clusterVersionStr indicates the Kubernetes cluster version to use when
	// generating a testing environment and allows the caller to provide a specific
	// version. If no version is provided the default version for the cluster
	// provisioner in the testing framework will be used.
	clusterVersionStr = os.Getenv("KONG_CLUSTER_VERSION")

	// controllerImageOverride is the controller image to use in lieu of the default.
	controllerImageOverride = os.Getenv("TEST_KONG_CONTROLLER_IMAGE_OVERRIDE")

	// imageLoad is a boolean flag that indicates whether the controller and kong images should be loaded into the cluster.
	imageLoad = os.Getenv("TEST_KONG_LOAD_IMAGES")

	// kongImageOverride is the Kong image to use in lieu of the default.
	kongImageOverride     = os.Getenv("TEST_KONG_IMAGE_OVERRIDE")
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

// shouldLoadImages tells whether the controller and kong images should be loaded into the cluster.
func shouldLoadImages() bool {
	return imageLoad == "true"
}
