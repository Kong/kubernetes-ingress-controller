package builder

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func NewAllowedRoutesFromSameNamespaces() *gatewayv1beta1.AllowedRoutes {
	return &gatewayv1beta1.AllowedRoutes{
		Namespaces: NewRouteNamespaces().FromSame().Build(),
	}
}

func NewAllowedRoutesFromAllNamespaces() *gatewayv1beta1.AllowedRoutes {
	return &gatewayv1beta1.AllowedRoutes{
		Namespaces: NewRouteNamespaces().FromAll().Build(),
	}
}

func NewAllowedRoutesFromSelectorNamespace(selector *metav1.LabelSelector) *gatewayv1beta1.AllowedRoutes {
	return &gatewayv1beta1.AllowedRoutes{
		Namespaces: NewRouteNamespaces().FromSelector(selector).Build(),
	}
}
