package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	container "cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"google.golang.org/api/option"

	"github.com/kong/kubernetes-ingress-controller/v2/test/e2e"
)

const timeUntilClusterOrphaned = time.Hour

func cleanupGKEClusters(ctx context.Context) error {
	var creds map[string]string
	if err := json.Unmarshal([]byte(gkeCreds), &creds); err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}

	credsOpt := option.WithCredentialsJSON([]byte(gkeCreds))
	mgrc, err := container.NewClusterManagerClient(ctx, credsOpt)
	if err != nil {
		return fmt.Errorf("failed to create cluster manager client: %w", err)
	}
	defer mgrc.Close()

	clusterNames, err := findOrphanedClusters(ctx, mgrc)
	if err != nil {
		return fmt.Errorf("could not find orphaned clusters: %w", err)
	}

	if len(clusterNames) < 1 {
		log.Info("no clusters to clean up")
		return nil
	}

	var errs []error
	for _, clusterName := range clusterNames {
		log.Infof("cleaning up cluster %s\n", clusterName)
		err := deleteCluster(ctx, mgrc, gkeProject, gkeLocation, clusterName)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to cleanup all clusters: %w", errors.Join(errs...))
	}

	return nil
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
				log.Infof("cluster %s skipped (built in the last %s)\n", cluster.Name, timeUntilClusterOrphaned)
			}
		}
	}

	return orphanedClusterNames, nil
}
