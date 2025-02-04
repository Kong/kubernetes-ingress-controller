package diagnostics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-logr/logr"
)

const (
	defaultHTTPReadHeaderTimeout = 10 * time.Second
)

// Server is an HTTP server running exposing the pprof profiling tool, and processing diagnostic dumps of Kong configurations.
type Server struct {
	logger                   logr.Logger
	cfg                      ServerConfig
	configDiagnosticsHandler http.Handler
}

// ServerConfig contains configuration for the diagnostics server.
type ServerConfig struct {
	// ProfilingEnabled enables Golang profiling endpoints under /debug/pprof/ prefix.
	ProfilingEnabled bool

	// DumpSensitiveConfig makes config dumps to include sensitive information.
	DumpSensitiveConfig bool

	// ListenerPort is the port the diagnostics server will listen on.
	ListenerPort int
}

// ServerOption is a functional option for configuring the server.
type ServerOption func(*Server)

// WithConfigDiagnostics enables the config diagnostics handler (with /debug/config/ prefix).
func WithConfigDiagnostics(handler http.Handler) ServerOption {
	return func(s *Server) {
		s.configDiagnosticsHandler = handler
	}
}

// NewServer creates a diagnostics server. Listen() must be called to start the server.
func NewServer(logger logr.Logger, cfg ServerConfig, opts ...ServerOption) Server {
	s := Server{
		logger: logger,
		cfg:    cfg,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}

// Listen starts up the HTTP server and blocks until ctx expires.
func (s *Server) Listen(ctx context.Context) error {
	const (
		configPrefix = "/debug/config"
		pprofPrefix  = "/debug/pprof"
	)

	mux := http.NewServeMux()
	if h := s.configDiagnosticsHandler; h != nil {
		mux.Handle(fmt.Sprintf("%s/", configPrefix), http.StripPrefix(configPrefix, h))
	}
	if s.cfg.ProfilingEnabled {
		mux.Handle(fmt.Sprintf("%s/", pprofPrefix), http.StripPrefix(pprofPrefix, profilingHandler()))
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.cfg.ListenerPort),
		Handler:           mux,
		ReadHeaderTimeout: defaultHTTPReadHeaderTimeout,
	}
	errChan := make(chan error)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.logger.Error(err, "Could not start diagnostics server")
				errChan <- err
			}
		}
	}()

	s.logger.Info("Diagnostics server is starting to listen", "addr", s.cfg.ListenerPort)

	select {
	case <-ctx.Done():
		s.logger.Info("Shutting down diagnostics server")
		return httpServer.Shutdown(context.Background()) //nolint:contextcheck
	case err := <-errChan:
		return err
	}
}

// installProfilingHandlers adds the Profiling webservice to the given mux.
func profilingHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", pprof.Index)
	mux.HandleFunc("/heap", pprof.Index)
	mux.HandleFunc("/mutex", pprof.Index)
	mux.HandleFunc("/goroutine", pprof.Index)
	mux.HandleFunc("/threadcreate", pprof.Index)
	mux.HandleFunc("/block", pprof.Index)
	mux.HandleFunc("/cmdline", pprof.Cmdline)
	mux.HandleFunc("/profile", pprof.Profile)
	mux.HandleFunc("/symbol", pprof.Symbol)
	mux.HandleFunc("/trace", pprof.Trace)
	return mux
}
