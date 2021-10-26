package diagnostics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/go-logr/logr"
	"github.com/kong/deck/file"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// Server is an HTTP server running exposing the pprof profiling tool, and processing diagnostic dumps of Kong configurations.
type Server struct {
	Logger           logr.Logger
	ProfilingEnabled bool
	ConfigDumps      util.ConfigDumpDiagnostic
	ConfigLock       *sync.RWMutex
}

var successfulConfigDump file.Content
var failedConfigDump file.Content

// Listen starts up the HTTP server and blocks until ctx expires.
func (s *Server) Listen(ctx context.Context, port int) error {

	mux := http.NewServeMux()
	if s.ConfigDumps != (util.ConfigDumpDiagnostic{}) {
		s.installDumpHandlers(mux)
	}
	if s.ProfilingEnabled {
		installProfilingHandlers(mux)
	}

	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	errChan := make(chan error)

	go s.receiveConfig(ctx)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				s.Logger.Info("shutting down diagnostics server")
			} else {
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
		case dump := <-s.ConfigDumps.Configs:
			s.ConfigLock.Lock()
			if dump.Failed {
				failedConfigDump = dump.Config
			} else {
				successfulConfigDump = dump.Config
			}
			s.ConfigLock.Unlock()
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				s.Logger.Error(err, "shutting down diagnostic config collection: context completed with error")
			}
			s.Logger.V(util.InfoLevel).Info("shutting down diagnostic config collection: context completed")
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
	mux.HandleFunc("/debug/config/successful", s.lastConfig(&successfulConfigDump))
	mux.HandleFunc("/debug/config/failed", s.lastConfig(&failedConfigDump))
}

// redirectTo redirects request to a certain destination.
func redirectTo(to string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, to, http.StatusFound)
	}
}

func (s *Server) lastConfig(config *file.Content) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		s.ConfigLock.RLock()
		if err := json.NewEncoder(rw).Encode(*config); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		s.ConfigLock.RUnlock()
	}
}
