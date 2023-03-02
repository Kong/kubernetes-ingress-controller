//go:build e2e_tests
// +build e2e_tests

package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kuma"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/require"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func TestDeployAllInOneDBLESSKuma(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, env := setupE2ETest(t, kuma.New())

	t.Log("deploying kong components")
	manifest := getTestManifest(t, dblessPath)
	_ = deployKong(ctx, t, env, manifest)

	t.Log("adding Kuma mesh")
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "kong"))
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "default"))

	// scale to force a restart of pods and trigger mesh injection (we can't annotate the Kong namespace in advance,
	// it gets clobbered by deployKong()). is there a "rollout restart" in client-go? who knows!
	scaleDeployment(ctx, t, env, "proxy-kong", 0)
	scaleDeployment(ctx, t, env, "ingress-kong", 0)

	scaleDeployment(ctx, t, env, "proxy-kong", 2)
	scaleDeployment(ctx, t, env, "ingress-kong", 2)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)

	// use retry.RetryOnConflict to update service, to avoid conflicts from different source.
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
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

func scaleDeployment(ctx context.Context, t *testing.T, env environments.Environment, deploymentName string, replicas int32) {
	t.Helper()

	scale := &autoscalingv1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      deploymentName,
		},
		Spec: autoscalingv1.ScaleSpec{
			Replicas: replicas,
		},
	}
	_, err := env.Cluster().Client().AppsV1().Deployments(namespace).UpdateScale(ctx, deploymentName, scale, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		deployment, err := env.Cluster().Client().AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return deployment.Status.ReadyReplicas == replicas
	}, time.Minute*3, time.Second, "deployment %s did not scale to %d replicas", deploymentName, replicas)
}

func TestDeployAllInOnePostgresKuma(t *testing.T) {
	t.Log("configuring all-in-one-postgres.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	builder, err := getEnvironmentBuilder(ctx, t)
	require.NoError(t, err)
	builder = builder.WithAddons(kuma.New())
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	logClusterInfo(t, env.Cluster())

	defer func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	}()

	t.Log("deploying kong components")
	manifest := getTestManifest(t, postgresPath)
	_ = deployKong(ctx, t, env, manifest)

	t.Log("this deployment used a postgres backend, verifying that postgres migrations ran properly")
	verifyPostgres(ctx, t, env)

	t.Log("adding Kuma mesh")
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "kong"))
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "default"))

	// scale to force a restart of pods and trigger mesh injection (we can't annotate the Kong namespace in advance,
	// it gets clobbered by deployKong()). is there a "rollout restart" in client-go? who knows!
	scaleDeployment(ctx, t, env, "ingress-kong", 0)
	scaleDeployment(ctx, t, env, "ingress-kong", 2)

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
