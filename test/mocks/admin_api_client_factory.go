package mocks

import (
	"context"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
)

// AdminAPIClientFactory is a mock implementation of adminapi.ClientFactory.
type AdminAPIClientFactory struct {
	errorsToReturn map[string]error // map from address to error
}

func NewAdminAPIClientFactory(errorsToReturn map[string]error) *AdminAPIClientFactory {
	if errorsToReturn == nil {
		errorsToReturn = make(map[string]error)
	}
	return &AdminAPIClientFactory{
		errorsToReturn: errorsToReturn,
	}
}

func (m *AdminAPIClientFactory) CreateAdminAPIClient(_ context.Context, api adminapi.DiscoveredAdminAPI) (*adminapi.Client, error) {
	err, ok := m.errorsToReturn[api.Address]
	if !ok {
		return adminapi.NewTestClient(api.Address)
	}
	return nil, err
}
