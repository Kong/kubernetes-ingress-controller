//go:build integration_tests && knative
// +build integration_tests,knative

package integration

import (
	"context"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/knative"
)

// knativeMinKubernetesVersion indicates the minimum Kubernetes version
// required in order to successfully run Knative tests.
var knativeMinKubernetesVersion = semver.MustParse("1.22.0")

func DeployAddonsForCluster(ctx context.Context, cluster clusters.Cluster) error {
	v, err := cluster.Version()
	if err != nil {
		return err
	}

	if v.GE(knativeMinKubernetesVersion) {
		knativeAddon := knative.NewBuilder().Build()
		fmt.Println("INFO: deploying knative addon")
		err := env.Cluster().DeployAddon(ctx, knativeAddon)
		if err != nil {
			return err
		}
	}

	return nil
}
