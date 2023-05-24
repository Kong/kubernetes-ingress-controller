package builder

import (
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
)

// EndpointPortBuilder is a builder for discovery v1 EndpointPort.
// Primarily used for testing.
type EndpointPortBuilder struct {
	ep discoveryv1.EndpointPort
}

func NewEndpointPort(port int32) *EndpointPortBuilder {
	return &EndpointPortBuilder{
		ep: discoveryv1.EndpointPort{
			Port: lo.ToPtr(port),
		},
	}
}

// WithProtocol sets the protocol on the endpoint port.
func (b *EndpointPortBuilder) WithProtocol(proto corev1.Protocol) *EndpointPortBuilder {
	b.ep.Protocol = lo.ToPtr(proto)
	return b
}

// WithName sets the name on the endpoint port.
func (b *EndpointPortBuilder) WithName(name string) *EndpointPortBuilder {
	b.ep.Name = lo.ToPtr(name)
	return b
}

// Build returns the configured EndpointPort.
func (b *EndpointPortBuilder) Build() discoveryv1.EndpointPort {
	return b.ep
}

// IntoSlice returns the configured EndpointPort in a slice.
func (b *EndpointPortBuilder) IntoSlice() []discoveryv1.EndpointPort {
	return []discoveryv1.EndpointPort{b.ep}
}
