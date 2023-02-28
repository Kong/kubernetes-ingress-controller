package konnect

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strconv"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	tlsutil "github.com/kong/kubernetes-ingress-controller/v2/internal/util/tls"
)

// NodeAPIClient is used for sending requests to Konnect Node API.
// It can be used to register Nodes in Konnect's Runtime Groups.
type NodeAPIClient struct {
	Address        string
	RuntimeGroupID string
	Client         *http.Client
}

// KicNodeAPIPathPattern is the path pattern for KIC node operations.
var KicNodeAPIPathPattern = "%s/kic/api/runtime_groups/%s/v1/kic-nodes"

// NewNodeAPIClient creates a Konnect client.
func NewNodeAPIClient(cfg adminapi.KonnectConfig) (*NodeAPIClient, error) {
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
	defaultTransport := http.DefaultTransport.(*http.Transport)
	defaultTransport.TLSClientConfig = &tlsConfig
	c.Transport = defaultTransport

	return &NodeAPIClient{
		Address:        cfg.Address,
		RuntimeGroupID: cfg.RuntimeGroupID,
		Client:         c,
	}, nil
}

func (c *NodeAPIClient) kicNodeAPIEndpoint() string {
	return fmt.Sprintf(KicNodeAPIPathPattern, c.Address, c.RuntimeGroupID)
}

func (c *NodeAPIClient) kicNodeAPIEndpointWithNodeID(nodeID string) string {
	return fmt.Sprintf(KicNodeAPIPathPattern, c.Address, c.RuntimeGroupID) + "/" + nodeID
}

func (c *NodeAPIClient) CreateNode(ctx context.Context, req *CreateNodeRequest) (*CreateNodeResponse, error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create node request: %w", err)
	}
	reqReader := bytes.NewReader(buf)
	url := c.kicNodeAPIEndpoint()
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, reqReader)
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
	}

	resp := &CreateNodeResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON body: %w", err)
	}

	return resp, nil
}

func (c *NodeAPIClient) UpdateNode(ctx context.Context, nodeID string, req *UpdateNodeRequest) (*UpdateNodeResponse, error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update node request: %w", err)
	}
	reqReader := bytes.NewReader(buf)
	url := c.kicNodeAPIEndpointWithNodeID(nodeID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, reqReader)
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

func (c *NodeAPIClient) ListNodes(ctx context.Context, pageNumber int) (*ListNodeResponse, error) {
	url, _ := neturl.Parse(c.kicNodeAPIEndpoint())
	if pageNumber != 0 {
		q := url.Query()
		q.Set("page.number", strconv.Itoa(pageNumber))
		url.RawQuery = q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpResp, err := c.Client.Do(req)
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

// ListAllNodes call ListNodes() repeatedly to get all nodes in a runtime group.
func (c *NodeAPIClient) ListAllNodes(ctx context.Context) ([]*NodeItem, error) {
	nodes := []*NodeItem{}
	pageNum := 0
	for {
		resp, err := c.ListNodes(ctx, pageNum)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, resp.Items...)
		if resp.Page == nil || resp.Page.NextPageNum == 0 {
			return nodes, nil
		}
		// if konnect returns a non-0 NextPageNum, the node are not all listed
		// and we should start listing from the returned NextPageNum.
		pageNum = int(resp.Page.NextPageNum)
	}
}

func (c *NodeAPIClient) DeleteNode(ctx context.Context, nodeID string) error {
	url := c.kicNodeAPIEndpointWithNodeID(nodeID)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
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

// isOKStatusCode returns true if the input HTTP status code is 2xx, in [200,300).
func isOKStatusCode(code int) bool {
	return code >= 200 && code < 300
}
