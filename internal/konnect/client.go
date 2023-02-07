package konnect

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	tlsutil "github.com/kong/kubernetes-ingress-controller/v2/internal/util/tls"
)

// AdminClient is used for sending requests to Konnect APIs which are not included
// in Kong Admin APIs, like node registration APIs or runtime group operation APIs.
// TODO(naming): give a better type name to this client?
type AdminClient struct {
	Address        string
	RuntimeGroupID string
	Client         *http.Client
}

var (
	// KicAPIPathPattern is the pattern of paths to API for
	// operating runtime group with ID in AdminClient.
	KicAPIPathPattern = "%s/kic/api/runtime_groups/%s"
	// KicNodeAPIPathPattern is the path pattern for KIC node operations.
	KicNodeAPIPathPattern = "%s/kic/api/runtime_groups/%s/v1/kic-nodes"
)

// NewAdminClient creates a Konnect client.
func NewAdminClient(cfg adminapi.KonnectConfig) (*AdminClient, error) {
	tlsClientCert, err := tlsutil.ValueFromVariableOrFile([]byte(cfg.TLSClient.Cert), cfg.TLSClient.CertFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client cert: %w", err)
	}
	tlsClientKey, err := tlsutil.ValueFromVariableOrFile([]byte(cfg.TLSClient.Key), cfg.TLSClient.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client key: %w", err)
	}

	tlsConfig := tls.Config{ //nolint:gosec
		Certificates: []tls.Certificate{},
	}
	if len(tlsClientCert) > 0 && len(tlsClientKey) > 0 {
		// Read the key pair to create certificate
		cert, err := tls.X509KeyPair(tlsClientCert, tlsClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	c := &http.Client{}
	defaultTransport := http.DefaultTransport.(*http.Transport)
	defaultTransport.TLSClientConfig = &tlsConfig
	c.Transport = defaultTransport

	return &AdminClient{
		Address:        cfg.Address,
		RuntimeGroupID: cfg.RuntimeGroupID,
		Client:         c,
	}, nil
}

func (c *AdminClient) kicNodeAPIEndpoint() string {
	return fmt.Sprintf(KicNodeAPIPathPattern, c.Address, c.RuntimeGroupID)
}

func (c *AdminClient) kicNodeAPIEndpointWithNodeID(nodeID string) string {
	return fmt.Sprintf(KicNodeAPIPathPattern, c.Address, c.RuntimeGroupID) + "/" + nodeID
}

func (c *AdminClient) CreateNode(req *CreateNodeRequest) (*CreateNodeResponse, error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create node request: %w", err)
	}
	reqReader := bytes.NewReader(buf)
	url := c.kicNodeAPIEndpoint()
	httpReq, err := http.NewRequest("POST", url, reqReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if !isOKStatusCode(httpResp.StatusCode) {
		return nil, fmt.Errorf("non-success response code from Koko: %d, resp body: %s", httpResp.StatusCode, string(respBuf))
		// TODO: parse returned body to return a more detailed error
	}

	resp := &CreateNodeResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON body: %w", err)
	}

	return resp, nil
}

func (c *AdminClient) UpdateNode(nodeID string, req *UpdateNodeRequest) (*UpdateNodeResponse, error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update node request: %w", err)
	}
	reqReader := bytes.NewReader(buf)
	url := c.kicNodeAPIEndpointWithNodeID(nodeID)
	httpReq, err := http.NewRequest("PUT", url, reqReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request:%w", err)
	}
	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		err := fmt.Errorf("failed to read response body: %w", err)
		return nil, err
	}

	if !isOKStatusCode(httpResp.StatusCode) {
		return nil, fmt.Errorf("non-success response code from Koko: %d, resp body %s", httpResp.StatusCode, string(respBuf))
	}

	resp := &UpdateNodeResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON body: %w", err)
	}
	return resp, nil
}

func (c *AdminClient) ListNodes() (*ListNodeResponse, error) {
	url := c.kicNodeAPIEndpoint()
	httpResp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if !isOKStatusCode(httpResp.StatusCode) {
		return nil, fmt.Errorf("non-success response from Koko: %d, resp body %s", httpResp.StatusCode, string(respBuf))
	}

	resp := &ListNodeResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return resp, nil
}

func (c *AdminClient) DeleteNode(nodeID string) error {
	url := c.kicNodeAPIEndpointWithNodeID(nodeID)
	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request:%w", err)
	}
	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to get response: %w", err)
	}
	defer httpResp.Body.Close()

	if !isOKStatusCode(httpResp.StatusCode) {
		return fmt.Errorf("non-success response from Koko: %d", httpResp.StatusCode)
	}

	return nil
}

// returns true if the input HTTP status code is 2xx, in [200,300).
func isOKStatusCode(code int) bool {
	return code >= 200 && code < 300
}
