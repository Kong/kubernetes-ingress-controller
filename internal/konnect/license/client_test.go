package license_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
)

type mockKonnectLicenseServer struct {
	response   []byte
	statusCode int
	t          *testing.T
}

func newMockKonnectLicenseServer(t *testing.T, response []byte, statusCode int) *mockKonnectLicenseServer {
	return &mockKonnectLicenseServer{
		t:          t,
		response:   response,
		statusCode: statusCode,
	}
}

func (m *mockKonnectLicenseServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	require.Equal(m.t, metadata.UserAgent(), r.Header.Get("User-Agent"))
	w.WriteHeader(m.statusCode)
	_, _ = w.Write(m.response)
}

func TestLicenseClient(t *testing.T) {
	testCases := []struct {
		name       string
		response   []byte
		status     int
		assertions func(t *testing.T, c *license.Client)
	}{
		{
			name: "200 valid response",
			response: []byte(`{
				"items": [
					{
						"payload": "some-license-content",
						"updated_at": 1234567890,
						"id": "some-license-id"
					}
				]
			}`),
			status: http.StatusOK,
			assertions: func(t *testing.T, c *license.Client) {
				licenseOpt, err := c.Get(context.Background())
				require.NoError(t, err)

				l, ok := licenseOpt.Get()
				require.True(t, ok)
				require.Equal(t, "some-license-content", l.Payload)
				require.Equal(t, int64(1234567890), l.UpdatedAt.Unix())
			},
		},
		{
			name:     "200 but empty response",
			response: []byte(`{}`),
			status:   http.StatusOK,
			assertions: func(t *testing.T, c *license.Client) {
				_, err := c.Get(context.Background())
				require.ErrorContains(t, err, "no license item found in response")
			},
		},
		{
			name:     "200 but invalid response",
			response: []byte(`{invalid-json`),
			status:   http.StatusOK,
			assertions: func(t *testing.T, c *license.Client) {
				_, err := c.Get(context.Background())
				require.ErrorContains(t, err, "failed to parse response body")
			},
		},
		{
			name: "200 but empty license id",
			response: []byte(`{
				"items": [
					{
						"payload": "some-license-content",
						"updated_at": 1234567890,
						"id": ""
					}
				]
			}`),
			status: http.StatusOK,
			assertions: func(t *testing.T, c *license.Client) {
				_, err := c.Get(context.Background())
				require.ErrorContains(t, err, "empty id")
			},
		},
		{
			name: "200 but empty updated_at",
			response: []byte(`{
				"items": [
					{
						"payload": "some-license-content",
						"updated_at": 0,
						"id": "some-license-id"
					}
				]
			}`),
			status: http.StatusOK,
			assertions: func(t *testing.T, c *license.Client) {
				_, err := c.Get(context.Background())
				require.ErrorContains(t, err, "empty updated_at")
			},
		},
		{
			name: "200 but empty payload",
			response: []byte(`{
				"items": [
					{
						"payload": "",
						"updated_at": 1234567890,
						"id": "some-license-id"
					}
				]
			}`),
			status: http.StatusOK,
			assertions: func(t *testing.T, c *license.Client) {
				_, err := c.Get(context.Background())
				require.ErrorContains(t, err, "empty license")
			},
		},
		{
			name:     "404 returns empty license with no error",
			response: nil,
			status:   http.StatusNotFound,
			assertions: func(t *testing.T, c *license.Client) {
				l, err := c.Get(context.Background())
				require.NoError(t, err)
				require.False(t, l.IsPresent())
			},
		},
		{
			name:     "400 returns error",
			response: nil,
			status:   http.StatusBadRequest,
			assertions: func(t *testing.T, c *license.Client) {
				_, err := c.Get(context.Background())
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newMockKonnectLicenseServer(t, tc.response, tc.status)
			ts := httptest.NewServer(server)
			defer ts.Close()

			c, err := license.NewClient(adminapi.KonnectConfig{Address: ts.URL})
			require.NoError(t, err)
			tc.assertions(t, c)
		})
	}
}
