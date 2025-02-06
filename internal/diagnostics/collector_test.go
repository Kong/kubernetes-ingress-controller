package diagnostics

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

func TestCollector_EventsHandling(t *testing.T) {
	successfulDump := ConfigDump{
		Meta: DumpMeta{
			Failed:   false,
			Fallback: false,
			Hash:     "success-hash",
		},
		Config: file.Content{
			FormatVersion: "success", // Just for the sake of distinguishing between success and failure.
		},
	}
	failedDump := ConfigDump{
		Config: file.Content{
			FormatVersion: "failed", // Just for the sake of distinguishing between success and failure.
		},
		Meta: DumpMeta{
			Failed:   true,
			Fallback: false,
		},
		RawResponseBody: []byte("error body"),
	}
	fallbackMeta := fallback.GeneratedCacheMetadata{
		BrokenObjects: []fallback.ObjectHash{
			{
				Name: "object",
			},
		},
	}

	c := NewCollector(logr.Discard(), managercfg.Config{
		DumpSensitiveConfig: true,
	})
	t.Run("on successful config dump", func(t *testing.T) {
		c.onConfigDump(successfulDump)
		require.Equal(t, mo.Some(successfulDump.Config), c.lastSuccessfulConfigDump)
		require.Equal(t, successfulDump.Meta.Hash, c.lastSuccessHash)
	})
	t.Run("on failed config dump", func(t *testing.T) {
		c.onConfigDump(failedDump)
		require.Equal(t, mo.Some(failedDump.Config), c.lastFailedConfigDump)
		require.Equal(t, failedDump.Meta.Hash, c.lastFailedHash)
		require.Equal(t, failedDump.RawResponseBody, c.lastRawErrBody)
	})
	t.Run("on fallback cache metadata", func(t *testing.T) {
		c.onFallbackCacheMetadata(fallbackMeta)
		require.NotNilf(t, c.currentFallbackCacheMetadata, "expected fallback cache metadata to be set")
		meta, ok := c.currentFallbackCacheMetadata.Get()
		require.True(t, ok)
		require.Equal(t, fallbackMeta, meta)
	})
	t.Run("on successful config dump after fallback", func(t *testing.T) {
		c.onConfigDump(successfulDump)
		require.Equal(t, mo.Some(successfulDump.Config), c.lastSuccessfulConfigDump)
		require.Equal(t, successfulDump.Meta.Hash, c.lastSuccessHash)
		require.Equal(t, mo.None[fallback.GeneratedCacheMetadata](), c.currentFallbackCacheMetadata, "expected fallback cache metadata to be dropped as it'c no more relevant")
	})
}
