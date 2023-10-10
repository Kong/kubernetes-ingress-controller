package admission

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// CustomValidatorAdapter is an adapter for legacy validators in our codebase that makes them compatible with
// the new controller-runtime's CustomValidator interface.
type CustomValidatorAdapter struct {
	validateCreate func(ctx context.Context, obj runtime.Object) (bool, string, error)
	validateUpdate func(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (bool, string, error)
	validateDelete func(ctx context.Context, obj runtime.Object) (bool, string, error)
}

func (c CustomValidatorAdapter) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	if c.validateCreate == nil {
		return admission.Warnings{}, nil
	}
	return c.returnValues(c.validateCreate(ctx, obj))
}

func (c CustomValidatorAdapter) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (warnings admission.Warnings, err error) {
	if c.validateUpdate == nil {
		return admission.Warnings{}, nil
	}
	return c.returnValues(c.validateUpdate(ctx, oldObj, newObj))
}

func (c CustomValidatorAdapter) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	if c.validateDelete == nil {
		return admission.Warnings{}, nil
	}
	return c.returnValues(c.validateDelete(ctx, obj))
}

func (c CustomValidatorAdapter) returnValues(ok bool, message string, err error) (admission.Warnings, error) {
	if err != nil {
		return admission.Warnings{message}, err
	}
	if !ok {
		return admission.Warnings{message}, errors.New(message)
	}
	if message != "" {
		return admission.Warnings{message}, nil
	}
	return admission.Warnings{}, nil
}
