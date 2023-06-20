//go:build integration_tests && knative

package integration

import (
	"context"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/knative"
)

// knativeVersion indicates the version of Knative to use in tests.
const knativeVersion = "1.10.0"

// knativeMinKubernetesVersion indicates the minimum Kubernetes version
// required in order to successfully run Knative tests.
var knativeMinKubernetesVersion = semver.MustParse("1.24.0")

func DeployAddonsForCluster(ctx context.Context, cluster clusters.Cluster) error {
	v, err := cluster.Version()
	if err != nil {
		return err
	}

	if v.GE(knativeMinKubernetesVersion) {
		// knative.Builder.WithVersion expects a git tag, not a semantic version, therefore we must provide
		// a tag in a format that Knative uses.
		knativeTag := fmt.Sprintf("knative-v%s", knativeVersion)
		builder, err := knative.NewBuilder().WithVersion(knativeTag)
		if err != nil {
			return err
		}
		fmt.Println("INFO: deploying knative addon")
		if err := env.Cluster().DeployAddon(ctx, builder.Build()); err != nil {
			return err
		}
	}

	return nil
}
