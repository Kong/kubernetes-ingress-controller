package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/cmd/rootcmd"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
)

// -----------------------------------------------------------------------------
// Testing Utility Functions - Controller Manager
// -----------------------------------------------------------------------------

// crds provides the filesystem paths to all the base CRDs needed for controller
// manager functionality.
var crds = []string{
	"../../config/crd/bases/configuration.konghq.com_udpingresses.yaml",
	"../../config/crd/bases/configuration.konghq.com_tcpingresses.yaml",
	"../../config/crd/bases/configuration.konghq.com_kongplugins.yaml",
	"../../config/crd/bases/configuration.konghq.com_kongingresses.yaml",
	"../../config/crd/bases/configuration.konghq.com_kongconsumers.yaml",
	"../../config/crd/bases/configuration.konghq.com_kongclusterplugins.yaml",
}

// DeployControllerManagerForCluster deploys all the base CRDs needed for the
// controller manager to function, and also runs a copy of the controller
// manager on a provided test cluster.
//
// Controller managers started this way will run in the background in a goroutine:
// The caller must use the provided context.Context to stop the controller manager
// from running when they're done with it.
func DeployControllerManagerForCluster(ctx context.Context, cluster clusters.Cluster, additionalFlags ...string) error {
	// ensure that the provided test cluster has a kongAddon deployed to it
	var kongAddon *kong.Addon
	for _, addon := range cluster.ListAddons() {
		if addon.Name() == kong.AddonName {
			var ok bool
			kongAddon, ok = addon.(*kong.Addon)
			if !ok {
				return fmt.Errorf("an invalid kong addon was present in test environment")
			}
		}
	}
	if kongAddon == nil {
		return fmt.Errorf("no kong addon found loaded in cluster %s", cluster.Name())
	}

	// determine the proxy admin URL for the Kong Gateway on this cluster:
	proxyAdminURL, err := kongAddon.ProxyAdminURL(ctx, cluster)
	if err != nil {
		return fmt.Errorf("couldn't determine Kong Gateway Admin URL on cluster %s: %w", cluster.Name(), err)
	}

	// create a tempfile to hold the cluster kubeconfig that will be used for the controller
	// generate a temporary kubeconfig since we're going to be using the helm CLI
	kubeconfig, err := generators.TempKubeconfig(cluster)
	if err != nil {
		return err
	}

	// deploy our CRDs to the cluster
	for _, crd := range crds {
		cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig.Name(), "apply", "-f", crd) //nolint:gosec
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			os.Remove(kubeconfig.Name())
			return fmt.Errorf("%s: %w", stderr.String(), err)
		}
	}

	// render all controller manager flag options
	controllerManagerFlags := []string{
		fmt.Sprintf("--kong-admin-url=http://%s:8001", proxyAdminURL.Hostname()),
		fmt.Sprintf("--kubeconfig=%s", kubeconfig.Name()),
		"--election-id=integrationtests.konghq.com",
		"--publish-service=kong-system/ingress-controller-kong-proxy",
		"--log-format=text",
	}
	controllerManagerFlags = append(controllerManagerFlags, additionalFlags...)

	// parsing the controller configuration flags
	config := manager.Config{}
	flags := config.FlagSet()
	if err := flags.Parse(controllerManagerFlags); err != nil {
		os.Remove(kubeconfig.Name())
		return fmt.Errorf("failed to parse manager flags: %w", err)
	}

	// run the controller in the background
	go func() {
		defer os.Remove(kubeconfig.Name())
		fmt.Fprintf(os.Stderr, "INFO: Starting Controller Manager for Cluster %s with Configuration: %+v\n", cluster.Name(), config)
		if err := rootcmd.Run(ctx, &config); err != nil {
			panic(err)
		}
	}()

	return nil
}
