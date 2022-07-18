//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

func TestPluginEssentials(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/httpbin", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))

	defer func() {
		t.Log("ensuring that Ingress is cleaned up")
		if err := clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Log("waiting for routes from Ingress to be operational")
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
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	kongplugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
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
	require.NoError(t, err)
	kongplugin, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, kongplugin, metav1.CreateOptions{})
	require.NoError(t, err)
	kongclusterplugin, err = c.ConfigurationV1().KongClusterPlugins().Create(ctx, kongclusterplugin, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Log("cleaning up plugins")
		if err := c.ConfigurationV1().KongPlugins(ns.Name).Delete(ctx, kongplugin.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
		if err := c.ConfigurationV1().KongClusterPlugins().Delete(ctx, kongclusterplugin.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Logf("updating Ingress to use plugin %s", kongplugin.Name)
	require.Eventually(t, func() bool {
		switch obj := ingress.(type) {
		case *netv1.Ingress:
			ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
			if err != nil {
				return false
			}
			ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongplugin.Name
			_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
			return err == nil
		case *netv1beta1.Ingress:
			ingress, err := env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
			if err != nil {
				return false
			}
			ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongplugin.Name
			_, err = env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
			return err == nil
		}
		return false
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

	t.Logf("updating Ingress to use cluster plugin %s", kongclusterplugin.Name)
	require.Eventually(t, func() bool {
		switch obj := ingress.(type) {
		case *netv1.Ingress:
			ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
			if err != nil {
				return false
			}
			ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongclusterplugin.Name
			_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
			return err == nil
		case *netv1beta1.Ingress:
			ingress, err := env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Get(ctx, obj.Name, metav1.GetOptions{})
			if err != nil {
				return false
			}
			ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongclusterplugin.Name
			_, err = env.Cluster().Client().NetworkingV1beta1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
			return err == nil
		}
		return false
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

	t.Log("deleting Ingress and waiting for routes to be torn down")
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
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

func TestPluginOrdering(t *testing.T) {
	if !util.GetKongVersion().GTE(semver.MustParse("3.0.0")) || kongEnterpriseEnabled == "" {
		t.Skip("plugin ordering requires Kong Enterprise 3.0+")
	}
	t.Parallel()
	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	ingress := generators.NewIngressForService("/httpbin", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Log("waiting for routes from Ingress to be operational")
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
			n, err := b.ReadFrom(resp.Body)
			if err != nil || n == 0 {
				return false
			}
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	termplugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "teapot",
		},
		PluginName: "request-termination",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"status_code": 418}`),
		},
	}
	authplugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "auth",
		},
		PluginName: "key-auth",
	}
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	termplugin, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, termplugin, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(termplugin)
	authplugin, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, authplugin, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(authplugin)

	t.Logf("updating Ingress to use plugin %s", termplugin.Name)
	require.Eventually(t, func() bool {
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = termplugin.Name
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Logf("validating that plugin %s was successfully configured", termplugin.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusTeapot
	}, ingressWait, waitTick)

	t.Logf("updating Ingress to use plugin %s", authplugin.Name)
	require.Eventually(t, func() bool {
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = strings.Join([]string{authplugin.Name, termplugin.Name}, ",")
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Logf("validating that plugin %s was successfully configured", authplugin.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusUnauthorized
	}, ingressWait, waitTick)

	t.Logf("updating plugin %s order to run before key-auth in access", termplugin.Name)
	termplugin, err = c.ConfigurationV1().KongPlugins(ns.Name).Get(ctx, termplugin.Name, metav1.GetOptions{})
	require.NoError(t, err)
	termplugin.Ordering = &kong.PluginOrdering{
		Before: kong.PluginOrderingPhase{
			"access": []string{"key-auth"},
		},
	}
	_, err = c.ConfigurationV1().KongPlugins(ns.Name).Update(ctx, termplugin, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Logf("validating that plugin %s now takes precedence", termplugin.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusTeapot
	}, ingressWait, waitTick)

	t.Log("deleting Ingress and waiting for routes to be torn down")
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
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
