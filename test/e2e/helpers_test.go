//go:build e2e_tests || istio_tests || performance_tests

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
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
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
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	conststest "github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
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

	tcpListenerPort = ktfkong.DefaultTCPServicePort

	// controllerDeploymentName is the name of the controller deployment in all manifests variants.
	controllerDeploymentName = "ingress-kong"

	// controllerContainerName is the name of the controller container in all manifests variants.
	controllerContainerName = "ingress-controller"

	// proxyContainerName is the name of the proxy container in all manifests variants.
	proxyContainerName = "proxy"

	// migrationsJobName is the name of the migrations job in postgres manifests variant.
	migrationsJobName = "kong-migrations"

	// echoPath is the legit echo path to use for the echo service.
	echoPath = "/echo"

	// badEchoPath is a wrong path to use for testing ingress misconfiguration.
	badEchoPath = "/~/echo/**"
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

	t.Logf("deploying KIC Incubator CRDs from %s (since they are not packaged with base CRDs)", conststest.IncubatorCRDKustomizeDir)
	require.NoError(t, clusters.KustomizeDeployForCluster(ctx, env.Cluster(), conststest.IncubatorCRDKustomizeDir))

	t.Cleanup(func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	})

	return ctx, env
}

func getEnvironmentBuilder(ctx context.Context, t *testing.T) (*environments.Builder, error) {
	t.Helper()

	if testenv.ExistingClusterName() == "" {
		t.Logf("no existing cluster provided, creating a new one for %q type", testenv.ClusterProvider())
		switch testenv.ClusterProvider() {
		case string(gke.GKEClusterType):
			t.Log("creating a GKE cluster builder")
			return createGKEBuilder(t)
		default:
			t.Log("creating a Kind cluster builder")
			return createKINDBuilder(t), nil
		}
	}

	clusterParts := strings.Split(testenv.ExistingClusterName(), ":")
	if len(clusterParts) < 2 {
		return nil, fmt.Errorf("expected existing cluster in format <type>:<name>, got %s", testenv.ExistingClusterName())
	}

	clusterType, clusterName := clusterParts[0], clusterParts[1]
	if testenv.ClusterVersion() != "" {
		return nil, fmt.Errorf("cannot provide cluster version with existing cluster")
	}

	t.Logf("using existing %s cluster %s", clusterType, clusterName)
	switch clusterType {
	case string(kind.KindClusterType):
		return createExistingKINDBuilder(t, clusterName)
	case string(gke.GKEClusterType):
		return createExistingGKEBuilder(ctx, clusterName)
	default:
		return nil, fmt.Errorf("unrecognized cluster type %s", clusterType)
	}
}

// Since the main purpose of KIC is to set up Kong Gateway to properly route traffic to
// backends, ensure that discovering IP addresses of Pods works as expected, even in case
// of having multiple EndpointSlices per Service (by default it allows up to 100 endpoints
// per EndpointSlice, hence the below config decreases limit significantly).
const kindConfig = `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
kubeadmConfigPatches:
- |
  apiVersion: kubeadm.k8s.io/v1beta3
  kind: ClusterConfiguration
  controllerManager:
    extraArgs:
      max-endpoints-per-slice: "2"
`

func createKINDBuilder(t *testing.T) *environments.Builder {
	clusterBuilder := kind.NewBuilder().WithConfigReader(strings.NewReader(kindConfig))
	if v := testenv.ClusterVersion(); v != "" {
		clusterVersion := semver.MustParse(strings.TrimPrefix(v, "v"))
		clusterBuilder = clusterBuilder.WithClusterVersion(clusterVersion)
	}
	builder := environments.NewBuilder().WithClusterBuilder(clusterBuilder).WithAddons(metallb.New())
	if testenv.ClusterLoadImages() == "true" {
		builder = builder.WithAddons(buildImageLoadAddon(t, testenv.ControllerImageTag(), testenv.KongImageTag()))
	}
	return builder
}

