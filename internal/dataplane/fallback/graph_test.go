package fallback_test

import (
	"sort"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

func TestNewConfigGraphFromCacheStores(t *testing.T) {
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
				"Ingress:test-namespace/test-ingress": {},
				"IngressClass:test-ingress-class": {
					"Ingress:test-namespace/test-ingress",
				},
				"Service:test-namespace/test-service": {
					"Ingress:test-namespace/test-ingress",
				},
				"KongServiceFacade:test-namespace/test-kong-service-facade": {
					"Ingress:test-namespace/test-ingress",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := fallback.NewConfigGraphFromCacheStores(tc.cache)
			require.NoError(t, err)
			require.NotNil(t, g)
			require.Equal(t, tc.expectedAdjacencyMap, adjacencyGraphStrings(t, g))
		})
	}
}
