package admission

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/samber/mo"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/certwatcher"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/consts"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

const (
	DefaultAdmissionWebhookCertPath = "/admission-webhook/tls.crt"
	DefaultAdmissionWebhookKeyPath  = "/admission-webhook/tls.key"
)

type Server struct {
	s           *http.Server
	certWatcher mo.Option[*certwatcher.CertWatcher]
}

func MakeTLSServer(config managercfg.AdmissionServerConfig, handler http.Handler) (*Server, error) {
	const defaultHTTPReadHeaderTimeout = 10 * time.Second

	s := &Server{}
	tlsConfig, err := s.setupTLSConfig(config)
	if err != nil {
		return nil, err
	}

	s.s = &http.Server{
		Addr:              config.ListenAddr,
		TLSConfig:         tlsConfig,
		Handler:           handler,
		ReadHeaderTimeout: defaultHTTPReadHeaderTimeout,
	}
	return s, nil
}

// Start starts the admission server and blocks until the context is done.
func (s *Server) Start(ctx context.Context) error {
	logger := ctrllog.FromContext(ctx)
	go func() {
		if err := s.s.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err, "Failed to start admission server")
		}
	}()

	if cw, ok := s.certWatcher.Get(); ok {
		go func() {
			if err := cw.Start(ctx); err != nil {
				logger.Error(err, "Failed to start CertWatcher")
			}
		}()
	}

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), consts.DefaultGracefulShutdownTimeout) //nolint:contextcheck
	defer cancel()
	return s.s.Shutdown(ctx)
}

func (s *Server) setupTLSConfig(sc managercfg.AdmissionServerConfig) (*tls.Config, error) {
	var watcher *certwatcher.CertWatcher
	var cert, key []byte
	switch {
	// the caller provided certificates via the ENV (certwatcher can't be used here)
	case sc.CertPath == "" && sc.KeyPath == "" && sc.Cert != "" && sc.Key != "":
		cert, key = []byte(sc.Cert), []byte(sc.Key)
		keyPair, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, fmt.Errorf("X509KeyPair error: %w", err)
		}
		return &tls.Config{
			MinVersion:   tls.VersionTLS12,
			MaxVersion:   tls.VersionTLS13,
			Certificates: []tls.Certificate{keyPair},
		}, nil

	// the caller provided explicit file paths to the certs, enable certwatcher for these paths
	case sc.CertPath != "" && sc.KeyPath != "" && sc.Cert == "" && sc.Key == "":
		var err error
		watcher, err = certwatcher.New(sc.CertPath, sc.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create CertWatcher: %w", err)
		}

	// the caller provided no certificate configuration, assume the default paths and enable certwatcher for them
	case sc.CertPath == "" && sc.KeyPath == "" && sc.Cert == "" && sc.Key == "":
		var err error
		watcher, err = certwatcher.New(DefaultAdmissionWebhookCertPath, DefaultAdmissionWebhookKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create CertWatcher: %w", err)
		}

	default:
		return nil, fmt.Errorf("either cert/key files OR cert/key values must be provided, or none")
	}

	// If we have a watcher, we need to keep it to run it later in Start() method.
	if watcher != nil {
		s.certWatcher = mo.Some(watcher)
	}

	return &tls.Config{
		MinVersion:     tls.VersionTLS12,
		MaxVersion:     tls.VersionTLS13,
		GetCertificate: watcher.GetCertificate,
	}, nil
}
