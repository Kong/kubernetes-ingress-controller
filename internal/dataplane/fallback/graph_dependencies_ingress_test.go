package fallback_test

import (
	"testing"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

func TestResolveDependencies_Ingress(t *testing.T) {
	testIngressClass := helpers.WithTypeMeta(t, &netv1.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-ingress-class",
		},
	})
	testService := helpers.WithTypeMeta(t, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
	})
	testKongServiceFacade := helpers.WithTypeMeta(t, &incubatorv1alpha1.KongServiceFacade{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-kong-service-facade",
			Namespace: "test-namespace",
		},
	})
	testKongPlugin := helpers.WithTypeMeta(t, &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-plugin",
			Namespace: "test-namespace",
		},
	})
	testKongClusterPlugin := helpers.WithTypeMeta(t, &kongv1.KongClusterPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster-plugin",
			Namespace: "test-namespace",
		},
	})

	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
				},
			},
			cache:    cacheStoresFromObjs(t, testIngressClass),
			expected: []client.Object{},
		},
		{
			name: "Ingress -> IngressClass - annotation",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.IngressClassKey: "test-ingress-class",
					},
				},
			},
			cache:    cacheStoresFromObjs(t, testIngressClass),
			expected: []client.Object{testIngressClass},
		},
		{
			name: "Ingress -> IngressClass - field",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
				},
				Spec: netv1.IngressSpec{
					IngressClassName: lo.ToPtr("test-ingress-class"),
				},
			},
			cache: cacheStoresFromObjs(t, testIngressClass),
			expected: []client.Object{
				testIngressClass,
			},
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
													Name: "test-service",
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
			cache:    cacheStoresFromObjs(t, testService),
			expected: []client.Object{testService},
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
													Name:     "test-kong-service-facade",
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
			cache:    cacheStoresFromObjs(t, testKongServiceFacade),
			expected: []client.Object{testKongServiceFacade},
		},
		{
			name: "Ingress -> KongPlugin, KongClusterPlugin",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "test-plugin,test-cluster-plugin",
					},
				},
			},
			cache:    cacheStoresFromObjs(t, testKongPlugin, testKongClusterPlugin),
			expected: []client.Object{testKongPlugin, testKongClusterPlugin},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}
