//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

// -----------------------------------------------------------------------------
// Testing Timeouts
// -----------------------------------------------------------------------------

const (
	// waitTick is the default timeout tick interval for checking on ingress resources.
	waitTick = time.Second * 1

	// ingressWait is the default amount of time to wait for any particular ingress resource to be provisioned.
	ingressWait = time.Minute * 3

	// httpcTimeout is the default client timeout for HTTP clients used in tests.
	httpcTimeout = time.Second * 3

	// environmentCleanupTimeout is the amount of time that will be given by the test suite to the
	// testing environment to perform its cleanup when the test suite is shutting down.
	environmentCleanupTimeout = time.Minute * 3

	// statusWait is a const duration used in test assertions like .Eventually to
	// wait for object statuses to fulfill a provided predicate.
	statusWait = time.Minute * 3
)

// -----------------------------------------------------------------------------
// Testing Variables
// -----------------------------------------------------------------------------

const (
	// ingressClass indicates the ingress class name which the tests will use for supported object reconciliation.
	ingressClass = "kongtests"

	// redisImage is Redis. Pinned because of
	// https://github.com/Kong/kubernetes-ingress-controller/issues/2735#issuecomment-1194376496 breakage.
	redisImage = "bitnami/redis:7.0.4-debian-11-r3"

	// controllerNamespace is the Kubernetes namespace where the controller is deployed.
	controllerNamespace = "kong-system"
)

var (
	// httpc is the default HTTP client to use for tests.
	httpc = http.Client{Timeout: httpcTimeout}

	// env is the primary testing environment object which includes access to the Kubernetes cluster
	// and all the addons deployed in support of the tests.
	env environments.Environment

	// proxyURL provides access to the proxy endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyURL *url.URL

	// proxyAdminURL provides access to the Admin API endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyAdminURL *url.URL

	// proxyUDPURL provides access to the UDP API endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyUDPURL *url.URL

	// clusterVersion is a convenience var where the found version of the env.Cluster is stored.
	clusterVersion semver.Version

	// runInvalidConfigTests is set to true to run the test cases including invalid test cases.
	runInvalidConfigTests bool
)

const (
	// defaultFeatureGates is the default feature gates setting that should be
	// provided if none are provided by the user. This generally includes features
	// that are innocuous, or otherwise don't actually get triggered unless the
	// user takes further action.
	defaultFeatureGates = "GatewayAlpha=true"
)

// -----------------------------------------------------------------------------
// Testing Variables - Environment Overrides
// -----------------------------------------------------------------------------

var (
	// dbmode indicates the database backend of the test cluster ("off" and "postgres" are supported).
	dbmode = os.Getenv("TEST_DATABASE_MODE")

	// clusterVersion indicates the version of Kubernetes to use for the tests (if the cluster was not provided by the caller).
	clusterVersionStr = os.Getenv("KONG_CLUSTER_VERSION")

	// existingCluster indicates whether or not the caller is providing their own cluster for running the tests.
	// These need to come in the format <TYPE>:<NAME> (e.g. "kind:<NAME>", "gke:<NAME>", e.t.c.).
	existingCluster = os.Getenv("KONG_TEST_CLUSTER")

	// keepTestCluster indicates whether the caller wants the cluster created by the test suite
	// to persist after the test for inspection. This has a nil effect when an existing cluster
	// is provided, as cleanup is not performed for existing clusters.
	keepTestCluster = os.Getenv("KONG_TEST_CLUSTER_PERSIST")

	// kongEnterpriseEnabled enables Enterprise-specific tests when set to "true".
	kongEnterpriseEnabled = os.Getenv("TEST_KONG_ENTERPRISE")

	// kongRouterFlavor configures router mode of Kong. currently supports:
	// - `traditional`
	// - `traditional_compatible`.
	kongRouterFlavor = os.Getenv("TEST_KONG_ROUTER_FLAVOR")

	// kongImage is the Kong image to use in lieu of the default.
	kongImage = os.Getenv("TEST_KONG_IMAGE")

	// kongImage is the Kong image to use in lieu of the default.
	kongTag = os.Getenv("TEST_KONG_TAG")

	// kongPullUsername is the Docker username to use for the Kong image pull secret.
	kongPullUsername = os.Getenv("TEST_KONG_PULL_USERNAME")

	// kongPullPassword is the Docker password to use for the Kong image pull secret.
	kongPullPassword = os.Getenv("TEST_KONG_PULL_PASSWORD")

	// controllerFeatureGates contains the feature gates that should be enabled
	// for test runs.
	controllerFeatureGates = os.Getenv("KONG_CONTROLLER_FEATURE_GATES")
)

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

	// ExitCodeEnvSetupFailed is a generic exit code that can be used as a fallback for general
	// problems setting up the testing environment and/or cluster.
	ExitCodeEnvSetupFailed = 104

	// ExitCodeCentCreateLogger is a POSIX compliant exit code for the test suite to indicate
	// that a failure occurred when trying to create a logger for the test suite.
	ExitCodeCantCreateLogger = 105

	// kongTestPassword is used as a password only within the context of transient integration test runs
	// and is left static to help developers debug failures in those testing environments.
	kongTestPassword = "password"
)

