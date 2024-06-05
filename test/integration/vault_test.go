//go:build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func TestCustomVault(t *testing.T) {
	t.Parallel()

	RunWhenKongEnterprise(t)
	// TODO: run hcv vault to enable test with DBMode
	RunWhenKongDBMode(t, dpconf.DBModeOff, "Skipping because DBMode cannot support env vault")

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
	ingress := generators.NewIngressForService("/test_custom_vault", map[string]string{
		"konghq.com/strip-path": "true",
	}, service)
	ingress.Spec.IngressClassName = kong.String(consts.IngressClass)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Log("waiting for routes from Ingress to be operational")
	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_custom_vault", proxyHTTPURL))
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

	t.Logf("creating a Kong vault using env backend")

	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	_, err = c.ConfigurationV1alpha1().KongVaults().Create(ctx, &kongv1alpha1.KongVault{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-env-vault",
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1alpha1.KongVaultSpec{
			Backend:     "env",
			Prefix:      "test-env",
			Description: "env vault for test",
			Config: apiextensionsv1.JSON{
				Raw: []byte(`{"prefix":"kong_vault_test_"}`),
			},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("create a request-transformer-advanced plugin referencing the value from the vault")
	_, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "request-transformer-advanced",
		},
		PluginName: "request-transformer-advanced",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"add":{"headers":["{vault://test-env/add-header-1}"]}}`),
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("attach plugin to ingress and check if the config from vault takes effect")
	ingressName := ingress.Name
	require.Eventually(t, func() bool {
		ingClient := env.Cluster().Client().NetworkingV1().Ingresses(ns.Name)
		ingress, err = ingClient.Get(ctx, ingressName, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting %s: %v", client.ObjectKeyFromObject(ingress), err)
			return false
		}
		ingress.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = "request-transformer-advanced"
		_, err = ingClient.Update(ctx, ingress, metav1.UpdateOptions{})
		if err != nil {
			t.Logf("error annotating %s: %v", client.ObjectKeyFromObject(ingress), err)
			return false
		}
		return true
	}, ingressWait, waitTick)

	require.Eventuallyf(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("%s/test_custom_vault/headers", proxyHTTPURL))
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
			return strings.Contains(b.String(), `"H1": "v1"`)
		}
		return false
	},
		ingressWait, waitTick,
		"Cannot find added header in request",
	)
}
