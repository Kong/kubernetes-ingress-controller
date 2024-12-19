package gateway

import (
	"fmt"
	"reflect"

	"github.com/samber/lo"
	"github.com/samber/mo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
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
	case *gatewayapi.HTTPRoute,
		*gatewayapi.GRPCRoute:
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

// parentStatusesForRoute returns route parent statuses for the given route
// and its supported gateways.
// It returns a map of parent statuses and a boolean indicating whether any
// changes were made.
func parentStatusesForRoute[routeT gatewayapi.RouteT](
	route routeT,
	routeParentStatuses []gatewayapi.RouteParentStatus,
	gateways ...supportedGatewayWithCondition,
) (map[string]*gatewayapi.RouteParentStatus, bool) {
	// map the existing parentStatues to avoid duplications
	parentStatuses := getParentStatuses(route, routeParentStatuses)

	// overlay the parent ref statuses for all new gateway references
	statusChangesWereMade := false
	for _, gateway := range gateways {

		// build a new status for the parent Gateway
		gatewayParentStatus := gatewayParentStatusForRoute(route, gateway, withSectionName(gateway.listenerName))

		// if the reference already exists and doesn't require any changes
		// then just leave it alone.
		parentRefKey := routeParentStatusKey(route, gateway)
		if existing, ok := parentStatuses[parentRefKey]; ok {
			// check if the parentRef and controllerName are equal, and whether
			// the new condition is present in existing conditions
			if reflect.DeepEqual(existing.ParentRef, gatewayParentStatus.ParentRef) &&
				existing.ControllerName == gatewayParentStatus.ControllerName &&
				lo.ContainsBy(existing.Conditions, func(condition metav1.Condition) bool {
					return sameCondition(gatewayParentStatus.Conditions[0], condition)
				}) {
				continue
			}
		}

		// otherwise overlay the new status on top the list of parentStatuses
		parentStatuses[parentRefKey] = gatewayParentStatus
		statusChangesWereMade = true
	}
	return parentStatuses, statusChangesWereMade
}

func withSectionName(sectionName string) func(*gatewayapi.RouteParentStatus) {
	return func(routeParentStatus *gatewayapi.RouteParentStatus) {
		if sectionName != "" {
			routeParentStatus.ParentRef.SectionName = (*gatewayapi.SectionName)(&sectionName)
		}
	}
}

func gatewayParentStatusForRoute[routeT gatewayapi.RouteT](
	route routeT,
	parentGateway supportedGatewayWithCondition,
	opts ...func(*gatewayapi.RouteParentStatus),
) *gatewayapi.RouteParentStatus {
	parentGVK := parentGateway.gateway.GroupVersionKind()
	if parentGVK.Kind == "" {
		parentGVK.Kind = gatewayapi.V1GatewayTypeMeta.Kind
	}
	if parentGVK.Group == "" {
		parentGVK.Group = gatewayapi.V1GatewayTypeMeta.GroupVersionKind().Group
		parentGateway.gateway.SetGroupVersionKind(parentGVK)
	}

	var (
		parentRef = gatewayapi.ParentReference{
			Group:     util.StringToTypedPtr[*gatewayapi.Group](parentGateway.gateway.GroupVersionKind().Group),
			Kind:      util.StringToTypedPtr[*gatewayapi.Kind](parentGateway.gateway.Kind),
			Namespace: (*gatewayapi.Namespace)(&parentGateway.gateway.Namespace),
			Name:      gatewayapi.ObjectName(parentGateway.gateway.Name),
		}
		routeParentStatus = &gatewayapi.RouteParentStatus{
			ParentRef:      parentRef,
			ControllerName: GetControllerName(),
			Conditions: []metav1.Condition{
				{
					Type:               parentGateway.condition.Type,
					Status:             parentGateway.condition.Status,
					ObservedGeneration: route.GetGeneration(),
					LastTransitionTime: metav1.Now(),
					Reason:             parentGateway.condition.Reason,
				},
			},
		}
	)

	for _, opt := range opts {
		opt(routeParentStatus)
	}

	return routeParentStatus
}

func initializeParentStatusesWithProgrammedCondition[routeT gatewayapi.RouteT](
	route routeT,
	parentStatuses map[string]*gatewayapi.RouteParentStatus,
) bool {
	// do not update the condition if a "Programmed" condition is already present.
	changed := false
	programmedConditionUnknown := newProgrammedConditionUnknown(route)
	for _, ps := range parentStatuses {
		if !parentStatusHasProgrammedCondition(ps) {
			ps.Conditions = append(ps.Conditions, programmedConditionUnknown)
			changed = true
		}
	}
	return changed
}
