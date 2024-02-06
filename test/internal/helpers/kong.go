package helpers

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/utils/kongconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
)

// GetKongRootConfig gets version and root configurations of Kong from / endpoint of the provided Admin API URL.
func GetKongRootConfig(ctx context.Context, proxyAdminURL *url.URL, kongTestPassword string) (map[string]any, error) {
	httpClient, err := adminapi.MakeHTTPClient(&adminapi.HTTPClientOpts{}, kongTestPassword)
	if err != nil {
		return nil, fmt.Errorf("failed creating specific HTTP client for Kong API URL: %q: %w", proxyAdminURL, err)
	}
	kc, err := kong.NewClient(lo.ToPtr(proxyAdminURL.String()), httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed creating Kong API client for URL: %q: %w", proxyAdminURL, err)
	}
	return kc.Root(ctx)
}

// GetKongVersion returns kong version using the provided Admin API URL.
func GetKongVersion(ctx context.Context, proxyAdminURL *url.URL, kongTestPassword string) (kong.Version, error) {
	if override := os.Getenv("TEST_KONG_VERSION_OVERRIDE"); len(override) > 0 {
		if _, err := kong.ParseSemanticVersion(override); err != nil {
			return kong.Version{}, err
		}
		return kong.NewVersion(override)
	}
	jsonResp, err := GetKongRootConfig(ctx, proxyAdminURL, kongTestPassword)
	if err != nil {
		return kong.Version{}, err
	}
	return kongconfig.KongVersionFromRoot(jsonResp)
}

// ValidateMinimalSupportedKongVersion returns version of Kong Gateway running at the provided Admin API URL.
// In case the version is below the minimal supported version versions.KICv3VersionCutoff (3.4.1), it returns an error.
func ValidateMinimalSupportedKongVersion(ctx context.Context, proxyAdminURL *url.URL, kongTestPassword string) (kong.Version, error) {
	kongVersion, err := GetKongVersion(ctx, proxyAdminURL, kongTestPassword)
	if err != nil {
		return kong.Version{}, err
	}
	kongSemVersion := semver.Version{Major: kongVersion.Major(), Minor: kongVersion.Minor(), Patch: kongVersion.Patch()}
	if kongSemVersion.LT(versions.KICv3VersionCutoff) {
		return kong.Version{}, TooOldKongGatewayError{
			actualVersion:   kongSemVersion,
			expectedVersion: versions.KICv3VersionCutoff,
		}
	}
	return kongVersion, nil
}

type TooOldKongGatewayError struct {
	actualVersion   semver.Version
	expectedVersion semver.Version
}

func (e TooOldKongGatewayError) Error() string {
	return fmt.Sprintf(
		"version: %q is not supported by Kong Kubernetes Ingress Controller in version >=3.0.0, the lowest supported version is: %q",
		e.actualVersion, e.expectedVersion,
	)
}

// GetKongDBMode returns kong dbmode using the provided Admin API URL.
func GetKongDBMode(ctx context.Context, proxyAdminURL *url.URL, kongTestPassword string) (dpconf.DBMode, error) {
	jsonResp, err := GetKongRootConfig(ctx, proxyAdminURL, kongTestPassword)
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
func GetKongRouterFlavor(ctx context.Context, proxyAdminURL *url.URL, kongTestPassword string) (dpconf.RouterFlavor, error) {
	jsonResp, err := GetKongRootConfig(ctx, proxyAdminURL, kongTestPassword)
	if err != nil {
		return "", err
	}
	routerFlavor, err := kongconfig.RouterFlavorFromRoot(jsonResp)
	if err != nil {
		return "", fmt.Errorf("%w (for URL: %s)", err, proxyAdminURL)
	}
	return routerFlavor, nil
}

// GetKongLicenses fetches all licenses applied to Kong gateway.
func GetKongLicenses(ctx context.Context, proxyAdminURL *url.URL, kongTestPassword string) ([]*kong.License, error) {
	httpClient, err := adminapi.MakeHTTPClient(&adminapi.HTTPClientOpts{}, kongTestPassword)
	if err != nil {
		return nil, err
	}
	kc, err := kong.NewClient(lo.ToPtr(proxyAdminURL.String()), httpClient)
	if err != nil {
		return nil, err
	}
	return kc.Licenses.ListAll(ctx)
}
