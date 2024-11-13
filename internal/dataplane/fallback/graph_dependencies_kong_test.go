package fallback_test

import (
	"testing"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
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
		{
			name: "KongConsumer -> Secret from credentials",
			object: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongconsumer",
					Namespace: testNamespace,
				},
				Credentials: []string{"1"},
			},
			cache: cacheStoresFromObjs(t,
				testSecret(t, "1"),
				testKongPlugin(t, "1"),
			),
			expected: []client.Object{testSecret(t, "1")},
		},
		{
			name: "KongConsumer -> non existing Secret from credentials",
			object: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongconsumer",
					Namespace: testNamespace,
				},
				Credentials: []string{"non-existing"},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
			),
			expected: []client.Object{},
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

func TestResolveDependencies_UDPIngress(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &kongv1beta1.UDPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-UDPIngress",
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
			name: "UDPIngress -> plugins - annotation (KongPlugin and KongClusterPlugin with the same name) and referenced Service",
			object: &kongv1beta1.UDPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-UDPIngress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1, 2",
					},
				},
				Spec: kongv1beta1.UDPIngressSpec{
					Rules: []kongv1beta1.UDPIngressRule{
						{
							Backend: kongv1beta1.IngressBackend{
								ServiceName: "1",
							},
						},
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testService(t, "1"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongPlugin(t, "2"), testService(t, "1")},
		},
		{
			name: "UDPIngress -> plugins - annotation (KongPlugin and KongClusterPlugin with different names) and duplicated Service in different namespaces",
			object: &kongv1beta1.UDPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-UDPIngress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1, 3",
					},
				},
				Spec: kongv1beta1.UDPIngressSpec{
					Rules: []kongv1beta1.UDPIngressRule{
						{
							Backend: kongv1beta1.IngressBackend{
								ServiceName: "1",
							},
						},
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "3"),
				testService(t, "1"),
				testService(t, "1", func(s *corev1.Service) {
					s.Namespace = "other-namespace"
				}),
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongClusterPlugin(t, "3"), testService(t, "1")},
		},
		{
			name: "UDPIngress -> plugins - annotation (KongClusterPlugin) and referenced Services",
			object: &kongv1beta1.UDPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-UDPIngress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "3",
						kongv1beta1.KongUpstreamPolicyAnnotationKey:           "3",
					},
				},
				Spec: kongv1beta1.UDPIngressSpec{
					Rules: []kongv1beta1.UDPIngressRule{
						{
							Backend: kongv1beta1.IngressBackend{
								ServiceName: "1",
							},
						},
						{
							Backend: kongv1beta1.IngressBackend{
								ServiceName: "2",
							},
						},
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "3"),
				testKongServiceFacade(t, "1"),
				testKongServiceFacade(t, "2"),
				testService(t, "1"),
				testService(t, "2"),
			),
			expected: []client.Object{testKongClusterPlugin(t, "3"), testService(t, "1"), testService(t, "2")},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}

func TestResolveDependencies_TCPIngress(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &kongv1beta1.TCPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-TCPIngress",
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
			name: "TCPIngress -> plugins - annotation (KongPlugin and KongClusterPlugin with the same name) and referenced Service",
			object: &kongv1beta1.TCPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-TCPIngress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1, 2",
					},
				},
				Spec: kongv1beta1.TCPIngressSpec{
					Rules: []kongv1beta1.IngressRule{
						{
							Backend: kongv1beta1.IngressBackend{
								ServiceName: "1",
							},
						},
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testService(t, "1"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongPlugin(t, "2"), testService(t, "1")},
		},
		{
			name: "TCPIngress -> plugins - annotation (KongPlugin and KongClusterPlugin with different names) and duplicated Service in different namespaces",
			object: &kongv1beta1.TCPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-TCPIngress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1, 3",
					},
				},
				Spec: kongv1beta1.TCPIngressSpec{
					Rules: []kongv1beta1.IngressRule{
						{
							Backend: kongv1beta1.IngressBackend{
								ServiceName: "1",
							},
						},
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "3"),
				testService(t, "1"),
				testService(t, "1", func(s *corev1.Service) {
					s.Namespace = "other-namespace"
				}),
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongClusterPlugin(t, "3"), testService(t, "1")},
		},
		{
			name: "TCPIngress -> plugins - annotation (KongClusterPlugin) and referenced Services",
			object: &kongv1beta1.TCPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-TCPIngress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "3",
						kongv1beta1.KongUpstreamPolicyAnnotationKey:           "3",
					},
				},
				Spec: kongv1beta1.TCPIngressSpec{
					Rules: []kongv1beta1.IngressRule{
						{
							Backend: kongv1beta1.IngressBackend{
								ServiceName: "1",
							},
						},
						{
							Backend: kongv1beta1.IngressBackend{
								ServiceName: "2",
							},
						},
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "3"),
				testKongServiceFacade(t, "1"),
				testKongServiceFacade(t, "2"),
				testService(t, "1"),
				testService(t, "2"),
			),
			expected: []client.Object{testKongClusterPlugin(t, "3"), testService(t, "1"), testService(t, "2")},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}

func TestResolveDependencies_KongCustomEntity(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies - parent reference is nil",
			object: &kongv1alpha1.KongCustomEntity{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-custom-entity",
					Namespace: testNamespace,
				},
				Spec: kongv1alpha1.KongCustomEntitySpec{
					EntityType: "test-entity",
					ParentRef:  nil,
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{},
		},
		{
			name: "parent reference to KongPlugin",
			object: &kongv1alpha1.KongCustomEntity{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-custom-entity",
					Namespace: testNamespace,
				},
				Spec: kongv1alpha1.KongCustomEntitySpec{
					EntityType: "test-entity",
					ParentRef: &kongv1alpha1.ObjectReference{
						Name:  "1",
						Kind:  lo.ToPtr("KongPlugin"),
						Group: lo.ToPtr(kongv1alpha1.GroupVersion.Group),
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{
				testKongPlugin(t, "1"),
			},
		},
		{
			name: "parent reference to KongClusterPlugin",
			object: &kongv1alpha1.KongCustomEntity{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-custom-entity",
					Namespace: testNamespace,
				},
				Spec: kongv1alpha1.KongCustomEntitySpec{
					EntityType: "test-entity",
					ParentRef: &kongv1alpha1.ObjectReference{
						Name:  "1",
						Kind:  lo.ToPtr("KongClusterPlugin"),
						Group: lo.ToPtr(kongv1alpha1.GroupVersion.Group),
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{
				testKongClusterPlugin(t, "1"),
			},
		},
		{
			name: "parent reference to KongPlugin in a different namespace",
			object: &kongv1alpha1.KongCustomEntity{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-custom-entity",
					Namespace: testNamespace,
				},
				Spec: kongv1alpha1.KongCustomEntitySpec{
					EntityType: "test-entity",
					ParentRef: &kongv1alpha1.ObjectReference{
						Name:      "1",
						Namespace: lo.ToPtr("other-namespace"),
						Kind:      lo.ToPtr("KongPlugin"),
						Group:     lo.ToPtr(kongv1alpha1.GroupVersion.Group),
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1", func(p *kongv1.KongPlugin) {
					p.Namespace = "other-namespace"
				}),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}
