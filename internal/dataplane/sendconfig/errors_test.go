package sendconfig_test

import (
	"errors"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
)

type testError struct{}

func (t testError) Error() string {
	return "test error"
}

func TestUpdateError(t *testing.T) {
	someResourceFailure := lo.Must(failures.NewResourceFailure("some reason", &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "name",
			Namespace: "namespace",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
	}))

	updateErr := sendconfig.NewUpdateError([]failures.ResourceFailure{someResourceFailure}, testError{})
	require.Equal(t, "test error", updateErr.Error())
	require.Len(t, updateErr.ResourceFailures(), 1)
	unwraps := errors.As(updateErr, &testError{})
	require.True(t, unwraps, "UpdateError should unwrap to inner error")
}
