package credentials_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi/validators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi/validators/consumer/credentials"
)

func TestUniqueConstraintsValidation(t *testing.T) {
	t.Log("setting up an index of existing credentials which have unique constraints")
	index := make(credentials.Index)
	require.NoError(t, index.Add(credentials.Credential{
		Key:   "username",
		Value: "batman",
		Type:  "basic-auth",
	}))
	require.NoError(t, index.Add(credentials.Credential{
		Key:   "username",
		Value: "robin",
		Type:  "basic-auth",
	}))

	t.Log("verifying that a new basic-auth credential with a unique username doesn't violate constraints")
	nonviolatingCredential := credentials.Credential{
		Key:   "username",
		Value: "nightwing",
		Type:  "basic-auth",
	}
	assert.NoError(t, index.Add(nonviolatingCredential))

	t.Log("verifying that a new basic-auth credential with a username that's already in use violates constraints")
	violatingCredential := credentials.Credential{
		Key:   "username",
		Value: "batman",
		Type:  "basic-auth",
	}
	assert.True(t, credentials.IsKeyUniqueConstrained(violatingCredential.Type, violatingCredential.Key))
	err := index.Add(violatingCredential)
	assert.Error(t, err)
	assert.IsType(t, validators.UniqueConstraintViolationError{}, err)

	t.Log("setting up a list of existing credentials which have no unique constraints")
	index = make(credentials.Index)
	assert.NoError(t, index.Add(credentials.Credential{
		Key:   "key",
		Value: "test",
		Type:  "acl",
	}))

	t.Log("verifying that non-unique constrained credentials don't trigger a violation")
	duplicate := credentials.Credential{
		Key:   "key",
		Value: "test",
		Type:  "acl",
	}
	assert.False(t, credentials.IsKeyUniqueConstrained(duplicate.Type, duplicate.Key))
	assert.NoError(t, index.Add(duplicate))

	t.Log("verifying that unconstrained keys for types with constraints don't flag as violated")
	assert.False(t, credentials.IsKeyUniqueConstrained("basic-auth", "unconstrained-key"))
}
