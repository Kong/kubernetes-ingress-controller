package diagnostics_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
)

type MockDiagnosticsProvider struct {
	lastSuccessfulConfigDump mo.Option[file.Content]

	lastFailedConfigDump mo.Option[file.Content]
	lastErrorBody        mo.Option[[]byte]

	currentFallbackCacheMeta mo.Option[fallback.GeneratedCacheMetadata]

	lastConfigDiffHash mo.Option[string]
	configDiffsByHash  map[string]diagnostics.ConfigDiff
}

func (m MockDiagnosticsProvider) LastSuccessfulConfigDump() (file.Content, string, bool) {
	if d, ok := m.lastSuccessfulConfigDump.Get(); ok {
		return d, "success-hash", true
	}
	return file.Content{}, "", false
}

func (m MockDiagnosticsProvider) LastFailedConfigDump() (file.Content, string, bool) {
	if d, ok := m.lastFailedConfigDump.Get(); ok {
		return d, "failed-hash", true
	}
	return file.Content{}, "", false
}

func (m MockDiagnosticsProvider) LastErrorBody() ([]byte, bool) {
	if errBody, ok := m.lastErrorBody.Get(); ok {
		return errBody, true
	}
	return nil, false
}

func (m MockDiagnosticsProvider) CurrentFallbackCacheMetadata() mo.Option[fallback.GeneratedCacheMetadata] {
	return m.currentFallbackCacheMeta
}

func (m MockDiagnosticsProvider) LastConfigDiffHash() string {
	if h, ok := m.lastConfigDiffHash.Get(); ok {
		return h
	}
	return ""
}

func (m MockDiagnosticsProvider) ConfigDiffByHash(s string) (diagnostics.ConfigDiff, bool) {
	c, ok := m.configDiffsByHash[s]
	return c, ok
}

func (m MockDiagnosticsProvider) AvailableConfigDiffsHashes() []diagnostics.DiffIndex {
	var result []diagnostics.DiffIndex
	for k := range m.configDiffsByHash {
		result = append(result, diagnostics.DiffIndex{ConfigHash: k})
	}
	return result
}

func TestConfigDiagnosticsHTTPHandler(t *testing.T) {
	testCases := []struct {
		name               string
		provider           MockDiagnosticsProvider
		endpoint           string
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "no successful config dump",
			provider:           MockDiagnosticsProvider{},
			endpoint:           "/successful",
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "successful config dump",
			provider: MockDiagnosticsProvider{
				lastSuccessfulConfigDump: mo.Some(file.Content{
					FormatVersion: "2.0",
				}),
			},
			endpoint:           "/successful",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"_format_version":"2.0"}`,
		},
		{
			name:               "no failed config dump",
			provider:           MockDiagnosticsProvider{},
			endpoint:           "/failed",
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "failed config dump",
			provider: MockDiagnosticsProvider{
				lastFailedConfigDump: mo.Some(file.Content{
					FormatVersion: "2.0",
				}),
			},
			endpoint:           "/failed",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"config": {"_format_version":"2.0"}, "hash": "failed-hash"}`,
		},
		{
			name:               "no error body",
			provider:           MockDiagnosticsProvider{},
			endpoint:           "/raw-error",
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "error body",
			provider: MockDiagnosticsProvider{
				lastErrorBody: mo.Some([]byte("error body")),
			},
			endpoint:           "/raw-error",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "error body",
		},
		{
			name:               "no fallback cache metadata",
			provider:           MockDiagnosticsProvider{},
			endpoint:           "/fallback",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "fallback cache metadata",
			provider: MockDiagnosticsProvider{
				currentFallbackCacheMeta: mo.Some(fallback.GeneratedCacheMetadata{}),
			},
			endpoint:           "/fallback",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "no config diff by hash",
			provider:           MockDiagnosticsProvider{},
			endpoint:           "/diff-report?hash=missing",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"available": null, "diffs": null, "hash": "", "message": "no diffs available", "timestamp": ""}`,
		},
		{
			name: "config diff by hash",
			provider: MockDiagnosticsProvider{
				configDiffsByHash: map[string]diagnostics.ConfigDiff{
					"existing": {
						Hash: "existing",
						Entities: []diagnostics.EntityDiff{
							{
								Generated: diagnostics.GeneratedEntity{
									Name: "name",
									Kind: "kind",
								},
								Action: "action",
								Diff:   "+diff",
							},
						},
						Timestamp: "2025-01-05T00:00:00Z",
					},
				},
			},
			endpoint:           "/diff-report?hash=existing",
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
  "hash": "existing",
  "timestamp": "2025-01-05T00:00:00Z",
  "diffs": [
    {
      "kongEntity": {
        "name": "name",
        "kind": "kind"
      },
      "action": "action",
      "diff": "+diff"
    }
  ],
  "available": [
    {
      "hash": "existing",
      "timestamp": ""
    }
  ]
}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := diagnostics.NewConfigDiagnosticsHTTPHandler(tc.provider, true)
			s := httptest.NewServer(h)
			defer s.Close()
			client := s.Client()

			resp, err := client.Get(fmt.Sprintf("%s%s", s.URL, tc.endpoint))
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedResponse != "" {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				dummy := map[string]interface{}{}
				if err := json.Unmarshal(body, &dummy); err == nil {
					// If the expected response is JSON, we can use JSONEq to compare.
					require.JSONEq(t, tc.expectedResponse, string(body))
				} else {
					// Otherwise, we should compare the raw string.
					require.Equal(t, tc.expectedResponse, string(body))
				}
			}
		})
	}
}
