package konnect

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	tlsutil "github.com/kong/kubernetes-ingress-controller/v2/internal/util/tls"
)

// Client is a Client for the Konnect API, scoped to a single runtime group.
type Client struct {
	Address        string
	RuntimeGroupID string
	Client         *http.Client
	Common         konnectResourceClient
	Node           AbstractNodeAPI
	License        AbstractLicenseAPI
}

type konnectResourceClient struct {
	Client *Client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

// KICNodeAPIPathPattern is the path pattern for KIC node operations.
var KICNodeAPIPathPattern = "%s/kic/api/runtime_groups/%s/v1/kic-nodes"

// NewKonnectAPIClient creates a Konnect Client.
func NewKonnectAPIClient(cfg adminapi.KonnectConfig) (*Client, error) {
	tlsConfig := tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	cert, err := tlsutil.ExtractClientCertificates([]byte(cfg.TLSClient.Cert), cfg.TLSClient.CertFile, []byte(cfg.TLSClient.Key), cfg.TLSClient.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to extract Client certificates: %w", err)
	}
	if cert != nil {
		tlsConfig.Certificates = append(tlsConfig.Certificates, *cert)
	}

	c := &http.Client{}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tlsConfig
	c.Transport = transport

	kapi := Client{
		Address:        cfg.Address,
		RuntimeGroupID: cfg.RuntimeGroupID,
		Client:         c,
	}

	kapi.Common.Client = &kapi
	kapi.Node = (*NodeAPIClient)(&kapi.Common)
	kapi.License = (*LicenseAPIClient)(&kapi.Common)

	return &Client{
		Address:        cfg.Address,
		RuntimeGroupID: cfg.RuntimeGroupID,
		Client:         c,
		Node:           &NodeAPIClient{},
	}, nil
}

// isOKStatusCode returns true if the input HTTP status code is 2xx, in [200,300).
func isOKStatusCode(code int) bool {
	return code >= 200 && code < 300
}
