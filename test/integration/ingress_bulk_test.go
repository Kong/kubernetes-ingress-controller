//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/internal/annotations"
)

const testBulkIngressNamespace = "ingress-bulk-testing"

// TestIngressBulk attempts to validate functionality at scale by rapidly deploying a large number of ingress resources.
func TestIngressBulk(t *testing.T) {
	ctx := context.Background()

	t.Logf("creating namespace %s for testing", testBulkIngressNamespace)
	namespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testBulkIngressNamespace}}
	namespace, err := env.Cluster().Client().CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up namespace %s", testBulkIngressNamespace)
		require.NoError(t, env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, namespace.Name, metav1.DeleteOptions{}))
		require.Eventually(t, func() bool {
			_, err := env.Cluster().Client().CoreV1().Namespaces().Get(ctx, namespace.Name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					return true
				}
			}
			return false
		}, ingressWait, waitTick)
	}()

	t.Log("deploying a minimal HTTP container to be exoposed via ingress")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(testBulkIngressNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via ingress %d times", deployment.Name, maxBatchSize)
	for i := 0; i < maxBatchSize; i++ {
		name := fmt.Sprintf("bulk-httpbin-%d", i)
		path := fmt.Sprintf("/%s", name)

		service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		service.Name = name
		service.Spec.Type = corev1.ServiceTypeClusterIP
		_, err = env.Cluster().Client().CoreV1().Services(testBulkIngressNamespace).Create(ctx, service, metav1.CreateOptions{})
		require.NoError(t, err)

		ingress := generators.NewIngressForService(path, map[string]string{
			annotations.IngressClassKey: ingressClass,
			"konghq.com/strip-path":     "true",
		}, service)
		ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testBulkIngressNamespace).Create(ctx, ingress, metav1.CreateOptions{})
		require.NoError(t, err)
	}

	t.Logf("verifying that all %d ingresses route properly", maxBatchSize)
	for i := 0; i < maxBatchSize; i++ {
		name := fmt.Sprintf("bulk-httpbin-%d", i)
		path := fmt.Sprintf("/%s", name)

		require.Eventually(t, func() bool {
			resp, err := httpc.Get(fmt.Sprintf("%s/%s", proxyURL, path))
			if err != nil {
				t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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
		}, ingressWait, waitTick)
	}

	t.Log("cleaning up last batch of resources")
	for i := 0; i < maxBatchSize; i++ {
		name := fmt.Sprintf("bulk-httpbin-%d", i)
		require.NoError(t, env.Cluster().Client().CoreV1().Services(testBulkIngressNamespace).Delete(ctx, name, metav1.DeleteOptions{}))
		require.NoError(t, env.Cluster().Client().NetworkingV1().Ingresses(testBulkIngressNamespace).Delete(ctx, name, metav1.DeleteOptions{}))
		require.Eventually(t, func() bool {
			_, err := env.Cluster().Client().CoreV1().Services(testBulkIngressNamespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return errors.IsNotFound(err)
			}
			return false
		}, ingressWait, waitTick)
	}

	t.Log("staggering ingress deployments over several seconds")
	maxStaggeredBatchSize := maxBatchSize
	for i := 0; i < maxStaggeredBatchSize; i++ {
		name := fmt.Sprintf("bulk-staggered-httpbin-%d", i)
		path := fmt.Sprintf("/%s", name)

		service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		service.Name = name
		service.Spec.Type = corev1.ServiceTypeClusterIP
		_, err = env.Cluster().Client().CoreV1().Services(testBulkIngressNamespace).Create(ctx, service, metav1.CreateOptions{})
		require.NoError(t, err)

		ingress := generators.NewIngressForService(path, map[string]string{
			annotations.IngressClassKey: ingressClass,
			"konghq.com/strip-path":     "true",
		}, service)
		ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(testBulkIngressNamespace).Create(ctx, ingress, metav1.CreateOptions{})
		require.NoError(t, err)

		// every 10 items sleep for 1 second to stagger the updates to ~10/s
		if (i + 1%10) == 0 {
			time.Sleep(time.Second * 1)
		}
	}

	t.Logf("verifying that all %d staggered ingresses route properly", maxStaggeredBatchSize)
	for i := 0; i < maxStaggeredBatchSize; i++ {
		name := fmt.Sprintf("bulk-staggered-httpbin-%d", i)
		path := fmt.Sprintf("/%s", name)

		require.Eventually(t, func() bool {
			resp, err := httpc.Get(fmt.Sprintf("%s/%s", proxyURL, path))
			if err != nil {
				t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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
		}, ingressWait, waitTick)
	}

	t.Log("cleaning up last batch of resources")
	for i := 0; i < maxBatchSize; i++ {
		name := fmt.Sprintf("bulk-staggered-httpbin-%d", i)
		require.NoError(t, env.Cluster().Client().CoreV1().Services(testBulkIngressNamespace).Delete(ctx, name, metav1.DeleteOptions{}))
		require.NoError(t, env.Cluster().Client().NetworkingV1().Ingresses(testBulkIngressNamespace).Delete(ctx, name, metav1.DeleteOptions{}))
		require.Eventually(t, func() bool {
			_, err := env.Cluster().Client().CoreV1().Services(testBulkIngressNamespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return errors.IsNotFound(err)
			}
			return false
		}, ingressWait, waitTick)
	}
}
