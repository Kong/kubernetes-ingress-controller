package fallback_test

import (
	"sort"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func TestDefaultCacheGraphProvider_CacheToGraph(t *testing.T) {
	// adjacencyGraphStrings returns a map of stringified vertices and their neighbours
	// in the given graph for easy comparison.
	adjacencyGraphStrings := func(t *testing.T, g *fallback.ConfigGraph) map[string][]string {
		am, err := g.AdjacencyMap()
		require.NoError(t, err)
		adjacencyMapStrings := make(map[string][]string, len(am))
		for v, neighbours := range am {
			neighboursStrings := lo.Map(neighbours, func(n fallback.ObjectHash, _ int) string {
				return n.String()
			})
			sort.Strings(neighboursStrings) // Sort for deterministic output.
			adjacencyMapStrings[v.String()] = neighboursStrings
		}
		return adjacencyMapStrings
	}

	testCases := []struct {
		name                 string
		cache                store.CacheStores
		expectedAdjacencyMap map[string][]string
	}{
		{
			name:                 "empty cache",
			cache:                store.NewCacheStores(),
			expectedAdjacencyMap: map[string][]string{},
		},
		{
			name:                 "cache with stores being nil",
			cache:                store.CacheStores{},
			expectedAdjacencyMap: map[string][]string{},
		},
		{
			name: "cache with Ingress and its dependencies",
			cache: cacheStoresFromObjs(t,
				&netv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-ingress",
						Namespace: "test-namespace",
					},
					Spec: netv1.IngressSpec{
						IngressClassName: lo.ToPtr("test-ingress-class"),
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
											{
												Backend: netv1.IngressBackend{
													Resource: &corev1.TypedLocalObjectReference{
														Name:     "test-kong-service-facade",
														Kind:     "KongServiceFacade",
														APIGroup: lo.ToPtr(incubatorv1alpha1.GroupVersion.Group),
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
				&netv1.IngressClass{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-ingress-class",
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service",
						Namespace: "test-namespace",
					},
				},
				&incubatorv1alpha1.KongServiceFacade{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-kong-service-facade",
						Namespace: "test-namespace",
					},
				},
			),
			expectedAdjacencyMap: map[string][]string{
				"networking.k8s.io/Ingress:test-namespace/test-ingress": {},
				"networking.k8s.io/IngressClass:test-ingress-class": {
					"networking.k8s.io/Ingress:test-namespace/test-ingress",
				},
				"core/Service:test-namespace/test-service": {
					"networking.k8s.io/Ingress:test-namespace/test-ingress",
				},
				"incubator.ingress-controller.konghq.com/KongServiceFacade:test-namespace/test-kong-service-facade": {
					"networking.k8s.io/Ingress:test-namespace/test-ingress",
				},
			},
		},
		{
			name: "cache with HTTPRoute and its dependencies",
			cache: cacheStoresFromObjs(t,
				&gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-route",
						Namespace: "test-namespace",
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.PluginsKey: "1,cluster-1",
						},
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
								},
							},
						},
					},
				},
				testService(t, "1"),
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "cluster-1"),
			),
			expectedAdjacencyMap: map[string][]string{
				"gateway.networking.k8s.io/HTTPRoute:test-namespace/test-route": {},
				"core/Service:test-namespace/1": {
					"gateway.networking.k8s.io/HTTPRoute:test-namespace/test-route",
				},
				"configuration.konghq.com/KongPlugin:test-namespace/1": {
					"gateway.networking.k8s.io/HTTPRoute:test-namespace/test-route",
				},
				"configuration.konghq.com/KongClusterPlugin:test-namespace/cluster-1": {
					"gateway.networking.k8s.io/HTTPRoute:test-namespace/test-route",
				},
			},
		},
		{
			name: "cache with TLSRoute and its dependencies",
			cache: cacheStoresFromObjs(t,
				&gatewayapi.TLSRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-route",
						Namespace: "test-namespace",
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.PluginsKey: "1,cluster-1",
						},
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
								},
							},
						},
					},
				},
				testService(t, "1"),
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "cluster-1"),
			),
			expectedAdjacencyMap: map[string][]string{
				"gateway.networking.k8s.io/TLSRoute:test-namespace/test-route": {},
				"core/Service:test-namespace/1": {
					"gateway.networking.k8s.io/TLSRoute:test-namespace/test-route",
				},
				"configuration.konghq.com/KongPlugin:test-namespace/1": {
					"gateway.networking.k8s.io/TLSRoute:test-namespace/test-route",
				},
				"configuration.konghq.com/KongClusterPlugin:test-namespace/cluster-1": {
					"gateway.networking.k8s.io/TLSRoute:test-namespace/test-route",
				},
			},
		},
		{
			name: "cache with TCPRoute and its dependencies",
			cache: cacheStoresFromObjs(t,
				&gatewayapi.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-route",
						Namespace: "test-namespace",
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.PluginsKey: "1,cluster-1",
						},
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
								},
							},
						},
					},
				},
				testService(t, "1"),
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "cluster-1"),
			),
			expectedAdjacencyMap: map[string][]string{
				"gateway.networking.k8s.io/TCPRoute:test-namespace/test-route": {},
				"core/Service:test-namespace/1": {
					"gateway.networking.k8s.io/TCPRoute:test-namespace/test-route",
				},
				"configuration.konghq.com/KongPlugin:test-namespace/1": {
					"gateway.networking.k8s.io/TCPRoute:test-namespace/test-route",
				},
				"configuration.konghq.com/KongClusterPlugin:test-namespace/cluster-1": {
					"gateway.networking.k8s.io/TCPRoute:test-namespace/test-route",
				},
			},
		},
		{
			name: "cache with UDPRoute and its dependencies",
			cache: cacheStoresFromObjs(t,
				&gatewayapi.UDPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-route",
						Namespace: "test-namespace",
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.PluginsKey: "1,cluster-1",
						},
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
								},
							},
						},
					},
				},
				testService(t, "1"),
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "cluster-1"),
			),
			expectedAdjacencyMap: map[string][]string{
				"gateway.networking.k8s.io/UDPRoute:test-namespace/test-route": {},
				"core/Service:test-namespace/1": {
					"gateway.networking.k8s.io/UDPRoute:test-namespace/test-route",
				},
				"configuration.konghq.com/KongPlugin:test-namespace/1": {
					"gateway.networking.k8s.io/UDPRoute:test-namespace/test-route",
				},
				"configuration.konghq.com/KongClusterPlugin:test-namespace/cluster-1": {
					"gateway.networking.k8s.io/UDPRoute:test-namespace/test-route",
				},
			},
		},
		{
			name: "cache with GRPCRoute and its dependencies",
			cache: cacheStoresFromObjs(t,
				&gatewayapi.GRPCRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-route",
						Namespace: "test-namespace",
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.PluginsKey: "1,cluster-1",
						},
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
								},
							},
						},
					},
				},
				testService(t, "1"),
				testKongPlugin(t, "1"),
				testKongClusterPlugin(t, "cluster-1"),
			),
			expectedAdjacencyMap: map[string][]string{
				"gateway.networking.k8s.io/GRPCRoute:test-namespace/test-route": {},
				"core/Service:test-namespace/1": {
					"gateway.networking.k8s.io/GRPCRoute:test-namespace/test-route",
				},
				"configuration.konghq.com/KongPlugin:test-namespace/1": {
					"gateway.networking.k8s.io/GRPCRoute:test-namespace/test-route",
				},
				"configuration.konghq.com/KongClusterPlugin:test-namespace/cluster-1": {
					"gateway.networking.k8s.io/GRPCRoute:test-namespace/test-route",
				},
			},
		},
		{
			name: "cache with KongCustomEntities and its dependencies",
			cache: cacheStoresFromObjs(t,
				&kongv1alpha1.KongCustomEntity{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-entity-kong-plugin",
						Namespace: testNamespace,
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						ParentRef: &kongv1alpha1.ObjectReference{
							Kind:  lo.ToPtr("KongPlugin"),
							Group: lo.ToPtr(kongv1alpha1.GroupVersion.Group),
							Name:  "test-plugin",
						},
					},
				},
				&kongv1alpha1.KongCustomEntity{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-entity-kong-cluster-plugin",
						Namespace: testNamespace,
					},
					Spec: kongv1alpha1.KongCustomEntitySpec{
						ParentRef: &kongv1alpha1.ObjectReference{
							Kind:  lo.ToPtr("KongClusterPlugin"),
							Group: lo.ToPtr(kongv1alpha1.GroupVersion.Group),
							Name:  "test-cluster-plugin",
						},
					},
				},
				testKongPlugin(t, "test-plugin"),
				testKongClusterPlugin(t, "test-cluster-plugin"),
			),
			expectedAdjacencyMap: map[string][]string{
				"configuration.konghq.com/KongPlugin:test-namespace/test-plugin": {
					"configuration.konghq.com/KongCustomEntity:test-namespace/test-entity-kong-plugin",
				},
				"configuration.konghq.com/KongClusterPlugin:test-namespace/test-cluster-plugin": {
					"configuration.konghq.com/KongCustomEntity:test-namespace/test-entity-kong-cluster-plugin",
				},
				"configuration.konghq.com/KongCustomEntity:test-namespace/test-entity-kong-plugin":         {},
				"configuration.konghq.com/KongCustomEntity:test-namespace/test-entity-kong-cluster-plugin": {},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := fallback.NewDefaultCacheGraphProvider()
			g, err := p.CacheToGraph(tc.cache)
			require.NoError(t, err)
			require.NotNil(t, g)
			require.Equal(t, tc.expectedAdjacencyMap, adjacencyGraphStrings(t, g))
		})
	}
}
