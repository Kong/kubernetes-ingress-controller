package admission

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
)

func (validator KongHTTPValidator) Gateway() CustomValidatorAdapter {
	return CustomValidatorAdapter{
		validateCreate: func(ctx context.Context, obj runtime.Object) (bool, string, error) {
			gateway, ok := obj.(*gatewayapi.Gateway)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *gatewayapi.Gateway, got %T", obj)
			}
			return validator.ValidateGateway(ctx, *gateway)
		},
		validateUpdate: func(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (bool, string, error) {
			gateway, ok := newObj.(*gatewayapi.Gateway)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *gatewayapi.Gateway, got %T", newObj)
			}
			return validator.ValidateGateway(ctx, *gateway)
		},
	}
}

func (validator KongHTTPValidator) ValidateGateway(
	ctx context.Context, gateway gatewayapi.Gateway,
) (bool, string, error) {
	// check if the gateway declares a gateway class
	if gateway.Spec.GatewayClassName == "" {
		return true, "", nil
	}

	// validate the gatewayclass reference
	gwc := gatewayapi.GatewayClass{}
	if err := validator.ManagerClient.Get(ctx, client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, &gwc); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return true, "", nil // not managed by this controller
		}
		return false, ErrTextCantRetrieveGatewayClass, err
	}

	// validate whether the gatewayclass is a supported class, if not
	// then this gateway belongs to another controller.
	if gwc.Spec.ControllerName != gatewaycontroller.GetControllerName() {
		return true, "", nil
	}

	return true, "", nil
}
