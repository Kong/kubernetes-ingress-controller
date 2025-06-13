package license_test

import (
	"encoding/base64"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	konnectlicense "github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/license"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
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
			err := s.Store(t.Context(), tc.license)
			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			secret := &corev1.Secret{}
			err = cl.Get(t.Context(), client.ObjectKeyFromObject(tc.secret), secret)
			require.NoError(t, err)
			// fake client stores stringData of secret as-is.
			require.Equal(t, tc.license.Payload, secret.StringData["payload"])
		})
	}
}

func TestSecretLicenseStore_Load(t *testing.T) {
	timeNowUnix := time.Now().Unix()
	testCases := []struct {
		name        string
		secret      *corev1.Secret
		license     license.KonnectLicense
		expectError bool
	}{
		{
			name: "load license successfully",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "konnect-license-test-cp",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"payload":    []byte(base64.StdEncoding.EncodeToString([]byte("some-license-payload"))),
					"id":         []byte(base64.StdEncoding.EncodeToString([]byte("some-license-id"))),
					"updated_at": []byte(base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(timeNowUnix, 10)))),
				},
			},
			license: license.KonnectLicense{
				Payload:   "some-license-payload",
				ID:        "some-license-id",
				UpdatedAt: time.Unix(timeNowUnix, 0),
			},
		},
		{
			name: "secret not found",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "some-other-secret",
					Namespace: "default",
				},
			},
			expectError: true,
		},
		{
			name: "missing payload",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "konnect-license-test-cp",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"id":         []byte(base64.StdEncoding.EncodeToString([]byte("some-license-id"))),
					"updated_at": []byte(base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(timeNowUnix, 10)))),
				},
			},
			expectError: true,
		},
		{
			name: "cannot parse update_at",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "konnect-license-test-cp",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"payload":    []byte(base64.StdEncoding.EncodeToString([]byte("some-license-payload"))),
					"id":         []byte(base64.StdEncoding.EncodeToString([]byte("some-license-id"))),
					"updated_at": []byte(base64.StdEncoding.EncodeToString([]byte("not-a-timestamp"))),
				},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithObjects(tc.secret).Build()
			s := konnectlicense.NewSecretLicenseStore(cl, "default", "test-cp")

			l, err := s.Load(t.Context())
			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.license, l)
		})
	}
}
