package manager

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

// The file provides a standalone health check server instead of the server
// inside controller-runtime.manager because the manager is dependent on
// initial kong clients, but we want the liveness probe be OK if
// gateway discovery enabled, but 0 ready kong gateway endpoints detected.
// https://github.com/Kong/kubernetes-ingress-controller/issues/3592

// TODO: let the manager not dependent on initial kong clients
// then we could move back to the health check server inside manager:
// https://github.com/Kong/kubernetes-ingress-controller/issues/3590

// healthCheckHandler provides health checks for
// liveness probe (/healthz) and readiness probe (/readyz).
type healthCheckServer struct {
	lock         sync.RWMutex
	healthzCheck healthz.Checker
	readyzCheck  healthz.Checker
}

// getHealthzCheck gets the checker function for liveness probe.
func (s *healthCheckServer) getHealthzCheck() healthz.Checker {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.healthzCheck
}

// getReadyzCheck gets the check function for readiness probe.
func (s *healthCheckServer) getReadyzCheck() healthz.Checker {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.readyzCheck
}

// setHealthzCheck sets the checker function for liveness probe. The old checker function is replaced.
func (s *healthCheckServer) setHealthzCheck(checker healthz.Checker) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.healthzCheck = checker
}

// setReadyzCheck sets the checker function for readiness probe. The old checker function is replaced.
func (s *healthCheckServer) setReadyzCheck(checker healthz.Checker) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.readyzCheck = checker
}

// ServeHTTP serves for liveness probe (/healthz) and readiness probe (/readyz).
func (s *healthCheckServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var check healthz.Checker
	switch req.URL.Path {
	case "/healthz", "/healthz/":
		check = s.getHealthzCheck()
	case "/readyz", "/readyz/":
		check = s.getReadyzCheck()
	}
	// checker function not set or invalid path, return 404 not found
	if check == nil {
		http.NotFoundHandler().ServeHTTP(rw, req)
		return
	}

	if err := check(req); err != nil {
		// check failed, return 500.
		http.Error(rw, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}
	// check passed, return 200 OK
	fmt.Fprint(rw, "ok")
}

// Start starts the HTTP server serving healthz and readyz endpoints in a separate goroutine.
func (s *healthCheckServer) Start(ctx context.Context, addr string, logger logr.Logger) {
	server := &http.Server{
		Addr:              addr,
		Handler:           s,
		ReadHeaderTimeout: 3 * time.Second,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				logger.Info("Healthz server closed")
			} else {
				logger.Error(err, "Healthz server failed")
			}
		}
	}()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// We don't use the original context here as it's already done.
		//nolint:contextcheck
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Error(err, "Healthz server shutdown failed")
		}
	}()
}
