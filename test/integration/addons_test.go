//go:build integration_tests

package integration

import (
	"context"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
)

func DeployAddonsForCluster(_ context.Context, _ clusters.Cluster) error {
	return nil
}
