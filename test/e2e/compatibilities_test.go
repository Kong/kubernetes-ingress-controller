//go:build e2e_tests

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
	k8stypes "k8s.io/apimachinery/pkg/types"
)

// TestKongRouterCompatibility verifies that KIC behaves consistently with all
// Kong routers:
// - `expressions`
// - `traditional`
// - and `traditional_compatible`.
func TestKongRouterFlavorCompatibility(t *testing.T) {
	t.Parallel()

	ctx, env := setupE2ETest(t)
	cluster := env.Cluster()

	routerFlavors := []string{"expressions", "traditional_compatible", "traditional"}
	for _, rf := range routerFlavors {
		t.Run(rf, func(t *testing.T) {
			deploy := ManifestDeploy{
				Path: dblessPath,
				Patches: []ManifestPatch{
					patchKongRouterFlavorFn(rf),
				},
			}
			deployments := deploy.Run(ctx, t, env)
			t.Cleanup(func() { deploy.Delete(ctx, t, env) })
			proxyDeploymentNN := deployments.ProxyNN

			setGatewayRouterFlavor(ctx, t, cluster, proxyDeploymentNN, rf)
			t.Logf("waiting for Kong with %s router to start", rf)
			ensureGatewayDeployedWithRouterFlavor(ctx, t, env, proxyDeploymentNN, rf)
			t.Logf("running ingress tests to verify that KIC with %s Kong router works", rf)
			deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
			verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
		})
	}
}

func setGatewayRouterFlavor(
	ctx context.Context,
	t *testing.T,
	cluster clusters.Cluster,
	proxyDeploymentNN k8stypes.NamespacedName,
	flavor string,
) {
	t.Helper()

	// Since we cannot replace env vars in kustomize, here we update the deployment to set KONG_ROUTER_FLAVOR to traditional_compatible.
	t.Log("update deployment to modify Kong's router to traditional_compatible")
	deployments := cluster.Client().AppsV1().Deployments(proxyDeploymentNN.Namespace)
	gatewayDeployment, err := deployments.Get(ctx, proxyDeploymentNN.Name, metav1.GetOptions{})
	require.NoError(t, err)
	container := getContainerInPodSpec(&gatewayDeployment.Spec.Template.Spec, proxyContainerName)
	require.NotNil(t, container)
	for i, env := range container.Env {
		if env.Name == "KONG_ROUTER_FLAVOR" {
			container.Env[i].Value = flavor
		}
	}
	_, err = deployments.Update(ctx, gatewayDeployment, metav1.UpdateOptions{})
	require.NoError(t, err)
}

func ensureGatewayDeployedWithRouterFlavor(
	ctx context.Context,
	t *testing.T,
	env environments.Environment,
	proxyDeploymentNN k8stypes.NamespacedName,
	expectedFlavor string,
) {
	t.Helper()

	labelsForDeployment := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", proxyDeploymentNN.Name),
	}
	require.Eventually(t, func() bool {
		podList, err := env.Cluster().Client().CoreV1().Pods(namespace).List(ctx, labelsForDeployment)
		require.NoError(t, err)
		if len(podList.Items) < 1 {
			return false
		}

		allPodsMatch := true
		for _, pod := range podList.Items {
			proxyContainer := getContainerInPodSpec(&pod.Spec, proxyContainerName)
			if proxyContainer == nil {
				t.Logf("proxy container not found for Pod %s", pod.Name)
				allPodsMatch = false
				continue
			}
			if v := getEnvValueInContainer(proxyContainer, "KONG_ROUTER_FLAVOR"); v != expectedFlavor {
				t.Logf("KONG_ROUTER_FLAVOR is not set to expected value for Pod %s, actual: %s, expected: %s",
					pod.Name, v, expectedFlavor,
				)
				allPodsMatch = false
			}
		}

		return allPodsMatch
	}, kongComponentWait, time.Second)
}
