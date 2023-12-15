package graph_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	graph2 "github.com/dominikbraun/graph"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/graph"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

const sampleKongConfig = `
_format_version: "3.0"
ca_certificates:
- cert: |
    -----BEGIN CERTIFICATE-----
    MIIDtTCCAp2gAwIBAgIBATANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMCVVMx
    EzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xJTAj
    BgNVBAkTHDE1MCBTcGVhciBTdHJlZXQsIFN1aXRlIDE2MDAxDjAMBgNVBBETBTk0
    MTA1MRAwDgYDVQQKEwdLb25nIEhRMB4XDTIzMDkyNjA4NDQzNVoXDTI0MDkyNjA4
    NDQzNVowgYMxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYD
    VQQHEw1TYW4gRnJhbmNpc2NvMSUwIwYDVQQJExwxNTAgU3BlYXIgU3RyZWV0LCBT
    dWl0ZSAxNjAwMQ4wDAYDVQQREwU5NDEwNTEQMA4GA1UEChMHS29uZyBIUTCCASIw
    DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANcnsjbovNdW1HSKaq9ZJ9MwTP0h
    DvShbh9VldLol3Am47xBJYny10EQU4yNF7KhBjbQGAg1hhjDGMp5wPrT66syt4gZ
    IY5xW/6j4GL3E3DNfAgNo+xruEnVHjoz3z6qkt9oAC+T2Gt0BKVtPNQlUqhRBN4f
    YBYoe08K79KSJpjLf96/H8eNJmw5WDzfTH0HdNgZRmcUQfWKgE+iZzAC4ppp8vxx
    YDlXX24GN9bylWcn6TKUkSTolsxJ8mKYDR8zj4Sk6e3z9K14cdIKP3rpXIBiTIrr
    vPsNZAzWgDArTGTS13NC7IzAwkK5iCB582CGJZ8TKqrMHtE+dGwofHJ1Mw0CAwEA
    AaMyMDAwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUniBh32i4njEZO9HxGPmY
    k65EzZswDQYJKoZIhvcNAQELBQADggEBAArcmyihIjIin0nfeL6tJI08Ps2FKIQ9
    7KnKHzkQPQvqgEmmBbzX15w8/YJul7LlspgXSJPWlJB+U3i00mbWZHN16uE2QcZG
    b9leMr37xKz45199x9p0TFA8NC5MFmJOsHD60mxxoa35es0R21N6fykAj6YTrbvx
    qUD+rfiJiS6k21Wt8ZreYIUK+8KNJGAXhBp2wGP7zUaxfMZtbuskoPca9pIyjX/C
    MK0iwnVwlXkSqVBu7lizFJ07iuqZaPXbCPzVdiu2b9hNIp64bYAFL324xpBWmhTE
    czuk5435Us8zYG1LGqa5S5CDWf2avx3Rfc3p6/IVSAwlqqLemKiCkZs=
    -----END CERTIFICATE-----
  id: secret-id
  tags:
  - k8s-name:kong-ca
  - k8s-kind:Secret
  - k8s-version:v1
plugins:
- config:
    header_name: kong-id-2
  instance_name: correlation-id-728157fcb
  name: correlation-id
  route: .httpbin.httpbin..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
- config:
    header_name: kong-id-2
  instance_name: correlation-id-c1ebced53
  name: correlation-id
  route: .httpbin-other.httpbin..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
- config:
    header_name: kong-id-2
  instance_name: correlation-id-b8a0ddb44
  name: correlation-id
  route: .httpbin-other.httpbin-other..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
services:
- connect_timeout: 60000
  host: httpbin..80.svc
  name: .httpbin.80
  path: /
  port: 80
  protocol: http
  read_timeout: 60000
  retries: 5
  routes:
  - https_redirect_status_code: 426
    name: .httpbin.httpbin..80
    path_handling: v0
    paths:
    - /httpbin
    preserve_host: true
    protocols:
    - http
    - https
    regex_priority: 0
    request_buffering: true
    response_buffering: true
    strip_path: false
    tags:
    - k8s-name:httpbin
    - k8s-kind:Ingress
    - k8s-group:networking.k8s.io
    - k8s-version:v1
  - https_redirect_status_code: 426
    name: .httpbin-other.httpbin..80
    path_handling: v0
    paths:
    - /httpbin-diff
    preserve_host: true
    protocols:
    - http
    - https
    regex_priority: 0
    request_buffering: true
    response_buffering: true
    strip_path: false
    tags:
    - k8s-name:httpbin-other
    - k8s-kind:Ingress
    - k8s-group:networking.k8s.io
    - k8s-version:v1
  tags:
  - k8s-name:httpbin
  - k8s-kind:Service
  - k8s-version:v1
  write_timeout: 60000
- connect_timeout: 60000
  host: httpbin-other..80.svc
  name: .httpbin-other.80
  path: /
  port: 80
  protocol: http
  read_timeout: 60000
  retries: 5
  routes:
  - https_redirect_status_code: 426
    name: .httpbin-other.httpbin-other..80
    path_handling: v0
    paths:
    - /httpbin-other
    preserve_host: true
    protocols:
    - http
    - https
    regex_priority: 0
    request_buffering: true
    response_buffering: true
    strip_path: false
    tags:
    - k8s-name:httpbin-other
    - k8s-kind:Ingress
    - k8s-group:networking.k8s.io
    - k8s-version:v1
  tags:
  - k8s-name:httpbin-other
  - k8s-kind:Service
  - k8s-version:v1
  write_timeout: 60000
upstreams:
- algorithm: round-robin
  name: httpbin..80.svc
  tags:
  - k8s-name:httpbin
  - k8s-kind:Service
  - k8s-version:v1
- algorithm: round-robin
  name: httpbin-other..80.svc
  tags:
  - k8s-name:httpbin-other
  - k8s-kind:Service
  - k8s-version:v1
`

