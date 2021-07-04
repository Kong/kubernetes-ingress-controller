package rootcmd

import (
	"context"
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/railgun/manager"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/config"
)

func Run(ctx context.Context, c *config.Config) error {
	if err := StartAdmissionServer(ctx, c); err != nil {
		return fmt.Errorf("StartAdmissionServer: %w", err)
	}
	if err := StartProfilingServer(ctx, c); err != nil {
		return fmt.Errorf("StartProfilingServer: %w", err)
	}
	err := c.UpdateKongAdminURL(ctx)
	if err != nil {
		return fmt.Errorf("failed retrieving kong admin URL. err %v", err)
	}
	return manager.Run(ctx, c)
}
