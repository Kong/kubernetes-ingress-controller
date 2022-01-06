//go:build performance_tests
// +build performance_tests

package performance

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

func TestIngressBulk(t *testing.T) {
	ctx := context.Background()
	ns := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: corev1.NamespaceDefault}}

	t.Log("deploying a minimal HTTP container to be exoposed via ingress")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("checking the cluster version to determine which ingress version to use")
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)

	t.Logf("exposing deployment %s via ingress %d times", deployment.Name, maxBatchSize)
	ingresses := make([]runtime.Object, maxBatchSize)
	for i := 0; i < maxBatchSize; i++ {
		name := fmt.Sprintf("bulk-httpbin-%d", i)
		path := fmt.Sprintf("/%s", name)

		service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		service.Name = name
		service.Spec.Type = corev1.ServiceTypeClusterIP
		_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
		require.NoError(t, err)

		ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, path, map[string]string{
			annotations.IngressClassKey: ingressClass,
			"konghq.com/strip-path":     "true",
		}, service)
		require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
		ingresses[i] = ingress
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
		require.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, name, metav1.DeleteOptions{}))
		require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingresses[i]))
		require.Eventually(t, func() bool {
			_, err := env.Cluster().Client().CoreV1().Services(ns.Name).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return errors.IsNotFound(err)
			}
			return false
		}, ingressWait, waitTick)
	}

	t.Log("staggering ingress deployments over several seconds")
	maxStaggeredBatchSize := maxBatchSize
	ingresses = make([]runtime.Object, maxStaggeredBatchSize)
	for i := 0; i < maxStaggeredBatchSize; i++ {
		name := fmt.Sprintf("bulk-staggered-httpbin-%d", i)
		path := fmt.Sprintf("/%s", name)

		service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		service.Name = name
		service.Spec.Type = corev1.ServiceTypeClusterIP
		_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
		require.NoError(t, err)

		ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, path, map[string]string{
			annotations.IngressClassKey: ingressClass,
			"konghq.com/strip-path":     "true",
		}, service)
		require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
		ingresses[i] = ingress

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
		require.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, name, metav1.DeleteOptions{}))
		require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingresses[i]))
		require.Eventually(t, func() bool {
			_, err := env.Cluster().Client().CoreV1().Services(ns.Name).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return errors.IsNotFound(err)
			}
			return false
		}, ingressWait, waitTick)
	}
}
