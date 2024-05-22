package fallback_test

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

func TestResolveDependencies_KongPlugin(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongPlugin",
					Namespace: testNamespace,
				},
			},
			cache: cacheStoresFromObjs(t,
				testSecret(t, "1"),
				testSecret(t, "2"),
			),
			expected: []client.Object{},
		},
		{
			name: "KongPlugin -> Secret referenced by ConfigFrom (secret with the same name exists in multiple namespaces)",
			object: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongPlugin",
					Namespace: testNamespace,
				},
				ConfigFrom: &kongv1.ConfigSource{
					SecretValue: kongv1.SecretValueFromSource{
						Secret: "1",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testSecret(t, "1"),
				testSecret(t, "1", func(s *corev1.Secret) {
					s.Namespace = "another-namespace"
				}),
			),
			expected: []client.Object{testSecret(t, "1")},
		},
		{
			name: "KongPlugin -> Secret referenced by ConfigFrom does not exist in the same namespace as the KongPlugin",
			object: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongPlugin",
					Namespace: testNamespace,
				},
				ConfigFrom: &kongv1.ConfigSource{
					SecretValue: kongv1.SecretValueFromSource{
						Secret: "2",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testSecret(t, "1"),
				testSecret(t, "2", func(s *corev1.Secret) {
					s.Namespace = "another-namespace"
				}),
			),
			expected: []client.Object{},
		},
		{
			name: "KongPlugin -> two Secrets referenced by ConfigPatches (Secret with the same name exists in multiple namespaces)",
			object: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongPlugin",
					Namespace: testNamespace,
				},
				ConfigPatches: []kongv1.ConfigPatch{
					{
						ValueFrom: kongv1.ConfigSource{
							SecretValue: kongv1.SecretValueFromSource{
								Secret: "1",
							},
						},
					},
					{
						ValueFrom: kongv1.ConfigSource{
							SecretValue: kongv1.SecretValueFromSource{
								Secret: "2",
							},
						},
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testSecret(t, "1"),
				testSecret(t, "1", func(s *corev1.Secret) {
					s.Namespace = "another-namespace"
				}),
				testSecret(t, "2"),
			),
			expected: []client.Object{testSecret(t, "1"), testSecret(t, "2")},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}

func TestResolveDependencies_KongClusterPlugin(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-KongClusterPlugin",
				},
			},
			cache: cacheStoresFromObjs(t,
				testSecret(t, "1"),
				testSecret(t, "2"),
			),
			expected: []client.Object{},
		},
		{
			name: "KongClusterPlugin -> Secret referenced by ConfigFrom (Secret with the same name exists in multiple namespaces)",
			object: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-KongClusterPlugin",
				},
				ConfigFrom: &kongv1.NamespacedConfigSource{
					SecretValue: kongv1.NamespacedSecretValueFromSource{
						Namespace: testNamespace,
						Secret:    "1",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testSecret(t, "1"),
				testSecret(t, "1", func(s *corev1.Secret) {
					s.Namespace = "another-namespace"
				}),
			),
			expected: []client.Object{testSecret(t, "1")},
		},
		{
			name: "KongClusterPlugin -> Secret referenced by ConfigFrom does not exists",
			object: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-KongClusterPlugin",
				},
				ConfigFrom: &kongv1.NamespacedConfigSource{
					SecretValue: kongv1.NamespacedSecretValueFromSource{
						Namespace: testNamespace,
						Secret:    "1",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testSecret(t, "1", func(s *corev1.Secret) {
					s.Namespace = "another-namespace"
				}),
				testSecret(t, "2"),
			),
			expected: []client.Object{},
		},
		{
			name: "KongClusterPlugin -> two Secrets referenced by ConfigPatches (Secret with the same name exists in multiple namespaces)",
			object: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-KongClusterPlugin",
				},
				ConfigPatches: []kongv1.NamespacedConfigPatch{
					{
						ValueFrom: kongv1.NamespacedConfigSource{
							SecretValue: kongv1.NamespacedSecretValueFromSource{
								Namespace: "another-namespace",
								Secret:    "1",
							},
						},
					},
					{
						ValueFrom: kongv1.NamespacedConfigSource{
							SecretValue: kongv1.NamespacedSecretValueFromSource{
								Namespace: testNamespace,
								Secret:    "2",
							},
						},
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testSecret(t, "1"),
				testSecret(t, "1", func(s *corev1.Secret) {
					s.Namespace = "another-namespace"
				}),
				testSecret(t, "2"),
			),
			expected: []client.Object{
				testSecret(t, "1", func(s *corev1.Secret) {
					s.Namespace = "another-namespace"
				}),
				testSecret(t, "2"),
			},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}

