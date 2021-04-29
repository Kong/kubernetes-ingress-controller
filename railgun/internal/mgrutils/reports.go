package mgrutils

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
)

// RunReport is runs the full report wanted for new KIC controller setups.
func RunReport(ctx context.Context, kubeCFG *rest.Config, kongCFG sendconfig.Kong, kicVersion string) error {
	// if anonymous reports are enabled this helps provide Kong with insights about usage of the ingress controller
	// which is non-sensitive and predominantly informs us of the controller and cluster versions in use.
	// This data helps inform us what versions, features, e.t.c. end-users are actively using which helps to inform
	// our prioritization of work and we appreciate when our end-users provide them, however if you do feel
	// uncomfortable and would rather turn them off run the controller with the "--anonymous-reports false" flag.

	// record the system hostname
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to fetch hostname: %w", err)
	}

	// create a universal unique identifer for this system
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("failed to generate a random uuid: %w", err)
	}

	// record the current Kubernetes server version
	kc, err := kubernetes.NewForConfig(kubeCFG)
	if err != nil {
		return fmt.Errorf("could not create client-go for Kubernetes discovery: %w", err)
	}
	k8sVersion, err := kc.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to fetch k8s api-server version: %w", err)
	}

	// gather versioning information from the kong client
	root, err := kongCFG.Client.Root(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Kong root config data: %w", err)
	}
	kongVersion, ok := root["version"].(string)
	if !ok {
		return fmt.Errorf("malformed Kong version found in Kong client root")
	}
	cfg, ok := root["configuration"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("malformed Kong configuration found in Kong client root")
	}
	kongDB, ok := cfg["database"].(string)
	if !ok {
		return fmt.Errorf("malformed database configuration found in Kong client root")
	}

	// build the final report
	info := util.Info{
		KongVersion:       kongVersion,
		KICVersion:        kicVersion,
		KubernetesVersion: k8sVersion.String(),
		Hostname:          hostname,
		ID:                uuid,
		KongDB:            kongDB,
	}

	// run the reporter in the background
	reporter := util.Reporter{Info: info}
	go reporter.Run(ctx.Done())

	return nil
}
