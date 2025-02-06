package workflows

import (
	"context"
	"errors"
	"fmt"

	"github.com/kong/kubernetes-telemetry/pkg/provider"
	"github.com/kong/kubernetes-telemetry/pkg/telemetry"
	"github.com/kong/kubernetes-telemetry/pkg/types"
)

const GatewayDiscoveryWorkflowName = "gateway_discovery"

// DiscoveredGatewaysCounter is an interface that allows to count currently discovered Gateways.
type DiscoveredGatewaysCounter interface {
	GatewayClientsCount() int
}

func NewGatewayDiscoveryWorkflow(gatewaysCounter DiscoveredGatewaysCounter) (telemetry.Workflow, error) {
	w := telemetry.NewWorkflow(GatewayDiscoveryWorkflowName)

	discoveredGatewaysCountProvider, err := NewDiscoveredGatewaysCountProvider(gatewaysCounter)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovered gateways count provider: %w", err)
	}
	w.AddProvider(discoveredGatewaysCountProvider)

	return w, nil
}

// DiscoveredGatewaysCountProvider is a provider that reports the number of currently discovered Gateways.
type DiscoveredGatewaysCountProvider struct {
	counter DiscoveredGatewaysCounter
}

func NewDiscoveredGatewaysCountProvider(counter DiscoveredGatewaysCounter) (*DiscoveredGatewaysCountProvider, error) {
	if counter == nil {
		return nil, errors.New("discovered gateways counter is required")
	}

	return &DiscoveredGatewaysCountProvider{counter: counter}, nil
}

const (
	DiscoveredGatewaysCountProviderName = "discovered_gateways_count"
	DiscoveredGatewaysCountProviderKind = provider.Kind(DiscoveredGatewaysCountProviderName)
	DiscoveredGatewaysCountKey          = types.ProviderReportKey(DiscoveredGatewaysCountProviderName)
)

func (d *DiscoveredGatewaysCountProvider) Name() string {
	return DiscoveredGatewaysCountProviderName
}

func (d *DiscoveredGatewaysCountProvider) Kind() provider.Kind {
	return DiscoveredGatewaysCountProviderKind
}

func (d *DiscoveredGatewaysCountProvider) Provide(context.Context) (types.ProviderReport, error) {
	return types.ProviderReport{
		DiscoveredGatewaysCountKey: d.counter.GatewayClientsCount(),
	}, nil
}
