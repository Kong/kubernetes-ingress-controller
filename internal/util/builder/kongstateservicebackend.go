package builder

import (
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
)

// KongstateServiceBackendBuilder is a builder for KongstateServiceBackend.
// Primarily used for testing.
type KongstateServiceBackendBuilder struct {
	kongstateServiceBackend kongstate.ServiceBackend
}

func NewKongstateServiceBackend(name string) *KongstateServiceBackendBuilder {
	return &KongstateServiceBackendBuilder{
		kongstateServiceBackend: kongstate.ServiceBackend{
			Name: name,
		},
	}
}

func (b *KongstateServiceBackendBuilder) WithNamespace(namespace string) *KongstateServiceBackendBuilder {
	b.kongstateServiceBackend.Namespace = namespace
	return b
}

func (b *KongstateServiceBackendBuilder) WithWeight(weight int) *KongstateServiceBackendBuilder {
	b.kongstateServiceBackend.Weight = lo.ToPtr(int32(weight))
	return b
}

func (b *KongstateServiceBackendBuilder) WithPortNumber(port int) *KongstateServiceBackendBuilder {
	b.kongstateServiceBackend.PortDef = kongstate.PortDef{
		Number: int32(port),
		Mode:   kongstate.PortModeByNumber,
	}
	return b
}

func (b *KongstateServiceBackendBuilder) Build() kongstate.ServiceBackend {
	return b.kongstateServiceBackend
}
