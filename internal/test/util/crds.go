package util

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
)

// -----------------------------------------------------------------------------
// Testing Utility Functions - CRDs
// -----------------------------------------------------------------------------

const (
	kongCRDsKustomize    = "../../config/crd/"
	gatewayCRDsKustomize = "https://github.com/kubernetes-sigs/gateway-api/config/crd"
)

func DeployCRDsForCluster(ctx context.Context, cluster clusters.Cluster) error {
	// create a tempfile to hold the cluster kubeconfig that will be used for the controller
	// generate a temporary kubeconfig since we're going to be using the helm CLI
	kubeconfig, err := clusters.TempKubeconfig(cluster)
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
		return fmt.Errorf("failed to deploy kong CRDs STDOUT=(%s) STDERR=(%s): %w", stdout.String(), stderr.String(), err)
	}
	kongCRDYAML := stdout.String()

	// gather the YAML to deploy Gateway CRDs
	stdout, stderr = new(bytes.Buffer), new(bytes.Buffer)
	args = []string{"--kubeconfig", kubeconfig.Name(), "kustomize", gatewayCRDsKustomize}
	cmd = exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy kong CRDs STDOUT=(%s) STDERR=(%s): %w", stdout.String(), stderr.String(), err)
	}
	gatewayCRDYAML := stdout.String()

	// deploy all CRDs required for testing
	for _, yaml := range []string{kongCRDYAML, gatewayCRDYAML} {
		if err := clusters.ApplyYAML(ctx, cluster, yaml); err != nil {
			return err
		}
	}

	return nil
}
