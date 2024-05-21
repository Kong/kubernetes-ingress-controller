package fallback_test

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

func TestResolveDependencies_KongConsumer(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongConsumer",
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
			name: "KongConsumer -> plugins - annotation (KongPlugin and KongClusterPlugin with the same name)",
			object: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongconsumer",
					Namespace: "test-namespace",
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
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongClusterPlugin(t, "3")},
		},
		{
			name: "KongConsumer -> plugins - annotation (KongClusterPlugin)",
			object: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kongconsumer",
					Namespace: "test-namespace",
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
			name: "KongConsumerGroup -> plugins - annotation (KongPlugin and KongClusterPlugin with the same name)",
			object: &kongv1beta1.KongConsumerGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongConsumerGroup",
					Namespace: "test-namespace",
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
			),
			expected: []client.Object{testKongPlugin(t, "1"), testKongClusterPlugin(t, "3")},
		},
		{
			name: "KongConsumerGroup -> plugins - annotation (KongClusterPlugin)",
			object: &kongv1beta1.KongConsumerGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-KongConsumerGroup",
					Namespace: "test-namespace",
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
