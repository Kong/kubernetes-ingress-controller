package fallback_test

import (
	"testing"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

func TestResolveDependencies_Ingress(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
				},
			},
			cache: cacheStoresFromObjs(t,
				testIngressClass(t, "1"),
				testService(t, "1"),
				testKongServiceFacade(t, "1"),
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "1"),
			),
			expected: []client.Object{},
		},
		{
			name: "Ingress -> IngressClass - annotation",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.IngressClassKey: "1",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testIngressClass(t, "1"),
				testIngressClass(t, "2"),
			),
			expected: []client.Object{testIngressClass(t, "1")},
		},
		{
			name: "Ingress -> IngressClass - field",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
				},
				Spec: netv1.IngressSpec{
					IngressClassName: lo.ToPtr("1"),
				},
			},
			cache: cacheStoresFromObjs(t,
				testIngressClass(t, "1"),
				testIngressClass(t, "2"),
			),
			expected: []client.Object{testIngressClass(t, "1")},
		},
		{
			name: "Ingress -> Service",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "1",
												},
											},
										},
										{
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "2",
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
			cache: cacheStoresFromObjs(t,
				testService(t, "1"),
				testService(t, "2"),
				testService(t, "3"),
			),
			expected: []client.Object{
				testService(t, "1"),
				testService(t, "2"),
			},
		},
		{
			name: "Ingress -> KongServiceFacade",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Backend: netv1.IngressBackend{
												Resource: &corev1.TypedLocalObjectReference{
													Name:     "1",
													Kind:     "KongServiceFacade",
													APIGroup: lo.ToPtr(incubatorv1alpha1.SchemeGroupVersion.Group),
												},
											},
										},
										{
											Backend: netv1.IngressBackend{
												Resource: &corev1.TypedLocalObjectReference{
													Name:     "2",
													Kind:     "KongServiceFacade",
													APIGroup: lo.ToPtr(incubatorv1alpha1.SchemeGroupVersion.Group),
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
			cache: cacheStoresFromObjs(t,
				testKongServiceFacade(t, "1"),
				testKongServiceFacade(t, "2"),
				testKongServiceFacade(t, "3"),
			),
			expected: []client.Object{
				testKongServiceFacade(t, "1"),
				testKongServiceFacade(t, "2"),
			},
		},
		{
			name: "Ingress -> KongPlugin, KongClusterPlugin",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1,2,cluster-1,cluster-2",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongPlugin(t, "3"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
				testKongClusterPlugin(t, "cluster-3"),
			),
			expected: []client.Object{
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			},
		},
		{
			name: "Ingress -> all dependencies at once",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1,cluster-1",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Backend: netv1.IngressBackend{
												Resource: &corev1.TypedLocalObjectReference{
													Name:     "1",
													Kind:     "KongServiceFacade",
													APIGroup: lo.ToPtr(incubatorv1alpha1.SchemeGroupVersion.Group),
												},
											},
										},
										{
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "1",
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
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "cluster-1"),
				testService(t, "1"),
				testKongServiceFacade(t, "1"),
			),
			expected: []client.Object{
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "cluster-1"),
				testService(t, "1"),
				testKongServiceFacade(t, "1"),
			},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}
