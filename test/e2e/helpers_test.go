//go:build e2e_tests || istio_tests
// +build e2e_tests istio_tests

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/google/uuid"
	"github.com/kong/deck/dump"
	gokong "github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/loadimage"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/gke"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

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

	ingressClass     = "kong"
	namespace        = "kong"
	adminServiceName = "kong-admin-lb"

	tcpEchoPort    = 1025
	tcpListnerPort = 8888

	// controllerDeploymentName is the name of the controller deployment in all manifests variants.
	controllerDeploymentName = "ingress-kong"

	// controllerContainerName is the name of the controller container in all manifests variants.
	controllerContainerName = "ingress-controller"

	// proxyContainerName is the name of the proxy container in all manifests variants.
	proxyContainerName = "proxy"
)

// setupE2ETest builds a testing environment for the E2E test. It also sets up the environment's teardown and test
// context cancellation. It can accept optional addons to be passed to the environment builder.
func setupE2ETest(t *testing.T, addons ...clusters.Addon) (context.Context, environments.Environment) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Log("building test cluster and environment")
	builder, err := getEnvironmentBuilder(ctx, t)
	require.NoError(t, err)
	env, err := builder.WithAddons(addons...).Build(ctx)
	require.NoError(t, err)
	logClusterInfo(t, env.Cluster())

	t.Cleanup(func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	})

	return ctx, env
}

func getEnvironmentBuilder(ctx context.Context, t *testing.T) (*environments.Builder, error) {
	t.Helper()

	if existingCluster == "" {
		t.Logf("no existing cluster provided, creating a new one for %q type", clusterProvider)
		switch clusterProvider {
		case string(gke.GKEClusterType):
			t.Log("creating a GKE cluster builder")
			return createGKEBuilder(t)
		default:
			t.Log("creating a Kind cluster builder")
			return createKINDBuilder(), nil
		}
	}

	clusterParts := strings.Split(existingCluster, ":")
	if len(clusterParts) < 2 {
		return nil, fmt.Errorf("expected existing cluster in format <type>:<name>, got %s", existingCluster)
	}

	clusterType, clusterName := clusterParts[0], clusterParts[1]
	if clusterVersionStr != "" {
		return nil, fmt.Errorf("cannot provide cluster version with existing cluster")
	}

	t.Logf("using existing %s cluster %s", clusterType, clusterName)
	switch clusterType {
	case string(kind.KindClusterType):
		return createExistingKINDBuilder(clusterName)
	case string(gke.GKEClusterType):
		return createExistingGKEBuilder(ctx, clusterName)
	default:
		return nil, fmt.Errorf("unrecognized cluster type %s", clusterType)
	}
}

func createKINDBuilder() *environments.Builder {
	builder := environments.NewBuilder()
	clusterBuilder := kind.NewBuilder()
	if clusterVersionStr != "" {
		clusterVersion := semver.MustParse(strings.TrimPrefix(clusterVersionStr, "v"))
		clusterBuilder = clusterBuilder.WithClusterVersion(clusterVersion)
	}
	builder = builder.WithClusterBuilder(clusterBuilder)
	builder = builder.WithAddons(metallb.New())
	builder = builder.WithAddons(buildImageLoadAddons(imageLoad, kongImageLoad)...)
	return builder
}

func createExistingKINDBuilder(name string) (*environments.Builder, error) {
	builder := environments.NewBuilder()
	cluster, err := kind.NewFromExisting(name)
	if err != nil {
		return nil, err
	}
	builder = builder.WithExistingCluster(cluster)
	builder = builder.WithAddons(metallb.New())
	builder = builder.WithAddons(buildImageLoadAddons(imageLoad, kongImageLoad)...)
	return builder, nil
}

func createExistingGKEBuilder(ctx context.Context, name string) (*environments.Builder, error) {
	cluster, err := gke.NewFromExistingWithEnv(ctx, name)
	if err != nil {
		return nil, err
	}
	builder := environments.NewBuilder()
	builder = builder.WithExistingCluster(cluster)
	return builder, nil
}

func createGKEBuilder(t *testing.T) (*environments.Builder, error) {
	t.Helper()

	var (
		name        = "e2e-" + uuid.NewString()
		gkeCreds    = os.Getenv(gke.GKECredsVar)
		gkeProject  = os.Getenv(gke.GKEProjectVar)
		gkeLocation = os.Getenv(gke.GKELocationVar)
	)

	t.Logf("creating GKE cluster, name: %s", name)

	clusterBuilder := gke.
		NewBuilder([]byte(gkeCreds), gkeProject, gkeLocation).
		WithName(name).
		WithWaitForTeardown(testenv.WaitForClusterDelete()).
		WithCreateSubnet(true).
		WithLabels(gkeTestClusterLabels())

	if clusterVersionStr != "" {
		k8sVersion, err := semver.Parse(strings.TrimPrefix(clusterVersionStr, "v"))
		if err != nil {
			return nil, err
		}

		t.Logf("creating GKE cluster, with requested version: %s", k8sVersion)
		clusterBuilder = clusterBuilder.WithClusterMinorVersion(k8sVersion.Major, k8sVersion.Minor)
	}

	return environments.NewBuilder().WithClusterBuilder(clusterBuilder), nil
}

