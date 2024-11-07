package konnect

import (
	"os"
	"testing"

	sdkkonnectgo "github.com/Kong/sdk-konnect-go"
)

// SkipIfMissingRequiredKonnectEnvVariables skips the test if the required Konnect environment variables are missing.
func SkipIfMissingRequiredKonnectEnvVariables(t *testing.T) {
	if accessToken() == "" {
		t.Skip("missing TEST_KONG_KONNECT_ACCESS_TOKEN")
	}
}

// accessToken returns the access token to be used for Konnect API requests.
func accessToken() string {
	return os.Getenv("TEST_KONG_KONNECT_ACCESS_TOKEN")
}

// konnectRolesBaseURL returns the base URL for Konnect Roles API.
// NOTE: This is a temporary solution until we migrate all the Konnect API calls to the new SDK.
func konnectRolesBaseURL() string {
	const konnectDefaultRolesBaseURL = "https://global.api.konghq.tech/v2"
	return konnectDefaultRolesBaseURL
}

// konnectControlPlaneAdminAPIBaseURL returns the base URL for Konnect Control Plane Admin API.
// NOTE: This is a temporary solution until we migrate all the Konnect API calls to the new SDK.
func konnectControlPlaneAdminAPIBaseURL() string {
	const konnectDefaultControlPlaneAdminAPIBaseURL = "https://us.kic.api.konghq.tech"

	serverURL := os.Getenv("TEST_KONG_KONNECT_SERVER_URL")
	switch serverURL {
	case "https://eu.api.konghq.tech":
		return "https://eu.kic.api.konghq.tech"
	case "https://ap.api.konghq.tech":
		return "https://ap.kic.api.konghq.tech"
	case "https://us.api.konghq.tech":
		return konnectDefaultControlPlaneAdminAPIBaseURL
	default:
		return konnectDefaultControlPlaneAdminAPIBaseURL
	}
}

func serverURLOpt() sdkkonnectgo.SDKOption {
	const konnectDefaultSDKServerURL = "https://us.api.konghq.tech"

	serverURL := os.Getenv("TEST_KONG_KONNECT_SERVER_URL")
	if serverURL != "" {
		return sdkkonnectgo.WithServerURL(serverURL)
	}
	return sdkkonnectgo.WithServerURL(konnectDefaultSDKServerURL)
}
