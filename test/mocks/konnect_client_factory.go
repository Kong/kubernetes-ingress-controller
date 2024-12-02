package mocks

import (
	"context"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
)

type KonnectClientFactory struct {
	Client *adminapi.KonnectClient
}

func (f *KonnectClientFactory) NewKonnectClient(context.Context) (*adminapi.KonnectClient, error) {
	return f.Client, nil
}
