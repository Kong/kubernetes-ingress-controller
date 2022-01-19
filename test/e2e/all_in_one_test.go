//go:build e2e_tests
// +build e2e_tests

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/sethvargo/go-password/password"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/loadimage"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
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
)

var (
	imageOverride = os.Getenv("TEST_KONG_CONTROLLER_IMAGE_OVERRIDE")
	imageLoad     = os.Getenv("TEST_KONG_CONTROLLER_IMAGE_LOAD")
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
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.Parse(clusterVersionStr)
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
	_ = deployKong(ctx, t, env, manifest)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
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
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.Parse(clusterVersionStr)
		require.NoError(t, err)
		builder.WithKubernetesVersion(clusterVersion)
	}
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, env.Cleanup(ctx))
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

// Unsatisfied LoadBalancers have special handling, see
// https://github.com/Kong/kubernetes-ingress-controller/issues/2001
func TestDeployAllInOneDBLESSNoLoadBalancer(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	addons := []clusters.Addon{}
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.Parse(clusterVersionStr)
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
	_ = deployKong(ctx, t, env, manifest)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
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
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.Parse(clusterVersionStr)
		require.NoError(t, err)
		builder.WithKubernetesVersion(clusterVersion)
	}
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, env.Cleanup(ctx))
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
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.Parse(clusterVersionStr)
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
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.Parse(clusterVersionStr)
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

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
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
	startPortForwarder(forwardCtx, t, env, secondary, "9777", "cmetrics")
	require.Never(t, func() bool {
		req, err := http.NewRequest("GET", "http://localhost:9777/metrics", nil)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		// if we are not the leader, we run no config pushes, and this metric string will not appear
		return strings.Contains(string(body), metrics.MetricNameConfigPushCount)
	}, time.Minute, time.Second*10)

	t.Log("deleting the original replica and current leader")
	err = env.Cluster().Client().CoreV1().Pods(initialPod.Namespace).Delete(ctx, initialPod.Name, metav1.DeleteOptions{})
	require.NoError(t, err)

	t.Log("confirming the second replica becomes the leader and starts pushing configuration")
	require.Eventually(t, func() bool {
		req, err := http.NewRequest("GET", "http://localhost:9777/metrics", nil)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		return strings.Contains(string(body), metrics.MetricNameConfigPushCount)
	}, time.Minute, time.Second)
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
	if b, err := loadimage.NewBuilder().WithImage(imageLoad); err == nil {
		addons = append(addons, b.Build())
	}
	builder := environments.NewBuilder().WithAddons(addons...)
	if clusterVersionStr != "" {
		clusterVersion, err := semver.Parse(clusterVersionStr)
		require.NoError(t, err)
		builder.WithKubernetesVersion(clusterVersion)
	}
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, env.Cleanup(ctx))
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

// -----------------------------------------------------------------------------
// Private Functions - Test Helpers
// -----------------------------------------------------------------------------

const (
	httpBinImage     = "kennethreitz/httpbin"
	ingressClass     = "kong"
	namespace        = "kong"
	adminServiceName = "kong-admin"
)

func deployKong(ctx context.Context, t *testing.T, env environments.Environment, manifest io.Reader, additionalSecrets ...*corev1.Secret) *appsv1.Deployment {
	t.Log("creating a tempfile for kubeconfig")
	kubeconfig, err := generators.NewKubeConfigForRestConfig(env.Name(), env.Cluster().Config())
	require.NoError(t, err)
	kubeconfigFile, err := os.CreateTemp(os.TempDir(), "manifest-tests-kubeconfig-")
	require.NoError(t, err)
	defer os.Remove(kubeconfigFile.Name())
	defer kubeconfigFile.Close()

	t.Log("dumping kubeconfig to tempfile")
	written, err := kubeconfigFile.Write(kubeconfig)
	require.NoError(t, err)
	require.Equal(t, len(kubeconfig), written)
	kubeconfigFilename := kubeconfigFile.Name()

	t.Log("waiting for testing environment to be ready")
	require.NoError(t, <-env.WaitForReady(ctx))

	t.Log("creating the kong namespace")
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kong"}}
	_, err = env.Cluster().Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if !kerrors.IsAlreadyExists(err) {
		require.NoError(t, err)
	}

	t.Logf("deploying any supplemental secrets (found: %d)", len(additionalSecrets))
	for _, secret := range additionalSecrets {
		_, err := env.Cluster().Client().CoreV1().Secrets("kong").Create(ctx, secret, metav1.CreateOptions{})
		if !kerrors.IsAlreadyExists(err) {
			require.NoError(t, err)
		}
	}

	t.Log("deploying the manifest to the cluster")
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfigFilename, "apply", "-f", "-")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = manifest
	require.NoError(t, cmd.Run(), fmt.Sprintf("STDOUT=(%s), STDERR=(%s)", stdout.String(), stderr.String()))

	t.Log("waiting for kong to be ready")
	var deployment *appsv1.Deployment
	require.Eventually(t, func() bool {
		deployment, err = env.Cluster().Client().AppsV1().Deployments(namespace).Get(ctx, "ingress-kong", metav1.GetOptions{})
		require.NoError(t, err)
		return deployment.Status.ReadyReplicas == *deployment.Spec.Replicas
	}, kongComponentWait, time.Second)
	return deployment
}

