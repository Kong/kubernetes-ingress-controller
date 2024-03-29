package helpers

import (
	"context"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

// Setup is a test helper function which:
//   - creates a cluster cleaner which will be used to to clean up test resources
//     automatically after the test finishes and creates a new namespace for the test
//     to use.
//   - creates a namespace for the provided test and adds it to the cleaner for
//     automatic cleanup using the previously created cleaner.
//
// TODO move this into a shared library https://github.com/Kong/kubernetes-testing-framework/issues/302
func Setup(ctx context.Context, t *testing.T, env environments.Environment) (*corev1.Namespace, *clusters.Cleaner) {
	t.Helper()

	t.Log("performing test setup")
	cluster := env.Cluster()
	cleaner := clusters.NewCleaner(cluster)
	t.Cleanup(func() { //nolint:contextcheck
		// We still want to dump the diagnostics and perform the cleanup so use
		// a separate context.
		ctx := context.Background()
		t.Logf("Start cleanup for test %s", t.Name())
		dumpDiagnosticsIfFailed(ctx, t, cluster)
		assert.NoError(t, cleaner.Cleanup(ctx))
	})

	t.Log("creating a testing namespace")
	namespace, err := clusters.GenerateNamespace(ctx, cluster, LabelValueForTest(t))
	require.NoError(t, err)
	cleaner.AddNamespace(namespace)
	t.Logf("created namespace %s for test case %s", namespace.Name, t.Name())

	return namespace, cleaner
}
