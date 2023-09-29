package builder

import (
	"github.com/samber/lo"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
)

// RouteGroupKindBuilder is a builder for gateway api RouteGroupKind.
// Will set default values, as specified in the gateway API, for fields that are not set.
// Primarily used for testing.
type RouteGroupKindBuilder struct {
	routeGroupKind gatewayapi.RouteGroupKind
}

func NewRouteGroupKind() *RouteGroupKindBuilder {
	return &RouteGroupKindBuilder{
		routeGroupKind: gatewayapi.RouteGroupKind{
			Group: lo.ToPtr(gatewayapi.Group(gatewayv1beta1.GroupVersion.Group)),
		},
	}
}

// Build returns the configured RouteGroupKind.
func (b *RouteGroupKindBuilder) Build() gatewayapi.RouteGroupKind {
	return b.routeGroupKind
}

// IntoSlice returns the configured RouteGroupKind in a slice.
func (b *RouteGroupKindBuilder) IntoSlice() []gatewayapi.RouteGroupKind {
	return []gatewayapi.RouteGroupKind{b.routeGroupKind}
}

func (b *RouteGroupKindBuilder) TCPRoute() *RouteGroupKindBuilder {
	b.routeGroupKind.Kind = "TCPRoute"
	return b
}

func (b *RouteGroupKindBuilder) HTTPRoute() *RouteGroupKindBuilder {
	b.routeGroupKind.Kind = "HTTPRoute"
	return b
}

func (b *RouteGroupKindBuilder) UDPRoute() *RouteGroupKindBuilder {
	b.routeGroupKind.Kind = "UDPRoute"
	return b
}

func (b *RouteGroupKindBuilder) TLSRoute() *RouteGroupKindBuilder {
	b.routeGroupKind.Kind = "TLSRoute"
	return b
}

func (b *RouteGroupKindBuilder) GRPCRoute() *RouteGroupKindBuilder {
	b.routeGroupKind.Kind = "GRPCRoute"
	return b
}
