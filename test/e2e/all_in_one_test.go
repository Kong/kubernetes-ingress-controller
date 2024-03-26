//go:build e2e_tests

package e2e

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/require"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

// -----------------------------------------------------------------------------
// All-In-One Manifest Tests - Suite
//
// The following tests ensure that the local "all-in-one" style deployment manifests
// (which are predominantly used for testing, whereas the helm chart is meant for
// production use cases) are functional by deploying them to a cluster and verifying
// some of the fundamental functionality of the ingress controller and the proxy to
// ensure that things are up and running.
// -----------------------------------------------------------------------------

func TestDeployAllInOneEnterpriseDBLESS(t *testing.T) {
	const entDBLESSPath = "manifests/all-in-one-dbless-k4k8s-enterprise.yaml"

	t.Logf("configuring %s manifest test", entDBLESSPath)
	if os.Getenv(kong.LicenseDataEnvVar) == "" {
		t.Skipf("no license available to test enterprise: %s was not provided", kong.LicenseDataEnvVar)
	}
	t.Parallel()
	ctx, env := setupE2ETest(t)

	createKongImagePullSecret(ctx, t, env)

	t.Log("generating a superuser password")
	adminPassword, adminPasswordSecretYAML, err := generateAdminPasswordSecret()
	require.NoError(t, err)

	t.Log("generating a license secret")
	licenseSecret, err := kong.GetLicenseSecretFromEnv()
	require.NoError(t, err)

	t.Log("deploying kong components")
	deployments := ManifestDeploy{
		Path:              entDBLESSPath,
		AdditionalSecrets: []*corev1.Secret{licenseSecret, adminPasswordSecretYAML},
	}.Run(ctx, t, env)

	t.Log("exposing the admin api so that enterprise features can be verified")
	exposeAdminAPI(ctx, t, env, deployments.ProxyNN)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)

	t.Log("verifying enterprise mode was enabled properly")
	verifyEnterprise(ctx, t, env, adminPassword)
}

const postgresPath = "manifests/all-in-one-postgres.yaml"