func deployIngress(ctx context.Context, t *testing.T, env environments.Environment) {
	c, err := clientset.NewForConfig(env.Cluster().Config())
	assert.NoError(t, err)
	t.Log("deploying an HTTP service to test the ingress controller and proxy")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	getString := "GET"
	king := &kongv1.KongIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testki",
			Namespace: corev1.NamespaceDefault,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Route: &kongv1.KongIngressRoute{
			Methods: []*string{&getString},
		},
	}
	king, err = c.ConfigurationV1().KongIngresses(corev1.NamespaceDefault).Create(ctx, king,
		metav1.CreateOptions{})
	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/httpbin", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
		"konghq.com/override":       "testki",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), corev1.NamespaceDefault, ingress))
}

func verifyIngress(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Log("finding the kong proxy service ip")
	svc, err := env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, "kong-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	proxyIP := ""
	require.NotEqual(t, svc.Spec.Type, svc.Spec.ClusterIP)
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			proxyIP = svc.Status.LoadBalancer.Ingress[0].IP
		}
	}
	// the above failed to find an address. either the LB didn't provision or we're using a NodePort
	if proxyIP == "" {
		var port int32
		for _, sport := range svc.Spec.Ports {
			if sport.Name == "kong-proxy" || sport.Name == "proxy" {
				port = sport.NodePort
			}
		}
		var extAddrs []string
		var intAddrs []string
		nodes, err := env.Cluster().Client().CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		require.NoError(t, err)
		for _, node := range nodes.Items {
			for _, naddr := range node.Status.Addresses {
				if naddr.Type == corev1.NodeExternalIP {
					extAddrs = append(extAddrs, naddr.Address)
				}
				if naddr.Type == corev1.NodeInternalIP {
					extAddrs = append(intAddrs, naddr.Address)
				}
			}
		}
		// local clusters (KIND, minikube) typically provide no external addresses, but their internal addresses are
		// routeable from their host. We prefer external addresses if they're available, but fall back to internal
		// in their absence
		if len(extAddrs) > 0 {
			proxyIP = fmt.Sprintf("%v:%v", extAddrs[0], port)
		} else if len(intAddrs) > 0 {
			proxyIP = fmt.Sprintf("%v:%v", intAddrs[0], port)
		} else {
			assert.Fail(t, "both extAddrs and intAddrs are empty")
		}
	}

	t.Logf("waiting for route from Ingress to be operational at http://%s/httpbin", proxyIP)
	httpc := http.Client{Timeout: time.Second * 10}
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("http://%s/httpbin", proxyIP))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			if !strings.Contains(b.String(), "<title>httpbin.org</title>") {
				return false
			}
		}
		// verify the KongIngress method restriction
		fakeData := url.Values{}
		fakeData.Set("foo", "bar")
		resp, err = httpc.PostForm(fmt.Sprintf("http://%s/httpbin", proxyIP), fakeData)
		if err != nil {
			return false
		}
		if resp.StatusCode == http.StatusNotFound {
			return true
		}
		return false
	}, ingressWait, time.Second)
}

// verifyEnterprise performs some basic tests of the Kong Admin API in the provided
// environment to ensure that the Admin API that responds is in fact the enterprise
// version of Kong.
func verifyEnterprise(ctx context.Context, t *testing.T, env environments.Environment, adminPassword string) {
	t.Log("finding the ip address for the admin API")
	service, err := env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, adminServiceName, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, 1, len(service.Status.LoadBalancer.Ingress))
	adminIP := service.Status.LoadBalancer.Ingress[0].IP

	t.Log("building a GET request to gather admin api information")
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/", adminIP), nil)
	require.NoError(t, err)
	req.Header.Set("Kong-Admin-Token", adminPassword)

	t.Log("pulling the admin api information")
	adminOutput := struct {
		Version string `json:"version"`
	}{}
	httpc := http.Client{Timeout: time.Second * 10}
	require.Eventually(t, func() bool {
		// at the time of writing it was seen that the admin API had
		// brief timing windows where it could respond 200 OK but
		// the API version data would not be populated and the JSON
		// decode would fail. Thus this check actually waits until
		// the response body is fully decoded with a non-empty value
		// before considering this complete.
		resp, err := httpc.Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		if resp.StatusCode != http.StatusOK {
			return false
		}
		if err := json.Unmarshal(body, &adminOutput); err != nil {
			return false
		}
		return adminOutput.Version != ""
	}, adminAPIWait, time.Second)
	require.True(t, strings.Contains(adminOutput.Version, "enterprise-edition"))
}

