package admission

import (
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

// GatewayClientsProvider returns the most recent set of Gateway Admin API clients.
type GatewayClientsProvider interface {
	GatewayClients() []*adminapi.Client
}

// DefaultAdminAPIServicesProvider allows getting Admin API services that require having at least one Gateway discovered.
// In the case there's no Gateways, it will return `false` from every method, signalling there's no Gateway available.
type DefaultAdminAPIServicesProvider struct {
	gatewayClientsProvider GatewayClientsProvider
}

func NewDefaultAdminAPIServicesProvider(gatewaysProvider GatewayClientsProvider) *DefaultAdminAPIServicesProvider {
	return &DefaultAdminAPIServicesProvider{gatewayClientsProvider: gatewaysProvider}
}

func (p DefaultAdminAPIServicesProvider) GetConsumersService() (kong.AbstractConsumerService, bool) {
	c, ok := p.designatedAdminAPIClient()
	if !ok {
		return nil, ok
	}
	return c.Consumers, true
}

func (p DefaultAdminAPIServicesProvider) GetPluginsService() (kong.AbstractPluginService, bool) {
	c, ok := p.designatedAdminAPIClient()
	if !ok {
		return nil, ok
	}
	return c.Plugins, true
}

func (p DefaultAdminAPIServicesProvider) designatedAdminAPIClient() (*kong.Client, bool) {
	gwClients := p.gatewayClientsProvider.GatewayClients()
	if len(gwClients) == 0 {
		return nil, false
	}

	// For now using first client is kind of OK. Using Consumer and Plugin
	// services from first kong client should theoretically return the same
	// results as for all other clients. There might be instances where
	// configurations in different Kong Gateways are ever so slightly
	// different but that shouldn't cause a fatal failure.
	//
	// TODO: We should take a look at this sooner rather than later.
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3363
	return gwClients[0].AdminAPIClient(), true
}
