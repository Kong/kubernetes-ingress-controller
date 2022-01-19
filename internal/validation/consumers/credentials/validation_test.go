package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniqueConstraintsValidation(t *testing.T) {
	t.Log("setting up an index of existing credentials which have unique constraints")
	index := make(Index)
	require.NoError(t, index.add(Credential{
		Key:   "username",
		Value: "batman",
		Type:  "basic-auth",
	}))
	require.NoError(t, index.add(Credential{
		Key:   "username",
		Value: "robin",
		Type:  "basic-auth",
	}))

	t.Log("verifying that a new basic-auth credential with a unique username doesn't violate constraints")
	nonviolatingCredential := Credential{
		Key:   "username",
		Value: "nightwing",
		Type:  "basic-auth",
	}
	assert.NoError(t, index.add(nonviolatingCredential))

	t.Log("verifying that a new basic-auth credential with a username that's already in use violates constraints")
	violatingCredential := Credential{
		Key:   "username",
		Value: "batman",
		Type:  "basic-auth",
	}
	assert.True(t, IsKeyUniqueConstrained(violatingCredential.Type, violatingCredential.Key))
	err := index.add(violatingCredential)
	assert.Error(t, err)

	t.Log("setting up a list of existing credentials which have no unique constraints")
	index = make(Index)
	assert.NoError(t, index.add(Credential{
		Key:   "key",
		Value: "test",
		Type:  "acl",
	}))

	t.Log("verifying that non-unique constrained credentials don't trigger a violation")
	duplicate := Credential{
		Key:   "key",
		Value: "test",
		Type:  "acl",
	}
	assert.False(t, IsKeyUniqueConstrained(duplicate.Type, duplicate.Key))
	assert.NoError(t, index.add(duplicate))

	t.Log("verifying that unconstrained keys for types with constraints don't flag as violated")
	assert.False(t, IsKeyUniqueConstrained("basic-auth", "unconstrained-key"))
}
