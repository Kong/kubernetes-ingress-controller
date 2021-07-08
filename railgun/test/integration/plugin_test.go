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
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPluginEssentials(t *testing.T) {
	ctx := context.Background()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	ingress := generators.NewIngressForService("/httpbin", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Create(ctx, ingress, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", ingress.Name)
		if err := env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("waiting for routes from Ingress %s to be operational", ingress.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
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
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		PluginName: "request-termination",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"status_code": 451}`),
		},
	}
	c, err := clientset.NewForConfig(env.Cluster().Config())
	assert.NoError(t, err)
	kongplugin, err = c.ConfigurationV1().KongPlugins(corev1.NamespaceDefault).Create(ctx, kongplugin, metav1.CreateOptions{})
	assert.NoError(t, err)
	kongclusterplugin, err = c.ConfigurationV1().KongClusterPlugins().Create(ctx, kongclusterplugin, metav1.CreateOptions{})
	assert.NoError(t, err)

	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress.Name, metav1.GetOptions{})
	assert.NoError(t, err)

	t.Logf("updating Ingress %s to use plugin %s", ingress.Name, kongplugin.Name)
	require.Eventually(t, func() bool {
		ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongplugin.Name
		ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Update(ctx, ingress, metav1.UpdateOptions{})
		if err != nil {
			return false
		}
		return true
	}, ingressWait, waitTick)

	t.Logf("validating that plugin %s was successfully configured", kongplugin.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusTeapot
	}, ingressWait, waitTick)

	t.Logf("updating Ingress %s to use cluster plugin %s", ingress.Name, kongclusterplugin.Name)
	require.Eventually(t, func() bool {
		ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongclusterplugin.Name
		ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Update(ctx, ingress, metav1.UpdateOptions{})
		if err != nil {
			return false
		}
		return true
	}, ingressWait, waitTick)

	t.Logf("validating that clusterplugin %s was successfully configured", kongclusterplugin.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusUnavailableForLegalReasons
	}, ingressWait, waitTick)

	t.Logf("deleting Ingress %s and waiting for routes to be torn down", ingress.Name)
	assert.NoError(t, env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return expect404WithNoRoute(t, proxyURL.String(), resp)
	}, ingressWait, waitTick)
}
