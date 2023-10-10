package kongintegration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v2/test/kongintegration/containers"
)

func TestExpressionsRouterMatchers_GenerateValidExpressions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	kongC := containers.NewKong(ctx, t)
	kongClient, err := kong.NewClient(lo.ToPtr(kongC.AdminURL(ctx, t)), &http.Client{})
	require.NoError(t, err)

	httpBinC := containers.NewHTTPBin(ctx, t)

	testCases := []struct {
		name            string
		matcher         atc.Matcher
		matchRequests   []*http.Request
		unmatchRequests []*http.Request
	}{
		{
			name:    "exact match on path",
			matcher: atc.NewPredicateHTTPPath(atc.OpEqual, "/foo"),
			matchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com").String(), "foo", nil),
			},
			unmatchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com").String(), "foobar", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com").String(), "foo/", nil),
			},
		},
		{
			name: "exact match on path and host",
			matcher: atc.And(
				atc.NewPredicateHTTPPath(atc.OpEqual, "/foo"),
				atc.NewPrediacteHTTPHost(atc.OpEqual, "a.foo.com"),
			),
			matchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com").String(), "foo", nil),
			},
			unmatchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com").String(), "foobar", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://b.foo.com").String(), "foo", nil),
			},
		},
		{
			name: "exact match on path and wildcard match on host",
			matcher: atc.And(
				atc.NewPredicateHTTPPath(atc.OpEqual, "/foo"),
				atc.NewPrediacteHTTPHost(atc.OpSuffixMatch, ".foo.com"),
			),
			matchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com").String(), "foo", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://b.foo.com").String(), "foo", nil),
			},
			unmatchRequests: []*http.Request{
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.foo.com").String(), "foobar", nil),
				helpers.MustHTTPRequest(t, "GET", helpers.MustParseURL(t, "http://a.bar.com").String(), "foo", nil),
			},
		},
	}

	proxyParsedURL, err := url.Parse(kongC.ProxyURL(ctx, t))
	require.NoError(t, err)
	proxyClient := helpers.DefaultHTTPClient()
	proxyClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyParsedURL),
	}

	s := &kong.Service{
		Host: kong.String(httpBinC.IP(ctx, t)),
		Path: kong.String("/"),
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := &kong.Route{
				StripPath: kong.Bool(true),
			}
			atc.ApplyExpression(r, tc.matcher, 1)
			req, err := kongClient.NewRequest("POST", "/config", nil, marshalKongConfig(t, *s, *r))
			require.NoError(t, err)

			resp, err := kongClient.DoRAW(ctx, req)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())
			require.Equal(t, http.StatusCreated, resp.StatusCode)

			// matched requests should access upstream service
			require.Eventually(t, func() bool {
				for _, req := range tc.matchRequests {
					resp, err := proxyClient.Do(req)
					if err != nil {
						t.Logf("error happened on getting response from kong: %v", err)
						return false
					}
					resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						return false
					}
				}
				return true
			}, time.Minute, 5*time.Second)

			// unmatched requests should get a 404 from Kong
			for _, req := range tc.unmatchRequests {
				resp, err := proxyClient.Do(req)
				require.NoError(t, err)
				require.NoError(t, resp.Body.Close())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			}
		})
	}
}

func marshalKongConfig(t *testing.T, s kong.Service, r kong.Route) io.Reader {
	t.Helper()
	content := &file.Content{
		FormatVersion: "3.0",
		Services: []file.FService{
			{
				Service: s,
				Routes: []*file.FRoute{
					{
						Route: r,
					},
				},
			},
		},
	}
	config, err := json.Marshal(content)
	require.NoError(t, err)

	return bytes.NewReader(config)
}
