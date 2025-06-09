package license_test

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	konnectlicense "github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
	"github.com/stretchr/testify/require"
)

func TestSecretLicenseStore_Store(t *testing.T) {
	testCases := []struct {
		name        string
		secret      *corev1.Secret
		license     license.KonnectLicense
		expectError bool
	}{
		{
			name: "stored secrets",
			secret: &corev1.Secret{

				ObjectMeta: metav1.ObjectMeta{
					Name:      "konnect-license-test-cp",
					Namespace: "default",
				},
			},
			license: license.KonnectLicense{
				Payload:   "some-license-payload",
				UpdatedAt: time.Now(),
				ID:        "some-license-id",
			},
			expectError: false,
		},
		{
			name: "secret not found",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "another-secret",
					Namespace: "default",
				},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithObjects(tc.secret).Build()
			s := konnectlicense.NewSecretLicenseStore(cl, "default", "test-cp")
			err := s.Store(context.Background(), tc.license)
			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			secret := &corev1.Secret{}
			err = cl.Get(context.Background(), client.ObjectKeyFromObject(tc.secret), secret)
			require.NoError(t, err)
			// fake client stores stringData of secret as-is.
			require.Equal(t, tc.license.Payload, secret.StringData["payload"])
		})
	}
}
