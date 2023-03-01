//go:build e2e_tests
// +build e2e_tests

package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestKongRouterCompatibility verifies that KIC behaves consistently with Kong routers
// `traditional` and `traditional_compatible`.
func TestKongRouterFlavorCompatibility(t *testing.T) {
	t.Parallel()
	ctx, env := setupE2ETest(t)
	cluster := env.Cluster()

	t.Log("deploying kong components with traditional Kong router")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	_ = deployKong(ctx, t, env, manifest)
	ensureGatewayDeployedWithRouterFlavor(ctx, t, env, "traditional")

	t.Log("running ingress tests to verify that KIC with traditonal Kong router works")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)

	setGatewayRouterFlavor(ctx, t, cluster, "traditional_compatible")

	t.Log("waiting for Kong with traditional_compatible router to start")
	ensureGatewayDeployedWithRouterFlavor(ctx, t, env, "traditional_compatible")

	t.Log("running ingress tests to verify that KIC with traditonal_compatible Kong router works")
	verifyIngress(ctx, t, env)
}

func setGatewayRouterFlavor(ctx context.Context, t *testing.T, cluster clusters.Cluster, flavor string) {
	// Since we cannot replace env vars in kustomize, here we update the deployment to set KONG_ROUTER_FLAVOR to traditional_compatible.
	t.Log("update deployment to modify Kong's router to traditional_compatible")
	gatewayDeployment, err := cluster.Client().AppsV1().Deployments(namespace).Get(ctx, "proxy-kong", metav1.GetOptions{})
	require.NoError(t, err)
	container := getContainerInPodSpec(&gatewayDeployment.Spec.Template.Spec, "proxy")
	require.NotNil(t, container)
	for i, env := range container.Env {
		if env.Name == "KONG_ROUTER_FLAVOR" {
			container.Env[i].Value = flavor
		}
	}
	_, err = cluster.Client().AppsV1().Deployments(namespace).Update(ctx, gatewayDeployment, metav1.UpdateOptions{})
	require.NoError(t, err)
}

func ensureGatewayDeployedWithRouterFlavor(ctx context.Context, t *testing.T, env environments.Environment, expectedFlavor string) {
	labelsForDeployment := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", "proxy-kong"),
	}
	require.Eventually(t, func() bool {
		podList, err := env.Cluster().Client().CoreV1().Pods(namespace).List(ctx, labelsForDeployment)
		require.NoError(t, err)
		if len(podList.Items) < 1 {
			return false
		}

		allPodsMatch := true
		for _, pod := range podList.Items {
			proxyContainer := getContainerInPodSpec(&pod.Spec, "proxy")
			if proxyContainer == nil {
				t.Logf("proxy container not found for Pod %s", pod.Name)
				allPodsMatch = false
				continue
			}
			if getEnvValueInContainer(proxyContainer, "KONG_ROUTER_FLAVOR") != expectedFlavor {
				t.Logf("KONG_ROUTER_FLAVOR is not set to expected value for Pod %s", pod.Name)
				allPodsMatch = false
			}
		}

		return allPodsMatch
	}, kongComponentWait, time.Second)
}
