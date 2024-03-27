package gateway

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/samber/mo"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/types"
)

// getParentStatuses creates a parent status map for the provided route given the
// route parent status slice.
func getParentStatuses[routeT types.RouteT](
	route routeT, parentStatuses []gatewayv1.RouteParentStatus,
) map[string]*gatewayv1.RouteParentStatus {
	m := make(map[string]*gatewayv1.RouteParentStatus)

	for _, existingParent := range parentStatuses {
		parentRef := getParentRef(existingParent)
		key := routeParentStatusKey(route, parentRef)

		existingParentCopy := existingParent
		m[key] = &existingParentCopy
	}
	return m
}

type namespacedNamer interface {
	GetNamespace() string
	GetName() string
	GetSectionName() mo.Option[string]
}

func routeParentStatusKey[routeT types.RouteT](
	route routeT, parentRef namespacedNamer,
) string {
	namespace := route.GetNamespace()
	if ns := parentRef.GetNamespace(); ns != "" {
		namespace = ns
	}

	switch any(route).(type) {
	case *gatewayv1.HTTPRoute:
		return fmt.Sprintf("%s/%s/%s",
			namespace,
			parentRef.GetName(),
			parentRef.GetSectionName().OrEmpty())
	default:
		return fmt.Sprintf("%s/%s", namespace, parentRef.GetName())
	}
}

type parentRef struct {
	Namespace   *string
	Name        string
	SectionName *string
}

func (p parentRef) GetName() string {
	return p.Name
}

func (p parentRef) GetNamespace() string {
	if p.Namespace != nil {
		return *p.Namespace
	}
	return ""
}

func (p parentRef) GetSectionName() mo.Option[string] {
	if p.SectionName != nil {
		return mo.Some(*p.SectionName)
	}
	return mo.None[string]()
}

// getParentRef serves as glue code to generically get parentRef from either
// gatewayv1alpha2.RouteParentStatus or gatewayv1.RouteParentStatus.
func getParentRef(parentStatus gatewayv1.RouteParentStatus) parentRef {
	var sectionName *string

	if parentStatus.ParentRef.SectionName != nil {
		sectionName = lo.ToPtr(string(*parentStatus.ParentRef.SectionName))
	}
	return parentRef{
		Namespace:   lo.ToPtr(string(*parentStatus.ParentRef.Namespace)),
		Name:        string(parentStatus.ParentRef.Name),
		SectionName: sectionName,
	}
}
