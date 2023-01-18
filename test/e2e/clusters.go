package e2e

import "strings"

const (
	// GKETestClusterNamePrefix defines a prefix that is used for naming GKE clusters created in tests.
	// It allows hack/cleanup_gke_clusters.go script to deterministically tell whether a cluster should
	// be taken into account in the cleanup procedure.
	gkeTestClusterNamePrefix = "e2e-"
)

func IsCreatedByE2ETests(clusterName string) bool {
	return strings.HasPrefix(clusterName, gkeTestClusterNamePrefix)
}
