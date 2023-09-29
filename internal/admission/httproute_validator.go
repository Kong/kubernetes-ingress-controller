package admission

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gatewayvalidation "github.com/kong/kubernetes-ingress-controller/v2/internal/admission/validation/gateway"
	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
)

func (validator KongHTTPValidator) HTTPRoute() CustomValidatorAdapter {
	return CustomValidatorAdapter{
		validateCreate: func(ctx context.Context, obj runtime.Object) (bool, string, error) {
			route, ok := obj.(*gatewayapi.HTTPRoute)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *gatewayapi.HTTPRoute, got %T", obj)
			}
			return validator.ValidateHTTPRoute(ctx, *route)
		},
		validateUpdate: func(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (bool, string, error) {
			route, ok := newObj.(*gatewayapi.HTTPRoute)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *gatewayapi.HTTPRoute, got %T", newObj)
			}
			return validator.ValidateHTTPRoute(ctx, *route)
		},
	}
}

func (validator KongHTTPValidator) ValidateHTTPRoute(
	ctx context.Context, httproute gatewayapi.HTTPRoute,
) (bool, string, error) {
	// in order to be sure whether or not an HTTPRoute resource is managed by this
	// controller we disallow references to Gateway resources that do not exist.
	var managedGateways []*gatewayapi.Gateway
	for _, parentRef := range httproute.Spec.ParentRefs {
		// determine the namespace of the gateway referenced via parentRef. If no
		// explicit namespace is provided, assume the namespace of the route.
		namespace := httproute.Namespace
		if parentRef.Namespace != nil {
			namespace = string(*parentRef.Namespace)
		}

		// gather the Gateway resource referenced by parentRef and fail validation
		// if there is no such Gateway resource.
		gateway := gatewayapi.Gateway{}
		if err := validator.ManagerClient.Get(ctx, client.ObjectKey{
			Namespace: namespace,
			Name:      string(parentRef.Name),
		}, &gateway); err != nil {
			return false, fmt.Sprintf("couldn't retrieve referenced gateway %s/%s", namespace, parentRef.Name), err
		}

		// pull the referenced GatewayClass object from the Gateway
		gatewayClass := gatewayapi.GatewayClass{}
		if err := validator.ManagerClient.Get(ctx, client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, &gatewayClass); err != nil {
			return false, fmt.Sprintf("couldn't retrieve referenced gatewayclass %s", gateway.Spec.GatewayClassName), err
		}

		// determine ultimately whether the Gateway is managed by this controller implementation
		if gatewayClass.Spec.ControllerName == gatewaycontroller.GetControllerName() {
			managedGateways = append(managedGateways, &gateway)
		}
	}

	// if there are no managed Gateways this is not a supported HTTPRoute
	if len(managedGateways) == 0 {
		return true, "", nil
	}

	// Now that we know whether or not the HTTPRoute is linked to a managed
	// Gateway we can run it through full validation.
	var routeValidator routeValidator = noOpRoutesValidator{}
	if routesSvc, ok := validator.AdminAPIServicesProvider.GetRoutesService(); ok {
		routeValidator = routesSvc
	}
	return gatewayvalidation.ValidateHTTPRoute(
		ctx, routeValidator, validator.ParserFeatures, &httproute, managedGateways...,
	)
}
