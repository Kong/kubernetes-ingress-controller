//go:build integration_tests
// +build integration_tests

package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
)

// TestTranslationFailures ensures that proper warning Kubernetes events are recorded in case of translation failures
// encountered.
func TestTranslationFailures(t *testing.T) {
	type expectedTranslationFailure struct {
		object          client.Object
		messageContains string
	}

	testCases := []struct {
		name string
		// translationFailureScenario should create objects that trigger translation failure and return the objects
		// that we expect translation failure warning events to be created for.
		translationFailureScenario func(t *testing.T, cleaner *clusters.Cleaner, ns string) []expectedTranslationFailure
	}{
		{
			name: "invalid CA secret",
			translationFailureScenario: func(t *testing.T, cleaner *clusters.Cleaner, ns string) []expectedTranslationFailure {
				createdSecret, err := env.Cluster().Client().CoreV1().Secrets(ns).Create(ctx, invalidCASecret(ns), metav1.CreateOptions{})
				require.NoError(t, err)

				return []expectedTranslationFailure{{object: createdSecret}}
			},
		},
		{
			name: "invalid CA secret referred by a plugin",
			translationFailureScenario: func(t *testing.T, cleaner *clusters.Cleaner, ns string) []expectedTranslationFailure {
				createdSecret, err := env.Cluster().Client().CoreV1().Secrets(ns).Create(ctx, invalidCASecret(ns), metav1.CreateOptions{})
				require.NoError(t, err)

				c, err := clientset.NewForConfig(env.Cluster().Config())
				require.NoError(t, err)
				createdPlugin, err := c.ConfigurationV1().KongPlugins(ns).Create(ctx, pluginUsingInvalidCACert(ns), metav1.CreateOptions{})
				require.NoError(t, err)

				// expect events for both: a faulty secret and a plugin referring it
				return []expectedTranslationFailure{{object: createdSecret}, {object: createdPlugin}}
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ns, cleaner := setup(t)
			defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

			expectedTranslationFailures := tt.translationFailureScenario(t, cleaner, ns.GetName())

			require.Eventually(t, func() bool {
				eventsForAllObjectsFound := true

				for _, expected := range expectedTranslationFailures {
					events, err := env.Cluster().Client().CoreV1().Events(ns.GetName()).List(ctx, metav1.ListOptions{
						FieldSelector: fmt.Sprintf(
							"reason=%s,type=%s,involvedObject.name=%s",
							dataplane.KongConfigurationTranslationFailedEventReason,
							corev1.EventTypeWarning,
							expected.object.GetName(),
						),
					})
					if err != nil {
						t.Logf("failed to list events: %s", err)
						eventsForAllObjectsFound = false
						continue
					}

					if len(events.Items) == 0 {
						t.Logf("waiting for events related to '%s' to be created", expected.object.GetName())
						eventsForAllObjectsFound = false
						continue
					}

					if actualMsg := events.Items[0].Message; !strings.Contains(actualMsg, expected.messageContains) {
						t.Logf("expected '%s' not found in the actual event message: '%s'", expected.messageContains, actualMsg)
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
