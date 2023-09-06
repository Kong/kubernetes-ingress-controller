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

// GetKongRootConfig gets version and root configurations of Kong from / endpoint of the provided Admin API URL.
func GetKongRootConfig(proxyAdminURL *url.URL, kongTestPassword string) (map[string]any, error) {
	req, err := http.NewRequest(http.MethodGet, proxyAdminURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request for %s: %w", proxyAdminURL, err)
	}
	req.Header.Set("kong-admin-token", kongTestPassword)
	resp, err := DefaultHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed issuing HTTP request for %s: %w", proxyAdminURL, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body from %s: %w", proxyAdminURL, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed getting kong version from %s: %s: %s", proxyAdminURL, resp.Status, body)
	}

	var jsonResp map[string]any
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return nil, fmt.Errorf("failed parsing response body from %s: %w", proxyAdminURL, err)
	}
	return jsonResp, nil
}

// GetKongVersion returns kong version using the provided Admin API URL.
func GetKongVersion(proxyAdminURL *url.URL, kongTestPassword string) (kong.Version, error) {
	if override := os.Getenv("TEST_KONG_VERSION_OVERRIDE"); len(override) > 0 {
		if _, err := kong.ParseSemanticVersion(override); err != nil {
			return kong.Version{}, err
		}
		return kong.NewVersion(override)
	}
	jsonResp, err := GetKongRootConfig(proxyAdminURL, kongTestPassword)
	if err != nil {
		return kong.Version{}, err
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
	jsonResp, err := GetKongRootConfig(proxyAdminURL, kongTestPassword)
	if err != nil {
		return "", err
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

// GetKongRouterFlavor gets router flavor of Kong using the provided Admin API URL.
func GetKongRouterFlavor(proxyAdminURL *url.URL, kongTestPassword string) (string, error) {
	const (
		// ExpressionRouterMinMajorVersion is the lowest major version of Kong that supports expression router.
		// Kong below this version supports only "traditional" router, and does not contain "router_flavor" field in root configuration.
		ExpressionRouterMinMajorVersion = 3
	)
	kongVersion, err := GetKongVersion(proxyAdminURL, kongTestPassword)
	if err != nil {
		return "", err
	}

	if kongVersion.Major() < ExpressionRouterMinMajorVersion {
		return "traditional", nil
	}

	jsonResp, err := GetKongRootConfig(proxyAdminURL, kongTestPassword)
	if err != nil {
		return "", err
	}

	rootConfig, ok := jsonResp["configuration"].(map[string]any)
	if !ok {
		return "", fmt.Errorf(
			"unexpected root configuration type %T for kong (URL: %s)",
			jsonResp["configuration"], proxyAdminURL,
		)
	}
	routerFlavor, ok := rootConfig["router_flavor"]
	if !ok {
		return "", fmt.Errorf("missing 'router_flavor' key in kong's (URL: %s) configuration", proxyAdminURL)
	}

	routerFlavorStr, ok := routerFlavor.(string)
	if !ok {
		return "", fmt.Errorf(
			"'router_flavor' key is of unexpected type - %T - in kong's (URL: %s) configuration, value: %v",
			routerFlavor, proxyAdminURL, routerFlavor,
		)
	}
	return routerFlavorStr, nil
}
