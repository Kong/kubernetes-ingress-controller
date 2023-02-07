//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

func TestConfigErrorEventGeneration(t *testing.T) {
	// this test is NOT parallel. the broken configuration prevents all updates and will break unrelated tests
	if testenv.DBMode() != "off" {
		t.Skip("config errors are only supported on DB-less mode")
	}
	if !versions.GetKongVersion().MajorMinorOnly().GTE(versions.FlattenedErrorCutoff) {
		t.Skip("flattened errors require Kong 3.2 or higher")
	} else {
		t.Logf("kong version is %s >= 3.2, testing config error parsing", versions.GetKongVersion().MajorMinorOnly().String())
	}
	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)
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
	service.ObjectMeta.Annotations["connect_timeout"] = "mankurt"
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating an ingress for service %s with invalid configuration", service.Name)
	kubernetesVersion, err := env.Cluster().Version()
	require.NoError(t, err)
	ingress := generators.NewIngressForServiceWithClusterVersion(kubernetesVersion, "/bar", map[string]string{
		annotations.IngressClassKey: consts.IngressClass,
		"konghq.com/strip-path":     "true",
		"konghq.com/protocols":      "grpcs",
		"konghq.com/methods":        "GET",
	}, service)

	t.Log("deploying ingress")
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	helpers.AddIngressToCleaner(cleaner, ingress)

	t.Log("checking event creation")

	// check broken route generates event
	require.Eventually(t, func() bool {
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

	// check broken service also generates event
	require.Eventually(t, func() bool {
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
