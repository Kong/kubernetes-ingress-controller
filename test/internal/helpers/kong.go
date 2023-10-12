package helpers

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/utils/kongconfig"
)

// GetKongRootConfig gets version and root configurations of Kong from / endpoint of the provided Admin API URL.
func GetKongRootConfig(proxyAdminURL *url.URL, kongTestPassword string) (map[string]any, error) {
	httpClient, err := adminapi.MakeHTTPClient(&adminapi.HTTPClientOpts{}, kongTestPassword)
	if err != nil {
		return nil, fmt.Errorf("failed creating specific HTTP client for Kong API URL: %q: %w", proxyAdminURL, err)
	}
	kc, err := kong.NewClient(lo.ToPtr(proxyAdminURL.String()), httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed creating Kong API client for URL: %q: %w", proxyAdminURL, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return kc.Root(ctx)
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
	return kongconfig.KongVersionFromRoot(jsonResp)
}

// GetKongDBMode returns kong dbmode using the provided Admin API URL.
func GetKongDBMode(proxyAdminURL *url.URL, kongTestPassword string) (string, error) {
	jsonResp, err := GetKongRootConfig(proxyAdminURL, kongTestPassword)
	if err != nil {
		return "", err
	}
	dbMode, err := kongconfig.DBModeFromRoot(jsonResp)
	if err != nil {
		return "", fmt.Errorf("%w (for URL: %s)", err, proxyAdminURL)
	}
	return dbMode, nil
}

// GetKongRouterFlavor gets router flavor of Kong using the provided Admin API URL.
func GetKongRouterFlavor(proxyAdminURL *url.URL, kongTestPassword string) (string, error) {
	jsonResp, err := GetKongRootConfig(proxyAdminURL, kongTestPassword)
	if err != nil {
		return "", err
	}
	routerFlavor, err := kongconfig.RouterFlavorFromRoot(jsonResp)
	if err != nil {
		return "", fmt.Errorf("%w (for URL: %s)", err, proxyAdminURL)
	}
	return routerFlavor, nil
}
