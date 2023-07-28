package envtest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers/conditions"
)

func TestKongCRDs_ProgrammedCondition(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme, WithInstallKongCRDs(true))
	ctrlClient := NewControllerClient(t, scheme, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	RunManager(ctx, t, envcfg, func(cfg *manager.Config) {
		cfg.UpdateStatus = true
		cfg.PublishStatusAddress = []string{"http://localhost:8080"}
	})

	ns := CreateNamespace(ctx, t, ctrlClient)

	testCases := []struct {
		name                        string
		objects                     []client.Object
		getExpectedObjectConditions func(ctrlClient client.Client) ([]metav1.Condition, error)
		expectedProgrammedStatus    metav1.ConditionStatus
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
			getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
				var consumer kongv1.KongConsumer
				err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
					Name:      "consumer",
					Namespace: ns.Name,
				}, &consumer)
				if err != nil {
					return nil, err
				}
				return consumer.Status.Conditions, nil
			},
			expectedProgrammedStatus: metav1.ConditionTrue,
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
			getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
				var consumer kongv1.KongConsumer
				err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
					Name:      "consumer-with-secret",
					Namespace: ns.Name,
				}, &consumer)
				if err != nil {
					return nil, err
				}
				return consumer.Status.Conditions, nil
			},
			expectedProgrammedStatus: metav1.ConditionFalse,
		},
		{
			name: "valid KongConsumerGroup",
			objects: []client.Object{
				&kongv1beta1.KongConsumerGroup{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "consumer-group",
						Namespace: ns.Name,
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
				},
			},
			getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
				var consumerGroup kongv1beta1.KongConsumerGroup
				err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
					Name:      "consumer-group",
					Namespace: ns.Name,
				}, &consumerGroup)
				if err != nil {
					return nil, err
				}
				return consumerGroup.Status.Conditions, nil
			},
			expectedProgrammedStatus: metav1.ConditionTrue,
		},
		{
			name: "valid KongPlugin",
			objects: []client.Object{
				&kongv1.KongPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "kong-plugin",
						Namespace:   ns.Name,
						Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
					},
					PluginName: "plugin",
				},
				&kongv1.KongConsumer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "consumer-for-plugin",
						Namespace: ns.Name,
						Annotations: map[string]string{
							annotations.IngressClassKey:                           annotations.DefaultIngressClass,
							annotations.AnnotationPrefix + annotations.PluginsKey: "kong-plugin",
						},
					},
					Username: "foo",
				},
			},
			getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
				var plugin kongv1.KongPlugin
				err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
					Name:      "kong-plugin",
					Namespace: ns.Name,
				}, &plugin)
				if err != nil {
					return nil, err
				}
				return plugin.Status.Conditions, nil
			},
			expectedProgrammedStatus: metav1.ConditionTrue,
		},
		{
			name: "valid KongClusterPlugin",
			objects: []client.Object{
				&kongv1.KongClusterPlugin{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "kong-cluster-plugin",
						Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
					},
					PluginName: "plugin",
				},
				&kongv1.KongConsumer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "consumer-for-cluster-plugin",
						Namespace: ns.Name,
						Annotations: map[string]string{
							annotations.IngressClassKey:                           annotations.DefaultIngressClass,
							annotations.AnnotationPrefix + annotations.PluginsKey: "kong-cluster-plugin",
						},
					},
					Username: "foo",
				},
			},
			getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
				var clusterPlugin kongv1.KongClusterPlugin
				err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
					Name: "kong-cluster-plugin",
				}, &clusterPlugin)
				if err != nil {
					return nil, err
				}
				return clusterPlugin.Status.Conditions, nil
			},
			expectedProgrammedStatus: metav1.ConditionTrue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, obj := range tc.objects {
				require.NoError(t, ctrlClient.Create(ctx, obj))
				t.Cleanup(func() { _ = ctrlClient.Delete(ctx, obj) })
			}

			require.Eventually(t, func() bool {
				cs, err := tc.getExpectedObjectConditions(ctrlClient)
				if err != nil {
					t.Logf("error getting expected object conditions: %v", err)
					return false
				}
				if !conditions.Contain(
					cs,
					conditions.WithType(string(kongv1.ConditionProgrammed)),
					conditions.WithStatus(tc.expectedProgrammedStatus),
				) {
					t.Logf("Programmed condition not found, actual: %v", cs)
					return false
				}
				return true
			}, 10*time.Second, 50*time.Millisecond)
		})
	}
}
