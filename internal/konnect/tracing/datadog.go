package tracing

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
)

// DoRequest is a helper function that sends an HTTP request and logs the result with DataDog trace ID.
func DoRequest(ctx context.Context, httpClient *http.Client, req *http.Request) (*http.Response, error) {
	httpResp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	logger := loggerWithDataDogTraceID(log.FromContext(ctx), httpResp)
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		// In case of a non-2xx response, log at Error level for visibility.
		logger.V(logging.ErrorLevel).Info(
			"Request failed",
			"method", req.Method,
			"url", req.URL,
			"status_code", httpResp.StatusCode,
		)
	} else {
		// For 2xx responses, log at Trace level.
		logger.V(logging.TraceLevel).Info(
			"Request completed",
			"method", req.Method,
			"url", req.URL,
			"status_code", httpResp.StatusCode,
		)
	}

	return httpResp, nil
}

// loggerWithDataDogTraceID creates a new logger with the DataDog tracing information extracted from the HTTP response's
// headers. This data is useful for correlating logs with traces and logs in DataDog.
func loggerWithDataDogTraceID(logger logr.Logger, resp *http.Response) logr.Logger {
	headersToLog := []string{
		// This one is used for indexing in DataDog and can be used to correlate logs with traces.
		"X-B3-TraceId",
		// Logging these headers as well for completeness and in case they are needed for debugging.
		"X-B3-SpanId",
		"X-Datadog-Trace-Id",
		"X-Datadog-Parent-Id",
	}
	for _, header := range headersToLog {
		if value := resp.Header.Get(header); value != "" {
			logKey := strings.ReplaceAll(strings.ToLower(header), "-", "_")
			logger = logger.WithValues(logKey, value)
		}
	}
	return logger
}
