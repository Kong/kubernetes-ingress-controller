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
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
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
	adminServiceName = "kong-admin"

	tcpEchoPort    = 1025
	tcpListnerPort = 8888
)

func getEnvironmentBuilder(ctx context.Context) (*environments.Builder, error) {
	if existingCluster == "" {
		fmt.Printf("INFO: no existing cluster provided, creating a new one for %q type\n", clusterProvider)
		switch clusterProvider {
		case string(gke.GKEClusterType):
			fmt.Println("INFO: creating a GKE cluster builder")
			return createGKEBuilder()
		default:
			fmt.Println("INFO: creating a Kind cluster builder")
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

	fmt.Printf("INFO: using existing %s cluster %s\n", clusterType, clusterName)
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

func createGKEBuilder() (*environments.Builder, error) {
	var (
		name        = "e2e-" + uuid.NewString()
		gkeCreds    = os.Getenv(gke.GKECredsVar)
		gkeProject  = os.Getenv(gke.GKEProjectVar)
		gkeLocation = os.Getenv(gke.GKELocationVar)
	)

	fmt.Printf("INFO: cluster name: %q\n", name)

	clusterBuilder := gke.
		NewBuilder([]byte(gkeCreds), gkeProject, gkeLocation).
		WithName(name).
		WithWaitForTeardown(true).
		WithCreateSubnet(true).
		WithLabels(gkeTestClusterLabels())

	if clusterVersionStr != "" {
		k8sVersion, err := semver.Parse(strings.TrimPrefix(clusterVersionStr, "v"))
		if err != nil {
			return nil, err
		}

		clusterBuilder = clusterBuilder.WithClusterMinorVersion(k8sVersion.Major, k8sVersion.Minor)
	}

	return environments.NewBuilder().WithClusterBuilder(clusterBuilder), nil
}

func deployKong(ctx context.Context, t *testing.T, env environments.Environment, manifest io.Reader, additionalSecrets ...*corev1.Secret) *appsv1.Deployment {
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
		_, err := env.Cluster().Client().CoreV1().Secrets("kong").Create(ctx, secret, metav1.CreateOptions{})
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

	t.Log("waiting for kong to be ready")
	var deployment *appsv1.Deployment
	require.Eventually(t, func() bool {
		deployment, err = env.Cluster().Client().AppsV1().Deployments(namespace).Get(ctx, "ingress-kong", metav1.GetOptions{})
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
	return deployment
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
		if status.Name == "proxy" {
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
			if status.Name == "proxy" {
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
