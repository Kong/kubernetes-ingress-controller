package main

// TODO: this is temporary: it was created for speed but will be replaced
//       by upstream functionality in KTF.
//       See: https://github.com/Kong/kubernetes-testing-framework/issues/61

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/gke"
)

const (
	k8sNameVar  = "KUBERNETES_CLUSTER_NAME"
	k8sMajorVar = "KUBERNETES_MAJOR_VERSION"
	k8sMinorVar = "KUBERNETES_MINOR_VERSION"
)

var (
	ctx     = context.Background()
	cluster clusters.Cluster

	gkeCreds    = os.Getenv(gke.GKECredsVar)
	gkeProject  = os.Getenv(gke.GKEProjectVar)
	gkeLocation = os.Getenv(gke.GKELocationVar)
	k8sName     = os.Getenv(k8sNameVar)
	k8sMajor    = os.Getenv(k8sMajorVar)
	k8sMinor    = os.Getenv(k8sMinorVar)
)

func main() {
	fmt.Println("INFO: configuring GKE cloud environment for tests")
	mustNotBeEmpty(gke.GKECredsVar, gkeCreds)
	mustNotBeEmpty(gke.GKEProjectVar, gkeProject)
	mustNotBeEmpty(gke.GKELocationVar, gkeLocation)
	mustNotBeEmpty(k8sNameVar, k8sName)
	mustNotBeEmpty(k8sMajorVar, k8sMajor)
	mustNotBeEmpty(k8sMinorVar, k8sMinor)

	fmt.Println("INFO: validating cluster version requirements")
	major, err := strconv.Atoi(k8sMajor)
	mustNotError(err)
	minor, err := strconv.Atoi(k8sMinor)
	mustNotError(err)

	if len(os.Args) > 1 && os.Args[1] == "cleanup" {
		fmt.Printf("INFO: cleanup called, deleting GKE cluster %s\n", k8sName)
		cluster, err := gke.NewFromExistingWithEnv(ctx, k8sName)
		mustNotError(err)
		mustNotError(cluster.Cleanup(ctx))
		fmt.Printf("INFO: GKE cluster %s successfully cleaned up\n", k8sName)
		os.Exit(0)
	}

	fmt.Printf("INFO: configuring the GKE cluster NAME=(%s) VERSION=(v%d.%d) PROJECT=(%s) LOCATION=(%s)\n", k8sName, major, minor, gkeProject, gkeLocation)
	builder := gke.NewBuilder([]byte(gkeCreds), gkeProject, gkeLocation).WithName(k8sName)
	builder.WithClusterMinorVersion(uint64(major), uint64(minor))

	fmt.Printf("INFO: building cluster %s (this can take some time)\n", builder.Name)
	cluster, err = builder.Build(ctx)
	mustNotError(err)

	fmt.Println("INFO: verifying that the cluster can be communicated with")
	version, err := cluster.Client().ServerVersion()
	mustNotError(err)

	fmt.Printf("INFO: server version found: %s\n", version)
}

func mustNotBeEmpty(name, value string) {
	if value == "" {
		if cluster != nil {
			if err := cluster.Cleanup(ctx); err != nil {
				panic(fmt.Sprintf("%s was empty, and then cleanup failed: %s", name, err))
			}
		}
		panic(fmt.Sprintf("%s was empty", name))
	}
}

func mustNotError(err error) {
	if err != nil {
		if cluster != nil {
			if cleanupErr := cluster.Cleanup(ctx); cleanupErr != nil {
				panic(fmt.Sprintf("deployment failed with %s, and then cleanup failed: %s", err, cleanupErr))
			}
		}
		panic(fmt.Errorf("failed to deploy e2e environment: %w", err))
	}
}
