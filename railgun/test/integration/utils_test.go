//+build integration_tests

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Testing Utility Functions
// -----------------------------------------------------------------------------

// expect404WithNoRoute is used to check whether a given http response is (specifically) a Kong 404.
func expect404WithNoRoute(t *testing.T, proxyURL string, resp *http.Response) bool {
	if resp.StatusCode == http.StatusNotFound {
		// once the route is torn down and returning 404's, ensure that we got the expected response body back from Kong
		// Expected: {"message":"no Route matched with those values"}
		b := new(bytes.Buffer)
		_, err := b.ReadFrom(resp.Body)
		require.NoError(t, err)
		body := struct {
			Message string `json:"message"`
		}{}
		if err := json.Unmarshal(b.Bytes(), &body); err != nil {
			t.Logf("WARNING: error decoding JSON from proxy while waiting for %s: %v", proxyURL, err)
			return false
		}
		return body.Message == "no Route matched with those values"
	}
	return false
}

// determineMaxBatchSize provides a size limit for the number of resources to POST in a single second during tests, and can be overridden with an ENV var if desired.
func determineMaxBatchSize() int {
	if v := os.Getenv("KONG_BULK_TESTING_BATCH_SIZE"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("Error: invalid batch size %s: %s", v, err))
		}
		return i
	}
	return 50
}

// -----------------------------------------------------------------------------
// Test Suite Exit Codes
// -----------------------------------------------------------------------------

const (
	// ExitCodeIncompatibleOptions is a POSIX compliant exit code for the test suite to
	// indicate that some combination of provided configurations were not compatible.
	ExitCodeIncompatibleOptions = 100

	// ExitCodeInvalidOptions is a POSIX compliant exit code for the test suite to indicate
	// that some of the provided runtime options were not valid and the tests could not run.
	ExitCodeInvalidOptions = 101

	// ExitCodeCantUseExistingCluster is a POSIX compliant exit code for the test suite to
	// indicate that an existing cluster provided for the tests was not usable.
	ExitCodeCantUseExistingCluster = 101

	// ExitCodeCantCreateCluster is a POSIX compliant exit code for the test suite to indicate
	// that a failure occurred when trying to create a Kubernetes cluster to run the tests.
	ExitCodeCantCreateCluster = 102

	// ExitCodeCleanupFailed is a POSIX compliant exit code for the test suite to indicate
	// that a failure occurred during cluster cleanup.
	ExitCodeCleanupFailed = 103
)
