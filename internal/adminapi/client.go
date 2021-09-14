package adminapi

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

// HTTPClientOpts defines parameters that configure an HTTP client.
type HTTPClientOpts struct {
	TLSSkipVerify bool
	TLSServerName string
	CACertPath    string
	CACert        string
	Headers       []string
}

// MakeHTTPClient returns an HTTP client with the specified mTLS/headers configuration.
// BUG: This function overwrites the default transport and client in package http!
// This problem is being left as-is during refactoring to avoid regression of untested code.
// https://github.com/Kong/kubernetes-ingress-controller/issues/1233
func MakeHTTPClient(opts *HTTPClientOpts) (*http.Client, error) {
	defaultTransport := http.DefaultTransport.(*http.Transport)

	var tlsConfig tls.Config

	if opts.TLSSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	if opts.TLSServerName != "" {
		tlsConfig.ServerName = opts.TLSServerName
	}

	if opts.CACertPath != "" && opts.CACert != "" {
		return nil, fmt.Errorf("both --kong-admin-ca-cert-path and --kong-admin-ca-cert" +
			"are set; please remove one or the other")
	}
	if opts.CACert != "" {
		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM([]byte(opts.CACert))
		if !ok {
			// TODO give user an error to make this actionable
			return nil, fmt.Errorf("failed to load kong-admin-ca-cert")
		}
		tlsConfig.RootCAs = certPool
	}
	if opts.CACertPath != "" {
		certPath := opts.CACertPath
		certPool := x509.NewCertPool()
		cert, err := os.ReadFile(certPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read kong-admin-ca-cert from path '%s': %w", certPath, err)
		}
		ok := certPool.AppendCertsFromPEM(cert)
		if !ok {
			// TODO give user an error to make this actionable
			return nil, fmt.Errorf("failed to load kong-admin-ca-cert from path '%s'", certPath)
		}
		tlsConfig.RootCAs = certPool
	}
	defaultTransport.TLSClientConfig = tlsConfig.Clone()
	c := http.DefaultClient
	// BUG: this overwrites the DefaultClient instance!
	c.Transport = &HeaderRoundTripper{
		headers: opts.Headers,
		rt:      defaultTransport,
	}

	return c, nil
}
