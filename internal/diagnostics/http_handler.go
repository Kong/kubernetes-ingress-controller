package diagnostics

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/samber/mo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
)

// ConfigDiagnosticsProvider is an interface representing a provider of diagnostic information.
type ConfigDiagnosticsProvider interface {
	LastSuccessfulConfigDump() (file.Content, string, bool)
	LastFailedConfigDump() (file.Content, string, bool)
	LastErrorBody() ([]byte, bool)
	CurrentFallbackCacheMetadata() mo.Option[fallback.GeneratedCacheMetadata]
	LastConfigDiffHash() string
	ConfigDiffByHash(string) (ConfigDiff, bool)
	AvailableConfigDiffsHashes() []DiffIndex
}

// ConfigDiagnosticsHTTPHandler is a handler for the diagnostic HTTP endpoints.
type ConfigDiagnosticsHTTPHandler struct {
	diagnosticsProvider   ConfigDiagnosticsProvider
	dumpsIncludeSensitive bool
	mux                   *http.ServeMux
}

func NewConfigDiagnosticsHTTPHandler(diagnosticsProvider ConfigDiagnosticsProvider, dumpsIncludeSensitive bool) *ConfigDiagnosticsHTTPHandler {
	h := &ConfigDiagnosticsHTTPHandler{
		diagnosticsProvider:   diagnosticsProvider,
		dumpsIncludeSensitive: dumpsIncludeSensitive,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/successful", h.handleLastValidConfig)
	mux.HandleFunc("/failed", h.handleLastFailedConfig)
	mux.HandleFunc("/fallback", h.handleCurrentFallback)
	mux.HandleFunc("/raw-error", h.handleLastErrBody)
	mux.HandleFunc("/diff-report", h.handleDiffReport)

	h.mux = mux
	return h
}

func (h *ConfigDiagnosticsHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *ConfigDiagnosticsHTTPHandler) handleLastValidConfig(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	config, configHash, ok := h.diagnosticsProvider.LastSuccessfulConfigDump()
	if !ok {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	if err := json.NewEncoder(rw).Encode(
		ConfigDumpResponse{
			Config:     config,
			ConfigHash: configHash,
		}.Config); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ConfigDiagnosticsHTTPHandler) handleLastFailedConfig(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	config, configHash, ok := h.diagnosticsProvider.LastFailedConfigDump()
	if !ok {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(rw).Encode(
		ConfigDumpResponse{
			Config:     config,
			ConfigHash: configHash,
		}); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ConfigDiagnosticsHTTPHandler) handleCurrentFallback(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	fallbackCacheMeta := h.diagnosticsProvider.CurrentFallbackCacheMetadata()
	resp := mapFallbackCacheMetadataIntoFallbackResponse(fallbackCacheMeta)
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ConfigDiagnosticsHTTPHandler) handleLastErrBody(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")

	errBody, ok := h.diagnosticsProvider.LastErrorBody()
	if !ok || len(errBody) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	if _, err := rw.Write(errBody); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ConfigDiagnosticsHTTPHandler) handleDiffReport(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	// GDR has no notion of sensitive data, so its raw diffs will include credentials and certificates when they
	// change. We could make this fancier by walking through the entity types to exclude them if sensitive is not
	// enabled, but would need to maintain a list of such types. Filter would probably happen on the producer (DB
	// update strategy) side, since that'h where we currently filter for the dump.
	if !h.dumpsIncludeSensitive {
		if err := json.NewEncoder(rw).Encode(DiffResponse{
			Message: "diffs include sensitive data: set CONTROLLER_DUMP_SENSITIVE_CONFIG=true in environment to enable",
		}); err == nil {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	availableDiffs := h.diagnosticsProvider.AvailableConfigDiffsHashes()

	if len(availableDiffs) == 0 {
		if err := json.NewEncoder(rw).Encode(DiffResponse{
			Message: "no diffs available",
		}); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var requestedHash string
	var message string
	requestedHashQuery := r.URL.Query()[diffHashQuery]
	if len(requestedHashQuery) == 0 {
		requestedHash = h.diagnosticsProvider.LastConfigDiffHash()
	} else {
		if len(requestedHashQuery) > 1 {
			message = "this endpoint does not support requesting multiple diffs, using the first hash provided"
		}
		requestedHash = requestedHashQuery[0]
	}

	requestedConfigDiff, ok := h.diagnosticsProvider.ConfigDiffByHash(requestedHash)
	if !ok {
		message = fmt.Sprintf("diff with hash %q not found", requestedHash)
		rw.WriteHeader(http.StatusNotFound)
	}

	response := DiffResponse{
		Message:    message,
		ConfigHash: requestedHash,
		Timestamp:  requestedConfigDiff.Timestamp,
		Diffs:      requestedConfigDiff.Entities,
		Available:  availableDiffs,
	}

	if err := json.NewEncoder(rw).Encode(response); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
