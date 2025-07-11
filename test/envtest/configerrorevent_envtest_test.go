//go:build envtest

package envtest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"text/template"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	configurationv1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func TestConfigErrorEventGenerationInMemoryMode(t *testing.T) {
	// Can't be run in parallel because we're using t.Setenv() below which doesn't allow it.

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	scheme := Scheme(t, WithKong, WithGatewayAPI)

	restConfig := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, restConfig)

	ns := CreateNamespace(ctx, t, ctrlClient)
	ingressClassName := "kongenvtest"
	deployIngressClass(ctx, t, ingressClassName, ctrlClient)

	const podName = "kong-ingress-controller-tyjh1"
	t.Setenv("POD_NAMESPACE", ns.Name)
	t.Setenv("POD_NAME", podName)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Namespace = ns.Name
	require.NoError(t, ctrlClient.Create(ctx, deployment))

	t.Log("creating a KongUpstreamPolicy with sticky sessions configuration")
	upstreamPolicy := &configurationv1beta1.KongUpstreamPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "echo-drain-policy",
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClassName,
			},
		},
		Spec: configurationv1beta1.KongUpstreamPolicySpec{
			Algorithm: lo.ToPtr("sticky-sessions"),
			Slots:     lo.ToPtr(100),
			HashOn: &configurationv1beta1.KongUpstreamHash{
				Input: lo.ToPtr(configurationv1beta1.HashInput("none")),
			},
			StickySessions: &configurationv1beta1.KongUpstreamStickySessions{
				Cookie:     "session-id",
				CookiePath: lo.ToPtr("/"),
			},
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, upstreamPolicy))

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.Annotations = map[string]string{
		// TCP services cannot have paths, and we don't catch this as a translation error
		"konghq.com/protocol":        "tcp",
		"konghq.com/path":            "/aitmatov",
		"konghq.com/upstream-policy": upstreamPolicy.Name,
		// Referencing non-existent KongPlugins.
		"konghq.com/plugins": "foo,bar,n1:p1",
	}
	service.Namespace = ns.Name
	require.NoError(t, ctrlClient.Create(ctx, service))

	t.Logf("creating an ingress for service %s with invalid configuration", service.Name)
	// GRPC routes cannot have methods, only HTTP, and we don't catch this as a translation error
	ingress := generators.NewIngressForService("/bar", map[string]string{
		"konghq.com/strip-path": "true",
		"konghq.com/protocols":  "grpcs",
		"konghq.com/methods":    "GET",
		"konghq.com/plugins":    "baz",
	}, service)
	ingress.Spec.IngressClassName = lo.ToPtr(ingressClassName)
	ingress.Namespace = ns.Name
	t.Logf("deploying ingress %s", ingress.Name)
	require.NoError(t, ctrlClient.Create(ctx, ingress))

	RunManager(ctx, t, restConfig,
		AdminAPIOptFns(
			mocks.WithConfigPostError(formatErrBody(t, ns.Name, ingress, service)),
		),
		WithPublishService(ns.Name),
		WithIngressClass(ingressClassName),
		WithKongUpstreamPolicyEnabled(),
		WithProxySyncSeconds(0.1),
	)

	const numberOfExpectedEvents = 13
	collectedEvents := collectGeneratedEvents(
		ctx, t, ctrlClient, ns, t.Name(), numberOfExpectedEvents,
	)

	predicatesToCheck := []func(e corev1.Event) bool{
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationApplyFailedEventReason, "Ingress", ingress.Name, `^invalid methods: cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs'$`),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationApplyFailedEventReason, "Service", service.Name, `^invalid path: value must be null$`),
		predicate(corev1.EventTypeWarning, dataplane.FallbackKongConfigurationApplyFailedEventReason, "Ingress", ingress.Name, `^invalid methods: cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs'$`),
		predicate(corev1.EventTypeWarning, dataplane.FallbackKongConfigurationApplyFailedEventReason, "Service", service.Name, `^invalid path: value must be null$`),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationApplyFailedEventReason, "Service", service.Name, `^invalid service:httpbin\.httpbin\.80: failed conditional validation given value of field 'protocol'$`),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationApplyFailedEventReason, "Pod", podName, `failed to apply Kong configuration to http://[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+: HTTP status 400 \(message: "failed posting new config to /config"\)`),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationTranslationFailedEventReason, "Service", service.Name, `^referenced KongPlugin or KongClusterPlugin "foo" does not exist$`),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationTranslationFailedEventReason, "Service", service.Name, `^referenced KongPlugin or KongClusterPlugin "bar" does not exist$`),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationTranslationFailedEventReason, "Ingress", ingress.Name, `^referenced KongPlugin or KongClusterPlugin "baz" does not exist$`),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationTranslationFailedEventReason, "Service", service.Name, `^no grant found to referenced "n1:p1" plugin in the requested remote KongPlugin bind$`),
		predicate(corev1.EventTypeWarning, dataplane.FallbackKongConfigurationApplyFailedEventReason, "Service", service.Name, `^invalid service:httpbin\.httpbin\.80: failed conditional validation given value of field 'protocol'$`),
		predicate(corev1.EventTypeWarning, dataplane.FallbackKongConfigurationApplyFailedEventReason, "Pod", podName, `failed to apply fallback Kong configuration to http://[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+: HTTP status 400 \(message: "failed posting new config to /config"\)`),
		// TODO: Remove this event once we start using Kong Gateway >= 3.11.0, because sticky sessions type is supported there.
		// Adjust also numberOfExpectedEvents constant above.
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationTranslationFailedEventReason, "Service", service.Name, `^sticky sessions algorithm specified in KongUpstreamPolicy 'echo-drain-policy' is not supported with Kong Gateway versions < 3\.11\.0$`),
	}

	assertExpectedEvents(t, predicatesToCheck, collectedEvents)
}

