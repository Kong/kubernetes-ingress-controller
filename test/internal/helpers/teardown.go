//go:build integration_tests || e2e_tests
// +build integration_tests e2e_tests

package helpers

import (
	"context"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
)

// TeardownCluster dumps the diagnostics from the test cluster if the test failed
// and performs a cluster teardown.
func TeardownCluster(ctx context.Context, t *testing.T, cluster clusters.Cluster) {
	t.Helper()

	DumpDiagnosticsIfFailed(ctx, t, cluster)
	assert.NoError(t, cluster.Cleanup(ctx))
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
