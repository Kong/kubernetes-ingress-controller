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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

const testTCPIngressNamespace = "tcpingress"

func TestTCPIngress(t *testing.T) {
	t.Log("setting up the TCPIngress tests")
	p := proxyReady()
	testName := "tcpingress"
	c, err := clientset.NewForConfig(cluster.Config())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
	defer cancel()

	t.Logf("creating namespace %s for testing TCPIngress", testTCPIngressNamespace)
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testTCPIngressNamespace}}
	ns, err = cluster.Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up namespace %s", testTCPIngressNamespace)
		assert.NoError(t, cluster.Client().CoreV1().Namespaces().Delete(ctx, ns.Name, metav1.DeleteOptions{}))
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	deployment := k8sgen.NewDeploymentForContainer(k8sgen.NewContainer(testName, httpBinImage, 80))
	deployment, err = cluster.Client().AppsV1().Deployments(testTCPIngressNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, cluster.Client().AppsV1().Deployments(testTCPIngressNamespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = cluster.Client().CoreV1().Services(testTCPIngressNamespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services(testTCPIngressNamespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
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

	t.Logf("waiting for routes from Ingress %s to be operational", tcp.Name)
	tcpProxyURL, err := url.Parse(fmt.Sprintf("http://%s:8888/", p.ProxyURL.Hostname()))
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
			b.ReadFrom(resp.Body)
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
