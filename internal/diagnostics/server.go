package diagnostics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
)

const (
	defaultHTTPReadHeaderTimeout = 10 * time.Second

	// diagnosticConfigBufferDepth is the size of the channel buffer for receiving diagnostic
	// config dumps from the proxy sync loop. The chosen size is essentially arbitrary: we don't
	// expect that the receive end will get backlogged (it only assigns the value to a local
	// variable) but do want a small amount of leeway to account for goroutine scheduling, so it
	// is not zero.
	diagnosticConfigBufferDepth = 3

	// diffHistorySize is the number of diffs to keep in history.
	diffHistorySize = 5

	// diffHashQuery is the query string key for requesting a specific hash's diff.
	diffHashQuery = "hash"
)

// Server is an HTTP server running exposing the pprof profiling tool, and processing diagnostic dumps of Kong configurations.
type Server struct {
	logger           logr.Logger
	profilingEnabled bool
	clientDiagnostic ClientDiagnostic

	lastSuccessfulConfigDump file.Content
	lastSuccessHash          string

	lastFailedConfigDump file.Content
	lastFailedHash       string
	lastRawErrBody       []byte

	currentFallbackCacheMetadata *fallback.GeneratedCacheMetadata

	diffs diffMap

	configLock   *sync.RWMutex
	fallbackLock *sync.RWMutex
	diffLock     *sync.RWMutex
}

// ServerConfig contains configuration for the diagnostics server.
type ServerConfig struct {
	// ProfilingEnabled enables profiling endpoints.
	ProfilingEnabled bool

	// ConfigDumpsEnabled enables config dumps endpoints.
	ConfigDumpsEnabled bool

	// DumpSensitiveConfig makes config dumps to include sensitive information.
	DumpSensitiveConfig bool
}

// NewServer creates a diagnostics server ready to start listening.
func NewServer(logger logr.Logger, cfg ServerConfig) Server {
	s := Server{
		logger:           logger,
		profilingEnabled: cfg.ProfilingEnabled,
		configLock:       &sync.RWMutex{},
		fallbackLock:     &sync.RWMutex{},
		diffLock:         &sync.RWMutex{},
		diffs:            newDiffMap(diffHistorySize),
	}

	if cfg.ConfigDumpsEnabled {
		s.clientDiagnostic = ClientDiagnostic{
			DumpsIncludeSensitive: cfg.DumpSensitiveConfig,
			Configs:               make(chan ConfigDump, diagnosticConfigBufferDepth),
			FallbackCacheMetadata: make(chan fallback.GeneratedCacheMetadata, diagnosticConfigBufferDepth),
			Diffs:                 make(chan ConfigDiff, diagnosticConfigBufferDepth),
		}
	}

	return s
}

// ConfigDumps returns an object allowing dumping succeeded and failed configuration updates.
// It will return a zero value of the type in case the config dumps are not enabled.
func (s *Server) ConfigDumps() ClientDiagnostic {
	return s.clientDiagnostic
}

// Listen starts up the HTTP server and blocks until ctx expires.
func (s *Server) Listen(ctx context.Context, port int) error {
	mux := http.NewServeMux()
	if s.clientDiagnostic != (ClientDiagnostic{}) {
		s.installConfigDebugHandlers(mux)
	}
	if s.profilingEnabled {
		installProfilingHandlers(mux)
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: defaultHTTPReadHeaderTimeout,
	}
	errChan := make(chan error)

	go s.receiveDiagnostics(ctx)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.logger.Error(err, "Could not start diagnostics server")
				errChan <- err
			}
		}
	}()

	s.logger.Info("Diagnostics server is starting to listen", "addr", port)

	select {
	case <-ctx.Done():
		s.logger.Info("Shutting down diagnostics server")
		return httpServer.Shutdown(context.Background()) //nolint:contextcheck
	case err := <-errChan:
		return err
	}
}

// receiveDiagnostics watches the diagnostic update channels.
func (s *Server) receiveDiagnostics(ctx context.Context) {
	for {
		select {
		case dump := <-s.clientDiagnostic.Configs:
			s.onConfigDump(dump)
		case meta := <-s.clientDiagnostic.FallbackCacheMetadata:
			s.onFallbackCacheMetadata(meta)
		case diff := <-s.clientDiagnostic.Diffs:
			s.onDiff(diff)
		case <-ctx.Done():
			if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
				s.logger.Error(err, "Shutting down diagnostic collection: context completed with error")
				return
			}
			s.logger.V(logging.InfoLevel).Info("Shutting down diagnostic collection: context completed")
			return
		}
	}
}

