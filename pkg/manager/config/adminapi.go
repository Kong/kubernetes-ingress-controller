package config

// AdminAPIClientConfig defines parameters that configure a client for Kong Admin API.
type AdminAPIClientConfig struct {
	// Disable verification of TLS certificate of Kong's Admin endpoint.
	TLSSkipVerify bool
	// SNI name to use to verify the certificate presented by Kong in TLS.
	TLSServerName string
	// Path to PEM-encoded CA certificate file to verify Kong's Admin SSL certificate.
	CACertPath string
	// PEM-encoded CA certificate to verify Kong's Admin SSL certificate.
	CACert string
	// Array of headers added to every Admin API call.
	Headers []string
	// TLSClient is TLS client config.
	TLSClient TLSClientConfig
}

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
