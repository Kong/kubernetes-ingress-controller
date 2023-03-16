package atc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func marshalKongConfig(t *testing.T, s kong.Service, r kong.Route) io.Reader {
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
	assert.NoError(t, err)

	return bytes.NewReader(config)
}

func TestApplyExpressionDBLess(t *testing.T) {
	// TODO: run Kong locally
	s := &kong.Service{
		// TODO: run a local service and use it here
		Host: kong.String("www.example.com"),
		Path: kong.String("/"),
	}

	kongAdminURL := "http://127.0.0.1:8001"
	kongClient, err := kong.NewClient(&kongAdminURL, nil)
	require.NoError(t, err)

	// kongProxyURL := "http://127.0.0.1:8000"

	testCases := []struct {
		name            string
		matcher         Matcher
		matchRequests   []*http.Request
		unmatchRequests []*http.Request
	}{
		{
			name: "prefix match of path",
			matcher: Or(
				NewPredicateHTTPPath(OpEqual, "/foo"),
				NewPredicateHTTPPath(OpPrefixMatch, "/foo/"),
			),

			// TODO: add test requests
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := &kong.Route{
				StripPath: kong.Bool(true),
			}
			ApplyExpression(r, tc.matcher, 1)
			req, err := kongClient.NewRequest("POST", "/config", nil, marshalKongConfig(t, *s, *r))
			require.NoError(t, err)

			resp, err := kongClient.DoRAW(context.Background(), req)
			require.NoError(t, err)
			// we do not need resp body here, so close it immediately.
			resp.Body.Close()

			require.Equal(t, http.StatusCreated, resp.StatusCode)

			// matched requests should access upstream service
			for _, req := range tc.matchRequests {
				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				resp.Body.Close()
				require.Equal(t, http.StatusOK, resp.StatusCode)
			}
			// unmatched requests should get a 404 from Kong
			for _, req := range tc.unmatchRequests {
				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				resp.Body.Close()
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			}
		})

	}

}
