package e2e

import "os"

var (
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

	// githubServerURL, githubRepo, githubRunID are used to locate the run of github wokflow.
	// See: https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
	githubServerURL = os.Getenv("GITHUB_SERVER_URL")
	githubRepo      = os.Getenv("GITHUB_REPOSITORY")
	githubRunID     = os.Getenv("GITHUB_RUN_ID")
)

// shouldLoadImages tells whether the controller and kong images should be loaded into the cluster.
func shouldLoadImages() bool {
	return imageLoad == "true"
}
