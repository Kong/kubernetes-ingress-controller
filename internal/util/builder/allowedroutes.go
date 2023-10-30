package builder

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func NewAllowedRoutesFromSameNamespaces() *gatewayapi.AllowedRoutes {
	return &gatewayapi.AllowedRoutes{
		Namespaces: NewRouteNamespaces().FromSame().Build(),
	}
}

func NewAllowedRoutesFromAllNamespaces() *gatewayapi.AllowedRoutes {
	return &gatewayapi.AllowedRoutes{
		Namespaces: NewRouteNamespaces().FromAll().Build(),
	}
}

func NewAllowedRoutesFromSelectorNamespace(selector *metav1.LabelSelector) *gatewayapi.AllowedRoutes {
	return &gatewayapi.AllowedRoutes{
		Namespaces: NewRouteNamespaces().FromSelector(selector).Build(),
	}
}
