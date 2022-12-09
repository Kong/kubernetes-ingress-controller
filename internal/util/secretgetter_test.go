package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestSecretGetter(t *testing.T) {
	// setup kubernetes related configurations
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "valid-secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"config": []byte("carp"),
		},
	}
	fakeK8sClient := fake.NewClientBuilder().WithObjects(secret).Build()
	secretGetter := &SecretGetterFromK8s{
		Reader: fakeK8sClient,
	}

	t.Log("should return the secret object if it is found")

	res, err := secretGetter.GetSecret("default", "valid-secret")

	assert.NoError(t, err)
	assert.Equal(t, secret.Data, res.Data)
	assert.Equal(t, secret.ObjectMeta.GetName(), "valid-secret")

	t.Log("should return error if it is not found")

	_, err = secretGetter.GetSecret("default", "valid-secret-old")
	assert.Error(t, err)
}
