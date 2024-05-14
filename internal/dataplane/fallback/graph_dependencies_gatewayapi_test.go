package fallback_test

import (
	"testing"

	"github.com/samber/lo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestResolveDependencies_HTTPRoute(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
				},
			},
			cache: cacheStoresFromObjs(t,
				testService(t, "1"),
				testService(t, "2"),
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			),
			expected: []client.Object{},
		},
		{
			name: "HTTPRoute -> Service",
			object: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Rules: []gatewayapi.HTTPRouteRule{
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								{
									BackendRef: gatewayapi.BackendRef{
										BackendObjectReference: gatewayapi.BackendObjectReference{
											Name: "1",
											Kind: lo.ToPtr(gatewayapi.Kind("Service")),
										},
									},
								},
								{
									BackendRef: gatewayapi.BackendRef{
										BackendObjectReference: gatewayapi.BackendObjectReference{
											Name: "2",
											Kind: lo.ToPtr(gatewayapi.Kind("Service")),
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
			),
			expected: []client.Object{
				testService(t, "1"),
				testService(t, "2"),
			},
		},
		{
			name: "HTTPRoute -> KongPlugin, KongClusterPlugin",
			object: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1,2,cluster-1,cluster-2",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			),
			expected: []client.Object{
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}

func TestResolveDependencies_TLSRoute(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &gatewayapi.TLSRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
				},
			},
			cache: cacheStoresFromObjs(t,
				testService(t, "1"),
				testService(t, "2"),
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			),
			expected: []client.Object{},
		},
		{
			name: "TLSRoute -> Service",
			object: &gatewayapi.TLSRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
				},
				Spec: gatewayapi.TLSRouteSpec{
					Rules: []gatewayapi.TLSRouteRule{
						{
							BackendRefs: []gatewayapi.BackendRef{
								{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "1",
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
									},
								},
								{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "2",
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
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
			),
			expected: []client.Object{
				testService(t, "1"),
				testService(t, "2"),
			},
		},
		{
			name: "TLSRoute -> KongPlugin, KongClusterPlugin",
			object: &gatewayapi.TLSRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1,2,cluster-1,cluster-2",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			),
			expected: []client.Object{
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}

func TestResolveDependencies_TCPRoute(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &gatewayapi.TCPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
				},
			},
			cache: cacheStoresFromObjs(t,
				testService(t, "1"),
				testService(t, "2"),
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			),
			expected: []client.Object{},
		},
		{
			name: "TCPRoute -> Service",
			object: &gatewayapi.TCPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
				},
				Spec: gatewayapi.TCPRouteSpec{
					Rules: []gatewayapi.TCPRouteRule{
						{
							BackendRefs: []gatewayapi.BackendRef{
								{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "1",
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
									},
								},
								{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "2",
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
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
			),
			expected: []client.Object{
				testService(t, "1"),
				testService(t, "2"),
			},
		},
		{
			name: "TCPRoute -> KongPlugin, KongClusterPlugin",
			object: &gatewayapi.TCPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1,2,cluster-1,cluster-2",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			),
			expected: []client.Object{
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}

func TestResolveDependencies_UDPRoute(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &gatewayapi.UDPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
				},
			},
			cache: cacheStoresFromObjs(t,
				testService(t, "1"),
				testService(t, "2"),
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			),
			expected: []client.Object{},
		},
		{
			name: "UDPRoute -> Service",
			object: &gatewayapi.UDPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
				},
				Spec: gatewayapi.UDPRouteSpec{
					Rules: []gatewayapi.UDPRouteRule{
						{
							BackendRefs: []gatewayapi.BackendRef{
								{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "1",
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
									},
								},
								{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Name: "2",
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
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
			),
			expected: []client.Object{
				testService(t, "1"),
				testService(t, "2"),
			},
		},
		{
			name: "UDPRoute -> KongPlugin, KongClusterPlugin",
			object: &gatewayapi.UDPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1,2,cluster-1,cluster-2",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			),
			expected: []client.Object{
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}

func TestResolveDependencies_GRPCRoute(t *testing.T) {
	testCases := []resolveDependenciesTestCase{
		{
			name: "no dependencies",
			object: &gatewayapi.GRPCRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
				},
			},
			cache: cacheStoresFromObjs(t,
				testService(t, "1"),
				testService(t, "2"),
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			),
			expected: []client.Object{},
		},
		{
			name: "GRPCRoute -> Service",
			object: &gatewayapi.GRPCRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
				},
				Spec: gatewayapi.GRPCRouteSpec{
					Rules: []gatewayapi.GRPCRouteRule{
						{
							BackendRefs: []gatewayapi.GRPCBackendRef{
								{
									BackendRef: gatewayapi.BackendRef{
										BackendObjectReference: gatewayapi.BackendObjectReference{
											Name: "1",
											Kind: lo.ToPtr(gatewayapi.Kind("Service")),
										},
									},
								},
								{
									BackendRef: gatewayapi.BackendRef{
										BackendObjectReference: gatewayapi.BackendObjectReference{
											Name: "2",
											Kind: lo.ToPtr(gatewayapi.Kind("Service")),
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
			),
			expected: []client.Object{
				testService(t, "1"),
				testService(t, "2"),
			},
		},
		{
			name: "GRPCRoute -> KongPlugin, KongClusterPlugin",
			object: &gatewayapi.GRPCRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "test-namespace",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "1,2,cluster-1,cluster-2",
					},
				},
			},
			cache: cacheStoresFromObjs(t,
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			),
			expected: []client.Object{
				testKongPlugin(t, "1"),
				testKongPlugin(t, "2"),
				testKongClusterPlugin(t, "cluster-1"),
				testKongClusterPlugin(t, "cluster-2"),
			},
		},
	}

	for _, tc := range testCases {
		runResolveDependenciesTest(t, tc)
	}
}
