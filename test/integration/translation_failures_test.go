//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	"github.com/kong/kubernetes-configuration/pkg/clientset"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	testutils "github.com/kong/kubernetes-ingress-controller/v3/internal/util/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

const testTranslationFailuresObjectsPrefix = "translation-failures-"

// TestTranslationFailures ensures that proper warning Kubernetes events are recorded in case of translation failures
// encountered.
func TestTranslationFailures(t *testing.T) {
	ctx := context.Background()

	type expectedTranslationFailure struct {
		causingObjects []client.Object
		reasonContains string
	}

	testCases := []struct {
		name string
		// translationFailureTrigger should create objects that trigger translation failure and return the objects
		// that we expect translation failure warning events to be created for.
		translationFailureTrigger func(t *testing.T, cleaner *clusters.Cleaner, ns string) expectedTranslationFailure
	}{
		{
			name: "invalid CA secret",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) expectedTranslationFailure {
				createdSecret, err := env.Cluster().Client().CoreV1().Secrets(ns).Create(ctx, invalidCASecret(ns), metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(createdSecret)

				return expectedTranslationFailure{
					causingObjects: []client.Object{createdSecret},
					reasonContains: "invalid CA certificate: missing 'cert' field in data",
				}
			},
		},
		{
			name: "invalid CA secret referred by a plugin",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) expectedTranslationFailure {
				createdSecret, err := env.Cluster().Client().CoreV1().Secrets(ns).Create(ctx, invalidCASecret(ns), metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(createdSecret)

				c, err := clientset.NewForConfig(env.Cluster().Config())
				require.NoError(t, err)
				createdPlugin, err := c.ConfigurationV1().KongPlugins(ns).Create(ctx, pluginUsingInvalidCACert(ns), metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(createdPlugin)

				return expectedTranslationFailure{
					// expect events for both: a faulty secret and a plugin referring it
					causingObjects: []client.Object{createdSecret, createdPlugin},
					reasonContains: "invalid CA certificate: missing 'cert' field in data",
				}
			},
		},
		{
			name: "grouped services annotations do not match",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) expectedTranslationFailure {
				gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
				require.NoError(t, err)

				gatewayClassName := testutils.RandomName(testTranslationFailuresObjectsPrefix)
				gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
				require.NoError(t, err)
				cleaner.Add(gwc)

				gatewayName := testutils.RandomName(testTranslationFailuresObjectsPrefix)
				gateway, err := helpers.DeployGateway(ctx, gatewayClient, ns, gatewayClassName, func(gw *gatewayapi.Gateway) {
					gw.Name = gatewayName
				})
				require.NoError(t, err)
				cleaner.Add(gateway)

				container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
				deployment := generators.NewDeploymentForContainer(container)
				deployment, err = env.Cluster().Client().AppsV1().Deployments(ns).Create(ctx, deployment, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(deployment)

				service1 := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
				service1.Name = testutils.RandomName(testTranslationFailuresObjectsPrefix)
				// adding the annotation to trigger conflict
				service1.Annotations = map[string]string{annotations.AnnotationPrefix + annotations.HostHeaderKey: "example.com"}
				service1, err = env.Cluster().Client().CoreV1().Services(ns).Create(ctx, service1, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(service1)

				service2 := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
				service2.Name = testutils.RandomName(testTranslationFailuresObjectsPrefix)
				service2, err = env.Cluster().Client().CoreV1().Services(ns).Create(ctx, service2, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(service2)

				httpRoute := httpRouteWithBackends(gatewayName, service1, service2)
				httpRoute, err = gatewayClient.GatewayV1().HTTPRoutes(ns).Create(ctx, httpRoute, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(httpRoute)

				return expectedTranslationFailure{
					// expect event for service2 as it doesn't have annotations that service1 has
					causingObjects: []client.Object{service2},
					reasonContains: "This annotation must have the same value across all Services in the backend.",
				}
			},
		},
		{
			name: "missing client-cert for service",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) expectedTranslationFailure {
				service := validService()
				service.ObjectMeta.Annotations = map[string]string{
					"konghq.com/client-cert": "not-existing-secret",
				}
				service, err := env.Cluster().Client().CoreV1().Services(ns).Create(ctx, service, metav1.CreateOptions{})
				require.NoError(t, err)

				ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns).Create(
					ctx,
					ingressWithPathBackedByService(service),
					metav1.CreateOptions{},
				)
				require.NoError(t, err)
				cleaner.Add(ingress)

				return expectedTranslationFailure{
					causingObjects: []client.Object{service},
					reasonContains: "Failed to fetch secret",
				}
			},
		},
		{
			name: "missing ingress backing service",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) expectedTranslationFailure {
				nonExistingService := &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "non-existing-service"}}
				ingress := ingressWithPathBackedByService(nonExistingService)
				ingress, err := env.Cluster().Client().NetworkingV1().Ingresses(ns).Create(ctx, ingress, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(ingress)

				return expectedTranslationFailure{
					causingObjects: []client.Object{ingress},
					reasonContains: "can't add target for backend non-existing-service: no kubernetes service found",
				}
			},
		},
		{
			name: "missing port for service",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) expectedTranslationFailure {
				service, err := env.Cluster().Client().CoreV1().Services(ns).Create(ctx, validService(), metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(service)

				ingress := ingressWithPathBackedByService(service)
				const notMatchingPort = 90
				ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Backend.Service.Port = netv1.ServiceBackendPort{
					Number: notMatchingPort,
				}
				ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(ns).Create(ctx, ingress, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(ingress)

				return expectedTranslationFailure{
					causingObjects: []client.Object{ingress, service},
					reasonContains: "can't find port for backend kubernetes service: no suitable port found",
				}
			},
		},
		{
			name: "ingress referring a non-existing TLS secret",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) expectedTranslationFailure {
				service, err := env.Cluster().Client().CoreV1().Services(ns).Create(ctx, validService(), metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(service)

				ingress := ingressWithPathBackedByService(service)
				ingress.Spec.TLS = []netv1.IngressTLS{
					{
						Hosts:      []string{"example.com"},
						SecretName: "non-existing-secret",
					},
				}
				ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(ns).Create(ctx, ingress, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(ingress)

				return expectedTranslationFailure{
					causingObjects: []client.Object{ingress},
					reasonContains: "Failed to fetch the secret",
				}
			},
		},
		{
			name: "ingress referring a secret with no valid TLS key-pair",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) expectedTranslationFailure {
				secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: testutils.RandomName(testTranslationFailuresObjectsPrefix)}}
				secret, err := env.Cluster().Client().CoreV1().Secrets(ns).Create(ctx, secret, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(secret)

				service, err := env.Cluster().Client().CoreV1().Services(ns).Create(ctx, validService(), metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(service)

				ingress := ingressWithPathBackedByService(service)
				ingress.Spec.TLS = []netv1.IngressTLS{
					{
						Hosts:      []string{"example.com"},
						SecretName: secret.Name,
					},
				}
				ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(ns).Create(ctx, ingress, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(ingress)

				return expectedTranslationFailure{
					causingObjects: []client.Object{ingress, secret},
					reasonContains: "failed to construct certificate from secret",
				}
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ns, cleaner := helpers.Setup(ctx, t, env)

			expected := tt.translationFailureTrigger(t, cleaner, ns.GetName())

			require.Eventually(t, func() bool {
				eventsForAllObjectsFound := true
				var receivedEvents []corev1.Event

				for _, expectedCausingObject := range expected.causingObjects {
					events, eventForCurrentObjectFound := findTranslationFailureEventsForObject(
						ctx,
						t,
						ns,
						expectedCausingObject,
						expected.reasonContains,
					)

					receivedEvents = append(receivedEvents, events...)
					if !eventForCurrentObjectFound {
						eventsForAllObjectsFound = false
					}

				}

				logReceivedEvents(t, receivedEvents, eventsForAllObjectsFound)
				return eventsForAllObjectsFound
			}, time.Minute*5, time.Second)
		})
	}
}

// findTranslationFailureEventsForObject finds out ALL warning events related to the input object
// and returns true if any of the event's message contains the expected reason.
func findTranslationFailureEventsForObject(
	ctx context.Context,
	t *testing.T,
	ns *corev1.Namespace,
	obj client.Object,
	expectedReasonContains string,
) ([]corev1.Event, bool) {
	relatedEvents, err := env.Cluster().Client().CoreV1().Events(ns.GetName()).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf(
			"reason=%s,type=%s,involvedObject.name=%s",
			dataplane.KongConfigurationTranslationFailedEventReason,
			corev1.EventTypeWarning,
			obj.GetName(),
		),
	})
	if err != nil {
		t.Logf("Failed to list events: %v", err)
		return []corev1.Event{}, false
	}
	if len(relatedEvents.Items) == 0 {
		t.Logf("waiting for events related to %q to be created", obj.GetName())
		return []corev1.Event{}, false
	}
	for _, event := range relatedEvents.Items {
		actualMsg := event.Message
		if strings.Contains(actualMsg, expectedReasonContains) {
			t.Logf("received event's message (%q) contains the expected reason: %q", actualMsg, expectedReasonContains)
			return relatedEvents.Items, true
		}
		t.Logf("received event's message (%q) does not contain the expected reason: %q", actualMsg, expectedReasonContains)
	}
	return relatedEvents.Items, false
}

