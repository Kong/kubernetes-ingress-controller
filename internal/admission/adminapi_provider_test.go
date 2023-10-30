package admission_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission"
)

type fakeGatewayClientsProvider struct {
	clients []*adminapi.Client
}

func (f fakeGatewayClientsProvider) GatewayClients() []*adminapi.Client {
	return f.clients
}

func TestDefaultAdminAPIServicesProvider(t *testing.T) {
	t.Run("no clients available should return false from methods", func(t *testing.T) {
		p := admission.NewDefaultAdminAPIServicesProvider(fakeGatewayClientsProvider{})

		_, ok := p.GetConsumersService()
		require.False(t, ok)

		_, ok = p.GetPluginsService()
		require.False(t, ok)

		_, ok = p.GetConsumerGroupsService()
		require.False(t, ok)

		_, ok = p.GetInfoService()
		require.False(t, ok)
	})

	t.Run("when clients available should return first one", func(t *testing.T) {
		firstClient := lo.Must(adminapi.NewTestClient("localhost:8080"))
		p := admission.NewDefaultAdminAPIServicesProvider(fakeGatewayClientsProvider{
			clients: []*adminapi.Client{
				firstClient,
				lo.Must(adminapi.NewTestClient("localhost:8081")),
			},
		})

		consumersSvc, ok := p.GetConsumersService()
		require.True(t, ok)
		require.Equal(t, firstClient.AdminAPIClient().Consumers, consumersSvc)

		pluginsSvc, ok := p.GetPluginsService()
		require.True(t, ok)
		require.Equal(t, firstClient.AdminAPIClient().Plugins, pluginsSvc)

		consumerGroupsSvc, ok := p.GetConsumerGroupsService()
		require.True(t, ok)
		require.Equal(t, firstClient.AdminAPIClient().ConsumerGroups, consumerGroupsSvc)
	})
}
