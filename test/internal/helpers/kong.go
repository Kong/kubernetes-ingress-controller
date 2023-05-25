package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/kong/go-kong/kong"
)

// GetKongVersion returns kong version using the provided Admin API URL.
func GetKongVersion(proxyAdminURL *url.URL, kongTestPassword string) (kong.Version, error) {
	if override := os.Getenv("TEST_KONG_VERSION_OVERRIDE"); len(override) > 0 {
		_, err := kong.ParseSemanticVersion(override)
		if err != nil {
			return kong.Version{}, err
		}
		return kong.NewVersion(override)
	}

	req, err := http.NewRequest("GET", proxyAdminURL.String(), nil)
	if err != nil {
		return kong.Version{}, fmt.Errorf("failed creating request for %s: %w", proxyAdminURL, err)
	}
	req.Header.Set("kong-admin-token", kongTestPassword)
	resp, err := DefaultHTTPClient().Do(req)
	if err != nil {
		return kong.Version{}, fmt.Errorf("failed issuing HTTP request for %s: %w", proxyAdminURL, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return kong.Version{}, fmt.Errorf("failed reading response body from %s: %w", proxyAdminURL, err)
	}
	var jsonResp map[string]interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return kong.Version{}, fmt.Errorf("failed parsing response body from %s: %w", proxyAdminURL, err)
	}

	m := kong.VersionFromInfo(jsonResp)
	version, err := kong.ParseSemanticVersion(m)
	if err != nil {
		return kong.Version{}, fmt.Errorf("failed parsing kong (URL: %s) semver from body: %s: %w", proxyAdminURL, m, err)
	}
	return version, nil
}

// GetKongDBMode returns kong dbmode using the provided Admin API URL.
func GetKongDBMode(proxyAdminURL *url.URL, kongTestPassword string) (string, error) {
	req, err := http.NewRequest("GET", proxyAdminURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed creating request for %s: %w", proxyAdminURL, err)
	}
	req.Header.Set("kong-admin-token", kongTestPassword)
	resp, err := DefaultHTTPClient().Do(req)
	if err != nil {
		return "", fmt.Errorf("failed issuing HTTP request for %s: %w", proxyAdminURL, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading response body from %s: %w", proxyAdminURL, err)
	}

	var jsonResp map[string]interface{}
	if err = json.Unmarshal(body, &jsonResp); err != nil {
		return "", fmt.Errorf("failed parsing response body from %s: %w", proxyAdminURL, err)
	}

	rootConfig, ok := jsonResp["configuration"].(map[string]any)
	if !ok {
		return "", fmt.Errorf(
			"unexpected root configuration type %T for kong (URL: %s)",
			jsonResp["configuration"], proxyAdminURL,
		)
	}

	db, ok := rootConfig["database"]
	if !ok {
		return "", fmt.Errorf("missing 'database' key in kong's (URL: %s) configuration", proxyAdminURL)
	}

	dbStr, ok := db.(string)
	if !ok {
		return "", fmt.Errorf(
			"'database' key is of unexpected type - %T - in kong's (URL: %s) configuration, value: %v",
			db, proxyAdminURL, db,
		)
	}
	return dbStr, nil
}