const twoServicesSampleConfig = `
_format_version: "3.0"
plugins:
- config:
    header_name: kong-id
  instance_name: correlation-id-7f3599b13
  name: correlation-id
  route: .ingress1.httpbin..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
services:
- connect_timeout: 60000
  host: httpbin..80.svc
  name: .httpbin.80
  path: /
  port: 80
  protocol: http
  read_timeout: 60000
  retries: 5
  routes:
  - https_redirect_status_code: 426
    name: .ingress1.httpbin..80
    path_handling: v0
    paths:
    - /httpbin-diff
    preserve_host: true
    protocols:
    - http
    - https
    regex_priority: 0
    request_buffering: true
    response_buffering: true
    strip_path: false
    tags:
    - k8s-name:ingress1
    - k8s-kind:Ingress
    - k8s-group:networking.k8s.io
    - k8s-version:v1
  tags:
  - k8s-name:httpbin
  - k8s-kind:Service
  - k8s-version:v1
  write_timeout: 60000
- connect_timeout: 60000
  host: httpbin-other..80.svc
  name: .httpbin-other.80
  path: /
  port: 80
  protocol: http
  read_timeout: 60000
  retries: 5
  routes:
  - https_redirect_status_code: 426
    name: .ingress2.httpbin-other..80
    path_handling: v0
    paths:
    - /httpbin-other
    preserve_host: true
    protocols:
    - http
    - https
    regex_priority: 0
    request_buffering: true
    response_buffering: true
    strip_path: false
    tags:
    - k8s-name:ingress2
    - k8s-kind:Ingress
    - k8s-group:networking.k8s.io
    - k8s-version:v1
  tags:
  - k8s-name:httpbin-other
  - k8s-kind:Service
  - k8s-version:v1
  write_timeout: 60000
upstreams:
- algorithm: round-robin
  name: httpbin..80.svc
  tags:
  - k8s-name:httpbin
  - k8s-kind:Service
  - k8s-version:v1
- algorithm: round-robin
  name: httpbin-other..80.svc
  tags:
  - k8s-name:httpbin-other
  - k8s-kind:Service
  - k8s-version:v1
`

