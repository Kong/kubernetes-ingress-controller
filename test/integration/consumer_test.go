//go:build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	"github.com/kong/kubernetes-configuration/pkg/clientset"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func TestConsumerCredential(t *testing.T) {
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
	ingress := generators.NewIngressForService("/test_consumer_credential", map[string]string{
		"konghq.com/strip-path": "true",
	}, service)
	ingress.Spec.IngressClassName = kong.String(consts.IngressClass)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Log("waiting for routes from Ingress to be operational")
	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_consumer_credential", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
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
			Name:      "basicbasic",
		},
		PluginName: "basic-auth",
	}
	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	kongplugin, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, kongplugin, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(kongplugin)

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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_consumer_credential", proxyHTTPURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyHTTPURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusUnauthorized
	}, ingressWait, waitTick)

	// we test basic-auth in particular because it's one of the plugins where the credentials are transformed and
	// stored with values different from the configuration provided (the password gets converted to its HTTP Basic Auth
	// form)
	credential := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Labels: map[string]string{
				labels.CredentialTypeLabel: "basic-auth",
			},
		},
		StringData: map[string]string{
			"username": "test_consumer_credential",
			"password": "test_consumer_credential",
		},
	}
	_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, credential, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(credential)

	consumer := &kongv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Username:    uuid.NewString(),
		Credentials: []string{credential.Name},
	}
	_, err = c.ConfigurationV1().KongConsumers(ns.Name).Create(ctx, consumer, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(consumer)

	t.Logf("validating that consumer has access")
	assert.Eventually(t, func() bool {
		req := helpers.MustHTTPRequest(t, http.MethodGet, proxyHTTPURL.Host, "/test_consumer_credential", nil)
		req.SetBasicAuth("test_consumer_credential", "test_consumer_credential")
		resp, err := helpers.DefaultHTTPClient(helpers.WithResolveHostTo(proxyHTTPURL.Host)).Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, ingressWait, waitTick)

	t.Logf("updating credential to confirm references update store copy")
	credential, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Get(ctx, credential.Name, metav1.GetOptions{})
	require.NoError(t, err)
	credential.StringData = map[string]string{"username": "new_consumer_credential"}
	_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Update(ctx, credential, metav1.UpdateOptions{})
	require.NoError(t, err)
	assert.Eventually(t, func() bool {
		req := helpers.MustHTTPRequest(t, "GET", proxyHTTPURL.String(), "/test_consumer_credential", nil)
		req.SetBasicAuth("new_consumer_credential", "test_consumer_credential")
		resp, err := helpers.DefaultHTTPClient().Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, ingressWait, waitTick)

	t.Log("deleting Ingress and waiting for routes to be torn down")
	require.NoError(t, clusters.DeleteIngress(ctx, env.Cluster(), ns.Name, ingress))
	helpers.EventuallyExpectHTTP404WithNoRoute(t, proxyHTTPURL, proxyHTTPURL.Host, "/test_plugin_essentials", ingressWait, waitTick, nil)
}
