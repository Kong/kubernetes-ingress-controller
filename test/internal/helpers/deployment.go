package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kongComponentRolloutTimeout = 15 * time.Minute // Increased from 7 to 15 minutes to allow more time for Kuma pods to start
)

// WaitForDeploymentRollout waits for the deployment to roll out in the cluster. It fails the test if the deployment
// doesn't roll out in time.
func WaitForDeploymentRollout(ctx context.Context, t *testing.T, cluster clusters.Cluster, namespace, name string) {
	require.Eventuallyf(t, func() bool {
		deployment, err := cluster.Client().AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false
		}

		if err := allUpdatedReplicasRolledOutAndReady(deployment); err != nil {
			t.Logf("%s/%s deployment not ready: %s", namespace, name, err)
			return false
		}

		return true
	}, kongComponentRolloutTimeout, time.Second, "deployment %s/%s didn't roll out in time", namespace, name)
}

// allUpdatedReplicasRolledOutAndReady ensures that all updated replicas are rolled out and ready. It is to make sure
// that the deployment rollout is finished and all the new replicas are ready to serve traffic before we proceed with
// the test. It returns an error with a reason if the deployment is not ready yet.
func allUpdatedReplicasRolledOutAndReady(d *appsv1.Deployment) error {
	if newReplicasRolledOut := d.Spec.Replicas != nil && d.Status.UpdatedReplicas < *d.Spec.Replicas; newReplicasRolledOut {
		return fmt.Errorf(
			"%d out of %d new replicas have been updated",
			d.Status.UpdatedReplicas,
			*d.Spec.Replicas,
		)
	}

	if oldReplicasPendingTermination := d.Status.Replicas > d.Status.UpdatedReplicas; oldReplicasPendingTermination {
		return fmt.Errorf(
			"%d old replicas pending termination",
			d.Status.Replicas-d.Status.UpdatedReplicas,
		)
	}

	if rolloutFinished := d.Status.AvailableReplicas == d.Status.UpdatedReplicas; !rolloutFinished {
		return fmt.Errorf(
			"%d of %d updated replicas are available",
			d.Status.AvailableReplicas,
			d.Status.UpdatedReplicas,
		)
	}

	return nil
}
