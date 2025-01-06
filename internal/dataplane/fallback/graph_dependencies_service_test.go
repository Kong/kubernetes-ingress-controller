package fallback_test

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
)

func TestResolveDependencies_Service(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
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
			name: "Service -> plugins - annotation (KongPlugin and KongClusterPlugin with the same name) and KongUpstreamPolicy",
			object: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
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
			name: "Service -> plugins - annotation (KongPlugin and KongClusterPlugin with different names)",
			object: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
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
			name: "Service -> plugins - annotation (KongClusterPlugin) and KongUpstreamPolicy",
			object: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
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
			name: "Service -> KongUpstreamPolicy - the same name in different namespaces",
			object: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
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
