package diagnostics

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/manager"
)

// Server is an HTTP server running exposing the pprof profiling tool.
type Server struct {
	Logger logr.Logger
}

// Listen starts up the HTTP server and blocks until ctx expires.
func (s *Server) Listen(ctx context.Context) error {
	mux := http.NewServeMux()
	installHandlers(mux)

	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", manager.DiagnosticsPort), Handler: mux}
	errChan := make(chan error)
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

	s.Logger.Info("diagnostics server is starting to listen", "addr", manager.DiagnosticsPort)

	select {
	case <-ctx.Done():
		s.Logger.Info("shutting down diagnostics server")
		return httpServer.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

// installHandlers adds the Profiling webservice to the given mux.
func installHandlers(mux *http.ServeMux) {
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

// redirectTo redirects request to a certain destination.
func redirectTo(to string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, to, http.StatusFound)
	}
}
