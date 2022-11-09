package main

// TODO: this is temporary: it was created for speed but will be replaced
//       by upstream functionality in KTF.
//       See: https://github.com/Kong/kubernetes-testing-framework/issues/61

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/gke"
)

const (
	k8sNameVar    = "KUBERNETES_CLUSTER_NAME"
	k8sVersionVar = "KONG_CLUSTER_VERSION"
)

var (
	ctx     = context.Background()
	cluster clusters.Cluster

	gkeCreds    = os.Getenv(gke.GKECredsVar)
	gkeProject  = os.Getenv(gke.GKEProjectVar)
	gkeLocation = os.Getenv(gke.GKELocationVar)
	k8sName     = os.Getenv(k8sNameVar)
	k8sVersion  = semver.MustParse(strings.TrimPrefix(os.Getenv(k8sVersionVar), "v"))
)

func main() {
	mustNotBeEmpty(gke.GKECredsVar, gkeCreds)
	mustNotBeEmpty(gke.GKEProjectVar, gkeProject)
	mustNotBeEmpty(gke.GKELocationVar, gkeLocation)
	mustNotBeEmpty(k8sVersionVar, k8sVersion.String())
	if k8sName == "" {
		k8sName = "kic-" + uuid.NewString()
		fmt.Println("INFO: no cluster name provided, using generated name " + k8sName)
	}

	fmt.Println("INFO: validating cluster version requirements")
	major := k8sVersion.Major
	minor := k8sVersion.Minor

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
	builder.WithClusterMinorVersion(major, minor)

	fmt.Printf("INFO: building cluster %s (this can take some time)\n", builder.Name)
	cluster, err := builder.Build(ctx)
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