func (s *Server) onConfigDump(dump ConfigDump) {
	s.configLock.Lock()
	defer s.configLock.Unlock()

	if dump.Meta.Failed {
		// If the config push failed, we need to keep the failed config dump and the raw error body.
		s.lastFailedConfigDump = dump.Config
		s.lastFailedHash = dump.Meta.Hash
		s.lastRawErrBody = dump.RawResponseBody
	} else {
		// If the config push was successful, we need to keep successful config dump and the hash.
		s.lastSuccessfulConfigDump = dump.Config
		s.lastSuccessHash = dump.Meta.Hash

		// If the regular config push was successful, we can drop the fallback cache metadata as it is
		// no longer relevant.
		if !dump.Meta.Fallback {
			s.fallbackLock.Lock()
			s.currentFallbackCacheMetadata = nil
			s.fallbackLock.Unlock()
		}
	}
}

func (s *Server) onFallbackCacheMetadata(meta fallback.GeneratedCacheMetadata) {
	s.fallbackLock.Lock()
	defer s.fallbackLock.Unlock()
	s.currentFallbackCacheMetadata = &meta
}

func (s *Server) onDiff(diff ConfigDiff) {
	s.diffLock.Lock()
	defer s.diffLock.Unlock()
	s.diffs.Update(diff)
}

// installProfilingHandlers adds the Profiling webservice to the given mux.
func installProfilingHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof", redirectTo("/debug/pprof/"))
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/heap", pprof.Index)
	mux.HandleFunc("/debug/pprof/mutex", pprof.Index)
	mux.HandleFunc("/debug/pprof/goroutine", pprof.Index)
	mux.HandleFunc("/debug/pprof/threadcreate", pprof.Index)
	mux.HandleFunc("/debug/pprof/block", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

// installConfigDebugHandlers adds the config dump webservice to the given mux.
func (s *Server) installConfigDebugHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/debug/config/successful", s.handleLastValidConfig)
	mux.HandleFunc("/debug/config/failed", s.handleLastFailedConfig)
	mux.HandleFunc("/debug/config/fallback", s.handleCurrentFallback)
	mux.HandleFunc("/debug/config/raw-error", s.handleLastErrBody)
	mux.HandleFunc("/debug/config/diff-report", s.handleDiffReport)
}

// redirectTo redirects request to a certain destination.
func redirectTo(to string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, to, http.StatusFound)
	}
}

func (s *Server) handleLastValidConfig(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	if err := json.NewEncoder(rw).Encode(
		ConfigDumpResponse{
			Config:     s.lastSuccessfulConfigDump,
			ConfigHash: s.lastSuccessHash,
		}); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleLastFailedConfig(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	if err := json.NewEncoder(rw).Encode(
		ConfigDumpResponse{
			Config:     s.lastFailedConfigDump,
			ConfigHash: s.lastFailedHash,
		}); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleCurrentFallback(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	resp := mapFallbackCacheMetadataIntoFallbackResponse(s.currentFallbackCacheMetadata)
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		rw.WriteHeader(http.StatusOK)
	}
}

func (s *Server) handleLastErrBody(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	raw := s.lastRawErrBody
	if len(raw) == 0 {
		raw = []byte("No raw error body available.\n")
	}
	if _, err := rw.Write(raw); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleDiffReport(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	s.diffLock.RLock()
	defer s.diffLock.RUnlock()

	// GDR has no notion of sensitive data, so its raw diffs will include credentials and certificates when they
	// change. We could make this fancier by walking through the entity types to exclude them if sensitive is not
	// enabled, but would need to maintain a list of such types. Filter would probably happen on the producer (DB
	// update strategy) side, since that's where we currently filter for the dump.
	if !s.clientDiagnostic.DumpsIncludeSensitive {
		if err := json.NewEncoder(rw).Encode(DiffResponse{
			Message: "diffs include sensitive data: set CONTROLLER_DUMP_SENSITIVE_CONFIG=true in environment to enable",
		}); err == nil {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if s.diffs.Len() == 0 {
		if err := json.NewEncoder(rw).Encode(DiffResponse{
			Message: "no diffs available",
		}); err == nil {
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var requestedHash string
	var message string
	requestedHashQuery := r.URL.Query()[diffHashQuery]
	if len(requestedHashQuery) == 0 {
		requestedHash = s.diffs.Latest()
	} else {
		if len(requestedHashQuery) > 1 {
			message = "this endpoint does not support requesting multiple diffs, using the first hash provided"
		}
		requestedHash = requestedHashQuery[0]
	}

	diffs, err := s.diffs.ByHash(requestedHash)
	if err != nil {
		message = err.Error()
		rw.WriteHeader(http.StatusNotFound)
	}

	response := DiffResponse{
		Message:    message,
		ConfigHash: requestedHash,
		Diffs:      diffs,
	}

	if err := json.NewEncoder(rw).Encode(response); err == nil {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