func deployKong(ctx context.Context, t *testing.T, env environments.Environment, manifest io.Reader, additionalSecrets ...*corev1.Secret) {
	t.Helper()

	t.Log("waiting for testing environment to be ready")
	require.NoError(t, <-env.WaitForReady(ctx))

	t.Log("creating the kong namespace")
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kong"}}
	_, err := env.Cluster().Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if !apierrors.IsAlreadyExists(err) {
		require.NoError(t, err)
	}

	t.Logf("deploying any supplemental secrets (found: %d)", len(additionalSecrets))
	for _, secret := range additionalSecrets {
		_, err := env.Cluster().Client().CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
		if !apierrors.IsAlreadyExists(err) {
			require.NoError(t, err)
		}
	}

	t.Log("deploying the manifest to the cluster")
	kubeconfigFilename := getTemporaryKubeconfig(t, env)
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfigFilename, "apply", "-f", "-")
	cmd.Stdin = manifest
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	t.Log("waiting for controller to be ready")
	var deployment *appsv1.Deployment
	require.Eventually(t, func() bool {
		deployment, err = env.Cluster().Client().AppsV1().Deployments(namespace).Get(ctx, controllerDeploymentName, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return deployment.Status.ReadyReplicas == *deployment.Spec.Replicas
	}, kongComponentWait, time.Second,
		func() string {
			if deployment == nil {
				return ""
			}
			return fmt.Sprintf(
				"deployment %s: ready replicas %d, spec replicas: %d",
				deployment.Name, deployment.Status.ReadyReplicas, *deployment.Spec.Replicas,
			)
		}(),
	)
}

// Deployments represent the deployments that are deployed by the all-in-one manifests.
type Deployments struct {
	ProxyNN      types.NamespacedName
	ControllerNN types.NamespacedName
}

// GetProxy gets the proxy deployment from the cluster.
func (d Deployments) GetProxy(ctx context.Context, t *testing.T, env environments.Environment) *appsv1.Deployment {
	t.Helper()
	deployment, err := env.Cluster().Client().AppsV1().Deployments(d.ProxyNN.Namespace).Get(ctx, d.ProxyNN.Name, metav1.GetOptions{})
	require.NoError(t, err)
	return deployment
}

// GetController gets the controller deployment from the cluster.
func (d Deployments) GetController(ctx context.Context, t *testing.T, env environments.Environment) *appsv1.Deployment {
	t.Helper()
	deployment, err := env.Cluster().Client().AppsV1().Deployments(d.ControllerNN.Namespace).Get(ctx, d.ControllerNN.Name, metav1.GetOptions{})
	require.NoError(t, err)
	return deployment
}

// getManifestDeployments returns the deployments for the proxy and controller that are expected to be deployed for a given
// manifest.
func getManifestDeployments(manifestPath string) Deployments {
	proxyDeploymentName := getProxyDeploymentName(manifestPath)
	return Deployments{
		ProxyNN: types.NamespacedName{
			Namespace: namespace,
			Name:      proxyDeploymentName,
		},
		ControllerNN: types.NamespacedName{
			Namespace: namespace,
			Name:      controllerDeploymentName,
		},
	}
}

// getProxyDeploymentName returns the name of the proxy deployment that is expected to be deployed for a given manifest.
func getProxyDeploymentName(manifestPath string) string {
	const (
		singlePodDeploymentName = controllerDeploymentName
		multiPodDeploymentName  = "proxy-kong"
	)

	// special case for the old dbless-legacy that's still using a single-pod deployment
	if strings.Contains(manifestPath, "dbless-legacy") {
		return singlePodDeploymentName
	}
	// all non-legacy dbless manifests use a multi-pod deployment
	if strings.Contains(manifestPath, "dbless") {
		return multiPodDeploymentName
	}

	return singlePodDeploymentName
}

func deployIngress(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Helper()

	c, err := clientset.NewForConfig(env.Cluster().Config())
	assert.NoError(t, err)
	t.Log("deploying an HTTP service to test the ingress controller and proxy")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
	_, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	kongIngressName := uuid.NewString()
	king := &kongv1.KongIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kongIngressName,
			Namespace: corev1.NamespaceDefault,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Route: &kongv1.KongIngressRoute{
			Methods: []*string{lo.ToPtr("GET")},
		},
	}
	_, err = c.ConfigurationV1().KongIngresses(corev1.NamespaceDefault).Create(ctx, king, metav1.CreateOptions{})
	require.NoError(t, err)
	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/httpbin", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
		"konghq.com/override":       kongIngressName,
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), corev1.NamespaceDefault, ingress))
}

