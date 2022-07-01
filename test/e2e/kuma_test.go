//go:build e2e_tests
// +build e2e_tests

package e2e

import (
	"context"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kuma"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/loadimage"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDeployAllInOneDBLESSKuma(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	addons = append(addons, kuma.New())
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.ParseTolerant(clusterVersionStr)
		require.NoError(t, err)
		builder.WithKubernetesVersion(clusterVersion)
	}
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, env.Cleanup(ctx))
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)

	t.Log("adding Kuma mesh")
	kuma.EnableMeshForNamespace(ctx, env.Cluster(), "kong")
	kuma.EnableMeshForNamespace(ctx, env.Cluster(), "default")

	// scale to force a restart of pods and trigger mesh injection (we can't annotate the Kong namespace in advance,
	// it gets clobbered by deployKong()). is there a "rollout restart" in client-go? who knows!
	scale := &autoscalingv1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
		},
		Spec: autoscalingv1.ScaleSpec{
			Replicas: 0,
		},
	}
	_, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).UpdateScale(ctx,
		deployment.Name, scale, metav1.UpdateOptions{})
	require.NoError(t, err)
	scale.Spec.Replicas = 2
	_, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).UpdateScale(ctx,
		deployment.Name, scale, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	service, err := env.Cluster().Client().CoreV1().Services("default").Get(ctx, "httpbin", metav1.GetOptions{})
	service.ObjectMeta.Annotations["ingress.kubernetes.io/service-upstream"] = "true"
	service, err = env.Cluster().Client().CoreV1().Services("default").Update(ctx, service, metav1.UpdateOptions{})
	verifyIngress(ctx, t, env)
}

func TestDeployAllInOnePostgresKuma(t *testing.T) {
	t.Log("configuring all-in-one-postgres.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())
	addons = append(addons, kuma.New())
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.ParseTolerant(clusterVersionStr)
		require.NoError(t, err)
		builder.WithKubernetesVersion(clusterVersion)
	}
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, env.Cleanup(ctx))
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, postgresPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)

	t.Log("this deployment used a postgres backend, verifying that postgres migrations ran properly")
	verifyPostgres(ctx, t, env)

	t.Log("adding Kuma mesh")
	kuma.EnableMeshForNamespace(ctx, env.Cluster(), "kong")
	kuma.EnableMeshForNamespace(ctx, env.Cluster(), "default")

	// scale to force a restart of pods and trigger mesh injection (we can't annotate the Kong namespace in advance,
	// it gets clobbered by deployKong()). is there a "rollout restart" in client-go? who knows!
	scale := &autoscalingv1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
		},
		Spec: autoscalingv1.ScaleSpec{
			Replicas: 0,
		},
	}
	_, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).UpdateScale(ctx,
		deployment.Name, scale, metav1.UpdateOptions{})
	require.NoError(t, err)
	scale.Spec.Replicas = 2
	_, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).UpdateScale(ctx,
		deployment.Name, scale, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	service, err := env.Cluster().Client().CoreV1().Services("default").Get(ctx, "httpbin", metav1.GetOptions{})
	require.NoError(t, err)
	service.ObjectMeta.Annotations["ingress.kubernetes.io/service-upstream"] = "true"
	service, err = env.Cluster().Client().CoreV1().Services("default").Update(ctx, service, metav1.UpdateOptions{})
	require.NoError(t, err)
	verifyIngress(ctx, t, env)
}
