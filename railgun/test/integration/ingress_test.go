//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestMinimalIngress(t *testing.T) {
	ctx := context.Background()

	// deploy a minimal deployment to test Ingress routes to
	container := k8sgen.NewContainer("nginx", "nginx", 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	deployment, err := cluster.Client().AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	// ensure cleanup of the deployment
	defer func() {
		assert.NoError(t, cluster.Client().AppsV1().Deployments("default").Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	// expose the deployment via service
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = cluster.Client().CoreV1().Services("default").Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	// route to the service via Kong Ingress
	ingress := k8sgen.NewIngressForService("/nginx", map[string]string{
		"kubernetes.io/ingress.class": "kong",
		"konghq.com/strip-path":       "true",
	}, service)
	cluster.Client().NetworkingV1().Ingresses("default").Create(ctx, ingress, metav1.CreateOptions{})

	// wait for the ingress backend to be routable
	u := proxyURL()
	assert.Eventually(t, func() bool {
		resp, err := http.Get(fmt.Sprintf("%s/nginx", u.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", u.String(), err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
			// Expected: Welcome to nginx!
			b := new(bytes.Buffer)
			b.ReadFrom(resp.Body)
			return strings.Contains(b.String(), "Welcome to nginx!")
		}
		return false
	}, ingressWait, waitTick)

	// ensure that a deleted ingress results in the route being torn down
	assert.NoError(t, cluster.Client().NetworkingV1().Ingresses("default").Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	assert.Eventually(t, func() bool {
		resp, err := http.Get(fmt.Sprintf("%s/nginx", u.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", u.String(), err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			// once the route is torn down and returning 404's, ensure that we got the expected response body back from Kong
			// Expected: {"message":"no Route matched with those values"}
			b := new(bytes.Buffer)
			b.ReadFrom(resp.Body)
			body := struct {
				Message string `json:"message"`
			}{}
			if err := json.Unmarshal(b.Bytes(), &body); err != nil {
				t.Logf("WARNING: error decoding JSON from proxy while waiting for %s: %v", u.String(), err)
				return false
			}
			return body.Message == "no Route matched with those values"
		}
		return false
	}, ingressWait, waitTick)
}

func TestHTTPSRedirect(t *testing.T) {
	ctx := context.Background()
	opts := metav1.CreateOptions{}

	container := k8sgen.NewContainer("alsonginx", "nginx", 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	_, err := cluster.Client().AppsV1().Deployments("default").Create(ctx, deployment, opts)
	assert.NoError(t, err)

	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
	_, err = cluster.Client().CoreV1().Services("default").Create(ctx, service, opts)
	assert.NoError(t, err)

	ingress := k8sgen.NewIngressForService("/example", map[string]string{
		"kubernetes.io/ingress.class":           "kong",
		"konghq.com/protocols":                  "https",
		"konghq.com/https-redirect-status-code": "301",
	}, service)
	_, err = cluster.Client().NetworkingV1().Ingresses("default").Create(ctx, ingress, opts)
	assert.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	assert.Eventually(t, func() bool {
		u := proxyURL()
		resp, err := client.Get(fmt.Sprintf("%s/example", u.String()))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusMovedPermanently
	}, ingressWait, waitTick)
}