// -----------------------------------------------------------------------------
// Testing Utility Functions - Kong
// -----------------------------------------------------------------------------

func getKongVersion() (semver.Version, error) {
	if override := os.Getenv("TEST_KONG_VERSION_OVERRIDE"); len(override) > 0 {
		version, err := kong.ParseSemanticVersion(override)
		if err != nil {
			return semver.Version{}, err
		}
		return semver.Version{Major: version.Major(), Minor: version.Minor(), Patch: version.Patch()}, nil
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", proxyAdminURL.String(), nil)
	if err != nil {
		return semver.Version{}, err
	}
	req.Header.Set("kong-admin-token", kongTestPassword)
	resp, err := client.Do(req)
	if err != nil {
		return semver.Version{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return semver.Version{}, err
	}
	var jsonResp map[string]interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return semver.Version{}, err
	}
	version, err := kong.ParseSemanticVersion(kong.VersionFromInfo(jsonResp))
	if err != nil {
		return semver.Version{}, err
	}
	return semver.Version{Major: version.Major(), Minor: version.Minor(), Patch: version.Patch()}, nil
}

// -----------------------------------------------------------------------------
// Testing Utility Functions - Namespaces
// -----------------------------------------------------------------------------

// namespace provides the namespace provisioned for each test case given their t.Name as the "testCase".
func namespace(ctx context.Context, t *testing.T) *corev1.Namespace {
	namespace, err := clusters.GenerateNamespace(ctx, env.Cluster(), t.Name())
	require.NoError(t, err)
	t.Cleanup(func() {
		// Use context.Background() to ensure the namespace got removed when ctx
		// gets cancelled.
		assert.NoError(t, clusters.CleanupGeneratedResources(context.Background(), env.Cluster(), t.Name()))
	})

	return namespace
}

// -----------------------------------------------------------------------------
// Testing Utility Functions - HTTP Requests
// -----------------------------------------------------------------------------

// eventuallyGETPath makes a GET request to the Kong proxy multiple times until
// either the request starts to respond with the given status code and contents
// present in the response body, or until timeout occurrs according to
// ingressWait time limits. This uses only the path of for the request and does
// not pay attention to hostname or other routing rules. This uses a "require"
// for the desired conditions so if this request doesn't eventually succeed the
// calling test will fail and stop.
func eventuallyGETPath(t *testing.T, path string, statusCode int, bodyContents string, headers map[string]string) {
	req := newRequest(t, http.MethodGet, path, headers)

	require.Eventually(t, func() bool {
		resp, err := httpc.Do(req)
		if err != nil {
			t.Logf("WARNING: http request failed for GET %s/%s: %v", proxyURL, path, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == statusCode {
			if bodyContents == "" {
				return true
			}
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), bodyContents)
		}
		return false
	}, ingressWait, waitTick)
}

// responseMatcher is a function that returns match-name and whether the response
// matches the provided criteria.
type responseMatcher func(resp *http.Response, respBody string) (key string, ok bool)

// matchRespByStatusAndContent returns a responseMatcher that matches the given status code
// and body contents.
func matchRespByStatusAndContent(
	responseName string,
	expectedStatusCode int,
	expectedBodyContents string,
) responseMatcher {
	return func(resp *http.Response, respBody string) (string, bool) {
		if resp.StatusCode != expectedStatusCode {
			return responseName, false
		}
		ok := strings.Contains(respBody, expectedBodyContents)
		return responseName, ok
	}
}

type countHTTPResponsesConfig struct {
	Method      string
	Path        string
	Headers     map[string]string
	Duration    time.Duration
	RequestTick time.Duration
}

func countHTTPGetResponses(t *testing.T,
	cfg countHTTPResponsesConfig,
	matchers ...responseMatcher,
) (matchedResponseCounter map[string]int) {
	req := newRequest(t, cfg.Method, cfg.Path, cfg.Headers)
	finished := time.After(cfg.Duration)
	matchedResponseCounter = make(map[string]int)

	for {
		select {
		case <-time.Tick(cfg.RequestTick):
			countHTTPGetResponse(t, req, matchedResponseCounter, matchers...)
		case <-finished:
			return matchedResponseCounter
		}
	}
}

func countHTTPGetResponse(t *testing.T, req *http.Request, matchCounter map[string]int, matchers ...responseMatcher) {
	resp, err := httpc.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Logf("failed to read response body: %v", err)
	}

	body := string(bytes)

	for _, matcher := range matchers {
		if key, ok := matcher(resp, body); ok {
			matchCounter[key]++
			t.Logf("response %s matched", key)
			return
		}
	}
}

