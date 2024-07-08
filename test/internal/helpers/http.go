package helpers

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type httpClientCfg struct {
	resolveHostTo string
	rootCAs       *x509.CertPool
}

// HTTPClientOption is a functional option for configuring the HTTP client.
type HTTPClientOption func(*httpClientCfg)

// WithResolveHostTo sets the host to resolve to (equivalent of `curl --resolve`).
func WithResolveHostTo(host string) HTTPClientOption {
	return func(opts *httpClientCfg) {
		opts.resolveHostTo = host
	}
}

// WithRootCAs sets the root CAs for the client.
func WithRootCAs(rootCAs *x509.CertPool) HTTPClientOption {
	return func(opts *httpClientCfg) {
		opts.rootCAs = rootCAs
	}
}

// DefaultHTTPClient returns a client that should be used by default in tests.
// All defaults that should be propagated to tests for use should be changed in here.
func DefaultHTTPClient(opts ...HTTPClientOption) *http.Client {
	var cfg httpClientCfg
	for _, opt := range opts {
		opt(&cfg)
	}

	tr := cleanhttp.DefaultPooledTransport()
	if cfg.rootCAs != nil {
		tr.TLSClientConfig = &tls.Config{
			RootCAs:    cfg.rootCAs,
			MinVersion: tls.VersionTLS12,
		}
	}
	// It provides the equivalent of `curl --resolve` for the client.
	if cfg.resolveHostTo != "" {
		tr.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
			return (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext(ctx, network, cfg.resolveHostTo)
		}
	}

	client := cleanhttp.DefaultClient()
	client.Transport = tr
	client.Timeout = 10 * time.Second

	return client
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
func MustHTTPRequest(t *testing.T, method string, host, path string, headers map[string]string) *http.Request {
	scheme := "http"
	if strings.HasPrefix(host, "https://") {
		scheme = "https"
		host = strings.TrimPrefix(host, "https://")
	} else if strings.HasPrefix(host, "http://") {
		scheme = "http"
		host = strings.TrimPrefix(host, "http://")
	}
	host = strings.TrimRight(host, "/")
	path = strings.TrimLeft(path, "/")
	req, err := http.NewRequest(method, fmt.Sprintf("%s://%s/%s", scheme, host, path), nil)
	require.NoError(t, err)
	for header, value := range headers {
		req.Header.Set(header, value)
	}
	return req
}

// -----------------------------------------------------------------------------
// Testing Utility Functions - Various HTTP related
// -----------------------------------------------------------------------------

// EventuallyGETPath makes a GET request to the Kong proxy multiple times until
// either the request starts to respond with the given status code and contents
// present in the response body, or until timeout occurs according to ingressWait
// time limits. This uses a "require" for the desired conditions so if this request
// doesn't eventually succeed the calling test will fail and stop.
// Parameter proxyURL is the URL of Kong Gateway proxy (set nil when it's not different
// from parameter host). Parameter host, path and headers are used to make the GET request.
// Parameter certPool is used to validate the server certificate (nil to use system one).
// Response is expected to have the given statusCode and contain the passed bodyContent.
func EventuallyGETPath(
	t *testing.T,
	proxyURL *url.URL,
	host string,
	path string,
	certPool *x509.CertPool,
	statusCode int,
	bodyContent string,
	requestHeaders map[string]string,
	waitDuration time.Duration,
	waitTick time.Duration,
	responseMatchers ...ResponseMatcher,
) {
	t.Helper()
	var clientOptions []HTTPClientOption
	if proxyURL != nil {
		clientOptions = append(clientOptions, WithResolveHostTo(proxyURL.Host))
	}
	if certPool != nil {
		clientOptions = append(clientOptions, WithRootCAs(certPool))
	}
	client := DefaultHTTPClient(clientOptions...)

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		resp, err := client.Do(MustHTTPRequest(t, http.MethodGet, host, path, requestHeaders))
		if !assert.NoError(c, err) {
			return
		}
		defer resp.Body.Close()

		if !assert.Equal(c, statusCode, resp.StatusCode) {
			return
		}
		if bodyContent == "" {
			return
		}

		b := new(bytes.Buffer)
		n, err := b.ReadFrom(resp.Body)
		if !assert.NoError(c, err) {
			return
		}
		if !assert.Greater(c, n, int64(0)) {
			return
		}
		assert.Contains(c, b.String(), bodyContent)
		for _, matcher := range responseMatchers {
			reason, ok := matcher(resp, b.String())
			assert.Truef(t, ok, "response did not match %s", reason)
		}
	}, waitDuration, waitTick)
}

// EventuallyExpectHTTP404WithNoRoute is used to check whether a given http response is (specifically) a Kong 404.
func EventuallyExpectHTTP404WithNoRoute(
	t *testing.T,
	proxyURL *url.URL,
	host string,
	path string,
	waitDuration time.Duration,
	waitTick time.Duration,
	headers map[string]string,
) {
	EventuallyGETPath(
		t,
		proxyURL,
		host,
		path,
		nil,
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
	Host        string
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
	req := MustHTTPRequest(t, cfg.Method, cfg.Host, cfg.Path, cfg.Headers)
	matchedResponseCounter = make(map[string]int)

	finished := time.NewTimer(cfg.Duration)
	defer finished.Stop()
	ticker := time.NewTicker(cfg.RequestTick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			countHTTPGetResponse(t, req, proxyURL, matchedResponseCounter, matchers...)
		case <-finished.C:
			return matchedResponseCounter
		}
	}
}

func countHTTPGetResponse(t *testing.T, req *http.Request, proxyURL *url.URL, matchCounter map[string]int, matchers ...ResponseMatcher) {
	resp, err := DefaultHTTPClient(WithResolveHostTo(proxyURL.Host)).Do(req)
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
