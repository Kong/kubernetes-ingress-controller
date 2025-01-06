package builder

import (
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// HTTPBackendRefBuilder is a builder for gateway api HTTPBackendRef.
// Will set default values, as specified in the gateway API, for fields that are not set.
// Primarily used for testing.
type HTTPBackendRefBuilder struct {
	httpBackendRef gatewayapi.HTTPBackendRef
}

func NewHTTPBackendRef(name string) *HTTPBackendRefBuilder {
	return &HTTPBackendRefBuilder{
		httpBackendRef: gatewayapi.HTTPBackendRef{
			BackendRef: gatewayapi.BackendRef{
				BackendObjectReference: gatewayapi.BackendObjectReference{
					Name: gatewayapi.ObjectName(name),
					Kind: util.StringToGatewayAPIKindPtr("Service"), // default value
					Port: lo.ToPtr(gatewayapi.PortNumber(80)),
				},
			},
		},
	}
}

func (b *HTTPBackendRefBuilder) Build() gatewayapi.HTTPBackendRef {
	return b.httpBackendRef
}

func (b *HTTPBackendRefBuilder) ToSlice() []gatewayapi.HTTPBackendRef {
	return []gatewayapi.HTTPBackendRef{b.httpBackendRef}
}

func (b *HTTPBackendRefBuilder) WithPort(port int) *HTTPBackendRefBuilder {
	val := gatewayapi.PortNumber(port)
	b.httpBackendRef.Port = &val
	return b
}

func (b *HTTPBackendRefBuilder) WithWeight(weight int) *HTTPBackendRefBuilder {
	b.httpBackendRef.Weight = lo.ToPtr(int32(weight))
	return b
}

func (b *HTTPBackendRefBuilder) WithKind(kind string) *HTTPBackendRefBuilder {
	val := gatewayapi.Kind(kind)
	b.httpBackendRef.Kind = &val
	return b
}

func (b *HTTPBackendRefBuilder) WithGroup(group string) *HTTPBackendRefBuilder {
	val := gatewayapi.Group(group)
	b.httpBackendRef.Group = &val
	return b
}

func (b *HTTPBackendRefBuilder) WithNamespace(namespace string) *HTTPBackendRefBuilder {
	val := gatewayapi.Namespace(namespace)
	b.httpBackendRef.Namespace = &val
	return b
}
