package fallback_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
)

func TestGeneratedCacheMetadataCollector(t *testing.T) {
	excluded1 := testService(t, "excluded-1")
	excluded2 := testService(t, "excluded-2")
	backfilled1 := testService(t, "backfilled-1")
	backfilled2 := testService(t, "backfilled-2")
	causing1 := testService(t, "causing-1")
	causing2 := testService(t, "causing-2")

	t.Run("only excluded", func(t *testing.T) {
		c := fallback.NewGenerateCacheMetadataCollector(
			fallback.GetObjectHash(causing1),
			fallback.GetObjectHash(causing2),
		)
		c.CollectExcluded(excluded1, fallback.GetObjectHash(causing1))
		c.CollectExcluded(excluded2, fallback.GetObjectHash(causing1))
		c.CollectExcluded(excluded2, fallback.GetObjectHash(causing2)) // Duplicate with another causing object.

		meta := c.Metadata()
		require.ElementsMatch(t, []fallback.ObjectHash{
			fallback.GetObjectHash(causing1),
			fallback.GetObjectHash(causing2),
		},
			meta.BrokenObjects,
		)
		require.ElementsMatch(t, []fallback.AffectedCacheObjectMetadata{
			{
				Object:         excluded1,
				CausingObjects: []fallback.ObjectHash{fallback.GetObjectHash(causing1)},
			},
			{
				Object:         excluded2,
				CausingObjects: []fallback.ObjectHash{fallback.GetObjectHash(causing1), fallback.GetObjectHash(causing2)},
			},
		},
			meta.ExcludedObjects,
		)
	})

	t.Run("excluded and backfilled", func(t *testing.T) {
		c := fallback.NewGenerateCacheMetadataCollector(
			fallback.GetObjectHash(causing1),
			fallback.GetObjectHash(causing2),
		)
		c.CollectExcluded(excluded1, fallback.GetObjectHash(causing1))
		c.CollectExcluded(excluded2, fallback.GetObjectHash(causing1))
		c.CollectExcluded(excluded2, fallback.GetObjectHash(causing2)) // Duplicate with another causing object.

		c.CollectBackfilled(backfilled1, fallback.GetObjectHash(causing1))
		c.CollectBackfilled(backfilled2, fallback.GetObjectHash(causing1))
		c.CollectBackfilled(backfilled2, fallback.GetObjectHash(causing2)) // Duplicate with another causing object.

		meta := c.Metadata()
		require.ElementsMatch(t, []fallback.ObjectHash{
			fallback.GetObjectHash(causing1),
			fallback.GetObjectHash(causing2),
		},
			meta.BrokenObjects,
		)
		require.ElementsMatch(t, []fallback.AffectedCacheObjectMetadata{
			{
				Object:         excluded1,
				CausingObjects: []fallback.ObjectHash{fallback.GetObjectHash(causing1)},
			},
			{
				Object:         excluded2,
				CausingObjects: []fallback.ObjectHash{fallback.GetObjectHash(causing1), fallback.GetObjectHash(causing2)},
			},
		},
			meta.ExcludedObjects,
		)
		require.ElementsMatch(t, []fallback.AffectedCacheObjectMetadata{
			{
				Object:         backfilled1,
				CausingObjects: []fallback.ObjectHash{fallback.GetObjectHash(causing1)},
			},
			{
				Object:         backfilled2,
				CausingObjects: []fallback.ObjectHash{fallback.GetObjectHash(causing1), fallback.GetObjectHash(causing2)},
			},
		},
			meta.BackfilledObjects,
		)
	})
}
