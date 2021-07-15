//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
)

const testBulkIngressNamespace = "ingress-bulk-testing"

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

	t.Logf("exposing deployment %s via ingress 50 times", deployment.Name)
	for i := 0; i < 50; i++ {
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

	/*
		t.Log("verifying that all 50 ingresses receive status updates properly")
		for i := 0; i < 50; i++ {
			name := fmt.Sprintf("bulk-httpbin-%d", i)
			path := fmt.Sprintf("/%s", name)

			require.Eventually(t, func() bool {
				ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns).Get(ctx, name, metav1.GetOptions{})
				if err != nil {
					return false
				}
				return len(ingress.Status.LoadBalancer.Ingress) > 0
			}, ingressWait, waitTick)
		}
	*/

	t.Log("verifying that all 50 ingresses route properly")
	for i := 0; i < 50; i++ {
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
}
