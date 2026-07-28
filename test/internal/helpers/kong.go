package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	kftkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/kongconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

// GetKongRootConfig gets version and root configurations of Kong from / endpoint of the provided Admin API URL.
func GetKongRootConfig(ctx context.Context, proxyAdminURL *url.URL, kongTestPassword string) (map[string]any, error) {
	kc, err := adminapi.NewKongAPIClient(proxyAdminURL.String(), managercfg.AdminAPIClientConfig{}, kongTestPassword)
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
	kc, err := adminapi.NewKongAPIClient(proxyAdminURL.String(), managercfg.AdminAPIClientConfig{}, kongTestPassword)
	if err != nil {
		return nil, err
	}
	return kc.Licenses.ListAll(ctx)
}

// GetLicenseSecretFromEnv returns a secret object containing the license data from the environment variable KONG_LICENSE_DATA.
// It validates the license data and returns an error if the license is invalid or missing.
func GetLicenseSecretFromEnv() (*corev1.Secret, error) {
	licenseData := testenv.KongLicenseData()
	if err := ValidateKongLicense(licenseData); err != nil {
		return nil, err
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong-enterprise-license",
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"license": []byte(licenseData),
		},
	}, nil
}

// ValidateKongLicense validates the license data and returns an error if the license is invalid or expired.
func ValidateKongLicense(licenseData string) error {
	licenseObj := &kftkong.License{}
	if err := json.Unmarshal([]byte(licenseData), licenseObj); err != nil {
		return fmt.Errorf("invalid license JSON: %w", err)
	}

	// partialRFC3339Regex is a regex to match on timestamps that are only the date
	// portion of an RFC3339 timestamp, this is commonly used in Kong license timestamps.
	partialRFC3339Regex := regexp.MustCompile("^[0-9]+-[0-9]+-[0-9]+$")

	// validate license expiration date
	expirationDateStr := licenseObj.Data.Payload.ExpirationDate
	if partialRFC3339Regex.MatchString(expirationDateStr) {
		// allow for shorthand dates which don't match the full RFC3339 spec,
		// but assume the very beginning of the day.
		expirationDateStr = fmt.Sprintf("%sT00:00:01Z", expirationDateStr)
	}
	t, err := time.Parse(time.RFC3339, expirationDateStr)
	if err != nil {
		return fmt.Errorf("invalid license date (%s): %w", expirationDateStr, err)
	}

	// validate expiration
	if time.Now().UTC().After(t) {
		return fmt.Errorf("the provided license is expired")
	}
	return nil
}
