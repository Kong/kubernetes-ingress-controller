//go:build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	"github.com/kong/kubernetes-configuration/pkg/clientset"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func TestPluginEssentials(t *testing.T) {
	ctx := context.Background()

	t.Parallel()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
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
	ingress := generators.NewIngressForService("/test_plugin_essentials", map[string]string{
		"konghq.com/strip-path": "true",
	}, service)
	ingress.Spec.IngressClassName = kong.String(consts.IngressClass)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Log("waiting for routes from Ingress to be operational")
	assert.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_essentials", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
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
		InstanceName: "example",
		PluginName:   "request-termination",
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
		InstanceName: "example-cluster",
		PluginName:   "request-termination",
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
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongplugin.Name
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Logf("validating that plugin %s was successfully configured", kongplugin.Name)
	assert.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_essentials", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusTeapot
	}, ingressWait, waitTick)

	t.Logf("updating Ingress to use cluster plugin %s", kongclusterplugin.Name)
	require.Eventually(t, func() bool {
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongclusterplugin.Name
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Logf("validating that clusterplugin %s was successfully configured", kongclusterplugin.Name)
	assert.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_essentials", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusUnavailableForLegalReasons
	}, ingressWait, waitTick)

	t.Log("deleting Ingress and waiting for routes to be torn down")
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
	helpers.EventuallyExpectHTTP404WithNoRoute(t, proxyHTTPURL, proxyHTTPURL.Host, "/test_plugin_essentials", ingressWait, waitTick, nil)
}

func TestPluginConfigPatch(t *testing.T) {
	ctx := context.Background()

	t.Parallel()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
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
	ingress := generators.NewIngressForService("/test_plugin_essentials", map[string]string{
		"konghq.com/strip-path": "true",
	}, service)
	ingress.Spec.IngressClassName = kong.String(consts.IngressClass)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Log("waiting for routes from Ingress to be operational")
	assert.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_essentials", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
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

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "kongplugin-config",
		},
		StringData: map[string]string{
			"teapot-message":    `"I am a little teapot"`,
			"forbidden-message": `"This service is forbidden"`,
		},
	}
	secret, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, secret, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(secret)

	kongplugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "teapot",
		},
		InstanceName: "example",
		PluginName:   "request-termination",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"status_code": 418}`),
		},
		ConfigPatches: []kongv1.ConfigPatch{
			{
				Path: "/message",
				ValueFrom: kongv1.ConfigSource{
					SecretValue: kongv1.SecretValueFromSource{
						Secret: "kongplugin-config",
						Key:    "teapot-message",
					},
				},
			},
		},
	}
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	kongplugin, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, kongplugin, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(kongplugin)

	kongclusterplugin := &kongv1.KongClusterPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "forbidden",
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		InstanceName: "example-cluster",
		PluginName:   "request-termination",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"status_code": 403}`),
		},
		ConfigPatches: []kongv1.NamespacedConfigPatch{
			{
				Path: "/message",
				ValueFrom: kongv1.NamespacedConfigSource{
					SecretValue: kongv1.NamespacedSecretValueFromSource{
						Namespace: ns.Name,
						Secret:    "kongplugin-config",
						Key:       "forbidden-message",
					},
				},
			},
		},
	}
	kongclusterplugin, err = c.ConfigurationV1().KongClusterPlugins().Create(ctx, kongclusterplugin, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(kongclusterplugin)

	t.Logf("updating Ingress to use plugin %s", kongplugin.Name)
	require.Eventually(t, func() bool {
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongplugin.Name
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Logf("validating that plugin %s was successfully configured", kongplugin.Name)
	assert.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_essentials", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
			return false
		}
		defer resp.Body.Close()
		buf, _ := io.ReadAll(resp.Body)
		return resp.StatusCode == http.StatusTeapot && strings.Contains(string(buf), "I am a little teapot")
	}, ingressWait, waitTick)

	t.Logf("updating Ingress to use cluster plugin %s", kongclusterplugin.Name)
	require.Eventually(t, func() bool {
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = kongclusterplugin.Name
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Logf("validating that clusterplugin %s was successfully configured", kongclusterplugin.Name)
	assert.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_essentials", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
			return false
		}
		defer resp.Body.Close()
		buf, _ := io.ReadAll(resp.Body)
		return resp.StatusCode == http.StatusForbidden && strings.Contains(string(buf), "This service is forbidden")
	}, ingressWait, waitTick)

	t.Log("deleting Ingress and waiting for routes to be torn down")
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
	helpers.EventuallyExpectHTTP404WithNoRoute(t, proxyHTTPURL, proxyHTTPURL.Host, "/test_plugin_essentials", ingressWait, waitTick, nil)
}

