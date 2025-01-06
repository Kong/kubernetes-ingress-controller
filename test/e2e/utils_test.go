//go:build e2e_tests || istio_tests || performance_tests

package e2e

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/gke"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/sethvargo/go-password/password"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

const (
	// adminPasswordSecretName is the name of the secret which will house the admin
	// API admin password.
	adminPasswordSecretName = "kong-enterprise-superuser-password"

	dblessPath = "manifests/all-in-one-dbless.yaml"
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
func exposeAdminAPI(ctx context.Context, t *testing.T, env environments.Environment, proxyDeployment k8stypes.NamespacedName) {
	t.Log("updating the proxy container KONG_ADMIN_LISTEN to expose the admin api")
	deployment, err := env.Cluster().Client().AppsV1().Deployments(proxyDeployment.Namespace).Get(ctx, proxyDeployment.Name, metav1.GetOptions{})
	require.NoError(t, err)
	for i, containerSpec := range deployment.Spec.Template.Spec.Containers {
		if containerSpec.Name == proxyContainerName {
			for j, envVar := range containerSpec.Env {
				if envVar.Name == "KONG_ADMIN_LISTEN" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = "0.0.0.0:8001, 0.0.0.0:8444 ssl"
				}
			}
		}
	}

	deployment, err = env.Cluster().Client().AppsV1().Deployments(proxyDeployment.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
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
}

// getTestManifest gets a manifest io.Reader, applying optional patches to the base manifest provided.
// In case of any failure while patching, the base manifest is returned.
// If skipTestPatches is true, no patches are applied (useful when untouched manifest is needed, e.g. in upgrade tests).
func getTestManifest(t *testing.T, baseManifestPath string, skipTestPatches bool, testPatches ...ManifestPatch) io.Reader {
	t.Helper()
	t.Logf("getting test manifest from %v", baseManifestPath)

	var (
		manifestsReader io.Reader
		err             error
	)
	manifestsReader, err = os.Open(baseManifestPath)
	require.NoError(t, err)

	if !skipTestPatches {
		t.Logf("applying test patches to manifest %v", baseManifestPath)

		manifestsReader, err = patchControllerImageFromEnv(t, manifestsReader)
		if err != nil {
			t.Logf("failed patching controller image (%v), using default manifest %v", err, baseManifestPath)
			return manifestsReader
		}

		manifestsReader, err = patchGatewayImageFromEnv(t, manifestsReader)
		if err != nil {
			t.Logf("failed patching gateway image (%v), using default manifest %v", err, baseManifestPath)
			return manifestsReader
		}

		manifestsReader, err = patchControllerStartTimeout(manifestsReader, 120, time.Second*3)
		if err != nil {
			t.Logf("failed patching controller timeouts (%v), using default manifest %v", err, baseManifestPath)
			return manifestsReader
		}

		deployments := getManifestDeployments(baseManifestPath)
		manifestsReader, err = patchLivenessProbes(manifestsReader, deployments.ProxyNN, 10, time.Second*15, time.Second*3)
		if err != nil {
			t.Logf("failed patching kong liveness (%v), using default manifest %v", err, baseManifestPath)
			return manifestsReader
		}

		manifestsReader, err = patchLivenessProbes(manifestsReader, deployments.ControllerNN, 15, time.Second*3, time.Second*10)
		if err != nil {
			t.Logf("failed patching controller liveness (%v), using default manifest %v", err, baseManifestPath)
			return manifestsReader
		}

	}

	for _, patch := range testPatches {
		manifestsReader, err = patch(manifestsReader)
		if err != nil {
			t.Logf("failed patching manifest, using default manifest %v, err: %v", baseManifestPath, err)
		}
	}

	return manifestsReader
}

// extractVersionFromImage extracts semver of image from image tag. If tag is not given,
// or is not in a semver format, it returns an error.
// for example: kong/kubernetes-ingress-controller:2.9.3 => semver.Version{Major:2,Minor:9,Patch:3}.
//
//lint:ignore U1000 retained for future use
func extractVersionFromImage(imageName string) (semver.Version, error) {
	split := strings.Split(imageName, ":")
	if len(split) < 2 {
		return semver.Version{}, fmt.Errorf("could not parse override image '%s', expected <repo>:<tag> format", imageName)
	}
	// parse version from image tag, like kong/kubernetes-ingress-controller:2.9.3 => 2.9.3
	tag := split[len(split)-1]
	v, err := semver.ParseTolerant(tag)
	if err != nil {
		return semver.Version{}, fmt.Errorf("failed to parse version from image tag %s: %w", tag, err)
	}
	return v, nil
}

// skipTestIfControllerVersionBelow skips the test case if version of override KIC image is
// below the minVersion.
// if the override KIC image is not set, it assumes that the latest image is used, so it never skips
// the test if override image is not given.
//
//lint:ignore U1000 retained for future use
func skipTestIfControllerVersionBelow(t *testing.T, minVersion semver.Version) {
	if testenv.ControllerImageTag() == "" {
		return
	}
	v, err := extractVersionFromImage(testenv.ControllerImageTag())
	// assume using latest version if failed to extract version from image tag.
	if err != nil {
		t.Logf("could not extract version from controller image: %v, assume using the latest version", err)
		return
	}
	if v.LE(minVersion) {
		t.Skipf("skipped the test because version of KIC %s is below the minimum version %s",
			v.String(), minVersion.String())
	}
}

// patchGatewayImageFromEnv will optionally replace a default controller image in manifests with env overrides.
func patchGatewayImageFromEnv(t *testing.T, manifestsReader io.Reader) (io.Reader, error) {
	t.Helper()

	if testenv.KongImageTag() != "" {
		t.Logf("replace kong image with %s", testenv.KongImageTag())
		manifestsReader, err := patchKongImage(manifestsReader, testenv.KongImage(), testenv.KongTag())
		if err != nil {
			return nil, fmt.Errorf("failed patching override image '%v'", testenv.KongImageTag())
		}
		return manifestsReader, nil
	}

	t.Log("kong image override undefined, using defaults")
	return manifestsReader, nil
}

// patchControllerImageFromEnv will optionally replace a default controller image in manifests with env override
// if it's set.
func patchControllerImageFromEnv(t *testing.T, manifestReader io.Reader) (io.Reader, error) {
	t.Helper()

	if testenv.ControllerImageTag() != "" {
		manifestReader, err := patchControllerImage(manifestReader, testenv.ControllerImage(), testenv.ControllerTag())
		if err != nil {
			return nil, fmt.Errorf("failed patching override image '%v': %w", testenv.ControllerImageTag(), err)
		}
		return manifestReader, nil
	}

	t.Log("controller image override undefined, using defaults")
	return manifestReader, nil
}

// getKongVersionFromOverrideTag parses Kong version from env effective version or override tag. The effective version
// takes precedence.
//
//lint:ignore U1000 retained for future use
func getKongVersionFromOverrideTag() (kong.Version, error) {
	if kongEffectiveVersion := testenv.KongEffectiveVersion(); kongEffectiveVersion != "" {
		return kong.ParseSemanticVersion(kongEffectiveVersion)
	}

	if testenv.KongImageTag() == "" {
		return kong.Version{}, errors.New("No Kong tag provided")
	}

	return kong.ParseSemanticVersion(testenv.KongTag())
}

// getKongProxyIP takes a Service with Kong proxy ports and returns and its IP, or fails the test if it cannot.
func getKongProxyIP(ctx context.Context, t *testing.T, env environments.Environment) string {
	t.Helper()

	refreshService := func() *corev1.Service {
		svc, err := env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, "kong-proxy", metav1.GetOptions{})
		require.NoError(t, err)
		return svc
	}

	svc := refreshService()
	require.NotEqual(t, svc.Spec.Type, corev1.ServiceTypeClusterIP, "ClusterIP service is not supported")

	switch svc.Spec.Type {
	case corev1.ServiceTypeLoadBalancer:
		return getKongProxyLoadBalancerIP(t, refreshService)
	case corev1.ServiceTypeNodePort:
		return getKongProxyNodePortIP(ctx, t, env, svc)
	default:
		t.Fatalf("unknown service type: %q", svc.Spec.Type)
		return ""
	}
}

