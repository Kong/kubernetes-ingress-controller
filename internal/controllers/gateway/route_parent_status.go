package gateway

import (
	"fmt"

	"github.com/samber/lo"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

type RouteParentStatusT interface {
	gatewayv1beta1.RouteParentStatus | gatewayv1alpha2.RouteParentStatus
}

type namespacedObjectT interface {
	GetNamespace() string
}

// getParentStatuses creates a parent status map for the provided route given the
// route parent status slice.
func getParentStatuses[routeT namespacedObjectT, parentStatusT RouteParentStatusT](
	route routeT, parentStatuses []parentStatusT,
) map[string]*parentStatusT {
	var (
		namespace = route.GetNamespace()
		m         = make(map[string]*parentStatusT)
	)

	for _, existingParent := range parentStatuses {
		parentRef := getParentRef(existingParent)

		if parentRef.Namespace != nil {
			namespace = *parentRef.Namespace
		}
		var sectionName string
		if parentRef.SectionName != nil {
			sectionName = *parentRef.SectionName
		}

		var key string
		switch any(route).(type) {
		case *gatewayv1beta1.HTTPRoute:
			key = fmt.Sprintf("%s/%s/%s", namespace, parentRef.Name, sectionName)
		default:
			key = fmt.Sprintf("%s/%s", namespace, parentRef.Name)
		}

		existingParentCopy := existingParent
		m[key] = &existingParentCopy
	}
	return m
}

type parentRef struct {
	Namespace   *string
	Name        string
	SectionName *string
}

// getParentRef serves as glue code to generically get parentRef from either
// gatewayv1alpha2.RouteParentStatus or gatewayv1beta1.RouteParentStatus.
func getParentRef[T RouteParentStatusT](parentStatus T) parentRef {
	var sectionName *string

	switch ps := any(parentStatus).(type) {
	case gatewayv1beta1.RouteParentStatus:
		if ps.ParentRef.SectionName != nil {
			sectionName = lo.ToPtr(string(*ps.ParentRef.SectionName))
		}
		return parentRef{
			Namespace:   lo.ToPtr(string(*ps.ParentRef.Namespace)),
			Name:        string(ps.ParentRef.Name),
			SectionName: sectionName,
		}
	case gatewayv1alpha2.RouteParentStatus:
		if ps.ParentRef.SectionName != nil {
			sectionName = lo.ToPtr(string(*ps.ParentRef.SectionName))
		}
		return parentRef{
			Namespace:   lo.ToPtr(string(*ps.ParentRef.Namespace)),
			Name:        string(ps.ParentRef.Name),
			SectionName: sectionName,
		}
	}
	return parentRef{}
}
