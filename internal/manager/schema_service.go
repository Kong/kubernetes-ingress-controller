package manager

import (
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
)

// GatewayClientsProvider is an interface that provides clients for the currently discovered Gateway instances.
type GatewayClientsProvider interface {
	GatewayClients() []*adminapi.Client
}

// SchemaServiceGetter returns schema service of an admin API client if there is any client available.
type SchemaServiceGetter struct {
	clientsManager GatewayClientsProvider
}

// NewSchemaServiceGetter creates a schema service getter that uses given client manager to maintain admin API clients.
func NewSchemaServiceGetter(cm GatewayClientsProvider) SchemaServiceGetter {
	return SchemaServiceGetter{
		clientsManager: cm,
	}
}

// GetSchemaService returns the schema service for admin API client.
// It uses the configured clients manager to get the clients and then it uses one of those to obtain the service.
func (ssg SchemaServiceGetter) GetSchemaService() kong.AbstractSchemaService {
	clients := ssg.clientsManager.GatewayClients()
	if len(clients) > 0 {
		return clients[0].AdminAPIClient().Schemas
	}
	// returns a fake schema service when no gateway clients available.
	return translator.UnavailableSchemaService{}
}
