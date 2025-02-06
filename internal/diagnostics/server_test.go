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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	testhelpers "github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

// TestDiagnosticsServer_ConfigDumps tests that the diagnostics server can receive and serve config dumps.
// It's primarily to test that write and read operations run simultaneously do not fall into a race condition.
func TestDiagnosticsServer_ConfigDumps(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	client, port := setupTestServer(ctx, t)
	configsCh := client.Configs

	// Use a WaitGroup to ensure that both the write and read operations are run simultaneously.
	readWriteWg := sync.WaitGroup{}
	readWriteWg.Add(2)

	// Write 1000 config dumps to the Server.
	const configDumpsToWrite = 1000
	go func() {
		readWriteWg.Done()
		readWriteWg.Wait()

		failed := false
		for range configDumpsToWrite {
			failed = !failed // Toggle failed flag.
			configsCh <- ConfigDump{
				Config:          file.Content{},
				Meta:            DumpMeta{Failed: failed},
				RawResponseBody: []byte("fake error body"),
			}
		}
	}()
	t.Log("Started writing config dumps")

	readWriteWg.Done()
	readWriteWg.Wait()

	httpClient := &http.Client{}
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%d/debug/config/successful", port))
		require.NoError(t, err)
		_ = resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = httpClient.Get(fmt.Sprintf("http://localhost:%d/debug/config/failed", port))
		require.NoError(t, err)
		_ = resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = httpClient.Get(fmt.Sprintf("http://localhost:%d/debug/config/raw-error", port))
		require.NoError(t, err)
		_ = resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	}, time.Second*5, time.Millisecond*10)
}

// TestDiagnosticsServer_Diffs tests the diff endpoint using fake data.
func TestDiagnosticsServer_Diffs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, port := setupTestServer(ctx, t)
	diffCh := client.Diffs

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

// setupTestServer sets up a diagnostics server for testing. It returns a client attached to it and the port it's running on.
func setupTestServer(ctx context.Context, t *testing.T) (Client, int) {
	diagnosticsCollector := NewCollector(logr.Discard(), managercfg.Config{
		DumpSensitiveConfig: true,
	})
	diagnosticsHandler := NewConfigDiagnosticsHTTPHandler(diagnosticsCollector, true)

	port := testhelpers.GetFreePort(t)
	t.Logf("Obtained a free port: %d", port)

	s := NewServer(logr.Discard(), ServerConfig{
		ListenerPort: port,
	}, WithConfigDiagnostics(diagnosticsHandler))
	client := diagnosticsCollector.Client()

	go func() {
		err := s.Listen(ctx)
		require.NoError(t, err)
	}()
	t.Log("Started diagnostics server")

	go func() {
		err := diagnosticsCollector.Start(ctx)
		require.NoError(t, err)
	}()
	t.Log("Started diagnostics collector")

	return client, port
}
