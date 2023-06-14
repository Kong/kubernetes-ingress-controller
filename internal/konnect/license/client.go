package license

import (
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

// ErrNotFound is returned when backend returns 404.
var ErrNotFound = fmt.Errorf("license not found")

// Client interacts with the Konnect license API.
type Client struct {
	address        string
	runtimeGroupID string
	httpClient     *http.Client
}

// KICLicenseAPIPathPattern is the path pattern for KIC license operations.
var KICLicenseAPIPathPattern = "%s/kic/api/runtime_groups/%s/v1/licenses"

// NewClient creates a License API Konnect client.
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
	c.Transport = transport

	return &Client{
		address:        cfg.Address,
		runtimeGroupID: cfg.RuntimeGroupID,
		httpClient:     c,
	}, nil
}

func (c *Client) kicLicenseAPIEndpoint() string {
	return fmt.Sprintf(KICLicenseAPIPathPattern, c.address, c.runtimeGroupID)
}

func (c *Client) List(ctx context.Context, pageNumber int) (*ListLicenseResponse, error) {
	// TODO this is another case where we have a pseudo-unary object. The page is always 0 in practice, but if we have
	// separate functions per entity, we end up with effectively dead code for some
	url, _ := neturl.Parse(c.kicLicenseAPIEndpoint())
	if pageNumber != 0 {
		q := url.Query()
		q.Set("page.number", strconv.Itoa(pageNumber))
		url.RawQuery = q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if httpResp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if !isOKStatusCode(httpResp.StatusCode) {
		return nil, fmt.Errorf("non-success response from Koko: %d, resp body %s", httpResp.StatusCode, string(respBuf))
	}

	resp := &ListLicenseResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}
	return resp, nil
}

// isOKStatusCode returns true if the input HTTP status code is 2xx, in [200,300).
func isOKStatusCode(code int) bool {
	return code >= 200 && code < 300
}
