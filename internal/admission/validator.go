package admission

import (
	"github.com/golang/glog"
	"github.com/hbagdi/go-kong/kong"
	configuration "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	"github.com/pkg/errors"
)

// KongValidator validates Kong entities.
type KongValidator interface {
	ValidateConsumer(consumer configuration.KongConsumer) (bool, string, error)
	ValidatePlugin(consumer configuration.KongPlugin) (bool, string, error)
}

// KongHTTPValidator implements KongValidator interface to validate Kong
// entities using the Admin API of Kong.
type KongHTTPValidator struct {
	Client *kong.Client
}

// ValidateConsumer checks if consumer has a Username and a consumer with
// the same username doesn't exist in Kong.
// If an error occurs during validation, it is returned as the last argument.
// The first boolean communicates if the consumer is valid or not and string
// holds a message if the entity is not valid.
func (validator KongHTTPValidator) ValidateConsumer(
	consumer configuration.KongConsumer) (bool, string, error) {
	if consumer.Username == "" {
		return false, "username cannot be empty", nil
	}
	c, err := validator.Client.Consumers.Get(nil, &consumer.Username)
	if err != nil {
		if kong.IsNotFoundErr(err) {
			return true, "", nil
		}
		glog.Errorf("admission controller: "+
			"error getting consumer from Kong: %v", err)
		return false, "", errors.Wrap(err, "fetching consumer from Kong")
	}
	if c != nil {
		return false, "consumer already exists", nil
	}
	return true, "", nil
}

// ValidatePlugin checks if k8sPlugin is valid. It does so by performing
// an HTTP request to Kong's Admin API entity validation endpoints.
// If an error occurs during validation, it is returned as the last argument.
// The first boolean communicates if k8sPluign is valid or not and string
// holds a message if the entity is not valid.
func (validator KongHTTPValidator) ValidatePlugin(
	k8sPlugin configuration.KongPlugin) (bool, string, error) {
	if k8sPlugin.PluginName == "" {
		return false, "plugin name cannot be empty", nil
	}
	var plugin kong.Plugin
	plugin.Name = kong.String(k8sPlugin.PluginName)
	if k8sPlugin.Config != nil {
		plugin.Config = kong.Configuration(k8sPlugin.Config)
	}
	if k8sPlugin.RunOn != "" {
		plugin.RunOn = kong.String(k8sPlugin.RunOn)
	}
	if len(k8sPlugin.Protocols) > 0 {
		plugin.Protocols = kong.StringSlice(k8sPlugin.Protocols...)
	}
	req, err := validator.Client.NewRequest("POST", "/schemas/plugins/validate",
		nil, &plugin)
	if err != nil {
		return false, "", err
	}
	resp, err := validator.Client.Do(nil, req, nil)
	if err != nil {
		return false, err.Error(), nil
	}
	if resp.StatusCode == 201 {
		return true, "", nil
	}
	if err != nil {
		return false, "", err
	}
	return true, "", nil
}

func empty(s *string) bool {
	return s == nil && *s == ""
}
