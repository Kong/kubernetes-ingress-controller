package admission

import (
	"context"
	"fmt"

	"github.com/kong/go-kong/kong"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

func (validator KongHTTPValidator) KongPlugin() CustomValidatorAdapter {
	return CustomValidatorAdapter{
		validateCreate: func(ctx context.Context, obj runtime.Object) (bool, string, error) {
			plugin, ok := obj.(*kongv1.KongPlugin)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *kongv1.KongPlugin, got %T", obj)
			}
			return validator.ValidatePlugin(ctx, *plugin)
		},
		validateUpdate: func(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (bool, string, error) {
			newPlugin, ok := newObj.(*kongv1.KongPlugin)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *kongv1.KongPlugin, got %T", newObj)
			}
			return validator.ValidatePlugin(ctx, *newPlugin)
		},
	}
}

// ValidatePlugin checks if k8sPlugin is valid. It does so by performing
// an HTTP request to Kong's Admin API entity validation endpoints.
// If an error occurs during validation, it is returned as the last argument.
// The first boolean communicates if k8sPluign is valid or not and string
// holds a message if the entity is not valid.
func (validator KongHTTPValidator) ValidatePlugin(
	ctx context.Context,
	k8sPlugin kongv1.KongPlugin,
) (bool, string, error) {
	if k8sPlugin.PluginName == "" {
		return false, ErrTextPluginNameEmpty, nil
	}
	var plugin kong.Plugin
	plugin.Name = kong.String(k8sPlugin.PluginName)
	var err error
	plugin.Config, err = kongstate.RawConfigToConfiguration(k8sPlugin.Config)
	if err != nil {
		return false, ErrTextPluginConfigInvalid, err
	}
	if k8sPlugin.ConfigFrom != nil {
		if len(plugin.Config) > 0 {
			return false, ErrTextPluginUsesBothConfigTypes, nil
		}
		config, err := kongstate.SecretToConfiguration(validator.SecretGetter, (*k8sPlugin.ConfigFrom).SecretValue, k8sPlugin.Namespace)
		if err != nil {
			return false, ErrTextPluginSecretConfigUnretrievable, err
		}
		plugin.Config = config
	}
	if k8sPlugin.RunOn != "" {
		plugin.RunOn = kong.String(k8sPlugin.RunOn)
	}
	if k8sPlugin.Ordering != nil {
		plugin.Ordering = k8sPlugin.Ordering
	}
	if len(k8sPlugin.Protocols) > 0 {
		plugin.Protocols = kong.StringSlice(kongv1.KongProtocolsToStrings(k8sPlugin.Protocols)...)
	}
	errText, err := validator.validatePluginAgainstGatewaySchema(ctx, plugin)
	if err != nil || errText != "" {
		return false, errText, err
	}

	return true, "", nil
}

func (validator KongHTTPValidator) validatePluginAgainstGatewaySchema(ctx context.Context, plugin kong.Plugin) (string, error) {
	pluginService, hasClient := validator.AdminAPIServicesProvider.GetPluginsService()
	if hasClient {
		isValid, msg, err := pluginService.Validate(ctx, &plugin)
		if err != nil {
			return ErrTextPluginConfigValidationFailed, err
		}
		if !isValid {
			return fmt.Sprintf(ErrTextPluginConfigViolatesSchema, msg), nil
		}
	}

	// if there's no client, do not verify with data-plane as there's none available
	return "", nil
}