func TestDeployAllInOnePostgres(t *testing.T) {
	t.Log("configuring all-in-one-postgres.yaml manifest test")
	t.Parallel()
	ctx, env := setupE2ETest(t)

	t.Log("deploying kong components")
	ManifestDeploy{Path: postgresPath}.Run(ctx, t, env)

	t.Log("this deployment used a postgres backend, verifying that postgres migrations ran properly")
	verifyPostgres(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
}

func TestDeployAllInOnePostgresWithMultipleReplicas(t *testing.T) {
	t.Log("configuring all-in-one-postgres.yaml manifest test")
	t.Parallel()
	ctx, env := setupE2ETest(t)

	t.Log("deploying kong components")
	deployments := ManifestDeploy{Path: postgresPath}.Run(ctx, t, env)
	deployment := deployments.ControllerNN

	t.Log("this deployment used a postgres backend, verifying that postgres migrations ran properly")
	verifyPostgres(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)

	t.Log("verifying that kong pods deployed properly and gathering a sample pod")
	forDeployment := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deployment.Name),
	}
	podList, err := env.Cluster().Client().CoreV1().Pods(deployment.Namespace).List(ctx, forDeployment)
	require.NoError(t, err)
	require.Equal(t, 1, len(podList.Items))
	initialPod := podList.Items[0]

	t.Log("adding a second replica to the Kong deployment")
	scale := &autoscalingv1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
		},
		Spec: autoscalingv1.ScaleSpec{
			Replicas: 2,
		},
	}
	_, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).UpdateScale(ctx,
		deployment.Name, scale, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("verifying that scaling completes and the additional replicas come up")
	require.Eventually(t, func() bool {
		deployment, err := env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return deployment.Status.ReadyReplicas == *deployment.Spec.Replicas
	}, kongComponentWait, time.Second)

	t.Log("gathering another sample pod to verify leadership is configured appropriately")
	podList, err = env.Cluster().Client().CoreV1().Pods(deployment.Namespace).List(ctx, forDeployment)
	require.NoError(t, err)
	var secondary corev1.Pod
	for _, pod := range podList.Items {
		if pod.Name != initialPod.Name {
			secondary = pod
			break
		}
	}

	client := &http.Client{Timeout: time.Second * 30}
	t.Log("confirming the second replica is not the leader and is not pushing configuration")
	forwardCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	localPort := startPortForwarder(forwardCtx, t, env, secondary.Namespace, secondary.Name, "cmetrics")

	require.Never(t, func() bool {
		// if we are not the leader, we run no config pushes, and this metric string will not appear.
		return httpGetResponseContains(t, fmt.Sprintf("http://localhost:%d/metrics", localPort), client, metrics.MetricNameConfigPushCount)
	}, time.Minute, time.Second*10)

	// since leader election is time sensitive, we log the time here.
	t.Logf("deleting the original replica and current leader at %v", time.Now())
	err = env.Cluster().Client().CoreV1().Pods(initialPod.Namespace).Delete(ctx, initialPod.Name, metav1.DeleteOptions{})
	require.NoError(t, err)

	t.Logf("waiting for the initial pod disappear and new pod to be recreated and up")
	require.Eventually(t, func() bool {
		podList, err = env.Cluster().Client().CoreV1().Pods(initialPod.Namespace).List(ctx, forDeployment)
		require.NoError(t, err)
		podNum := 0
		// we wait for the number of running pod excluding the initial one to be 2
		// since the replicas is set to 2 in the deployment.
		// So if there are exactly 2 running pods except the initial pod, we can know
		// that the new pod is recreated and up after the initial one is deleted,
		// and the status of deployment runs into a stable state.
		for _, pod := range podList.Items {
			if pod.Name != initialPod.Name && pod.Status.Phase == corev1.PodRunning {
				podNum++
			}
		}
		return podNum == 2
	}, time.Minute, time.Second)

	var (
		rebuiltPod       corev1.Pod
		rebuiltLocalPort int
	)
	for _, pod := range podList.Items {
		if pod.Name != initialPod.Name && pod.Name != secondary.Name {
			rebuiltPod = pod
			rebuiltLocalPort = startPortForwarder(forwardCtx, t, env, rebuiltPod.Namespace, rebuiltPod.Name, "cmetrics")
			break
		}
	}

	// Pass the test if exactly one of the pod becomes the leader, not limited to the original secondary pod.
	// Because in several times, the rebuilt pod (new pod created after initial pod deleted) became the leader.
	t.Logf("confirming there is exactly one pod that becomes leader and starts pushing configuration at %v", time.Now())
	require.Eventually(t, func() bool {
		leaderCount := 0
		if httpGetResponseContains(t, fmt.Sprintf("http://localhost:%d/metrics", localPort), client, metrics.MetricNameConfigPushCount) {
			t.Logf("secondary pod %s is the leader at %v", secondary.Name, time.Now())
			leaderCount++
		}
		if httpGetResponseContains(t, fmt.Sprintf("http://localhost:%d/metrics", rebuiltLocalPort), client, metrics.MetricNameConfigPushCount) {
			t.Logf("rebuilt pod %s is the leader at %v", rebuiltPod.Name, time.Now())
			leaderCount++
		}
		t.Logf("expected exactly one leader, actual %d", leaderCount)
		return leaderCount == 1
	}, 2*time.Minute, time.Second)
}

const entPostgresPath = "manifests/all-in-one-postgres-enterprise.yaml"

