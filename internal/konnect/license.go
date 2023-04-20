package konnect

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strconv"
)

// AbstractLicenseService provides functions for interacting with the Konnect license API.
type AbstractLicenseAPI interface {
	// ListLicenses returns a page of licenses.
	List(ctx context.Context, pageNumber int) (*ListLicenseResponse, error)
	// TODO yet more pseudo-unary: there's no ListAll here because there's no need, and apparently not even page
	// counts in the response, so you can't really implement it even.
}

// LicenseAPIClient is used for sending requests to Konnect License API.
// It can be used to register Licenses in Konnect's Runtime Groups.
type LicenseAPIClient konnectResourceClient

// KICLicenseAPIPathPattern is the path pattern for KIC license operations.
var KICLicenseAPIPathPattern = "%s/kic/api/runtime_groups/%s/v1/licenses"

func (c *LicenseAPIClient) kicLicenseAPIEndpoint() string {
	return fmt.Sprintf(KICLicenseAPIPathPattern, c.Client.Address, c.Client.RuntimeGroupID)
}

func (c *LicenseAPIClient) kicLicenseAPIEndpointWithLicenseID(licenseID string) string {
	return fmt.Sprintf(KICLicenseAPIPathPattern, c.Client.Address, c.Client.RuntimeGroupID) + "/" + licenseID
}

func (c *LicenseAPIClient) List(ctx context.Context, pageNumber int) (*ListLicenseResponse, error) {
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

	resp := &ListLicenseResponse{}
	err = json.Unmarshal(respBuf, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return resp, nil
}