func TestConfigErrorEventGenerationDBMode(t *testing.T) {
	// Can't be run in parallel because we're using t.Setenv() below which doesn't allow it.

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	scheme := Scheme(t, WithKong)
	restConfig := Setup(t, scheme)
	ctrlClientGlobal := NewControllerClient(t, scheme, restConfig)
	ns := CreateNamespace(ctx, t, ctrlClientGlobal)
	ctrlClient := client.NewNamespacedClient(ctrlClientGlobal, ns.Name)

	ingressClassName := "kongenvtest"
	deployIngressClass(ctx, t, ingressClassName, ctrlClient)

	const podName = "kong-ingress-controller-tyjh1"
	t.Setenv("POD_NAMESPACE", ns.Name)
	t.Setenv("POD_NAME", podName)

	t.Logf("creating a static consumer in %s namespace which will be used to test global validation", ns.Name)
	consumer := &configurationv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name: "donenbai",
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClassName,
				// Referencing non-existent KongPlugin.
				"konghq.com/plugins": "foo, n1:p1",
			},
		},
		Username: "donenbai",
	}
	require.NoError(t, ctrlClient.Create(ctx, consumer))
	t.Cleanup(func() {
		if err := ctrlClient.Delete(ctx, consumer); err != nil && !apierrors.IsNotFound(err) && !errors.Is(err, context.Canceled) {
			assert.NoError(t, err)
		}
	})

	RunManager(ctx, t, restConfig,
		AdminAPIOptFns(
			mocks.WithRoot(formatDBRootResponse("999.999.999")),
		),
		WithPublishService(ns.Name),
		WithIngressClass(ingressClassName),
		WithProxySyncSeconds(0.1),
	)

	const numberOfExpectedEvents = 6
	collectedEvents := collectGeneratedEvents(
		ctx, t, ctrlClient, ns, t.Name(), numberOfExpectedEvents,
	)

	predicatesToCheck := []func(e corev1.Event) bool{
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationApplyFailedEventReason, "KongConsumer", consumer.Name, fmt.Sprintf(`^invalid consumer:%s: HTTP status 400 \(message: "2 schema violations \(at least one of these fields must be non-empty: 'custom_id', 'username'; fake: unknown field\)"\)$`, consumer.Name)),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationTranslationFailedEventReason, "KongConsumer", consumer.Name, `^referenced KongPlugin or KongClusterPlugin "foo" does not exist$`),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationTranslationFailedEventReason, "KongConsumer", consumer.Name, `^no grant found to referenced "n1:p1" plugin in the requested remote KongPlugin bind$`),
		predicate(corev1.EventTypeNormal, dataplane.KongConfigurationApplySucceededEventReason, "Pod", podName, `successfully applied Kong configuration to http://[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+`),
		predicate(corev1.EventTypeNormal, dataplane.FallbackKongConfigurationApplySucceededEventReason, "Pod", podName, `successfully applied fallback Kong configuration to http://[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+`),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationApplyFailedEventReason, "Pod", podName, `failed to apply Kong configuration to http://[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+: 1 errors occurred:\s+while processing event: Create consumer donenbai failed: HTTP status 400 \(message: "2 schema violations \(at least one of these fields must be non-empty: 'custom_id', 'username'; fake: unknown field\)\"\)`),
	}
	// Check that all expected events are present
	assertExpectedEvents(t, predicatesToCheck, collectedEvents)
}

