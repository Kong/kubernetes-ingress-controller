package test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
)

// -----------------------------------------------------------------------------
// Testing Utility Functions - CRDs
// -----------------------------------------------------------------------------

const (
	kongCRDsKustomize = "../../config/crd/"
)

func DeployCRDsForCluster(ctx context.Context, cluster clusters.Cluster) error {
	// create a tempfile to hold the cluster kubeconfig that will be used for the controller
	// generate a temporary kubeconfig since we're going to be using the helm CLI
	kubeconfig, err := clusters.TempKubeconfig(cluster)
	if err != nil {
		return err
	}
	apiextClient, err := apiextclient.NewForConfig(cluster.Config())
	if err != nil {
		return err
	}
	defer os.Remove(kubeconfig.Name())

	// gather the YAML to deploy our CRDs
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	args := []string{"--kubeconfig", kubeconfig.Name(), "kustomize", kongCRDsKustomize}
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy kong CRDs STDOUT=(%s) STDERR=(%s): %w", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err)
	}
	fmt.Printf("INFO: running kubectl kustomize for Kong CRDs (args: %v)\n", args)
	kongCRDYAML := stdout.String()

	// gather the YAML to deploy Gateway CRDs
	stdout, stderr = new(bytes.Buffer), new(bytes.Buffer)
	args = []string{"--kubeconfig", kubeconfig.Name(), "kustomize", consts.GatewayExperimentalCRDsKustomizeURL}
	cmd = exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	fmt.Printf("INFO: running kubectl kustomize for Gateway CRDs (args: %v)\n", args)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy gateway CRDs STDOUT=(%s) STDERR=(%s): %w", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err)
	}
	gatewayCRDYAML := stdout.String()

	// deploy all CRDs required for testing
	for _, yaml := range []string{kongCRDYAML, gatewayCRDYAML} {
		if err := clusters.ApplyManifestByYAML(ctx, cluster, yaml); err != nil {
			return err
		}
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
