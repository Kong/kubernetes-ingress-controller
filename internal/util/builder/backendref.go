package builder

import (
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// BackendRefBuilder is a builder for gateway api BackendRef.
// Will set default values, as specified in the gateway API, for fields that are not set.
// Primarily used for testing.
type BackendRefBuilder struct {
	backendRef gatewayapi.BackendRef
}

func NewBackendRef(name string) *BackendRefBuilder {
	return &BackendRefBuilder{
		backendRef: gatewayapi.BackendRef{
			BackendObjectReference: gatewayapi.BackendObjectReference{
				Name: gatewayapi.ObjectName(name),
				Kind: util.StringToGatewayAPIKindPtr("Service"), // default value
			},
		},
	}
}

func (b *BackendRefBuilder) Build() gatewayapi.BackendRef {
	return b.backendRef
}

func (b *BackendRefBuilder) ToSlice() []gatewayapi.BackendRef {
	return []gatewayapi.BackendRef{b.backendRef}
}

func (b *BackendRefBuilder) WithPort(port int) *BackendRefBuilder {
	val := gatewayapi.PortNumber(port)
	b.backendRef.Port = &val
	return b
}

func (b *BackendRefBuilder) WithWeight(weight int) *BackendRefBuilder {
	b.backendRef.Weight = lo.ToPtr(int32(weight))
	return b
}

func (b *BackendRefBuilder) WithKind(kind string) *BackendRefBuilder {
	val := gatewayapi.Kind(kind)
	b.backendRef.Kind = &val
	return b
}

func (b *BackendRefBuilder) WithGroup(group string) *BackendRefBuilder {
	val := gatewayapi.Group(group)
	b.backendRef.Group = &val
	return b
}

func (b *BackendRefBuilder) WithNamespace(namespace string) *BackendRefBuilder {
	val := gatewayapi.Namespace(namespace)
	b.backendRef.Namespace = &val
	return b
}
