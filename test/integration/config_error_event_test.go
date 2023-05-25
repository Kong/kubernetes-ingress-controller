//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"fmt"
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
)

func TestConfigErrorEventGeneration(t *testing.T) {
	// This test is NOT parallel.
	// The broken configuration prevents all updates and will break unrelated tests

	skipTestForExpressionRouter(t)

	RunWhenKongDBMode(t, "off", "config errors are only supported on DB-less mode")
	RunWhenKongVersion(t, fmt.Sprintf(">=%s", versions.FlattenedErrorCutoff))

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.ObjectMeta.Annotations = map[string]string{}
	// TCP services cannot have paths, and we don't catch this as a translation error
	service.ObjectMeta.Annotations["konghq.com/protocol"] = "tcp"
	service.ObjectMeta.Annotations["konghq.com/path"] = "/aitmatov"
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating an ingress for service %s with invalid configuration", service.Name)
	// GRPC routes cannot have methods, only HTTP, and we don't catch this as a translation error
	ingress := generators.NewIngressForService("/bar", map[string]string{
		annotations.IngressClassKey: consts.IngressClass,
		"konghq.com/strip-path":     "true",
		"konghq.com/protocols":      "grpcs",
		"konghq.com/methods":        "GET",
	}, service)

	t.Log("deploying ingress")
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Log("checking ingress event creation")
	require.Eventually(t, func() bool {
		events, err := env.Cluster().Client().CoreV1().Events(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			return false
		}
		for _, event := range events.Items {
			if event.Reason == dataplane.KongConfigurationApplyFailedEventReason {
				if event.InvolvedObject.Kind == "Ingress" {
					// this is a runtime.Object because of v1/v1beta1 handling, so no ObjectMeta or other obvious way
					// to get the name. we can reasonably assume it's the only Ingress in the namespace
					return true
				}
			}
		}
		return false
	}, statusWait, waitTick, true)

	t.Log("checking service event creation")
	require.Eventually(t, func() bool {
		events, err := env.Cluster().Client().CoreV1().Events(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			return false
		}
		for _, event := range events.Items {
			if event.Reason == dataplane.KongConfigurationApplyFailedEventReason {
				if event.InvolvedObject.Kind == "Service" {
					if event.InvolvedObject.Name == service.ObjectMeta.Name {
						return true
					}
				}
			}
		}
		return false
	}, statusWait, waitTick, true)
	t.Log("push failure events recorded successfully")
}
