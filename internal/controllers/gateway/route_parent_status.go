package gateway

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/samber/mo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// getParentStatuses creates a parent status map for the provided route given the
// route parent status slice.
func getParentStatuses[routeT gatewayapi.RouteT](
	route routeT, parentStatuses []gatewayapi.RouteParentStatus,
) map[string]*gatewayapi.RouteParentStatus {
	m := make(map[string]*gatewayapi.RouteParentStatus)

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

func routeParentStatusKey[routeT gatewayapi.RouteT](
	route routeT, parentRef namespacedNamer,
) string {
	namespace := route.GetNamespace()
	if ns := parentRef.GetNamespace(); ns != "" {
		namespace = ns
	}

	switch any(route).(type) {
	case *gatewayapi.HTTPRoute:
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
// gatewayapi.RouteParentStatus or gatewayapi.RouteParentStatus.
func getParentRef(parentStatus gatewayapi.RouteParentStatus) parentRef {
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

func getRouteStatusParents[T gatewayapi.RouteT](route T) []gatewayapi.RouteParentStatus {
	switch r := any(route).(type) {
	case *gatewayapi.HTTPRoute:
		return r.Status.Parents
	case *gatewayapi.TCPRoute:
		return r.Status.Parents
	case *gatewayapi.UDPRoute:
		return r.Status.Parents
	case *gatewayapi.TLSRoute:
		return r.Status.Parents
	case *gatewayapi.GRPCRoute:
		return r.Status.Parents
	default:
		return nil
	}
}

func setRouteStatusParents[T gatewayapi.RouteT](route T, parents []gatewayapi.RouteParentStatus) {
	switch r := any(route).(type) {
	case *gatewayapi.HTTPRoute:
		r.Status.Parents = parents
	case *gatewayapi.TCPRoute:
		r.Status.Parents = parents
	case *gatewayapi.UDPRoute:
		r.Status.Parents = parents
	case *gatewayapi.TLSRoute:
		r.Status.Parents = parents
	case *gatewayapi.GRPCRoute:
		r.Status.Parents = parents
	}
}
