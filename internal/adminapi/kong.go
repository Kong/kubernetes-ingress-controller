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

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	tlsutil "github.com/kong/kubernetes-ingress-controller/v3/internal/util/tls"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/metadata"
)

// KongClientNotReadyError is returned when the Kong client is not ready to be used yet.
// This can happen if the Kong Admin API is not reachable, or if it's reachable but `GET /status` does not return 200.
type KongClientNotReadyError struct {
	Err error
}

func (e KongClientNotReadyError) Error() string {
	return fmt.Sprintf("client not ready: %s", e.Err)
}

func (e KongClientNotReadyError) Unwrap() error {
	return e.Err
}

type KongGatewayUnsupportedVersionError struct {
	msg string
}

func (e KongGatewayUnsupportedVersionError) Error() string {
	return fmt.Sprintf("Kong Gateway version is not supported: %s", e.msg)
}

// NewKongAPIClient returns a Kong API client for a given root API URL.
// It ensures that proper User-Agent is set. Do not use kong.NewClient directly.
func NewKongAPIClient(adminURL string, kongAdminAPIConfig managercfg.AdminAPIClientConfig, kongAdminToken string) (*kong.Client, error) {
	httpClient, err := makeHTTPClient(kongAdminAPIConfig, kongAdminToken)
	if err != nil {
		return nil, err
	}

	client, err := kong.NewClient(kong.String(adminURL), httpClient) //nolint:forbidigo
	if err != nil {
		return nil, fmt.Errorf("creating Kong client: %w", err)
	}
	client.UserAgent = metadata.UserAgent()
	return client, nil
}

// NewKongClientForWorkspace returns a Kong API client for a given root API URL and workspace.
// It ensures that the client is ready to be used by performing a status check, returns KongClientNotReadyError if not
// or KongGatewayUnsupportedVersionError if it can't check Kong Gateway's version or it is not >= 3.4.1.
// If the workspace does not already exist, NewKongClientForWorkspace will create it.
func NewKongClientForWorkspace(
	ctx context.Context, adminURL string, wsName string, kongAdminAPIConfig managercfg.AdminAPIClientConfig, kongAdminToken string,
) (*Client, error) {
	// Create the base client, and if no workspace was provided then return that.
	client, err := NewKongAPIClient(adminURL, kongAdminAPIConfig, kongAdminToken)
	if err != nil {
		return nil, fmt.Errorf("creating Kong client: %w", err)
	}

	// Ensure that the client is ready to be used by performing a status check.
	// Only run the check when workspace is not given,
	// because the client may not be granted to call /status and only allowed to access the given workspace.
	if wsName == "" {
		if _, err := client.Status(ctx); err != nil {
			return nil, KongClientNotReadyError{Err: err}
		}
	} else {
		// If a workspace was provided, verify whether or not it exists.
		exists, err := client.Workspaces.ExistsByName(ctx, kong.String(wsName))
		if err != nil {
			return nil, fmt.Errorf("looking up workspace: %w", err)
		}

		// If the provided workspace does not exist, for convenience we create it.
		if !exists {
			workspace := kong.Workspace{
				Name: kong.String(wsName),
			}
			if _, err := client.Workspaces.Create(ctx, &workspace); err != nil {
				return nil, fmt.Errorf("creating workspace: %w", err)
			}
		}
		// Ensure that we set the workspace appropriately.
		client.SetWorkspace(wsName)

		// Now that we have set the workspace, ensure that the client is ready
		// to be used with said workspace.
		if _, err := client.Status(ctx); err != nil {
			return nil, KongClientNotReadyError{Err: err}
		}
	}

	cl := NewClient(client)

	fetchedKongVersion, err := cl.GetKongVersion(ctx)
	if err != nil {
		return nil, KongGatewayUnsupportedVersionError{msg: fmt.Sprintf("getting Kong version: %v", err)}
	}
	kongVersion, err := kong.NewVersion(fetchedKongVersion)
	if err != nil {
		return nil, KongGatewayUnsupportedVersionError{msg: fmt.Sprintf("invalid Kong version: %v", err)}
	}
	kongSemVersion := semver.Version{Major: kongVersion.Major(), Minor: kongVersion.Minor(), Patch: kongVersion.Patch()}
	if kongSemVersion.LT(versions.KICv3VersionCutoff) {
		return nil, KongGatewayUnsupportedVersionError{msg: fmt.Sprintf(
			"version: %q is not supported by Kong Kubernetes Ingress Controller in version >=3.0.0, the lowest supported version is: %q",
			kongSemVersion, versions.KICv3VersionCutoff,
		)}
	}

	return cl, nil
}

const (
	HeaderNameAdminToken = "Kong-Admin-Token"
)

// makeHTTPClient returns an HTTP client with the specified mTLS/headers configuration.
func makeHTTPClient(opts managercfg.AdminAPIClientConfig, kongAdminToken string) (*http.Client, error) {
	var tlsConfig tls.Config

	if opts.TLSSkipVerify {
		tlsConfig.InsecureSkipVerify = true //nolint:gosec
	}

	tlsConfig.ServerName = opts.TLSServerName

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
		[]byte(opts.TLSClient.Cert), opts.TLSClient.CertFile, []byte(opts.TLSClient.Key), opts.TLSClient.KeyFile,
	)
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
			return strings.HasPrefix(header, HeaderNameAdminToken+":")
		})

		if !contains {
			headers = append(headers, HeaderNameAdminToken+":"+kongAdminToken)
		}
	}
	return headers
}
