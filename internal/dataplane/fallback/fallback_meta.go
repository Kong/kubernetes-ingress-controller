package fallback

import (
	"github.com/samber/lo"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GeneratedCacheMetadata contains metadata generated during the fallback process.
type GeneratedCacheMetadata struct {
	// BrokenObjects are objects that were reported as broken by the Kong Admin API.
	BrokenObjects []ObjectHash
	// ExcludedObjects are objects that were excluded from the fallback configuration as they were broken or either of their
	// dependencies was broken.
	ExcludedObjects []AffectedCacheObjectMetadata
	// BackfilledObjects are objects that were backfilled from the last valid cache state as they were broken or either of
	// their dependencies was broken.
	BackfilledObjects []AffectedCacheObjectMetadata
}

// GeneratedCacheMetadataCollector is a collector for cache metadata generated during the fallback process.
// It's primarily used to deduplicate the metadata and make it easier to work with.
type GeneratedCacheMetadataCollector struct {
	brokenObjects     []ObjectHash
	excludedObjects   map[ObjectHash]AffectedCacheObjectMetadata
	backfilledObjects map[ObjectHash]AffectedCacheObjectMetadata
}

// AffectedCacheObjectMetadata contains an object and a list of objects that caused it to be excluded or backfilled
// during the fallback process.
type AffectedCacheObjectMetadata struct {
	Object         client.Object
	CausingObjects []ObjectHash
}

// NewGenerateCacheMetadataCollector creates a new GeneratedCacheMetadataCollector instance.
func NewGenerateCacheMetadataCollector(brokenObjects ...ObjectHash) *GeneratedCacheMetadataCollector {
	return &GeneratedCacheMetadataCollector{
		brokenObjects:     brokenObjects,
		excludedObjects:   make(map[ObjectHash]AffectedCacheObjectMetadata),
		backfilledObjects: make(map[ObjectHash]AffectedCacheObjectMetadata),
	}
}

// CollectExcluded collects an excluded object (an object that was excluded from the fallback configuration as it was
// broken or one of its dependencies was broken).
func (m *GeneratedCacheMetadataCollector) CollectExcluded(excluded client.Object, causing ObjectHash) {
	objHash := GetObjectHash(excluded)
	if existingEntry, ok := m.excludedObjects[objHash]; ok {
		existingEntry.CausingObjects = append(existingEntry.CausingObjects, causing)
		m.excludedObjects[objHash] = existingEntry
	} else {
		m.excludedObjects[objHash] = AffectedCacheObjectMetadata{Object: excluded, CausingObjects: []ObjectHash{causing}}
	}
}

// CollectBackfilled collects a backfilled object (an object that was backfilled from the last valid cache state as that or
// one of its dependencies was broken).
func (m *GeneratedCacheMetadataCollector) CollectBackfilled(backfilled client.Object, causing ObjectHash) {
	objHash := GetObjectHash(backfilled)
	if existingEntry, ok := m.backfilledObjects[objHash]; ok {
		existingEntry.CausingObjects = append(existingEntry.CausingObjects, causing)
		m.backfilledObjects[objHash] = existingEntry
	} else {
		m.backfilledObjects[objHash] = AffectedCacheObjectMetadata{Object: backfilled, CausingObjects: []ObjectHash{causing}}
	}
}

// Metadata generates the final cache metadata from the collected data.
func (m *GeneratedCacheMetadataCollector) Metadata() GeneratedCacheMetadata {
	return GeneratedCacheMetadata{
		BrokenObjects:     m.brokenObjects,
		ExcludedObjects:   lo.Values(m.excludedObjects),
		BackfilledObjects: lo.Values(m.backfilledObjects),
	}
}
