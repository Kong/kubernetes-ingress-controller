package test

import "os"

// KonnectServerURL returns the Konnect server URL to be used for Konnect API
// requests in tests and CI.
// It is driven by the TEST_KONG_KONNECT_SERVER_URL environment variable.
func KonnectServerURL() string {
	serverURL := os.Getenv("TEST_KONG_KONNECT_SERVER_URL")
	if serverURL != "" {
		return serverURL
	}
	return konnectDefaultDevServerURL
}
