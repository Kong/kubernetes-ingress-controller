package credentials

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/labels"
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

func TestValidateCredentials(t *testing.T) {
	tests := []struct {
		name    string
		secret  *corev1.Secret
		wantErr error
	}{
		{
			name: "valid credential",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.LabelPrefix + labels.CredentialKey: "key-auth",
					},
				},
				Data: map[string][]byte{
					"key": []byte("little-rabbits-be-good"),
				},
			},
			wantErr: nil,
		},
		{
			// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/4853 to be removed after deprecation window
			name: "valid credential with deprectated field",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"key":          []byte("little-rabbits-be-good"),
					"kongCredType": []byte("key-auth"),
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid credential type",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.LabelPrefix + labels.CredentialKey: "bee-auth",
					},
				},
				Data: map[string][]byte{
					"key": []byte("little-rabbits-be-good"),
				},
			},
			wantErr: fmt.Errorf("invalid credential type bee-auth"),
		},
		{
			name: "no credential type",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"key": []byte("little-rabbits-be-good"),
				},
			},
			wantErr: fmt.Errorf("secret has no credential type, add a %s label", labels.LabelPrefix+labels.CredentialKey),
		},
		{
			name: "missing required field",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.LabelPrefix + labels.CredentialKey: "key-auth",
					},
				},
				Data: map[string][]byte{
					"bee": []byte("little-rabbits-be-good"),
				},
			},
			wantErr: fmt.Errorf("missing required field(s): key"),
		},
		{
			name: "empty required field",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.LabelPrefix + labels.CredentialKey: "key-auth",
					},
				},
				Data: map[string][]byte{
					"key": []byte(""),
				},
			},
			wantErr: fmt.Errorf("some fields were invalid due to missing data: key"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCredentials(tt.secret)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