func logReceivedEvents(t *testing.T, events []corev1.Event, eventsForAllObjectsFound bool) {
	eventsString := eventsToString(events)
	if eventsForAllObjectsFound {
		t.Logf("received all events:\n%s", eventsString)
	} else {
		t.Logf("waiting for events, received so far:\n%s", eventsString)
	}
}

func eventsToString(events []corev1.Event) string {
	eventRow := func(e corev1.Event) string {
		return fmt.Sprintf(`* %s/%s: "%s"`, e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Message)
	}

	rows := make([]string, 0, len(events))
	for _, e := range events {
		rows = append(rows, eventRow(e))
	}

	return strings.Join(rows, "\n")
}

const invalidCASecretID = "8214a145-a328-4c56-ab72-2973a56d4eae"

func invalidCASecret(ns string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testutils.RandomName(testTranslationFailuresObjectsPrefix),
			Namespace: ns,
			Labels: map[string]string{
				"konghq.com/ca-cert": "true",
			},
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Data: map[string][]byte{
			"id": []byte(invalidCASecretID),
			// missing cert key
		},
	}
}

func pluginUsingInvalidCACert(ns string) *kongv1.KongPlugin {
	return &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testutils.RandomName(testTranslationFailuresObjectsPrefix),
			Namespace: ns,
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Config:     apiextensionsv1.JSON{Raw: []byte(fmt.Sprintf(`{"ca_certificates": ["%s"]}`, invalidCASecretID))},
		PluginName: "mtls-auth",
	}
}

