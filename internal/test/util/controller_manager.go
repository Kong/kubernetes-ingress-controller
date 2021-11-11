package util

import (
	"context"
	"fmt"
	"os"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/cmd/rootcmd"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

// -----------------------------------------------------------------------------
// Testing Utility Functions - Controller Manager
// -----------------------------------------------------------------------------

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

	// deploy all CRDs required for testing to the cluster
	if err := DeployCRDsForCluster(ctx, cluster); err != nil {
		return fmt.Errorf("failed to deploy CRDs: %w", err)
	}

	// create a tempfile to hold the cluster kubeconfig that will be used for the controller
	// generate a temporary kubeconfig since we're going to be using the helm CLI
	kubeconfig, err := clusters.TempKubeconfig(cluster)
	if err != nil {
		return err
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
