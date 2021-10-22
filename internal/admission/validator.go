package admission

import (
	"context"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/kongstate"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

// KongValidator validates Kong entities.
type KongValidator interface {
	ValidateConsumer(ctx context.Context, consumer configurationv1.KongConsumer) (bool, string, error)
	ValidatePlugin(ctx context.Context, plugin configurationv1.KongPlugin) (bool, string, error)
	ValidateClusterPlugin(ctx context.Context, plugin configurationv1.KongClusterPlugin) (bool, string, error)
	ValidateCredential(secret corev1.Secret) (bool, string, error)
}

// KongHTTPValidator implements KongValidator interface to validate Kong
// entities using the Admin API of Kong.
type KongHTTPValidator struct {
	ConsumerSvc  kong.AbstractConsumerService
	PluginSvc    kong.AbstractPluginService
	Logger       logrus.FieldLogger
	SecretGetter kongstate.SecretGetter
}

// ValidateConsumer checks if consumer has a Username and a consumer with
// the same username doesn't exist in Kong.
// If an error occurs during validation, it is returned as the last argument.
// The first boolean communicates if the consumer is valid or not and string
// holds a message if the entity is not valid.
func (validator KongHTTPValidator) ValidateConsumer(ctx context.Context,
	consumer configurationv1.KongConsumer) (bool, string, error) {
	if consumer.Username == "" {
		return false, ErrTextConsumerUsernameEmpty, nil
	}
	c, err := validator.ConsumerSvc.Get(ctx, &consumer.Username)
	if err != nil {
		if kong.IsNotFoundErr(err) {
			return true, "", nil
		}
		validator.Logger.Errorf("failed to fetch consumer from kong: %v", err)
		return false, ErrTextConsumerUnretrievable, err
	}
	if c != nil {
		return false, ErrTextConsumerExists, nil
	}
	return true, "", nil
}

// ValidatePlugin checks if k8sPlugin is valid. It does so by performing
// an HTTP request to Kong's Admin API entity validation endpoints.
// If an error occurs during validation, it is returned as the last argument.
// The first boolean communicates if k8sPluign is valid or not and string
// holds a message if the entity is not valid.
func (validator KongHTTPValidator) ValidatePlugin(ctx context.Context,
	k8sPlugin configurationv1.KongPlugin) (bool, string, error) {
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
		config, err := kongstate.SecretToConfiguration(validator.SecretGetter,
			(*k8sPlugin.ConfigFrom).SecretValue, k8sPlugin.Namespace)
		if err != nil {
			return false, ErrTextPluginSecretConfigUnretrievable, err
		}
		plugin.Config = config
	}
	if k8sPlugin.RunOn != "" {
		plugin.RunOn = kong.String(k8sPlugin.RunOn)
	}
	if len(k8sPlugin.Protocols) > 0 {
		plugin.Protocols = kong.StringSlice(k8sPlugin.Protocols...)
	}
	isValid, err := validator.PluginSvc.Validate(ctx, &plugin)
	if err != nil {
		return false, ErrTextPluginConfigViolatesSchema, err
	}
	return isValid, "", nil
}

// ValidateClusterPlugin transfers relevant fields from a KongClusterPlugin into a KongPlugin and then returns
// the result of ValidatePlugin for the derived KongPlugin
func (validator KongHTTPValidator) ValidateClusterPlugin(ctx context.Context,
	k8sPlugin configurationv1.KongClusterPlugin) (bool, string, error) {
	derived := configurationv1.KongPlugin{
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
		ref := configurationv1.ConfigSource{
			SecretValue: configurationv1.SecretValueFromSource{
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

var (
	keyAuthFields   = []string{"key"}
	basicAuthFields = []string{"username", "password"}
	hmacAuthFields  = []string{"username", "secret"}
	jwtAuthFields   = []string{"algorithm", "rsa_public_key", "key", "secret"}
	mtlsAuthFields  = []string{"subject_name"}

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
		"mtls-auth":            mtlsAuthFields,
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
