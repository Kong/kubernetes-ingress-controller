package admission

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

func (validator KongHTTPValidator) KongClusterPlugin() CustomValidatorAdapter {
	return CustomValidatorAdapter{
		validateCreate: func(ctx context.Context, obj runtime.Object) (bool, string, error) {
			plugin, ok := obj.(*kongv1.KongClusterPlugin)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *kongv1.KongClusterPlugin, got %T", obj)
			}
			return validator.ValidateClusterPlugin(ctx, *plugin)
		},
		validateUpdate: func(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (bool, string, error) {
			newPlugin, ok := newObj.(*kongv1.KongClusterPlugin)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *kongv1.KongClusterPlugin, got %T", newObj)
			}
			return validator.ValidateClusterPlugin(ctx, *newPlugin)
		},
	}
}

// ValidateClusterPlugin transfers relevant fields from a KongClusterPlugin into a KongPlugin and then returns
// the result of ValidatePlugin for the derived KongPlugin.
func (validator KongHTTPValidator) ValidateClusterPlugin(
	ctx context.Context,
	k8sPlugin kongv1.KongClusterPlugin,
) (bool, string, error) {
	derived := kongv1.KongPlugin{
		TypeMeta:    k8sPlugin.TypeMeta,
		ObjectMeta:  k8sPlugin.ObjectMeta,
		ConsumerRef: k8sPlugin.ConsumerRef,
		Disabled:    k8sPlugin.Disabled,
		Config:      k8sPlugin.Config,
		PluginName:  k8sPlugin.PluginName,
		RunOn:       k8sPlugin.RunOn,
		Protocols:   k8sPlugin.Protocols,
	}
	if k8sPlugin.ConfigFrom != nil {
		ref := kongv1.ConfigSource{
			SecretValue: kongv1.SecretValueFromSource{
				Secret: k8sPlugin.ConfigFrom.SecretValue.Secret,
				Key:    k8sPlugin.ConfigFrom.SecretValue.Key,
			},
		}
		derived.ConfigFrom = &ref
		derived.ObjectMeta.Namespace = k8sPlugin.ConfigFrom.SecretValue.Namespace
	} else {
		derived.ObjectMeta.Namespace = "default"
	}
	return validator.ValidatePlugin(ctx, derived)
}
