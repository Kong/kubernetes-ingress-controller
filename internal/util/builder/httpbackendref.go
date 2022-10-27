package builder

import (
	"k8s.io/utils/pointer"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// HTTPBackendRefBuilder is a builder for gateway api HTTPBackendRef.
// Will set default values, as specified in the gateway API, for fields that are not set.
// Primarily used for testing.
type HTTPBackendRefBuilder struct {
	httpBackendRef gatewayv1beta1.HTTPBackendRef
}

func NewHTTPBackendRef(name string) *HTTPBackendRefBuilder {
	return &HTTPBackendRefBuilder{
		httpBackendRef: gatewayv1beta1.HTTPBackendRef{
			BackendRef: gatewayv1beta1.BackendRef{
				BackendObjectReference: gatewayv1beta1.BackendObjectReference{
					Name: gatewayv1beta1.ObjectName(name),
					Kind: util.StringToGatewayAPIKindPtr("Service"), // default value
				},
			},
		},
	}
}

func (b *HTTPBackendRefBuilder) Build() gatewayv1beta1.HTTPBackendRef {
	return b.httpBackendRef
}

func (b *HTTPBackendRefBuilder) WithPort(port int) *HTTPBackendRefBuilder {
	val := gatewayv1beta1.PortNumber(port)
	b.httpBackendRef.Port = &val
	return b
}

func (b *HTTPBackendRefBuilder) WithWeight(weight int) *HTTPBackendRefBuilder {
	b.httpBackendRef.Weight = pointer.Int32(int32(weight))
	return b
}

func (b *HTTPBackendRefBuilder) WithKind(kind string) *HTTPBackendRefBuilder {
	val := gatewayv1beta1.Kind(kind)
	b.httpBackendRef.Kind = &val
	return b
}

func (b *HTTPBackendRefBuilder) WithGroup(group string) *HTTPBackendRefBuilder {
	val := gatewayv1beta1.Group(group)
	b.httpBackendRef.Group = &val
	return b
}

func (b *HTTPBackendRefBuilder) WithNamespace(namespace string) *HTTPBackendRefBuilder {
	val := gatewayv1beta1.Namespace(namespace)
	b.httpBackendRef.Namespace = &val
	return b
}
