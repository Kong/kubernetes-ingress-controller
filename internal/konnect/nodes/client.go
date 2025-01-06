package nodes

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"

	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/tracing"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/useragent"
	tlsutil "github.com/kong/kubernetes-ingress-controller/v3/internal/util/tls"
)

// Client is used for sending requests to Konnect Node API.
// It can be used to register Nodes in Konnect's Control Planes.
type Client struct {
	address        string
	controlPlaneID string
	httpClient     *http.Client
}

// KicNodeAPIPathPattern is the path pattern for KIC node operations.
var KicNodeAPIPathPattern = "%s/kic/api/control-planes/%s/v1/kic-nodes"

// NewClient creates a Node API Konnect client.
func NewClient(cfg adminapi.KonnectConfig) (*Client, error) {
	tlsConfig := tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	cert, err := tlsutil.ExtractClientCertificates([]byte(cfg.TLSClient.Cert), cfg.TLSClient.CertFile, []byte(cfg.TLSClient.Key), cfg.TLSClient.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to extract client certificates: %w", err)
	}
	if cert != nil {
		tlsConfig.Certificates = append(tlsConfig.Certificates, *cert)
	}

	c := &http.Client{}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tlsConfig
	c.Transport = useragent.NewTransport(transport)

	return &Client{
		address:        cfg.Address,
		controlPlaneID: cfg.ControlPlaneID,
		httpClient:     c,
	}, nil
}

func (c *Client) kicNodeAPIEndpoint() string {
	return fmt.Sprintf(KicNodeAPIPathPattern, c.address, c.controlPlaneID)
}

func (c *Client) kicNodeAPIEndpointWithNodeID(nodeID string) string {
	return fmt.Sprintf(KicNodeAPIPathPattern, c.address, c.controlPlaneID) + "/" + nodeID
}

func (c *Client) CreateNode(ctx context.Context, req *CreateNodeRequest) (*CreateNodeResponse, error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create node request: %w", err)
	}
	reqReader := bytes.NewReader(buf)
	url := c.kicNodeAPIEndpoint()
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, reqReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for url %s: %w", url, err)
	}
	httpResp, err := tracing.DoRequest(ctx, c.httpClient, httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from url %s: %w", url, err)
	}
	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from url %s: %w", url, err)
	}

	if !isOKStatusCode(httpResp.StatusCode) {
		return nil, fmt.Errorf("non-success response code from url %s: %d, resp body: %s", url, httpResp.StatusCode, string(respBuf))
	}

	resp := &CreateNodeResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal create node response from body %q: %w", maxFirst64Bytes(respBuf), err)
	}

	return resp, nil
}

func (c *Client) UpdateNode(ctx context.Context, nodeID string, req *UpdateNodeRequest) (*UpdateNodeResponse, error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update node request: %w", err)
	}
	reqReader := bytes.NewReader(buf)
	url := c.kicNodeAPIEndpointWithNodeID(nodeID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, reqReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for url %s: %w", url, err)
	}
	httpResp, err := tracing.DoRequest(ctx, c.httpClient, httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from url %s: %w", url, err)
	}
	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		err := fmt.Errorf("failed to read response body from url %s: %w", url, err)
		return nil, err
	}

	if !isOKStatusCode(httpResp.StatusCode) {
		return nil, fmt.Errorf("non-success response code from url %s: %d, resp body %s", url, httpResp.StatusCode, string(respBuf))
	}

	resp := &UpdateNodeResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal update node response from body %q: %w", maxFirst64Bytes(respBuf), err)
	}
	return resp, nil
}

// ListAllNodes call ListNodes() repeatedly to get all nodes in a control plane.
func (c *Client) ListAllNodes(ctx context.Context) ([]*NodeItem, error) {
	nodes := []*NodeItem{}
	var nextCursor string
	for {
		resp, err := c.listNodes(ctx, nextCursor)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, resp.Items...)
		if resp.Page == nil || !resp.Page.HasNextPage {
			return nodes, nil
		}
		// if konnect returns that there is a next page, the nodes are not all
		// listed and we should start listing from the returned NextCursor.
		nextCursor = resp.Page.NextCursor
	}
}

func (c *Client) listNodes(ctx context.Context, nextCursor string) (*ListNodeResponse, error) {
	url, _ := neturl.Parse(c.kicNodeAPIEndpoint())
	if nextCursor != "" {
		q := url.Query()
		q.Set("page.next_cursor", nextCursor)
		url.RawQuery = q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for url %s: %w", url, err)
	}

	httpResp, err := tracing.DoRequest(ctx, c.httpClient, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from url %s: %w", url, err)
	}

	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from url %s: %w", url, err)
	}

	if !isOKStatusCode(httpResp.StatusCode) {
		return nil, fmt.Errorf("non-success response from url %s: %d, resp body %s", url, httpResp.StatusCode, string(respBuf))
	}

	resp := &ListNodeResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node items from body %q: %w", maxFirst64Bytes(respBuf), err)
	}
	return resp, nil
}

func (c *Client) DeleteNode(ctx context.Context, nodeID string) error {
	url := c.kicNodeAPIEndpointWithNodeID(nodeID)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for url %s: %w", url, err)
	}
	httpResp, err := tracing.DoRequest(ctx, c.httpClient, httpReq)
	if err != nil {
		return fmt.Errorf("failed to get response from url %s: %w", url, err)
	}
	defer httpResp.Body.Close()

	if !isOKStatusCode(httpResp.StatusCode) {
		return fmt.Errorf("non-success response from url %s: %d", url, httpResp.StatusCode)
	}

	return nil
}

func (c *Client) GetNode(ctx context.Context, nodeID string) (*NodeItem, error) {
	url := c.kicNodeAPIEndpointWithNodeID(nodeID)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for url %s: %w", url, err)
	}
	httpResp, err := tracing.DoRequest(ctx, c.httpClient, httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from url %s: %w", url, err)
	}
	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from url %s: %w", url, err)
	}

	if !isOKStatusCode(httpResp.StatusCode) {
		return nil, fmt.Errorf("non-success response from url %s: %d, resp body %s", url, httpResp.StatusCode, maxFirst64Bytes(respBuf))
	}

	resp := &NodeItem{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node item from body %q: %w", maxFirst64Bytes(respBuf), err)
	}
	return resp, nil
}

// isOKStatusCode returns true if the input HTTP status code is 2xx, in [200,300).
func isOKStatusCode(code int) bool {
	return code >= 200 && code < 300
}

// maxFirst64Bytes returns the first 64 bytes of the input byte slice as a string for debug purposes.
func maxFirst64Bytes(b []byte) string {
	return string(b[:lo.Clamp(len(b), 0, 64)])
}
