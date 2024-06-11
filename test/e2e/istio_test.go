//go:build istio_tests

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/istio"
	kongaddon "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	testutils "github.com/kong/kubernetes-ingress-controller/v3/internal/util/test"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

var (
	// istioVersion allows the version of Istio to be overridden by ENV.
	// If not provided, the latest version of Istio will be tested.
	istioVersionStr = os.Getenv("ISTIO_VERSION")

	// kialiAPIPort is the port number that the Kiali API will use.
	kialiAPIPort = 20001

	// perHourRateLimit is a default rate-limit configuration for tests.
	//
	// See: https://docs.konghq.com/hub/kong-inc/rate-limiting/
	perHourRateLimit = 3
	// workloadEndpointIstioVersionCutoff is the lowest version that supports Kiali API /namespaces/<ns>/workloads/<workload>
	// that returns the metrics of a workload.
	workloadEndpointIstioVersionCutoff = semver.MustParse("1.18.0")
)

// TestIstioWithKongIngressGateway verifies integration of Kong Gateway as an Ingress
// Gateway for application traffic into an Istio mesh network.
//
// See: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/getting-started-istio/
// See: https://docs.konghq.com/kubernetes-ingress-controller/latest/references/version-compatibility/#istio
func TestIstioWithKongIngressGateway(t *testing.T) {
	t.Parallel()

	// Logger needs to be configured before anything else happens.
	// This is because the controller manager has a timeout for
	// logger initialization, and if the logger isn't configured
	// after 30s from the start of controller manager package init function,
	// the controller manager will set up a no op logger and continue.
	// The logger cannot be configured after that point.
	logger, logOutput, err := testutils.SetupLoggers("trace", "text")
	require.NoError(t, err, "failed to configure logger")
	if logOutput != "" {
		t.Logf("INFO: writing manager logs to %s", logOutput)
	}

	t.Log("configuring cluster addons for the testing environment")
	kongBuilder := kongaddon.NewBuilder().
		WithControllerDisabled().
		WithProxyAdminServiceTypeLoadBalancer().
		WithNamespace(consts.ControllerNamespace)
	kongAddon := kongBuilder.Build()

	t.Log("configuring istio cluster addon for the testing environment")
	istioBuilder := istio.NewBuilder().
		WithPrometheus().
		WithKiali()
	if istioVersionStr != "" {
		t.Logf("a specific version of istio was requested: %s", istioVersionStr)
		istioVersion, err := semver.ParseTolerant(istioVersionStr)
		require.NoError(t, err)
		istioBuilder.WithVersion(istioVersion)
	}
	istioAddon := istioBuilder.Build()

	ctx, env := setupE2ETest(t, istioAddon, kongAddon)

	t.Log("waiting for test cluster to be ready")
	envReadyCtx, envReadyCancel := context.WithTimeout(ctx, testenv.EnvironmentReadyTimeout())
	defer envReadyCancel()
	require.NoError(t, <-env.WaitForReady(envReadyCtx))

	t.Logf("istio version %s was deployed, enabling istio mesh network for the Kong Gateway's namespace", istioAddon.Version().String())
	require.NoError(t, istioAddon.EnableMeshForNamespace(ctx, env.Cluster(), kongAddon.Namespace()))

	t.Log("deleting kong pods to ensure istio sidecar injection")
	pods, err := env.Cluster().Client().CoreV1().Pods(kongAddon.Namespace()).List(ctx, metav1.ListOptions{})
	require.NoError(t, err)
	for _, pod := range pods.Items {
		require.NoError(t, env.Cluster().Client().CoreV1().Pods(kongAddon.Namespace()).Delete(ctx, pod.Name, metav1.DeleteOptions{}))
	}

	t.Log("waiting for kong pods to come back online with sidecar containers loaded")
	require.Eventually(t, func() bool {
		_, ready, err := kongAddon.Ready(ctx, env.Cluster())
		require.NoError(t, err)
		return ready
	}, time.Minute, time.Second)

	t.Log("Preparing the environment to run the controller manager")
	require.NoError(t, testutils.PrepareClusterForRunningControllerManager(ctx, env.Cluster()))
	t.Log("starting the controller manager")
	cancel, err := testutils.DeployControllerManagerForCluster(ctx, logger, env.Cluster(), kongAddon, []string{"--log-level=debug"})
	require.NoError(t, err)
	t.Cleanup(func() { cancel() })

	t.Log("creating a new mesh-enabled namespace for testing http traffic")
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "httpbin",
			Labels: map[string]string{
				"istio-injection": "enabled",
			},
		},
	}
	namespace, err = env.Cluster().Client().CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("creating a mesh enabled http deployment")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(namespace.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(namespace.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("creating an ingress resource for service %s with ingress.class %s", service.Name, ingressClass)
	ingress := generators.NewIngressForService("/httpbin", map[string]string{
		annotations.AnnotationPrefix + annotations.StripPathKey: "true",
	}, service)
	ingress.Spec.IngressClassName = kong.String(ingressClass)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), namespace.Name, ingress))

	t.Log("retrieving the kong proxy URL")
	proxyURL, err := kongAddon.ProxyHTTPURL(ctx, env.Cluster())
	require.NoError(t, err)

	t.Log("waiting for routes from Ingress to be operational")
	appURL := fmt.Sprintf("%s/httpbin", proxyURL)
	appStatusOKUrl := fmt.Sprintf("%s/status/200", appURL)
	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(appStatusOKUrl)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, time.Minute*5, time.Second)

	t.Log("exposing Kiali API via LoadBalancer type Service")
	kialiDeployment, err := env.Cluster().Client().AppsV1().Deployments(istioAddon.Namespace()).Get(ctx, "kiali", metav1.GetOptions{})
	require.NoError(t, err)
	service = generators.NewServiceForDeployment(kialiDeployment, corev1.ServiceTypeLoadBalancer)
	service.SetName("kiali-lb")
	service, err = env.Cluster().Client().CoreV1().Services(istioAddon.Namespace()).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		service, err = env.Cluster().Client().CoreV1().Services(istioAddon.Namespace()).Get(ctx, service.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return len(service.Status.LoadBalancer.Ingress) > 0
	}, time.Minute, time.Second)
	kialiAPIUrl := fmt.Sprintf("http://%s:%d/kiali/api", service.Status.LoadBalancer.Ingress[0].IP, kialiAPIPort)

	t.Logf("retrieving the Kiali workload metrics for deployment %s", deployment.Name)
	respData := kialiWorkloads{}
	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/namespaces/%s/apps/%s", kialiAPIUrl, namespace.Name, deployment.Name))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return false
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		return json.Unmarshal(b, &respData) == nil
	}, time.Minute*3, time.Second)

	t.Logf("verifying the contents of Kiali workload metrics for deployment %s", deployment.Name)
	require.Len(t, respData.Workloads, 1)
	require.Equal(t, deployment.Name, respData.Workloads[0].Name)
	require.True(t, respData.Workloads[0].IstioSidecar)
	workload := respData.Workloads[0]

	t.Logf("generating traffic and verifying health metrics for kiali workload %s", workload.Name)
	var health *workloadHealth
	var inboundHTTPRequests map[string]float64
	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(appStatusOKUrl)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return false
		}
		if health, err = getKialiWorkloadHealth(t, kialiAPIUrl, namespace.Name, workload.Name); err != nil {
			return false
		}
		inboundHTTPRequests = health.Requests.Inbound.HTTP
		return len(inboundHTTPRequests) == 1
	}, time.Minute*3, time.Second, "eventually metrics data for successful requests should start populating in kiali")
	require.Len(t, inboundHTTPRequests, 1, "no HTTP errors should have occurred")
	okResponsesFirstSample, ok := inboundHTTPRequests[strconv.Itoa(http.StatusOK)]
	require.True(t, ok, "a metric for 200 OK statuses should be present")
	require.Greater(t, okResponsesFirstSample, 0.0)

	t.Log("performing several http requests including 200's, 404's and 500 responses and verifying that health metrics get updated")
	notFoundURL := fmt.Sprintf("%s/status/404", appURL)
	serverErrorURL := fmt.Sprintf("%s/status/500", appURL)
	require.Eventually(t, func() bool {
		if err := verifyStatusForURL(appStatusOKUrl, http.StatusOK); err != nil {
			return false
		}
		if err := verifyStatusForURL(notFoundURL, http.StatusNotFound); err != nil {
			return false
		}
		if err := verifyStatusForURL(serverErrorURL, http.StatusInternalServerError); err != nil {
			return false
		}
		if health, err = getKialiWorkloadHealth(t, kialiAPIUrl, namespace.Name, workload.Name); err != nil {
			return false
		}
		inboundHTTPRequests = health.Requests.Inbound.HTTP
		return len(inboundHTTPRequests) == 3
	}, time.Minute*3, time.Second, "eventually metrics data for failed requests should start populating in kiali")

	t.Logf("verifying the contents of the health metrics for kiali workload %s after several failures", workload.Name)
	require.Len(t, inboundHTTPRequests, 3, "now expecting 200, 404 and 500's in the health metrics")
	okResponsesSecondSample, ok := inboundHTTPRequests[strconv.Itoa(http.StatusOK)]
	require.True(t, ok, "a metric for 200 OK statuses should be present")
	require.Greater(t, okResponsesSecondSample, 0.0)
	require.Greater(t, okResponsesSecondSample, okResponsesFirstSample)
	notFoundResponses, ok := inboundHTTPRequests[strconv.Itoa(http.StatusNotFound)]
	require.True(t, ok, "a metric for 404 Not Found statuses should be present")
	require.Greater(t, notFoundResponses, 0.0)
	serverErrorResponses, ok := inboundHTTPRequests[strconv.Itoa(http.StatusInternalServerError)]
	require.True(t, ok, "a metric for 404 Not Found statuses should be present")
	require.Greater(t, serverErrorResponses, 0.0)
	require.Greater(t, okResponsesSecondSample, serverErrorResponses)

	t.Logf("deploying a kong rate-limiter plugin for the %s deployment", deployment.Name)
	kongc, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	rateLimiterPlugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rate-limit",
			Namespace: namespace.Name,
		},
		PluginName: "rate-limiting",
		Config: apiextensionsv1.JSON{
			Raw: []byte(fmt.Sprintf(`{"hour":%d,"policy":"local"}`, perHourRateLimit)),
		},
	}
	rateLimiterPlugin, err = kongc.ConfigurationV1().KongPlugins(namespace.Name).Create(ctx, rateLimiterPlugin, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("enabling the rate-limiter plugin for deployment %s", deployment.Name)
	require.Eventually(t, func() bool {
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(namespace.Name).Get(ctx, "httpbin", metav1.GetOptions{})
		require.NoError(t, err)
		anns := ingress.ObjectMeta.GetAnnotations()
		anns["konghq.com/plugins"] = rateLimiterPlugin.Name
		ingress.ObjectMeta.SetAnnotations(anns)
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(namespace.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("waiting for the rate-limiter plugin to be active")
	var headers http.Header
	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(appStatusOKUrl)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		headers = resp.Header
		limitPerHour := headers.Get("X-Ratelimit-Limit-Hour")
		return limitPerHour != "" && (limitPerHour == strconv.Itoa(perHourRateLimit))
	}, time.Minute*3, time.Second)

	t.Log("intentionally using up the current rate-limit availability")
	require.Eventually(t, func() bool {
		return verifyStatusForURL(appStatusOKUrl, http.StatusTooManyRequests) == nil
	}, time.Minute*3, time.Second)
}

// -----------------------------------------------------------------------------
// Private Testing Functions - HTTP Request/Response Helpers
// -----------------------------------------------------------------------------

// verifyStatusForURL is a helper function which given a URL and a status code performs
// a GET and verifies the status code returning an error if the result isn't as expected.
func verifyStatusForURL(getURL string, statusCode int) error {
	resp, err := helpers.DefaultHTTPClient().Get(getURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != statusCode {
		return fmt.Errorf("expected status code %d got %d", statusCode, resp.StatusCode)
	}
	return nil
}

// getKialiWorkloadHealth produces the health metrics of a workload given the namespace and name of that workload.
func getKialiWorkloadHealth(t *testing.T, kialiAPIUrl string, namespace, workloadName string) (*workloadHealth, error) {
	istioVersion, err := semver.Parse(istioVersionStr)
	require.NoError(t, err, "failed to parse istio version")
	if istioVersion.GTE(workloadEndpointIstioVersionCutoff) {
		return getKialiWorkloadHealthIstioByWorkloadEndpoint(t, kialiAPIUrl, namespace, workloadName)
	}
	return getKialiWorkloadHealthByHealthEndpoint(t, kialiAPIUrl, namespace, workloadName)
}

// getKialiWorkloadHealthByHealthEndpoint gets health metrics of ALL workloads from /namespaces/<ns>/health?type=workload API.
// Used in istio 1.17 and prior. Istio 1.22 does not have this API.
func getKialiWorkloadHealthByHealthEndpoint(t *testing.T, kialiAPIUrl string, namespace, workloadName string) (*workloadHealth, error) {
	// generate the URL for the namespace health metrics
	kialiHealthURL := fmt.Sprintf("%s/namespaces/%s/health", kialiAPIUrl, namespace)
	req, err := http.NewRequest("GET", kialiHealthURL, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("type", "workload")
	req.URL.RawQuery = query.Encode()

	// make the health metrics request
	resp, err := helpers.DefaultHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// verify the health metrics response
	require.Equal(t, http.StatusOK, resp.StatusCode)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// decode the health metrics response
	healthData := map[string]workloadHealth{}
	if err := json.Unmarshal(b, &healthData); err != nil {
		return nil, err
	}

	// verify the presence of workload metrics data
	health, ok := healthData[workloadName]
	if !ok {
		return nil, fmt.Errorf("health metrics are not yet ready")
	}

	return &health, nil
}

// getKialiWorkloadHealthIstioByWorkloadEndpoint gets metrics of workload by /namespaces/<ns>/workloads/<workload> API.
// Used in Istio 1.18 and later. Istio 1.17 does not have this API.
func getKialiWorkloadHealthIstioByWorkloadEndpoint(t *testing.T, kialiAPIUrl string, namespace, workloadName string) (*workloadHealth, error) {
	kialiWorkloadURL := fmt.Sprintf("%s/namespaces/%s/workloads/%s", kialiAPIUrl, namespace, workloadName)
	resp, err := helpers.DefaultHTTPClient().Get(kialiWorkloadURL)
	require.NoErrorf(t, err, "failed to call Kiali workload API %s", kialiWorkloadURL)
	defer resp.Body.Close()
	// verify the workload response
	require.Equalf(t, http.StatusOK, resp.StatusCode, "Got code %d from Kiali workload API %s", resp.StatusCode, kialiWorkloadURL)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	status := workloadStatus{}
	err = json.Unmarshal(b, &status)
	if err != nil {
		return nil, err
	}
	return &status.Health, nil
}

// -----------------------------------------------------------------------------
// Private Testing Types - Kiali API Responses
// -----------------------------------------------------------------------------

type kialiWorkload struct {
	Name         string `json:"workloadName"`
	IstioSidecar bool   `json:"istioSidecar"`
}

type kialiWorkloads struct {
	Workloads []kialiWorkload `json:"workloads"`
}

type inbound struct {
	HTTP map[string]float64 `json:"http"`
}

type requests struct {
	Inbound inbound `json:"inbound"`
}

type workloadHealth struct {
	Requests requests `json:"requests"`
}

type workloadStatus struct {
	Name   string         `json:"name"`
	Health workloadHealth `json:"health"`
}
