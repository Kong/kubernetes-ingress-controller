package mgrutils

import (
	"context"
	"os"

	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
)

// RunReport is runs the full report wanted for new KIC controller setups.
func RunReport(ctx context.Context, kubeCFG *rest.Config, kongCFG sendconfig.Kong, kicVersion string) {
	// if anonymous reports are enabled this helps provide Kong with insights about usage of the ingress controller
	// which is non-sensitive and predominantly informs us of the controller and cluster versions in use.
	// This data helps inform us what versions, features, e.t.c. end-users are actively using which helps to inform
	// our prioritization of work and we appreciate when our end-users provide them, however if you do feel
	// uncomfortable and would rather turn them off run the controller with the "--anonymous-reports false" flag.
	reporterLogger := logrus.StandardLogger()

	// record the system hostname
	hostname, err := os.Hostname()
	if err != nil {
		reporterLogger.Error(err, "failed to fetch hostname")
	}

	// create a universal unique identifer for this system
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		reporterLogger.Error(err, "failed to generate a random uuid")
	}

	// record the current Kubernetes server version
	kc, err := kubernetes.NewForConfig(kubeCFG)
	if err != nil {
		reporterLogger.Error(err, "could not create client-go for Kubernetes discovery")
	}
	k8sVersion, err := kc.Discovery().ServerVersion()
	if err != nil {
		reporterLogger.Error(err, "failed to fetch k8s api-server version")
	}

	// gather versioning information from the kong client
	root, err := kongCFG.Client.Root(ctx)
	if err != nil {
		reporterLogger.Error(err, "failed to get Kong root config data")
	}
	kongVersion, ok := root["version"].(string)
	if !ok {
		reporterLogger.Error("malformed Kong version found in Kong client root")
	}
	cfg, ok := root["configuration"].(map[string]interface{})
	if !ok {
		reporterLogger.Error("malformed Kong configuration found in Kong client root")
	}
	kongDB, ok := cfg["database"].(string)
	if !ok {
		reporterLogger.Error("malformed database configuration found in Kong client root")
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
	reporter := util.Reporter{
		Info:   info,
		Logger: reporterLogger,
	}
	go reporter.Run(ctx.Done())
}
