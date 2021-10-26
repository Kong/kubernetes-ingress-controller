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
)

func TestIngressPerformance(t *testing.T) {
	t.Log("setting up the TestIngressPerf")

	ctx := context.Background()
	cluster := env.Cluster()
	cnt := 1
	cost := 0
	for cnt <= max_ingress {
		namespace := fmt.Sprintf("ingress-%d", cnt)
		err := CreateNamespace(ctx, namespace, t)
		assert.NoError(t, err)

		t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
		container := generators.NewContainer("httpbin", httpBinImage, 80)
		deployment := generators.NewDeploymentForContainer(container)
		deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Logf("exposing deployment %s via service", deployment.Name)
		service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		_, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Logf("[%s] creating an ingress for service httpbin with ingress.class %s", namespace, ingressClass)
		ingress := generators.NewIngressForService("/httpbin", map[string]string{
			annotations.IngressClassKey: ingressClass,
			"konghq.com/strip-path":     "true",
		}, service)
		ingress, err = cluster.Client().NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})
		assert.NoError(t, err)
		start_time := time.Now().Nanosecond()

		t.Logf("waiting for updated ingress status to include IP")
		require.Eventually(t, func() bool {
			ingress, err = cluster.Client().NetworkingV1().Ingresses(namespace).Get(ctx, ingress.Name, metav1.GetOptions{})
			if err != nil || ingress == nil {
				return false
			}
			t.Logf("debug ingress %v", ingress.Status.LoadBalancer.Ingress)
			if len(ingress.Status.LoadBalancer.Ingress) > 0 {
				return true
			}
			return false
		}, ingressWait, waitTick, true)

		t.Logf("ingress %v", ingress.Status.LoadBalancer.Ingress)
		end_time := time.Now().Nanosecond()
		loop := end_time - start_time
		t.Logf("networkingv1 hostname is ready to redirect traffic after %d nanosecond.", loop)
		cost += loop
		cnt += 1
	}
	t.Logf("ingress processing time %d nanosecond", cost/cnt)
}
