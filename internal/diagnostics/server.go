package diagnostics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/go-logr/logr"
	"github.com/kong/deck/file"

	"github.com/kong/kubernetes-ingress-controller/internal/util"
)

// Server is an HTTP server running exposing the pprof profiling tool.
type Server struct {
	Logger           logr.Logger
	ProfilingEnabled bool
	ConfigDumps      util.ConfigDumpDiagnostic
}

var successfulConfigDump file.Content
var failedConfigDump file.Content

// Listen starts up the HTTP server and blocks until ctx expires.
func (s *Server) Listen(ctx context.Context, port int) error {
	mux := http.NewServeMux()
	if s.ConfigDumps != (util.ConfigDumpDiagnostic{}) {
		installDumpHandlers(mux)
	}
	if s.ProfilingEnabled {
		installProfilingHandlers(mux)
	}

	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	errChan := make(chan error)

	go func() {
		s.receiveConfig(ctx)
	}()

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			switch err {
			case http.ErrServerClosed:
				s.Logger.Info("shutting down diagnostics server")
			default:
				s.Logger.Error(err, "could not start diagnostics server")
				errChan <- err
			}
		}
	}()

	s.Logger.Info("diagnostics server is starting to listen", "addr", port)

	select {
	case <-ctx.Done():
		s.Logger.Info("shutting down diagnostics server")
		return httpServer.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

// receiveConfig watches the config update channel
func (s *Server) receiveConfig(ctx context.Context) {
	for {
		select {
		case dump := <-s.ConfigDumps.SuccessfulConfigs:
			successfulConfigDump = dump
		case dump := <-s.ConfigDumps.FailedConfigs:
			failedConfigDump = dump
		case <-ctx.Done():
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
func installDumpHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/debug/config/successful", lastConfig(&successfulConfigDump))
	mux.HandleFunc("/debug/config/failed", lastConfig(&failedConfigDump))
}

// redirectTo redirects request to a certain destination.
func redirectTo(to string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, to, http.StatusFound)
	}
}

func lastConfig(config *file.Content) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(*config)
	}
}
