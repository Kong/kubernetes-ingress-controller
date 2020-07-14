package admission

import (
	"context"
	"fmt"
	"strings"

	"github.com/hbagdi/go-kong/kong"
	configuration "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// KongValidator validates Kong entities.
type KongValidator interface {
	ValidateConsumer(consumer configuration.KongConsumer) (bool, string, error)
	ValidatePlugin(consumer configuration.KongPlugin) (bool, string, error)
	ValidateCredential(secret corev1.Secret) (bool, string, error)
}

// KongHTTPValidator implements KongValidator interface to validate Kong
// entities using the Admin API of Kong.
type KongHTTPValidator struct {
	Client *kong.Client
	Logger logrus.FieldLogger
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
	c, err := validator.Client.Consumers.Get(context.TODO(), &consumer.Username)
	if err != nil {
		if kong.IsNotFoundErr(err) {
			return true, "", nil
		}
		validator.Logger.Errorf("failed to fetch consumer from kong: %v", err)
		return false, "", fmt.Errorf("fetching consumer from Kong: %w", err)
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
	resp, err := validator.Client.Do(context.Background(), req, nil)
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

var (
	keyAuthFields   = []string{"key"}
	basicAuthFields = []string{"username", "password"}
	hmacAuthFields  = []string{"username", "secret"}
	jwtAuthFields   = []string{"algorithm", "rsa_public_key", "key", "secret"}

	// TODO dynamically fetch these from Kong
	credTypeToFields = map[string][]string{
		"key-auth":             keyAuthFields,
		"keyauth_credential":   keyAuthFields,
		"basic-auth":           basicAuthFields,
		"basicauth_credential": basicAuthFields,
		"hmac-auth":            hmacAuthFields,
		"hmacauth_credential":  hmacAuthFields,
		"jwt":                  jwtAuthFields,
		"jwt_secret":           jwtAuthFields,
		"oauth2":               {"name", "client_id", "client_secret", "redirect_uris"},
		"acl":                  {"group"},
	}
)

// ValidateCredential checks if the secret contains a credential meant to
// be installed in Kong. If so, then it verifies if all the required fields
// are present in it or not. If valid, it returns true with an empty string,
// else it returns false with the error messsage. If an error happens during
// validation, error is returned.
func (validator KongHTTPValidator) ValidateCredential(
	secret corev1.Secret) (bool, string, error) {

	credTypeBytes, ok := secret.Data["kongCredType"]
	if !ok {
		// doesn't look like a credential resource
		return true, "", nil
	}
	credType := string(credTypeBytes)

	fields, ok := credTypeToFields[credType]
	if !ok {
		return false, "invalid credential type: " + credType, nil
	}

	var missingFields []string
	for _, field := range fields {
		if _, ok := secret.Data[field]; !ok {
			missingFields = append(missingFields, field)
		}
	}
	if len(missingFields) != 0 {
		return false, "missing required field(s): " +
			strings.Join(missingFields, ", "), nil
	}

	// TODO add unique key violation detection
	// For each credential, there is a unique column, like key for key-auth,
	// username for basic-auth; make an API call to Kong's Admin API
	// and verify if there will be a violation, similar to how it's done
	// for KongConsumer; return error if the resource is already present in
	// Kong.
	return true, "", nil
}
