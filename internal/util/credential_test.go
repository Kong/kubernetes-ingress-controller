package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
)

func TestExtractKongCredentialType(t *testing.T) {
	tests := []struct {
		name           string
		secret         *corev1.Secret
		credType       string
		credTypeSource CredentialTypeSource
	}{
		{
			name: "labeled credential",
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
			credType:       "key-auth",
			credTypeSource: CredentialTypeFromLabel,
		},
		{
			// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/4853 to be removed after deprecation window
			name: "kongCredType field credential",
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
			credType:       "key-auth",
			credTypeSource: CredentialTypeFromField,
		},
		{
			// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/4853 to be removed after deprecation window
			name: "kongCredType field and labeled credential, label takes precedence",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "key-auth",
					},
				},
				Data: map[string][]byte{
					"key":          []byte("little-rabbits-be-good"),
					"kongCredType": []byte("bee-auth"),
				},
			},
			credType:       "key-auth",
			credTypeSource: CredentialTypeFromLabel,
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
			credType:       "",
			credTypeSource: CredentialTypeAbsent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credType, credTypeSource := ExtractKongCredentialType(tt.secret)
			require.Equal(t, tt.credType, credType)
			require.Equal(t, tt.credTypeSource, credTypeSource)
		})
	}
}