func httpRouteWithBackends(gatewayName string, services ...*corev1.Service) *gatewayapi.HTTPRoute {
	backendRefs := make([]gatewayapi.HTTPBackendRef, 0, len(services))

	if len(services) > 0 {
		httpPort := gatewayapi.PortNumber(80)
		weight := int32(100 / len(services))

		for _, service := range services {
			backendRefs = append(backendRefs,
				gatewayapi.HTTPBackendRef{
					BackendRef: gatewayapi.BackendRef{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Name: gatewayapi.ObjectName(service.Name),
							Port: &httpPort,
							Kind: util.StringToGatewayAPIKindPtr("Service"),
						},
						Weight: &weight,
					},
				})
		}
	}

	pathMatchPrefix := gatewayapi.PathMatchPathPrefix
	return &gatewayapi.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: testutils.RandomName(testTranslationFailuresObjectsPrefix),
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.StripPathKey: "true",
			},
		},
		Spec: gatewayapi.HTTPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name: gatewayapi.ObjectName(gatewayName),
				}},
			},
			Rules: []gatewayapi.HTTPRouteRule{
				{
					Matches: []gatewayapi.HTTPRouteMatch{
						{
							Path: &gatewayapi.HTTPPathMatch{
								Type:  &pathMatchPrefix,
								Value: kong.String("/test"),
							},
						},
					},
					BackendRefs: backendRefs,
				},
			},
		},
	}
}

func ingressWithPathBackedByService(service *corev1.Service) *netv1.Ingress {
	pathType := netv1.PathTypePrefix
	return &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: testutils.RandomName(testTranslationFailuresObjectsPrefix),
		},
		Spec: netv1.IngressSpec{
			IngressClassName: kong.String(consts.IngressClass),
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: service.Name,
											Port: netv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func validService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: testutils.RandomName(testTranslationFailuresObjectsPrefix),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: 80,
				},
			},
		},
	}
}
