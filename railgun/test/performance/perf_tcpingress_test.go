//+build performance_tests

package performance

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestTCPPerformance(t *testing.T) {

	t.Log("setting up the TestTCPPerformance tests")
	c, err := clientset.NewForConfig(cluster.Config())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
	defer cancel()

	cnt := 1
	cost := 0
	for cnt <= max_ingress {
		namespace := fmt.Sprintf("tcpingress-%d", cnt)
		nsName := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		t.Logf("creating namespace %s for testing TCPIngress", namespace)
		_, err := cluster.Client().CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
		testName := "tcpingress"
		deployment := k8sgen.NewDeploymentForContainer(k8sgen.NewContainer(testName, httpBinImage, 80))
		deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
		require.NoError(t, err)

		t.Logf("exposing deployment %s via service", deployment.Name)
		service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		service, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
		require.NoError(t, err)

		t.Logf("routing to service %s via TCPIngress", service.Name)
		tcp := &kongv1beta1.TCPIngress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: namespace,
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
		tcp, err = c.ConfigurationV1beta1().TCPIngresses(namespace).Create(ctx, tcp, metav1.CreateOptions{})
		require.NoError(t, err)

		t.Logf("waiting for routes from Ingress %s to be operational", tcp.Name)
		s := time.Now().Nanosecond()
		tcpProxyURL, err := url.Parse(fmt.Sprintf("http://%s:8888/", KongInfo.ProxyURL.Hostname()))
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
				e := time.Now().Nanosecond()
				loop := e - s
				t.Logf("tcp ingress loop cost %d nanosecond", loop)
				cost += loop
				return strings.Contains(b.String(), "<title>httpbin.org</title>")
			}
			return false
		}, ingressWait, waitTick)
		cnt += 1
	}
	t.Logf("tcp ingress cost %d millisecond", cost/cnt/1000)
}
