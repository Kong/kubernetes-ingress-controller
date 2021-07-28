//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/internal/annotations"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/pkg/clientset"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
)

const testTCPIngressNamespace = "tcpingress"

func TestTCPIngressEssentials(t *testing.T) {
	t.Log("setting up the TCPIngress tests")
	testName := "tcpingress"
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
	defer cancel()

	t.Logf("creating namespace %s for testing", testTCPIngressNamespace)
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testTCPIngressNamespace}}
	ns, err = env.Cluster().Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up namespace %s", testTCPIngressNamespace)
		require.NoError(t, env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, ns.Name, metav1.DeleteOptions{}))
		require.Eventually(t, func() bool {
			_, err := env.Cluster().Client().CoreV1().Namespaces().Get(ctx, ns.Name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					return true
				}
			}
			return false
		}, ingressWait, waitTick)
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	deployment := generators.NewDeploymentForContainer(generators.NewContainer(testName, httpBinImage, 80))
	deployment, err = env.Cluster().Client().AppsV1().Deployments(testTCPIngressNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(testTCPIngressNamespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = env.Cluster().Client().CoreV1().Services(testTCPIngressNamespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(testTCPIngressNamespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("routing to service %s via TCPIngress", service.Name)
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: testTCPIngressNamespace,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Port: 8888,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: service.Name,
						ServicePort: 80,
					},
				},
			},
		},
	}
	tcp, err = c.ConfigurationV1beta1().TCPIngresses(testTCPIngressNamespace).Create(ctx, tcp, metav1.CreateOptions{})
	require.NoError(t, err)
	defer func() {
		t.Logf("ensuring that TCPIngress %s is cleaned up", tcp.Name)
		if err := c.ConfigurationV1beta1().TCPIngresses(testTCPIngressNamespace).Delete(ctx, tcp.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("checking tcpingress %s status readiness.", tcp.Name)
	ingCli := c.ConfigurationV1beta1().TCPIngresses(testTCPIngressNamespace)
	assert.Eventually(t, func() bool {
		curIng, err := ingCli.Get(ctx, tcp.Name, metav1.GetOptions{})
		if err != nil || curIng == nil {
			return false
		}
		ingresses := curIng.Status.LoadBalancer.Ingress
		for _, ingress := range ingresses {
			if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
				t.Logf("tcpingress hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
				return true
			}
		}
		return false
	}, 120*time.Second, 1*time.Second, true)

	t.Logf("verifying TCP Ingress %s operationalable", tcp.Name)
	tcpProxyURL, err := url.Parse(fmt.Sprintf("http://%s:8888/", proxyURL.Hostname()))
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(tcpProxyURL.String())
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
			// Expected: "<title>httpbin.org</title>"
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("tearing down TCPIngress %s and ensuring that the relevant backend routes are removed", tcp.Name)
	require.NoError(t, c.ConfigurationV1beta1().TCPIngresses(testTCPIngressNamespace).Delete(ctx, tcp.Name, metav1.DeleteOptions{}))
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(tcpProxyURL.String())
		if err != nil {
			return true
		}
		defer resp.Body.Close()
		return false
	}, ingressWait, waitTick)
}
