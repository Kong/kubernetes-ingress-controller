package admission

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	credsvalidation "github.com/kong/kubernetes-ingress-controller/v2/internal/admission/validation/consumers/credentials"
)

func (validator KongHTTPValidator) Secret() CustomValidatorAdapter {
	// secrets are only validated on update because they must be referenced by a
	// managed consumer in order for us to validate them, and because credentials
	// validation also happens at the consumer side of the reference so a
	// credentials secret can not be referenced without being validated.
	return CustomValidatorAdapter{
		validateUpdate: func(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (bool, string, error) {
			secret, ok := newObj.(*corev1.Secret)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *corev1.Secret, got %T", newObj)
			}
			return validator.ValidateCredential(ctx, *secret)
		},
	}
}

// ValidateCredential checks if the secret contains a credential meant to
// be installed in Kong. If so, then it verifies if all the required fields
// are present in it or not. If valid, it returns true with an empty string,
// else it returns false with the error message. If an error happens during
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

	// verify whether this secret is referenced by any managed consumer
	managedConsumersWithReferences := listManagedConsumersReferencingCredentialsSecret(secret, managedConsumers)
	if len(managedConsumersWithReferences) == 0 {
		// if no managed consumers reference this secret, its considered
		// unmanaged and we don't validate it unless it becomes referenced
		// by a managed consumer at a later time.
		return true, "", nil
	}

	// now that we know at least one managed consumer is referencing this
	// secret we perform the base-level credentials secret validation.
	if err := credsvalidation.ValidateCredentials(&secret); err != nil {
		return false, ErrTextConsumerCredentialValidationFailed, err
	}

	// if base-level validation passes we move on to create an index of
	// all managed credentials so that we can verify that the updates to
	// this secret are not in violation of any unique key constraints.
	ignoreSecrets := map[string]map[string]struct{}{secret.Namespace: {secret.Name: {}}}
	credentialsIndex, err := globalValidationIndexForCredentials(ctx, validator.ManagerClient, managedConsumers, ignoreSecrets)
	if err != nil {
		return false, ErrTextConsumerCredentialValidationFailed, err
	}

	// the index is built, now validate that the newly updated secret
	// is not in violation of any constraints.
	if err := credentialsIndex.ValidateCredentialsForUniqueKeyConstraints(&secret); err != nil {
		return false, ErrTextConsumerCredentialValidationFailed, err
	}

	return true, "", nil
}
