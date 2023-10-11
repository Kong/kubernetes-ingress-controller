//go:build integration_tests

package integration

import (
	"context"
	"regexp"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

// NOTE: there is an equivalent envtest based test for this but the idea is to
// keep that one and remove this one here when we introduce a new test that would
//   - test event generation for broken configurations
//   - not rely solely on mocked Admin API to prevent drifting away from changes
//     in Kong Gateway
//
// Hence when an appropriate test for the above is introduce this one below could
// be removed.
func TestConfigErrorEventGeneration(t *testing.T) {
	// This test is NOT parallel.
	// The broken configuration prevents all updates and will break unrelated tests

	skipTestForExpressionRouter(t)
	RunWhenKongDBMode(t, "off", "config errors are only supported on DB-less mode")

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	ctrlClient, err := ctrlclient.New(env.Cluster().Config(), ctrlclient.Options{})
	require.NoError(t, err)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Namespace = ns.Name
	require.NoError(t, ctrlClient.Create(ctx, deployment))
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.ObjectMeta.Annotations = map[string]string{
		// TCP services cannot have paths, and we don't catch this as a translation error
		"konghq.com/protocol": "tcp",
		"konghq.com/path":     "/aitmatov",
	}
	service.Namespace = ns.Name
	require.NoError(t, ctrlClient.Create(ctx, service))
	cleaner.Add(service)

	t.Logf("creating an ingress for service %s with invalid configuration", service.Name)
	// GRPC routes cannot have methods, only HTTP, and we don't catch this as a translation error
	ingress := generators.NewIngressForService("/bar", map[string]string{
		"konghq.com/strip-path": "true",
		"konghq.com/protocols":  "grpcs",
		"konghq.com/methods":    "GET",
	}, service)
	ingress.Spec.IngressClassName = kong.String(consts.IngressClass)
	ingress.Namespace = ns.Name
	t.Logf("deploying ingress %s", ingress.Name)
	require.NoError(t, ctrlClient.Create(ctx, ingress))
	cleaner.Add(ingress)

	t.Log("checking ingress and service event creation")
	require.Eventually(t, func() bool {
		var events corev1.EventList
		if err := ctrlClient.List(ctx, &events, &ctrlclient.ListOptions{Namespace: ns.Name}); err != nil {
			t.Logf("error listing events: %v", err)
			return false
		}
		t.Logf("got %d events", len(events.Items))

		matches := make([]bool, 3)
		matches[0] = lo.ContainsBy(events.Items, func(e corev1.Event) bool {
			return e.Reason == dataplane.KongConfigurationApplyFailedEventReason &&
				e.InvolvedObject.Kind == "Ingress" &&
				e.InvolvedObject.Name == ingress.Name &&
				e.Message == "invalid methods: cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs'"
		})
		matches[1] = lo.ContainsBy(events.Items, func(e corev1.Event) bool {
			return e.Reason == dataplane.KongConfigurationApplyFailedEventReason &&
				e.InvolvedObject.Kind == "Service" &&
				e.InvolvedObject.Name == service.Name &&
				e.Message == "invalid path: value must be null"
		})
		matches[2] = lo.ContainsBy(events.Items, func(e corev1.Event) bool {
			ok, err := regexp.MatchString(`invalid service:.+\.httpbin\.80: failed conditional validation given value of field 'protocol'`, e.Message)
			return e.Reason == dataplane.KongConfigurationApplyFailedEventReason &&
				e.InvolvedObject.Kind == "Service" &&
				e.InvolvedObject.Name == service.Name &&
				ok && err == nil
		})

		if lo.Count(matches, true) != 3 {
			t.Logf("not all events matched: %+v", matches)
			return false
		}

		return true
	}, statusWait, waitTick)

	t.Log("push failure events recorded successfully")
}
