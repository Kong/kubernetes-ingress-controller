//go:build e2e_tests
// +build e2e_tests

package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/loadimage"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestKongRouterCompatibility verifies that KIC behaves consistently with Kong routers
// `traditional` and `traditional_compatible`.
func TestKongRouterFlavorCompatibility(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	clusterBuilder := kind.NewBuilder()
	if clusterVersionStr != "" {
		clusterVersion, err := semver.ParseTolerant(clusterVersionStr)
		require.NoError(t, err)
		clusterBuilder.WithClusterVersion(clusterVersion)
	}
	cluster, err := clusterBuilder.Build(ctx)
	require.NoError(t, err)
	addons := []clusters.Addon{metallb.New()}

	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}

	builder := environments.NewBuilder().WithExistingCluster(cluster).WithAddons(addons...)
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	t.Logf("building cleaner to dump diagnostics...")
	cleaner := clusters.NewCleaner(env.Cluster())
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			if assert.NoError(t, err, "failed to dump diagnostics") {
				t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			}
		}
		assert.NoError(t, cluster.Cleanup(ctx))
	}()

	t.Log("deploying kong components with traditional Kong router")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)
	labelsForDeployment := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deployment.Name),
	}
	require.Eventually(t, func() bool {
		podList, err := env.Cluster().Client().CoreV1().Pods(deployment.Namespace).List(ctx, labelsForDeployment)
		require.NoError(t, err)
		if len(podList.Items) != 1 {
			return false
		}
		pod := podList.Items[0]
		proxyContainer := getContainerInPodSpec(&pod.Spec, "proxy")
		require.NotNil(t, proxyContainer)
		return getEnvValueInContainer(proxyContainer, "KONG_ROUTER_FLAVOR") == "traditional"
	}, kongComponentWait, time.Second)

	t.Log("running ingress tests to verify that KIC with traditonal Kong router works")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)

	// Since we cannot replace env vars in kustomize, here we update the deployment to set KONG_ROUTER_FLAVOR to traditional_compatible.
	t.Log("update deployment to modify Kong's router to traditional_compatible")
	deployment, err = cluster.Client().AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	require.NoError(t, err)
	container := getContainerInPodSpec(&deployment.Spec.Template.Spec, "proxy")
	require.NotNil(t, container)
	for i, env := range container.Env {
		if env.Name == "KONG_ROUTER_FLAVOR" {
			container.Env[i].Value = "traditional_compatible"
		}
	}
	_, err = cluster.Client().AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("waiting for Kong with traditional_compatible router to start")
	require.Eventually(t, func() bool {
		podList, err := env.Cluster().Client().CoreV1().Pods(deployment.Namespace).List(ctx, labelsForDeployment)
		require.NoError(t, err)
		if len(podList.Items) != 1 {
			return false
		}
		pod := podList.Items[0]
		proxyContainer := getContainerInPodSpec(&pod.Spec, "proxy")
		require.NotNil(t, proxyContainer)
		return getEnvValueInContainer(proxyContainer, "KONG_ROUTER_FLAVOR") == "traditional_compatible"
	}, 2*time.Minute, time.Second)
	t.Log("running ingress tests to verify that KIC with traditonal_compatible Kong router works")
	verifyIngress(ctx, t, env)
}
