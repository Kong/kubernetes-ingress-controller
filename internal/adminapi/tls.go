package adminapi

import (
	"crypto/tls"
	"fmt"
)

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
	clientCert, err := valueFromVariableOrFile(tlsClient.Cert, tlsClient.CertFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client cert")
	}
	clientKey, err := valueFromVariableOrFile(tlsClient.Key, tlsClient.KeyFile)
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
