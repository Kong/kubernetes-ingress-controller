//go:build integration_tests

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConsumersConflict(t *testing.T) {
	ctx := context.Background()

	t.Parallel()
	ns, cleaner := helpers.Setup(ctx, t, env)

	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	plugin1 := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "plugin-1",
		},
		InstanceName: "example",
		PluginName:   "request-termination",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"status_code": 418}`),
		},
	}
	plugin1, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, plugin1, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(plugin1)

	// Create a Plugin using the same PluginName but different Name.
	plugin2 := plugin1.DeepCopy()
	plugin2.Name = "plugin-2"
	plugin2, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, plugin2, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(plugin2)

	// Create an Ingress using two Plugins using the same PluginName.
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
	ingress := generators.NewIngressForService("/path", map[string]string{
		annotations.IngressClassKey: consts.IngressClass,
		"konghq.com/strip-path":     "true",
		"konghq.com/plugins":        plugin1.Name + "," + plugin2.Name,
	}, service)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	time.Sleep(10 * time.Minute)
}
