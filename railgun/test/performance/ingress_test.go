//+build integration_tests

package performance

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestIngressPerf(t *testing.T) {
	ctx := context.Background()

	cnt := 1
	cost := 0
	for cnt < 5 {
		namespace := fmt.Sprintf("ingress-%d", cnt)

		t.Log("[%s] deploying a minimal HTTP container deployment to test Ingress routes", namespace)
		container := k8sgen.NewContainer("httpbin", httpBinImage, 80)
		deployment := k8sgen.NewDeploymentForContainer(container)
		deployment, err := cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Logf("[%s] exposing deployment %s via service", namespace, deployment.Name)
		service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		_, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Logf("[%s] creating an ingress for service %s with ingress.class %s", namespace, service.Name, ingressClass)
		ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{
			annotations.IngressClassKey: ingressClass,
			"konghq.com/strip-path":     "true",
		}, service)
		ingress, err = cluster.Client().NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})
		assert.NoError(t, err)
		start_time := time.Now().Nanosecond()
		end_time := time.Now().Nanosecond()

		t.Logf("waiting for routes from Ingress %s to be operational", ingress.Name)
		p := proxyReady()
		assert.Eventually(t, func() bool {
			resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
			if err != nil {
				t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
				return false
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				end_time = time.Now().Nanosecond()
				// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
				// Expected: "<title>httpbin.org</title>"
				b := new(bytes.Buffer)
				b.ReadFrom(resp.Body)
				return strings.Contains(b.String(), "<title>httpbin.org</title>")
			}
			return false
		}, ingressWait, waitTick)
		cost += end_time - start_time
		t.Logf("loop %d cost %v", cnt, cost)
	}

	t.Logf("ingress processing time %v", cost/cnt)
}
