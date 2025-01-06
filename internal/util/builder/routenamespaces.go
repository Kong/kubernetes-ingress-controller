package builder

import (
	"github.com/samber/lo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// RouteNamespacesBuilder is a builder for gateway api RouteNamespaces.
// Will set default values, as specified in the gateway API, for fields that are not set.
// Primarily used for testing.
type RouteNamespacesBuilder struct {
	routeNamespaces gatewayapi.RouteNamespaces
}

func NewRouteNamespaces() *RouteNamespacesBuilder {
	return &RouteNamespacesBuilder{}
}

// Build returns the configured RouteNamespaces.
func (b *RouteNamespacesBuilder) Build() *gatewayapi.RouteNamespaces {
	return &b.routeNamespaces
}

func (b *RouteNamespacesBuilder) FromSame() *RouteNamespacesBuilder {
	b.routeNamespaces.From = lo.ToPtr(gatewayapi.NamespacesFromSame)
	return b
}

func (b *RouteNamespacesBuilder) FromAll() *RouteNamespacesBuilder {
	b.routeNamespaces.From = lo.ToPtr(gatewayapi.NamespacesFromAll)
	return b
}

func (b *RouteNamespacesBuilder) FromSelector(s *metav1.LabelSelector) *RouteNamespacesBuilder {
	b.routeNamespaces.From = lo.ToPtr(gatewayapi.NamespacesFromSelector)
	b.routeNamespaces.Selector = s
	return b
}
