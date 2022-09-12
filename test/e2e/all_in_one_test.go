//go:build e2e_tests
// +build e2e_tests

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

// -----------------------------------------------------------------------------
// All-In-One Manifest Tests - Vars
// -----------------------------------------------------------------------------

const (
	// kongComponentWait is the maximum amount of time to wait for components (such as
	// the ingress controller or the Kong Gateway) to become responsive after
	// deployment to the cluster has finished.
	kongComponentWait = time.Minute * 7

	// ingressWait is the maximum amount of time to wait for a basic HTTP service
	// (e.g. httpbin) to come online and for ingress to have properly configured
	// proxy traffic to route to it.
	ingressWait = time.Minute * 5

	// adminAPIWait is the maximum amount of time to wait for the Admin API to become
	// responsive after updating the KONG_ADMIN_LISTEN and adding a service for it.
	adminAPIWait = time.Minute * 2

	// gatewayUpdateWaitTime is the amount of time to wait for updates to the Gateway, or to its
	// parent Service to fully resolve into ready state.
	gatewayUpdateWaitTime = time.Minute * 3
)

var (
	imageOverride         = os.Getenv("TEST_KONG_CONTROLLER_IMAGE_OVERRIDE")
	imageLoad             = os.Getenv("TEST_KONG_CONTROLLER_IMAGE_LOAD")
	kongImageOverride     = os.Getenv("TEST_KONG_IMAGE_OVERRIDE")
	kongImageLoad         = os.Getenv("TEST_KONG_IMAGE_LOAD")
	kongImagePullUsername = os.Getenv("TEST_KONG_PULL_USERNAME")
	kongImagePullPassword = os.Getenv("TEST_KONG_PULL_PASSWORD")
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

const (
	dblessPath = "../../deploy/single/all-in-one-dbless.yaml"
	dblessURL  = "https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/%v.%v.x/deploy/single/all-in-one-dbless.yaml"
)

func TestDeployAllInOneDBLESS(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())

	addons = append(addons, buildImageLoadAddons(t, imageLoad, kongImageLoad)...)

	builder := setBuilderKubernetesVersion(t,
		environments.NewBuilder().WithAddons(addons...), clusterVersionStr)
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, env.Cleanup(ctx))
	}()

	t.Logf("build a cleaner to dump diagnostics...")
	cleaner := clusters.NewCleaner(env.Cluster())
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)

	forDeployment := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deployment.Name),
	}
	podList, err := env.Cluster().Client().CoreV1().Pods(deployment.Namespace).List(ctx, forDeployment)
	require.NoError(t, err)
	require.Equal(t, 1, len(podList.Items))
	pod := podList.Items[0]

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)

	t.Log("killing Kong process to simulate a crash and container restart")
	killKong(ctx, t, env, &pod)

	t.Log("confirming that routes are restored after crash")
	verifyIngress(ctx, t, env)
}

func TestDeployAndUpgradeAllInOneDBLESS(t *testing.T) {
	curTag, err := getCurrentGitTag("")
	require.NoError(t, err)
	preTag, err := getPreviousGitTag("", curTag)
	require.NoError(t, err)
	if curTag.Patch != 0 || len(curTag.Pre) > 0 {
		t.Skipf("%v not a new minor version, skipping upgrade test", curTag)
	}
	oldManifest, err := http.Get(fmt.Sprintf(dblessURL, preTag.Major, preTag.Minor))
	require.NoError(t, err)
	defer oldManifest.Body.Close()

	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())

	addons = append(addons, buildImageLoadAddons(t, imageLoad, kongImageLoad)...)

	builder := setBuilderKubernetesVersion(t,
		environments.NewBuilder().WithAddons(addons...), clusterVersionStr)

	env, err := builder.Build(ctx)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, env.Cleanup(ctx))
	}()

	t.Logf("build a cleaner to dump diagnostics...")
	cleaner := clusters.NewCleaner(env.Cluster())
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Logf("deploying previous version %s kong manifest", preTag)
	deployKong(ctx, t, env, oldManifest.Body)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)

	t.Logf("deploying current version %s kong manifest", curTag)

	manifest, err := getTestManifest(t, dblessPath)
	require.NoError(t, err)
	deployKong(ctx, t, env, manifest)
	verifyIngress(ctx, t, env)
}

const entDBLESSPath = "../../deploy/single/all-in-one-dbless-k4k8s-enterprise.yaml"

func TestDeployAllInOneEnterpriseDBLESS(t *testing.T) {
	t.Log("configuring all-in-one-dbless-k4k8s-enterprise.yaml manifest test")
	if os.Getenv(kong.LicenseDataEnvVar) == "" {
		t.Skipf("no license available to test enterprise: %s was not provided", kong.LicenseDataEnvVar)
	}
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())

	addons = append(addons, buildImageLoadAddons(t, imageLoad, kongImageLoad)...)

	builder := setBuilderKubernetesVersion(t,
		environments.NewBuilder().WithAddons(addons...), clusterVersionStr)
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	createKongImagePullSecret(ctx, t, env)

	t.Logf("build a cleaner to dump diagnostics...")
	cleaner := clusters.NewCleaner(env.Cluster())
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Log("generating a superuser password")
	adminPassword, adminPasswordSecretYAML, err := generateAdminPasswordSecret()
	require.NoError(t, err)

	t.Log("generating a license secret")
	licenseSecret, err := kong.GetLicenseSecretFromEnv()
	require.NoError(t, err)

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, entDBLESSPath)
	require.NoError(t, err)
	_ = deployKong(ctx, t, env, manifest, licenseSecret, adminPasswordSecretYAML)

	t.Log("exposing the admin api so that enterprise features can be verified")
	exposeAdminAPI(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)

	t.Log("verifying enterprise mode was enabled properly")
	verifyEnterprise(ctx, t, env, adminPassword)
}

