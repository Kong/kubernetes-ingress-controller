package test

import (
	"context"
	"fmt"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
)

const (
	kongRBACsKustomize        = "../../config/rbac/"
	kongKnativeRBACsKustomize = "../../config/rbac/knative"
	kongGatewayRBACsKustomize = "../../config/rbac/gateway"
	kongCRDsRBACsKustomize    = "../../config/rbac/crds"
)

func DeployRBACsForCluster(ctx context.Context, cluster clusters.Cluster) error {
	fmt.Printf("INFO: deploying Kong RBACs to cluster")
	if err := clusters.KustomizeDeployForCluster(ctx, cluster, kongRBACsKustomize); err != nil {
		return err
	}

	fmt.Printf("INFO: deploying Kong knative RBACs to cluster")
	if err := clusters.KustomizeDeployForCluster(ctx, cluster, kongKnativeRBACsKustomize); err != nil {
		return err
	}

	fmt.Printf("INFO: deploying Kong gateway RBACs to cluster")
	if err := clusters.KustomizeDeployForCluster(ctx, cluster, kongGatewayRBACsKustomize); err != nil {
		return err
	}

	fmt.Printf("INFO: deploying Kong CRDs RBACs to cluster")
	if err := clusters.KustomizeDeployForCluster(ctx, cluster, kongCRDsRBACsKustomize); err != nil {
		return err
	}

	return nil
}
