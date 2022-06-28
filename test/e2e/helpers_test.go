//go:build e2e_tests
// +build e2e_tests

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

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

const (
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
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
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
	_, err = c.ConfigurationV1().KongIngresses(corev1.NamespaceDefault).Create(ctx, king, metav1.CreateOptions{})
	require.NoError(t, err)
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
	proxyIP := getKongProxyIP(ctx, t, env, svc)

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
		} else {
			return false
		}
		// verify the KongIngress method restriction
		fakeData := url.Values{}
		fakeData.Set("foo", "bar")
		resp, err = httpc.PostForm(fmt.Sprintf("http://%s/httpbin", proxyIP), fakeData)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusNotFound
	}, ingressWait, time.Second)
}

func deployGateway(ctx context.Context, t *testing.T, env environments.Environment) *gatewayv1alpha2.Gateway {
	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a supported gatewayclass to the test cluster")
	supportedGatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	supportedGatewayClass, err = gc.GatewayV1alpha2().GatewayClasses().Create(ctx, supportedGatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode")
	gw := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong",
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.GatewayUnmanagedAnnotation: "true", // trigger the unmanaged gateway mode
			},
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(supportedGatewayClass.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	gw, err = gc.GatewayV1alpha2().Gateways(corev1.NamespaceDefault).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)
	return gw
}

func verifyGateway(ctx context.Context, t *testing.T, env environments.Environment, gw *gatewayv1alpha2.Gateway) {
	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("verifying that the gateway receives a final ready condition once reconciliation completes")
	require.Eventually(t, func() bool {
		gw, err = gc.GatewayV1alpha2().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range gw.Status.Conditions {
			if cond.Reason == string(gatewayv1alpha2.GatewayReasonReady) {
				return true
			}
		}
		return false
	}, gatewayUpdateWaitTime, time.Second)
}

func deployHTTPRoute(ctx context.Context, t *testing.T, env environments.Environment, gw *gatewayv1alpha2.Gateway) {
	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	assert.NoError(t, err)
	t.Log("deploying an HTTP service to test the ingress controller and proxy")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("creating an HTTPRoute for service %s with Gateway %s", service.Name, gw.Name)
	pathMatchPrefix := gatewayv1alpha2.PathMatchPathPrefix
	path := "/httpbin"
	httpPort := gatewayv1alpha2.PortNumber(80)
	httproute := &gatewayv1alpha2.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gw.Name),
				}},
			},
			Rules: []gatewayv1alpha2.HTTPRouteRule{{
				Matches: []gatewayv1alpha2.HTTPRouteMatch{{
					Path: &gatewayv1alpha2.HTTPPathMatch{
						Type:  &pathMatchPrefix,
						Value: &path,
					},
				}},
				BackendRefs: []gatewayv1alpha2.HTTPBackendRef{{
					BackendRef: gatewayv1alpha2.BackendRef{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name: gatewayv1alpha2.ObjectName(service.Name),
							Port: &httpPort,
						},
					},
				}},
			}},
		},
	}
	_, err = gc.GatewayV1alpha2().HTTPRoutes(corev1.NamespaceDefault).Create(ctx, httproute, metav1.CreateOptions{})
	require.NoError(t, err)
}

// verifyHTTPRoute verifies an HTTPRoute exposes a route at /httpbin
// TODO this is not actually specific to HTTPRoutes. It is verifyIngress with the KongIngress removed
// Once we support HTTPMethod HTTPRouteMatch handling, we can combine the two into a single generic function
func verifyHTTPRoute(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Log("finding the kong proxy service ip")
	svc, err := env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, "kong-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	proxyIP := getKongProxyIP(ctx, t, env, svc)

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
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
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
	require.Contains(t, adminOutput.Version, "enterprise-edition")
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

// killKong kills the Kong container in a given Pod and returns when it has restarted
func killKong(ctx context.Context, t *testing.T, env environments.Environment, pod *corev1.Pod) {
	var orig, after int32
	for _, status := range pod.Status.ContainerStatuses {
		if status.Name == "proxy" {
			orig = status.RestartCount
		}
	}
	t.Logf("kong container has %v restart currently", orig)
	kubeconfig, err := generators.NewKubeConfigForRestConfig(env.Name(), env.Cluster().Config())
	require.NoError(t, err)
	kubeconfigFile, err := os.CreateTemp(os.TempDir(), "kill-tests-kubeconfig-")
	require.NoError(t, err)
	defer os.Remove(kubeconfigFile.Name())
	defer kubeconfigFile.Close()
	written, err := kubeconfigFile.Write(kubeconfig)
	require.NoError(t, err)
	require.Equal(t, len(kubeconfig), written)
	cmd := exec.Command("kubectl", "--kubeconfig", kubeconfigFile.Name(), "exec", "-n", pod.Namespace, pod.Name, "--", "kill", "1") //nolint:gosec
	require.NoError(t, cmd.Run())
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
