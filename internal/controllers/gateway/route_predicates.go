package gateway

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func isRouteAttachedToReconciledGateway[routeT gatewayapi.RouteT](
	cl client.Client, log logr.Logger, gatewayNN controllers.OptionalNamespacedName, obj client.Object,
) bool {
	route, ok := obj.(routeT)
	if !ok {
		kind := obj.GetObjectKind().GroupVersionKind().Kind
		log.Error(
			fmt.Errorf("unexpected object type"),
			"Route watch predicate received unexpected object type",
			"expected", kind, "found", reflect.TypeOf(obj),
		)
		return false
	}

	parentRefs := getRouteParentRefs(route)

	// If the reconciler has a GatewayNN set, only HTTPRoutes attached to that Gateway are reconciled.
	if gNN, ok := gatewayNN.Get(); ok {
		for _, parentRef := range parentRefs {
			if parentRef.Namespace != nil && string(*parentRef.Namespace) != gNN.Namespace {
				continue
			}
			if string(parentRef.Name) != gNN.Name {
				continue
			}
			if parentRef.Kind != nil && *parentRef.Kind != "Gateway" {
				continue
			}
			if parentRef.Group != nil && *parentRef.Group != gatewayapi.Group(gatewayapi.GroupVersion.Group) {
				continue
			}
			return true
		}
		return false
	}

	// If the GatewayNN is not set, all HTTPRoutes are reconciled.
	// Hence we need to check if the HTTPRoute is attached to a Gateway that is managed by this controller.
	for _, parentRef := range parentRefs {
		namespace := route.GetNamespace()
		if parentRef.Namespace != nil {
			namespace = string(*parentRef.Namespace)
		}

		kind := gatewayapi.Kind("Gateway")
		if parentRef.Kind != nil {
			kind = *parentRef.Kind
		}

		group := gatewayapi.GroupVersion.Group
		if parentRef.Group != nil {
			group = string(*parentRef.Group)
		}

		switch {
		case kind == "Gateway" && group == gatewayapi.GroupVersion.Group:
			var gateway gatewayapi.Gateway
			err := cl.Get(context.Background(), k8stypes.NamespacedName{Namespace: namespace, Name: string(parentRef.Name)}, &gateway)
			if err != nil {
				log.Error(err, "Failed to get Gateway in HTTPRoute watch")
				return false
			}

			var gatewayClass gatewayapi.GatewayClass
			err = cl.Get(context.Background(), k8stypes.NamespacedName{Name: string(gateway.Spec.GatewayClassName)}, &gatewayClass)
			if err != nil {
				log.Error(err, "Failed to get GatewayClass in HTTPRoute watch")
				return false
			}

			if isGatewayClassControlled(&gatewayClass) {
				return true
			}
		default:
			log.Error(
				fmt.Errorf("unsupported parentRef kind %s and group %s", kind, group),
				"Got an unexpected kind and group when checking route's parentRefs",
			)
			return false
		}
	}

	return false
}

func isOrWasRouteAttachedToReconciledGateway[routeT gatewayapi.RouteT](
	cl client.Client, log logr.Logger, gatewayNN controllers.OptionalNamespacedName, e event.UpdateEvent,
) bool {
	oldObj, newObj := e.ObjectOld, e.ObjectNew
	return isRouteAttachedToReconciledGateway[routeT](cl, log, gatewayNN, oldObj) ||
		isRouteAttachedToReconciledGateway[routeT](cl, log, gatewayNN, newObj)
}
