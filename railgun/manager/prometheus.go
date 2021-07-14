package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/mgrutils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunHTTP(log logr.Logger) {
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	mux := http.NewServeMux()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ServeHTTP(true,
			ctrlutils.PROMTHPORT, mux, stopCh,
			log,
			&wg)
	}()
}

// ServeHTTP enable HTTP Server
// - prometheus
func ServeHTTP(enableProfiling bool,
	port int,
	mux *http.ServeMux,
	stop <-chan struct{},
	logger logr.Logger,
	wg *sync.WaitGroup) {
	defer wg.Done()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/build", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(mgrutils.Version("UNKNOWN", "UNKNOWN", mgrutils.KICREPO))
		if _, err := w.Write(b); err != nil {
			logger.Error(err, " endpoint /build failed to write response.")
		}
	})

	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		if err != nil {
			logger.Error(err, "failed to send SIGTERM to self.")
		}
	})

	if enableProfiling {
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

	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      300 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	serveDone := make(chan struct{})
	var wglocal sync.WaitGroup
	wglocal.Add(1)
	go func() {
		defer wglocal.Done()
		select {
		case <-stop:
			if err := server.Shutdown(context.Background()); err != nil {
				logger.Error(err, "failed to shut down server.")
			}
		case <-serveDone:
		}
	}()
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error(err, "server stopped with err.")
		close(serveDone)
	}
	wg.Wait()
}
