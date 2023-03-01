package envtest

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateNamespace creates namespace using the provided client and returns it.
func CreateNamespace(ctx context.Context, t *testing.T, client ctrlclient.Client) corev1.Namespace {
	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Labels: map[string]string{
				"test": "envtest",
			},
		},
	}
	require.NoError(t, client.Create(ctx, &ns, &ctrlclient.CreateOptions{}))
	// No need to remove the namespace since envtest cannot delete namespaces:
	// https://book.kubebuilder.io/reference/envtest.html#namespace-usage-limitation
	return ns
}

// CreatePod creates pod using the provided client and returns it.
func CreatePod(ctx context.Context, t *testing.T, client ctrlclient.Client, ns string) corev1.Pod {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      uuid.NewString(),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "kong",
					Image: "kong",
				},
			},
		},
	}
	require.NoError(t, client.Create(ctx, &pod, &ctrlclient.CreateOptions{}))
	return pod
}
