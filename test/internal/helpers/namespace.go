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

// Namespace provides the Namespace provisioned for each test case given their t.Name as the "testCase".
func Namespace(ctx context.Context, t *testing.T, env environments.Environment) *corev1.Namespace {
	namespace, err := clusters.GenerateNamespace(ctx, env.Cluster(), t.Name())
	require.NoError(t, err)
	t.Cleanup(func() { //nolint:contextcheck
		// Use context.Background() to ensure the namespace got removed when ctx
		// gets cancelled.
		assert.NoError(t, clusters.CleanupGeneratedResources(context.Background(), env.Cluster(), t.Name()))
	})

	return namespace
}
