package diagnostics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kong/kubernetes-ingress-controller/railgun/internal/mgrutils"
	"github.com/sirupsen/logrus"
)

const diagnosticsServerPort = 10254

type diagnosticsServer struct {
	logger logrus.FieldLogger
	mux    *http.ServeMux
}

func NewDiagnosticsServer(enableProfiling bool, log logrus.FieldLogger) *diagnosticsServer {
	mux := http.NewServeMux()
	if enableProfiling {
		mgrutils.Install(mux)
	}
	return &diagnosticsServer{
		logger: log,
		mux:    mux,
	}
}

func (s *diagnosticsServer) Start(ctx context.Context) error {
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", diagnosticsServerPort), Handler: s.mux}
	errChan := make(chan error)
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			switch err {
			case http.ErrServerClosed:
				s.logger.Info("shutting down diagnostics server")
			default:
				s.logger.Error(err, "could not start a diagnostics server")
				errChan <- err
			}
		}
	}()

	s.logger.Info("started diagonistics server at port ", diagnosticsServerPort)

	select {
	case <-ctx.Done():
		s.logger.Info("stopping down diagnostics server")
		return httpServer.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}