func verifyIngress(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Helper()

	t.Log("finding the kong proxy service ip")
	proxyIP := getKongProxyIP(ctx, t, env)

	t.Logf("waiting for route from Ingress to be operational at http://%s/httpbin", proxyIP)
	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("http://%s/httpbin", proxyIP))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			_, err := b.ReadFrom(resp.Body)
			if err != nil {
				return false
			}
			if !strings.Contains(b.String(), "<title>httpbin.org</title>") {
				return false
			}
		} else {
			return false
		}
		// verify the KongIngress method restriction
		fakeData := url.Values{}
		fakeData.Set("foo", "bar")
		resp, err = helpers.DefaultHTTPClient().PostForm(fmt.Sprintf("http://%s/httpbin", proxyIP), fakeData)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusNotFound
	}, ingressWait, 100*time.Millisecond)
}

// requireIngressConfiguredInAdminAPIEventually ensures all expected Kong Admin API resources are created for the Ingress
// deployed with deployIngress helper function.
func requireIngressConfiguredInAdminAPIEventually(
	ctx context.Context,
	t *testing.T,
	kongClient *gokong.Client,
) {
	t.Helper()

	require.Eventually(t, func() bool {
		d, err := dump.Get(ctx, kongClient, dump.Config{})
		if err != nil {
			t.Logf("failed dumping config: %s", err)
			return false
		}
		if len(d.Services) != 1 {
			t.Log("still no service found...")
			return false
		}
		if len(d.Routes) != 1 {
			t.Log("still no route found...")
			return false
		}
		if d.Services[0].ID == nil ||
			d.Routes[0].Service.ID == nil ||
			*d.Services[0].ID != *d.Routes[0].Service.ID {
			t.Log("still no matching service found...")
			return false
		}
		if len(d.Targets) != 1 {
			t.Log("still no target found...")
			return false
		}
		if len(d.Upstreams) != 1 {
			t.Log("still no upstream found...")
			return false
		}
		return true
	}, time.Minute*3, time.Second, "%q didn't get the config", kongClient.BaseRootURL())
}

