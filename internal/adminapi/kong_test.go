package adminapi_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func TestMakeHTTPClientWithTLSOpts(t *testing.T) {
	cert, key := certificate.MustGenerateSelfSignedCertPEMFormat()
	caCert := cert

	opts := adminapi.HTTPClientOpts{
		TLSSkipVerify: true,
		TLSServerName: "",
		CACertPath:    "",
		CACert:        string(caCert),
		Headers:       nil,
		TLSClient: adminapi.TLSClientConfig{
			Cert: string(cert),
			Key:  string(key),
		},
	}

	t.Run("without kong admin token", func(t *testing.T) {
		c, err := adminapi.MakeHTTPClient(&opts, "")
		require.NoError(t, err)
		require.NotNil(t, c)
		validate(t, c, caCert, cert, key, "")
	})

	t.Run("with kong admin token", func(t *testing.T) {
		const kongAdminToken = "my-token"
		c, err := adminapi.MakeHTTPClient(&opts, kongAdminToken)
		require.NoError(t, err)
		require.NotNil(t, c)
		validate(t, c, caCert, cert, key, kongAdminToken)
	})
}

func TestMakeHTTPClientWithTLSOptsAndFilePaths(t *testing.T) {
	cert, key := certificate.MustGenerateSelfSignedCertPEMFormat()
	caCert := cert

	certDir := t.TempDir()

	caFile, err := os.CreateTemp(certDir, "ca.crt")
	require.NoError(t, err)
	writtenBytes, err := caFile.Write(caCert)
	require.NoError(t, err)
	require.Len(t, caCert, writtenBytes)

	certFile, err := os.CreateTemp(certDir, "cert.crt")
	require.NoError(t, err)
	writtenBytes, err = certFile.Write(cert)
	require.NoError(t, err)
	require.Len(t, cert, writtenBytes)

	certPrivateKeyFile, err := os.CreateTemp(certDir, "cert.key")
	require.NoError(t, err)
	writtenBytes, err = certPrivateKeyFile.Write(key)
	require.NoError(t, err)
	require.Len(t, key, writtenBytes)

	opts := adminapi.HTTPClientOpts{
		TLSSkipVerify: true,
		TLSServerName: "",
		CACertPath:    caFile.Name(),
		CACert:        "",
		Headers:       nil,
		TLSClient: adminapi.TLSClientConfig{
			CertFile: certFile.Name(),
			KeyFile:  certPrivateKeyFile.Name(),
		},
	}

	t.Run("without kong admin token", func(t *testing.T) {
		c, err := adminapi.MakeHTTPClient(&opts, "")
		require.NoError(t, err)
		require.NotNil(t, c)
		validate(t, c, caCert, cert, key, "")
	})

	t.Run("with kong admin token", func(t *testing.T) {
		const kongAdminToken = "my-token"
		c, err := adminapi.MakeHTTPClient(&opts, kongAdminToken)
		require.NoError(t, err)
		require.NotNil(t, c)
		validate(t, c, caCert, cert, key, kongAdminToken)
	})
}

