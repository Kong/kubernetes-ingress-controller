package admission

import (
	"context"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	credsvalidation "github.com/kong/kubernetes-ingress-controller/v3/internal/admission/validation/consumers/credentials"
)

// -----------------------------------------------------------------------------
// KongHTTPValidator - Private Functions
// -----------------------------------------------------------------------------

// listManagedConsumersReferencingCredentialsSecret takes a Secret and a list of KongConsumers.
// It returns a list of KongConsumers that reference that Secret as a credential.
func listManagedConsumersReferencingCredentialsSecret(secret corev1.Secret, managedConsumers []*kongv1.KongConsumer) []*kongv1.KongConsumer {
	// determine if this credential is being actively referenced by a consumer
	consumersWhichReferenceSecret := make([]*kongv1.KongConsumer, 0)
	for _, consumer := range managedConsumers {
		// verify that the secret is actually in the same namespace (its possible for
		// there to be name duplication across multiple namespaces).
		if consumer.Namespace == secret.Namespace {
			// verify whether the consumer in this same namespace as the secret
			// actually references it as a credential.
			for _, secretName := range consumer.Credentials {
				if secretName == secret.Name { // this credential is referred to from a consumer
					consumersWhichReferenceSecret = append(consumersWhichReferenceSecret, consumer)
				}
			}
		}
	}
	return consumersWhichReferenceSecret
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
