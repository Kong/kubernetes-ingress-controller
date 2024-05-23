package fallback

import (
	"fmt"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

type CacheGraphProvider interface {
	// CacheToGraph returns a new ConfigGraph instance built from the given cache snapshot.
	CacheToGraph(cache store.CacheStores) (*ConfigGraph, error)
}

// Generator is responsible for generating fallback cache snapshots.
type Generator struct {
	cacheGraphProvider CacheGraphProvider
	logger             logr.Logger
}

func NewGenerator(cacheGraphProvider CacheGraphProvider, logger logr.Logger) *Generator {
	return &Generator{
		cacheGraphProvider: cacheGraphProvider,
		logger:             logger.WithName("fallback-generator"),
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
			g.logger.V(util.DebugLevel).Info("Excluded object from fallback cache",
				"object_kind", obj.GetObjectKind(),
				"object_name", obj.GetName(),
				"object_namespace", obj.GetNamespace(),
			)
		}
	}

	return fallbackCache, nil
}
