package gateway

import (
	"context"

	"github.com/samber/lo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func newCondition(
	typ string,
	status metav1.ConditionStatus,
	reason string,
	generation int64,
) metav1.Condition {
	return metav1.Condition{
		Type:               typ,
		Status:             status,
		Reason:             reason,
		ObservedGeneration: generation,
		LastTransitionTime: metav1.Now(),
	}
}

func newProgrammedConditionUnknown(obj interface{ GetGeneration() int64 }) metav1.Condition {
	return newCondition(
		ConditionTypeProgrammed,
		metav1.ConditionUnknown,
		string(ConditionReasonProgrammedUnknown),
		obj.GetGeneration(),
	)
}

// sameCondition returns true if the conditions in parameter has the same type, status, reason and message.
func sameCondition(a, b metav1.Condition) bool {
	return a.Type == b.Type &&
		a.Status == b.Status &&
		a.Reason == b.Reason &&
		a.Message == b.Message &&
		a.ObservedGeneration == b.ObservedGeneration
}

func setRouteParentStatusCondition(parentStatus *gatewayapi.RouteParentStatus, newCondition metav1.Condition) bool {
	var conditionFound, changed bool
	for i, condition := range parentStatus.Conditions {
		if condition.Type == newCondition.Type {
			conditionFound = true
			if !sameCondition(condition, newCondition) {
				parentStatus.Conditions[i] = newCondition
				changed = true
			}
		}
	}

	if !conditionFound {
		parentStatus.Conditions = append(parentStatus.Conditions, newCondition)
		changed = true
	}
	return changed
}

func parentStatusHasProgrammedCondition(parentStatus *gatewayapi.RouteParentStatus) bool {
	return lo.ContainsBy(parentStatus.Conditions, func(c metav1.Condition) bool {
		return c.Type == ConditionTypeProgrammed
	})
}

// ensureParentsProgrammedCondition ensures that provided route's parent statuses
// have Programmed condition set properly. It returns a boolean flag indicating
// whether an update to the provided route has been performed.
//
// Use the condition argument to specify the Reason, Status and Message.
// Type will be set to Programmed whereas ObservedGeneration and LastTransitionTime
// will be set accordingly based on the route's generation and current time.
func ensureParentsProgrammedCondition[
	routeT gatewayapi.RouteT,
](
	ctx context.Context,
	client client.SubResourceWriter,
	route routeT,
	routeParentStatuses []gatewayapi.RouteParentStatus,
	gateways []supportedGatewayWithCondition,
	condition metav1.Condition,
) (bool, error) {
	// map the existing parentStatues to avoid duplications
	parentStatuses := getParentStatuses(route, routeParentStatuses)

	condition.Type = ConditionTypeProgrammed
	condition.ObservedGeneration = route.GetGeneration()
	condition.LastTransitionTime = metav1.Now()

	statusChanged := false
	for _, g := range gateways {
		gateway := g.gateway

		parentRefKey := routeParentStatusKey(route, g)
		parentStatus, ok := parentStatuses[parentRefKey]
		if ok {
			// update existing parent in status.
			changed := setRouteParentStatusCondition(parentStatus, condition)
			if changed {
				parentStatuses[parentRefKey] = parentStatus
				setRouteParentInStatusForParent(route, *parentStatus, g)
			}
			statusChanged = statusChanged || changed
		} else {
			// add a new parent if the parent is not found in status.
			newParentStatus := gatewayapi.RouteParentStatus{
				ParentRef: gatewayapi.ParentReference{
					Namespace: lo.ToPtr(gatewayapi.Namespace(gateway.Namespace)),
					Name:      gatewayapi.ObjectName(gateway.Name),
					Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
					Group:     lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
					// We don't need to check whether the listener matches route's spec
					// because that should already be done via getSupportedGatewayForRoute
					// at this point.
					SectionName: lo.EmptyableToPtr(gatewayapi.SectionName(g.listenerName)),

					// TODO: set port after gateway port matching implemented:
					// https://github.com/Kong/kubernetes-ingress-controller/issues/3016
				},
				ControllerName: GetControllerName(),
				Conditions: []metav1.Condition{
					condition,
				},
			}
			setRouteParentInStatusForParent(route, newParentStatus, g)

			routeParentStatuses = append(routeParentStatuses, newParentStatus)
			parentStatuses[parentRefKey] = &newParentStatus
			statusChanged = true
		}
	}

	// update status if needed.
	if statusChanged {
		if err := client.Update(ctx, route); err != nil {
			return false, err
		}
		return true, nil
	}
	// no need to update if no status is changed.
	return false, nil
}
