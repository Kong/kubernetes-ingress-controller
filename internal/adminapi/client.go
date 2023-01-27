package adminapi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
)

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
	if opts.TLSClient.CertFile != "" && opts.TLSClient.Cert != "" {
		return nil, fmt.Errorf("both --kong-admin-tls-client-cert-file and --kong-admin-tls-client-cert are set; " +
			"please remove one or the other")
	}
	if opts.TLSClient.KeyFile != "" && opts.TLSClient.Key != "" {
		return nil, fmt.Errorf("both --kong-admin-tls-client-key-file and --kong-admin-tls-client-key are set; " +
			"please remove one or the other")
	}

	// if the caller has supplied either the cert or the key but not both, this is
	// erroneous input.
	if opts.TLSClient.Cert != "" && opts.TLSClient.Key == "" {
		return nil, fmt.Errorf("client certificate was provided, but the client key was not")
	}
	if opts.TLSClient.Key != "" && opts.TLSClient.Cert == "" {
		return nil, fmt.Errorf("client key was provided, but the client certificate was not")
	}

	var clientCert, clientKey []byte
	var err error

	// if a path to the certificate or key has been provided, retrieve the file contents
	if opts.TLSClient.CertFile != "" {
		tlsClientCertPath := opts.TLSClient.CertFile
		clientCert, err = os.ReadFile(tlsClientCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read certificate file %s: %w", tlsClientCertPath, err)
		}
	}
	if opts.TLSClient.KeyFile != "" {
		tlsClientKeyPath := opts.TLSClient.KeyFile
		clientKey, err = os.ReadFile(tlsClientKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file %s: %w", tlsClientKeyPath, err)
		}
	}
	if opts.TLSClient.Cert != "" {
		clientCert = []byte(opts.TLSClient.Cert)
	}
	if opts.TLSClient.Key != "" {
		clientKey = []byte(opts.TLSClient.Key)
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

type KonnectConfig struct {
	ConfigSynchronizationEnabled bool
	RuntimeGroup                 string
	Address                      string
	ClientTLS                    TLSClientConfig
}

func NewKongClientForKonnect(c KonnectConfig) (*Client, error) {
	tlsClientCert, err := valueFromVariableOrFile(c.ClientTLS.Cert, c.ClientTLS.CertFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client cert")
	}
	tlsClientKey, err := valueFromVariableOrFile(c.ClientTLS.Key, c.ClientTLS.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client key")
	}

	client, err := deckutils.GetKongClient(deckutils.KongClientConfig{
		Address:       fmt.Sprintf("%s/%s/%s", c.Address, "kic/api/runtime_groups", c.RuntimeGroup),
		TLSClientCert: tlsClientCert,
		TLSClientKey:  tlsClientKey,
	})
	if err != nil {
		return nil, err
	}

	// Konnect supports tags, we don't need to verify that.
	client.Tags = tagsStub{}

	return NewKonnectClient(client, c.RuntimeGroup), nil
}

type tagsStub struct{}

func (t tagsStub) Exists(context.Context) (bool, error) {
	return true, nil
}

// valueFromVariableOrFile uses v value if it's not empty, and falls back to reading a file content when value is missing.
func valueFromVariableOrFile(v string, file string) (string, error) {
	if v != "" {
		return v, nil
	}

	b, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