func createExistingKINDBuilder(t *testing.T, name string) (*environments.Builder, error) {
	builder := environments.NewBuilder()
	cluster, err := kind.NewFromExisting(name)
	require.NoError(t, err)

	builder = builder.WithExistingCluster(cluster)
	builder = builder.WithAddons(metallb.New())
	if testenv.ClusterLoadImages() == "true" {
		builder = builder.WithAddons(buildImageLoadAddon(t, testenv.ControllerImageTag(), testenv.KongImageTag()))
	}
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

	if v := testenv.ClusterVersion(); v != "" {
		k8sVersion, err := semver.Parse(strings.TrimPrefix(v, "v"))
		if err != nil {
			return nil, err
		}

		t.Logf("creating GKE cluster, with requested version: %s", k8sVersion)
		clusterBuilder.WithClusterVersion(k8sVersion)
	}
	if ch := testenv.GKEClusterReleaseChannel(); ch != "" {
		clusterBuilder.WithReleaseChannel(gke.ReleaseChannel(ch))
	}

	return environments.NewBuilder().WithClusterBuilder(clusterBuilder), nil
}

type ManifestPatch func(io.Reader) (io.Reader, error)

type ManifestDeploy struct {
	// Path is the path to the manifest to deploy.
	Path string

	// SkipTestPatches is a flag that controls whether to apply standard test patches (e.g. replace controller
	// image when TEST_CONTROLLER_IMAGE set, etc.) to the manifests before deploying them.
	SkipTestPatches bool

	// AdditionalSecrets is a list of additional secrets to create before deploying the manifest.
	AdditionalSecrets []*corev1.Secret

	// Patches contain additionall patches that will be applied before deploying the manifest.
	Patches []ManifestPatch
}

func (d ManifestDeploy) Run(ctx context.Context, t *testing.T, env environments.Environment) Deployments {
	t.Helper()

	t.Log("waiting for testing environment to be ready")
	envReadyCtx, envReadyCancel := context.WithTimeout(ctx, testenv.EnvironmentReadyTimeout())
	defer envReadyCancel()
	require.NoError(t, <-env.WaitForReady(envReadyCtx))

	t.Log("creating the kong namespace")
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kong"}}
	_, err := env.Cluster().Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if !apierrors.IsAlreadyExists(err) {
		require.NoError(t, err)
	}

	t.Logf("deploying any supplemental secrets (found: %d)", len(d.AdditionalSecrets))
	for _, secret := range d.AdditionalSecrets {
		_, err := env.Cluster().Client().CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
		if !apierrors.IsAlreadyExists(err) {
			require.NoError(t, err)
		}
	}

	t.Logf("deploying %s manifest to the cluster", d.Path)
	manifest := getTestManifest(t, d.Path, d.SkipTestPatches, d.Patches...)
	kubeconfigFilename := getTemporaryKubeconfig(t, env)
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfigFilename, "apply", "-f", "-")
	cmd.Stdin = manifest
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	t.Log("waiting for controller to be ready")
	deployments := getManifestDeployments(d.Path)

	helpers.WaitForDeploymentRollout(ctx, t, env.Cluster(), deployments.ControllerNN.Namespace, deployments.ControllerNN.Name)
	helpers.WaitForDeploymentRollout(ctx, t, env.Cluster(), deployments.ProxyNN.Namespace, deployments.ProxyNN.Name)

	return deployments
}

func (d ManifestDeploy) Delete(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Helper()

	t.Log("waiting for testing environment to be ready")
	envReadyCtx, envReadyCancel := context.WithTimeout(ctx, testenv.EnvironmentReadyTimeout())
	defer envReadyCancel()
	require.NoError(t, <-env.WaitForReady(envReadyCtx))

	t.Logf("deleting any supplemental secrets (found: %d)", len(d.AdditionalSecrets))
	for _, secret := range d.AdditionalSecrets {
		err := env.Cluster().Client().CoreV1().Secrets(namespace).Delete(ctx, secret.Name, metav1.DeleteOptions{})
		require.NoError(t, err)
	}

	t.Logf("deleting %s manifest from the cluster", d.Path)
	manifest := getTestManifest(t, d.Path, d.SkipTestPatches)
	kubeconfigFilename := getTemporaryKubeconfig(t, env)
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfigFilename, "delete", "-f", "-")
	cmd.Stdin = manifest
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
}

