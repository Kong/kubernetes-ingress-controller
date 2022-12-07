package builder

import (
	"k8s.io/utils/pointer"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// BackendRefBuilder is a builder for gateway api BackendRef.
// Will set default values, as specified in the gateway API, for fields that are not set.
// Primarily used for testing.
type BackendRefBuilder struct {
	backendRef gatewayv1alpha2.BackendRef
}

func NewBackendRef(name string) *BackendRefBuilder {
	return &BackendRefBuilder{
		backendRef: gatewayv1alpha2.BackendRef{
			BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
				Name: gatewayv1alpha2.ObjectName(name),
				Kind: util.StringToGatewayAPIKindV1Alpha2Ptr("Service"), // default value
			},
		},
	}
}

func (b *BackendRefBuilder) Build() gatewayv1alpha2.BackendRef {
	return b.backendRef
}

func (b *BackendRefBuilder) ToSlice() []gatewayv1alpha2.BackendRef {
	return []gatewayv1alpha2.BackendRef{b.backendRef}
}

func (b *BackendRefBuilder) WithPort(port int) *BackendRefBuilder {
	val := gatewayv1alpha2.PortNumber(port)
	b.backendRef.Port = &val
	return b
}

func (b *BackendRefBuilder) WithWeight(weight int) *BackendRefBuilder {
	b.backendRef.Weight = pointer.Int32(int32(weight))
	return b
}

func (b *BackendRefBuilder) WithKind(kind string) *BackendRefBuilder {
	val := gatewayv1alpha2.Kind(kind)
	b.backendRef.Kind = &val
	return b
}

func (b *BackendRefBuilder) WithGroup(group string) *BackendRefBuilder {
	val := gatewayv1alpha2.Group(group)
	b.backendRef.Group = &val
	return b
}

func (b *BackendRefBuilder) WithNamespace(namespace string) *BackendRefBuilder {
	val := gatewayv1alpha2.Namespace(namespace)
	b.backendRef.Namespace = &val
	return b
}
