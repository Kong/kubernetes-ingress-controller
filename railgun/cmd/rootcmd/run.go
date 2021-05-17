package rootcmd

import (
	"context"

	"github.com/kong/kubernetes-ingress-controller/railgun/manager"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/config"
)

func Run(ctx context.Context, c *config.Config) error {
	return manager.Run(ctx, c)
}
