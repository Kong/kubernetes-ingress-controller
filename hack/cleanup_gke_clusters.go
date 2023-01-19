// This script cleans up all GKE clusters that could be potentially
// orphaned by our tests (e.g. unexpected crash that didn't allow
// a test's teardown to be completed correctly). It's meant to be
// installed as a cronjob and run repeatedly throughout the day to
// catch any orphaned clusters: however tests should be trying to
// delete the clusters they create themselves.
//
// A cluster is considered orphaned when all conditions are satisfied:
// 1. Its name begins with a predefined prefix (`gke-e2e-`).
// 2. It was created more than 1h ago.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	container "cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/gke"
	"google.golang.org/api/option"

	"github.com/kong/kubernetes-ingress-controller/v2/test/e2e"
)

const timeUntilClusterOrphaned = time.Hour

var (
	gkeCreds    = os.Getenv(gke.GKECredsVar)
	gkeProject  = os.Getenv(gke.GKEProjectVar)
	gkeLocation = os.Getenv(gke.GKELocationVar)
)

func main() {
	mustNotBeEmpty(gke.GKECredsVar, gkeCreds)
	mustNotBeEmpty(gke.GKEProjectVar, gkeProject)
	mustNotBeEmpty(gke.GKELocationVar, gkeLocation)

	var creds map[string]string
	if err := json.Unmarshal([]byte(gkeCreds), &creds); err != nil {
		fmt.Fprintf(os.Stderr, "invalid credentials: %s\n", err)
		os.Exit(10)
	}

	if len(os.Args) < 1 {
		fmt.Fprintln(os.Stdout, "Usage: cleanup all | <list of cluster names...>")
		os.Exit(1)
	}

	credsOpt := option.WithCredentialsJSON([]byte(gkeCreds))
	mgrc, err := container.NewClusterManagerClient(context.Background(), credsOpt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create cluster manager client: %s", err)
		os.Exit(2)
	}
	defer mgrc.Close()

	clusterNames, err := findOrphanedClusters(context.Background(), mgrc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not find orphaned clusters: %s", err)
		os.Exit(2)
	}

	if len(clusterNames) < 1 {
		fmt.Println("INFO: no clusters to clean up")
		os.Exit(0)
	}

	var errs []error
	for _, clusterName := range clusterNames {
		fmt.Printf("INFO: cleaning up cluster %s\n", clusterName)
		err := deleteCluster(context.Background(), mgrc, gkeProject, gkeLocation, clusterName)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "failed to cleanup all clusters: %v\n", errs)
		os.Exit(3)
	}
}

func mustNotBeEmpty(name, value string) {
	if value == "" {
		panic(fmt.Sprintf("%s was empty", name))
	}
}

func deleteCluster(ctx context.Context, mgrc *container.ClusterManagerClient, project, location, name string) error {
	fullname := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, location, name)
	op, err := mgrc.DeleteCluster(ctx, &containerpb.DeleteClusterRequest{Name: fullname})
	if err != nil {
		return fmt.Errorf("failed to call delete cluster for %q: %w", name, err)
	}
	if op.Error != nil {
		return fmt.Errorf("failed to remove cluster %q: %s", name, op.Error)
	}

	return nil
}

func findOrphanedClusters(ctx context.Context, mgrc *container.ClusterManagerClient) ([]string, error) {
	clusterListReq := containerpb.ListClustersRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", gkeProject, gkeLocation),
	}
	clusterListResp, err := mgrc.ListClusters(ctx, &clusterListReq)
	if err != nil {
		return nil, err
	}

	var orphanedClusterNames []string
	for _, cluster := range clusterListResp.Clusters {
		if e2e.IsGKETestCluster(cluster) {
			createdAt, err := time.Parse(time.RFC3339, cluster.CreateTime)
			if err != nil {
				return nil, err
			}

			orphanTime := createdAt.Add(timeUntilClusterOrphaned)
			if time.Now().UTC().After(orphanTime) {
				orphanedClusterNames = append(orphanedClusterNames, cluster.Name)
			} else {
				fmt.Printf("INFO: cluster %s skipped (built in the last %s)\n", cluster.Name, timeUntilClusterOrphaned)
			}
		}
	}

	return orphanedClusterNames, nil
}
