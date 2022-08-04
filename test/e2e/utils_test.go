//go:build e2e_tests
// +build e2e_tests

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
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/sethvargo/go-password/password"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

var (
	// clusterVersionStr indicates the Kubernetes cluster version to use when
	// generating a testing environment and allows the caller to provide a specific
	// version. If no version is provided the default version for the cluster
	// provisioner in the testing framework will be used.
	clusterVersionStr = os.Getenv("KONG_CLUSTER_VERSION")

	// httpc is a standard HTTP client for tests to use that has a low default
	// timeout instead of the longer default provided by the http stdlib.
	httpc = http.Client{Timeout: time.Second * 10}
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
	if len(split) < 2 {
		t.Logf("could not parse override image '%v', using default manifest %v", imagetag, baseManifestPath)
		return os.Open(baseManifestPath)
	}
	modified, err := patchControllerImage(baseManifestPath, strings.Join(split[0:len(split)-1], ":"),
		split[len(split)-1])
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
// kong/kubernetes-ingress-controller image with the provided image. It returns the location of kustomize's output.
func patchControllerImage(baseManifestPath string, image string, tag string) (io.Reader, error) {
	workDir, err := os.MkdirTemp("", "kictest.")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workDir)
	orig, err := os.ReadFile(baseManifestPath)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(filepath.Join(workDir, "base.yaml"), orig, 0o600)
	if err != nil {
		return nil, err
	}
	kustomization := []byte(fmt.Sprintf(imageKustomizationContents, image, tag))
	err = os.WriteFile(filepath.Join(workDir, "kustomization.yaml"), kustomization, 0o600)
	if err != nil {
		return nil, err
	}
	kustomized, err := kustomizeManifest(workDir)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(kustomized), nil
}

// kustomizeManifest runs kustomize on a path and returns the YAML output.
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
func getKongProxyIP(ctx context.Context, t *testing.T, env environments.Environment, svc *corev1.Service) string {
	proxyIP := ""
	require.NotEqual(t, svc.Spec.Type, svc.Spec.ClusterIP)
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			proxyIP = svc.Status.LoadBalancer.Ingress[0].IP
			t.Logf("found loadbalancer IP for the Kong Proxy: %s", proxyIP)
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
	return proxyIP
}

// startPortForwarder runs "kubectl port-forward" in the background. It stops the forward when the provided context
// ends.
func startPortForwarder(ctx context.Context, t *testing.T, env environments.Environment, namespace, name, localPort,
	targetPort string,
) {
	kubeconfig, err := generators.NewKubeConfigForRestConfig(env.Name(), env.Cluster().Config())
	require.NoError(t, err)
	kubeconfigFile, err := os.CreateTemp(os.TempDir(), "portforward-tests-kubeconfig-")
	require.NoError(t, err)
	defer os.Remove(kubeconfigFile.Name())
	defer kubeconfigFile.Close()
	written, err := kubeconfigFile.Write(kubeconfig)
	require.NoError(t, err)
	require.Equal(t, len(kubeconfig), written)
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfigFile.Name(), "port-forward", "-n", namespace, name, fmt.Sprintf("%s:%s", localPort, targetPort)) //nolint:gosec
	t.Logf("forwarding port %s to %s/%s:%s", localPort, namespace, name, targetPort)
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