func TestNewKongClientForWorkspace(t *testing.T) {
	const workspace = "workspace"

	testCases := []struct {
		name            string
		adminAPIReady   bool
		adminAPIVersion string
		workspace       string
		workspaceExists bool
		expectError     error
	}{
		{
			name:          "admin api is ready for default workspace",
			adminAPIReady: true,
			workspace:     "",
		},
		{
			name:            "admin api is ready and workspace exists",
			adminAPIReady:   true,
			workspace:       workspace,
			workspaceExists: true,
		},
		{
			name:            "admin api is ready and workspace doesn't exist",
			adminAPIReady:   true,
			workspace:       workspace,
			workspaceExists: false,
		},
		{
			name:          "admin api is not ready",
			adminAPIReady: false,
			workspace:     "",
			expectError:   adminapi.KongClientNotReadyError{},
		},
		{
			name:            "admin api is not ready for an existing workspace",
			adminAPIReady:   false,
			workspace:       workspace,
			workspaceExists: true,
			expectError:     adminapi.KongClientNotReadyError{},
		},
		{
			name:            "admin api is in too old version",
			adminAPIReady:   true,
			adminAPIVersion: "3.4.0",
			workspace:       "",
			expectError:     adminapi.KongGatewayUnsupportedVersionError{},
		},
		{
			name:            "admin api is in supported OSS version",
			adminAPIReady:   true,
			workspace:       "",
			adminAPIVersion: versions.KICv3VersionCutoff.String(),
		},
		{
			name:            "admin api has malformed version",
			adminAPIReady:   true,
			adminAPIVersion: "3-malformed-version",
			workspace:       "",
			expectError:     adminapi.KongGatewayUnsupportedVersionError{},
		},
		{
			name:            "admin api has malformed version for workspace",
			adminAPIReady:   true,
			adminAPIVersion: "3-malformed-version",
			workspace:       workspace,
			workspaceExists: true,
			expectError:     adminapi.KongGatewayUnsupportedVersionError{},
		},
		{
			name:            "admin api has enterprise version",
			adminAPIReady:   true,
			workspace:       "",
			adminAPIVersion: "3.4.1.2",
		},
		{
			name:            "admin api has enterprise version for workspace",
			adminAPIReady:   true,
			workspace:       workspace,
			workspaceExists: true,
			adminAPIVersion: "3.4.1.2",
		},
		{
			name:            "admin api has too old enterprise version",
			adminAPIReady:   true,
			adminAPIVersion: "3.4.0.2",
			workspace:       "",
			expectError:     adminapi.KongGatewayUnsupportedVersionError{},
		},
		{
			name:            "admin api has too old enterprise version for workspace",
			adminAPIReady:   true,
			adminAPIVersion: "3.4.0.2",
			workspace:       workspace,
			workspaceExists: true,
			expectError:     adminapi.KongGatewayUnsupportedVersionError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adminAPIHandler := mocks.NewAdminAPIHandler(
				t,
				mocks.WithWorkspaceExists(tc.workspaceExists),
				mocks.WithReady(tc.adminAPIReady),
				mocks.WithVersion(tc.adminAPIVersion),
			)
			adminAPIServer := httptest.NewServer(adminAPIHandler)
			t.Cleanup(func() { adminAPIServer.Close() })

			client, err := adminapi.NewKongClientForWorkspace(
				context.Background(),
				adminAPIServer.URL,
				tc.workspace,
				adminAPIServer.Client(),
			)

			if tc.expectError != nil {
				require.IsType(t, err, tc.expectError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, client)

			if tc.workspace != "" && !tc.workspaceExists {
				require.True(t, adminAPIHandler.WasWorkspaceCreated(), "expected workspace to be created")
			}

			require.Equal(t, client.AdminAPIClient().Workspace(), tc.workspace)
			_, ok := client.PodReference()
			require.False(t, ok, "expected no pod reference to be attached to the client")
		})
	}
}

// validate spins up a test server with the given TLS configuration and verifies
// whether the passed client can connect to it successfully.
func validate(
	t *testing.T,
	httpClient *http.Client,
	caPEM []byte,
	certPEM []byte,
	certPrivateKeyPEM []byte,
	kongAdminToken string,
) {
	serverCert, err := tls.X509KeyPair(certPEM, certPrivateKeyPEM)
	require.NoError(t, err, "fail to load server certificates")

	certPool := x509.NewCertPool()
	require.True(t, certPool.AppendCertsFromPEM(caPEM))

	serverTLSConf := &tls.Config{
		RootCAs:      certPool,
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAnyClientCert,
		Certificates: []tls.Certificate{serverCert},
		MinVersion:   tls.VersionTLS12,
	}

	const successMessage = "connection successful"
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if kongAdminToken != "" {
			v, ok := r.Header[http.CanonicalHeaderKey(adminapi.HeaderNameAdminToken)]
			if !ok {
				fmt.Fprintf(w, "%s header not found", adminapi.HeaderNameAdminToken)
				return
			}
			if len(v) != 1 {
				fmt.Fprintf(w, "%s header expected to contain %s but found %v",
					adminapi.HeaderNameAdminToken, kongAdminToken, v)
				return
			}
			if v[0] != kongAdminToken {
				fmt.Fprintf(w, "%s header expected to contain %s but found %s",
					adminapi.HeaderNameAdminToken, kongAdminToken, v[0])
				return
			}
		}
		fmt.Fprintln(w, successMessage)
	}))
	server.TLS = serverTLSConf
	server.StartTLS()
	defer server.Close()

	response, err := httpClient.Get(server.URL)
	require.NoError(t, err, "HTTP client failed to issue a GET request")
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	require.NoError(t, err, "failed to read response body")
	require.Equal(t, strings.TrimSpace(string(data)), successMessage, "unexpected content of response body")
}