func TestStickySessionsNotSupportedEventGeneration(t *testing.T) {
	// Can't be run in parallel because we're using t.Setenv() below which doesn't allow it.

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	scheme := Scheme(t, WithKong)
	restConfig := Setup(t, scheme)
	ctrlClientGlobal := NewControllerClient(t, scheme, restConfig)
	ns := CreateNamespace(ctx, t, ctrlClientGlobal)
	ctrlClient := client.NewNamespacedClient(ctrlClientGlobal, ns.Name)

	ingressClassName := "kongenvtest"
	deployIngressClass(ctx, t, ingressClassName, ctrlClient)

	const podName = "kong-ingress-controller-tyjh1"
	t.Setenv("POD_NAMESPACE", ns.Name)
	t.Setenv("POD_NAME", podName)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Namespace = ns.Name
	require.NoError(t, ctrlClient.Create(ctx, deployment))

	t.Log("creating a KongUpstreamPolicy with sticky sessions configuration")
	upstreamPolicy := &configurationv1beta1.KongUpstreamPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "echo-drain-policy",
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClassName,
			},
		},
		Spec: configurationv1beta1.KongUpstreamPolicySpec{
			Algorithm: lo.ToPtr("sticky-sessions"),
			Slots:     lo.ToPtr(100),
			HashOn: &configurationv1beta1.KongUpstreamHash{
				Input: lo.ToPtr(configurationv1beta1.HashInput("none")),
			},
			StickySessions: &configurationv1beta1.KongUpstreamStickySessions{
				Cookie:     "session-id",
				CookiePath: lo.ToPtr("/"),
			},
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, upstreamPolicy))

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.Annotations = map[string]string{
		"konghq.com/upstream-policy": upstreamPolicy.Name,
	}
	service.Namespace = ns.Name
	require.NoError(t, ctrlClient.Create(ctx, service))

	t.Logf("creating an ingress for service %s with invalid configuration", service.Name)
	ingress := generators.NewIngressForService("/bar", nil, service)
	ingress.Spec.IngressClassName = lo.ToPtr(ingressClassName)
	ingress.Namespace = ns.Name
	t.Logf("deploying ingress %s", ingress.Name)
	require.NoError(t, ctrlClient.Create(ctx, ingress))

	kongContainer := runKongGatewayWithoutStickySessionsSupport(ctx, t)
	RunManager(ctx, t, restConfig,
		AdminAPIOptFns(),
		WithPublishService(ns.Name),
		WithIngressClass(ingressClassName),
		WithKongServiceFacadeFeatureEnabled(),
		WithProxySyncSeconds(0.1),
		WithKongAdminURLs(kongContainer.AdminURL(ctx, t)),
	)

	const numberOfExpectedEvents = 2
	collectedEvents := collectGeneratedEvents(
		ctx, t, ctrlClient, ns, t.Name(), numberOfExpectedEvents,
	)

	predicatesToCheck := []func(e corev1.Event) bool{
		predicate(corev1.EventTypeNormal, dataplane.KongConfigurationApplySucceededEventReason, "Pod", podName, `successfully applied Kong configuration to http://([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+|localhost:[0-9]+)`),
		predicate(corev1.EventTypeWarning, dataplane.KongConfigurationTranslationFailedEventReason, "Service", service.Name, `^sticky sessions algorithm specified in KongUpstreamPolicy 'echo-drain-policy' is not supported with Kong Gateway versions < 3\.11\.0$`),
	}

	assertExpectedEvents(t, predicatesToCheck, collectedEvents)
}

func predicate(eventType, eventReason, invObjKind, invObjName, msgToMatch string) func(e corev1.Event) bool {
	return func(e corev1.Event) bool {
		ok, err := regexp.MatchString(msgToMatch, e.Message)
		return e.Type == eventType &&
			e.Reason == eventReason &&
			e.InvolvedObject.Kind == invObjKind &&
			e.InvolvedObject.Name == invObjName &&
			ok && err == nil
	}
}

func collectGeneratedEvents(
	ctx context.Context, t *testing.T, ctrlClient client.Client, ns corev1.Namespace, expectedInstanceID string, numberOfExpectedEvents int,
) []corev1.Event {
	t.Helper()
	t.Log("checking for events generated by the controller")
	const (
		waitTime = time.Minute
		tickTime = 100 * time.Millisecond
	)
	var collectedEvents []corev1.Event
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		var events corev1.EventList
		// Filter out events that are not related to the current test instance.
		require.NoError(c, ctrlClient.List(ctx, &events, client.InNamespace(ns.Name)))
		collectedEvents = lo.Filter(events.Items, func(e corev1.Event, _ int) bool {
			return e.Annotations[consts.InstanceIDAnnotationKey] == expectedInstanceID
		})
		require.Len(c, collectedEvents, numberOfExpectedEvents, "number of events mismatch")
	}, waitTime, tickTime)
	return collectedEvents
}

