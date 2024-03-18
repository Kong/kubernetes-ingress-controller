//go:build performance_tests

package e2e

import (
	"bytes"
	"context"
	"encoding/base64"
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
	// resourcesNumberString is the number of resource rules to be created
	// if not set, defaultResNum will be used
	resourcesNumberString = os.Getenv("PERF_RES_NUMBER")

	consumerUsername = "consumer-key-auth-name-%d"

	// rulesTpl is the template of resource rules to be created
	// %d is the index of the resource rule
	// In order to simplify the testing process,
	// individual secrets are not specifically assigned for each consumer here.
	// Instead, a fixed secret `consumer-key-auth-secret` is used.
	rulesTpl = `
---
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: auth-plugin-%d
  annotations:
    kubernetes.io/ingress.class: kong
plugin: key-auth

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ingress-%d
  annotations:
    konghq.com/plugins: auth-plugin-%d
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

---
apiVersion: v1
data:
  key: %s
kind: Secret
metadata:
  labels:
    konghq.com/credential: key-auth
  name: consumer-key-auth-secret-%d
type: Opaque

---
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: consumer-%d
  annotations:
    kubernetes.io/ingress.class: kong
username: %s
credentials:
- consumer-key-auth-secret-%d

`

	updatedIngressTpl = `
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ingress-%d
  annotations:
    konghq.com/plugins: auth-plugin-%d
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

	// Use the following file to store the time-consuming reports for each action.
	allResourceApplyReport      = "all_resource_apply_%d.txt"
	allResourceTakeEffectReport = "all_resource_take_effect_%d.txt"
	oneResourceUpdateReport     = "one_resource_update_%d.txt"
	oneResourceTakeEffectReport = "one_resource_take_effect_%d.txt"
)

// -----------------------------------------------------------------------------
// E2E Performance tests
// -----------------------------------------------------------------------------

// TestResourceApplyAndUpdatePerf tests the performance of KIC,
// when creat a specified number of resources, how long does it take for them to take effect in the Gateway.
// And when update one of the resources, how long does it take for the update to take effect in the Gateway.
func TestResourceApplyAndUpdatePerf(t *testing.T) {
	t.Log("configuring all-in-one-dbless.yaml manifest test")
	t.Parallel()
	ctx, env := setupE2ETest(t)

	if resourcesNumberString != "" {
		if num, err := strconv.Atoi(resourcesNumberString); err == nil {
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

	batchSize := 500
	startTime := time.Now()
	for i := 0; i < defaultResNum; i += batchSize {
		resourceYaml := ""
		for j := i; j < i+batchSize && j < defaultResNum; j++ {
			resourceYaml += fmt.Sprintf(rulesTpl, j, j, j, j, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(consumerUsername, j))), j, j, fmt.Sprintf(consumerUsername, j), j)

		}
		err = applyResourceWithKubectl(ctx, t, kubeconfig, resourceYaml)
		require.NoError(t, err)
	}
	completionTime := time.Now()
	t.Logf("time to apply %d rules(including Ingress, plugin, and consumer): %v", defaultResNum, completionTime.Sub(startTime))
	writeResultToTempFile(t, allResourceApplyReport, defaultResNum, int(completionTime.Sub(startTime).Milliseconds()))

	t.Log("getting kong proxy IP after LB provisioning")
	proxyURLForDefaultIngress := "http://" + getKongProxyIP(ctx, t, env)

	t.Log("waiting for routes from Ingress to be operational")

	// create wait group to wait for all resource rules to take effect
	randomList := getRandomList(defaultResNum)
	var wg sync.WaitGroup
	wg.Add(len(randomList))

	for _, i := range randomList {
		go func(i int) {
			defer wg.Done()

			require.Eventually(t, func() bool {
				return isRouteActive(ctx, t, helpers.DefaultHTTPClient(), proxyURLForDefaultIngress, "get", i)
			}, ingressWait*10, time.Millisecond*500)
		}(i)
	}

	wg.Wait()

	effectTime := time.Now()

	t.Logf("time to make %d ingress rules take effect: %v", defaultResNum, effectTime.Sub(completionTime))
	writeResultToTempFile(t, allResourceTakeEffectReport, defaultResNum, int(effectTime.Sub(completionTime).Milliseconds()))

	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(10000)

	startTime = time.Now()
	err = applyResourceWithKubectl(ctx, t, kubeconfig, fmt.Sprintf(updatedIngressTpl, randomInt, randomInt, randomInt))
	require.NoError(t, err)
	completionTime = time.Now()

	require.Eventually(t, func() bool {
		return isRouteActive(ctx, t, helpers.DefaultHTTPClient(), proxyURLForDefaultIngress, "ip", randomInt)
	}, ingressWait*5, time.Millisecond*500)

	effectTime = time.Now()

	t.Logf("time to update 1 ingress rules when %d ingress exists: %v", defaultResNum, completionTime.Sub(startTime))
	t.Logf("time to make 1 ingress rules take effect when %d ingress exists: %v", defaultResNum, effectTime.Sub(completionTime))
	writeResultToTempFile(t, oneResourceUpdateReport, defaultResNum, int(completionTime.Sub(startTime).Milliseconds()))
	writeResultToTempFile(t, oneResourceTakeEffectReport, defaultResNum, int(effectTime.Sub(completionTime).Milliseconds()))

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

func isRouteActive(ctx context.Context, t *testing.T, client *http.Client, proxyIP, path string, index int) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", proxyIP, path), nil)
	require.NoError(t, err)
	req.Host = fmt.Sprintf("example-%d.com", index)
	req.Header.Set("apikey", fmt.Sprintf(consumerUsername, index))

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
	err := cmd.Run()
	return err
}

// we will store the result in the `/tmp/kic-perf/` directory
func writeResultToTempFile(t *testing.T, filename string, resourceNum, time int) {
	// create a file to /tmp/kic-perf/all_apply.txt
	// if the file already exists, it will be overwritten
	// if the file does not exist, it will be created
	defaultResultDir := "/tmp/kic-perf/"
	if err := os.MkdirAll(defaultResultDir, 0755); err != nil {
		t.Logf("failed to create directory: %v", err)
		return
	}

	file, err := os.Create(defaultResultDir + fmt.Sprintf(filename, resourceNum))
	if err != nil {
		t.Logf("failed to create file: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d %d\n", resourceNum, time))
	if err != nil {
		t.Logf("failed to write to file: %v", err)
		return
	}
}
