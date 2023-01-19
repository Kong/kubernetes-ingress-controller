package adminapi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
)

var clientSetup sync.Mutex

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
	// mTLS client certificate file for authentication.
	TLSClientCertPath string
	// mTLS client key file for authentication.
	TLSClientCert string
	// mTLS client certificate for authentication.
	TLSClientKeyPath string
	// mTLS client key for authentication.
	TLSClientKey string
}

// MakeHTTPClient returns an HTTP client with the specified mTLS/headers configuration.
func MakeHTTPClient(opts *HTTPClientOpts) (*http.Client, error) {
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

	// don't allow the caller to specify both the literal and path versions to supply the
	// certificate and key, they must choose one or the other for each.
	if opts.TLSClientCertPath != "" && opts.TLSClientCert != "" {
		return nil, fmt.Errorf("both --kong-admin-tls-client-cert-file and --kong-admin-tls-client-cert are set; " +
			"please remove one or the other")
	}
	if opts.TLSClientKeyPath != "" && opts.TLSClientKey != "" {
		return nil, fmt.Errorf("both --kong-admin-tls-client-key-file and --kong-admin-tls-client-key are set; " +
			"please remove one or the other")
	}

	// if the caller has supplied either the cert or the key but not both, this is
	// erroneous input.
	if opts.TLSClientCert != "" && opts.TLSClientKey == "" {
		return nil, fmt.Errorf("client certificate was provided, but the client key was not")
	}
	if opts.TLSClientKey != "" && opts.TLSClientCert == "" {
		return nil, fmt.Errorf("client key was provided, but the client certificate was not")
	}

	var clientCert, clientKey []byte
	var err error

	// if a path to the certificate or key has been provided, retrieve the file contents
	if opts.TLSClientCertPath != "" {
		tlsClientCertPath := opts.TLSClientCertPath
		clientCert, err = os.ReadFile(tlsClientCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read certificate file %s: %w", tlsClientCertPath, err)
		}
	}
	if opts.TLSClientKeyPath != "" {
		tlsClientKeyPath := opts.TLSClientKeyPath
		clientKey, err = os.ReadFile(tlsClientKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file %s: %w", tlsClientKeyPath, err)
		}
	}
	if opts.TLSClientCert != "" {
		clientCert = []byte(opts.TLSClientCert)
	}
	if opts.TLSClientKey != "" {
		clientKey = []byte(opts.TLSClientKey)
	}

	if len(clientCert) != 0 && len(clientKey) != 0 {
		// Read the key pair to create certificate
		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tlsConfig
	return &http.Client{
		Transport: &HeaderRoundTripper{
			headers: opts.Headers,
			rt:      transport,
		},
	}, nil
}

type Client struct {
	*kong.Client

	isKonnect           bool
	konnectRuntimeGroup string
}

func NewClient(c *kong.Client) *Client {
	return &Client{Client: c}
}

func NewKonnectClient(c *kong.Client, runtimeGroup string) *Client {
	return &Client{
		Client:              c,
		isKonnect:           true,
		konnectRuntimeGroup: runtimeGroup,
	}
}

func (c *Client) IsKonnect() bool {
	return c.isKonnect
}

func (c *Client) KonnectRuntimeGroup() string {
	return c.konnectRuntimeGroup
}

// GetKongClientForWorkspace returns a Kong API client for a given root API URL and workspace.
// If the workspace does not already exist, GetKongClientForWorkspace will create it.
func GetKongClientForWorkspace(ctx context.Context, adminURL string, wsName string,
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
	clientSetup.Lock()
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
	clientSetup.Unlock()

	// ensure that we set the workspace appropriately
	client.SetWorkspace(wsName)

	return NewClient(client), nil
}

func NewKongClientForKonnect(token, address, runtimeGroup string) (*Client, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}
	if address == "" {
		return nil, errors.New("empty address")
	}
	if runtimeGroup == "" {
		return nil, errors.New("empty runtime group")
	}

	// TODO: when https://github.com/Kong/koko-private/pull/1384 is done we should switch to using KIC-specific
	// endpoints and pass client cert in headers.
	headers := []string{"Authorization:Bearer " + token}
	client, err := deckutils.GetKongClient(deckutils.KongClientConfig{
		Address:    address + "/konnect-api/api/runtime_groups/" + runtimeGroup,
		HTTPClient: &http.Client{},
		Debug:      false,
		Headers:    headers,
		Retryable:  true,
	})
	if err != nil {
		return nil, err
	}

	return NewKonnectClient(client, runtimeGroup), nil
}
