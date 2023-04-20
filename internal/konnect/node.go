package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strconv"
)

// AbstractNodeAPI provides functions for interacting with the Konnect node API.
type AbstractNodeAPI interface {
	CreateNode(ctx context.Context, req *CreateNodeRequest) (*CreateNodeResponse, error)
	UpdateNode(ctx context.Context, nodeID string, req *UpdateNodeRequest) (*UpdateNodeResponse, error)
	ListNodes(ctx context.Context, pageNumber int) (*ListNodeResponse, error)
	ListAllNodes(ctx context.Context) ([]*NodeItem, error)
	DeleteNode(ctx context.Context, nodeID string) error
}

// NodeAPIClient is used for sending requests to Konnect Node API.
// It can be used to register Nodes in Konnect's Runtime Groups.
type NodeAPIClient konnectResourceClient

func (c *NodeAPIClient) kicNodeAPIEndpoint() string {
	return fmt.Sprintf(KICNodeAPIPathPattern, c.Client.Address, c.Client.RuntimeGroupID)
}

func (c *NodeAPIClient) kicNodeAPIEndpointWithNodeID(nodeID string) string {
	return fmt.Sprintf(KICNodeAPIPathPattern, c.Client.Address, c.Client.RuntimeGroupID) + "/" + nodeID
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
