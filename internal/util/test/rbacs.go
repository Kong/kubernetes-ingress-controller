package test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
)

const (
	kongRBACsKustomize        = "../../config/rbac/"
	kongKnativeRBACsKustomize = "../../config/rbac/knative"
	kongGatewayRBACsKustomize = "../../config/rbac/gateway"
	kongCRDsRBACsKustomize    = "../../config/rbac/crds"
)

func DeployRBACsForCluster(ctx context.Context, cluster clusters.Cluster) error {
	// create a tempfile to hold the cluster kubeconfig that will be used for the controller
	// generate a temporary kubeconfig since we're going to be using the helm CLI
	kubeconfig, err := clusters.TempKubeconfig(cluster)
	if err != nil {
		return err
	}
	defer os.Remove(kubeconfig.Name())

	// gather the YAML to deploy our RBACs
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	args := []string{"--kubeconfig", kubeconfig.Name(), "kustomize", kongRBACsKustomize}
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy kong RBACs STDOUT=(%s) STDERR=(%s): %w", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err)
	}
	fmt.Printf("INFO: running kubectl kustomize for Kong RBACs (args: %v)\n", args)
	kongRBACsYAML := stdout.String()

	// gather the YAML to deploy our RBACs
	stdout, stderr = new(bytes.Buffer), new(bytes.Buffer)
	args = []string{"--kubeconfig", kubeconfig.Name(), "kustomize", kongKnativeRBACsKustomize}
	cmd = exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy kong knative RBACs STDOUT=(%s) STDERR=(%s): %w", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err)
	}
	fmt.Printf("INFO: running kubectl kustomize for Kong knative RBACs (args: %v)\n", args)
	kongKnativeRBACsYAML := stdout.String()

	// gather the YAML to deploy our RBACs
	stdout, stderr = new(bytes.Buffer), new(bytes.Buffer)
	args = []string{"--kubeconfig", kubeconfig.Name(), "kustomize", kongGatewayRBACsKustomize}
	cmd = exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy kong gateway RBACs STDOUT=(%s) STDERR=(%s): %w", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err)
	}
	fmt.Printf("INFO: running kubectl kustomize for Kong gateway RBACs (args: %v)\n", args)
	kongGatewayRBACsYAML := stdout.String()

	// gather the YAML to deploy our RBACs
	stdout, stderr = new(bytes.Buffer), new(bytes.Buffer)
	args = []string{"--kubeconfig", kubeconfig.Name(), "kustomize", kongCRDsRBACsKustomize}
	cmd = exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy kong CRDs RBACs STDOUT=(%s) STDERR=(%s): %w", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err)
	}
	fmt.Printf("INFO: running kubectl kustomize for Kong CRDs RBACs (args: %v)\n", args)
	kongCRDsRBACsYAML := stdout.String()

	// deploy all RBACs required for testing
	for _, yaml := range []string{kongRBACsYAML, kongKnativeRBACsYAML, kongGatewayRBACsYAML, kongCRDsRBACsYAML} {
		if err := clusters.ApplyManifestByYAML(ctx, cluster, yaml); err != nil {
			return err
		}
	}

	return nil
}
