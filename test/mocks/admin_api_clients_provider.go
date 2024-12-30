package mocks

import (
	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
)

// MockGatewayClientsProvider is a mock implementation of dataplane.AdminAPIClientsProvider.
type MockGatewayClientsProvider struct {
	GatewayClientList     []*adminapi.Client
	KonnectClientInstance *adminapi.KonnectClient
	DBMode                dpconf.DBMode
}

func (p *MockGatewayClientsProvider) KonnectClient() *adminapi.KonnectClient {
	return p.KonnectClientInstance
}

func (p *MockGatewayClientsProvider) GatewayClients() []*adminapi.Client {
	return p.GatewayClientList
}

func (p *MockGatewayClientsProvider) GatewayClientsToConfigure() []*adminapi.Client {
	if p.DBMode.IsDBLessMode() {
		return p.GatewayClientList
	}
	if len(p.GatewayClientList) == 0 {
		return []*adminapi.Client{}
	}
	return p.GatewayClientList[:1]
}
