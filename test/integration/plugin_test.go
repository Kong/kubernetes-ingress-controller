//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

func TestPluginEssentials(t *testing.T) {
	ctx := context.Background()

	t.Parallel()
	ns, cleaner := helpers.Setup(ctx, t, env)

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

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, consts.IngressClass)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/test_plugin_essentials", map[string]string{
		annotations.IngressClassKey: consts.IngressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	helpers.AddIngressToCleaner(cleaner, ingress)

	t.Log("waiting for routes from Ingress to be operational")
	assert.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_essentials", proxyURL))
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
				annotations.IngressClassKey: consts.IngressClass,
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
	cleaner.Add(kongplugin)
	cleaner.Add(kongclusterplugin)

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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_essentials", proxyURL))
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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_essentials", proxyURL))
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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_essentials", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return helpers.ExpectHTTP404WithNoRoute(t, proxyURL, resp)
	}, ingressWait, waitTick)
}

func TestPluginOrdering(t *testing.T) {
	t.Parallel()

	// the manager runs in a goroutine and may not have pulled the version before this test starts
	require.Eventually(t, func() bool {
		return !versions.GetKongVersion().Full().EQ(semver.MustParse("0.0.0"))
	}, time.Minute, time.Second)

	{
		v := versions.GetKongVersion()
		enterprise := testenv.KongEnterpriseEnabled()
		if !v.MajorOnly().GTE(versions.PluginOrderingVersionCutoff) || !enterprise {
			t.Skipf("plugin ordering requires Kong Enterprise 3.0+, detected: %s (enterprise: %v)", v.Full(), enterprise)
		}
	}

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

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

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, consts.IngressClass)
	ingress := generators.NewIngressForService("/test_plugin_ordering", map[string]string{
		annotations.IngressClassKey: consts.IngressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Log("waiting for routes from Ingress to be operational")
	assert.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_ordering", proxyURL))
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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_ordering", proxyURL))
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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_ordering", proxyURL))
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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_ordering", proxyURL))
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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_ordering", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return helpers.ExpectHTTP404WithNoRoute(t, proxyURL, resp)
	}, ingressWait, waitTick)
}