func verifyEnterpriseWithPostgres(ctx context.Context, t *testing.T, env environments.Environment, adminPassword string) {
	t.Log("finding the ip address for the admin API")
	service, err := env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, adminServiceName, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, 1, len(service.Status.LoadBalancer.Ingress))
	adminIP := service.Status.LoadBalancer.Ingress[0].IP

	t.Log("building a POST request to create a new kong workspace")
	form := url.Values{"name": {"kic-e2e-tests"}}
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://%s/workspaces", adminIP), strings.NewReader(form.Encode()))
	require.NoError(t, err)
	req.Header.Set("Kong-Admin-Token", adminPassword)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	t.Log("creating a workspace to validate enterprise functionality")
	httpc := http.Client{Timeout: time.Second * 10}
	resp, err := httpc.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode, fmt.Sprintf("STATUS=(%s), BODY=(%s)", resp.Status, string(body)))
}

func verifyPostgres(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Log("verifying that postgres pod was deployed and is running")
	postgresPod, err := env.Cluster().Client().CoreV1().Pods(namespace).Get(ctx, "postgres-0", metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, corev1.PodRunning, postgresPod.Status.Phase)

	t.Log("verifying that all migrations ran properly")
	migrationJob, err := env.Cluster().Client().BatchV1().Jobs(namespace).Get(ctx, "kong-migrations", metav1.GetOptions{})
	require.NoError(t, err)
	require.GreaterOrEqual(t, migrationJob.Status.Succeeded, int32(1))
}

// startPortForwarder runs "kubectl port-forward" in the background. It stops the forward when the provided context
// ends
func startPortForwarder(ctx context.Context, t *testing.T, env environments.Environment, pod corev1.Pod, localPort,
	targetPort string) {
	kubeconfig, err := generators.NewKubeConfigForRestConfig(env.Name(), env.Cluster().Config())
	require.NoError(t, err)
	kubeconfigFile, err := os.CreateTemp(os.TempDir(), "portforward-tests-kubeconfig-")
	require.NoError(t, err)
	defer os.Remove(kubeconfigFile.Name())
	defer kubeconfigFile.Close()
	written, err := kubeconfigFile.Write(kubeconfig)
	require.NoError(t, err)
	require.Equal(t, len(kubeconfig), written)
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfigFile.Name(), "port-forward", "-n", pod.Namespace, pod.Name, "9777:cmetrics")
	t.Logf("forwarding port %s to %s/%s:%s", localPort, pod.Namespace, pod.Name, targetPort)
	if startErr := cmd.Start(); startErr != nil {
		startOutput, outputErr := cmd.Output()
		assert.NoError(t, outputErr)
		require.NoError(t, startErr, string(startOutput))
	}
	require.Eventually(t, func() bool {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%s", localPort))
		if err == nil {
			conn.Close()
			return true
		}
		return false
	}, kongComponentWait, time.Second)
}

// -----------------------------------------------------------------------------
// Private Functions - Utilities
// -----------------------------------------------------------------------------

const (
	// adminPasswordSecretName is the name of the secret which will house the admin
	// API admin password.
	adminPasswordSecretName = "kong-enterprise-superuser-password"
)

func generateAdminPasswordSecret() (string, *corev1.Secret, error) {
	adminPassword, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return "", nil, err
	}

	return adminPassword, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: adminPasswordSecretName,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"password": []byte(adminPassword),
		},
	}, nil
}

