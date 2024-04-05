package configuration

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
)

func TestCoreV1SecretReconciler_shouldReconcileSecret(t *testing.T) {
	tests := []struct {
		name   string
		secret *corev1.Secret
		want   bool
	}{
		{
			name: "Secret with no labels",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{},
			},
			want: false,
		},
		{
			name: "Secret with konghq.com/ca-cert:true label",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"konghq.com/ca-cert": "true",
					},
				},
			},
			want: true,
		},
		{
			name: "Secret with konghq.com/ca-cert:false label",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"konghq.com/ca-cert": "false",
					},
				},
			},
			want: false,
		},
		{
			name: "Secret without konghq.com/ca-cert label",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"some-other-label": "true",
					},
				},
			},
			want: false,
		},
		{
			name: "Secret with labels.CredentialTypeLabel label",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
			},
			want: true,
		},
	}

	r := &CoreV1SecretReconciler{
		ReferenceIndexers: ctrlref.NewCacheIndexers(logr.Discard()),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.shouldReconcileSecret(tt.secret)
			require.Equal(t, tt.want, got)
		})
	}
}
