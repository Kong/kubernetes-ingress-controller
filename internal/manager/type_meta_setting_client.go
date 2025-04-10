package manager

import (
	"context"
	"fmt"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// TypeMetaSettingClient decorates client.Client so that it populates the TypeMeta field of the
// object after fetching it from the API server.
type TypeMetaSettingClient struct {
	client.Client
}

func NewTypeMetaSettingClient(c client.Client) TypeMetaSettingClient {
	return TypeMetaSettingClient{Client: c}
}

// Get retrieves an object from the Kubernetes API server and populates the TypeMeta field of the object.
func (c TypeMetaSettingClient) Get(
	ctx context.Context,
	key client.ObjectKey,
	obj client.Object,
	opts ...client.GetOption,
) error {
	if err := c.Client.Get(ctx, key, obj, opts...); err != nil {
		return err
	}
	if err := util.PopulateTypeMeta(obj, c.Scheme()); err != nil {
		return fmt.Errorf("failed to populate type meta: %w", err)
	}
	return nil
}

// newManagerClient generates a controller-runtime client and wraps it in our override decorator.
func newManagerClient(config *rest.Config, options client.Options) (client.Client, error) {
	base, err := client.New(config, options)
	if err != nil {
		return nil, err
	}
	metaSetter := NewTypeMetaSettingClient(base)
	return metaSetter, nil
}
