package admission

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/certwatcher"
)

const (
	DefaultAdmissionWebhookCertPath = "/admission-webhook/tls.crt"
	DefaultAdmissionWebhookKeyPath  = "/admission-webhook/tls.key"
)

type ServerConfig struct {
	ListenAddr string

	CertPath string
	Cert     string

	KeyPath string
	Key     string
}

// ServerConfigToTLSConfig converts a ServerConfig to a tls.Config.
// TODO: this could be handled by controller-runtime if we set its webhook.Options properly.
func ServerConfigToTLSConfig(ctx context.Context, sc *ServerConfig, log logr.Logger) (*tls.Config, error) {
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
			log.Error(err, "certificate watcher error")
		}
	}()
	return &tls.Config{
		MinVersion:     tls.VersionTLS12,
		MaxVersion:     tls.VersionTLS13,
		GetCertificate: watcher.GetCertificate,
	}, nil
}
