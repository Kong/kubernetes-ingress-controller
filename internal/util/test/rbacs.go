package test

import (
	"context"
	"fmt"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"

	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
)

func DeployRBACsForCluster(ctx context.Context, cluster clusters.Cluster) error {
	fmt.Printf("INFO: deploying Kong RBACs to cluster\n")
	if err := clusters.KustomizeDeployForCluster(ctx, cluster, kongRBACsKustomize, "-n", consts.ControllerNamespace); err != nil {
		return err
	}

	fmt.Printf("INFO: deploying Kong gateway RBACs to cluster\n")
	if err := clusters.KustomizeDeployForCluster(ctx, cluster, kongGatewayRBACsKustomize, "-n", consts.ControllerNamespace); err != nil {
		return err
	}

	fmt.Printf("INFO: deploying Kong CRDs RBACs to cluster\n")
	return clusters.KustomizeDeployForCluster(ctx, cluster, kongCRDsRBACsKustomize, "-n", consts.ControllerNamespace)
}
