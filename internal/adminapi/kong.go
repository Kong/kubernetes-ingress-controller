package adminapi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	tlsutil "github.com/kong/kubernetes-ingress-controller/v2/internal/util/tls"
)

// NewKongClientForWorkspace returns a Kong API client for a given root API URL and workspace.
// If the workspace does not already exist, NewKongClientForWorkspace will create it.
func NewKongClientForWorkspace(ctx context.Context, adminURL string, wsName string,
	httpclient *http.Client,
) (*Client, error) {
	// create the base client, and if no workspace was provided then return that.
	client, err := kong.NewClient(kong.String(adminURL), httpclient)
	if err != nil {
		return nil, fmt.Errorf("creating Kong client: %w", err)
	}
	if wsName == "" {
		return NewClient(client), nil
	}

	// if a workspace was provided, verify whether or not it exists.
	exists, err := client.Workspaces.ExistsByName(ctx, kong.String(wsName))
	if err != nil {
		return nil, fmt.Errorf("looking up workspace: %w", err)
	}

	// if the provided workspace does not exist, for convenience we create it.
	if !exists {
		workspace := kong.Workspace{
			Name: kong.String(wsName),
		}
		_, err := client.Workspaces.Create(ctx, &workspace)
		if err != nil {
			return nil, fmt.Errorf("creating workspace: %w", err)
		}
	}

	// ensure that we set the workspace appropriately
	client.SetWorkspace(wsName)

	return NewClient(client), nil
}

// HTTPClientOpts defines parameters that configure an HTTP client.
type HTTPClientOpts struct {
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

const (
	headerNameAdminToken = "Kong-Admin-Token"
)

// MakeHTTPClient returns an HTTP client with the specified mTLS/headers configuration.
func MakeHTTPClient(opts *HTTPClientOpts, kongAdminToken string) (*http.Client, error) {
	var tlsConfig tls.Config

	if opts.TLSSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	if opts.TLSServerName != "" {
		tlsConfig.ServerName = opts.TLSServerName
	}

	if opts.CACertPath != "" && opts.CACert != "" {
		return nil, fmt.Errorf("both --kong-admin-ca-cert-file and --kong-admin-ca-cert are set; " +
			"please remove one or the other")
	}
	if opts.CACert != "" {
		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM([]byte(opts.CACert))
		if !ok {
			return nil, errors.New("failed to load --kong-admin-ca-cert")
		}
		tlsConfig.RootCAs = certPool
	}
	if opts.CACertPath != "" {
		certPath := opts.CACertPath
		certPool := x509.NewCertPool()
		cert, err := os.ReadFile(certPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read --kong-admin-ca-cert from path '%s': %w", certPath, err)
		}
		ok := certPool.AppendCertsFromPEM(cert)
		if !ok {
			return nil, fmt.Errorf("failed to load --kong-admin-ca-cert from path '%s'", certPath)
		}
		tlsConfig.RootCAs = certPool
	}

	clientCertificate, err := tlsutil.ExtractClientCertificates(
		[]byte(opts.TLSClient.Cert), opts.TLSClient.CertFile, []byte(opts.TLSClient.Key), opts.TLSClient.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to extract client certificates: %w", err)
	}
	if clientCertificate != nil {
		tlsConfig.Certificates = append(tlsConfig.Certificates, *clientCertificate)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tlsConfig
	return &http.Client{
		Transport: &HeaderRoundTripper{
			headers: prepareHeaders(opts.Headers, kongAdminToken),
			rt:      transport,
		},
	}, nil
}

func prepareHeaders(headers []string, kongAdminToken string) []string {
	if kongAdminToken != "" {
		contains := lo.ContainsBy(headers, func(header string) bool {
			return strings.HasPrefix(header, headerNameAdminToken+":")
		})

		if !contains {
			headers = append(headers, headerNameAdminToken+":"+kongAdminToken)
		}
	}
	return headers
}
