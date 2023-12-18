package konnect

import (
	"os"
	"testing"
)

const (
	konnectControlPlaneAdminAPIBaseURL   = "https://us.kic.api.konghq.tech"
	konnectControlPlanesBaseURL          = "https://us.kic.api.konghq.tech/v2"
	konnectControlPlanesConfigBaseURLFmt = "https://us.api.konghq.tech/v2/control-planes/%s/"
	konnectRolesBaseURL                  = "https://global.api.konghq.tech/v2"
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
