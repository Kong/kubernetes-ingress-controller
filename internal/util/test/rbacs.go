package test

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"

	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
)

var (
	kongRBACsKustomize        = "config/rbac/"
	kongGatewayRBACsKustomize = "config/rbac/gateway"
	kongCRDsRBACsKustomize    = "config/rbac/crds"
	rbacsOnce                 sync.Once
)

func DeployRBACsForCluster(ctx context.Context, cluster clusters.Cluster) error {
	var err error
	rbacsOnce.Do(func() {
		var dir string
		// We need the repo root directory to be able to run this  from anywhere in the repository.
		dir, err = getRepoRoot(ctx)
		if err != nil {
			return
		}

		kongRBACsKustomize = filepath.Join(dir, kongRBACsKustomize)
		kongGatewayRBACsKustomize = filepath.Join(dir, kongGatewayRBACsKustomize)
		kongCRDsRBACsKustomize = filepath.Join(dir, kongCRDsRBACsKustomize)
	})
	if err != nil {
		return err
	}

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