func TestBuildKongConfigGraph(t *testing.T) {
	testCases := []struct {
		Name                 string
		KongConfig           string
		ExpectedAdjacencyMap map[graph.EntityHash]map[graph.EntityHash]graph2.Edge[graph.EntityHash]
	}{
		{
			Name:       "plugins of the same type with two services",
			KongConfig: sampleKongConfig,
			ExpectedAdjacencyMap: map[graph.EntityHash]map[graph.EntityHash]graph2.Edge[graph.EntityHash]{
				"upstream:httpbin..80.svc":       {},
				"upstream:httpbin-other..80.svc": {},
				"route:.httpbin.httpbin..80": {
					"service:.httpbin.80":             {},
					"plugin:correlation-id-728157fcb": {},
				},
				"route:.httpbin-other.httpbin..80": {
					"plugin:correlation-id-c1ebced53": {},
					"service:.httpbin.80":             {},
				},
				"plugin:correlation-id-c1ebced53": {
					"route:.httpbin-other.httpbin..80": {},
				},
				"plugin:correlation-id-b8a0ddb44": {
					"route:.httpbin-other.httpbin-other..80": {},
				},
				"service:.httpbin-other.80": {
					"route:.httpbin-other.httpbin-other..80": {},
				},
				"route:.httpbin-other.httpbin-other..80": {
					"service:.httpbin-other.80":       {},
					"plugin:correlation-id-b8a0ddb44": {},
				},
				"plugin:correlation-id-728157fcb": {
					"route:.httpbin.httpbin..80": {},
				},
				"ca-certificate:secret-id": {},
				"service:.httpbin.80": {
					"route:.httpbin.httpbin..80":       {},
					"route:.httpbin-other.httpbin..80": {},
				},
			},
		},
		{
			Name:       "two connected components",
			KongConfig: twoServicesSampleConfig,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			g := mustGraphFromRawYAML(t, tc.KongConfig)

			adjacencyMap, err := g.AdjacencyMap()
			require.NoError(t, err)
			t.Logf("adjacency map:\n%s", adjacencyMapString(adjacencyMap))

			// assert.Equal(t, adjacencyMapString(tc.ExpectedAdjacencyMap), adjacencyMapString(adjacencyMap))

			svg, err := graph.RenderGraphDOT(g, "")
			require.NoError(t, err)
			t.Logf("graph: %s", svg)
		})
	}
}

func adjacencyMapString(am map[graph.EntityHash]map[graph.EntityHash]graph2.Edge[graph.EntityHash]) string {
	entries := make([]string, 0, len(am))
	for k, v := range am {
		adjacent := lo.Map(lo.Keys(v), func(h graph.EntityHash, _ int) string { return string(h) })
		// Make the order deterministic.
		sort.Strings(adjacent)
		entries = append(entries, fmt.Sprintf("%s -> [%s]", k, strings.Join(adjacent, ", ")))
	}
	// Make the order deterministic.
	sort.Strings(entries)
	return strings.Join(entries, "\n")
}

