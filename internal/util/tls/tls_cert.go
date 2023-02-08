package tls

import (
	"crypto/tls"
	"fmt"
	"os"
)

// ExtractClientCertificates extracts tls.Certificates from TLSClientConfig.
// It returns nil in case there was no client cert and/or client key provided.
// REVIEW: in case of no certs specified, return nil, nil, OR return non-nil error, OR add a boolean return value?
func ExtractClientCertificates(cert []byte, certFile string, key []byte, keyFile string) (*tls.Certificate, error) {
	clientCert, err := ValueFromVariableOrFile(cert, certFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client cert")
	}
	clientKey, err := ValueFromVariableOrFile(key, keyFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client key")
	}

	if len(clientCert) != 0 && len(clientKey) != 0 {
		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		return &cert, nil
	}

	return nil, nil
}

// ValueFromVariableOrFile uses v value if it's not empty, and falls back to reading a file content when value is missing.
// When both are empty, nil is returned.
func ValueFromVariableOrFile(v []byte, file string) ([]byte, error) {
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
