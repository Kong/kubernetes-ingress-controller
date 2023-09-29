package admission

import (
	"context"
	"fmt"

	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	credsvalidation "github.com/kong/kubernetes-ingress-controller/v2/internal/admission/validation/consumers/credentials"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

func (validator KongHTTPValidator) KongConsumer() CustomValidatorAdapter {
	return CustomValidatorAdapter{
		validateCreate: func(ctx context.Context, obj runtime.Object) (bool, string, error) {
			consumer, ok := obj.(*kongv1.KongConsumer)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *kongv1.KongConsumer, got %T", obj)
			}
			return validator.ValidateConsumer(ctx, *consumer)
		},
		validateUpdate: func(ctx context.Context, oldObj, newObj runtime.Object) (bool, string, error) {
			oldConsumer, ok := oldObj.(*kongv1.KongConsumer)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *kongv1.KongConsumer, got %T", oldObj)
			}
			newConsumer, ok := newObj.(*kongv1.KongConsumer)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *kongv1.KongConsumer, got %T", newObj)
			}
			// Only validate if the username has changed.
			if oldConsumer.Username != newConsumer.Username {
				return false, "", nil
			}

			return validator.ValidateConsumer(ctx, *newConsumer)
		},
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
	if !validator.ingressClassMatcher(&consumer.ObjectMeta, annotations.IngressClassKey, annotations.ExactClassMatch) {
		return true, "", nil
	}

	// a consumer without a username is not valid
	if consumer.Username == "" {
		return false, ErrTextConsumerUsernameEmpty, nil
	}

	errText, err := validator.ensureConsumerDoesNotExistInGateway(ctx, consumer.Username)
	if err != nil || errText != "" {
		return false, errText, err
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
	ignoredSecrets := make(map[string]map[string]struct{})
	for _, secretName := range consumer.Credentials {
		// retrieve the credentials secret
		secret, err := validator.SecretGetter.GetSecret(consumer.Namespace, secretName)
		if err != nil {
			if apierrors.IsNotFound(err) {
				return false, ErrTextConsumerCredentialSecretNotFound, err
			}
			return false, ErrTextFailedToRetrieveSecret, err
		}

		// do the basic credentials validation
		if err := credsvalidation.ValidateCredentials(secret); err != nil {
			return false, ErrTextConsumerCredentialValidationFailed, err
		}

		// if valid, store it so we can index it for upcoming constraints validation
		credentials = append(credentials, secret)

		// later we'll build a global index of all credentials which is needed to
		// validate unique key constraints. That index should omit the secrets that
		// are referenced by this consumer to avoid duplication.
		if _, ok := ignoredSecrets[consumer.Namespace]; !ok {
			ignoredSecrets[consumer.Namespace] = make(map[string]struct{}, len(consumer.Credentials))
		}
		ignoredSecrets[consumer.Namespace][secretName] = struct{}{}
	}

	// unique constraints on consumer credentials are global to all consumers
	// and credentials, so we must build an index based on all existing credentials.
	// we ignore the secrets referenced by this consumer so that the index is not
	// testing them against themselves.
	credentialsIndex, err := globalValidationIndexForCredentials(ctx, validator.ManagerClient, managedConsumers, ignoredSecrets)
	if err != nil {
		return false, ErrTextConsumerCredentialValidationFailed, err
	}

	// validate the consumer's credentials against the index of all managed
	// credentials to ensure they're not in violation of any unique constraints.
	for _, secret := range credentials {
		// do the unique constraints validation of the credentials using the credentials index
		if err := credentialsIndex.ValidateCredentialsForUniqueKeyConstraints(secret); err != nil {
			return false, ErrTextConsumerCredentialValidationFailed, err
		}
	}

	return true, "", nil
}

func (validator KongHTTPValidator) ensureConsumerDoesNotExistInGateway(ctx context.Context, username string) (string, error) {
	if consumerSvc, hasClient := validator.AdminAPIServicesProvider.GetConsumersService(); hasClient {
		// verify that the consumer is not already present in the data-plane
		c, err := consumerSvc.Get(ctx, &username)
		if err != nil {
			if !kong.IsNotFoundErr(err) {
				validator.Logger.WithError(err).Error("failed to fetch consumer from kong")
				return ErrTextConsumerUnretrievable, err
			}
		}
		if c != nil {
			return ErrTextConsumerExists, nil
		}
	}

	// if there's no client, do not verify existence with data-plane as there's none available
	return "", nil
}

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
		consumer := consumer
		if !validator.ingressClassMatcher(&consumer.ObjectMeta, annotations.IngressClassKey,
			annotations.ExactClassMatch) {
			// ignore consumers (and subsequently secrets) that are managed by other controllers
			continue
		}
		consumerCopy := consumer
		managedConsumers = append(managedConsumers, &consumerCopy)
	}

	return managedConsumers, nil
}

// globalValidationIndexForCredentials builds an index of all consumer credentials
// using a given controller-runtime client. This provides an index based on
// ALL namespaces in the cluster. This can be very expensive with high numbers
// of consumer credentials, particularly if the client you provide is not cached
//
// if the caller is building the index to validate updates for specific secrets
// and those secrets should be excluded from the index because they will be added
// later, a map of the namespace and name of those secrets can be provided to exclude them.
func globalValidationIndexForCredentials(ctx context.Context, managerClient client.Client, consumers []*kongv1.KongConsumer, ignoredSecrets map[string]map[string]struct{}) (credsvalidation.Index, error) {
	// pull the reference secrets for credentials from each consumer in the list
	index := make(credsvalidation.Index)
	for _, consumer := range consumers {
		for _, secretName := range consumer.Credentials {
			// if its been requested that this secret be specifically ignored
			// (e.g. that secret is being updated and will soon have new values)
			// then don't add it to the index.
			if secrets, namespaceContainsSkippedSecrets := ignoredSecrets[consumer.Namespace]; namespaceContainsSkippedSecrets {
				if _, secretShouldBeSkipped := secrets[secretName]; secretShouldBeSkipped {
					continue
				}
			}

			// grab a copy of the credential secret
			secret := &corev1.Secret{}
			if err := managerClient.Get(ctx, client.ObjectKey{
				Namespace: consumer.Namespace,
				Name:      secretName,
			}, secret); err != nil {
				if apierrors.IsNotFound(err) { // ignore missing secrets
					continue
				}
				return nil, err
			}

			// add the credential secret to the index
			if err := index.ValidateCredentialsForUniqueKeyConstraints(secret); err != nil {
				return nil, err
			}
		}
	}

	return index, nil
}
