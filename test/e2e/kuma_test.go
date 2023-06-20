//go:build e2e_tests

package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kuma"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func TestDeployAllInOneDBLESSKuma(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, env := setupE2ETest(t, kuma.New())

	t.Log("deploying kong components")
	manifest := getTestManifest(t, dblessPath)
	deployKong(ctx, t, env, manifest)
	deployments := getManifestDeployments(dblessPath)

	t.Log("adding Kuma mesh")
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "kong"))
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "default"))

	// scale to force a restart of pods and trigger mesh injection (we can't annotate the Kong namespace in advance,
	// it gets clobbered by deployKong()). is there a "rollout restart" in client-go? who knows!
	scaleDeployment(ctx, t, env, deployments.ProxyNN, 0)
	scaleDeployment(ctx, t, env, deployments.ControllerNN, 0)

	scaleDeployment(ctx, t, env, deployments.ProxyNN, 2)
	scaleDeployment(ctx, t, env, deployments.ControllerNN, 2)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyKuma(ctx, t, env)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
}

func TestDeployAllInOnePostgresKuma(t *testing.T) {
	t.Log("configuring all-in-one-postgres.yaml manifest test")
	t.Parallel()

	ctx, env := setupE2ETest(t, kuma.New())

	t.Log("deploying kong components")
	manifest := getTestManifest(t, postgresPath)
	deployKong(ctx, t, env, manifest)
	deployments := getManifestDeployments(postgresPath)

	t.Log("this deployment used a postgres backend, verifying that postgres migrations ran properly")
	verifyPostgres(ctx, t, env)

	t.Log("adding Kuma mesh")
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "kong"))
	require.NoError(t, kuma.EnableMeshForNamespace(ctx, env.Cluster(), "default"))

	// scale to force a restart of pods and trigger mesh injection (we can't annotate the Kong namespace in advance,
	// it gets clobbered by deployKong()). is there a "rollout restart" in client-go? who knows!
	scaleDeployment(ctx, t, env, deployments.ControllerNN, 0)
	scaleDeployment(ctx, t, env, deployments.ControllerNN, 2)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyKuma(ctx, t, env)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
}

func verifyKuma(ctx context.Context, t *testing.T, env environments.Environment) {
	svcClient := env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault)
	const svcName = "echo"
	// Use retry.RetryOnConflict to update service, to avoid conflicts from different source.
	err := retry.RetryOnConflict(retry.DefaultRetry,
		func() error {
			service, err := svcClient.Get(ctx, svcName, metav1.GetOptions{})
			if err != nil {
				return err
			}

			if service.ObjectMeta.Annotations == nil {
				service.ObjectMeta.Annotations = map[string]string{}
			}
			service.ObjectMeta.Annotations["ingress.kubernetes.io/service-upstream"] = "true"
			_, err = svcClient.Update(ctx, service, metav1.UpdateOptions{})
			return err
		},
	)
	require.NoError(t, err,
		// dump the status of service if the error happens on updating service.
		func() string {
			service, err := svcClient.Get(ctx, svcName, metav1.GetOptions{})
			if err != nil {
				return fmt.Sprintf("failed to dump service, error %v", err)
			}
			return fmt.Sprintf("current status of service: %#v", service)
		}(),
	)
}