const postgresPath = "../../deploy/single/all-in-one-postgres.yaml"

func TestDeployAllInOnePostgres(t *testing.T) {
	t.Log("configuring all-in-one-postgres.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())

	addons = append(addons, buildImageLoadAddons(t, imageLoad, kongImageLoad)...)

	builder := setBuilderKubernetesVersion(t,
		environments.NewBuilder().WithAddons(addons...), clusterVersionStr)
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	t.Logf("build a cleaner to dump diagnostics...")
	cleaner := clusters.NewCleaner(env.Cluster())
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, postgresPath)
	require.NoError(t, err)
	_ = deployKong(ctx, t, env, manifest)

	t.Log("this deployment used a postgres backend, verifying that postgres migrations ran properly")
	verifyPostgres(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)
}

func TestDeployAllInOnePostgresWithMultipleReplicas(t *testing.T) {
	t.Log("configuring all-in-one-postgres.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())

	addons = append(addons, buildImageLoadAddons(t, imageLoad, kongImageLoad)...)

	builder := setBuilderKubernetesVersion(t,
		environments.NewBuilder().WithAddons(addons...), clusterVersionStr)
	env, err := builder.Build(ctx)
	require.NoError(t, err)

	t.Logf("build a cleaner to dump diagnostics...")
	cleaner := clusters.NewCleaner(env.Cluster())
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, postgresPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)
	// dump diagnostics and print out logs of KIC pod to a temporary directory, if the test failed.
	defer func() {
		if t.Failed() {
			outputDir, err := cleaner.DumpDiagnostics(ctx, t.Name())
			assert.NoError(t, err, "failed to dump diagnostics")
			t.Logf("%s failed, dumped diagnostics to directory %s", t.Name(), outputDir)
		}
	}()

	t.Log("this deployment used a postgres backend, verifying that postgres migrations ran properly")
	verifyPostgres(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)

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
		deployment, err = env.Cluster().Client().AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
		require.NoError(t, err)
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

	startPortForwarder(forwardCtx, t, env, secondary.Namespace, secondary.Name, "9777", "cmetrics")

	require.Never(t, func() bool {
		// if we are not the leader, we run no config pushes, and this metric string will not appear.
		return httpGetResponseContains(t, "http://localhost:9777/metrics", client, metrics.MetricNameConfigPushCount)
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

	var rebuiltPod corev1.Pod
	for _, pod := range podList.Items {
		if pod.Name != initialPod.Name && pod.Name != secondary.Name {
			rebuiltPod = pod
			startPortForwarder(forwardCtx, t, env, rebuiltPod.Namespace, rebuiltPod.Name, "9778", "cmetrics")
			break
		}
	}

	// Pass the test if exactly one of the pod becomes the leader, not limited to the original secondary pod.
	// Because in several times, the rebuilt pod (new pod created after initial pod deleted) became the leader.
	t.Logf("confirming there is exactly one pod that becomes leader and starts pushing configuration at %v", time.Now())
	require.Eventually(t, func() bool {
		leaderCount := 0
		if httpGetResponseContains(t, "http://localhost:9777/metrics", client, metrics.MetricNameConfigPushCount) {
			t.Logf("secondary pod %s is the leader at %v", secondary.Name, time.Now())
			leaderCount++
		}
		if httpGetResponseContains(t, "http://localhost:9778/metrics", client, metrics.MetricNameConfigPushCount) {
			t.Logf("rebuilt pod %s is the leader at %v", rebuiltPod.Name, time.Now())
			leaderCount++
		}
		t.Logf("expected exactly one leader, actual %d", leaderCount)
		return leaderCount == 1
	}, 2*time.Minute, time.Second)
}

const entPostgresPath = "../../deploy/single/all-in-one-postgres-enterprise.yaml"

func TestDeployAllInOneEnterprisePostgres(t *testing.T) {
	t.Log("configuring all-in-one-postgres-enterprise.yaml manifest test")
	if os.Getenv(kong.LicenseDataEnvVar) == "" {
		t.Skipf("no license available to test enterprise: %s was not provided", kong.LicenseDataEnvVar)
	}
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	addons = append(addons, metallb.New())

	addons = append(addons, buildImageLoadAddons(t, imageLoad, kongImageLoad)...)

	builder := setBuilderKubernetesVersion(t,
		environments.NewBuilder().WithAddons(addons...), clusterVersionStr)
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	createKongImagePullSecret(ctx, t, env)

	t.Logf("build a cleaner to dump diagnostics...")
	cleaner := clusters.NewCleaner(env.Cluster())
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Log("generating a superuser password")
	adminPassword, adminPasswordSecret, err := generateAdminPasswordSecret()
	require.NoError(t, err)

	t.Log("generating a license secret")
	licenseSecret, err := kong.GetLicenseSecretFromEnv()
	require.NoError(t, err)

	t.Log("deploying kong components")
	manifest, err := getTestManifest(t, entPostgresPath)
	require.NoError(t, err)
	_ = deployKong(ctx, t, env, manifest, licenseSecret, adminPasswordSecret)

	t.Log("this deployment used a postgres backend, verifying that postgres migrations ran properly")
	verifyPostgres(ctx, t, env)

	t.Log("running ingress tests to verify ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)

	t.Log("this deployment used enterprise kong, verifying that enterprise functionality was set up properly")
	verifyEnterprise(ctx, t, env, adminPassword)
	verifyEnterpriseWithPostgres(ctx, t, env, adminPassword)
}
