//go:build integration_tests && !knative
// +build integration_tests,!knative

package integration

import (
	"context"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
)

func DeployAddonsForCluster(ctx context.Context, cluster clusters.Cluster) error {
	return nil
}
