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
	existingCluster       = os.Getenv("KONG_TEST_CLUSTER")
)
