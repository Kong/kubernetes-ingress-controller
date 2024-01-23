//go:build performance_tests

package e2e

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

var (
	defaultResNum = 10000
	resNumStr     = os.Getenv("PERF_RES_NUMBER")
)

// -----------------------------------------------------------------------------
// E2E Performance tests
// -----------------------------------------------------------------------------

// TestBasicHTTPRoute will create a basic HTTP route and test its functionality
// against a Kong proxy. This test will be used to measure the performance of
// the KIC with OpenTelemetry.
func TestBasicPerf(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, env := setupE2ETest(t)

	if resNumStr != "" {
		if num, err := strconv.Atoi(resNumStr); err == nil {
			defaultResNum = num
		}
	}

	t.Log("deploying kong components")
	ManifestDeploy{Path: dblessPath}.Run(ctx, t, env)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services("default").Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	kubeconfig := getTemporaryKubeconfig(t, env)

	ingressTpl := `
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ingress-%d
spec:
  ingressClassName: kong
  rules:
  - host: example-%d.com
    http:
      paths:
      - backend:
          service:
            name: httpbin
            port:
              number: 80
        path: /get
        pathType: Exact

`

	ingressYaml := ""
	for i := 0; i < defaultResNum; i++ {
		ingressYaml += fmt.Sprintf(ingressTpl, i, i)
	}

	startTime := time.Now()
	err = applyResourceWithKubectl(ctx, t, kubeconfig, ingressYaml)
	require.NoError(t, err)
	completionTime := time.Now()

	t.Log("getting kong proxy IP after LB provisioning")
	proxyURLForDefaultIngress := "http://" + getKongProxyIP(ctx, t, env)

	t.Log("waiting for routes from Ingress to be operational")

	// create wait group to wait for all ingress rules to take effect
	randomList := getRandomList(defaultResNum)
	var wg sync.WaitGroup
	wg.Add(len(randomList))

	for _, i := range randomList {
		go func(i int) {
			defer wg.Done()

			require.Eventually(t, func() bool {
				return isRouteActive(ctx, t, helpers.DefaultHTTPClient(), proxyURLForDefaultIngress, "get", fmt.Sprintf("example-%d.com", i))
			}, ingressWait, time.Millisecond*500)
		}(i)
	}

	wg.Wait()

	effectTime := time.Now()

	t.Logf("time to apply %d ingress rules: %v", defaultResNum, completionTime.Sub(startTime))
	t.Logf("time to make %d ingress rules take effect: %v", defaultResNum, effectTime.Sub(completionTime))

	updatedIngressTpl := `
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ingress-%d
spec:
  ingressClassName: kong
  rules:
  - host: example-%d.com
    http:
      paths:
      - backend:
          service:
            name: httpbin
            port:
              number: 80
        path: /ip
        pathType: Exact

`
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(10000)

	startTime = time.Now()
	err = applyResourceWithKubectl(ctx, t, kubeconfig, fmt.Sprintf(updatedIngressTpl, randomInt, randomInt))
	require.NoError(t, err)
	completionTime = time.Now()

	require.Eventually(t, func() bool {
		return isRouteActive(ctx, t, helpers.DefaultHTTPClient(), proxyURLForDefaultIngress, "ip", fmt.Sprintf("example-%d.com", randomInt))
	}, ingressWait, time.Millisecond*500)

	effectTime = time.Now()

	t.Logf("time to apply 1 ingress rules when %d ingress exists: %v", defaultResNum, completionTime.Sub(startTime))
	t.Logf("time to make 1 ingress rules take effect when %d ingress exists: %v", defaultResNum, effectTime.Sub(completionTime))

}

func getRandomList(n int) []int {
	if n <= 10 {
		return []int{0, n}
	}
	randPerm := rand.Perm(n)
	randPerm = randPerm[:10]
	randPerm = append(randPerm, 0, n-1)

	return randPerm
}

func isRouteActive(ctx context.Context, t *testing.T, client *http.Client, proxyIP, path, hostname string) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", proxyIP, path), nil)
	require.NoError(t, err)
	req.Host = hostname

	resp, err := client.Do(req)
	if err != nil {
		t.Logf("WARNING: error while waiting for %s: %v", proxyIP, err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		// now that the ingress backend is routable
		b := new(bytes.Buffer)
		n, err := b.ReadFrom(resp.Body)
		require.NoError(t, err)
		require.True(t, n > 0)
		return strings.Contains(b.String(), "origin")
	}
	return false
}

func applyResourceWithKubectl(ctx context.Context, t *testing.T, kubeconfig, resourceYAML string) error {
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig, "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(resourceYAML)
	_, err := cmd.CombinedOutput()
	return err
}