// Deployments represent the deployments that are deployed by the all-in-one manifests.
type Deployments struct {
	ProxyNN      k8stypes.NamespacedName
	ControllerNN k8stypes.NamespacedName
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

// RestartController triggers the KIC pods' recreation by deleting them.
func (d Deployments) RestartController(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Helper()
	err := env.Cluster().Client().CoreV1().Pods(d.ControllerNN.Namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": "ingress-kong",
			},
		}),
	})
	require.NoError(t, err)
}

// getManifestDeployments returns the deployments for the proxy and controller that are expected to be deployed for a given
// manifest.
func getManifestDeployments(manifestPath string) Deployments {
	proxyDeploymentName := getProxyDeploymentName(manifestPath)
	return Deployments{
		ProxyNN: k8stypes.NamespacedName{
			Namespace: namespace,
			Name:      proxyDeploymentName,
		},
		ControllerNN: k8stypes.NamespacedName{
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
	if strings.Contains(manifestPath, "dbless") || strings.Contains(manifestPath, "multiple-gateways") {
		return multiPodDeploymentName
	}

	return singlePodDeploymentName
}

// For Kind clusters that have option --max-endpoints-per-slice=2, 3 gives
// one fully filled EndpointSlice and one filled in half.
const numberOfEchoBackends = 3

//nolint:unparam
func deployIngressWithEchoBackends(ctx context.Context, t *testing.T, env environments.Environment, noReplicas int) *netv1.Ingress {
	t.Helper()

	t.Log("deploying an HTTP service to test the ingress controller and proxy")
	container := generators.NewContainer("echo", test.EchoImage, test.EchoHTTPPort)
	container.Env = append(container.Env, corev1.EnvVar{
		Name: "POD_IP",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"},
		},
	})
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Spec.Replicas = lo.ToPtr(int32(noReplicas))
	deployment, err := env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	})

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
	for i := range service.Spec.Ports {
		// Set Service's appProtocol to http so that Kuma load-balances requests on HTTP-level instead of TCP.
		// TCP load-balancing proved to cause issues with not distributing load evenly in short time spans like in this test.
		// Ref: https://github.com/Kong/kubernetes-ingress-controller/issues/5498
		service.Spec.Ports[i].AppProtocol = lo.ToPtr("http")
	}

	_, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	})

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	ingress := generators.NewIngressForService(echoPath, map[string]string{
		annotations.AnnotationPrefix + annotations.StripPathKey: "true",
		annotations.AnnotationPrefix + annotations.MethodsKey:   http.MethodGet,
	}, service)
	ingress.Spec.IngressClassName = kong.String(ingressClass)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), corev1.NamespaceDefault, ingress))
	t.Cleanup(func() {
		assert.NoError(t, env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	})
	return ingress
}

func reconfigureExistingIngress(ctx context.Context, t *testing.T, env environments.Environment, ingress *netv1.Ingress, options ...func(*netv1.Ingress)) {
	for _, opt := range options {
		opt(ingress)
	}
	_, err := env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Update(ctx, ingress, metav1.UpdateOptions{})
	require.NoError(t, err)
}

//nolint:unparam
func verifyIngressWithEchoBackends(
	ctx context.Context,
	t *testing.T,
	env environments.Environment,
	noReplicas int,
) {
	t.Helper()
	verifyIngressWithEchoBackendsPath(ctx, t, env, noReplicas, echoPath)
}

