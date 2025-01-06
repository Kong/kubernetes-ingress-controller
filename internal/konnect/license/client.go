package license

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"time"

	"github.com/samber/mo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/tracing"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/useragent"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
	tlsutil "github.com/kong/kubernetes-ingress-controller/v3/internal/util/tls"
)

// Client interacts with the Konnect license API.
type Client struct {
	address        string
	controlPlaneID string
	httpClient     *http.Client
}

// KICLicenseAPIPathPattern is the path pattern for KIC license operations.
var KICLicenseAPIPathPattern = "%s/kic/api/control-planes/%s/v1/licenses"

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
	c.Transport = useragent.NewTransport(transport)

	return &Client{
		address:        cfg.Address,
		controlPlaneID: cfg.ControlPlaneID,
		httpClient:     c,
	}, nil
}

func (c *Client) kicLicenseAPIEndpoint() string {
	return fmt.Sprintf(KICLicenseAPIPathPattern, c.address, c.controlPlaneID)
}

func (c *Client) Get(ctx context.Context) (mo.Option[license.KonnectLicense], error) {
	// Make a request to the Konnect license API to list all licenses.
	response, err := c.listLicenses(ctx)
	if err != nil {
		return mo.None[license.KonnectLicense](), fmt.Errorf("failed to list licenses: %w", err)
	}

	// Convert the response to a KonnectLicense - we're expecting only one license.
	l, err := listLicensesResponseToKonnectLicense(response)
	if err != nil {
		return mo.None[license.KonnectLicense](), fmt.Errorf("failed to convert list licenses response: %w", err)
	}

	return l, nil
}

// isOKStatusCode returns true if the input HTTP status code is 2xx, in [200,300).
func isOKStatusCode(code int) bool {
	return code >= 200 && code < 300
}

// listLicenses calls the Konnect license API to list all licenses.
func (c *Client) listLicenses(ctx context.Context) (*ListLicenseResponse, error) {
	url, _ := neturl.Parse(c.kicLicenseAPIEndpoint())
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpResp, err := tracing.DoRequest(ctx, c.httpClient, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer httpResp.Body.Close()

	respBuf, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if httpResp.StatusCode == http.StatusNotFound {
		// 404 means no license is found which is a valid response.
		return nil, nil
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

// listLicensesResponseToKonnectLicense converts a ListLicenseResponse to a KonnectLicense.
// It validates the response and returns an error if the response is invalid.
func listLicensesResponseToKonnectLicense(response *ListLicenseResponse) (mo.Option[license.KonnectLicense], error) {
	if response == nil {
		// If the response is nil, it means no license was found.
		return mo.None[license.KonnectLicense](), nil
	}
	if len(response.Items) == 0 {
		return mo.None[license.KonnectLicense](), errors.New("no license item found in response")
	}

	// We're expecting only one license.
	item := response.Items[0]
	if item.License == "" {
		return mo.None[license.KonnectLicense](), errors.New("license item has empty license")
	}
	if item.UpdatedAt == 0 {
		return mo.None[license.KonnectLicense](), errors.New("license item has empty updated_at")
	}
	if item.ID == "" {
		return mo.None[license.KonnectLicense](), errors.New("license item has empty id")
	}

	return mo.Some(license.KonnectLicense{
		ID:        item.ID,
		UpdatedAt: time.Unix(int64(item.UpdatedAt), 0),
		Payload:   item.License,
	}), nil
}
