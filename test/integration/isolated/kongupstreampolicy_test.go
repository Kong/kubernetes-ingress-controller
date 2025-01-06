//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"
	"github.com/kong/kubernetes-configuration/pkg/clientset"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestKongUpstreamPolicyStatus(t *testing.T) {
	f := features.
		New("essentials").
		WithLabel(testlabels.Kind, testlabels.KindKongUpstreamPolicy).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
		)).
		WithSetup("prepare clients", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cluster := GetClusterFromCtx(ctx)

			kongClients, err := clientset.NewForConfig(cluster.Config())
			require.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, kongClients)

			gatewayClient, err := gatewayclient.NewForConfig(cluster.Config())
			require.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, gatewayClient)

			return ctx
		}).
		WithSetup("deploy required resources", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
			cluster := GetClusterFromCtx(ctx)
			namespace := GetNamespaceForT(ctx, t)
			ingressClass := GetIngressClassFromCtx(ctx)
			clients := GetFromCtxForT[*clientset.Clientset](ctx, t)
			gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)

			t.Log("creating KongUpstreamPolicies")
			upstreamPolicies := []*kongv1beta1.KongUpstreamPolicy{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "upstream-policy-1",
					},
					Spec: kongv1beta1.KongUpstreamPolicySpec{
						Algorithm: lo.ToPtr("round-robin"),
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "upstream-policy-2",
					},
					Spec: kongv1beta1.KongUpstreamPolicySpec{
						Algorithm: lo.ToPtr("consistent-hashing"),
					},
				},
			}
			for _, upstreamPolicy := range upstreamPolicies {
				_, err := clients.ConfigurationV1beta1().KongUpstreamPolicies(namespace).Create(ctx, upstreamPolicy, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(upstreamPolicy)
			}

			t.Log("creating Services")
			container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
			deployment := generators.NewDeploymentForContainer(container)
			deployment, err := cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
			require.NoError(t, err)
			cleaner.Add(deployment)

			service1 := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
			service1.Name = "service-1"
			service1.Annotations = map[string]string{
				kongv1beta1.KongUpstreamPolicyAnnotationKey: "upstream-policy-1",
			}

			service2 := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
			service2.Name = "service-2"
			service2.Annotations = map[string]string{
				kongv1beta1.KongUpstreamPolicyAnnotationKey: "upstream-policy-2",
			}

			services := []*corev1.Service{service1, service2}
			for _, service := range services {
				_, err := cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
				require.NoError(t, err)
				cleaner.Add(service)
			}

			serviceFacade := &incubatorv1alpha1.KongServiceFacade{
				ObjectMeta: metav1.ObjectMeta{
					Name: "service-facade",
					Annotations: map[string]string{
						kongv1beta1.KongUpstreamPolicyAnnotationKey: "upstream-policy-1",
						annotations.IngressClassKey:                 ingressClass,
					},
				},
				Spec: incubatorv1alpha1.KongServiceFacadeSpec{
					Backend: incubatorv1alpha1.KongServiceFacadeBackend{
						Name: "service-1",
						Port: 80,
					},
				},
			}
			_, err = clients.IncubatorV1alpha1().KongServiceFacades(namespace).Create(ctx, serviceFacade, metav1.CreateOptions{})
			require.NoError(t, err)
			cleaner.Add(serviceFacade)

			t.Log("creating Ingress")
			ingress := &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ingress",
				},
				Spec: netv1.IngressSpec{
					IngressClassName: &ingressClass,
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path:     "/s1",
											PathType: lo.ToPtr(netv1.PathTypePrefix),
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "service-1",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
										{
											Path:     "/s2",
											PathType: lo.ToPtr(netv1.PathTypePrefix),
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "service-2",
													Port: netv1.ServiceBackendPort{
														Number: 80,
													},
												},
											},
										},
										{
											Path:     "/sf",
											PathType: lo.ToPtr(netv1.PathTypePrefix),
											Backend: netv1.IngressBackend{
												Resource: &corev1.TypedLocalObjectReference{
													APIGroup: lo.ToPtr(incubatorv1alpha1.SchemeGroupVersion.Group),
													Kind:     incubatorv1alpha1.KongServiceFacadeKind,
													Name:     "service-facade",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}
			_, err = cluster.Client().NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})
			require.NoError(t, err)
			cleaner.Add(ingress)

			t.Log("creating IngressClass")
			gatewayClassName := uuid.NewString()
			gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
			require.NoError(t, err)
			cleaner.Add(gwc)

			t.Log("creating Gateway")
			gw, err := helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayClassName)
			require.NoError(t, err)
			cleaner.Add(gw)

			t.Log("creating HTTPRoute")
			httpRoute := &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: "http-route",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Rules: []gatewayapi.HTTPRouteRule{
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								{
									BackendRef: gatewayapi.BackendRef{
										BackendObjectReference: gatewayapi.BackendObjectReference{
											Kind: lo.ToPtr(gatewayapi.Kind("Service")),
											Name: "service-1",
											Port: lo.ToPtr(gatewayapi.PortNumber(80)),
										},
									},
								},
							},
						},
					},
				},
			}
			_, err = gatewayClient.GatewayV1().HTTPRoutes(namespace).Create(ctx, httpRoute, metav1.CreateOptions{})
			require.NoError(t, err)
			cleaner.Add(httpRoute)

			return ctx
		}).
		Assess("all ancestors are Accepted and Programmed", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			namespace := GetNamespaceForT(ctx, t)
			clients := GetFromCtxForT[*clientset.Clientset](ctx, t)

			t.Log("checking conditions expecting Accepted and Programmed to be True for all ancestors")
			require.Eventually(t, func() bool {
				upstreamPolicy, err := clients.ConfigurationV1beta1().KongUpstreamPolicies(namespace).Get(ctx, "upstream-policy-1", metav1.GetOptions{})
				require.NoError(t, err)
				err = requireAncestorWithCondition(upstreamPolicy.Status.Ancestors, "service-1", "Accepted", metav1.ConditionTrue)
				if err != nil {
					t.Logf("no matching ancestor condition found: %s", err)
					return false
				}
				err = requireAncestorWithCondition(upstreamPolicy.Status.Ancestors, "service-1", "Programmed", metav1.ConditionTrue)
				if err != nil {
					t.Logf("no matching ancestor condition found: %s", err)
					return false
				}
				err = requireAncestorWithCondition(upstreamPolicy.Status.Ancestors, "service-facade", "Accepted", metav1.ConditionTrue)
				if err != nil {
					t.Logf("no matching ancestor condition found: %s", err)
					return false
				}
				err = requireAncestorWithCondition(upstreamPolicy.Status.Ancestors, "service-facade", "Programmed", metav1.ConditionTrue)
				if err != nil {
					t.Logf("no matching ancestor condition found: %s", err)
					return false
				}
				return true
			}, time.Minute, time.Second)

			require.Eventually(t, func() bool {
				upstreamPolicy, err := clients.ConfigurationV1beta1().KongUpstreamPolicies(namespace).Get(ctx, "upstream-policy-2", metav1.GetOptions{})
				require.NoError(t, err)
				err = requireAncestorWithCondition(upstreamPolicy.Status.Ancestors, "service-2", "Accepted", metav1.ConditionTrue)
				if err != nil {
					t.Logf("no matching ancestor condition found: %s", err)
					return false
				}
				err = requireAncestorWithCondition(upstreamPolicy.Status.Ancestors, "service-2", "Programmed", metav1.ConditionTrue)
				if err != nil {
					t.Logf("no matching ancestor condition found: %s", err)
					return false
				}
				return true
			}, time.Minute, time.Second)

			return ctx
		}).
		Assess("when HTTPRoute rule uses Services with different KongUpstreamPolicies, the Services are not Accepted and Programmed", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			namespace := GetNamespaceForT(ctx, t)
			gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
			clients := GetFromCtxForT[*clientset.Clientset](ctx, t)

			httpRoute, err := gatewayClient.GatewayV1().HTTPRoutes(namespace).Get(ctx, "http-route", metav1.GetOptions{})
			require.NoError(t, err)

			t.Log("updating HTTPRoute to use Services with different KongUpstreamPolicies in a single rule")
			httpRoute.Spec.Rules[0].BackendRefs = []gatewayapi.HTTPBackendRef{
				{
					BackendRef: gatewayapi.BackendRef{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Kind: lo.ToPtr(gatewayapi.Kind("Service")),
							Name: "service-1",
							Port: lo.ToPtr(gatewayapi.PortNumber(80)),
						},
					},
				},
				{
					BackendRef: gatewayapi.BackendRef{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Kind: lo.ToPtr(gatewayapi.Kind("Service")),
							Name: "service-2",
							Port: lo.ToPtr(gatewayapi.PortNumber(80)),
						},
					},
				},
			}
			_, err = gatewayClient.GatewayV1().HTTPRoutes(namespace).Update(ctx, httpRoute, metav1.UpdateOptions{})
			require.NoError(t, err)

			t.Log("ensuring conflicted Services are not Accepted and Programmed")
			require.Eventually(t, func() bool {
				upstreamPolicy, err := clients.ConfigurationV1beta1().KongUpstreamPolicies(namespace).Get(ctx, "upstream-policy-1", metav1.GetOptions{})
				require.NoError(t, err)
				err = requireAncestorWithCondition(upstreamPolicy.Status.Ancestors, "service-1", "Accepted", metav1.ConditionFalse)
				if err != nil {
					t.Logf("no matching ancestor condition found: %s", err)
					return false
				}
				err = requireAncestorWithCondition(upstreamPolicy.Status.Ancestors, "service-1", "Programmed", metav1.ConditionFalse)
				if err != nil {
					t.Logf("no matching ancestor condition found: %s", err)
					return false
				}
				return true
			}, time.Minute, time.Second)
			require.Eventually(t, func() bool {
				upstreamPolicy, err := clients.ConfigurationV1beta1().KongUpstreamPolicies(namespace).Get(ctx, "upstream-policy-2", metav1.GetOptions{})
				require.NoError(t, err)
				err = requireAncestorWithCondition(upstreamPolicy.Status.Ancestors, "service-2", "Accepted", metav1.ConditionFalse)
				if err != nil {
					t.Logf("no matching ancestor condition found: %s", err)
					return false
				}
				err = requireAncestorWithCondition(upstreamPolicy.Status.Ancestors, "service-2", "Programmed", metav1.ConditionFalse)
				if err != nil {
					t.Logf("no matching ancestor condition found: %s", err)
					return false
				}
				return true
			}, time.Minute, time.Second)

			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}

func requireAncestorWithCondition(
	ancestors []gatewayapi.PolicyAncestorStatus,
	ancestorName string,
	conditionType string,
	expectedStatus metav1.ConditionStatus,
) error {
	ancestor, ok := lo.Find(ancestors, func(ancestor gatewayapi.PolicyAncestorStatus) bool {
		return string(ancestor.AncestorRef.Name) == ancestorName
	})
	if !ok {
		return fmt.Errorf("ancestor named %q not found", ancestorName)
	}

	condition, ok := lo.Find(ancestor.Conditions, func(condition metav1.Condition) bool {
		return condition.Type == conditionType
	})
	if !ok {
		return fmt.Errorf("ancestor named %q does not have condition %q", ancestorName, conditionType)
	}
	if condition.Status != expectedStatus {
		return fmt.Errorf("ancestor named %q has condition %q with status %q, expected %q", ancestorName, conditionType, condition.Status, expectedStatus)
	}

	return nil
}
