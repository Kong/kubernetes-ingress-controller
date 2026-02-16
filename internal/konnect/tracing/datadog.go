package tracing

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
)

// ContextKey is the key to carry values to trace Konnect configuration sync in the context.
type ContextKey string

const (
	// SynchronizerIDKey is the context key for the ID of KIC's Konnect configuration synchronizer instance.
	SynchronizerIDKey ContextKey = "KonnectSynchronizerID"
	// SyncSerialNumberKey is the context key for serial number to mark a round of configuration sync to Konnect.
	SyncSerialNumberKey ContextKey = "KonnectSyncSerialNumber"
	// SyncRoundIDKey is the context key to mark the ID of a round of configuration sync.
	SyncRoundIDKey ContextKey = "KonnectSyncRoundID"
	// SyncStartTimestampKey is the context key for timestamp (in seconds) of starting a round of configuration sync.
	SyncStartTimestampKey ContextKey = "KonnectSyncStartTimestamp"
)

const (
	// InstanceIDHeader is the header to mark the ID of KIC's Konnect configuration synchronizer instance.
	InstanceIDHeader = "X-Kic-Konnect-Sync-Instance-Id"
	// SyncSerialNumberHeader is the header for serial number to mark a round of configuration sync to Konnect.
	SyncSerialNumberHeader = "X-Kic-Konnect-Sync-Serial-Number"
	// SyncStartTimestampHeader is the header for timestamp (in seconds) of starting a round of configuration sync.
	SyncStartTimestampHeader = "X-Kic-Konnect-Sync-Start-Timestamp"
	// SyncRoundIDHeader is the header to mark the ID of a round of configuration sync.
	SyncRoundIDHeader = "X-Kic-Konnect-Sync-Round-Id"

	// B3TraceIDHeader is the header used by the B3 propagation format to pass the trace ID.
	B3TraceIDHeader = "X-B3-TraceId"
	// B3SpanIDHeader is the header used by the B3 propagation format to pass the span ID.
	B3SpanIDHeader = "X-B3-SpanId"

	// DatadogTraceIDHeader is the header used by the Datadog tracing system to pass the trace ID.
	DatadogTraceIDHeader = "X-Datadog-Trace-Id"
	// DatadogParentIDHeader is the header used by the Datadog tracing system to pass the parent ID.
	DatadogParentIDHeader = "X-Datadog-Parent-Id"
)

// DoRequest is a helper function that sends an HTTP request and logs the result with DataDog trace ID.
func DoRequest(ctx context.Context, httpClient *http.Client, req *http.Request) (*http.Response, error) {
	req = addHeaderFromContext(ctx, req)

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

// addHeaderFromContext extracts the values to mark a configuration sync round in context
// and set the headers in the request for tracing.
func addHeaderFromContext(ctx context.Context, req *http.Request) *http.Request {
	if instanceID, ok := ctx.Value(SynchronizerIDKey).(string); ok {
		req.Header.Add(InstanceIDHeader, instanceID)
	}
	if syncRoundID, ok := ctx.Value(SyncRoundIDKey).(string); ok {
		req.Header.Add(SyncRoundIDHeader, syncRoundID)
	}
	if serialNumber, ok := ctx.Value(SyncSerialNumberKey).(uint32); ok {
		req.Header.Add(SyncSerialNumberHeader, strconv.FormatUint(uint64(serialNumber), 10))
	}
	if startTimestamp, ok := ctx.Value(SyncStartTimestampKey).(int64); ok {
		req.Header.Add(SyncStartTimestampHeader, strconv.FormatInt(startTimestamp, 10))
	}

	return req
}

// loggerWithDataDogTraceID creates a new logger with the DataDog tracing information extracted from the HTTP response's
// headers. This data is useful for correlating logs with traces and logs in DataDog.
func loggerWithDataDogTraceID(logger logr.Logger, resp *http.Response) logr.Logger {
	headersToLog := []string{
		// This one is used for indexing in DataDog and can be used to correlate logs with traces.
		B3TraceIDHeader,
		// Logging these headers as well for completeness and in case they are needed for debugging.
		B3SpanIDHeader,
		DatadogTraceIDHeader,
		DatadogParentIDHeader,
	}
	for _, header := range headersToLog {
		if value := resp.Header.Get(header); value != "" {
			// Convert to lower and replace dashes with underscores to make the key consistent with our logging conventions.
			logKey := strings.ReplaceAll(strings.ToLower(header), "-", "_")
			logger = logger.WithValues(logKey, value)
		}
	}
	return logger
}
