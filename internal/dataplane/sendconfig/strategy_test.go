package sendconfig_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
)

type clientMock struct {
	isKonnect bool

	konnectRuntimeGroupWasCalled bool
	adminAPIClientWasCalled      bool
}

func (c *clientMock) IsKonnect() bool {
	return c.isKonnect
}

func (c *clientMock) KonnectRuntimeGroup() string {
	c.konnectRuntimeGroupWasCalled = true
	return uuid.NewString()
}

func (c *clientMock) AdminAPIClient() *kong.Client {
	c.adminAPIClientWasCalled = true
	return &kong.Client{}
}

func TestResolveUpdateStrategy(t *testing.T) {
	testCases := []struct {
		isKonnect                     bool
		inMemory                      bool
		expectedStrategy              sendconfig.UpdateStrategy
		expectKonnectRuntimeGroupCall bool
	}{
		{
			isKonnect:                     true,
			inMemory:                      false,
			expectedStrategy:              sendconfig.UpdateStrategyDBMode{},
			expectKonnectRuntimeGroupCall: true,
		},
		{
			isKonnect:                     true,
			inMemory:                      true,
			expectedStrategy:              sendconfig.UpdateStrategyDBMode{},
			expectKonnectRuntimeGroupCall: true,
		},
		{
			isKonnect:        false,
			inMemory:         false,
			expectedStrategy: sendconfig.UpdateStrategyDBMode{},
		},
		{
			isKonnect:        false,
			inMemory:         true,
			expectedStrategy: sendconfig.UpdateStrategyInMemory{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("isKonnect=%v inMemory=%v", tc.isKonnect, tc.inMemory), func(t *testing.T) {
			client := &clientMock{
				isKonnect: tc.isKonnect,
			}

			strategy := sendconfig.ResolveUpdateStrategy(client, sendconfig.Config{
				InMemory: tc.inMemory,
			})
			require.IsType(t, tc.expectedStrategy, strategy)
			assert.True(t, client.adminAPIClientWasCalled)
			assert.Equal(t, tc.expectKonnectRuntimeGroupCall, client.konnectRuntimeGroupWasCalled)
		})
	}
}
