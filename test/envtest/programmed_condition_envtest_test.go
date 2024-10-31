package envtest

import (
	"context"
	"testing"
	"time"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/conditions"
)

func TestKongCRDs_ProgrammedCondition(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrlClient := NewControllerClient(t, scheme, envcfg)
	ns := CreateNamespace(ctx, t, ctrlClient)
	healthProbePort := helpers.GetFreePort(t)

	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithUpdateStatus(),
		WithHealthProbePort(healthProbePort),
		WithPublishService(ns.Name),
		WithKongServiceFacadeFeatureEnabled(),
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
		{
			name: "valid KongServiceFacade with Ingress",
			objects: []client.Object{
				&incubatorv1alpha1.KongServiceFacade{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-facade",
						Namespace: ns.Name,
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: incubatorv1alpha1.KongServiceFacadeSpec{
						Backend: incubatorv1alpha1.KongServiceFacadeBackend{
							Name: "svc",
							Port: 80,
						},
					},
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc",
						Namespace: ns.Name,
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Port: 80,
							},
						},
					},
				},
				&netv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ingress",
						Namespace: ns.Name,
					},
					Spec: netv1.IngressSpec{
						IngressClassName: lo.ToPtr(annotations.DefaultIngressClass),
						Rules: []netv1.IngressRule{{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{{
										Path:     "/",
										PathType: lo.ToPtr(netv1.PathTypeImplementationSpecific),
										Backend: netv1.IngressBackend{
											Resource: &corev1.TypedLocalObjectReference{
												APIGroup: lo.ToPtr(incubatorv1alpha1.SchemeGroupVersion.Group),
												Kind:     incubatorv1alpha1.KongServiceFacadeKind,
												Name:     "svc-facade",
											},
										},
									}},
								},
							},
						}},
					},
				},
			},
			getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
				var serviceFacade incubatorv1alpha1.KongServiceFacade
				err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
					Name:      "svc-facade",
					Namespace: ns.Name,
				}, &serviceFacade)
				if err != nil {
					return nil, err
				}
				return serviceFacade.Status.Conditions, nil
			},
			expectedProgrammedStatus: metav1.ConditionTrue,
			expectedProgrammedReason: kongv1.ReasonProgrammed,
		},
		{
			name: "KongServiceFacade with Ingress referring to non-existent Service",
			objects: []client.Object{
				&incubatorv1alpha1.KongServiceFacade{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-facade-with-non-existent-service",
						Namespace: ns.Name,
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: incubatorv1alpha1.KongServiceFacadeSpec{
						Backend: incubatorv1alpha1.KongServiceFacadeBackend{
							Name: "non-existent-service",
							Port: 80,
						},
					},
				},
				&netv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ingress",
						Namespace: ns.Name,
					},
					Spec: netv1.IngressSpec{
						IngressClassName: lo.ToPtr(annotations.DefaultIngressClass),
						Rules: []netv1.IngressRule{{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{{
										Path:     "/",
										PathType: lo.ToPtr(netv1.PathTypeImplementationSpecific),
										Backend: netv1.IngressBackend{
											Resource: &corev1.TypedLocalObjectReference{
												APIGroup: lo.ToPtr(incubatorv1alpha1.SchemeGroupVersion.Group),
												Kind:     incubatorv1alpha1.KongServiceFacadeKind,
												Name:     "svc-facade-with-non-existent-service",
											},
										},
									}},
								},
							},
						}},
					},
				},
			},
			getExpectedObjectConditions: func(ctrlClient client.Client) ([]metav1.Condition, error) {
				var serviceFacade incubatorv1alpha1.KongServiceFacade
				err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
					Name:      "svc-facade-with-non-existent-service",
					Namespace: ns.Name,
				}, &serviceFacade)
				if err != nil {
					return nil, err
				}
				return serviceFacade.Status.Conditions, nil
			},
			expectedProgrammedStatus: metav1.ConditionFalse,
			expectedProgrammedReason: kongv1.ReasonInvalid,
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