// verifyEnterprise performs some basic tests of the Kong Admin API in the provided
// environment to ensure that the Admin API that responds is in fact the enterprise
// version of Kong.
func verifyEnterprise(ctx context.Context, t *testing.T, env environments.Environment, adminPassword string) {
	t.Helper()

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

	require.Eventually(t, func() bool {
		// at the time of writing it was seen that the admin API had
		// brief timing windows where it could respond 200 OK but
		// the API version data would not be populated and the JSON
		// decode would fail. Thus this check actually waits until
		// the response body is fully decoded with a non-empty value
		// before considering this complete.
		resp, err := helpers.DefaultHTTPClient().Do(req)
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
	if string(adminOutput.Version[0]) == "3" {
		// 3.x removed the "-enterprise-edition" string but provided no other indication that something is enterprise
		require.Len(t, strings.Split(adminOutput.Version, "."), 4,
			fmt.Sprintf("actual kong version: %s", adminOutput.Version))
	} else {
		require.Contains(t, adminOutput.Version, "enterprise-edition",
			fmt.Sprintf("actual kong version: %s", adminOutput.Version))
	}
}

func verifyEnterpriseWithPostgres(ctx context.Context, t *testing.T, env environments.Environment, adminPassword string) {
	t.Helper()

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

	resp, err := helpers.DefaultHTTPClient().Do(req)
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

// killKong kills the Kong container in a given Pod and returns when it has restarted.
func killKong(ctx context.Context, t *testing.T, env environments.Environment, pod *corev1.Pod) {
	t.Helper()

	var orig, after int32
	for _, status := range pod.Status.ContainerStatuses {
		if status.Name == proxyContainerName {
			orig = status.RestartCount
		}
	}

	kubeconfig := getTemporaryKubeconfig(t, env)
	cmd := exec.Command("kubectl", "--kubeconfig", kubeconfig, "exec", "-n", pod.Namespace, pod.Name, "--", "bash", "-c", "kill 1")
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	require.NoErrorf(t, err, "kill failed: STDOUT(%s) STDERR(%s)", stdout.String(), stderr.String())
	require.Eventually(t, func() bool {
		pod, err = env.Cluster().Client().CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, status := range pod.Status.ContainerStatuses {
			if status.Name == proxyContainerName {
				if status.RestartCount > orig {
					after = status.RestartCount
					return true
				}
			}
		}
		return false
	}, kongComponentWait, time.Second)
	t.Logf("kong container has %v restart after kill", after)
}

// buildImageLoadAddons creates addons to load KIC and kong images.
func buildImageLoadAddons(images ...string) []clusters.Addon {
	addons := []clusters.Addon{}
	for _, image := range images {
		if image != "" {
			// https://github.com/Kong/kubernetes-testing-framework/issues/440 this error only occurs if image == ""
			// it will eventually be removed from the WithImage return signature
			b, _ := loadimage.NewBuilder().WithImage(image)
			addons = append(addons, b.Build())
		}
	}
	return addons
}

// createKongImagePullSecret creates the image pull secret
// `kong-enterprise-edition-docker` for kong enterprise image
// from env TEST_KONG_PULL_USERNAME and TEST_KONG_PULL_PASSWORD.
func createKongImagePullSecret(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Helper()

	if kongImagePullUsername == "" || kongImagePullPassword == "" {
		return
	}
	kubeconfigFilename := getTemporaryKubeconfig(t, env)

	const secretName = "kong-enterprise-edition-docker"
	cmd := exec.CommandContext(
		ctx,
		"kubectl", "--kubeconfig", kubeconfigFilename,
		"create", "secret", "docker-registry", secretName,
		"--docker-username="+kongImagePullUsername,
		"--docker-password="+kongImagePullPassword,
	)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "command output: "+string(out))
}

// getContainerInPodSpec returns the spec of container having the given name.
// returns nil if there is no such container.
func getContainerInPodSpec(podSpec *corev1.PodSpec, name string) *corev1.Container {
	for i, container := range podSpec.Containers {
		if container.Name == name {
			return &podSpec.Containers[i]
		}
	}
	return nil
}

// getEnvValueInContainer returns the value of specified environment variable in the container.
// if there are multiple envs with same value, return the last one which is actually effective.
// returns empty string if the env os not found.
func getEnvValueInContainer(container *corev1.Container, name string) string {
	if container == nil {
		return ""
	}
	value := ""
	for _, env := range container.Env {
		if env.Name == name {
			value = env.Value
		}
	}
	return value
}

// getTemporaryKubeconfig dumps an environment's kubeconfig to a temporary file.
func getTemporaryKubeconfig(t *testing.T, env environments.Environment) string {
	t.Helper()

	t.Log("creating a tempfile for kubeconfig")
	kubeconfig, err := generators.NewKubeConfigForRestConfig(env.Name(), env.Cluster().Config())
	require.NoError(t, err)
	kubeconfigFile, err := os.CreateTemp(os.TempDir(), "manifest-tests-kubeconfig-")
	require.NoError(t, err)
	defer kubeconfigFile.Close()
	t.Cleanup(func() {
		assert.NoError(t, os.Remove(kubeconfigFile.Name()))
	})

	t.Logf("dumping kubeconfig to tempfile: %s", kubeconfigFile.Name())
	written, err := kubeconfigFile.Write(kubeconfig)
	require.NoError(t, err)
	require.Equal(t, len(kubeconfig), written)

	return kubeconfigFile.Name()
}

func runOnlyOnKindClusters(t *testing.T) {
	t.Helper()

	existingClusterIsKind := strings.Split(existingCluster, ":")[0] == string(kind.KindClusterType)
	clusterProviderIsKind := clusterProvider == "" || clusterProvider == string(kind.KindClusterType)

	if !existingClusterIsKind || !clusterProviderIsKind {
		t.Skip("test is supported only on Kind clusters")
	}
}

// listPodsByLabels returns a list of pods in the given namespace that match the given labels map.
func listPodsByLabels(
	ctx context.Context, env environments.Environment, namespace string, podLabels map[string]string,
) ([]corev1.Pod, error) {
	podClient := env.Cluster().Client().CoreV1().Pods(namespace)
	selector := labels.NewSelector()

	for k, v := range podLabels {
		req, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			return nil, err
		}
		selector = selector.Add(*req)
	}

	podList, err := podClient.List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}

// scaleDeployment scales the deployment to the given number of replicas and waits for the replicas to be ready.
func scaleDeployment(ctx context.Context, t *testing.T, env environments.Environment, deployment types.NamespacedName, replicas int32) {
	t.Helper()

	scale := &autoscalingv1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: deployment.Namespace,
			Name:      deployment.Name,
		},
		Spec: autoscalingv1.ScaleSpec{
			Replicas: replicas,
		},
	}
	deployments := env.Cluster().Client().AppsV1().Deployments(deployment.Namespace)
	_, err := deployments.UpdateScale(ctx, deployment.Name, scale, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		deployment, err := deployments.Get(ctx, deployment.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return deployment.Status.ReadyReplicas == replicas
	}, time.Minute*3, time.Second, "deployment %s did not scale to %d replicas", deployment.Name, replicas)
}
