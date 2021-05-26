//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMinimalPlugin(t *testing.T) {
	ctx := context.Background()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := k8sgen.NewContainer("httpbin", httpBinImage, 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	deployment, err := cluster.Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, cluster.Client().AppsV1().Deployments(corev1.NamespaceDefault).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = cluster.Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services(corev1.NamespaceDefault).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	ingress := k8sgen.NewIngressForService("/httpbin", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Create(ctx, ingress, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", ingress.Name)
		if err := cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}); err != nil {
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

	kongplugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      "teapot",
		},
		PluginName: "request-termination",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"status_code": 418}`),
		},
	}
	kongclusterplugin := &kongv1.KongClusterPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "legal",
		},
		PluginName: "request-termination",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"status_code": 451}`),
		},
	}
	c, err := clientset.NewForConfig(cluster.Config())
	assert.NoError(t, err)
	kongplugin, err = c.ConfigurationV1().KongPlugins(corev1.NamespaceDefault).Create(ctx, kongplugin, metav1.CreateOptions{})
	assert.NoError(t, err)
	kongclusterplugin, err = c.ConfigurationV1().KongClusterPlugins(corev1.NamespaceDefault).Create(ctx, kongclusterplugin, metav1.CreateOptions{})
	assert.NoError(t, err)

	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress.Name, metav1.GetOptions{})
	assert.NoError(t, err)

	t.Logf("updating Ingress %s to use plugin %s", ingress.Name, kongplugin.Name)
	ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongplugin.Name
	ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Update(ctx, ingress, metav1.UpdateOptions{})
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusTeapot
	}, ingressWait, waitTick)

	// TODO as of yet, this does not work. The controller does not detect the plugin and cannot apply it.
	// t.Logf("updating Ingress %s to use cluster plugin %s", ingress.Name, kongclusterplugin.Name)
	// ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongclusterplugin.Name
	// ingress, err = cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Update(ctx, ingress, metav1.UpdateOptions{})
	// assert.NoError(t, err)
	// assert.Eventually(t, func() bool {
	// 	resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
	// 	if err != nil {
	// 		t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
	// 		return false
	// 	}
	// 	defer resp.Body.Close()
	// 	return resp.StatusCode == http.StatusUnavailableForLegalReasons
	// }, ingressWait, waitTick)

	t.Logf("deleting Ingress %s and waiting for routes to be torn down", ingress.Name)
	assert.NoError(t, cluster.Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", p.ProxyURL.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", p.ProxyURL.String(), err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, p.ProxyURL.String(), resp)
	}, ingressWait, waitTick)
}