// exposeAdminAPI will override the KONG_ADMIN_LISTEN for the cluster's proxy to expose the
// Admin API via a service. Some deployments only expose this on localhost by default as there's
// no authentication, so note that this is only for testing environment purposes.
func exposeAdminAPI(ctx context.Context, t *testing.T, env environments.Environment) *corev1.Service {
	t.Log("updating the proxy container KONG_ADMIN_LISTEN to expose the admin api")
	deployment, err := env.Cluster().Client().AppsV1().Deployments(namespace).Get(ctx, "ingress-kong", metav1.GetOptions{})
	require.NoError(t, err)
	for i, containerSpec := range deployment.Spec.Template.Spec.Containers {
		if containerSpec.Name == "proxy" {
			for j, envVar := range containerSpec.Env {
				if envVar.Name == "KONG_ADMIN_LISTEN" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = "0.0.0.0:8001, 0.0.0.0:8444 ssl"
				}
			}
		}
	}
	deployment, err = env.Cluster().Client().AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("creating a loadbalancer service for the admin API")
	svcPorts := []corev1.ServicePort{{
		Name:       "proxy",
		Protocol:   corev1.ProtocolTCP,
		TargetPort: intstr.IntOrString{IntVal: 8001},
		Port:       80,
	}}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: adminServiceName,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeLoadBalancer,
			Selector: deployment.Spec.Selector.MatchLabels,
			Ports:    svcPorts,
		},
	}
	service, err = env.Cluster().Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("waiting for loadbalancer ip to provision")
	require.Eventually(t, func() bool {
		service, err = env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, service.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return len(service.Status.LoadBalancer.Ingress) == 1
	}, time.Minute, time.Second)

	return service
}

// getTestManifest checks if a controller image override is set. If not, it returns the original provided path.
// If an override is set, it runs a kustomize patch that replaces the controller image with the override image and
// returns the modified manifest path. If there is any issue patching the manifest, it will log the issue and return
// the original provided path
func getTestManifest(t *testing.T, baseManifestPath string) (io.Reader, error) {
	var imagetag string
	if imageLoad != "" {
		imagetag = imageLoad
	} else {
		imagetag = imageOverride
	}
	if imagetag == "" {
		return os.Open(baseManifestPath)
	}
	split := strings.Split(imagetag, ":")
	if len(split) != 2 {
		t.Logf("could not parse override image '%v', using default manifest %v", imagetag, baseManifestPath)
		return os.Open(baseManifestPath)
	}
	modified, err := patchControllerImage(baseManifestPath, split[0], split[1])
	if err != nil {
		t.Logf("failed patching override image '%v' (%v), using default manifest %v", imagetag, err, baseManifestPath)
		return os.Open(baseManifestPath)
	}
	t.Logf("using modified %v manifest", baseManifestPath)
	return modified, nil
}

const imageKustomizationContents = `resources:
- base.yaml
images:
- name: kong/kubernetes-ingress-controller
  newName: %v
  newTag: '%v'
`

// patchControllerImage takes a manifest, image, and tag and runs kustomize to replace the
// kong/kubernetes-ingress-controller image with the provided image. It returns the location of kustomize's output
func patchControllerImage(baseManifestPath string, image string, tag string) (io.Reader, error) {
	workDir, err := os.MkdirTemp("", "kictest.")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workDir)
	orig, err := ioutil.ReadFile(baseManifestPath)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(filepath.Join(workDir, "base.yaml"), orig, 0600)
	if err != nil {
		return nil, err
	}
	kustomization := []byte(fmt.Sprintf(imageKustomizationContents, image, tag))
	err = os.WriteFile(filepath.Join(workDir, "kustomization.yaml"), kustomization, 0600)
	if err != nil {
		return nil, err
	}
	kustomized, err := kustomizeManifest(workDir)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(kustomized), nil
}

// kustomizeManifest runs kustomize on a path and returns the YAML output
func kustomizeManifest(path string) ([]byte, error) {
	k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	m, err := k.Run(filesys.MakeFsOnDisk(), path)
	if err != nil {
		return []byte{}, err
	}
	return m.AsYaml()
}

func getCurrentGitTag(path string) (semver.Version, error) {
	cmd := exec.Command("git", "describe", "--tags")
	cmd.Dir = path
	tagBytes, _ := cmd.Output()
	tag, err := semver.ParseTolerant(string(tagBytes))
	if err != nil {
		return semver.Version{}, err
	}
	return tag, nil
}

func getPreviousGitTag(path string, cur semver.Version) (semver.Version, error) {
	var tags []semver.Version
	cmd := exec.Command("git", "tag")
	cmd.Dir = path
	tagsBytes, _ := cmd.Output()
	foo := strings.Split(string(tagsBytes), "\n")
	for _, tag := range foo {
		ver, err := semver.ParseTolerant(tag)
		if err == nil {
			tags = append(tags, ver)
		}
	}
	sort.Slice(tags, func(i, j int) bool { return tags[i].LT(tags[j]) })
	curIndex := sort.Search(len(tags), func(i int) bool { return tags[i].EQ(cur) })
	if curIndex == 0 {
		return tags[curIndex], nil
	}
	return tags[curIndex-1], nil
}
