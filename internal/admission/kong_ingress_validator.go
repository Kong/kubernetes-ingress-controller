package admission

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

type KongIngressValidator struct{}

func (k KongIngressValidator) ValidateCreate(_ context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return k.validate(obj)
}

func (k KongIngressValidator) ValidateUpdate(_ context.Context, _, newObj runtime.Object) (warnings admission.Warnings, err error) {
	return k.validate(newObj)
}

func (k KongIngressValidator) ValidateDelete(context.Context, runtime.Object) (warnings admission.Warnings, err error) {
	return admission.Warnings{}, nil
}

func (k KongIngressValidator) validate(obj runtime.Object) (warnings admission.Warnings, err error) {
	ingress, ok := obj.(*kongv1.KongIngress)
	if !ok {
		return admission.Warnings{}, fmt.Errorf("unexpected type, expected *kongv1.KongIngress, got %T", obj)
	}

	if ingress.Proxy != nil {
		const warning = "'proxy' ids DEPRECATED. It will have no effect. Use Service's annotations instead."
		warnings = append(warnings, warning)
	}
	if ingress.Route != nil {
		const warning = "'route' is DEPRECATED. It will have no effect. Use Ingress' annotations instead."
		warnings = append(warnings, warning)
	}

	return warnings, nil
}
