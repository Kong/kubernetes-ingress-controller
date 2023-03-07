//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func TestKongIngressEssentials(t *testing.T) {
	ctx := context.Background()

	t.Parallel()
	ns := helpers.Namespace(ctx, t, env)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	testName := "minking"
	deployment := generators.NewDeploymentForContainer(generators.NewContainer(testName, test.HTTPBinImage, 80))
	_, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	c, err := clientset.NewForConfig(env.Cluster().Config())
	assert.NoError(t, err)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.Annotations = map[string]string{"konghq.com/override": testName}
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("routing to service %s via Ingress", service.Name)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/test_kongingress_essentials", map[string]string{
		annotations.IngressClassKey: consts.IngressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))

	defer func() {
		t.Log("ensuring that Ingress resources are cleaned up")
		assert.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
	}()

	t.Logf("applying service overrides to Service %s via KongIngress", service.Name)
	king := &kongv1.KongIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Proxy: &kongv1.KongIngressService{
			ReadTimeout: kong.Int(1000),
		},
	}
	king, err = c.ConfigurationV1().KongIngresses(ns.Name).Create(ctx, king, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring that KongIngress %s is cleaned up", king.Name)
		if err := c.ConfigurationV1().KongIngresses(ns.Name).Delete(ctx, king.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Log("waiting for routes from Ingress to be operational and that overrides are in place")

	assert.Eventually(t, func() bool {
		// Even though the HTTP client has a timeout of 10s, it should never be hit,
		// we expect a 504 from the proxy within 1000ms
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_kongingress_essentials/delay/5", proxyURL))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusGatewayTimeout
	}, ingressWait, waitTick)

	t.Logf("removing Service %s overrides", service.Name)
	svc, err := env.Cluster().Client().CoreV1().Services(ns.Name).Get(ctx, service.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	anns := svc.GetAnnotations()
	delete(anns, "konghq.com/override")
	svc.SetAnnotations(anns)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Update(ctx, svc, metav1.UpdateOptions{})
	assert.NoError(t, err)

	t.Logf("ensuring that Service %s overrides are eventually removed", service.Name)
	assert.Eventually(t, func() bool {
		url := fmt.Sprintf("%s/test_kongingress_essentials/delay/5", proxyURL)
		resp, err := helpers.DefaultHTTPClient().Get(url)
		if err != nil {
			t.Logf("failed issuing http GET for %q: %v", url, err)
			return false
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			return true
		default:
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Logf("failed reading response body (url: %q, status code: %d): %v",
					url, resp.StatusCode, err,
				)
				return false
			}

			t.Logf("response from %q: status code: %d; body: %s",
				url, resp.StatusCode, string(b),
			)
			return false
		}
	}, ingressWait, waitTick)
}
