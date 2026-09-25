//go:build e2e_tests || istio_tests || performance_tests

package e2e

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEnterpriseWorkspace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		wantOK     bool
	}{
		{
			name:       "created response succeeds",
			statusCode: http.StatusCreated,
			wantOK:     true,
		},
		{
			name:       "conflict response also succeeds",
			statusCode: http.StatusConflict,
			wantOK:     true,
		},
		{
			name:       "unexpected response retries",
			statusCode: http.StatusBadGateway,
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/workspaces", r.URL.Path)
				require.Equal(t, "secret", r.Header.Get("Kong-Admin-Token"))
				require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, "name=workspace-name", string(body))

				w.WriteHeader(tt.statusCode)
				_, err = w.Write([]byte("response-body"))
				require.NoError(t, err)
			}))
			defer server.Close()

			ok, response, err := createEnterpriseWorkspace(
				context.Background(),
				server.Client(),
				server.URL+"/workspaces",
				"secret",
				"workspace-name",
			)
			require.NoError(t, err)
			require.Equal(t, tt.wantOK, ok)
			require.Contains(t, response, "response-body")
		})
	}
}