func TestDeployAllInOneEnterprisePostgres(t *testing.T) {
	t.Log("configuring all-in-one-postgres-enterprise.yaml manifest test")
	if os.Getenv(kong.LicenseDataEnvVar) == "" {
		t.Skipf("no license available to test enterprise: %s was not provided", kong.LicenseDataEnvVar)
	}
	t.Parallel()
	ctx, env := setupE2ETest(t)

	createKongImagePullSecret(ctx, t, env)

	t.Log("generating a superuser password")
	adminPassword, adminPasswordSecret, err := generateAdminPasswordSecret()
	require.NoError(t, err)

	t.Log("generating a license secret")
	licenseSecret, err := kong.GetLicenseSecretFromEnv()
	require.NoError(t, err)

	t.Log("deploying kong components")
	deployments := ManifestDeploy{
		Path:              entPostgresPath,
		AdditionalSecrets: []*corev1.Secret{licenseSecret, adminPasswordSecret},
	}.Run(ctx, t, env)

	t.Log("this deployment used a postgres backend, verifying that postgres migrations ran properly")
	verifyPostgres(ctx, t, env)

	t.Log("running ingress tests to verify ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)

	t.Log("exposing the admin api so that enterprise features can be verified")
	exposeAdminAPI(ctx, t, env, deployments.ProxyNN)

	t.Log("this deployment used enterprise kong, verifying that enterprise functionality was set up properly")
	verifyEnterprise(ctx, t, env, adminPassword)
	verifyEnterpriseWithPostgres(ctx, t, env, adminPassword)
}

func TestDeployAllInOnePostgresGatewayDiscovery(t *testing.T) {
	t.Parallel()

	const manifestFilePath = "manifests/all-in-one-postgres-multiple-gateways.yaml"

	t.Logf("configuring %s manifest test", manifestFilePath)
	ctx, env := setupE2ETest(t)

	t.Log("deploying kong components")
	deployments := ManifestDeploy{Path: manifestFilePath}.Run(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	ensureAllProxyReplicasAreConfigured(ctx, t, env, deployments.ProxyNN)
}

func TestDeployAllInOneDBLESS(t *testing.T) {
	t.Parallel()

	const manifestFilePath = "manifests/all-in-one-dbless.yaml"

	t.Logf("configuring %s manifest test", manifestFilePath)
	ctx, env := setupE2ETest(t)

	t.Log("deploying kong components")
	deployments := ManifestDeploy{Path: manifestFilePath}.Run(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	ingress := deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	ensureAllProxyReplicasAreConfigured(ctx, t, env, deployments.ProxyNN)

	t.Log("scale proxy to 0 replicas")
	scaleDeployment(ctx, t, env, deployments.ProxyNN, 0)

	t.Log("wait for 10 seconds to let controller reconcile")
	<-time.After(10 * time.Second)

	t.Log("ensure that controller pods didn't crash after scaling proxy to 0")
	expectedControllerReplicas := *(deployments.GetController(ctx, t, env).Spec.Replicas)
	readyControllerReplicas := deployments.GetController(ctx, t, env).Status.ReadyReplicas
	require.Equal(t, expectedControllerReplicas, readyControllerReplicas,
		"controller replicas count should not change after scaling proxy to 0")
	ensureNoneOfDeploymentPodsHasCrashed(ctx, t, env, deployments.ControllerNN)

	t.Log("scale proxy to 3 replicas and wait for all instances to be ready")
	scaleDeployment(ctx, t, env, deployments.ProxyNN, 3)
	ensureAllProxyReplicasAreConfigured(ctx, t, env, deployments.ProxyNN)

	t.Log("scale proxy to 1 replica")
	scaleDeployment(ctx, t, env, deployments.ProxyNN, 1)

	t.Log("misconfigure the ingress")
	reconfigureExistingIngress(ctx, t, env, ingress, func(_ *netv1.Ingress) {
		ingress.Spec.Rules[0].HTTP.Paths[0].Path = badEchoPath
	})

	t.Log("scale proxy to 2 replicas and verify that the new replica gets the old good configuration")
	scaleDeployment(ctx, t, env, deployments.ProxyNN, 2)
	// Verify all the proxy replicas have the last good configuration.
	ensureAllProxyReplicasAreConfigured(ctx, t, env, deployments.ProxyNN)

	t.Log("restart the controller")
	deployments.RestartController(ctx, t, env)
	helpers.WaitForDeploymentRollout(ctx, t, env.Cluster(), namespace, controllerDeploymentName)

	t.Log("scale proxy to 3 replicas and verify that the new replica gets the old good configuration")
	scaleDeployment(ctx, t, env, deployments.ProxyNN, 3)
	// Verify all the proxy replicas have the last good configuration.
	ensureAllProxyReplicasAreConfigured(ctx, t, env, deployments.ProxyNN)
}

func ensureAllProxyReplicasAreConfigured(ctx context.Context, t *testing.T, env environments.Environment, proxyDeploymentNN k8stypes.NamespacedName) {
	pods, err := listPodsByLabels(ctx, env, proxyDeploymentNN.Namespace, map[string]string{"app": proxyDeploymentNN.Name})
	require.NoError(t, err)

	t.Logf("ensuring all %d proxy replicas are configured", len(pods))
	wg := sync.WaitGroup{}
	for _, pod := range pods {
		pod := pod
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{
				Timeout: time.Second * 30,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}

			forwardCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			localPort := startPortForwarder(forwardCtx, t, env, proxyDeploymentNN.Namespace, pod.Name, "8444")
			address := fmt.Sprintf("https://localhost:%d", localPort)

			kongClient, err := adminapi.NewKongAPIClient(address, client)
			require.NoError(t, err)

			verifyIngressWithEchoBackendsInAdminAPI(ctx, t, kongClient, numberOfEchoBackends)
			t.Logf("proxy pod %s/%s: got the config", pod.Namespace, pod.Name)
		}()
	}

	wg.Wait()
}