func TestPluginOrdering(t *testing.T) {
	t.Parallel()

	RunWhenKongEnterprise(t)

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
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
		"konghq.com/strip-path": "true",
	}, service)
	ingress.Spec.IngressClassName = kong.String(consts.IngressClass)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Log("waiting for routes from Ingress to be operational")
	assert.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_ordering", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_ordering", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_ordering", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_ordering", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusTeapot
	}, ingressWait, waitTick)

	t.Log("deleting Ingress and waiting for routes to be torn down")
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
	helpers.EventuallyExpectHTTP404WithNoRoute(t, proxyHTTPURL, proxyHTTPURL.Host, "/test_plugin_ordering", ingressWait, waitTick, nil)
}

// TestPluginCrossNamespaceReference creates an Ingress in one namespace and a KongConsumer in another. It tests that
// the controller can generate a plugin attached to the generated route and consumer if and only if a ReferenceGrant
// allows access to the KongPlugin from KongConsumers in a remote namespace.
func TestPluginCrossNamespaceReference(t *testing.T) {
	ctx := context.Background()

	t.Parallel()
	ns, cleaner := helpers.Setup(ctx, t, env)
	cluster := env.Cluster()
	remote, err := clusters.GenerateNamespace(ctx, cluster, helpers.LabelValueForTest(t))
	require.NoError(t, err)
	cleaner.AddNamespace(remote)

	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, consts.IngressClass)
	ingress := generators.NewIngressForService("/test_plugin_reference", map[string]string{
		"konghq.com/strip-path": "true",
	}, service)
	ingress.Spec.IngressClassName = kong.String(consts.IngressClass)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Log("waiting for routes from Ingress to be operational")
	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_plugin_reference", proxyHTTPURL))
		if !assert.NoError(c, err) {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			_, err := b.ReadFrom(resp.Body)
			assert.NoError(c, err)
			assert.True(c, strings.Contains(b.String(), "<title>httpbin.org</title>"))
		}
	}, ingressWait, waitTick)

	t.Log("creating plugins")
	kongplugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "teapot",
		},
		InstanceName: "example",
		PluginName:   "request-termination",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"status_code": 418}`),
		},
	}
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	kongplugin, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, kongplugin, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(kongplugin)

	// this isn't strictly necessary, but it's convenient for diagnosing issues if the test breaks, as it confirms
	// that a request was indeed recognized as being from the consumer. unlike the other consumer plugin, it's associated
	// with the consumer alone and works regardless of the grant state or any cross-namespace shenanigans
	kongslugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "snurch",
		},
		InstanceName: "snurch",
		PluginName:   "correlation-id",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"header_name": "snurch", "echo_downstream": true}`),
		},
	}
	kongslugin, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, kongslugin, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(kongslugin)

	t.Log("creating consumer and credential")
	credential := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Labels: map[string]string{
				labels.LabelPrefix + labels.CredentialKey: "key-auth",
			},
		},
		StringData: map[string]string{
			"key": "thirtytangas",
		},
	}
	_, err = env.Cluster().Client().CoreV1().Secrets(remote.Name).Create(ctx, credential, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(credential)

	consumer := &kongv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey:                           consts.IngressClass,
				annotations.AnnotationPrefix + annotations.PluginsKey: strings.Join([]string{kongslugin.Name, fmt.Sprintf("%s:%s", kongplugin.Namespace, kongplugin.Name)}, ","),
			},
		},
		Username:    uuid.NewString(),
		Credentials: []string{credential.Name},
	}
	_, err = c.ConfigurationV1().KongConsumers(remote.Name).Create(ctx, consumer, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(consumer)

	t.Log("creating auth plugin to identify consumer accessing route")
	authplugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "keykey",
		},
		PluginName: "key-auth",
	}
	authplugin, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, authplugin, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(authplugin)

	cleaner.Add(kongplugin)
	t.Logf("updating Ingress to use auth plugin %s", kongplugin.Name)
	require.Eventually(t, func() bool {
		ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = strings.Join([]string{kongplugin.Name, authplugin.Name}, ",")
		_, err = env.Cluster().Client().NetworkingV1().Ingresses(ns.Name).Update(ctx, ingress, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	const (
		// negativeCheckWait is the duration used in `Never` for verifying that the plugin is not configured
		// without reference grant or with incorrect reference grant.
		// Set it to 1 minute since we have to wait until the end in each `Never` if the test is OK.
		negativeCheckWait = time.Minute
	)

	t.Logf("validating that plugin %s is not configured without a grant", kongplugin.Name)
	assert.Never(t, func() bool {
		req := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.String(), "/test_plugin_reference?key=thirtytangas", nil)
		resp, err := helpers.DefaultHTTPClient(helpers.WithResolveHostTo(proxyHTTPURL.Host)).Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusTeapot
	}, negativeCheckWait, waitTick)

	t.Logf("creating a ReferenceGrant that does not permit kongconsumer access to kongplugins")
	grant := &gatewayapi.ReferenceGrant{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayapi.ReferenceGrantSpec{
			From: []gatewayapi.ReferenceGrantFrom{
				{
					Group:     gatewayapi.Group("configuration.konghq.com"),
					Kind:      gatewayapi.Kind("KongConsumer"),
					Namespace: gatewayapi.Namespace(remote.Name),
				},
			},
			To: []gatewayapi.ReferenceGrantTo{
				{
					Group: gatewayapi.Group("configuration.konghq.com"),
					Kind:  gatewayapi.Kind("KongPlugin"),
				},
			},
		},
	}
	// Not the namespace as the plugin.
	_, err = gatewayClient.GatewayV1beta1().ReferenceGrants(remote.Name).Create(ctx, grant, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(grant)

	t.Logf("validating that plugin %s is not configured with an incorrectly configured referencegrant", kongplugin.Name)
	assert.Never(t, func() bool {
		req := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.String(), "/test_plugin_reference?key=thirtytangas", nil)
		resp, err := helpers.DefaultHTTPClient(helpers.WithResolveHostTo(proxyHTTPURL.Host)).Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusTeapot
	}, negativeCheckWait, waitTick)

	t.Logf("creating a ReferenceGrant that permits kongconsumer access from %s to kongplugins in %s", remote.Name, ns.Name)
	grant = &gatewayapi.ReferenceGrant{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayapi.ReferenceGrantSpec{
			From: []gatewayapi.ReferenceGrantFrom{
				{
					Group:     gatewayapi.Group("configuration.konghq.com"),
					Kind:      gatewayapi.Kind("KongConsumer"),
					Namespace: gatewayapi.Namespace(remote.Name),
				},
			},
			To: []gatewayapi.ReferenceGrantTo{
				{
					Group: gatewayapi.Group("configuration.konghq.com"),
					Kind:  gatewayapi.Kind("KongPlugin"),
				},
			},
		},
	}
	_, err = gatewayClient.GatewayV1beta1().ReferenceGrants(ns.Name).Create(ctx, grant, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(grant)

	t.Logf("validating that plugin %s was successfully configured", kongplugin.Name)
	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		req := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.String(), "/test_plugin_reference?apikey=thirtytangas", nil)
		resp, err := helpers.DefaultHTTPClient(helpers.WithResolveHostTo(proxyHTTPURL.Host)).Do(req)
		if !assert.NoError(c, err) {
			return
		}
		defer resp.Body.Close()
		assert.True(c, resp.StatusCode == http.StatusTeapot)
	}, ingressWait, waitTick)
}
