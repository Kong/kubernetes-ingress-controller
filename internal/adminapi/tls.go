package adminapi

import (
	"crypto/tls"
	"fmt"
	"os"
)

// TLSClientConfig contains TLS client certificate and client key to be used when connecting with Admin APIs.
// It's validated with manager.validateClientTLS before passing it further down. It guarantees that only the
// allowed combinations of variables will be passed:
// - only one of Cert / CertFile,
// - only one of Key / KeyFile,
// - if any of Cert / CertFile is set, one of Key / KeyFile has to be set,
// - if any of Key / KeyFile is set, one of Cert / CertFile has to be set.
type TLSClientConfig struct {
	// Cert is a client certificate.
	Cert string
	// CertFile is a client certificate file path.
	CertFile string

	// Key is a client key.
	Key string
	// KeyFile is a client key file path.
	KeyFile string
}

func (c TLSClientConfig) IsZero() bool {
	return c == TLSClientConfig{}
}

// extractClientCertificates extracts tls.Certificates from TLSClientConfig.
// It returns an empty slice in case there was no client cert and/or client key provided.
func extractClientCertificates(tlsClient TLSClientConfig) ([]tls.Certificate, error) {
	clientCert, err := valueFromVariableOrFile([]byte(tlsClient.Cert), tlsClient.CertFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client cert")
	}
	clientKey, err := valueFromVariableOrFile([]byte(tlsClient.Key), tlsClient.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client key")
	}

	if len(clientCert) != 0 && len(clientKey) != 0 {
		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		return []tls.Certificate{cert}, nil
	}

	return nil, nil
}

// valueFromVariableOrFile uses v value if it's not empty, and falls back to reading a file content when value is missing.
// When both are empty, nil is returned.
func valueFromVariableOrFile(v []byte, file string) ([]byte, error) {
	if len(v) > 0 {
		return v, nil
	}
	if file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		return b, nil
	}

	return nil, nil
}
