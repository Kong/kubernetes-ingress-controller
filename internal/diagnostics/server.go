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

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

const (
	defaultHTTPReadHeaderTimeout = 10 * time.Second

	// diagnosticConfigBufferDepth is the size of the channel buffer for receiving diagnostic
	// config dumps from the proxy sync loop. The chosen size is essentially arbitrary: we don't
	// expect that the receive end will get backlogged (it only assigns the value to a local
	// variable) but do want a small amount of leeway to account for goroutine scheduling, so it
	// is not zero.
	diagnosticConfigBufferDepth = 3
)

// Server is an HTTP server running exposing the pprof profiling tool, and processing diagnostic dumps of Kong configurations.
type Server struct {
	logger           logr.Logger
	profilingEnabled bool
	configDumps      ConfigDumpDiagnostic

	successfulConfigDump file.Content
	failedConfigDump     file.Content
	problemObjects       []AffectedObject
	failedHash           string
	successHash          string
	rawErrBody           []byte
	configLock           *sync.RWMutex
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
	}

	if cfg.ConfigDumpsEnabled {
		s.configDumps = ConfigDumpDiagnostic{
			DumpsIncludeSensitive: cfg.DumpSensitiveConfig,
			Configs:               make(chan ConfigDump, diagnosticConfigBufferDepth),
		}
	}

	return s
}

// ConfigDumps returns an object allowing dumping succeeded and failed configuration updates.
// It will return a zero value of the type in case the config dumps are not enabled.
func (s *Server) ConfigDumps() ConfigDumpDiagnostic {
	return s.configDumps
}

// Listen starts up the HTTP server and blocks until ctx expires.
func (s *Server) Listen(ctx context.Context, port int) error {
	mux := http.NewServeMux()
	if s.configDumps != (ConfigDumpDiagnostic{}) {
		s.installDumpHandlers(mux)
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

	go s.receiveConfig(ctx)

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

// receiveConfig watches the config update channel.
func (s *Server) receiveConfig(ctx context.Context) {
	for {
		select {
		case dump := <-s.configDumps.Configs:
			s.configLock.Lock()
			if dump.Meta.Failed {
				s.failedConfigDump = dump.Config
				s.rawErrBody = dump.RawResponseBody
				s.problemObjects = dump.Meta.AffectedObjects
				s.failedHash = dump.Meta.Hash
			} else {
				s.successfulConfigDump = dump.Config
				s.successHash = dump.Meta.Hash
			}
			s.configLock.Unlock()
		case <-ctx.Done():
			if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
				s.logger.Error(err, "Shutting down diagnostic config collection: context completed with error")
				return
			}
			s.logger.V(util.InfoLevel).Info("Shutting down diagnostic config collection: context completed")
			return
		}
	}
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

// installDumpHandlers adds the config dump webservice to the given mux.
func (s *Server) installDumpHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/debug/config/successful", s.handleLastValidConfig)
	mux.HandleFunc("/debug/config/failed", s.handleLastFailedConfig)
	mux.HandleFunc("/debug/config/problems", s.handleLastFailedProblemObjects)
	mux.HandleFunc("/debug/config/raw-error", s.handleLastErrBody)
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
		configDumpResponse{
			Config:     s.successfulConfigDump,
			ConfigHash: s.successHash,
		}); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleLastFailedConfig(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	if err := json.NewEncoder(rw).Encode(
		configDumpResponse{
			Config:     s.failedConfigDump,
			ConfigHash: s.failedHash,
		}); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleLastFailedProblemObjects(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	if err := json.NewEncoder(rw).Encode(
		problemObjectsResponse{
			ConfigHash: s.failedHash,
			Objects:    s.problemObjects,
		}); err != nil {
		rw.WriteHeader(http.StatusOK)
	}
}

func (s *Server) handleLastErrBody(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	raw := s.rawErrBody
	if len(raw) == 0 {
		raw = []byte("No raw error body available.\n")
	}
	if _, err := rw.Write(raw); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