func TestBuildingFallbackConfig(t *testing.T) {
	lastKnownGoodConfig := `_format_version: "3.0"
ca_certificates:
- cert: |
    -----BEGIN CERTIFICATE-----
    MIIDtTCCAp2gAwIBAgIBATANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMCVVMx
    EzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xJTAj
    BgNVBAkTHDE1MCBTcGVhciBTdHJlZXQsIFN1aXRlIDE2MDAxDjAMBgNVBBETBTk0
    MTA1MRAwDgYDVQQKEwdLb25nIEhRMB4XDTIzMDkyNjA4NDQzNVoXDTI0MDkyNjA4
    NDQzNVowgYMxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYD
    VQQHEw1TYW4gRnJhbmNpc2NvMSUwIwYDVQQJExwxNTAgU3BlYXIgU3RyZWV0LCBT
    dWl0ZSAxNjAwMQ4wDAYDVQQREwU5NDEwNTEQMA4GA1UEChMHS29uZyBIUTCCASIw
    DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANcnsjbovNdW1HSKaq9ZJ9MwTP0h
    DvShbh9VldLol3Am47xBJYny10EQU4yNF7KhBjbQGAg1hhjDGMp5wPrT66syt4gZ
    IY5xW/6j4GL3E3DNfAgNo+xruEnVHjoz3z6qkt9oAC+T2Gt0BKVtPNQlUqhRBN4f
    YBYoe08K79KSJpjLf96/H8eNJmw5WDzfTH0HdNgZRmcUQfWKgE+iZzAC4ppp8vxx
    YDlXX24GN9bylWcn6TKUkSTolsxJ8mKYDR8zj4Sk6e3z9K14cdIKP3rpXIBiTIrr
    vPsNZAzWgDArTGTS13NC7IzAwkK5iCB582CGJZ8TKqrMHtE+dGwofHJ1Mw0CAwEA
    AaMyMDAwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUniBh32i4njEZO9HxGPmY
    k65EzZswDQYJKoZIhvcNAQELBQADggEBAArcmyihIjIin0nfeL6tJI08Ps2FKIQ9
    7KnKHzkQPQvqgEmmBbzX15w8/YJul7LlspgXSJPWlJB+U3i00mbWZHN16uE2QcZG
    b9leMr37xKz45199x9p0TFA8NC5MFmJOsHD60mxxoa35es0R21N6fykAj6YTrbvx
    qUD+rfiJiS6k21Wt8ZreYIUK+8KNJGAXhBp2wGP7zUaxfMZtbuskoPca9pIyjX/C
    MK0iwnVwlXkSqVBu7lizFJ07iuqZaPXbCPzVdiu2b9hNIp64bYAFL324xpBWmhTE
    czuk5435Us8zYG1LGqa5S5CDWf2avx3Rfc3p6/IVSAwlqqLemKiCkZs=
    -----END CERTIFICATE-----
  id: a5ea1ead-82cd-4b41-8eea-d7e396b8124d
  tags:
  - k8s-name:kong-ca
  - k8s-kind:Secret
  - k8s-version:v1
plugins:
- config:
    header_name: kong-id-2
  instance_name: correlation-id-728157fcb
  name: correlation-id
  route: .httpbin.httpbin..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
- config:
    header_name: kong-id-2
  instance_name: correlation-id-c1ebced53
  name: correlation-id
  route: .httpbin-other.httpbin..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
- config:
    header_name: kong-id-2
  instance_name: correlation-id-b8a0ddb44
  name: correlation-id
  route: .httpbin-other.httpbin-other..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
services:
- connect_timeout: 60000
  host: httpbin.80.svc
  name: .httpbin.80
  path: /
  port: 80
  protocol: http
  read_timeout: 60000
  retries: 5
  routes:
  - https_redirect_status_code: 426
    name: .httpbin.httpbin..80
    path_handling: v0
    paths:
    - /httpbin
    preserve_host: true
    protocols:
    - http
    - https
    regex_priority: 0
    request_buffering: true
    response_buffering: true
    strip_path: false
    tags:
    - k8s-name:httpbin
    - k8s-kind:Ingress
    - k8s-group:networking.k8s.io
    - k8s-version:v1
  - https_redirect_status_code: 426
    name: .httpbin-other.httpbin..80
    path_handling: v0
    paths:
    - /httpbin-diff
    preserve_host: true
    protocols:
    - http
    - https
    regex_priority: 0
    request_buffering: true
    response_buffering: true
    strip_path: false
    tags:
    - k8s-name:httpbin-other
    - k8s-kind:Ingress
    - k8s-group:networking.k8s.io
    - k8s-version:v1
  tags:
  - k8s-name:httpbin
  - k8s-kind:Service
  - k8s-version:v1
  write_timeout: 60000
- connect_timeout: 60000
  host: httpbin-other.80.svc
  name: .httpbin-other.80
  path: /
  port: 80
  protocol: http
  read_timeout: 60000
  retries: 5
  routes:
  - https_redirect_status_code: 426
    name: .httpbin-other.httpbin-other..80
    path_handling: v0
    paths:
    - /httpbin-other
    preserve_host: true
    protocols:
    - http
    - https
    regex_priority: 0
    request_buffering: true
    response_buffering: true
    strip_path: false
    tags:
    - k8s-name:httpbin-other
    - k8s-kind:Ingress
    - k8s-group:networking.k8s.io
    - k8s-version:v1
  tags:
  - k8s-name:httpbin-other
  - k8s-kind:Service
  - k8s-version:v1
  write_timeout: 60000
upstreams:
- algorithm: round-robin
  name: httpbin.80.svc
  tags:
  - k8s-name:httpbin
  - k8s-kind:Service
  - k8s-version:v1
- algorithm: round-robin
  name: httpbin-other.80.svc
  tags:
  - k8s-name:httpbin-other
  - k8s-kind:Service
  - k8s-version:v1`
	lastKnownGoodConfigGraph := mustGraphFromRawYAML(t, lastKnownGoodConfig)

	// This is the current Kong config parser has generated.
	currentConfig := `_format_version: "3.0"
plugins:
- config:
    header_name: kong-id-2
  instance_name: correlation-id-728157fcb
  name: correlation-id
  route: .httpbin.httpbin..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
- config:
    header_name: kong-id-2
  name: correlation-id
  route: .httpbin-other.httpbin-other..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
- config:
    header_name: kong-id-2
  name: correlation-id
  route: .httpbin-other.httpbin-other..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
  name: key-auth
  route: .httpbin-other.httpbin-other..80
  tags:
  - k8s-name:kong-id
  - k8s-kind:KongPlugin
  - k8s-group:configuration.konghq.com
  - k8s-version:v1
services:
- connect_timeout: 60000
  host: httpbin..80.svc
  name: .httpbin.80
  path: /
  port: 80
  protocol: http
  read_timeout: 60000
  retries: 5
  routes:
  - https_redirect_status_code: 426
    name: .httpbin.httpbin..80
    path_handling: v0
    paths:
    - /httpbin
    preserve_host: true
    protocols:
    - http
    - https
    regex_priority: 0
    request_buffering: true
    response_buffering: true
    strip_path: false
    tags:
    - k8s-name:httpbin
    - k8s-kind:Ingress
    - k8s-group:networking.k8s.io
    - k8s-version:v1
  tags:
  - k8s-name:httpbin
  - k8s-kind:Service
  - k8s-version:v1
  write_timeout: 60000
- connect_timeout: 60000
  host: httpbin-other..80.svc
  name: .httpbin-other.80
  path: /
  port: 80
  protocol: http
  read_timeout: 60000
  retries: 5
  routes:
  - https_redirect_status_code: 426
    name: .httpbin-other.httpbin-other..80
    path_handling: v0
    paths:
    - /httpbin-other
    preserve_host: true
    protocols:
    - http
    - https
    regex_priority: 0
    request_buffering: true
    response_buffering: true
    strip_path: false
    tags:
    - k8s-name:httpbin-other
    - k8s-kind:Ingress
    - k8s-group:networking.k8s.io
    - k8s-version:v1
  tags:
  - k8s-name:httpbin-other
  - k8s-kind:Service
  - k8s-version:v1
  write_timeout: 60000
upstreams:
- algorithm: round-robin
  name: httpbin..80.svc
  tags:
  - k8s-name:httpbin
  - k8s-kind:Service
  - k8s-version:v1
- algorithm: round-robin
  name: httpbin-other..80.svc
  tags:
  - k8s-name:httpbin-other
  - k8s-kind:Service
  - k8s-version:v1`
	currentConfigGraph := mustGraphFromRawYAML(t, currentConfig)

	entitiesErrors := []sendconfig.FlatEntityError{
		{
			Type: "route",
			Name: ".httpbin-other.httpbin-other..80",
		},
	}

	fallbackConfig, err := graph.BuildFallbackKongConfig(lastKnownGoodConfigGraph, currentConfigGraph, entitiesErrors)
	require.NoError(t, err)

	lastGoodSvg := dumpGraphAsDOT(t, lastKnownGoodConfigGraph)
	currentSvg := dumpGraphAsDOT(t, currentConfigGraph)
	fallbackSvg := dumpGraphAsDOT(t, fallbackConfig)
	t.Logf("open %s %s %s", lastGoodSvg, currentSvg, fallbackSvg)

	fallbackKongConfig, err := graph.BuildKongConfigFromGraph(fallbackConfig)
	require.NoError(t, err)

	b, err := yaml.Marshal(fallbackKongConfig)
	require.NoError(t, err)
	t.Logf("fallback config:\n%s", string(b))
}

func mustGraphFromRawYAML(t *testing.T, y string) graph.KongConfigGraph {
	t.Helper()
	kongConfig := &file.Content{}
	err := yaml.Unmarshal([]byte(y), kongConfig)
	require.NoError(t, err)

	g, err := graph.BuildKongConfigGraph(kongConfig)
	require.NoError(t, err)
	return g
}

func dumpGraphAsDOT(t *testing.T, g graph.KongConfigGraph) string {
	svg, err := graph.RenderGraphDOT(g, "")
	require.NoError(t, err)
	return svg
}
