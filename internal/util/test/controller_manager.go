package test

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/cmd/rootcmd"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
)

// logOutput is a file to use for manager log output other than stderr.
var logOutput = os.Getenv("TEST_KONG_KIC_MANAGER_LOG_OUTPUT")

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
func DeployControllerManagerForCluster(
	ctx context.Context,
	deprecatedLogger logrus.FieldLogger,
	logger logr.Logger,
	cluster clusters.Cluster,
	additionalFlags ...string,
) error {
	// ensure that the provided test cluster has a kongAddon deployed to it
	var kongAddon *ktfkong.Addon
	for _, addon := range cluster.ListAddons() {
		if addon.Name() == ktfkong.AddonName {
			var ok bool
			kongAddon, ok = addon.(*ktfkong.Addon)
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
		fmt.Sprintf("--publish-service=%s/ingress-controller-kong-proxy", consts.ControllerNamespace),
		fmt.Sprintf("--publish-service-udp=%s/ingress-controller-kong-udp-proxy", consts.ControllerNamespace),
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
		if err := rootcmd.RunWithLogger(ctx, &config, deprecatedLogger, logger); err != nil {
			panic(err)
		}
	}()

	return nil
}

// SetupLoggers sets up the loggers for the controller manager.
// The controller manager logs needs to be setup before the 30s timeout passes.
// Given the cluster deployment takes time, there is a chance that the timeout
// will pass before the controller manager logs are setup.
// This function can be used to sets up the loggers for the controller manager
// before the cluster deployment.
func SetupLoggers(logLevel string, logFormat string, logReduceRedundancy bool) (logrus.FieldLogger, logr.Logger, string, error) {
	// construct the config for the logger
	var err error
	output := os.Stderr
	if logOutput != "" {
		output, err = os.OpenFile(logOutput, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
		if err != nil {
			return nil, logr.Logger{}, logOutput, err
		}
	}
	config := manager.Config{
		LogLevel:            logLevel,
		LogFormat:           logFormat,
		LogReduceRedundancy: logReduceRedundancy,
	}

	deprecated, logger, err := manager.SetupLoggers(&config, output)
	return deprecated, logger, "", err
}