func verifyIngressWithEchoBackendsPath(
	ctx context.Context,
	t *testing.T,
	env environments.Environment,
	noReplicas int,
	path string,
) {
	t.Helper()

	t.Log("finding the service URL (through Kong proxy service ip)")
	echoURL := fmt.Sprintf("http://%s%s", getKongProxyIP(ctx, t, env), path)

	t.Logf(
		"waiting for route from Ingress to be operational at %s and forward traffic to all %d backends",
		echoURL, numberOfEchoBackends,
	)
	uniqueResponses := make(map[string]struct{})
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		resp, err := helpers.DefaultHTTPClient().Get(echoURL)
		if !assert.NoError(c, err) {
			return
		}
		defer resp.Body.Close()

		if !assert.Equal(c, http.StatusOK, resp.StatusCode) {
			return
		}
		b := new(bytes.Buffer)
		_, err = b.ReadFrom(resp.Body)
		if !assert.NoError(c, err) {
			return
		}
		// Every backend responds with its own (different) IP address.
		if msg := b.String(); strings.Contains(msg, "With IP address") {
			uniqueResponses[msg] = struct{}{}
		}
		assert.Len(c, uniqueResponses, noReplicas)
	},
		ingressWait, 10*time.Millisecond,
	)

	t.Log("verifying the KongIngress method restriction")
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		fakeData := url.Values{}
		fakeData.Set("foo", "bar")
		resp, err := helpers.DefaultHTTPClient().PostForm(echoURL, fakeData)
		if !assert.NoError(c, err) {
			return
		}
		defer resp.Body.Close()
		assert.Equal(c, http.StatusNotFound, resp.StatusCode)
	},
		ingressWait, 10*time.Millisecond,
	)
}

