//go:build performance_tests
// +build performance_tests

package performance

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
)

func TestTCPIngressPerformance(t *testing.T) {

	t.Log("setting up the TestTCPPerformance tests")
	cluster := env.Cluster()
	c, err := clientset.NewForConfig(cluster.Config())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
	defer cancel()

	cnt := 1
	cost := 0
	for cnt <= max_ingress {
		namespace := fmt.Sprintf("tcpingress-%d", cnt)
		err := CreateNamespace(ctx, namespace, t)
		require.NoError(t, err)

		t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
		testName := "tcpingress"
		deployment := generators.NewDeploymentForContainer(generators.NewContainer(testName, httpBinImage, 80))
		deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
		require.NoError(t, err)

		t.Logf("exposing deployment %s via service", deployment.Name)
		service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
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

		start_time := time.Now().Nanosecond()
		t.Logf("checking tcpingress %s status readiness.", tcp.Name)
		ingCli := c.ConfigurationV1beta1().TCPIngresses(namespace)
		assert.Eventually(t, func() bool {
			curIng, err := ingCli.Get(ctx, tcp.Name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				return false
			}
			ingresses := curIng.Status.LoadBalancer.Ingress
			for _, ingress := range ingresses {
				if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
					end_time := time.Now().Nanosecond()
					loop := end_time - start_time
					t.Logf("tcpingress hostname %s or ip %s is ready to redirect traffic after %d nanoseconds .", ingress.Hostname, ingress.IP, loop)
					cost += loop
					return true
				}
			}
			return false
		}, 120*time.Second, 1*time.Second, true)
		cnt += 1
	}
	t.Logf("tcp ingress average cost %d millisecond", cost/cnt/1000)
}
