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
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers/conditions"
)

func TestKongCRDs_ProgrammedCondition(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrlClient := NewControllerClient(t, scheme, envcfg)
	ns := CreateNamespace(ctx, t, ctrlClient)

	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithUpdateStatus(),
		WithPublishService(ns.Name),
		WithPublishStatusAddress("http://localhost:8080"),
	)

	testCases := []struct {
		name                        string
		objects                     []client.Object
		getExpectedObjectConditions func(ctrlClient client.Client) ([]metav1.Condition, error)
		expectedProgrammedStatus    metav1.ConditionStatus
		expectedProgrammedReason    kongv1.ConditionReason
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
			expectedProgrammedReason: kongv1.ReasonProgrammed,
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
			expectedProgrammedReason: kongv1.ReasonInvalid,
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
			expectedProgrammedReason: kongv1.ReasonProgrammed,
		},
		// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/4578
		// if there are multiple KIC instances within a cluster, they will fight over setting this condition because the
		// controllers do not filter on ingress class. we need to limit them to only resources referenced from others,
		// similar to Secrets, to use this
		// {
		//	name: "valid KongPlugin",
		//	objects: []client.Object{
		//		&kongv1.KongPlugin{
		//			ObjectMeta: metav1.ObjectMeta{
		//				Name:        "kong-plugin",
		//				Namespace:   ns.Name,
		//				Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
		//			},
		//			PluginName: "plugin",
		//		},
		//		&kongv1.KongConsumer{
		//			ObjectMeta: metav1.ObjectMeta{
		//				Name:      "consumer-for-plugin",
		//				Namespace: ns.Name,
		//				Annotations: map[string]string{
		//					annotations.IngressClassKey:                           annotations.DefaultIngressClass,
		//					annotations.AnnotationPrefix + annotations.PluginsKey: "kong-plugin",
		//				},
		//			},
		//			Username: "foo",
		//		},
		//	},
		//	getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
		//		var plugin kongv1.KongPlugin
		//		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
		//			Name:      "kong-plugin",
		//			Namespace: ns.Name,
		//		}, &plugin)
		//		if err != nil {
		//			return nil, err
		//		}
		//		return plugin.Status.Conditions, nil
		//	},
		//	expectedProgrammedStatus: metav1.ConditionTrue,
		//	expectedProgrammedReason: kongv1.ReasonProgrammed,
		// },
		// {
		//	name: "invalid KongPlugin",
		//	objects: []client.Object{
		//		&kongv1.KongPlugin{
		//			ObjectMeta: metav1.ObjectMeta{
		//				Name:        "invalid-kong-plugin",
		//				Namespace:   ns.Name,
		//				Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
		//			},
		//			PluginName: "plugin",
		//			// Specifying both Config and ConfigFrom is invalid.
		//			Config: apiextensionsv1.JSON{Raw: []byte(`{"key": "value"}`)},
		//			ConfigFrom: &kongv1.ConfigSource{
		//				SecretValue: kongv1.SecretValueFromSource{
		//					Secret: "secret",
		//					Key:    "key",
		//				},
		//			},
		//		},
		//		&kongv1.KongConsumer{
		//			ObjectMeta: metav1.ObjectMeta{
		//				Name:      "consumer-for-invalid-plugin",
		//				Namespace: ns.Name,
		//				Annotations: map[string]string{
		//					annotations.IngressClassKey:                           annotations.DefaultIngressClass,
		//					annotations.AnnotationPrefix + annotations.PluginsKey: "invalid-kong-plugin",
		//				},
		//			},
		//			Username: "foo",
		//		},
		//	},
		//	getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
		//		var plugin kongv1.KongPlugin
		//		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
		//			Name:      "invalid-kong-plugin",
		//			Namespace: ns.Name,
		//		}, &plugin)
		//		if err != nil {
		//			return nil, err
		//		}
		//		return plugin.Status.Conditions, nil
		//	},
		//	expectedProgrammedStatus: metav1.ConditionFalse,
		//	expectedProgrammedReason: kongv1.ReasonInvalid,
		// },
		// {
		//	name: "valid KongClusterPlugin",
		//	objects: []client.Object{
		//		&kongv1.KongClusterPlugin{
		//			ObjectMeta: metav1.ObjectMeta{
		//				Name:        "kong-cluster-plugin",
		//				Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
		//			},
		//			PluginName: "plugin",
		//		},
		//		&kongv1.KongConsumer{
		//			ObjectMeta: metav1.ObjectMeta{
		//				Name:      "consumer-for-cluster-plugin",
		//				Namespace: ns.Name,
		//				Annotations: map[string]string{
		//					annotations.IngressClassKey:                           annotations.DefaultIngressClass,
		//					annotations.AnnotationPrefix + annotations.PluginsKey: "kong-cluster-plugin",
		//				},
		//			},
		//			Username: "foo",
		//		},
		//	},
		//	getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
		//		var clusterPlugin kongv1.KongClusterPlugin
		//		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
		//			Name: "kong-cluster-plugin",
		//		}, &clusterPlugin)
		//		if err != nil {
		//			return nil, err
		//		}
		//		return clusterPlugin.Status.Conditions, nil
		//	},
		//	expectedProgrammedStatus: metav1.ConditionTrue,
		//	expectedProgrammedReason: kongv1.ReasonProgrammed,
		// },
		// {
		//	name: "invalid KongClusterPlugin",
		//	objects: []client.Object{
		//		&kongv1.KongPlugin{
		//			ObjectMeta: metav1.ObjectMeta{
		//				Name:        "invalid-kong-cluster-plugin",
		//				Namespace:   ns.Name,
		//				Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
		//			},
		//			PluginName: "plugin",
		//			// Specifying both Config and ConfigFrom is invalid.
		//			Config: apiextensionsv1.JSON{Raw: []byte(`{"key": "value"}`)},
		//			ConfigFrom: &kongv1.ConfigSource{
		//				SecretValue: kongv1.SecretValueFromSource{
		//					Secret: "secret",
		//					Key:    "key",
		//				},
		//			},
		//		},
		//		&kongv1.KongConsumer{
		//			ObjectMeta: metav1.ObjectMeta{
		//				Name:      "consumer-for-invalid-cluster-plugin",
		//				Namespace: ns.Name,
		//				Annotations: map[string]string{
		//					annotations.IngressClassKey:                           annotations.DefaultIngressClass,
		//					annotations.AnnotationPrefix + annotations.PluginsKey: "invalid-kong-cluster-plugin",
		//				},
		//			},
		//			Username: "foo",
		//		},
		//	},
		//	getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
		//		var plugin kongv1.KongPlugin
		//		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
		//			Name:      "invalid-kong-cluster-plugin",
		//			Namespace: ns.Name,
		//		}, &plugin)
		//		if err != nil {
		//			return nil, err
		//		}
		//		return plugin.Status.Conditions, nil
		//	},
		//	expectedProgrammedStatus: metav1.ConditionFalse,
		//	expectedProgrammedReason: kongv1.ReasonInvalid,
		// },
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
					conditions.WithReason(string(tc.expectedProgrammedReason)),
					conditions.WithStatus(tc.expectedProgrammedStatus),
				) {
					t.Logf("Programmed condition not found, actual: %v", cs)
					return false
				}
				return true
			}, test.RequestTimeout, 50*time.Millisecond)
		})
	}
}
