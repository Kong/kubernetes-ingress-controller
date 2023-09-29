package admission

import (
	"context"
	"fmt"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ingressvalidation "github.com/kong/kubernetes-ingress-controller/v2/internal/admission/validation/ingress"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

type routeValidator interface {
	Validate(context.Context, *kong.Route) (bool, string, error)
}

type noOpRoutesValidator struct{}

func (noOpRoutesValidator) Validate(_ context.Context, _ *kong.Route) (bool, string, error) {
	return true, "", nil
}

func (validator KongHTTPValidator) Ingress() CustomValidatorAdapter {
	return CustomValidatorAdapter{
		validateCreate: func(ctx context.Context, obj runtime.Object) (bool, string, error) {
			ingress, ok := obj.(*netv1.Ingress)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *netv1.Ingress, got %T", obj)
			}
			return validator.ValidateIngress(ctx, *ingress)
		},
		validateUpdate: func(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (bool, string, error) {
			ingress, ok := newObj.(*netv1.Ingress)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *netv1.Ingress, got %T", newObj)
			}
			return validator.ValidateIngress(ctx, *ingress)
		},
	}
}

func (validator KongHTTPValidator) ValidateIngress(
	ctx context.Context, ingress netv1.Ingress,
) (bool, string, error) {
	// Ignore Ingresses that are being managed by another controller.
	if !validator.ingressClassMatcher(&ingress.ObjectMeta, annotations.IngressClassKey, annotations.ExactClassMatch) &&
		!validator.ingressV1ClassMatcher(&ingress, annotations.ExactClassMatch) {
		return true, "", nil
	}

	var routeValidator routeValidator = noOpRoutesValidator{}
	if routesSvc, ok := validator.AdminAPIServicesProvider.GetRoutesService(); ok {
		routeValidator = routesSvc
	}
	return ingressvalidation.ValidateIngress(ctx, routeValidator, validator.ParserFeatures, validator.KongVersion, &ingress)
}