func getKongProxyLoadBalancerIP(t *testing.T, refreshSvc func() *corev1.Service) string {
	t.Helper()

	var resIP string
	require.Eventually(t, func() bool {
		svc := refreshSvc()

		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			ip := svc.Status.LoadBalancer.Ingress[0].IP
			t.Logf("found loadbalancer IP for the Kong Proxy: %s", ip)
			resIP = ip
			return true
		}

		t.Log("no IP for LoadBalancer found yet")
		return false
	}, ingressWait, time.Second)

	return resIP
}

func getKongProxyNodePortIP(ctx context.Context, t *testing.T, env environments.Environment, svc *corev1.Service) string {
	t.Helper()

	var port corev1.ServicePort
	for _, sport := range svc.Spec.Ports {
		if sport.Name == "kong-proxy" || sport.Name == "proxy" {
			port = sport
		}
	}

	// GKE clusters by default do not allow ingress traffic to its nodes
	// TODO: consider adding an option to create firewall rules in KTF GKE provider
	if env.Cluster().Type() == gke.GKEClusterType {
		kongProxyLocalPort := startPortForwarder(ctx, t, env, svc.Namespace, fmt.Sprintf("service/%s", svc.Name), strconv.Itoa(int(port.Port)))
		return fmt.Sprintf("localhost:%d", kongProxyLocalPort)
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
				intAddrs = append(intAddrs, naddr.Address)
			}
		}
	}
	// local clusters (KIND, minikube) typically provide no external addresses, but their internal addresses are
	// routeable from their host. We prefer external addresses if they're available, but fall back to internal
	// in their absence
	if len(extAddrs) > 0 {
		t.Logf("picking an external NodePort address: %s", extAddrs[0])
		return fmt.Sprintf("%v:%v", extAddrs[0], port.NodePort)
	} else if len(intAddrs) > 0 {
		t.Logf("picking an internal NodePort address: %s", intAddrs[0])
		return fmt.Sprintf("%v:%v", intAddrs[0], port.NodePort)
	}

	assert.Fail(t, "both extAddrs and intAddrs are empty")
	return ""
}

