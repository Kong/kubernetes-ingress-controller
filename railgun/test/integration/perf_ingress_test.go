//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestIngressPerf(t *testing.T) {
	t.Log("setting up the TestIngressPerf")
	proxyReady()

	ctx := context.Background()

	cnt := 1
	cost := 0
	for cnt < 1 {

		namespace := fmt.Sprintf("ingress-%d", cnt)
		nsName := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		_, err := cluster.Client().CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Logf("[%s] deploying a minimal HTTP container deployment to test Ingress routes", namespace)
		err = configManifest(ctx, namespace, "httpbin.yaml", t)
		assert.NoError(t, err)

		t.Logf("[%s] List the service from the namespace.", namespace)
		svcs, err := cluster.Client().CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			t.Logf("Failed to list services.")
			continue
		}

		var service *v1.Service
		for _, svc := range svcs.Items {
			fmt.Printf("\n service %v \n", svc)
			if svc.Name == "httpbin" {

				service = &svc
				break
			}
		}

		t.Logf("[%s] creating an ingress for service httpbin with ingress.class %s", namespace, ingressClass)
		ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{
			annotations.IngressClassKey: ingressClass,
			"konghq.com/strip-path":     "true",
		}, service)
		ingress, err = cluster.Client().NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})
		assert.NoError(t, err)
		start_time := time.Now().Nanosecond()
		end_time := start_time

		t.Logf("checking networkingv1 %s status readiness.", service.Name)
		ingCli := cluster.Client().NetworkingV1().Ingresses(namespace)
		assert.Eventually(t, func() bool {
			curIng, err := ingCli.Get(ctx, service.Name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				return false
			}
			ingresses := curIng.Status.LoadBalancer.Ingress
			for _, ingress := range ingresses {
				if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
					end_time = time.Now().Nanosecond()
					t.Logf("networkingv1 hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
					return true
				}
			}
			return false
		}, 120*time.Second, 1*time.Second, true)
		cost += end_time - start_time
		t.Logf("KIC process the %d ingress cost %v", cnt, cost)

		t.Logf("checking routes from Ingress %s to be operational", ingress.Name)
		p := proxyReady()
		assert.Eventually(t, func() bool {
			resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
			if err != nil {
				t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
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
	}
	t.Logf("ingress processing time %v", cost/cnt)
	cnt += 1
}

func configManifest(ctx context.Context, namespace, yml string, t *testing.T) error {
	cmd := exec.CommandContext(ctx, "kubectl", "-n", namespace, "apply", "-f", yml)
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stdout, stdout.String())
		return err
	}
	t.Logf("successfully deploy manifest " + yml)
	return nil
}
