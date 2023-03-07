package e2e

import (
	"testing"

	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/require"
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

func logClusterInfo(t *testing.T, cluster clusters.Cluster) {
	t.Helper()

	v, err := cluster.Version()
	require.NoError(t, err)
	t.Logf("cluster %s (type: %s, v: %s) is up", cluster.Name(), cluster.Type(), v)
}