// verifyIngressWithEchoBackendsInAdminAPI ensures all expected Kong Admin API resources
// are created for the Ingress deployed with deployIngressWithEchoBackends helper function.
func verifyIngressWithEchoBackendsInAdminAPI(
	ctx context.Context,
	t *testing.T,
	kongClient *kong.Client,
	noReplicas int,
) {
	t.Helper()

	require.Eventually(t, func() bool {
		start := time.Now()
		defer t.Logf("Try fetching config from %q started at %s, duration %v", kongClient.BaseRootURL(), start.Format(time.RFC3339), time.Since(start))

		services, err := kongClient.Services.ListAll(ctx)
		if err != nil {
			t.Logf("failed to list services: %v", err)
			return false
		}
		if len(services) != 1 || services[0].ID == nil {
			t.Log("still no service found...")
			return false
		}

		routes, _, err := kongClient.Routes.ListForService(ctx, services[0].ID, &kong.ListOpt{Size: 10})
		if err != nil {
			t.Logf("failed to list routes for service %s: %v", *services[0].ID, err)
			return false
		}
		if len(routes) != 1 {
			t.Log("still no route found...")
			return false
		}

		upstreams, err := kongClient.Upstreams.ListAll(ctx)
		if err != nil {
			t.Logf("failed to list upstreams: %v", err)
			return false
		}
		if len(upstreams) != 1 || upstreams[0].ID == nil {
			t.Logf("still no upstreams found...")
			return false
		}

		targets, err := kongClient.Targets.ListAll(ctx, upstreams[0].ID)
		if err != nil {
			t.Logf("failed to list targets for upstream %s: %v", *upstreams[0].ID, err)
			return false
		}
		if len(targets) != noReplicas {
			t.Log("still no targets found...")
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

// licenseOutput is the license section of the admin API root response.
type licenseOutput struct {
	License struct {
		Customer   string `json:"customer"`
		Dataplanes string `json:"dataplanes"`
		Creation   string `json:"license_creation_date"`
		Seats      string `json:"admin_seats"`
		Product    string `json:"product_subscription"`
		Plan       string `json:"support_plan"`
		Expiration string `json:"license_expiration_date"`
	} `json:"license"`
}

func getLicenseFromAdminAPI(ctx context.Context, env environments.Environment, adminPassword string) (licenseOutput, error) {
	var license licenseOutput
	// find the admin service
	service, err := env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, adminServiceName, metav1.GetOptions{})
	if err != nil {
		return license, fmt.Errorf("could not retrieve admin service: %w", err)
	}
	if len(service.Status.LoadBalancer.Ingress) == 0 {
		return license, fmt.Errorf("service %s has no external IPs", service.Name)
	}
	adminIP := service.Status.LoadBalancer.Ingress[0].IP

	// pull the root response
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/", adminIP), nil)
	if err != nil {
		return license, fmt.Errorf("could not create license request: %w", err)
	}
	req.Header.Set("Kong-Admin-Token", adminPassword)

	// read, unmarshal, and return result
	resp, err := helpers.DefaultHTTPClient().Do(req)
	if err != nil {
		return license, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return license, fmt.Errorf("failed to read response body: %w; body: %s", err, body)
	}
	if resp.StatusCode != http.StatusOK {
		return license, fmt.Errorf("unexpected status code: %d; body: %s", resp.StatusCode, body)
	}
	if err = json.Unmarshal(body, &license); err != nil {
		return licenseOutput{}, fmt.Errorf("could not unmarshal license response: %w; body: %s", err, body)
	}
	return license, nil
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

// buildImageLoadAddon creates addon to load KIC and kong images.
func buildImageLoadAddon(t *testing.T, images ...string) clusters.Addon {
	t.Helper()
	t.Log("building image load addon")

	if len(images) == 0 {
		return nil
	}

	builder := loadimage.NewBuilder()
	for _, image := range images {
		if image != "" {
			t.Logf("adding image %q to load image addon", image)
			// https://github.com/Kong/kubernetes-testing-framework/issues/440 this error only occurs if image == ""
			// it will eventually be removed from the WithImage return signature
			builder, _ = builder.WithImage(image)
		}
	}
	return builder.Build()
}

// createKongImagePullSecret creates the image pull secret
// `kong-enterprise-edition-docker` for kong enterprise image
// from env TEST_KONG_PULL_USERNAME and TEST_KONG_PULL_PASSWORD.
func createKongImagePullSecret(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Helper()

	if testenv.KongPullUsername() == "" || testenv.KongPullPassword() == "" {
		return
	}
	kubeconfigFilename := getTemporaryKubeconfig(t, env)

	const secretName = "kong-enterprise-edition-docker"
	cmd := exec.CommandContext(
		ctx,
		"kubectl", "--kubeconfig", kubeconfigFilename,
		"create", "secret", "docker-registry", secretName,
		"--docker-username="+testenv.KongPullUsername(),
		"--docker-password="+testenv.KongPullPassword(),
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

	existingClusterIsKind := strings.Split(testenv.ExistingClusterName(), ":")[0] == string(kind.KindClusterType)
	if existingClusterIsKind {
		return
	}

	clusterProviderIsKind := testenv.ClusterProvider() == string(kind.KindClusterType)
	if clusterProviderIsKind {
		return
	}

	clusterProviderUnspecified := testenv.ClusterProvider() == ""
	existingClusterUnspecified := testenv.ExistingClusterName() == ""
	if clusterProviderUnspecified && existingClusterUnspecified {
		return
	}

	t.Skip("test is supported only on Kind clusters")
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
func scaleDeployment(ctx context.Context, t *testing.T, env environments.Environment, deployment k8stypes.NamespacedName, replicas int32) {
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

func (d Deployments) Restart(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Helper()

	err := env.Cluster().Client().CoreV1().Pods(d.ControllerNN.Namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": d.ControllerNN.Name,
			},
		}),
	})
	require.NoError(t, err, "failed to delete controller pods")

	err = env.Cluster().Client().CoreV1().Pods(d.ControllerNN.Namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": d.ProxyNN.Name,
			},
		}),
	})
	require.NoError(t, err, "failed to delete proxy pods")

	helpers.WaitForDeploymentRollout(ctx, t, env.Cluster(), d.ControllerNN.Namespace, d.ControllerNN.Name)
	helpers.WaitForDeploymentRollout(ctx, t, env.Cluster(), d.ProxyNN.Namespace, d.ProxyNN.Name)
}
