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
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func TestConfigErrorEventGenerationInMemoryMode(t *testing.T) {
	// Can't be run in parallel because we're using t.Setenv() below which doesn't allow it.

	const (
		waitTime = time.Minute
		tickTime = 100 * time.Millisecond
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	restConfig := Setup(t, scheme.Scheme)
	ctrlClient := NewControllerClient(t, scheme.Scheme, restConfig)

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

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.ObjectMeta.Annotations = map[string]string{
		// TCP services cannot have paths, and we don't catch this as a translation error
		"konghq.com/protocol": "tcp",
		"konghq.com/path":     "/aitmatov",
	}
	service.Namespace = ns.Name
	require.NoError(t, ctrlClient.Create(ctx, service))

	t.Logf("creating an ingress for service %s with invalid configuration", service.Name)
	// GRPC routes cannot have methods, only HTTP, and we don't catch this as a translation error
	ingress := generators.NewIngressForService("/bar", map[string]string{
		"konghq.com/strip-path": "true",
		"konghq.com/protocols":  "grpcs",
		"konghq.com/methods":    "GET",
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
		WithProxySyncSeconds(0.1),
	)

	t.Log("checking ingress and service event creation")
	require.Eventually(t, func() bool {
		var events corev1.EventList
		if err := ctrlClient.List(ctx, &events, &client.ListOptions{Namespace: ns.Name}); err != nil {
			t.Logf("error listing events: %v", err)
			return false
		}
		t.Logf("got %d events", len(events.Items))

		matches := make([]bool, 4)
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
			return e.Reason == dataplane.KongConfigurationApplyFailedEventReason &&
				e.InvolvedObject.Kind == "Service" &&
				e.InvolvedObject.Name == service.Name &&
				e.Message == "invalid service:httpbin.httpbin.80: failed conditional validation given value of field 'protocol'"
		})
		matches[3] = lo.ContainsBy(events.Items, func(e corev1.Event) bool {
			ok, err := regexp.MatchString(`failed to apply Kong configuration to http://[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+: HTTP status 400 \(message: "failed posting new config to /config"\)`, e.Message)
			return e.Reason == dataplane.KongConfigurationApplyFailedEventReason &&
				e.InvolvedObject.Kind == "Pod" &&
				e.InvolvedObject.Name == podName &&
				ok && err == nil
		})
		if lo.Count(matches, true) != 4 {
			t.Logf("not all events matched: %+v", matches)
			return false
		}
		return true
	}, waitTime, tickTime)

	t.Log("push failure events recorded successfully")
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

func TestConfigErrorEventGenerationDBMode(t *testing.T) {
	// Can't be run in parallel because we're using t.Setenv() below which doesn't allow it.

	const (
		waitTime = time.Minute
		tickTime = 100 * time.Millisecond
	)

	ctx, cancel := context.WithCancel(context.Background())
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
	consumer := &kongv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name: "donenbai",
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClassName,
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
			// TODO IDK where we're getting the version from normally but it shouldn't really matter for this.
			mocks.WithRoot(formatDBRootResponse("999.999.999")),
		),
		WithPublishService(ns.Name),
		WithIngressClass(ingressClassName),
		WithProxySyncSeconds(0.1),
	)

	t.Log("checking kongconsumer event creation")
	require.Eventually(t, func() bool {
		var events corev1.EventList
		if err := ctrlClient.List(ctx, &events, &client.ListOptions{Namespace: ns.Name}); err != nil {
			t.Logf("error listing events: %v", err)
			return false
		}
		t.Logf("got %d events", len(events.Items))

		matches := make([]bool, 1)
		matches[0] = lo.ContainsBy(events.Items, func(e corev1.Event) bool {
			return e.Reason == dataplane.KongConfigurationApplyFailedEventReason &&
				e.InvolvedObject.Kind == "KongConsumer" &&
				e.InvolvedObject.Name == consumer.Name &&
				e.Message == fmt.Sprintf("invalid consumer:%s: HTTP status 400 (message: \"2 schema violations (at least one of these fields must be non-empty: 'custom_id', 'username'; fake: unknown field)\")", consumer.Name)
		})
		if lo.Count(matches, true) != 1 {
			t.Logf("not all events matched: %+v", matches)
			return false
		}
		return true
	}, waitTime, tickTime)

	t.Log("push failure events recorded successfully")
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
	return []byte(fmt.Sprintf(defaultDBLessRootResponse, version))
}
