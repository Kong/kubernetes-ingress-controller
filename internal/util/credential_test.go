package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
)

func TestExtractKongCredentialType(t *testing.T) {
	tests := []struct {
		name     string
		secret   *corev1.Secret
		credType string
		wantErr  error
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
			credType: "key-auth",
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
			credType: "",
			wantErr: fmt.Errorf("Secret %s/%s used as credential, but lacks %s label",
				"default", "secret", labels.CredentialTypeLabel),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credType, err := ExtractKongCredentialType(tt.secret)
			require.Equal(t, tt.credType, credType)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
