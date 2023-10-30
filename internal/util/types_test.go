package util

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/scheme"
)

func TestPopulateTypeMeta(t *testing.T) {
	credential := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "corn",
		},
		StringData: map[string]string{
			"kongCredType": "basic-auth",
			"username":     "corn",
			"password":     "corn",
		},
	}

	require.Empty(t, credential.GetObjectKind().GroupVersionKind().Kind)

	err := PopulateTypeMeta(credential, lo.Must(scheme.Get()))

	require.NoError(t, err)
	require.NotEmpty(t, credential.GetObjectKind().GroupVersionKind().Kind)
}
