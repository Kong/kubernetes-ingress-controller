package builder

import (
	"github.com/samber/lo"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// HTTPBackendRefBuilder is a builder for gateway api HTTPBackendRef.
// Will set default values, as specified in the gateway API, for fields that are not set.
// Primarily used for testing.
type HTTPBackendRefBuilder struct {
	httpBackendRef gatewayv1.HTTPBackendRef
}

func NewHTTPBackendRef(name string) *HTTPBackendRefBuilder {
	return &HTTPBackendRefBuilder{
		httpBackendRef: gatewayv1.HTTPBackendRef{
			BackendRef: gatewayv1.BackendRef{
				BackendObjectReference: gatewayv1.BackendObjectReference{
					Name: gatewayv1.ObjectName(name),
					Kind: util.StringToGatewayAPIKindPtr("Service"), // default value
					Port: lo.ToPtr(gatewayv1.PortNumber(80)),
				},
			},
		},
	}
}

func (b *HTTPBackendRefBuilder) Build() gatewayv1.HTTPBackendRef {
	return b.httpBackendRef
}

func (b *HTTPBackendRefBuilder) ToSlice() []gatewayv1.HTTPBackendRef {
	return []gatewayv1.HTTPBackendRef{b.httpBackendRef}
}

func (b *HTTPBackendRefBuilder) WithPort(port int) *HTTPBackendRefBuilder {
	val := gatewayv1.PortNumber(port)
	b.httpBackendRef.Port = &val
	return b
}

func (b *HTTPBackendRefBuilder) WithWeight(weight int) *HTTPBackendRefBuilder {
	b.httpBackendRef.Weight = lo.ToPtr(int32(weight))
	return b
}

func (b *HTTPBackendRefBuilder) WithKind(kind string) *HTTPBackendRefBuilder {
	val := gatewayv1.Kind(kind)
	b.httpBackendRef.Kind = &val
	return b
}

func (b *HTTPBackendRefBuilder) WithGroup(group string) *HTTPBackendRefBuilder {
	val := gatewayv1.Group(group)
	b.httpBackendRef.Group = &val
	return b
}

func (b *HTTPBackendRefBuilder) WithNamespace(namespace string) *HTTPBackendRefBuilder {
	val := gatewayv1.Namespace(namespace)
	b.httpBackendRef.Namespace = &val
	return b
}
