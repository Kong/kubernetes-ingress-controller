package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

// TeardownCluster dumps the diagnostics from the test cluster if the test failed
// and performs a cluster removal.
func TeardownCluster(ctx context.Context, t *testing.T, cluster clusters.Cluster) {
	t.Helper()

	DumpDiagnosticsIfFailed(ctx, t, cluster)
	const environmentCleanupTimeout = 10 * time.Minute
	ctx, cancel := context.WithTimeout(ctx, environmentCleanupTimeout)
	defer cancel()

	assert.NoError(t, RemoveCluster(ctx, cluster))
}

// RemoveCluster removes the cluster if it was created by the test suite.
// Pass desired timeout through context.
func RemoveCluster(ctx context.Context, cluster clusters.Cluster) error {
	if testenv.KeepTestCluster() == "" && testenv.ExistingClusterName() == "" {
		fmt.Printf("INFO: cluster %s is being deleted\n", cluster.Name())
		return cluster.Cleanup(ctx)
	}
	return nil
}

// DumpDiagnosticsIfFailed dumps the diagnostics if the test failed.
func DumpDiagnosticsIfFailed(ctx context.Context, t *testing.T, cluster clusters.Cluster) {
	t.Helper()

	if t.Failed() {
		output, err := cluster.DumpDiagnostics(ctx, t.Name())
		t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
		assert.NoError(t, err)
	}
}
