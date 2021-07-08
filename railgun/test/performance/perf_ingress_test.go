//+build performance_tests

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

func TestIngressPerformance(t *testing.T) {
	t.Log("setting up the TestIngressPerf")

	ctx := context.Background()
	cnt := 1
	cost := 0
	for cnt <= max_ingress {
		namespace := fmt.Sprintf("ingress-%d", cnt)
		nsName := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		_, err := cluster.Client().CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
		container := k8sgen.NewContainer("httpbin", httpBinImage, 80)
		deployment := k8sgen.NewDeploymentForContainer(container)
		deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Logf("exposing deployment %s via service", deployment.Name)
		service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		_, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Logf("[%s] creating an ingress for service httpbin with ingress.class %s", namespace, ingressClass)
		ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{
			annotations.IngressClassKey: ingressClass,
			"konghq.com/strip-path":     "true",
		}, service)
		ingress, err = cluster.Client().NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})
		assert.NoError(t, err)
		start_time := time.Now().Nanosecond()

		t.Logf("checking routes from Ingress %s to be operational", ingress.Name)
		assert.Eventually(t, func() bool {
			resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", KongInfo.ProxyURL.String()))
			if err != nil {
				t.Logf("WARNING: error while waiting for %s: %v", KongInfo.ProxyURL.String(), err)
				return false
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
				// Expected: "<title>httpbin.org</title>"
				b := new(bytes.Buffer)
				b.ReadFrom(resp.Body)
				end_time := time.Now().Nanosecond()
				loop := end_time - start_time
				t.Logf("networkingv1 hostname is ready to redirect traffic after %d.", loop)
				cost += loop
				return strings.Contains(b.String(), "<title>httpbin.org</title>")
			}
			return false
		}, ingressWait, time.Millisecond * 1)
		cnt += 1
	}
	t.Logf("ingress processing time %d nanosecond", cost/cnt)
}