// startPortForwarder runs "kubectl port-forward" in the background. It returns a local port that the traffic gets forward to.
// It stops the forward when the provided context ends.
func startPortForwarder(ctx context.Context, t *testing.T, env environments.Environment, namespace, name, targetPort string) int {
	t.Helper()

	localPort := helpers.GetFreePort(t)

	kubeconfig := getTemporaryKubeconfig(t, env)
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig, "port-forward", "-n", namespace, name, fmt.Sprintf("%d:%s", localPort, targetPort))
	out := new(bytes.Buffer)
	cmd.Stderr = out
	cmd.Stdout = out

	t.Logf("forwarding port %d to %s/%s:%s", localPort, namespace, name, targetPort)
	if startErr := cmd.Start(); startErr != nil {
		require.NoError(t, startErr, out.String())
	}
	require.Eventually(t, func() bool {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", localPort))
		if err == nil {
			conn.Close()
			return true
		}

		t.Logf("port forwarding (via %q) not ready....", cmd.String())
		return false
	}, kongComponentWait, time.Second)

	return localPort
}

// httpGetResponseContains returns true if the response body of GETting the URL contains specified substring.
func httpGetResponseContains(t *testing.T, url string, client *http.Client, substring string) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Logf("Failed to create request: %v", err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Logf("Failed to get response: %v", err)
		return false
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Logf("failed to read response body: %v", err)
		return false
	}

	return strings.Contains(string(body), substring)
}

// getPodLogs gets logs af ALL containers inside pod.
// returns a non-nil error if we failed to get logs of the pod (for example, pod is not started yet)
// otherwise, it returns the combination of logs of all containers in the pod.
// if we failed to create a kubeconfig file, fail the test `t` immediately.
func getPodLogs(
	ctx context.Context, t *testing.T, env environments.Environment,
	namespace string, podName string,
) (string, error) {
	kubeconfig := getTemporaryKubeconfig(t, env)
	stderr := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig, "logs", podName, "-n", namespace, "--all-containers")
	cmd.Stderr = stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s", stderr.String())
	}
	return string(out), nil
}

// containerDidntCrash evaluates whether a container with a given containerName did not restart.
// In case name=containerName is not found in pod's containers, returns false.
func containerDidntCrash(pod corev1.Pod, containerName string) bool {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Name == containerName {
			return containerStatus.RestartCount == 0
		}
	}
	return false
}

// ensureNoneOfDeploymentPodsHasCrashed ensures that none of the pods of a deployment has crashed.
func ensureNoneOfDeploymentPodsHasCrashed(ctx context.Context, t *testing.T, env environments.Environment, deploymentNN k8stypes.NamespacedName) {
	t.Logf("ensuring none of %s deployment pods has crashed", deploymentNN.String())
	pods, err := listPodsByLabels(ctx, env, deploymentNN.Namespace, map[string]string{"app": deploymentNN.Name})
	require.NoError(t, err)
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			require.Truef(t, containerDidntCrash(pod, container.Name), "controller pod %s/%s crashed", pod.Namespace, pod.Name)
		}
	}
}

type setEnvParams struct {
	kubeCfgPath   string
	namespace     string
	target        string // e.g. deployment/name or pod/name
	containerName string
	variableName  string
	value         string
}

func setEnv(p setEnvParams) error {
	var envvar string
	if p.value == "" {
		envvar = fmt.Sprintf("%s-", p.variableName)
	} else {
		envvar = fmt.Sprintf("%s=%s", p.variableName, p.value)
	}
	//nolint:gosec
	cmd := exec.Command(
		"kubectl",
		"--kubeconfig", p.kubeCfgPath,
		"set", "env",
		"-n", p.namespace,
		"-c", p.containerName,
		p.target,
		envvar,
	)
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("updating envvar failed: STDOUT(%s) STDERR(%s): %w", stdout, stderr, err)
	}
	return nil
}

// dumpToTempFile dumps the contents of the reader to a temporary file and returns the path to the file.
func dumpToTempFile(t *testing.T, reader io.Reader) string {
	t.Helper()
	file, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)
	defer file.Close()
	_, err = io.Copy(file, reader)
	require.NoError(t, err)
	return file.Name()
}