func assertExpectedEvents(t *testing.T, predicatesToCheck []func(e corev1.Event) bool, collectedEvents []corev1.Event) {
	t.Helper()
	for pi, predicate := range predicatesToCheck {
		lenBefore := len(collectedEvents)
		collectedEvents = lo.Reject(collectedEvents, func(e corev1.Event, _ int) bool {
			return predicate(e)
		})
		lenAfter := len(collectedEvents)
		if !assert.Equalf(t, lenBefore-1, lenAfter, "expected one event to be removed, but predicate with index: %d doesn't do it", pi) {
			break
		}
	}
	if !assert.Equal(t, 0, len(collectedEvents), "expected all warning events to match test predicates, but some were left") {
		t.Logf("remaining events %d:", len(collectedEvents))
		for _, e := range collectedEvents {
			t.Logf("  - %s | %s | %s | %s | %s", e.Type, e.Reason, e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Message)
		}
	}
}

func formatErrBody(t *testing.T, namespace string, ingress *netv1.Ingress, service *corev1.Service) []byte {
	t.Helper()

	const errBody = `{
	"code": 14,
	"name": "invalid declarative configuration",
	"flattened_errors": [
		{
			"entity_name": "{{ .Ingress.Name }}.httpbin.httpbin..80",
			"entity_tags": [
				"k8s-name:httpbin",
				"k8s-namespace:{{ .Namespace }}",
				"k8s-kind:Ingress",
				"k8s-uid:{{ .Ingress.UID }}",
				"k8s-group:networking.k8s.io",
				"k8s-version:v1"
			],
			"errors": [
				{
					"field": "methods",
					"type": "field",
					"message": "cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs'"
				}
			],
			"entity": {
				"regex_priority": 0,
				"preserve_host": true,
				"name": "{{ .Ingress.Name }}.httpbin.httpbin..80",
				"protocols": [
					"grpcs"
				],
				"https_redirect_status_code": 426,
				"request_buffering": true,
				"tags": [
					"k8s-name:httpbin",
					"k8s-namespace:{{ .Namespace }}",
					"k8s-kind:Ingress",
					"k8s-uid:{{ .Ingress.UID }}",
					"k8s-group:networking.k8s.io",
					"k8s-version:v1"
				],
				"path_handling": "v0",
				"response_buffering": true,
				"methods": [
					"GET"
				],
				"paths": [
					"/bar/",
					"~/bar$"
				]
			},
			"entity_type": "route"
		},
		{
			"entity_name": "{{ .Ingress.Name }}.httpbin.80",
			"entity_tags": [
				"k8s-name:httpbin",
				"k8s-namespace:{{ .Namespace }}",
				"k8s-kind:Service",
				"k8s-uid:{{ .Service.UID }}",
				"k8s-version:v1"
			],
			"errors": [
				{
					"field": "path",
					"type": "field",
					"message": "value must be null"
				},
				{
					"type": "entity",
					"message": "failed conditional validation given value of field 'protocol'"
				}
			],
			"entity": {
				"read_timeout": 60000,
				"path": "/aitmatov",
				"write_timeout": 60000,
				"protocol": "tcp",
				"tags": [
					"k8s-name:httpbin",
					"k8s-namespace:{{ .Namespace }}",
					"k8s-kind:Service",
					"k8s-uid:{{ .Service.UID }}",
					"k8s-version:v1"
				],
				"retries": 5,
				"port": 80,
				"name": "{{ .Ingress.Name }}.httpbin.80",
				"host": "httpbin.{{ .Ingress.Name }}.80.svc",
				"connect_timeout": 60000
			},
			"entity_type": "service"
		}
	],
	"message": "declarative config is invalid: {}",
	"fields": {}
}`
	tmpl, err := template.New("body").Parse(errBody)
	require.NoError(t, err)

	type ErrBody struct {
		Namespace string
		Ingress   *netv1.Ingress
		Service   *corev1.Service
	}

	var b bytes.Buffer
	require.NoError(t, tmpl.Execute(&b, ErrBody{
		Namespace: namespace,
		Ingress:   ingress,
		Service:   service,
	}))

	return b.Bytes()
}

func formatDBRootResponse(version string) []byte {
	const defaultDBLessRootResponse = `{
		"version": "%s",
		"configuration": {
			"database": "postgres",
			"router_flavor": "traditional",
			"role": "traditional",
			"proxy_listeners": [
				{
					"ipv6only=on": false,
					"ipv6only=off": false,
					"ssl": false,
					"so_keepalive=off": false,
					"listener": "0.0.0.0:8000",
					"bind": false,
					"port": 8000,
					"deferred": false,
					"so_keepalive=on": false,
					"http2": false,
					"proxy_protocol": false,
					"ip": "0.0.0.0",
					"reuseport": false
				}
			]
		}
	}`
	return fmt.Appendf(nil, defaultDBLessRootResponse, version)
}
