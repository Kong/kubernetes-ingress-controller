package configfetcher

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func TestTryFetchingValidConfigFromGateways(t *testing.T) {
	const (
		zeroConfigHash = "00000000000000000000000000000000"
		configHash     = "8f1dd2f83bc2627cc6b71c76d1476592"
	)

	startAdminAPI := func(t *testing.T, opts ...mocks.AdminAPIHandlerOpt) *adminapi.Client {
		adminAPIHandler := mocks.NewAdminAPIHandler(t, opts...)
		adminAPIServer := httptest.NewServer(adminAPIHandler)
		t.Cleanup(func() { adminAPIServer.Close() })

		// NOTE: We use here adminapi.NewKongAPIClient() as that doesn't check
		// the status of the Kong Gateway but just returns the client.
		client, err := adminapi.NewKongAPIClient(
			adminAPIServer.URL,
			adminAPIServer.Client(),
		)
		require.NoError(t, err)
		require.NotNil(t, client)
		return adminapi.NewClient(client)
	}

	testCases := []struct {
		name                    string
		expectError             bool
		expectedLastValidStatus bool
		adminAPIClients         func(t *testing.T) []*adminapi.Client
	}{
		{
			name: "correct configuration hash",
			adminAPIClients: func(t *testing.T) []*adminapi.Client {
				return []*adminapi.Client{
					startAdminAPI(t, mocks.WithReady(true), mocks.WithConfigurationHash(configHash)),
					startAdminAPI(t, mocks.WithReady(true), mocks.WithConfigurationHash(configHash)),
				}
			},
			expectedLastValidStatus: true,
		},
		{
			name: "zero configuration hash",
			adminAPIClients: func(t *testing.T) []*adminapi.Client {
				return []*adminapi.Client{
					startAdminAPI(t, mocks.WithReady(true), mocks.WithConfigurationHash(zeroConfigHash)),
					startAdminAPI(t, mocks.WithReady(true), mocks.WithConfigurationHash(zeroConfigHash)),
				}
			},
		},
		{
			name: "none are ready",
			adminAPIClients: func(t *testing.T) []*adminapi.Client {
				return []*adminapi.Client{
					startAdminAPI(t, mocks.WithReady(false)),
					startAdminAPI(t, mocks.WithReady(false)),
				}
			},
			expectError: true,
		},
		{
			name: "one out of 2 is ready",
			adminAPIClients: func(t *testing.T) []*adminapi.Client {
				return []*adminapi.Client{
					startAdminAPI(t, mocks.WithReady(true), mocks.WithConfigurationHash(configHash)),
					startAdminAPI(t, mocks.WithReady(false)),
				}
			},
			expectedLastValidStatus: true,
		},
		{
			name: "one out of 2 is ready with zero config hash",
			adminAPIClients: func(t *testing.T) []*adminapi.Client {
				return []*adminapi.Client{
					startAdminAPI(t, mocks.WithReady(true), mocks.WithConfigurationHash(zeroConfigHash)),
					startAdminAPI(t, mocks.WithReady(false)),
				}
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fetcher := NewDefaultKongLastGoodConfigFetcher(false, "")
			state, ok := fetcher.LastValidConfig()
			require.False(t, ok)
			require.Nil(t, state)

			ctx := context.Background()
			clients := tc.adminAPIClients(t)
			logger := zapr.NewLogger(zap.NewNop())
			err := fetcher.TryFetchingValidConfigFromGateways(ctx, logger, clients, nil)
			if tc.expectError {
				require.Error(t, err)
				assert.False(t, ok)
				assert.Nil(t, state)
				return
			}

			require.NoError(t, err)

			state, ok = fetcher.LastValidConfig()
			if tc.expectedLastValidStatus {
				assert.True(t, ok)
				assert.NotNil(t, state)
			} else {
				assert.False(t, ok)
				assert.Nil(t, state)
			}
		})
	}
}
