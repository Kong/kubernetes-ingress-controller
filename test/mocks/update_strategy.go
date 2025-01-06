package mocks

import (
	"context"

	"github.com/samber/mo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

// UpdateStrategy is a mock implementation of sendconfig.UpdateStrategy.
type UpdateStrategy struct {
	onUpdate func(content sendconfig.ContentWithHash) (mo.Option[int], error)
}

func (m *UpdateStrategy) Update(_ context.Context, targetContent sendconfig.ContentWithHash) (n mo.Option[int], err error) {
	return m.onUpdate(targetContent)
}

func (m *UpdateStrategy) MetricsProtocol() metrics.Protocol {
	return metrics.ProtocolDBLess
}

func (m *UpdateStrategy) Type() string {
	return "Mock"
}
