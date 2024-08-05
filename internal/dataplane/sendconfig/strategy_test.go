package sendconfig_test

import (
	"fmt"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
)

type clientMock struct {
	isKonnect bool

	konnectControlPlaneWasCalled bool
	adminAPIClientWasCalled      bool
}

func (c *clientMock) IsKonnect() bool {
	return c.isKonnect
}

func (c *clientMock) KonnectControlPlane() string {
	c.konnectControlPlaneWasCalled = true
	return uuid.NewString()
}

func (c *clientMock) AdminAPIClient() *kong.Client {
	c.adminAPIClientWasCalled = true
	return &kong.Client{}
}

type clientWithBackoffMock struct {
	*clientMock
}

func (c clientWithBackoffMock) BackoffStrategy() adminapi.UpdateBackoffStrategy {
	return newMockBackoffStrategy(true)
}

func TestDefaultUpdateStrategyResolver_ResolveUpdateStrategy(t *testing.T) {
	testCases := []struct {
		isKonnect                     bool
		inMemory                      bool
		expectedStrategyType          string
		expectKonnectControlPlaneCall bool
	}{
		{
			isKonnect:                     true,
			inMemory:                      false,
			expectedStrategyType:          "WithBackoff(DBMode)",
			expectKonnectControlPlaneCall: true,
		},
		{
			isKonnect:                     true,
			inMemory:                      true,
			expectedStrategyType:          "WithBackoff(DBMode)",
			expectKonnectControlPlaneCall: true,
		},
		{
			isKonnect:            false,
			inMemory:             false,
			expectedStrategyType: "DBMode",
		},
		{
			isKonnect:            false,
			inMemory:             true,
			expectedStrategyType: "InMemory",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("isKonnect=%v inMemory=%v", tc.isKonnect, tc.inMemory), func(t *testing.T) {
			client := &clientMock{
				isKonnect: tc.isKonnect,
			}

			var updateClient sendconfig.UpdateClient
			if tc.isKonnect {
				updateClient = &clientWithBackoffMock{client}
			} else {
				updateClient = client
			}

			resolver := sendconfig.NewDefaultUpdateStrategyResolver(sendconfig.Config{
				InMemory: tc.inMemory,
			}, zapr.NewLogger(zap.NewNop()))

			strategy := resolver.ResolveUpdateStrategy(updateClient, nil)
			require.Equal(t, tc.expectedStrategyType, strategy.Type())
			assert.True(t, client.adminAPIClientWasCalled)
			assert.Equal(t, tc.expectKonnectControlPlaneCall, client.konnectControlPlaneWasCalled)
		})
	}
}