// distributionOfMapValues returns a map of the values in the given counter map
// and the relative frequency of each value.
func distributionOfMapValues(counter map[string]int) map[string]float64 {
	total := 0
	normalized := make(map[string]float64)

	for _, count := range counter {
		total += count
	}
	for key, count := range counter {
		normalized[key] = float64(count) / float64(total)
	}

	return normalized
}

func newRequest(t *testing.T, method, path string, headers map[string]string) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", proxyURL, path), nil)
	require.NoError(t, err)
	for header, value := range headers {
		req.Header.Set(header, value)
	}
	return req
}

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

// -----------------------------------------------------------------------------
// Test.MAIN Utility Functions
// -----------------------------------------------------------------------------

// exitOnErrWithCode is a helper function meant for us in the test.Main to simplify failing and exiting
// the tests under unrecoverable error conditions. It will also attempt to perform any cluster
// cleaning necessary before exiting.
func exitOnErrWithCode(ctx context.Context, err error, exitCode int) {
	if err == nil {
		return
	}

	fmt.Println("WARNING: failure occurred, performing test cleanup")
	if env != nil && existingCluster == "" && keepTestCluster == "" {
		ctx, cancel := context.WithTimeout(ctx, environmentCleanupTimeout)
		defer cancel()

		fmt.Printf("INFO: cluster %s is being deleted\n", env.Cluster().Name())
		if cleanupErr := env.Cleanup(ctx); cleanupErr != nil {
			err = fmt.Errorf("cleanup failed after test failure occurred CLEANUP_FAILURE=(%s): %w", cleanupErr, err)
		}
	}

	fmt.Fprintf(os.Stderr, "Error: tests failed: %s\n", err)
	os.Exit(exitCode)
}

// exitOnErr is a wrapper around exitOnErrorWithCode that defaults to using the ExitCodeEnvSetupFailed
// exit code. This function is meant for convenience to wrap errors in setup that are hard to predict.
func exitOnErr(ctx context.Context, err error) {
	if err == nil {
		return
	}
	exitOnErrWithCode(ctx, err, ExitCodeEnvSetupFailed)
}

// setup is a test helper function which:
//   - creates a cluster cleaner which will be used to to clean up test resources
//     automatically after the test finishes and creates a new namespace for the test
//     to use.
//   - creates a namespace for the provided test and adds it to the cleaner for
//     automatic cleanup using the previously created cleaner.
//
// TODO move this into a shared library https://github.com/Kong/kubernetes-testing-framework/issues/302
func setup(ctx context.Context, t *testing.T) (*corev1.Namespace, *clusters.Cleaner) {
	t.Helper()

	t.Log("performing test setup")
	cluster := env.Cluster()
	cleaner := clusters.NewCleaner(cluster)
	t.Cleanup(func() {
		// We still want to dump the diagnostics and perform the cleanup so use
		// a separate context.
		ctx := context.Background()

		helpers.DumpDiagnosticsIfFailed(ctx, t, cluster)
		assert.NoError(t, cleaner.Cleanup(ctx))
	})

	t.Log("creating a testing namespace")
	namespace, err := clusters.GenerateNamespace(ctx, cluster, labelValueForTest(t))
	require.NoError(t, err)
	cleaner.AddNamespace(namespace)
	t.Logf("created namespace %s for test case %s", namespace.Name, t.Name())

	return namespace, cleaner
}

func labelValueForTest(t *testing.T) string {
	s := strings.ReplaceAll(t.Name(), "/", ".")
	// Trim to adhere to k8s label requirements:
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set
	if len(s) > 63 {
		return s[:63]
	}
	return s
}

// -----------------------------------------------------------------------------
// Ingress Helpers
// -----------------------------------------------------------------------------

// addIngressToCleaner adds a runtime.Object to the cleanup list if it is a supported version of Ingress. It panics if the
// runtime.Object is something else.
func addIngressToCleaner(cleaner *clusters.Cleaner, obj runtime.Object) {
	switch i := obj.(type) {
	case *netv1.Ingress:
		cleaner.Add(i)
	case *netv1beta1.Ingress:
		cleaner.Add(i)
	default:
		panic(fmt.Sprintf("%s passed to addIngressToCleaner but is not an Ingress", obj.GetObjectKind()))
	}
}
