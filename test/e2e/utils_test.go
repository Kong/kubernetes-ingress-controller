//go:build e2e_tests || istio_tests
// +build e2e_tests istio_tests

package e2e

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/gke"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/phayes/freeport"
	"github.com/sethvargo/go-password/password"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/yaml"
)

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
// the original provided path.
func getTestManifest(t *testing.T, baseManifestPath string) (io.Reader, error) {
	var (
		manifestsReader io.Reader
		err             error
	)
	manifestsReader, err = os.Open(baseManifestPath)
	if err != nil {
		return nil, err
	}

	manifestsReader, err = patchControllerImageFromEnv(manifestsReader, baseManifestPath)
	if err != nil {
		t.Logf("failed patching controller image (%v), using default manifest %v", err, baseManifestPath)
		return manifestsReader, nil
	}

	var kongImageFullname string
	if kongImageLoad != "" {
		kongImageFullname = kongImageLoad
	} else {
		kongImageFullname = kongImageOverride
	}
	if kongImageFullname != "" {
		t.Logf("replace kong image to %s", kongImageFullname)
		split := strings.Split(kongImageFullname, ":")
		if len(split) < 2 {
			t.Logf("could not parse override image '%v', using default manifest %v", kongImageFullname, baseManifestPath)
			return manifestsReader, nil
		}
		repo := strings.Join(split[0:len(split)-1], ":")
		tag := split[len(split)-1]
		manifestsReader, err = patchKongImage(manifestsReader, repo, tag)
		if err != nil {
			t.Logf("failed patching override image '%v' (%v), using default manifest %v", kongImageFullname, err, baseManifestPath)
			return manifestsReader, nil
		}
	}

	manifestsReader, err = patchControllerStartTimeout(manifestsReader, 120, time.Second*3)
	if err != nil {
		t.Logf("failed patching controller timeouts (%v), using default manifest %v", err, baseManifestPath)
		return manifestsReader, nil
	}

	manifestsReader, err = patchLivenessProbes(manifestsReader, 0, 10, time.Second*15, time.Second*3)
	if err != nil {
		t.Logf("failed patching kong liveness (%v), using default manifest %v", err, baseManifestPath)
		return manifestsReader, nil
	}

	manifestsReader, err = patchLivenessProbes(manifestsReader, 1, 15, time.Second*3, time.Second*10)
	if err != nil {
		t.Logf("failed patching controller liveness (%v), using default manifest %v", err, baseManifestPath)
		return manifestsReader, nil
	}

	t.Logf("generated modified manifest at %v", baseManifestPath)
	return manifestsReader, nil
}

// patchControllerImageFromEnv will optionally replace a default controller image in manifests with one of `imageLoad` or `imageOverride` if any is set.
func patchControllerImageFromEnv(manifestReader io.Reader, baseManifestPath string) (io.Reader, error) {
	var imageFullname string
	if imageLoad != "" {
		imageFullname = imageLoad
	} else {
		imageFullname = imageOverride
	}

	if imageFullname != "" {
		split := strings.Split(imageFullname, ":")
		if len(split) < 2 {
			return manifestReader, fmt.Errorf("could not parse override image '%v', using default manifest %v", imageFullname, baseManifestPath)
		}
		repo := strings.Join(split[0:len(split)-1], ":")
		tag := split[len(split)-1]
		var err error
		manifestReader, err = patchControllerImage(manifestReader, repo, tag)
		if err != nil {
			return manifestReader, fmt.Errorf("failed patching override image '%v' (%w), using default manifest %v", imageFullname, err, baseManifestPath)
		}
	}
	return manifestReader, nil
}

func getCurrentGitTag(path string) (semver.Version, error) {
	cmd := exec.Command("git", "describe", "--tags")
	cmd.Dir = path
	tagBytes, err := cmd.Output()
	if err != nil {
		return semver.Version{}, fmt.Errorf("%q command failed: %w", cmd.String(), err)
	}
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
	tagsBytes, err := cmd.Output()
	if err != nil {
		return semver.Version{}, err
	}
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

// getKongProxyIP takes a Service with Kong proxy ports and returns and its IP, or fails the test if it cannot.
func getKongProxyIP(ctx context.Context, t *testing.T, env environments.Environment) string {
	refreshService := func() *corev1.Service {
		svc, err := env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, "kong-proxy", metav1.GetOptions{})
		require.NoError(t, err)
		return svc
	}

	svc := refreshService()
	require.NotEqual(t, svc.Spec.Type, corev1.ServiceTypeClusterIP, "ClusterIP service is not supported")

	//nolint: exhaustive
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
	localPort, err := freeport.GetFreePort()
	require.NoError(t, err)

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
		t.Logf("failed to create request: %v", err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Logf("failed to get response: %v", err)
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

// stripCRDs removes every CustomResourceDefinition from the manifest.
func stripCRDs(t *testing.T, manifest io.Reader) io.Reader {
	const sep = "---\n"

	in, err := io.ReadAll(manifest)
	require.NoError(t, err)

	var filteredObjs [][]byte
	for _, objYaml := range bytes.Split(in, []byte(sep)) {
		var obj struct {
			Kind string `yaml:"kind"`
		}
		err = yaml.Unmarshal(objYaml, &obj)
		require.NoError(t, err)

		if obj.Kind == "CustomResourceDefinition" {
			continue
		}

		filteredObjs = append(filteredObjs, objYaml)
	}

	outBytes := bytes.Join(filteredObjs, []byte(sep))
	return bytes.NewReader(outBytes)
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

// isPodReady evaluates whether a pod is in Ready state.
func isPodReady(pod corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}
