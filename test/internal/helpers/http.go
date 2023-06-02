package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/require"
)

// DefaultHTTPClient returns a client that should be used by default in tests.
// All defaults that should be propagated to tests for use should be changed in here.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

// RetryableHTTPClient wraps a client with retry logic. That should be used when calling external services that might
// temporarily fail (e.g. Konnect APIs), and we don't want them to affect the test results.
func RetryableHTTPClient(base *http.Client) *http.Client {
	retryable := retryablehttp.NewClient()
	retryable.HTTPClient = base
	return retryable.StandardClient()
}

// -----------------------------------------------------------------------------
// Testing Utility Functions - HTTP Requests
// -----------------------------------------------------------------------------

// MustHTTPRequest creates a request with provided parameters and it fails the
// test that it was called in when request creation fails.
func MustHTTPRequest(t *testing.T, method string, proxyURL *url.URL, path string, headers map[string]string) *http.Request {
	proxyURLStr := proxyURL.String()
	// trim suffix '/'s of proxy URL and prefix '/'s of path to avoid duplicate /s
	proxyURLStr = strings.TrimRight(proxyURLStr, "/")
	path = strings.TrimLeft(path, "/")
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", proxyURLStr, path), nil)
	require.NoError(t, err)
	for header, value := range headers {
		req.Header.Set(header, value)
	}
	return req
}

// MustParseURL parses a string format URL to *url.URL. If error happens, fails the test.
func MustParseURL(t *testing.T, urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	require.NoErrorf(t, err, "failed to parse URL %s: %v", urlStr, err)
	return u
}

// -----------------------------------------------------------------------------
// Testing Utility Functions - Various HTTP related
// -----------------------------------------------------------------------------

// EventuallyGETPath makes a GET request to the Kong proxy multiple times until
// either the request starts to respond with the given status code and contents
// present in the response body, or until timeout occurrs according to
// ingressWait time limits. This uses only the path of for the request and does
// not pay attention to hostname or other routing rules. This uses a "require"
// for the desired conditions so if this request doesn't eventually succeed the
// calling test will fail and stop.
func EventuallyGETPath(
	t *testing.T,
	proxyURL *url.URL,
	path string,
	statusCode int,
	bodyContents string,
	headers map[string]string,
	waitDuration time.Duration,
	waitTick time.Duration,
) {
	client := DefaultHTTPClient()

	require.Eventually(t, func() bool {
		req := MustHTTPRequest(t, http.MethodGet, proxyURL, path, headers)
		resp, err := client.Do(req)
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
	}, waitDuration, waitTick)
}

// EventuallyExpectHTTP404WithNoRoute is used to check whether a given http response is (specifically) a Kong 404.
func EventuallyExpectHTTP404WithNoRoute(
	t *testing.T,
	proxyURL *url.URL,
	path string,
	waitDuration time.Duration,
	waitTick time.Duration,
	headers map[string]string,
) {
	EventuallyGETPath(
		t,
		proxyURL,
		path,
		http.StatusNotFound,
		"no Route matched with those values",
		headers,
		waitDuration,
		waitTick,
	)
}

// ResponseMatcher is a function that returns match-name and whether the response
// matches the provided criteria.
type ResponseMatcher func(resp *http.Response, respBody string) (key string, ok bool)

// MatchRespByStatusAndContent returns a responseMatcher that matches the given status code
// and body contents.
func MatchRespByStatusAndContent(
	responseName string,
	expectedStatusCode int,
	expectedBodyContents string,
) ResponseMatcher {
	return func(resp *http.Response, respBody string) (string, bool) {
		if resp.StatusCode != expectedStatusCode {
			return responseName, false
		}
		ok := strings.Contains(respBody, expectedBodyContents)
		return responseName, ok
	}
}

type CountHTTPResponsesConfig struct {
	Method      string
	Path        string
	Headers     map[string]string
	Duration    time.Duration
	RequestTick time.Duration
}

func CountHTTPGetResponses(
	t *testing.T,
	proxyURL *url.URL,
	cfg CountHTTPResponsesConfig,
	matchers ...ResponseMatcher,
) (matchedResponseCounter map[string]int) {
	req := MustHTTPRequest(t, cfg.Method, proxyURL, cfg.Path, cfg.Headers)
	matchedResponseCounter = make(map[string]int)

	finished := time.NewTimer(cfg.Duration)
	defer finished.Stop()
	ticker := time.NewTicker(cfg.RequestTick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			countHTTPGetResponse(t, req, matchedResponseCounter, matchers...)
		case <-finished.C:
			return matchedResponseCounter
		}
	}
}

func countHTTPGetResponse(t *testing.T, req *http.Request, matchCounter map[string]int, matchers ...ResponseMatcher) {
	resp, err := DefaultHTTPClient().Do(req)
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

// DistributionOfMapValues returns a map of the values in the given counter map
// and the relative frequency of each value.
func DistributionOfMapValues(counter map[string]int) map[string]float64 {
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
