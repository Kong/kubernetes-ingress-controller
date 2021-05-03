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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestMinimalIngress(t *testing.T) {
	ctx := context.Background()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := k8sgen.NewContainer("httpbin", httpBinImage, 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	deployment, err := cluster.Client().AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, cluster.Client().AppsV1().Deployments("default").Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = cluster.Client().CoreV1().Services("default").Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services("default").Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("routing to service %s via Ingress", service.Name)
	ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{
		"kubernetes.io/ingress.class": ingressClass,
		"konghq.com/strip-path":       "true",
	}, service)
	ingress, err = cluster.Client().NetworkingV1().Ingresses("default").Create(ctx, ingress, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", ingress.Name)
		if err := cluster.Client().NetworkingV1().Ingresses("default").Delete(ctx, ingress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

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
			// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
			// Expected: "<title>httpbin.org</title>"
			b := new(bytes.Buffer)
			b.ReadFrom(resp.Body)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("deleting Ingress %s and waiting for routes to be torn down", ingress.Name)
	assert.NoError(t, cluster.Client().NetworkingV1().Ingresses("default").Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
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
				t.Logf("WARNING: error decoding JSON from proxy while waiting for %s: %v", p.ProxyURL.String(), err)
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

	t.Log("creating an HTTP container via deployment to test redirect functionality")
	container := k8sgen.NewContainer("alsohttpbin", httpBinImage, 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	_, err := cluster.Client().AppsV1().Deployments("default").Create(ctx, deployment, opts)
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up deployment %s", deployment.Name)
		assert.NoError(t, cluster.Client().AppsV1().Deployments("default").Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via Service", deployment.Name)
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
	service, err = cluster.Client().CoreV1().Services("default").Create(ctx, service, opts)
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up Service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services("default").Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing Service %s via Ingress", service.Name)
	ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{
		"kubernetes.io/ingress.class":           ingressClass,
		"konghq.com/protocols":                  "https",
		"konghq.com/https-redirect-status-code": "301",
	}, service)
	ingress, err = cluster.Client().NetworkingV1().Ingresses("default").Create(ctx, ingress, opts)
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up Ingress %s", ingress.Name)
		assert.NoError(t, cluster.Client().NetworkingV1().Ingresses("default").Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("waiting for Ingress %s to be operational and properly redirect", ingress.Name)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 3,
	}
	assert.Eventually(t, func() bool {
		p := proxyReady()
		resp, err := client.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusMovedPermanently
	}, ingressWait, waitTick)
}
