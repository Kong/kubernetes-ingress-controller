package builder

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// RouteNamespacesBuilder is a builder for gateway api RouteNamespaces.
// Will set default values, as specified in the gateway API, for fields that are not set.
// Primarily used for testing.
type RouteNamespacesBuilder struct {
	routeNamespaces gatewayv1beta1.RouteNamespaces
}

func NewRouteNamespaces() *RouteNamespacesBuilder {
	return &RouteNamespacesBuilder{}
}

// Build returns the configured RouteNamespaces.
func (b *RouteNamespacesBuilder) Build() *gatewayv1beta1.RouteNamespaces {
	return &b.routeNamespaces
}

func (b *RouteNamespacesBuilder) FromSame() *RouteNamespacesBuilder {
	b.routeNamespaces.From = addressOf(gatewayv1beta1.NamespacesFromSame)
	return b
}

func (b *RouteNamespacesBuilder) FromAll() *RouteNamespacesBuilder {
	b.routeNamespaces.From = addressOf(gatewayv1beta1.NamespacesFromAll)
	return b
}

func (b *RouteNamespacesBuilder) FromSelector(s *metav1.LabelSelector) *RouteNamespacesBuilder {
	b.routeNamespaces.From = addressOf(gatewayv1beta1.NamespacesFromSelector)
	b.routeNamespaces.Selector = s
	return b
}
