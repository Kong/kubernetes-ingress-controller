package envtest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers/conditions"
)

func TestKongConsumer_ProgrammedCondition(t *testing.T) {
	t.Parallel()

	err := kongv1.AddToScheme(scheme.Scheme)
	require.NoError(t, err)
	envcfg := Setup(t, scheme.Scheme, WithInstallKongCRDs(true))
	ctrlClient := NewControllerClient(t, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	RunManager(ctx, t, envcfg, func(cfg *manager.Config) {
		cfg.UpdateStatus = true
		cfg.PublishStatusAddress = []string{"http://localhost:8080"}
	})

	ns := CreateNamespace(ctx, t, ctrlClient)

	testCases := []struct {
		name    string
		objects []client.Object
		test    func(t *testing.T, ctrlClient client.Client)
	}{
		{
			name: "valid KongConsumer",
			objects: []client.Object{
				&kongv1.KongConsumer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "consumer",
						Namespace: ns.Name,
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Username: "username",
				},
			},
			test: func(t *testing.T, ctrlClient client.Client) {
				require.Eventually(t, func() bool {
					var consumer kongv1.KongConsumer
					err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
						Name:      "consumer",
						Namespace: ns.Name,
					}, &consumer)
					if err != nil {
						t.Logf("error getting consumer: %v", err)
						return false
					}

					if !conditions.Contain(
						consumer.Status.Conditions,
						conditions.WithType("Programmed"),
						conditions.WithStatus(metav1.ConditionTrue),
					) {
						t.Logf("Programmed condition not found, actual: %v", consumer.Status.Conditions)
						return false
					}
					return true
				}, 10*time.Second, 50*time.Millisecond)
			},
		},
		{
			name: "KongConsumer referencing non-existent secret",
			objects: []client.Object{
				&kongv1.KongConsumer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "consumer-with-secret",
						Namespace: ns.Name,
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Username:    "username",
					Credentials: []string{"non-existent-secret"},
				},
			},
			test: func(t *testing.T, ctrlClient client.Client) {
				require.Eventually(t, func() bool {
					var consumer kongv1.KongConsumer
					err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
						Name:      "consumer-with-secret",
						Namespace: ns.Name,
					}, &consumer)
					if err != nil {
						t.Logf("error getting consumer: %v", err)
						return false
					}

					if !conditions.Contain(
						consumer.Status.Conditions,
						conditions.WithType("Programmed"),
						conditions.WithStatus(metav1.ConditionFalse),
						conditions.WithReason(string(kongv1.ReasonInvalid)),
					) {
						t.Logf("Programmed condition not found, actual: %v", consumer.Status.Conditions)
						return false
					}
					return true
				}, 10*time.Second, 50*time.Millisecond)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, obj := range tc.objects {
				require.NoError(t, ctrlClient.Create(ctx, obj))
				t.Cleanup(func() { _ = ctrlClient.Delete(ctx, obj) })
			}
			tc.test(t, ctrlClient)
		})
	}
}
