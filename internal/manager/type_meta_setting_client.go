package manager

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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
	if err := util.PopulateTypeMeta(obj, c.Client.Scheme()); err != nil {
		return fmt.Errorf("failed to populate type meta: %w", err)
	}
	return nil
}
