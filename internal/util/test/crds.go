package test

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
)

// -----------------------------------------------------------------------------
// Testing Utility Functions - CRDs
// -----------------------------------------------------------------------------

var (
	kongCRDsKustomize = "config/crd/"
	crdsOnce          sync.Once
)

func DeployCRDsForCluster(ctx context.Context, cluster clusters.Cluster) error {
	var err error
	crdsOnce.Do(func() {
		var dir string
		// We need the repo root directory to be able to run this  from anywhere in the repository.
		dir, err = getRepoRoot(ctx)
		if err != nil {
			return
		}

		kongCRDsKustomize = filepath.Join(dir, kongCRDsKustomize)
	})
	if err != nil {
		return err
	}

	apiextClient, err := apiextclient.NewForConfig(cluster.Config())
	if err != nil {
		return err
	}

	fmt.Printf("INFO: deploying Kong CRDs to cluster\n")
	if err := clusters.KustomizeDeployForCluster(ctx, cluster, kongCRDsKustomize); err != nil {
		return err
	}

	fmt.Printf("INFO: deploying Gateway CRDs to cluster\n")
	if err := clusters.KustomizeDeployForCluster(ctx, cluster, consts.GatewayExperimentalCRDsKustomizeURL); err != nil {
		return err
	}

	for _, crd := range []string{
		"gatewayclasses.gateway.networking.k8s.io",
		"gateways.gateway.networking.k8s.io",
		"httproutes.gateway.networking.k8s.io",
		"referencegrants.gateway.networking.k8s.io",
		"tcproutes.gateway.networking.k8s.io",
		"tlsroutes.gateway.networking.k8s.io",
		"udproutes.gateway.networking.k8s.io",
	} {
		if err := retry.OnError(
			retry.DefaultRetry,
			apierrors.IsNotFound,
			func() error {
				_, err := apiextClient.CustomResourceDefinitions().Get(ctx, crd, metav1.GetOptions{})
				return err
			},
		); err != nil {
			return err
		}
	}

	return nil
}
