package kongadmin

import (
	"context"
	"fmt"

	"github.com/kong/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
)

// KongAdmin holds the configuration that should be posted to KongAdmin
type KongAdmin struct {
	EventHooks []kong.EventHooks
}

// CollectEventHooks collect event hooks configuration from kongIngress
func CollectEventHooks(kongIngress *configurationv1.KongIngress, kongAdminPost *KongAdmin) error {
	if kongIngress.EventHooks.Handler == nil || len(*kongIngress.EventHooks.Handler) == 0 {
		return fmt.Errorf("handler could not be empty")
	}

	if kongIngress.EventHooks.Source == nil || len(*kongIngress.EventHooks.Source) == 0 {
		return fmt.Errorf("source could not be empty")
	}

	if kongIngress.EventHooks.Config.URL != nil || len(*kongIngress.EventHooks.Config.URL) == 0 {
		return fmt.Errorf("config url should be set")
	}

	config := &kong.Config{}
	config.URL = kongIngress.EventHooks.Config.URL
	if kongIngress.EventHooks.Config.Method != nil && len(*kongIngress.EventHooks.Config.Method) > 0 {
		config.Method = kongIngress.EventHooks.Config.Method
	}
	if kongIngress.EventHooks.Config.Functions != nil && len(kongIngress.EventHooks.Config.Functions) > 0 {
		config.Functions = append(config.Functions, kongIngress.EventHooks.Config.Functions...)
	}
	aEventHooks := kong.EventHooks{}
	aEventHooks.Config = config
	*aEventHooks.Handler = *kongIngress.EventHooks.Handler
	*aEventHooks.Source = *kongIngress.EventHooks.Source

	if kongIngress.EventHooks.Event != nil && len(*kongIngress.EventHooks.Event) > 0 {
		*aEventHooks.Event = *kongIngress.EventHooks.Event
	}

	kongAdminPost.EventHooks = append(kongAdminPost.EventHooks, aEventHooks)
	return nil
}

// PostEventHooks post eventhooks to Kong Admin API
func PostEventHooks(client *kong.Client, kongAdminPost *KongAdmin) error {
	ctx := context.TODO()

	_, err := client.EventHooks.AddWebhook(ctx, &kongAdminPost.EventHooks)
	if err != nil {
		return fmt.Errorf("failed posting event hooks %v", err)
	}

	return nil
}
