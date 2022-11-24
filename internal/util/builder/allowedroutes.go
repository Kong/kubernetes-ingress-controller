package builder

import gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

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
