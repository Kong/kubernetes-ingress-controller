package e2e

import (
	"cloud.google.com/go/container/apiv1/containerpb"
)

const (
	gkeTestClusterLabel = "test-cluster"
	gkeLabelValueTrue   = "true"
)

// IsGKETestCluster tells if the GKE cluster has been created for test purposes.
func IsGKETestCluster(cluster *containerpb.Cluster) bool {
	if cluster == nil {
		return false
	}

	if labels := cluster.GetResourceLabels(); labels != nil {
		return labels[gkeTestClusterLabel] == gkeLabelValueTrue
	}

	return false
}

func gkeTestClusterLabels() map[string]string {
	return map[string]string{
		gkeTestClusterLabel: gkeLabelValueTrue,
	}
}