func TestResolveDependencies_KongConsumer(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongConsumer",
					Namespace: testNamespace,
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{},
		},
		{
			name: "KongConsumer -> plugins - annotation (KongPlugin and KongClusterPlugin with the same name)",
			object: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongconsumer",
					Namespace: testNamespace,
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1, 2",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongPlugin(t, "2")},
		},
		{
			name: "KongConsumer -> plugins - annotation (KongPlugin and KongClusterPlugin with different names)",
			object: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongconsumer",
					Namespace: testNamespace,
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1, 3",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "3"),
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongClusterPlugin(t, "3")},
		},
		{
			name: "KongConsumer -> plugins - annotation (KongClusterPlugin)",
			object: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongconsumer",
					Namespace: testNamespace,
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "3",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "3"),
			),
			expected: []client.Object{testKongClusterPlugin(t, "3")},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}

func TestResolveDependencies_KongConsumerGroup(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &kongv1beta1.KongConsumerGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongConsumerGroup",
					Namespace: testNamespace,
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{},
		},
		{
			name: "KongConsumerGroup -> plugins - annotation (KongPlugin and KongClusterPlugin with the same name)",
			object: &kongv1beta1.KongConsumerGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongConsumerGroup",
					Namespace: testNamespace,
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1, 2",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongPlugin(t, "2")},
		},
		{
			name: "KongConsumerGroup -> plugins - annotation (KongPlugin and KongClusterPlugin with different names)",
			object: &kongv1beta1.KongConsumerGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongConsumerGroup",
					Namespace: testNamespace,
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1, 3",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "3"),
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongClusterPlugin(t, "3")},
		},
		{
			name: "KongConsumerGroup -> plugins - annotation (KongClusterPlugin)",
			object: &kongv1beta1.KongConsumerGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongConsumerGroup",
					Namespace: testNamespace,
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "3",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "3"),
			),
			expected: []client.Object{testKongClusterPlugin(t, "3")},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}

func TestResolveDependencies_KongServiceFacade(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &incubatorv1alpha1.KongServiceFacade{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongServiceFacade",
					Namespace: "test-namespace",
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{},
		},
		{
			name: "KongServiceFacade -> plugins - annotation (KongPlugin and KongClusterPlugin with the same name) and KongUpstreamPolicy",
			object: &incubatorv1alpha1.KongServiceFacade{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongServiceFacade",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1, 2",
						kongv1beta1.KongUpstreamPolicyAnnotationKey:           "1",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongUpstreamPolicy(t, "1"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongPlugin(t, "2"), testKongUpstreamPolicy(t, "1")},
		},
		{
			name: "KongServiceFacade -> plugins - annotation (KongPlugin and KongClusterPlugin with different names)",
			object: &incubatorv1alpha1.KongServiceFacade{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongServiceFacade",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1, 3",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "3"),
				testKongUpstreamPolicy(t, "1"),
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongClusterPlugin(t, "3")},
		},
		{
			name: "KongServiceFacade -> plugins - annotation (KongClusterPlugin) and KongUpstreamPolicy",
			object: &incubatorv1alpha1.KongServiceFacade{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongServiceFacade",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "3",
						kongv1beta1.KongUpstreamPolicyAnnotationKey:           "3",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "3"),
				testKongUpstreamPolicy(t, "3"),
			),
			expected: []client.Object{testKongClusterPlugin(t, "3"), testKongUpstreamPolicy(t, "3")},
		},
		{
			name: "KongServiceFacade -> KongUpstreamPolicy - the same name in different namespaces",
			object: &incubatorv1alpha1.KongServiceFacade{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongServiceFacade",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						kongv1beta1.KongUpstreamPolicyAnnotationKey: "1",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongUpstreamPolicy(t, "1"),
				testKongUpstreamPolicy(t, "1", func(kup *kongv1beta1.KongUpstreamPolicy) {
					kup.Namespace = "other-namespace"
				}),
			),
			expected: []client.Object{testKongUpstreamPolicy(t, "1")},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}
