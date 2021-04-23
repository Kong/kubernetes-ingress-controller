//+build integration_tests

package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/go-kong/kong"
	kongv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestMinimalKongIngress(t *testing.T) {
	// test setup
	namespace := "default"
	testName := "minking"
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
	defer cancel()

	// push a minimal deployment to test KongIngress routes to
	deployment := k8sgen.NewDeploymentForContainer(k8sgen.NewContainer(testName, "kennethreitz/httpbin", 80))
	_, err := cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	// make sure we clean up after ourselves
	defer func() {
		assert.NoError(t, cluster.Client().AppsV1().Deployments(namespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	// initialize a clientset for the KongIngress API
	c, err := clientset.NewForConfig(cluster.Config())
	assert.NoError(t, err)

	// expose the deployment via service
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.Annotations = map[string]string{"konghq.com/override": testName}
	service, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	// make sure we clean up after ourselves
	defer func() {
		assert.NoError(t, cluster.Client().CoreV1().Services(namespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	// route to the service via Ingress
	ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{
		"kubernetes.io/ingress.class": "kong",
		"konghq.com/strip-path":       "true",
	}, service)
	ingress, err = cluster.Client().NetworkingV1().Ingresses("default").Create(ctx, ingress, metav1.CreateOptions{})
	assert.NoError(t, err)

	// ensure cleanup of the ingress
	defer func() {
		assert.NoError(t, cluster.Client().NetworkingV1().Ingresses("default").Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	}()

	// TODO: this is a workaround for https://github.com/Kong/kubernetes-ingress-controller/issues/1146
	time.Sleep(time.Second * 30)

	// deploy the KongIngress object to apply overrides
	king := &kongv1.KongIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "kong",
			},
		},
		Proxy: &kong.Service{
			ReadTimeout: kong.Int(1),
		},
	}
	king, err = c.ConfigurationV1().KongIngresses(namespace).Create(ctx, king, metav1.CreateOptions{})
	assert.NoError(t, err)

	// make sure we clean up after ourselves
	defer func() {
		assert.NoError(t, c.ConfigurationV1().KongIngresses(namespace).Delete(ctx, king.Name, metav1.DeleteOptions{}))
	}()

	// test that the read delay works properly
	assert.Eventually(t, func() bool {
		p := proxyReady()
		resp, err := http.Get(fmt.Sprintf("%s/httpbin/delay/5", p.ProxyURL.String()))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusGatewayTimeout
	}, ingressWait, waitTick)
}
