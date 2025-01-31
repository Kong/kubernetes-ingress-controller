package admission

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/certwatcher"

	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

const (
	DefaultAdmissionWebhookCertPath = "/admission-webhook/tls.crt"
	DefaultAdmissionWebhookKeyPath  = "/admission-webhook/tls.key"
)

func MakeTLSServer(
	ctx context.Context,
	config *config.AdmissionServerConfig,
	handler http.Handler,
	logger logr.Logger,
) (*http.Server, error) {
	const defaultHTTPReadHeaderTimeout = 10 * time.Second
	tlsConfig, err := serverConfigToTLSConfig(ctx, config, logger)
	if err != nil {
		return nil, err
	}
	return &http.Server{
		Addr:              config.ListenAddr,
		TLSConfig:         tlsConfig,
		Handler:           handler,
		ReadHeaderTimeout: defaultHTTPReadHeaderTimeout,
	}, nil
}

func serverConfigToTLSConfig(ctx context.Context, sc *config.AdmissionServerConfig, logger logr.Logger) (*tls.Config, error) {
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

	go func() {
		if err := watcher.Start(ctx); err != nil {
			logger.Error(err, "Certificate watcher error")
		}
	}()
	return &tls.Config{
		MinVersion:     tls.VersionTLS12,
		MaxVersion:     tls.VersionTLS13,
		GetCertificate: watcher.GetCertificate,
	}, nil
}
