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

	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/atc"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/kongintegration/containers"
)

func TestExpressionsRouterMatchers_GenerateValidExpressions(t *testing.T) {
	t.Parallel()

	const (
		timeout = time.Second * 5
		period  = time.Millisecond * 200
	)

	ctx := context.Background()

	kongC := containers.NewKong(ctx, t)
	kongClient, err := adminapi.NewKongAPIClient(kongC.AdminURL(ctx, t), helpers.DefaultHTTPClient())
	require.NoError(t, err)

	httpBinC := containers.NewHTTPBin(ctx, t)

	type request struct {
		host string
		path string
	}
	testCases := []struct {
		name            string
		matcher         atc.Matcher
		matchRequests   []request
		unmatchRequests []request
	}{
		{
			name:          "exact match on path",
			matcher:       atc.NewPredicateHTTPPath(atc.OpEqual, "/foo"),
			matchRequests: []request{{"a.foo.com", "foo"}},
			unmatchRequests: []request{
				{"a.foo.com", "foobar"},
				{"a.foo.com", "foo/"},
			},
		},
		{
			name: "exact match on path and host",
			matcher: atc.And(
				atc.NewPredicateHTTPPath(atc.OpEqual, "/foo"),
				atc.NewPrediacteHTTPHost(atc.OpEqual, "a.foo.com"),
			),
			matchRequests: []request{{"a.foo.com", "foo"}},
			unmatchRequests: []request{
				{"a.foo.com", "foobar"},
				{"b.foo.com", "foo"},
			},
		},
		{
			name: "exact match on path and wildcard match on host",
			matcher: atc.And(
				atc.NewPredicateHTTPPath(atc.OpEqual, "/foo"),
				atc.NewPrediacteHTTPHost(atc.OpSuffixMatch, ".foo.com"),
			),
			matchRequests: []request{
				{"a.foo.com", "foo"},
				{"b.foo.com", "foo"},
			},
			unmatchRequests: []request{
				{"a.foo.com", "foobar"},
				{"a.bar.com", "foo"},
			},
		},
	}

	proxyParsedURL, err := url.Parse(kongC.ProxyURL(ctx, t))
	require.NoError(t, err)
	s := &kong.Service{
		Host: kong.String(httpBinC.IP(ctx, t)),
		Path: kong.String("/"),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &kong.Route{
				StripPath: kong.Bool(true),
			}
			atc.ApplyExpression(r, tc.matcher, 1)
			req, err := kongClient.NewRequest(http.MethodPost, "/config", nil, marshalKongConfig(t, *s, *r))
			require.NoError(t, err)

			resp, err := kongClient.DoRAW(ctx, req)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())
			require.Equal(t, http.StatusCreated, resp.StatusCode)

			// Matched requests should access the upstream service.
			for _, req := range tc.matchRequests {
				helpers.EventuallyGETPath(t, proxyParsedURL, req.host, req.path, nil, http.StatusOK, "", nil, timeout, period)
			}

			// Unmatched requests should get a 404 from Kong.
			for _, req := range tc.unmatchRequests {
				helpers.EventuallyGETPath(t, proxyParsedURL, req.host, req.path, nil, http.StatusNotFound, "", nil, timeout, period)
			}
		})
	}
}

func marshalKongConfig(t *testing.T, s kong.Service, r kong.Route) io.Reader {
	t.Helper()
	content := &file.Content{
		FormatVersion: versions.DeckFileFormatVersion,
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
