package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
)

// GetKongVersion returns kong version using the provided Admin API URL.
func GetKongVersion(proxyAdminURL *url.URL, kongTestPassword string) (semver.Version, error) {
	if override := os.Getenv("TEST_KONG_VERSION_OVERRIDE"); len(override) > 0 {
		version, err := kong.ParseSemanticVersion(override)
		if err != nil {
			return semver.Version{}, err
		}
		return semver.Version{Major: version.Major(), Minor: version.Minor(), Patch: version.Patch()}, nil
	}

	req, err := http.NewRequest("GET", proxyAdminURL.String(), nil)
	if err != nil {
		return semver.Version{}, fmt.Errorf("failed creating request for %s: %w", proxyAdminURL, err)
	}
	req.Header.Set("kong-admin-token", kongTestPassword)
	resp, err := DefaultHTTPClient().Do(req)
	if err != nil {
		return semver.Version{}, fmt.Errorf("failed issuing HTTP request for %s: %w", proxyAdminURL, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return semver.Version{}, fmt.Errorf("failed reading response body from %s: %w", proxyAdminURL, err)
	}
	var jsonResp map[string]interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return semver.Version{}, fmt.Errorf("failed parsing response body from %s: %w", proxyAdminURL, err)
	}

	m := kong.VersionFromInfo(jsonResp)
	version, err := kong.ParseSemanticVersion(m)
	if err != nil {
		return semver.Version{}, fmt.Errorf("failed parsing kong (URL: %s) semver from body: %s: %w", proxyAdminURL, m, err)
	}
	return semver.Version{Major: version.Major(), Minor: version.Minor(), Patch: version.Patch()}, nil
}
