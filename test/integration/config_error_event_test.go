//go:build integration_tests
// +build integration_tests

package integration

import (
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/runtime/schema"
	//"k8s.io/client-go/tools/events"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

func TestConfigErrorEventGeneration(t *testing.T) {
	// this test is NOT parallel. the broken configuration prevents all updates and will break unrelated tests
	// TODO maybe use the same separate test group as TestIngressRecoverFromInvalidPath
	// TODO skip on DB-backed
	ns, cleaner := setup(t)
	defer func() {
		if t.Failed() {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating an ingress for service %s with invalid configuration", service.Name)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/bar", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
		"konghq.com/protocols":      "grpcs",
		"konghq.com/methods":        "GET",
	}, service)

	t.Log("deploying ingress")
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	addIngressToCleaner(cleaner, ingress)

	t.Log("checking event creation")
	//selector, err := events.GetFieldSelector(
	//	schema.GroupVersion{Version: "v1"},
	//	schema.GroupVersionKind{
	//		Group:   "networking.k8s.io", // TODO not backwards compatible, but don't really care
	//		Version: "v1",                // TODO ditto
	//		Kind:    "Ingress",
	//	}, service.Name, "") // this relies on NewIngressForServiceWithClusterVersion's naming scheme, since runtime.Object has no name
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		//events, err := env.Cluster().Client().CoreV1().Events(ns.Name).List(ctx, metav1.ListOptions{FieldSelector: selector.String()})
		events, err := env.Cluster().Client().CoreV1().Events(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			return false
		}
		for _, event := range events.Items {
			if event.Reason == dataplane.KongConfigurationTranslationFailedEventReason {
				return true
			}
		}
		return false
	}, statusWait, waitTick, true)
}
