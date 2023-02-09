package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

// TeardownCluster dumps the diagnostics from the test cluster if the test failed
// and performs a cluster teardown.
func TeardownCluster(ctx context.Context, t *testing.T, cluster clusters.Cluster) {
	t.Helper()

	const (
		environmentCleanupTimeout = 3 * time.Minute
	)

	DumpDiagnosticsIfFailed(ctx, t, cluster)

	if testenv.KeepTestCluster() == "" && testenv.ExistingClusterName() == "" {
		ctx, cancel := context.WithTimeout(ctx, environmentCleanupTimeout)
		defer cancel()
		t.Logf("INFO: cluster %s is being deleted\n", cluster.Name())
		assert.NoError(t, cluster.Cleanup(ctx))
		return
	}
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
