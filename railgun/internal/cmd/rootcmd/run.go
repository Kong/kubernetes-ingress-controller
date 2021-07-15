package rootcmd

import (
	"context"
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/railgun/internal/manager"
)

func Run(ctx context.Context, c *manager.Config) error {
	if err := StartAdmissionServer(ctx, c); err != nil {
		return fmt.Errorf("StartAdmissionServer: %w", err)
	}
	if err := StartProfilingServer(ctx, c); err != nil {
		return fmt.Errorf("StartProfilingServer: %w", err)
	}
	return manager.Run(ctx, c)
}
