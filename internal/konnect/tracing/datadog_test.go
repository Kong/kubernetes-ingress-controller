package tracing_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonglil/buflogr"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/tracing"
)

func TestDoRequest(t *testing.T) {
	testCases := []struct {
		name                       string
		returnedStatusCode         int
		returnedHeaders            map[string]string
		expectedLogMessageContains []string
	}{
		{
			name:                       "2xx response",
			returnedStatusCode:         http.StatusOK,
			expectedLogMessageContains: []string{"V[2] Request completed"},
		},
		{
			name:                       "non-2xx response",
			returnedStatusCode:         http.StatusNotFound,
			expectedLogMessageContains: []string{"INFO Request failed"},
		},
		{
			name:               "2xx response with tracing headers",
			returnedStatusCode: http.StatusOK,
			returnedHeaders: map[string]string{
				tracing.DatadogTraceIDHeader:  "123",
				tracing.DatadogParentIDHeader: "456",
				tracing.B3TraceIDHeader:       "789",
				tracing.B3SpanIDHeader:        "012",
			},
			expectedLogMessageContains: []string{
				"V[2] Request completed",
				"x_datadog_trace_id 123",
				"x_datadog_parent_id 456",
				"x_b3_traceid 789",
				"x_b3_spanid 012",
			},
		},
		{
			name:               "non-2xx response with tracing headers",
			returnedStatusCode: http.StatusBadRequest,
			returnedHeaders: map[string]string{
				tracing.DatadogTraceIDHeader:  "123",
				tracing.DatadogParentIDHeader: "456",
				tracing.B3TraceIDHeader:       "789",
				tracing.B3SpanIDHeader:        "012",
			},
			expectedLogMessageContains: []string{
				"INFO Request failed",
				"x_datadog_trace_id 123",
				"x_datadog_parent_id 456",
				"x_b3_traceid 789",
				"x_b3_spanid 012",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				for k, v := range tc.returnedHeaders {
					w.Header().Set(k, v)
				}
				w.WriteHeader(tc.returnedStatusCode)
			}))
			t.Cleanup(testServer.Close)

			client := testServer.Client()
			request, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
			require.NoError(t, err)

			loggerBuf := &bytes.Buffer{}
			logger := buflogr.NewWithBuffer(loggerBuf)
			ctx := log.IntoContext(context.Background(), logger)

			resp, err := tracing.DoRequest(ctx, client, request)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = resp.Body.Close()
			})
			require.Equal(t, tc.returnedStatusCode, resp.StatusCode)

			for _, expectedSubstring := range tc.expectedLogMessageContains {
				require.Contains(t, loggerBuf.String(), expectedSubstring)
			}
		})
	}
}
