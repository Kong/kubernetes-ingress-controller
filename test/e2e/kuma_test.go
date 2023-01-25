//go:build e2e_tests
// +build e2e_tests

package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kuma"
	"github.com/stretchr/testify/require"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func TestDeployAllInOneDBLESSKuma(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	builder, err := getEnvironmentBuilder(ctx)
	require.NoError(t, err)
	builder = builder.WithAddons(kuma.New())
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	defer func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)

	t.Log("adding Kuma mesh")
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "kong"))
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "default"))

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

	// use retry.RetryOnConflict to update service, to avoid conflicts from different source.
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		service, err := env.Cluster().Client().CoreV1().Services("default").Get(ctx, "httpbin", metav1.GetOptions{})
		if err != nil {
			return err
		}

		if service.ObjectMeta.Annotations == nil {
			service.ObjectMeta.Annotations = map[string]string{}
		}
		service.ObjectMeta.Annotations["ingress.kubernetes.io/service-upstream"] = "true"
		_, err = env.Cluster().Client().CoreV1().Services("default").Update(ctx, service, metav1.UpdateOptions{})
		return err
	})
	require.NoError(t, err,
		// dump the status of service if the error happens on updating service.
		func() string {
			service, err := env.Cluster().Client().CoreV1().Services("default").Get(ctx, "httpbin", metav1.GetOptions{})
			if err != nil {
				return fmt.Sprintf("failed to dump service, error %v", err)
			}
			return fmt.Sprintf("current status of service: %#v", service)
		}(),
	)
	verifyIngress(ctx, t, env)
}

func TestDeployAllInOnePostgresKuma(t *testing.T) {
	t.Log("configuring all-in-one-postgres.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	builder, err := getEnvironmentBuilder(ctx)
	require.NoError(t, err)
	builder = builder.WithAddons(kuma.New())
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	defer func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, postgresPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)

	t.Log("this deployment used a postgres backend, verifying that postgres migrations ran properly")
	verifyPostgres(ctx, t, env)

	t.Log("adding Kuma mesh")
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "kong"))
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "default"))

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
	// use retry.RetryOnConflict to update service, to avoid conflicts from different source.
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		service, err := env.Cluster().Client().CoreV1().Services("default").Get(ctx, "httpbin", metav1.GetOptions{})
		if err != nil {
			return err
		}

		if service.ObjectMeta.Annotations == nil {
			service.ObjectMeta.Annotations = map[string]string{}
		}
		service.ObjectMeta.Annotations["ingress.kubernetes.io/service-upstream"] = "true"
		_, err = env.Cluster().Client().CoreV1().Services("default").Update(ctx, service, metav1.UpdateOptions{})
		return err
	})
	require.NoError(t, err,
		// dump the status of service if the error happens on updating service.
		func() string {
			service, err := env.Cluster().Client().CoreV1().Services("default").Get(ctx, "httpbin", metav1.GetOptions{})
			if err != nil {
				return fmt.Sprintf("failed to dump service, error %v", err)
			}
			return fmt.Sprintf("current status of service: %#v", service)
		}(),
	)

	verifyIngress(ctx, t, env)
}
