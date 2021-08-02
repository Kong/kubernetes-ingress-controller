package adminapi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kong/go-kong/kong"
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
		cert, err := ioutil.ReadFile(certPath)
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

// GetKongClientForWorkspace returns a Kong API client for a given root API URL and workspace.
// If the workspace does not already exist, GetKongClientForWorkspace will create it.
func GetKongClientForWorkspace(ctx context.Context, adminURL string, wsName string,
	httpclient *http.Client) (*kong.Client, error) {
	// create the base client, and if no workspace was provided then return that.
	client, err := kong.NewClient(kong.String(adminURL), httpclient)
	if err != nil {
		return nil, fmt.Errorf("creating Kong client: %w", err)
	}
	if wsName == "" {
		return client, nil
	}

	// if a workspace was provided, verify whether or not it exists.
	exists, err := client.Workspaces.Exists(ctx, kong.String(wsName))
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

	return client, nil
}
