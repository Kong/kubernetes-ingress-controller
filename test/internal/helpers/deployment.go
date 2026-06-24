package helpers

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kongComponentRolloutTimeout = 7 * time.Minute
)

// WaitForDeploymentRollout waits for the deployment to roll out in the cluster. It fails the test if the deployment
// doesn't roll out in time.
func WaitForDeploymentRollout(ctx context.Context, t *testing.T, cluster clusters.Cluster, namespace, name string) {
	t.Helper()

	var iteration int
	require.Eventuallyf(t, func() bool {
		iteration++
		deployment, err := cluster.Client().AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false
		}

		if err := allUpdatedReplicasRolledOutAndReady(deployment); err != nil {
			t.Logf("%s/%s deployment not ready: %s", namespace, name, err)
			if iteration%30 == 1 {
				logPodsForDeployment(ctx, t, cluster, namespace, deployment)
			}
			return false
		}

		return true
	}, kongComponentRolloutTimeout, time.Second, "deployment %s/%s didn't roll out in time", namespace, name)
}

// logPodsForDeployment logs the status of pods belonging to the deployment.
func logPodsForDeployment(ctx context.Context, t *testing.T, cluster clusters.Cluster, namespace string, deployment *appsv1.Deployment) {
	t.Helper()

	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		t.Logf("could not build selector for %s/%s: %v", namespace, deployment.Name, err)
		return
	}

	pods, err := cluster.Client().CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		t.Logf("could not list pods for %s/%s: %v", namespace, deployment.Name, err)
		return
	}

	if len(pods.Items) == 0 {
		t.Logf("%s/%s: no pods found", namespace, deployment.Name)
		return
	}

	for _, pod := range pods.Items {
		t.Logf("%s/%s pod %s: phase=%s", namespace, deployment.Name, pod.Name, pod.Status.Phase)
		for _, cs := range pod.Status.ContainerStatuses {
			t.Logf("  container %s: ready=%v restarts=%d state=%s",
				cs.Name, cs.Ready, cs.RestartCount, containerStateString(cs.State))
		}
		for _, cs := range pod.Status.InitContainerStatuses {
			t.Logf("  init-container %s: ready=%v restarts=%d state=%s",
				cs.Name, cs.Ready, cs.RestartCount, containerStateString(cs.State))
		}
		for _, cond := range pod.Status.Conditions {
			if cond.Status != corev1.ConditionTrue {
				t.Logf("  condition %s=%s: %s", cond.Type, cond.Status, cond.Message)
			}
		}
	}
}

func containerStateString(s corev1.ContainerState) string {
	if s.Running != nil {
		return "Running"
	}
	if s.Waiting != nil {
		var b strings.Builder
		b.WriteString("Waiting")
		if s.Waiting.Reason != "" {
			fmt.Fprintf(&b, "(%s)", s.Waiting.Reason)
		}
		if s.Waiting.Message != "" {
			fmt.Fprintf(&b, ": %s", s.Waiting.Message)
		}
		return b.String()
	}
	if s.Terminated != nil {
		return fmt.Sprintf("Terminated(exit=%d reason=%s)", s.Terminated.ExitCode, s.Terminated.Reason)
	}
	return "Unknown"
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
