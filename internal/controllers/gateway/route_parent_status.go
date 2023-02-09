package gateway

import (
	"fmt"

	"github.com/samber/lo"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

type namespacedObjectT interface {
	GetNamespace() string
}

// getParentStatuses creates a parent status map for the provided route given the
// route parent status slice.
func getParentStatuses[routeT namespacedObjectT](
	route routeT, parentStatuses []gatewayv1beta1.RouteParentStatus,
) map[string]*gatewayv1beta1.RouteParentStatus {
	var (
		namespace = route.GetNamespace()
		m         = make(map[string]*gatewayv1beta1.RouteParentStatus)
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
func getParentRef(parentStatus gatewayv1beta1.RouteParentStatus) parentRef {
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
