package admission

import (
	"context"
	"fmt"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	credsvalidation "github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi/validators/consumer/credentials"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/kongstate"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

// KongValidator validates Kong entities.
type KongValidator interface {
	ValidateConsumer(ctx context.Context, consumer kongv1.KongConsumer) (bool, string, error)
	ValidatePlugin(ctx context.Context, plugin kongv1.KongPlugin) (bool, string, error)
	ValidateClusterPlugin(ctx context.Context, plugin kongv1.KongClusterPlugin) (bool, string, error)
	ValidateCredential(ctx context.Context, secret corev1.Secret) (bool, string, error)
}

// KongHTTPValidator implements KongValidator interface to validate Kong
// entities using the Admin API of Kong.
type KongHTTPValidator struct {
	ConsumerSvc   kong.AbstractConsumerService
	PluginSvc     kong.AbstractPluginService
	Logger        logrus.FieldLogger
	SecretGetter  kongstate.SecretGetter
	ManagerClient client.Client

	ingressClassMatcher func(*metav1.ObjectMeta, annotations.ClassMatching) bool
}

// NewKongHTTPValidator provides a new KongHTTPValidator object provided a
// controller-runtime client which will be used to retrieve reference objects
// such as consumer credentials secrets. If you do not pass a cached client
// here, the performance of this validator can get very poor at high scales.
func NewKongHTTPValidator(
	consumerSvc kong.AbstractConsumerService,
	pluginSvc kong.AbstractPluginService,
	logger logrus.FieldLogger,
	managerClient client.Client,
	ingressClass string,
) KongHTTPValidator {
	matcher := annotations.IngressClassValidatorFuncFromObjectMeta(ingressClass)
	return KongHTTPValidator{
		ConsumerSvc:   consumerSvc,
		PluginSvc:     pluginSvc,
		Logger:        logger,
		SecretGetter:  &managerClientSecretGetter{managerClient: managerClient},
		ManagerClient: managerClient,

		ingressClassMatcher: matcher,
	}
}

// ValidateConsumer checks if consumer has a Username and a consumer with
// the same username doesn't exist in Kong.
// If an error occurs during validation, it is returned as the last argument.
// The first boolean communicates if the consumer is valid or not and string
// holds a message if the entity is not valid.
func (validator KongHTTPValidator) ValidateConsumer(
	ctx context.Context,
	consumer kongv1.KongConsumer,
) (bool, string, error) {
	// ignore consumers that are being managed by another controller
	if !validator.ingressClassMatcher(&consumer.ObjectMeta, annotations.ExactClassMatch) {
		return true, "", nil
	}

	// a consumer without a username is not valid
	if consumer.Username == "" {
		return false, ErrTextConsumerUsernameEmpty, nil
	}

	// verify that the consumer is not already otherwise present in the data-plane
	c, err := validator.ConsumerSvc.Get(ctx, &consumer.Username)
	if err != nil {
		if !kong.IsNotFoundErr(err) {
			validator.Logger.Errorf("failed to fetch consumer from kong: %v", err)
			return false, ErrTextConsumerUnretrievable, err
		}
	}
	if c != nil {
		return false, ErrTextConsumerExists, nil
	}

	// if there are no credentials for this consumer, there's no need to move on
	// to credentials validation.
	if len(consumer.Credentials) == 0 {
		return true, "", nil
	}

	// pull all the managed consumers in order to build a validation index of
	// credentials so that the consumers credentials references can be validated.
	managedConsumers, err := validator.listManagedConsumers(ctx)
	if err != nil {
		return false, ErrTextConsumerUnretrievable, err
	}

	// retrieve the consumer's credentials secrets to validate them with the index
	credentials := make([]*corev1.Secret, 0, len(consumer.Credentials))
	for _, secretName := range consumer.Credentials {
		// retrieve the credentials secret
		secret, err := validator.SecretGetter.GetSecret(consumer.Namespace, secretName)
		if err != nil {
			if errors.IsNotFound(err) {
				return false, ErrTextConsumerCredentialSecretNotFound, err
			}
			return false, ErrTextFailedToRetrieveSecret, err
		}

		// do the basic credentials validation
		if err := credsvalidation.ValidateCredentials(consumer.Name, secret); err != nil {
			return false, ErrTextConsumerCredentialValidationFailed, err
		}

		// if valid, store it so we can index it for upcoming constraints validation
		credentials = append(credentials, secret)
	}

	// unique constraints on consumer credentials are global to all consumers
	// and credentials, so we must build an index based on all existing credentials.
	credentialsIndex, err := globalValidationIndexForCredentials(ctx, validator.ManagerClient, managedConsumers)
	if err != nil {
		return false, ErrTextConsumerCredentialValidationFailed, err
	}

	// validate the consumer's credentials against the index of all managed
	// credentials to ensure they're not in violation of any unique constraints.
	for _, secret := range credentials {
		// do the unique constraints validation of the credentials using the credentials index
		if err := credentialsIndex.ValidateCredentialsForUniqueKeyConstraints(consumer.Name, secret); err != nil {
			return false, ErrTextConsumerCredentialValidationFailed, err
		}
	}

	return true, "", nil
}

// ValidateCredential checks if the secret contains a credential meant to
// be installed in Kong. If so, then it verifies if all the required fields
// are present in it or not. If valid, it returns true with an empty string,
// else it returns false with the error messsage. If an error happens during
// validation, error is returned.
func (validator KongHTTPValidator) ValidateCredential(
	ctx context.Context,
	secret corev1.Secret,
) (bool, string, error) {
	// if the secret doesn't contain a type key it's not a credentials secret
	_, ok := secret.Data[credsvalidation.TypeKey]
	if !ok {
		return true, "", nil
	}

	// credentials are only validated if they are referenced by a managed consumer
	// in the namespace, as such we pull a list of all consumers from the cached
	// client to determine if the credentials are referenced.
	managedConsumers, err := validator.listManagedConsumers(ctx)
	if err != nil {
		return false, ErrTextConsumerUnretrievable, err
	}

	// find the list of managed consumers which actually reference this secret
	// as a credential.
	for _, consumerWhichReferencesThisSecret := range listManagedConsumersReferencingCredentialsSecret(secret, managedConsumers) {
		// perform basic credentials validation first
		if err := credsvalidation.ValidateCredentials(consumerWhichReferencesThisSecret.Name, &secret); err != nil {
			return false, ErrTextConsumerCredentialValidationFailed, err
		}

		// unique constraints on consumer credentials are global to all consumers
		// and credentials, so we must build an index based on all existing credentials.
		credentialsIndex, err := globalValidationIndexForCredentials(ctx, validator.ManagerClient, managedConsumers)
		if err != nil {
			return false, ErrTextConsumerCredentialValidationFailed, err
		}

		// validate this new credential against an index of existing credentials
		secretCopy := secret
		if err := credentialsIndex.ValidateCredentialsForUniqueKeyConstraints(consumerWhichReferencesThisSecret.Name, &secretCopy); err != nil {
			return false, ErrTextConsumerCredentialValidationFailed, err
		}
	}

	return true, "", nil
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
	if len(k8sPlugin.Protocols) > 0 {
		plugin.Protocols = kong.StringSlice(k8sPlugin.Protocols...)
	}
	isValid, msg, err := validator.PluginSvc.Validate(ctx, &plugin)
	if err != nil {
		return false, ErrTextPluginConfigValidationFailed, err
	}
	if !isValid {
		return isValid, fmt.Sprintf(ErrTextPluginConfigViolatesSchema, msg), nil
	}
	return isValid, "", nil
}

// ValidateClusterPlugin transfers relevant fields from a KongClusterPlugin into a KongPlugin and then returns
// the result of ValidatePlugin for the derived KongPlugin
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

// -----------------------------------------------------------------------------
// KongHTTPValidator - Private Methods
// -----------------------------------------------------------------------------

func (validator KongHTTPValidator) listManagedConsumers(ctx context.Context) ([]*kongv1.KongConsumer, error) {
	// gather a list of all consumers from the cached client
	consumers := &kongv1.KongConsumerList{}
	if err := validator.ManagerClient.List(ctx, consumers, &client.ListOptions{
		Namespace: corev1.NamespaceAll,
	}); err != nil {
		return nil, err
	}

	// reduce the consumer set to consumers managed by this controller
	managedConsumers := make([]*kongv1.KongConsumer, 0)
	for _, consumer := range consumers.Items {
		if !validator.ingressClassMatcher(&consumer.ObjectMeta, annotations.ExactClassMatch) {
			// ignore consumers (and subsequently secrets) that are managed by other controllers
			continue
		}
		consumerCopy := consumer
		managedConsumers = append(managedConsumers, &consumerCopy)
	}

	return managedConsumers, nil
}

// -----------------------------------------------------------------------------
// Private - Manager Client Secret Getter
// -----------------------------------------------------------------------------

type managerClientSecretGetter struct {
	managerClient client.Client
}

func (m *managerClientSecretGetter) GetSecret(namespace, name string) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	return secret, m.managerClient.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, secret)
}
