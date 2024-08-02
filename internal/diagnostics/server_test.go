package diagnostics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	testhelpers "github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

// TestDiagnosticsServer_ConfigDumps tests that the diagnostics server can receive and serve config dumps.
// It's primarily to test that write and read operations run simultaneously do not fall into a race condition.
func TestDiagnosticsServer_ConfigDumps(t *testing.T) {
	s := NewServer(logr.Discard(), ServerConfig{
		ConfigDumpsEnabled: true,
	})
	configsCh := s.clientDiagnostic.Configs

	port := testhelpers.GetFreePort(t)
	t.Logf("Obtained a free port: %d", port)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := s.Listen(ctx, port)
		require.NoError(t, err)
	}()
	t.Log("Started diagnostics server")

	// Use a WaitGroup to ensure that both the write and read operations are run simultaneously.
	readWriteWg := sync.WaitGroup{}
	readWriteWg.Add(2)

	// Write 1000 config dumps to the Server.
	const configDumpsToWrite = 1000
	go func() {
		readWriteWg.Done()
		readWriteWg.Wait()

		defer cancel()
		failed := false
		for i := 0; i < configDumpsToWrite; i++ {
			failed = !failed // Toggle failed flag.
			configsCh <- ConfigDump{
				Config:          file.Content{},
				Meta:            DumpMeta{Failed: failed},
				RawResponseBody: []byte("fake error body"),
			}
		}
	}()
	t.Log("Started writing config dumps")

	// Continuously read config dumps from the Server until context is cancelled.
	go func() {
		readWriteWg.Done()
		readWriteWg.Wait()

		httpClient := &http.Client{}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%d/debug/config/successful", port))
				if err == nil {
					_ = resp.Body.Close()
				}
				resp, err = httpClient.Get(fmt.Sprintf("http://localhost:%d/debug/config/failed", port))
				if err == nil {
					_ = resp.Body.Close()
				}
				resp, err = httpClient.Get(fmt.Sprintf("http://localhost:%d/debug/config/raw-error", port))
				if err == nil {
					_ = resp.Body.Close()
				}
			}
		}
	}()
	t.Log("Started reading config dumps")

	<-ctx.Done()
}

func TestServer_EventsHandling(t *testing.T) {
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

	s := NewServer(logr.Discard(), ServerConfig{
		ConfigDumpsEnabled: true,
	})

	t.Run("on successful config dump", func(t *testing.T) {
		s.onConfigDump(successfulDump)
		require.Equal(t, successfulDump.Config, s.lastSuccessfulConfigDump)
		require.Equal(t, successfulDump.Meta.Hash, s.lastSuccessHash)
	})
	t.Run("on failed config dump", func(t *testing.T) {
		s.onConfigDump(failedDump)
		require.Equal(t, failedDump.Config, s.lastFailedConfigDump)
		require.Equal(t, failedDump.Meta.Hash, s.lastFailedHash)
		require.Equal(t, failedDump.RawResponseBody, s.lastRawErrBody)
	})
	t.Run("on fallback cache metadata", func(t *testing.T) {
		s.onFallbackCacheMetadata(fallbackMeta)
		require.NotNilf(t, s.currentFallbackCacheMetadata, "expected fallback cache metadata to be set")
		require.Equal(t, fallbackMeta, *s.currentFallbackCacheMetadata)
	})
	t.Run("on successful config dump after fallback", func(t *testing.T) {
		s.onConfigDump(successfulDump)
		require.Equal(t, successfulDump.Config, s.lastSuccessfulConfigDump)
		require.Equal(t, successfulDump.Meta.Hash, s.lastSuccessHash)
		require.Nil(t, s.currentFallbackCacheMetadata, "expected fallback cache metadata to be dropped as it's no more relevant")
	})
}

// TestDiagnosticsServer_Diffs tests the diff endpoint using fake data.
func TestDiagnosticsServer_Diffs(t *testing.T) {
	s := NewServer(logr.Discard(), ServerConfig{
		ConfigDumpsEnabled:  true,
		DumpSensitiveConfig: true,
	})
	diffCh := s.clientDiagnostic.Diffs

	port := testhelpers.GetFreePort(t)
	t.Logf("Obtained a free port: %d", port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := s.Listen(ctx, port)
		require.NoError(t, err)
	}()
	t.Log("Started diagnostics server")

	// initially write the max number of cached diffs
	configDumpsToWrite := diffHistorySize
	configDiffs := map[string]ConfigDiff{}
	var first, last string
	init := sync.Once{}
	for i := 0; i < configDumpsToWrite; i++ {
		diff := testConfigDiff()
		configDiffs[diff.Hash] = diff
		init.Do(func() { first = diff.Hash })
		diffCh <- diff
		last = diff.Hash
	}

	// request the diff report
	httpClient := &http.Client{}
	resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%d/debug/config/diff-report", port))
	require.NoError(t, err)
	defer resp.Body.Close()

	b := new(bytes.Buffer)
	_, err = b.ReadFrom(resp.Body)
	require.NoError(t, err)
	got := DiffResponse{}
	require.NoError(t, json.Unmarshal(b.Bytes(), &got))

	// the diff returned should be the last one sent
	require.Equal(t, last, got.ConfigHash)

	// Having gotten a response, check that its available list contains all the diffs we've sent, and that we have the
	// expected number of diffs.
	actual := map[string]interface{}{}
	for _, available := range got.Available {
		actual[available.ConfigHash] = nil
	}
	require.Equal(t, len(actual), len(configDiffs))
	for expected := range configDiffs {
		_, ok := actual[expected]
		require.Truef(t, ok, "expected hash %s not found in report", expected)
	}

	// send an additional diff and confirm that the ring buffer clears out an item, so its length does not exceed the max
	extra := testConfigDiff()
	configDiffs[extra.Hash] = extra
	diffCh <- extra

	second, err := httpClient.Get(fmt.Sprintf("http://localhost:%d/debug/config/diff-report", port))
	require.NoError(t, err)
	defer second.Body.Close()

	b = new(bytes.Buffer)
	_, err = b.ReadFrom(second.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(b.Bytes(), &got))
	require.Equal(t, len(got.Available), diffHistorySize)

	// confirm that the by hash endpoints cannot retrieve the last diff sent, and get a 404 for the first (now discarded)
	// diff sent
	third, err := httpClient.Get(fmt.Sprintf("http://localhost:%d/debug/config/diff-report?hash=%s", port, extra.Hash))
	require.NoError(t, err)
	defer third.Body.Close()
	require.Equal(t, third.StatusCode, http.StatusOK)

	fourth, err := httpClient.Get(fmt.Sprintf("http://localhost:%d/debug/config/diff-report?hash=%s", port, first))
	require.NoError(t, err)
	defer fourth.Body.Close()
	require.Equal(t, fourth.StatusCode, http.StatusNotFound)
}

func testConfigDiff() ConfigDiff {
	return ConfigDiff{
		Hash:      uuid.Must(uuid.NewV7()).String(),
		Timestamp: time.Now().Format(time.RFC3339),
		Entities: []EntityDiff{
			{
				Action: "fakeaction1",
				Diff:   "fakediff1",
			},
			{
				Action: "fakeaction2",
				Diff:   "fakediff2",
			},
		},
	}
}
