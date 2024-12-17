package credentials

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
)

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
						labels.CredentialTypeLabel: "key-auth",
					},
				},
				Data: map[string][]byte{
					"key": []byte("little-rabbits-be-good"),
				},
			},
			wantErr: nil,
		},
		{
			name: "valid jwt credential with HS512",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				Data: map[string][]byte{
					"algorithm": []byte("HS512"),
					"key":       []byte("key-name"),
					"secret":    []byte("secret-name"),
				},
			},
			wantErr: nil,
		},
		{
			name: "valid jwt credential with HS384",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				Data: map[string][]byte{
					"algorithm": []byte("HS384"),
					"key":       []byte("key-name"),
					"secret":    []byte("secret-name"),
				},
			},
			wantErr: nil,
		},
		{
			name: "valid jwt credential with HS256",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				Data: map[string][]byte{
					"algorithm": []byte("HS256"),
					"key":       []byte("key-name"),
					"secret":    []byte("secret-name"),
				},
			},
			wantErr: nil,
		},
		{
			name: "valid jwt credential with RS256",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				Data: map[string][]byte{
					"algorithm": []byte("RS256"),
				},
			},
			wantErr: fmt.Errorf("missing required field(s): rsa_public_key, key"),
		},
		{
			name: "valid jwt credential with RS256",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				Data: map[string][]byte{
					"algorithm":      []byte("RS256"),
					"key":            []byte(""),
					"rsa_public_key": []byte(""),
				},
			},
			wantErr: fmt.Errorf("some fields were invalid due to missing data: rsa_public_key, key"),
		},
		{
			name: "invalid credential type",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "bee-auth",
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
			wantErr: fmt.Errorf("secret has no credential type, add a %s label", labels.CredentialTypeLabel),
		},
		{
			name: "missing required field",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "key-auth",
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
						labels.CredentialTypeLabel: "key-auth",
					},
				},
				Data: map[string][]byte{
					"key": []byte(""),
				},
			},
			wantErr: fmt.Errorf("some fields were invalid due to missing data: key"),
		},
		{
			name: "invalid jwt credential with HS256 missing secret",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				Data: map[string][]byte{
					"algorithm": []byte("HS256"),
					"key":       []byte("key-name"),
				},
			},
			wantErr: fmt.Errorf("missing required field(s): secret"),
		},
		{
			name: "invalid jwt credential with HS384 missing secret",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				Data: map[string][]byte{
					"algorithm": []byte("HS384"),
					"key":       []byte("key-name"),
				},
			},
			wantErr: fmt.Errorf("missing required field(s): secret"),
		},
		{
			name: "invalid jwt credential with HS512 missing secret",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				Data: map[string][]byte{
					"algorithm": []byte("HS512"),
					"key":       []byte("key-name"),
				},
			},
			wantErr: fmt.Errorf("missing required field(s): secret"),
		},
		{
			name: "valid jwt credential with RS256 missing secret",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				Data: map[string][]byte{
					"algorithm":      []byte("RS256"),
					"key":            []byte("key-name"),
					"rsa_public_key": []byte("-----BEGIN PUBLIC KEY----- AXAXAXAAXA... -----END PUBLIC KEY-----"),
				},
			},
		},
		{
			name: "valid jwt credential with RS512 missing secret",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				Data: map[string][]byte{
					"algorithm":      []byte("RS512"),
					"key":            []byte("key-name"),
					"rsa_public_key": []byte("-----BEGIN PUBLIC KEY----- AXAXAXAAXA... -----END PUBLIC KEY-----"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCredentials(tt.secret)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
