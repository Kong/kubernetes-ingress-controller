package multiinstance

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager"
)

const (
	// diagnosticsServerReadHeaderTimeout is the amount of time allowed to read request headers in the diagnostics server.
	diagnosticsServerReadHeaderTimeout = 5 * time.Second
)

// DiagnosticsServer is a server that provides diagnostics information for multiple instances managed by the manager.
// Each instance exposes its own diagnostics endpoints on `/{instanceID}/debug/config/` prefix. On every call to
// RegisterInstance or UnregisterInstance, the server rebuilds its mux to include the latest set of handlers.
type DiagnosticsServer struct {
	listenerPort int
	handlers     map[manager.ID]http.Handler

	muxLock sync.Mutex
	mux     *http.ServeMux
}

func NewDiagnosticsServer(listenerPort int) *DiagnosticsServer {
	return &DiagnosticsServer{
		listenerPort: listenerPort,
		handlers:     make(map[manager.ID]http.Handler),
		mux:          http.NewServeMux(),
	}
}

// Start starts the diagnostics server.
func (s *DiagnosticsServer) Start(ctx context.Context) error {
	errg, _ := errgroup.WithContext(ctx)
	errg.Go(func() error {
		server := http.Server{
			Addr:              fmt.Sprintf(":%d", s.listenerPort),
			Handler:           s,
			ReadHeaderTimeout: diagnosticsServerReadHeaderTimeout,
		}
		return server.ListenAndServe()
	})
	return errg.Wait()
}

// ServeHTTP serves the diagnostics server.
func (s *DiagnosticsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.muxLock.Lock()
	defer s.muxLock.Unlock()

	s.mux.ServeHTTP(w, r)
}

// RegisterInstance registers a new instance to the diagnostics server.
func (s *DiagnosticsServer) RegisterInstance(instanceID manager.ID, instanceDiagnosticsHandler http.Handler) {
	s.muxLock.Lock()
	defer s.muxLock.Unlock()

	s.handlers[instanceID] = instanceDiagnosticsHandler
	s.rebuildMux()
}

// UnregisterInstance unregisters an instance from the diagnostics server.
func (s *DiagnosticsServer) UnregisterInstance(instanceID manager.ID) {
	s.muxLock.Lock()
	defer s.muxLock.Unlock()

	delete(s.handlers, instanceID)
	s.rebuildMux()
}

// rebuildMux rebuilds the mux with the current handlers. It should be called with the muxLock held.
func (s *DiagnosticsServer) rebuildMux() {
	s.mux = http.NewServeMux()
	for instanceID, handler := range s.handlers {
		// It's possible an instance doesn't have a diagnostics handler. Handle that gracefully.
		if handler == nil {
			continue
		}

		prefix := fmt.Sprintf("/%s/debug/config", instanceID)
		s.mux.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, handler))
	}
}
