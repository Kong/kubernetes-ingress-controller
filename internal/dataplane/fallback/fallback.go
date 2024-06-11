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

// GenerateExcludingBrokenObjects generates a new cache snapshot that excludes all objects that depend on the broken objects.
func (g *Generator) GenerateExcludingBrokenObjects(
	cache store.CacheStores,
	brokenObjects []ObjectHash,
) (store.CacheStores, GeneratedCacheMetadata, error) {
	metadataCollector := NewGenerateCacheMetadataCollector(brokenObjects...)

	graph, err := g.cacheGraphProvider.CacheToGraph(cache)
	if err != nil {
		return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to build cache graph: %w", err)
	}

	fallbackCache, err := cache.TakeSnapshot()
	if err != nil {
		return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to take cache snapshot: %w", err)
	}

	for _, brokenObject := range brokenObjects {
		subgraphObjects, err := graph.SubgraphObjects(brokenObject)
		if err != nil {
			return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to find dependants for %s: %w", brokenObject, err)
		}
		for _, obj := range subgraphObjects {
			if err := fallbackCache.Delete(obj); err != nil {
				return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to delete %s from the cache: %w", GetObjectHash(obj), err)
			}
			metadataCollector.CollectExcluded(obj, brokenObject)
		}
	}

	return fallbackCache, metadataCollector.Metadata(), nil
}

func (g *Generator) GenerateBackfillingBrokenObjects(
	currentCache store.CacheStores,
	lastValidCacheSnapshot *store.CacheStores,
	brokenObjects []ObjectHash,
) (store.CacheStores, GeneratedCacheMetadata, error) {
	metadataCollector := NewGenerateCacheMetadataCollector(brokenObjects...)

	// Build a graph from the current cache.
	currentGraph, err := g.cacheGraphProvider.CacheToGraph(currentCache)
	if err != nil {
		return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to build current cache graph: %w", err)
	}

	// Take a snapshot of the current cache to use as a fallback.
	fallbackCache, err := currentCache.TakeSnapshot()
	if err != nil {
		return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to take current cache snapshot: %w", err)
	}

	// Exclude the affected objects from the fallback cache. Also, collect all the affected objects as they will be
	// subjects of backfilling.
	var affectedObjects []ObjectHash
	for _, brokenObject := range brokenObjects {
		subgraphObjects, err := currentGraph.SubgraphObjects(brokenObject)
		if err != nil {
			return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to find dependants for %s: %w", brokenObject, err)
		}
		for _, obj := range subgraphObjects {
			if err := fallbackCache.Delete(obj); err != nil {
				return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to delete %s from the fallback cache: %w", GetObjectHash(obj), err)
			}
			affectedObjects = append(affectedObjects, GetObjectHash(obj))
			metadataCollector.CollectExcluded(obj, brokenObject)
		}
	}

	if lastValidCacheSnapshot == nil {
		g.logger.V(util.DebugLevel).Info("No previous valid cache snapshot found, skipping backfilling")
		return fallbackCache, metadataCollector.Metadata(), nil
	}

	// Build a graph from the last valid cache snapshot.
	lastValidGraph, err := g.cacheGraphProvider.CacheToGraph(*lastValidCacheSnapshot)
	if err != nil {
		return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to build cache graph: %w", err)
	}

	// Backfill the affected objects from the last valid cache snapshot.
	for _, affectedObject := range affectedObjects {
		objectsToBackfill, err := lastValidGraph.SubgraphObjects(affectedObject)
		if err != nil {
			return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to find dependants for %s: %w", affectedObject, err)
		}

		for _, obj := range objectsToBackfill {
			if err := fallbackCache.Add(obj); err != nil {
				return store.CacheStores{}, GeneratedCacheMetadata{}, fmt.Errorf("failed to add %s to the cache: %w", GetObjectHash(obj), err)
			}
			metadataCollector.CollectBackfilled(obj, affectedObject)
		}
	}
	return fallbackCache, metadataCollector.Metadata(), nil
}
