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

type KonnectAdminClient struct {
	Address        string
	RuntimeGroupID string
	Client         *http.Client
}

func NewKonnectAdminClient(cfg adminapi.KonnectConfig) (*KonnectAdminClient, error) {
	tlsClientCert, err := tlsutil.ValueFromVariableOrFile([]byte(cfg.TLSClient.Cert), cfg.TLSClient.CertFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client cert: %w", err)
	}
	tlsClientKey, err := tlsutil.ValueFromVariableOrFile([]byte(cfg.TLSClient.Key), cfg.TLSClient.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not extract TLS client key: %w", err)
	}

	tlsConfig := tls.Config{
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

	return &KonnectAdminClient{
		Address:        cfg.Address,
		RuntimeGroupID: cfg.RuntimeGroupID,
		Client:         c,
	}, nil
}

func (c *KonnectAdminClient) CreateNode(req *CreateNodeRequest) (*CreateNodeResponse, error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create node request: %v", err)
	}
	reqReader := bytes.NewReader(buf)
	url := fmt.Sprintf("https://%s/kic/api/runtime_groups/%s/kic-nodes", c.Address, c.RuntimeGroupID)
	httpReq, err := http.NewRequest("POST", url, reqReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %v", err)
	}
	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if httpResp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-success response code from Koko: %d, resp body: %s", httpResp.StatusCode, string(respBuf))
		// TODO: parse returned body to return a more detailed error
	}

	resp := &CreateNodeResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON body: %v", err)
	}

	return resp, nil
}

func (c *KonnectAdminClient) UpdateNode(nodeID string, req *UpdateNodeRequest) (*UpdateNodeResponse, error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update node request: %v", err)
	}
	reqReader := bytes.NewReader(buf)
	url := fmt.Sprintf("https://%s/kic/api/runtime_groups/%s/kic-nodes/%s", c.Address, c.RuntimeGroupID, nodeID)
	httpReq, err := http.NewRequest("PUT", url, reqReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request:%v", err)
	}
	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %v", err)
	}
	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		err := fmt.Errorf("failed to read response body: %v", err)
		return nil, err
	}

	if httpResp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-success response code from Koko: %d, resp body %s", httpResp.StatusCode, string(respBuf))
	}

	resp := &UpdateNodeResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON body: %v", err)
	}
	return resp, nil
}

func (c *KonnectAdminClient) ListNodes() (*ListNodeResponse, error) {
	url := fmt.Sprintf("https://%s/kic/api/runtime_groups/%s/kic-nodes", c.Address, c.RuntimeGroupID)
	httpResp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %v", err)
	}

	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if httpResp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-success response from Koko: %d, resp body %s", httpResp.StatusCode, string(respBuf))
	}

	resp := &ListNodeResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return resp, nil
}

func (c *KonnectAdminClient) DeleteNode(nodeID string) error {
	url := fmt.Sprintf("https://%s/kic/api/runtime_groups/%s/kic-nodes/%s", c.Address, c.RuntimeGroupID, nodeID)
	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request:%v", err)
	}
	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to get response: %v", err)

	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode/100 != 2 {
		return fmt.Errorf("non-success response from Koko: %d", httpResp.StatusCode)
	}

	return nil
}
