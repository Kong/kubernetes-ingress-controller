package fallback

import (
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

type CacheGraphProvider interface {
	// CacheToGraph returns a new ConfigGraph instance built from the given cache snapshot.
	CacheToGraph(cache store.CacheStores) (*ConfigGraph, error)
}

// Generator is responsible for generating fallback cache snapshots.
type Generator struct {
	cacheGraphProvider CacheGraphProvider
}

func NewGenerator(cacheGraphProvider CacheGraphProvider) *Generator {
	return &Generator{
		cacheGraphProvider: cacheGraphProvider,
	}
}

// GenerateExcludingAffected generates a new cache snapshot that excludes all objects that depend on the broken objects.
func (g *Generator) GenerateExcludingAffected(
	cache store.CacheStores,
	brokenObjects []ObjectHash,
) (store.CacheStores, error) {
	graph, err := g.cacheGraphProvider.CacheToGraph(cache)
	if err != nil {
		return store.CacheStores{}, fmt.Errorf("failed to build cache graph: %w", err)
	}

	fallbackCache, err := cache.TakeSnapshot()
	if err != nil {
		return store.CacheStores{}, fmt.Errorf("failed to take cache snapshot: %w", err)
	}

	for _, brokenObject := range brokenObjects {
		subgraphObjects, err := graph.SubgraphObjects(brokenObject)
		if err != nil {
			return store.CacheStores{}, fmt.Errorf("failed to find dependants for %s: %w", brokenObject, err)
		}
		for _, obj := range subgraphObjects {
			if err := fallbackCache.Delete(obj); err != nil {
				return store.CacheStores{}, fmt.Errorf("failed to delete %s from the cache: %w", GetObjectHash(obj), err)
			}
		}
	}

	return fallbackCache, nil
}
