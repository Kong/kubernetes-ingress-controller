package builder

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func NewAllowedRoutesFromSameNamespaces() *gatewayv1.AllowedRoutes {
	return &gatewayv1.AllowedRoutes{
		Namespaces: NewRouteNamespaces().FromSame().Build(),
	}
}

func NewAllowedRoutesFromAllNamespaces() *gatewayv1.AllowedRoutes {
	return &gatewayv1.AllowedRoutes{
		Namespaces: NewRouteNamespaces().FromAll().Build(),
	}
}

func NewAllowedRoutesFromSelectorNamespace(selector *metav1.LabelSelector) *gatewayv1.AllowedRoutes {
	return &gatewayv1.AllowedRoutes{
		Namespaces: NewRouteNamespaces().FromSelector(selector).Build(),
	}
}
