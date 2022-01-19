// Package credentials includes validators for the credentials provided for KongConsumers.
package credentials

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

// -----------------------------------------------------------------------------
//  Validation - Public Functions
// -----------------------------------------------------------------------------

// ValidateCredentials performs basic validation on a credential secret given
// the Kubernetes secret which contains credentials data.
func ValidateCredentials(consumerName string, secret *corev1.Secret) error {
	// the indication of credential type is required to be present on all credentials.
	credentialTypeB, ok := secret.Data[TypeKey]
	if !ok {
		return fmt.Errorf("missing required key %s", TypeKey)
	}
	credentialType := string(credentialTypeB)

	// verify that the credential type provided is valid
	if !SupportedTypes.Has(credentialType) {
		return fmt.Errorf("invalid credential type %s", secret.Data[TypeKey])
	}

	// it's not valid to have a secret that ONLY has a type
	if len(secret.Data) == 1 {
		return fmt.Errorf("invalid credentials secret, no data present")
	}

	// verify that all required fields are present
	var missingFields []string
	var missingDataFields []string
	for _, field := range CredTypeToFields[credentialType] {
		// verify whether the required field is missing
		requiredData, ok := secret.Data[field]
		if !ok {
			missingFields = append(missingFields, field)
			continue
		}

		// verify whether the required field is present, but missing data
		if len(requiredData) < 1 {
			missingDataFields = append(missingDataFields, field)
		}
	}

	// report on any required fields that were missing
	if len(missingFields) > 0 {
		return fmt.Errorf("missing required field(s): %s", strings.Join(missingFields, ", "))
	}

	// report on any required fields that were present, but were missing actual data
	if len(missingDataFields) > 0 {
		return fmt.Errorf("some fields were invalid due to missing data: %s", strings.Join(missingDataFields, ", "))
	}

	return nil
}

// IsKeyUniqueConstrained indicates whether or not a given key and its type there
// are unique constraints in place.
func IsKeyUniqueConstrained(keyType, key string) (constrained bool) {
	constrainedKeys, credTypeHasConstraints := uniqueKeyConstraints[keyType]
	if !credTypeHasConstraints {
		return
	}

	for _, constrainedKey := range constrainedKeys {
		if key == constrainedKey {
			constrained = true
			return
		}
	}

	return
}

// -----------------------------------------------------------------------------
//  Validation - Credentials
// -----------------------------------------------------------------------------

// Credential is a metadata struct to help validate the contents of
// consumer credentials, particularly unique constraints on the underlying data.
type Credential struct {
	// ConsumerName indicates the name of the KongConsumer which this credential
	// is supplied for.
	ConsumerName string

	// ConsumerNamespace indicates the namespace that the KongConsumer which this
	// credential is supplied for.
	ConsumerNamespace string

	// Type indicates the credential type, which will reference one of the types
	// in the SupportedTypes set.
	Type string

	// Key is the key for the credentials data
	Key string

	// Value is the data provided for the key
	Value string
}

// -----------------------------------------------------------------------------
// Validation - Validating Index
// -----------------------------------------------------------------------------

// Index is a map of credentials types to a map of credential keys to the underlying
// values already seen for that type and key. This type is used as a history tracker
// for validation so that callers can keep track of the credentials they've seen thus
// far and validate whether new credentials they encounter are in violation of any
// constraints on their respective types.
type Index map[string]map[string]map[string]struct{}

// ValidateCredentialsForUniqueKeyConstraints will attempt to add a new Credential to the CredentialsTypeMap
// and will validate it for both normal structure validation and for
// unique key constraint violations.
func (cs Index) ValidateCredentialsForUniqueKeyConstraints(consumerName string, secret *corev1.Secret) error {
	// the indication of credential type is required to be present on all credentials.
	credentialTypeB, ok := secret.Data[TypeKey]
	if !ok {
		return fmt.Errorf("missing required key %s", TypeKey)
	}
	credentialType := string(credentialTypeB)

	// the additional key/values are optional, but must be validated
	// for unique constraint violations. Using an index of credentials
	// validation will be checked on any Add() to the index, so errors
	// from this include the unique key constraint errors.
	for k, v := range secret.Data {
		if err := cs.add(Credential{
			ConsumerName:      consumerName,
			ConsumerNamespace: secret.Namespace,
			Type:              credentialType,
			Key:               k,
			Value:             string(v),
		}); err != nil {
			return err
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Valdating Index - Private Methods
// -----------------------------------------------------------------------------

func (cs Index) add(newCred Credential) error {
	// retrieve all the keys which are constrained for this type
	constraints, ok := uniqueKeyConstraints[newCred.Type]
	if !ok {
		return nil // there are no constraints for this credType
	}

	// for each key which is constrained for this type check the existing list
	// to see if there are any violations of that constraint given the new credentials
	for _, constrainedKey := range constraints {
		if newCred.Key == constrainedKey { // this key has constraints on it, we need to check for violations
			if _, ok := cs[newCred.Type][newCred.Key][newCred.Value]; ok {
				return fmt.Errorf("unique key constraint violated for %s", newCred.Key)
			}
		}
	}

	// if needed, initialize the index
	if cs[newCred.Type] == nil {
		cs[newCred.Type] = map[string]map[string]struct{}{newCred.Key: {newCred.Value: {}}}
	}
	if cs[newCred.Type][newCred.Key] == nil {
		cs[newCred.Type][newCred.Key] = make(map[string]struct{})
	}

	// if we make it here there's been no constraint violation, add it to the index
	cs[newCred.Type][newCred.Key][newCred.Value] = struct{}{}

	return nil
}
