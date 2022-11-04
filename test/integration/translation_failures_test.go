//go:build integration_tests
// +build integration_tests

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

// TestTranslationFailures ensures that proper warning Kubernetes events are recorded in case of translation failures
// encountered.
func TestTranslationFailures(t *testing.T) {
	testCases := []struct {
		name string
		// translationFailureTrigger should create objects that trigger translation failure and return the objects
		// that we expect translation failure warning events to be created for.
		translationFailureTrigger func(t *testing.T, cleaner *clusters.Cleaner, ns string) []client.Object
	}{
		{
			name: "invalid CA secret",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) []client.Object {
				createdSecret, err := env.Cluster().Client().CoreV1().Secrets(ns).Create(ctx, invalidCASecret(ns), metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(createdSecret)

				return []client.Object{createdSecret}
			},
		},
		{
			name: "invalid CA secret referred by a plugin",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) []client.Object {
				createdSecret, err := env.Cluster().Client().CoreV1().Secrets(ns).Create(ctx, invalidCASecret(ns), metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(createdSecret)

				c, err := clientset.NewForConfig(env.Cluster().Config())
				require.NoError(t, err)
				createdPlugin, err := c.ConfigurationV1().KongPlugins(ns).Create(ctx, pluginUsingInvalidCACert(ns), metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(createdPlugin)

				// expect events for both: a faulty secret and a plugin referring it
				return []client.Object{createdSecret, createdPlugin}
			},
		},
		{
			name: "grouped services annotations do not match",
			translationFailureTrigger: func(t *testing.T, cleaner *clusters.Cleaner, ns string) []client.Object {
				container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
				d1 := generators.NewDeploymentForContainer(container)
				d1.Name = ""
				d1.GenerateName = "deployment-"
				d1, err := env.Cluster().Client().AppsV1().Deployments(ns).Create(ctx, d1, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(d1)

				d2 := generators.NewDeploymentForContainer(container)
				d2.Name = ""
				d2.GenerateName = "deployment-"
				d2, err = env.Cluster().Client().AppsV1().Deployments(ns).Create(ctx, d2, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(d2)

				service1 := generators.NewServiceForDeployment(d1, corev1.ServiceTypeClusterIP)
				service1.Annotations = map[string]string{"konghq.com/annotation": "true"}
				_, err = env.Cluster().Client().CoreV1().Services(ns).Create(ctx, service1, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(service1)

				service2 := generators.NewServiceForDeployment(d2, corev1.ServiceTypeClusterIP)
				service2.Annotations = map[string]string{"konghq.com/annotation": "false"}
				_, err = env.Cluster().Client().CoreV1().Services(ns).Create(ctx, service2, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(service2)

				pathType := netv1.PathTypePrefix
				ingress := &netv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: "ingress-",
						Annotations: map[string]string{
							"konghq.com/strip-path": "true",
						},
					},
					Spec: netv1.IngressSpec{
						IngressClassName: kong.String(ingressClass),
						Rules: []netv1.IngressRule{
							{
								IngressRuleValue: netv1.IngressRuleValue{
									HTTP: &netv1.HTTPIngressRuleValue{
										Paths: []netv1.HTTPIngressPath{
											{
												Path:     "/test_1",
												PathType: &pathType,
												Backend: netv1.IngressBackend{
													Service: &netv1.IngressServiceBackend{
														Name: service1.Name,
														Port: netv1.ServiceBackendPort{
															Number: service1.Spec.Ports[0].Port,
														},
													},
												},
											},
											{
												Path:     "/test_2",
												PathType: &pathType,
												Backend: netv1.IngressBackend{
													Service: &netv1.IngressServiceBackend{
														Name: service2.Name,
														Port: netv1.ServiceBackendPort{
															Number: service2.Spec.Ports[0].Port,
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
				require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns, ingress))
				cleaner.Add(ingress)

				return []client.Object{service1, service2}
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ns, cleaner := setup(t)
			defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

			expectedCausingObjects := tt.translationFailureTrigger(t, cleaner, ns.GetName())

			require.Eventually(t, func() bool {
				eventsForAllObjectsFound := true

				for _, expectedCausingObject := range expectedCausingObjects {
					events, err := env.Cluster().Client().CoreV1().Events(ns.GetName()).List(ctx, metav1.ListOptions{
						FieldSelector: fmt.Sprintf(
							"reason=%s,type=%s,involvedObject.name=%s",
							dataplane.KongConfigurationTranslationFailedEventReason,
							corev1.EventTypeWarning,
							expectedCausingObject.GetName(),
						),
					})
					if err != nil {
						t.Logf("failed to list events: %s", err)
						eventsForAllObjectsFound = false
					}

					if len(events.Items) == 0 {
						t.Logf("waiting for events related to '%s' to be created", expectedCausingObject.GetName())
						eventsForAllObjectsFound = false
					}
				}

				return eventsForAllObjectsFound
			}, time.Minute*5, time.Second)
		})
	}
}

const invalidCASecretID = "8214a145-a328-4c56-ab72-2973a56d4eae" //nolint:gosec

func invalidCASecret(ns string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "ca-secret-",
			Namespace:    ns,
			Labels: map[string]string{
				"konghq.com/ca-cert": "true",
			},
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
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
			GenerateName: "kong-plugin-",
			Namespace:    ns,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Config:     v1.JSON{Raw: []byte(fmt.Sprintf(`{"ca_certificates": ["%s"]}`, invalidCASecretID))},
		PluginName: "mtls-auth",
	}
}
